package main

import (
	"dv/base"
	"dv/config"
	"dv/m3u8"
	"dv/video"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type Task struct {
	base.Logger

	beginTime   time.Time      // 起始时间
	fileName    string         // 文件名（无尾缀）
	saveDir     string         // 存放位置
	fileUrl     string         // http地址
	videoScript string         // 视频格式（用于文件尾缀）
	videoType   base.VideoTpye // 视频类型
	Do          func()

	errorCount int // 连续错误数
}

func (t *Task) loadFunc() {
	index := strings.LastIndex(t.fileUrl, ".")
	vs := t.fileUrl[index+1:]
	t.videoScript = vs

	switch vs {
	case base.M3u8Type:
		t.videoType = base.M3u8Type
		t.Do = t.m3u8
	default:
		t.videoType = base.SingleType
		t.Do = t.video
	}

}

func (t *Task) m3u8() {
	// 提取所有分片地址
	d := m3u8.NewDownloader(t.fileName, config.GetConfig().SaveDir, t.fileUrl)
	d.SetHeader(config.Header)
	d.SetClient(config.Client)
	p, err := d.ExtractContain()
	if err != nil {
		t.errorCount++
		t.Fail(t.fileName, err.Error())
		return
	}

	segments := p.Segments
	tCore := NewCore(config.GetConfig())
	tCore.vacancy = make(chan struct{}, 10)
	tCore.SetGroup(len(segments))
	var playbackDuration float32
	for index, segment := range segments {
		fn := fmt.Sprintf("%s_part_%d", t.fileName, index)
		link := d.M3u8BaseLink + segment.URI
		task := NewTask(fn, d.SaveDir, link)
		playbackDuration += segment.Duration
		tCore.AddTask(task)
	}
	log.Printf("该电影时长 %s \n", m3u8.CalculationTime(playbackDuration))
	tCore.Run()
	tCore.Wait()

	// 合并所有分片
	if err := d.MergeFiles(); err != nil {
		t.Fail(t.fileName, err.Error())
		return
	}

	_ = os.RemoveAll(d.SaveDir) // 删除文件夹

	t.Done(t.fileName, fmt.Sprintf("任务完成,耗时 %s", time.Now().Sub(t.beginTime)))
}

func (t *Task) video() {
	// 执行
	d := video.NewDownloader(t.fileName, t.saveDir, t.fileUrl)
	d.SetHeader(config.Header)
	d.SetClient(config.Client)
	d.SetScript(t.videoScript)
	d.SetVideoType(t.videoType)
	if err := d.Execute(); err != nil {
		t.errorCount++
		if strings.HasSuffix(err.Error(), io.EOF.Error()) {
			// 若链接被关闭则不打印
			return
		}
		t.Fail(t.fileName, err.Error())
	} else {
		t.errorCount = 0
	}

	switch t.videoType {
	case base.M3u8Type:
		break
	case base.SingleType:
		t.Done(t.fileName, fmt.Sprintf("任务完成,耗时 %s", time.Now().Sub(t.beginTime)))
	}
}

func NewTask(name, saveDir, u string) *Task {
	t := &Task{
		beginTime: time.Now(),
		fileName:  name,
		saveDir:   saveDir,
		fileUrl:   u,
	}
	t.loadFunc()
	return t
}
