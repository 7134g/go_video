package task_control

import (
	"bytes"
	"dv/internel/serve/api/internal/util/calc"
	"dv/internel/serve/api/internal/util/m3u8"
	"dv/internel/serve/api/internal/util/table"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type download struct {
	key      string // 任务标识
	fileDir  string // 目录
	fileName string // 文件名

	fileSize      int64         // 目前文件大小
	totalFileSize int64         // 文件总大小
	stop          chan struct{} // 打印进度
}

func newDownload(key, fileDir, fileName string) *download {
	_ = os.MkdirAll(fileDir, 0700)
	return &download{
		key:      key,
		fileDir:  fileDir,
		fileName: fmt.Sprintf("%s.mp4", fileName),
		stop:     make(chan struct{}),
	}
}

func buildKey(id uint, name string, childFlag ...string) string {
	key := fmt.Sprintf("%d_%s", id, name)
	if len(childFlag) > 0 {
		for _, child := range childFlag {
			key = fmt.Sprintf("%s_%s", key, child)
		}
	}

	return key
}

func (d *download) getM3u8File(client *http.Client, req *http.Request) ([]*m3u8.Segment, error) {
	// 构建请求
	buf := bytes.NewBuffer(nil)
	if err := d.get(client, req, buf); err != nil {
		return nil, err
	}
	logx.Debug(buf.String())
	m3u8Data, err := m3u8.ParseM3u8Data(buf)
	if err != nil {
		return nil, err
	}

	if len(m3u8Data.MasterPlaylist) != 0 {
		// 下载最高清的视频
		index := m3u8Data.GetMaxBandWidth()
		if index < 0 {
			return nil, errors.New("解析失败")
		}
		link, err := url.Parse(req.URL.String())
		if err != nil {
			return nil, err
		}
		link, err = link.Parse(m3u8Data.MasterPlaylist[index].URI)
		if err != nil {
			return nil, err
		}

		request, err := http.NewRequest(http.MethodGet, link.String(), nil)
		if err != nil {
			return nil, err
		}
		request.Header = req.Header
		return d.getM3u8File(client, request)
	}

	for _, key := range m3u8Data.Keys {
		if key.Method == m3u8.CryptMethodNONE {
			continue
		}
		// 获取加密密匙
		link, err := url.Parse(req.URL.String())
		if err != nil {
			return nil, err
		}
		aesUrl, err := link.Parse(key.URI)
		if err != nil {
			return nil, err
		}
		aesBuf := bytes.NewBuffer(nil)
		request, err := http.NewRequest(http.MethodGet, aesUrl.String(), nil)
		if err != nil {
			return nil, err
		}
		request.Header = req.Header
		if err := d.get(client, request, aesBuf); err != nil {
			return nil, err
		}

		table.CryptoVideoTable.Set(d.key, aesBuf.Bytes())
		break
	}

	return m3u8Data.Segments, nil
}

func (d *download) get(client *http.Client, req *http.Request, write io.Writer) error {
	d.stop = make(chan struct{})
	// 构建请求
	logx.Debug(req.URL.String())
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("%d %s", resp.StatusCode, resp.Status))
	}

	if resp.ContentLength <= 0 {
		return errors.New(fmt.Sprintf("%s 跳过数据内容大于等于文件大小，因此不下载\n", d.fileName))
	}

	ctxRange := resp.Header.Get("Content-Range")
	if len(ctxRange) == 0 {
		// 记录文件总大小
		d.totalFileSize = resp.ContentLength
	} else {
		// Content-Range: bytes 1629222-5510871/5510872 取 5510872
		// 5510872 指的是文件总大小
		completeFileSizeString := ctxRange[strings.LastIndex(ctxRange, "/")+1:]
		completeFileSize, err := strconv.Atoi(completeFileSizeString)
		if err != nil {
			return err
		}
		d.totalFileSize = int64(completeFileSize)
	}

	go d.printDownloadMessage()
	return d.rw(resp.Body, write)
}

func (d *download) rw(read io.Reader, write io.Writer) error {
	defer close(d.stop)

	bs := make([]byte, 1048576) // 每次读取http内容的大小(1mb)
	for {
		rn, err := read.Read(bs)
		if err != nil {
			if err == io.EOF {
				// 完成
				d.fileSize += int64(rn)
				_, _ = write.Write(bs[:rn])
				return nil
			}
			return err
		}

		d.fileSize += int64(rn)
		_, err = write.Write(bs[:rn])
		if err != nil {
			return err
		}
	}
}

func (d *download) printDownloadMessage() {
	if len(strings.Split(d.key, "_")) > 2 {
		return
	}
	var now = time.Now()                                         // 记录耗时
	var fileSize = float64(d.totalFileSize) / 1024 / 1024 / 1024 // gb
	var lastNowRS float64                                        // 上一次打印消息的已读数据长度

	ticker := time.NewTicker(time.Second * 3) // 间隔时间打印
	for {
		var msg string
		select {
		case <-ticker.C:
			nowRS := float64(d.fileSize)
			score := nowRS / float64(d.totalFileSize) * 100
			dataByTime := (nowRS - lastNowRS) / float64(3) // 间隔时间内下载的数据, byte
			speed, unit := calc.UnitReturn(dataByTime)
			msg = fmt.Sprintf("百分比 %.2f 速度 %.3f %s/s | %.3f GB", score, speed, unit, fileSize)
			lastNowRS = nowRS
			logx.Infof("%s %s\n", d.fileName, msg)
		case <-d.stop:
			averageSpeed := float64(d.fileSize) / time.Since(now).Seconds() // 本次每秒下载字节数
			speed, unit := calc.UnitReturn(averageSpeed)
			msg = fmt.Sprintf("平均速度 %.2f %s/s <======== done", speed, unit)
			logx.Infof("%s %s\n", d.fileName, msg)
			return
		}
	}
}
