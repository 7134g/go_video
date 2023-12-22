package task_control

import (
	"dv/internel/serve/api/internal/svc/task_control/m3u8"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type cell struct {
	faCell *cell

	TaskId   uint
	TaskName string
	Name     string
	Url      string
	Dir      string

	client *http.Client
	header http.Header
	Do     func() error
}

func (c *cell) DownloadM3u8() error {
	now := time.Now()
	// 解析出所有的视频文件连接
	// 将解析的连接生成新的 cell 写回 TaskControl.Submit
	p, err := m3u8.ExtractContain(c.Name, c.Url)
	if err != nil {
		return err
	}
	segments := p.Segments

	// 创建存放所有分片的文件夹
	saveDir := filepath.Join(c.Dir, c.Name)
	info, err := os.Stat(saveDir)
	if err != nil || !info.IsDir() {
		_ = os.MkdirAll(saveDir, os.ModeDir)
	}

	var playbackDuration float32 // 该视频总时间
	for _, segment := range segments {
		playbackDuration += segment.Duration
	}
	logx.Info("该电影时长 %s \n", m3u8.CalculationTime(playbackDuration))

	// 并发控制
	vacancy := make(chan struct{}, tcConfig.cfg.ConcurrencyM3u8)

	wg := sync.WaitGroup{}
	for index, segment := range segments {
		link, err := url.Parse(c.Url)
		if err != nil {
			return err
		}
		link, err = link.Parse(segment.URI)
		if err != nil {
			return err
		}

		particle := &cell{
			faCell: c,
			TaskId: c.TaskId,

			TaskName: c.TaskName,
			Name:     fmt.Sprintf("%06d_part_%s", index, c.Name),
			Url:      link.String(),
			Dir:      saveDir,

			client: c.client,
			header: c.header,
			Do:     c.DownloadVideo,
		}

		wg.Add(1)
		vacancy <- struct{}{}
		tc.Submit(func() {
			defer func() {
				defer wg.Done()
				<-vacancy
			}()
			if err := particle.Do(); err != nil {
				logx.Error(err, saveErrorCellData(particle))
			}
		}, true)
	}
	wg.Wait()

	// 合并所有分片
	if tcConfig.cfg.UseFfmpeg {
		if err := m3u8.MergeFilesFfmpeg(saveDir, c.Name, tcConfig.cfg.FfmpegPath); err != nil {
			return err
		}
	} else {
		if err := m3u8.MergeFiles(c.Name, saveDir); err != nil {
			return err
		}
	}
	_ = os.RemoveAll(saveDir) // 删除文件夹

	logx.Info("%s ===================> 任务完成,耗时 %s\n", c.Name, time.Since(now))

	return nil
}

func (c *cell) DownloadVideo() error {
	// todo
	return nil
}
