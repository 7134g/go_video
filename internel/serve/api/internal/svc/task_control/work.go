package task_control

import (
	"dv/internel/serve/api/internal/util/m3u8"
	"dv/internel/serve/api/internal/util/model"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

func (w work) formatCurl() (string, http.Header) {
	// todo
	return "", nil
}

func (w work) parseTask() (*download, particleFunc) {
	var _url = w.task.Data
	var header = tcConfig.Headers
	switch w.task.Type {
	case model.TypeUrl:
		break
	case model.TypeCurl:
		_url, header = w.formatCurl()
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
		return d, w.getVideo(d, _url, header)
	case model.VideoTypeM3u8:
		return d, w.getM3u8(d, _url, header)
	default:
		return nil, fail(errors.New("video type error"))
	}

}

func (w work) getVideo(d *download, _url string, header http.Header) func() error {
	savePath := filepath.Join(tcConfig.SaveDir, d.fileName)
	file, err := os.OpenFile(savePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
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
		head := http.Header{}
		if err := copier.Copy(&head, &header); err != nil {
			return err
		}
		if d.fileSize > 0 {
			head.Set("Range", fmt.Sprintf("bytes=%d-", d.fileSize))
		}

		return d.get(tcConfig.Client, _url, head, file)
	}
}

func (w work) getM3u8(d *download, _url string, header http.Header) func() error {
	beginTime := time.Now()
	segments, err := d.getM3u8File(tcConfig.Client, _url, header)
	if err != nil {
		return fail(err)
	}

	var playbackDuration float32 // 该视频总时间
	for _, segment := range segments {
		playbackDuration += segment.Duration
	}
	logx.Infof("该电影时长 %s \n", m3u8.CalculationTime(playbackDuration))

	concurrency := tcConfig.ConcurrencyM3u8
	if uint(len(segments))/(concurrency*concurrency) > concurrency {
		concurrency = uint(len(segments)) / (concurrency * concurrency)
	}

	dir := filepath.Join(tcConfig.SaveDir, w.task.Name)
	core := NewTaskControl(concurrency)
	core.start()
	for index, segment := range segments {
		fileName := fmt.Sprintf("%s_%d", w.task.Name, index)

		d := newDownload(
			buildKey(w.task.ID, fileName),
			dir,
			fileName,
		)
		link, err := url.Parse(_url)
		if err != nil {
			return fail(err)
		}
		link, err = link.Parse(segment.URI)
		if err != nil {
			return fail(err)
		}

		core.Submit(w.getVideo(d, _url, header), d)
	}
	core.wg.Wait()

	// 合并所有分片
	if tcConfig.UseFfmpeg {
		if err := m3u8.MergeFilesFfmpeg(dir, w.task.Name, tcConfig.FfmpegPath); err != nil {
			return fail(err)
		}
	} else {
		if err := m3u8.MergeFiles(w.task.Name, dir); err != nil {
			return fail(err)
		}
	}

	_ = os.RemoveAll(dir) // 删除文件夹

	logx.Infof("%s ===================> 任务完成,耗时 %s\n", w.task.Name, time.Since(beginTime))
	return nil
}
