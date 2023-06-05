package main

import (
	"dv/base"
	"dv/config"
	"dv/m3u8"
	"dv/video"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Cell struct {
	Name string
	Link string

	ext         string // 该文件类型
	IsM3U8child bool   // 是否为m3u8的视频任务
}

func NewCell(name, link string) Cell {
	index := strings.LastIndex(link, ".")
	ext := link[index+1:]

	return Cell{
		Name: name,
		Link: link,
		ext:  ext,
	}
}

func (c *Cell) Do() error {
	switch c.ext {
	case "m3u8":
		return c.m3u8()
	default:
		return c.video()
	}
}

func (c *Cell) m3u8() error {
	beginTime := time.Now()
	// 解析m3u8，得到所有分片地址
	m3u8Downloader := m3u8.NewDownloader(c.Name, c.Link)
	p, err := m3u8Downloader.ExtractContain()
	if err != nil {
		return err
	}
	segments := p.Segments

	// 创建存放所有分片的文件夹
	saveDir := filepath.Join(config.GetConfig().SaveDir, c.Name)
	info, err := os.Stat(saveDir)
	if err != nil || !info.IsDir() {
		_ = os.MkdirAll(saveDir, os.ModeDir)
	}

	core := NewCore()
	core.SetGroupCount(len(segments))

	// 设置并发大小
	size := int(config.GetConfig().ConcurrencyM3u8)
	if len(segments)/(size*size) > size {
		core.vacancy = make(chan struct{}, len(segments)/(size*size))
	} else {
		core.vacancy = make(chan struct{}, size)
	}

	// 输出进度消息
	stop := make(chan struct{})
	defer close(stop)
	go base.NewTicker(stop, func() {
		if c.IsM3U8child {
			// m3u8组任务打印信息
			core.Doing(c.Name, fmt.Sprintf("分片下载进度(%d/%d) %.2f ",
				core.doneCount, core.groupCount, float64(core.doneCount)*100/float64(core.groupCount))+"%")
		}
	})

	var playbackDuration float32 // 该视频总时间
	for index, segment := range segments {
		playbackDuration += segment.Duration

		fn := fmt.Sprintf("%s_part_%d", c.Name, index)
		var link string
		if video.CompleteURL(segment.URI) {
			link = m3u8Downloader.M3u8BaseLink + segment.URI
		} else {
			link = segment.URI
		}
		// 构建每个分片的cell，执行
		t := NewCell(fn, link)
		t.IsM3U8child = true
		core.Submit(&t)
	}
	log.Printf("该电影时长 %s \n", m3u8.CalculationTime(playbackDuration))

	core.Wait()
	// 合并所有分片
	if err := m3u8Downloader.MergeFiles(saveDir); err != nil {
		return err
	}
	_ = os.RemoveAll(saveDir) // 删除文件夹

	log.Printf("%s 任务完成,耗时 %s\n", c.Name, time.Now().Sub(beginTime))

	return nil
}

func (c *Cell) video() error {
	var dir string
	if c.IsM3U8child {
		dir = filepath.Join(config.GetConfig().SaveDir, c.Name)
	} else {
		dir = config.GetConfig().SaveDir
	}

	savePath := filepath.Join(dir, fmt.Sprintf("%s.%s", c.Name))
	d := video.NewDownloader(c.Name, c.Link, savePath)
	if err := d.Execute(c.IsM3U8child); err != nil {
		return err
	}
	return nil
}

func ParseTaskList() ([]Cell, error) {
	cells := make([]Cell, 0)
	bs, err := os.ReadFile(config.GetConfig().TaskList)
	if err != nil {
		return nil, err
	}

	content := string(bs)
	if content == "" {
		return nil, errors.New("content is 0")
	}
	reHead, _ := regexp.Compile(`\s+`)
	content = reHead.ReplaceAllString(content, "\n")
	content = strings.TrimPrefix(content, "\n")
	content = strings.TrimSuffix(content, "\n")
	list := strings.Split(content, "\n")
	for i := 0; i < len(list); i++ {
		if i+1 == len(list) {
			break
		}
		key := list[i]
		value := list[i+1]
		if len(value) < 4 || "http" != value[:4] {
			log.Println("错误值：", value)
			continue
		}
		cells = append(cells, NewCell(key, value))

		i++
	}

	return cells, nil
}
