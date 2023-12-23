package task_control

import (
	"dv/internel/serve/api/internal/model"
	"dv/internel/serve/api/internal/svc/task_control/m3u8"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"os"
	"path/filepath"
	"sync"
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
	var _url = w.task.Data
	var header = tcConfig.Headers
	switch w.task.Type {
	case model.TypeUrl:
	case model.TypeCurl:
		_url, header = w.formatCurl()
	default:
		return nil, fail(errors.New("type error"))
	}

	d := &download{
		key:      buildKey(w.task.ID, w.task.Name),
		fileDir:  tcConfig.cfg.SaveDir,
		fileName: w.task.Name,
		fileSize: 0,
	}
	savePath := filepath.Join(tcConfig.cfg.SaveDir, w.task.Name)
	file, err := os.OpenFile(savePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return nil, fail(err)
	}
	info, err := file.Stat()
	if err != nil {
		return nil, fail(err)
	}
	d.fileSize = info.Size()

	switch w.task.VideoType {
	case model.VideoTypeMp4:
		return d, func() error {
			head := http.Header{}
			if err := copier.Copy(&head, &header); err != nil {
				return err
			}
			if d.fileSize > 0 {
				head.Set("Range", fmt.Sprintf("bytes=%d-", d.fileSize))
			}

			return d.get(tcConfig.Client, _url, head, file)
		}
	case model.VideoTypeM3u8:
		return d, func() error {
			segments, err := d.getM3u8File(tcConfig.Client, _url, header)
			if err != nil {
				return err
			}

			var playbackDuration float32 // 该视频总时间
			for _, segment := range segments {
				playbackDuration += segment.Duration
			}
			logx.Infof("该电影时长 %s \n", m3u8.CalculationTime(playbackDuration))

			core := &TaskControl{
				wg:      sync.WaitGroup{},
				mux:     sync.Mutex{},
				running: false,
			}
			size := int(tcConfig.cfg.ConcurrencyM3u8)
			if len(segments)/(size*size) > size {
				core.vacancy = make(chan struct{}, len(segments)/(size*size))
			} else {
				core.vacancy = make(chan struct{}, size)
			}

			// todo 派生小任务

			return nil
		}
	default:
		return nil, fail(errors.New("video type error"))
	}

}

func (w work) formatCurl() (string, http.Header) {
	// todo
	return "", nil
}
