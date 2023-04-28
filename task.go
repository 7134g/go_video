package main

import (
	"dv/base"
	"dv/config"
	"dv/m3u8"
	"dv/video"
	"fmt"
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
	switch t.fileUrl[index+1:] {
	case mp4Type:
		t.videoScript = mp4Type
		t.Do = t.video
	case tsType:
		t.videoScript = tsType
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
	d.SetHeader(config.HeaderMP4)
	d.SetClient(config.Client)
	segments, err := d.ExtractSegments()
	if err != nil {
		t.errorCount++
		t.Fail(t.fileName, err.Error())
		return
	}

	// 下载所有分片,分成五段
	tCore := NewCore(config.GetConfig())
	if len(segments) > 5 {
		tCore.vacancy = make(chan struct{}, len(segments)/5*5)
	} else {
		tCore.vacancy = make(chan struct{}, 5)
	}
	tCore.SetGroup(len(segments))
	for index, segment := range segments {
		fn := fmt.Sprintf("%s_%d", t.fileName, index)
		task := NewTask(fn, d.SaveDir, segment)
		tCore.AddTask(task)
	}
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

//func (t *Task) m3u8() {
//	cmdParams := []string{
//		config.GetConfig().M3U8Path,
//		t.fileUrl,
//		"--workDir",
//		config.GetConfig().SaveDir,
//		"--enableDelAfterDone",
//		"--saveName",
//		t.fileName,
//		"--headers",
//		config.HeaderM3U8,
//	}
//
//	if config.GetConfig().ProxyStatus {
//		cmdParams = append(cmdParams, "--proxyAddress", config.GetConfig().Proxy)
//	}
//
//	if err := m3u8.Execute(t.fileName, cmdParams); err != nil {
//		t.errorCount++
//		t.Fail(t.fileName, err.Error())
//	} else {
//		t.errorCount = 0
//	}
//}

func (t *Task) video() {
	// 执行
	d := video.NewDownloader(t.fileName, t.saveDir, t.fileUrl)
	d.SetHeader(config.HeaderMP4)
	d.SetClient(config.Client)
	d.SetScript(t.videoScript)
	if err := d.Execute(); err != nil {
		t.errorCount++
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
