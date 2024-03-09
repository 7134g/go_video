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

type particleFunc func(params []interface{}) error

func fail(err error) particleFunc {
	return func(params []interface{}) error {
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

func (w work) parseTask() (particleFunc, *download) {
	var request *http.Request
	var err error
	switch w.task.Type {
	case model.TypeUrl:
		request, err = http.NewRequest(http.MethodGet, w.task.Data, nil)
		if err != nil {
			return fail(err), nil
		}
		request.Header = tcConfig.Headers
	case model.TypeCurl:
		_url, header, err := curl.Parse(w.task.Data)
		if err != nil {
			return fail(err), nil
		}
		request, err = http.NewRequest(http.MethodGet, _url, nil)
		if err != nil {
			return fail(err), nil
		}
		request.Header = header
	case model.TypeProxy:
		request, err = http.NewRequest(http.MethodGet, w.task.Data, nil)
		if err != nil {
			return fail(err), nil
		}
		var header http.Header
		if err := json.Unmarshal([]byte(w.task.HeaderJson), &header); err != nil {
			return fail(err), nil
		}
		if len(header) == 0 {
			header = tcConfig.Headers
		}
		request.Header = header
	default:
		return fail(errors.New("type error")), nil

	}

	d := newDownload(
		w.task,
		tcConfig.SaveDir,
		w.task.Name,
		false,
	)
	d.req = request
	switch w.task.VideoType {
	case model.VideoTypeMp4:
		return w.getVideo, d
	case model.VideoTypeM3u8:
		return w.getM3u8, d
	default:
		return fail(errors.New("video type error")), nil
	}

}

func (w work) getVideo(params []interface{}) error {
	d := params[0].(*download)
	savePath := filepath.Join(d.fileDir, d.fileName)
	var flag = os.O_RDWR | os.O_CREATE | os.O_APPEND
	if w.task.VideoType != model.VideoTypeMp4 {
		flag = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	}

	file, err := os.OpenFile(savePath, flag, os.ModePerm)
	if err != nil {
		return err
	}
	info, err := file.Stat()
	if err != nil {
		return err
	}
	d.fileSize = info.Size()

	defer func(file *os.File) {
		info, _ := file.Stat()
		table.DownloadDataLen.Inc(w.task.ID, uint(info.Size()))
		_ = file.Close()
	}(file)
	if d.fileSize > 0 {
		d.req.Header.Set("Range", fmt.Sprintf("bytes=%d-", d.fileSize))
	}

	return d.get(tcConfig.Client, d.req, file)
}

func (w work) getM3u8(params []interface{}) error {
	d := params[0].(*download)
	beginTime := time.Now()
	segments, err := d.getM3u8File(tcConfig.Client, d.req)
	if err != nil {
		return err
	}

	var playbackDuration float32 // 该视频总时间
	for _, segment := range segments {
		playbackDuration += segment.Duration
	}
	logx.Infof("%v 该电影时长 %v \n", w.task.Name, m3u8.CalculationTime(playbackDuration))

	concurrency := tcConfig.ConcurrencyM3u8
	if uint(len(segments))/(concurrency*concurrency) > concurrency {
		concurrency = uint(len(segments)) / (concurrency * concurrency)
	}

	dir := filepath.Join(tcConfig.SaveDir, w.task.Name)
	core := NewTaskControl(concurrency)
	core.start()
	go core.printDownloadProgress(w.task, uint(len(segments)))
	for index, segment := range segments {
		link, err := url.Parse(d.req.URL.String())
		if err != nil {
			return err
		}
		link, err = link.Parse(segment.URI)
		if err != nil {
			return err
		}

		fileName := fmt.Sprintf("%s_%d", w.task.Name, index)
		pathPart := strings.Split(link.Path, ".")
		if len(pathPart) > 0 {
			fileName = fmt.Sprintf("%s.%s", fileName, pathPart[len(pathPart)-1])
		}
		dChild := newDownload(
			w.task,
			dir,
			fileName,
			true,
		)
		request, err := http.NewRequest(http.MethodGet, link.String(), nil)
		if err != nil {
			return err
		}
		request.Header = d.req.Header
		dChild.req = request

		if crypto, exist := table.CryptoVideoTable.Get(w.task.ID); exist {
			// 编码过的视频
			tf := func(params []any) error {
				buf := bytes.NewBuffer(nil)
				if err := dChild.get(tcConfig.Client, dChild.req, buf); err != nil {
					return err
				}
				//table.DownloadDataLen.Inc(w.task.ID, uint(buf.Len()))
				//_, err := os.Stat("./test.mp4")
				//if err != nil {
				//	f, _ := os.Create("./test.mp4")
				//	f.Write(buf.Bytes())
				//	f.Close()
				//}

				data := aes.AESDecrypt(buf.Bytes(), crypto)
				if data == nil {
					return errors.New("视频格式解析失败")
				}
				savePath := filepath.Join(dir, dChild.fileName)
				f, err := os.Create(savePath)
				if err != nil {
					return err
				}
				defer f.Close()

				_, err = io.Copy(f, io.NopCloser(bytes.NewReader(data)))

				return err
			}
			//dChild.req = d.req
			core.submit(tf, []any{dChild})
		} else {
			// 无编码
			tf := func(params []interface{}) error {
				err := w.getVideo(params)
				if err != nil {
					return err
				}

				return nil
			}
			core.submit(tf, []any{dChild})
		}

	}
	core.wg.Wait()
	core.Stop()
	_ = tasKModel.UpdateStatus(d.t.ID, model.StatusSuccess)
	logx.Infof("%s 任务完成 ！！！！！！！！", w.task.Name)

	// 合并所有分片
	if tcConfig.UseFfmpeg {
		if err := m3u8.MergeFilesFfmpeg(dir, d.fileName, tcConfig.FfmpegPath); err != nil {
			return err
		}
	} else {
		if err := m3u8.MergeFiles(dir); err != nil {
			return err
		}
	}

	//_ = os.RemoveAll(dir) // 删除文件夹

	logx.Infof("%s ===================> 任务完成,耗时 %s\n", w.task.Name, time.Since(beginTime))
	return nil
}
