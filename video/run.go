package video

import (
	"dv/base"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//const (
//	mp4 = "mp4"
//	ts  = "ts"
//)

type DownVideo struct {
	base.Downloader

	readSize              int   // 当前已读取的文件大小
	existSize             int64 // 已存在的文件断点续传大小
	responseContentLength int64 // 该请求仍需要下载长度
	fileFutureSize        int64 // 文件总大小

	stop chan struct{} // 停止打印下载信息
}

func NewDownloader(taskName, saveDir, httpUrl string) DownVideo {
	m := DownVideo{}
	m.TaskName = taskName
	m.SaveDir = saveDir
	m.Link = httpUrl
	m.stop = make(chan struct{})
	return m
}

func (d DownVideo) Execute() error {
	// 打开文件
	filePath := filepath.Join(d.SaveDir, fmt.Sprintf("%s.%s", d.TaskName, d.GetScript()))
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	d.existSize = info.Size()

	// 构建请求
	res, err := http.NewRequest(http.MethodGet, d.Link, nil)
	if err != nil {
		return err
	}
	res.Header = d.GetHeader()
	if d.GetVideoType() == base.SingleType && d.existSize != 0 {
		// 断点续传，跳过已经下载的内容
		res.Header.Set("Range", fmt.Sprintf("bytes=%d-", d.existSize))
	}
	resp, err := d.GetClient().Do(res)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	// 处理请求结果
	if resp.ContentLength <= 0 {
		d.Done(d.TaskName, "跳过数据内容大于等于文件大小，因此不下载")
		return nil
	}
	switch d.GetVideoType() {
	case base.M3u8Type:
		if resp.ContentLength == d.existSize {
			d.Done(d.TaskName, "该分片已经下载过")
		}
		// 从头开始写
		_, err := f.Seek(0, 0)
		if err != nil {
			return err
		}
		break
	case base.SingleType:
		ctxRange := resp.Header.Get("Content-Range")
		if len(ctxRange) == 0 {
			// 首次请求，记录文件总大小
			d.fileFutureSize = resp.ContentLength
		} else {
			// Content-Range: bytes 1629222-5510871/5510872 取 5510872
			// 5510872 指的是文件总大小
			completeFileSizeString := ctxRange[strings.LastIndex(ctxRange, "/")+1:]
			completeFileSize, err := strconv.Atoi(completeFileSizeString)
			if err != nil {
				return err
			}
			d.fileFutureSize = int64(completeFileSize)
		}
	}

	d.responseContentLength = resp.ContentLength // 剩余大小
	bs := make([]byte, 1048576)                  // 每次读取http内容的大小(1mb)
	defer close(d.stop)
	if d.GetVideoType() == base.SingleType {
		go d.printDownloadMessage()
	}
	for {
		rn, err := resp.Body.Read(bs)
		if err != nil {
			if err == io.EOF {
				// 完成
				d.readSize += rn
				_, _ = f.Write(bs[:rn])
				return nil
			}
			return err
		}

		d.readSize += rn
		_, err = f.Write(bs[:rn])
		if err != nil {
			return err
		}

	}

}

func (d *DownVideo) setExistSize(value int64) {
	d.existSize = value
}

func (d *DownVideo) printDownloadMessage() {
	ticker := d.Ticker()                                       // 间隔时间打印
	fileSize := float64(d.fileFutureSize) / 1024 / 1024 / 1024 // gb
	var lastNowRS float64                                      // 上一次打印消息的已读数据长度
	now := time.Now().Unix()                                   // 记录耗时
	for {
		var msg string
		select {
		case <-ticker.C:
			nowRS := float64(d.readSize)
			score := (nowRS + float64(d.existSize)) / float64(d.fileFutureSize) * 100
			dataByTime := (nowRS - lastNowRS) / float64(base.GetInterval()) // 间隔时间内下载的数据, byte
			speed, unit := unitReturn(dataByTime)
			msg = fmt.Sprintf("百分比 %.2f 速度 %.3f %s/s | %.3f GB", score, speed, unit, fileSize)
			lastNowRS = nowRS
			d.Doing(d.TaskName, msg)
		case <-d.stop:
			averageSpeed := float64(d.readSize) / float64(time.Now().Unix()-now) // 本次每秒下载字节数
			speed, unit := unitReturn(averageSpeed)
			msg = fmt.Sprintf("平均速度 %.2f %s/s <======== done", speed, unit)
			d.Done(d.TaskName, msg)
			return
		}
	}
}
