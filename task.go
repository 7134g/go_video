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
)

const (
	mp4Type  = "video"
	tsType   = "ts"
	m3u8Type = "m3u8"
	unknown  = "unknown"
)

type Task struct {
	base.Logger

	fileName    string // 文件名（无尾缀）
	saveDir     string // 存放位置
	fileUrl     string // http地址
	videoScript string // 视频类型（用于文件尾缀）
	Do          func()

	errorCount int // 连续错误数
}

func (t *Task) loadFunc() {
	index := strings.LastIndex(t.fileUrl, ".")
	vs := t.fileUrl[index+1:]
	switch vs {
	case mp4Type, tsType:
		t.videoScript = vs
		t.Do = t.video
	case m3u8Type:
		t.videoScript = m3u8Type
		t.Do = t.m3u8
	default:
		t.videoScript = unknown
		t.Do = func() {
			t.Fail(t.fileName, fmt.Sprintf("地址：%v, 解析该类型视频失败", t.fileUrl))
		}
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
	// 下载所有分片,同一时间下载最大下载数为总分片的五分之一，最小为五
	tCore := NewCore(config.GetConfig())
	tCore.vacancy = make(chan struct{}, 10)
	//if len(segments) > 5 {
	//	tCore.vacancy = make(chan struct{}, len(segments)/5*5)
	//} else {
	//	tCore.vacancy = make(chan struct{}, 5)
	//}
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

	t.Done(t.fileName, "任务完成")
}

func (t *Task) video() {
	// 执行
	d := video.NewDownloader(t.fileName, t.saveDir, t.fileUrl)
	d.SetHeader(config.Header)
	d.SetClient(config.Client)
	d.SetScript(t.videoScript)
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
}

func NewTask(name, saveDir, u string) *Task {
	t := &Task{fileName: name, saveDir: saveDir, fileUrl: u}
	t.loadFunc()
	return t
}
