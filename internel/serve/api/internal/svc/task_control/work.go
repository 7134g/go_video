package task_control

import (
	"bytes"
	"dv/internel/serve/api/internal/util/aes"
	"dv/internel/serve/api/internal/util/curl"
	"dv/internel/serve/api/internal/util/m3u8"
	"dv/internel/serve/api/internal/util/model"
	"dv/internel/serve/api/internal/util/table"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func fail(err error) func() error {
	return func() error {
		return err
	}
}

type work struct {
	task model.Task
}

func newWork(task model.Task) *work {
	return &work{
		task: task,
	}
}

type particleFunc func() error

func (w work) parseTask() (*download, particleFunc) {
	var request *http.Request
	var err error
	switch w.task.Type {
	case model.TypeUrl:
		request, err = http.NewRequest(http.MethodGet, w.task.Data, nil)
		if err != nil {
			return nil, fail(err)
		}
		request.Header = tcConfig.Headers
	case model.TypeCurl:
		_url, header, err := curl.Parse(w.task.Data)
		if err != nil {
			return nil, fail(err)
		}
		request, err = http.NewRequest(http.MethodGet, _url, nil)
		if err != nil {
			return nil, fail(err)
		}
		request.Header = header
	case model.TypeProxy:
		if err := json.Unmarshal([]byte(w.task.Data), request); err != nil {
			return nil, fail(err)
		}
	default:
		return nil, fail(errors.New("type error"))
	}

	d := newDownload(
		buildKey(w.task.ID, w.task.Name),
		tcConfig.SaveDir,
		w.task.Name,
	)
	switch w.task.VideoType {
	case model.VideoTypeMp4:
		//d.fileName = fmt.Sprintf("%s.mp4", d.fileName)
		return d, w.getVideo(d, request)
	case model.VideoTypeM3u8:
		d.key = buildKey(w.task.ID, w.task.Name, "m3u8")
		return d, w.getM3u8(d, request)
	default:
		return nil, fail(errors.New("video type error"))
	}

}

func (w work) getVideo(d *download, req *http.Request) func() error {
	savePath := filepath.Join(d.fileDir, d.fileName)
	var flag = os.O_RDWR | os.O_CREATE | os.O_APPEND
	if len(strings.Split(d.key, "_")) > 2 {
		flag = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	}

	file, err := os.OpenFile(savePath, flag, os.ModePerm)
	if err != nil {
		return fail(err)
	}
	info, err := file.Stat()
	if err != nil {
		return fail(err)
	}
	d.fileSize = info.Size()
	return func() error {
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		//head := http.Header{}
		//if err := copier.Copy(&head, &header); err != nil {
		//	return err
		//}
		if d.fileSize > 0 {
			req.Header.Set("Range", fmt.Sprintf("bytes=%d-", d.fileSize))
		}

		return d.get(tcConfig.Client, req, file)
	}
}

func (w work) getM3u8(d *download, req *http.Request) func() error {
	beginTime := time.Now()
	segments, err := d.getM3u8File(tcConfig.Client, req)
	if err != nil {
		return fail(err)
	}

	var playbackDuration float32 // 该视频总时间
	for _, segment := range segments {
		playbackDuration += segment.Duration
	}
	logx.Infof("%s 该电影时长 %s \n", w.task.Name, m3u8.CalculationTime(playbackDuration))

	concurrency := tcConfig.ConcurrencyM3u8
	if uint(len(segments))/(concurrency*concurrency) > concurrency {
		concurrency = uint(len(segments)) / (concurrency * concurrency)
	}

	dir := filepath.Join(tcConfig.SaveDir, w.task.Name)
	core := NewTaskControl(concurrency)
	core.start()
	go core.printDownloadProgress(uint(len(segments)))
	for index, segment := range segments {
		link, err := url.Parse(req.URL.String())
		if err != nil {
			return fail(err)
		}
		link, err = link.Parse(segment.URI)
		if err != nil {
			return fail(err)
		}

		fileName := fmt.Sprintf("%s_%d", w.task.Name, index)
		pathPart := strings.Split(link.Path, ".")
		if len(pathPart) > 0 {
			fileName = fmt.Sprintf("%s.%s", fileName, pathPart[len(pathPart)-1])
		}
		dChild := newDownload(
			buildKey(w.task.ID, fileName, "m3u8"),
			dir,
			fileName,
		)

		if crypto, exist := table.CryptoVideoTable.Get(d.key); exist {
			core.submit(func() error {
				buf := bytes.NewBuffer(nil)
				if err := dChild.get(tcConfig.Client, req, buf); err != nil {
					return err
				}
				table.M3u8DownloadDataLen.Set(d.key, uint(buf.Len()))

				data := aes.AESDecrypt(buf.Bytes(), crypto)
				if data == nil {
					return errors.New("视频格式解析失败")
				}
				savePath := filepath.Join(dir, dChild.fileName)
				f, err := os.Create(savePath)
				if err != nil {
					return err
				}

				_, err = io.Copy(f, buf)

				return err
			}, dChild)
		} else {
			request, err := http.NewRequest(http.MethodGet, link.String(), nil)
			if err != nil {
				return fail(err)
			}
			request.Header = req.Header
			core.submit(w.getVideo(dChild, request), dChild)
		}

	}
	core.wg.Wait()

	// 合并所有分片
	if tcConfig.UseFfmpeg {
		if err := m3u8.MergeFilesFfmpeg(dir, d.fileName, tcConfig.FfmpegPath); err != nil {
			return fail(err)
		}
	} else {
		if err := m3u8.MergeFiles(dir); err != nil {
			return fail(err)
		}
	}

	//_ = os.RemoveAll(dir) // 删除文件夹

	logx.Infof("%s ===================> 任务完成,耗时 %s\n", w.task.Name, time.Since(beginTime))
	return nil
}
