package main

import (
	"dv/base"
	"dv/config"
	"dv/m3u8"
	"dv/table"
	"dv/video"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Task struct {
	Name string
	Link string

	ext         string // 该文件类型
	IsM3U8child bool   // 是否为m3u8的视频任务
	M3U8Dir     string // 下载的 m3u8 文件目录
}

func NewTask(name, link string) Task {
	u, _ := url.Parse(link)
	index := strings.LastIndex(u.Path, ".")
	ext := u.Path[index+1:]

	return Task{
		Name: name,
		Link: link,
		ext:  ext,
	}
}

func (t *Task) Do() error {
	switch t.ext {
	case "m3u8":
		return t.m3u8()
	default:
		return t.video()
	}
}

func (t *Task) m3u8() error {
	beginTime := time.Now()
	// 解析m3u8，得到所有分片地址
	m3u8Downloader := m3u8.NewDownloader(t.Name, t.Link)
	p, err := m3u8Downloader.ExtractContain()
	if err != nil {
		return err
	}
	segments := p.Segments

	// 创建存放所有分片的文件夹
	saveDir := filepath.Join(config.GetConfig().SaveDir, t.Name)
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
	go func() {
		ticker := time.NewTicker(time.Second * 3)
		defer ticker.Stop()
		var lastNowRS uint
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				// m3u8组任务打印信息
				nv, _ := table.M3U8DownloadSpeed.Get(t.Name)
				dataByTime := float64(nv-lastNowRS) / float64(3)
				speed, unit := base.UnitReturn(dataByTime)
				log.Println(fmt.Sprintf("%s 分片下载进度(%d/%d) 速度：%.2f %s/s 完成度：%.2f ",
					t.Name,
					core.doneCount, core.groupCount,
					speed, unit,
					float64(core.doneCount)*100/float64(core.groupCount),
				) + "%")
				lastNowRS = nv
			}
		}
	}()

	var playbackDuration float32 // 该视频总时间
	for index, segment := range segments {
		playbackDuration += segment.Duration

		fn := fmt.Sprintf("%s_part_%d", t.Name, index)
		link, err := url.Parse(m3u8Downloader.Link)
		if err != nil {
			return err
		}
		link, err = link.Parse(segment.URI)
		if err != nil {
			return err
		}
		// 构建每个分片的task，执行
		t1 := NewTask(fn, link.String())
		t1.IsM3U8child = true
		t1.M3U8Dir = filepath.Join(config.GetConfig().SaveDir, t.Name)
		core.Submit(&t1)
	}
	log.Printf("该电影时长 %s \n", m3u8.CalculationTime(playbackDuration))

	core.Wait()

	ts := table.GetListErrorTask(t.Name)
	for _, errorTask := range ts {
		et := errorTask.(Task)
		core.Submit(&et)
	}

	if core.doneCount != len(segments) {
		log.Printf("任务下载不完整(%d\\%d)\n", core.doneCount, len(segments))
		return errors.New("任务失败了！！！！！！！！！！！！！！！")
	}

	// 合并所有分片
	if err := m3u8Downloader.MergeFiles(saveDir); err != nil {
		return err
	}
	_ = os.RemoveAll(saveDir) // 删除文件夹

	log.Printf("%s ===================> 任务完成,耗时 %s\n", t.Name, time.Now().Sub(beginTime))

	return nil
}

func (t *Task) video() error {
	var dir string
	if t.IsM3U8child {
		dir = t.M3U8Dir
	} else {
		dir = config.GetConfig().SaveDir
	}

	savePath := filepath.Join(dir, fmt.Sprintf("%s.%s", t.Name, t.ext))
	d := video.NewDownloader(t.Name, t.Link, savePath)
	if err := d.Execute(t.IsM3U8child); err != nil {
		return err
	}
	return nil
}

func ParseTaskList() ([]Task, error) {
	tasks := make([]Task, 0)
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
		tasks = append(tasks, NewTask(key, value))

		i++
	}

	return tasks, nil
}
