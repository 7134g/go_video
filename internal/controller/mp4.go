package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// progressReader 包装 io.Reader，把读取的字节数累加到任务进度上。
type progressReader struct {
	r io.Reader
	p *Progress
}

func (pr *progressReader) Read(buf []byte) (int, error) {
	n, err := pr.r.Read(buf)
	if n > 0 {
		pr.p.AddDone(int64(n))
	}
	return n, err
}

func (c *DownloadController) downloadMp4(task *DTask) error {
	filename := safeJoin(c.config.DownloadDir, task.Name) + ".mp4"

	var localSize int64 = 0
	if info, err := os.Stat(filename); err == nil {
		localSize = info.Size()
	}

	req, err := http.NewRequestWithContext(task.ctx, "GET", task.URL, nil)
	if err != nil {
		return err
	}
	header := task.Header
	if len(header) == 0 {
		header = convertHeaders(c.config.DefaultHeaders)
	}
	for k, v := range header {
		req.Header[k] = v
	}

	if localSize > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", localSize))
	}

	resp, err := c.httpClient.Get().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var totalSize int64
	if resp.StatusCode == http.StatusPartialContent {
		totalSize = localSize + resp.ContentLength
	} else {
		totalSize = resp.ContentLength
		localSize = 0
	}
	task.Progress.SetTotal(totalSize)
	task.Progress.AddDone(localSize)

	var file *os.File
	if localSize > 0 {
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	} else {
		file, err = os.Create(filename)
	}
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, &progressReader{r: resp.Body, p: task.Progress})
	return err
}
