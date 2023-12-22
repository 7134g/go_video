package task_control

import (
	"dv/internel/serve/api/internal/model"
	"net/http"
)

type work struct {
	task model.Task
}

func newWork(task model.Task) *work {
	return &work{
		task: task,
	}
}

func (w work) parseTask() *cell {
	c := &cell{
		TaskId:   w.task.ID,
		TaskName: w.task.Name,
		Name:     w.task.Name,
		Url:      w.task.Data,
		Dir:      tcConfig.cfg.SaveDir,

		client: &http.Client{Transport: tcConfig.Transport},
		header: tcConfig.Headers,
	}
	switch w.task.Type {
	case model.TypeUrl:
		return w.parseVideo(c)
	case model.TypeCurl:
		u, header := w.formatCurl()
		c.Url = u
		c.header = header
		return w.parseVideo(c)
	}

	return nil
}

func (w work) parseVideo(c *cell) *cell {
	switch w.task.VideoType {
	case model.VideoTypeMp4:
		c.Do = c.DownloadVideo
	case model.VideoTypeM3u8:
		c.Do = c.DownloadM3u8
	}
	return c
}

func (w work) formatCurl() (string, http.Header) {
	// todo
	return "", nil
}
