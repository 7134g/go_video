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

	beginTime time.Time // 起始时间
	fileName  string    // 文件名（无尾缀）
	saveDir   string    // 存放位置
	fileUrl   string    // http地址
	//videoScript string         // 视频格式（用于文件尾缀）
	//videoType   base.VideoType // 视频类型

	videoSetting video.VideoSetting
	Do           func()

	errorCount uint // 连续错误数
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
	} else {
		t.errorCount = 0
	}

	segments := p.Segments
	tCore := NewCore(config.GetConfig())
	// 设置并发大小
	size := int(config.GetConfig().ConcurrencyM3u8)
	if len(segments)/(size*size) > size {
		tCore.vacancy = make(chan struct{}, len(segments)/(size*size))
	} else {
		tCore.vacancy = make(chan struct{}, size)
	}
	//tCore.vacancy = make(chan struct{}, config.GetConfig().ConcurrencyM3u8) // 设置分片并发下载数
	tCore.SetGroupCount(len(segments))
	var playbackDuration float32 // 该视频总时间
	vSetting := video.VideoSetting{
		CryptoKey:    d.Crypto,
		CryptoMethod: string(d.CryptoMethod),
	}
	for index, segment := range segments {
		playbackDuration += segment.Duration

		fn := fmt.Sprintf("%s_part_%d", t.fileName, index)
		var link string
		if video.CompleteURL(segment.URI) {
			link = d.M3u8BaseLink + segment.URI
		} else {
			link = segment.URI
		}

		task := NewTask(fn, d.SaveDir, link, vSetting)
		task.setFunc(base.VideoListType)
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
	d.SetVideoSetting(t.videoSetting)
	if err := d.Execute(); err != nil {
		t.errorCount++
		if strings.HasSuffix(err.Error(), io.EOF.Error()) {
			// 若链接被关闭则不打印
			return
		}
		t.Fail(t.fileName, err.Error())
		return
	} else {
		t.errorCount = 0
	}

	//switch t.videoSetting.VideoType {
	//case base.M3u8Type:
	//	break
	//case base.SingleType:
	//	t.Done(t.fileName, fmt.Sprintf("任务完成,耗时 %s", time.Now().Sub(t.beginTime)))
	//}
}

func NewTask(name, saveDir, u string, vSetting video.VideoSetting) *Task {
	t := &Task{
		beginTime: time.Now(),
		fileName:  name,
		saveDir:   saveDir,
		fileUrl:   u,
	}
	index := strings.LastIndex(t.fileUrl, ".")
	vSetting.VideoExt = t.fileUrl[index+1:]
	t.videoSetting = vSetting

	return t
}

func (t *Task) setFunc(v string) {
	t.videoSetting.VideoCategory = v
	switch v {
	case base.AloneVideoType:
		switch t.videoSetting.VideoExt {
		case "m3u8":
			t.Do = t.m3u8
		default:
			t.Do = t.video
		}
	case base.VideoListType:
		t.Do = t.video
	}

}
