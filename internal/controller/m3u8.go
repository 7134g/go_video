package controller

import (
	"context"
	"fmt"
	"go_video/pkg/m3u8"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
)

func convertHeaders(m map[string]string) http.Header {
	h := make(http.Header)
	for k, v := range m {
		h.Set(k, v)
	}
	return h
}

func (c *DownloadController) downloadM3u8(task *DTask) error {
	BroadcastMessage(task.ID, "开始下载..."+task.Name)
	dir := safeJoin(c.config.DownloadDir, task.Name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	parsed, baseURL, err := c.parseM3u8(task.ctx, task.URL, task.Header)
	if err != nil {
		return err
	}

	// 处理 Master Playlist
	if parsed.IsMasterPlaylist() {
		idx := parsed.GetMaxBandWidth()
		if idx < 0 {
			return fmt.Errorf("no valid stream found")
		}
		streamURL := m3u8.ResolveURL(parsed.MasterPlaylist[idx].URI, baseURL)
		parsed, baseURL, err = c.parseM3u8(task.ctx, streamURL, task.Header)
		if err != nil {
			return err
		}
	}

	header := task.Header
	if len(header) == 0 {
		header = convertHeaders(c.config.DefaultHeaders)
	}
	keyCache := m3u8.NewKeyCache()
	if parsed.HasEncryption() {
		for _, key := range parsed.Keys {
			if key.Method == m3u8.CryptMethodAES && key.URI != "" {
				keyData, err := m3u8.DownloadKey(key.URI, baseURL, header)
				if err != nil {
					return fmt.Errorf("download key failed: %w", err)
				}
				keyCache.Set(key.URI, keyData)
				key.KeyData = keyData
			}
		}
	}

	segments := parsed.Segments
	task.Progress.SetTotal(int64(len(segments)))

	segURLs := make([]string, len(segments))
	segFiles := make([]string, len(segments))
	for i, seg := range segments {
		segURLs[i] = m3u8.ResolveURL(seg.URI, baseURL)
		segFiles[i] = filepath.Join(dir, fmt.Sprintf("%06d%s", i, m3u8.SegmentExtension(segURLs[i])))
	}

	// 清理旧版本遗留的、不属于当前播放列表的分片文件，避免合并时混入。
	cleanupSegments(dir, segFiles, parsed.InitSegment != nil)

	// 下载 fMP4 初始化分片（EXT-X-MAP）。
	if parsed.InitSegment != nil {
		initFile := filepath.Join(dir, m3u8.InitSegmentName)
		if _, err := os.Stat(initFile); err != nil {
			initURL := m3u8.ResolveURL(parsed.InitSegment.URI, baseURL)
			if err := c.downloadSegment(task.ctx, initURL, initFile, task.Header, parsed.InitSegment, parsed.GetSegmentKey(parsed.InitSegment), parsed.MediaSequence); err != nil {
				return fmt.Errorf("下载初始化分片失败: %w", err)
			}
		}
	}

	// 快照 pool，避免 ApplyConfig 中途替换。
	c.mu.RLock()
	pool := c.downloadPool
	c.mu.RUnlock()

	type failedSegment struct {
		idx  int
		url  string
		file string
		seg  *m3u8.Segment
		key  *m3u8.Key
	}
	var failedSegments []failedSegment
	var failedMu sync.Mutex

	group := pool.NewGroup()

	for i, seg := range segments {
		if task.ctx.Err() != nil {
			break
		}

		segFile := segFiles[i]
		if _, err := os.Stat(segFile); err == nil {
			task.Progress.IncrementDone()
			continue
		}

		segURL := segURLs[i]
		key := parsed.GetSegmentKey(seg)

		idx := i
		downloadURL := segURL
		file := segFile
		segment := seg
		segKey := key

		group.Submit(downloadURL, func() error {
			for attempt := 0; attempt < 3; attempt++ {
				if task.ctx.Err() != nil {
					return nil
				}
				if err := c.downloadSegment(task.ctx, downloadURL, file, task.Header, segment, segKey, parsed.MediaSequence); err != nil {
					if attempt == 2 {
						failedMu.Lock()
						failedSegments = append(failedSegments, failedSegment{idx, downloadURL, file, segment, segKey})
						failedMu.Unlock()
						slog.Warn("分片下载失败（已重试3次）", "idx", idx, "url", downloadURL, "file", file, "error", err)
					}
					continue
				}
				task.Progress.IncrementDone()
				return nil
			}
			return nil
		})
	}

	group.Wait()
	BroadcastMessage(task.ID, "下载结束..."+task.Name)

	if task.ctx.Err() != nil {
		return context.Canceled
	}

	if len(failedSegments) > 0 {
		BroadcastMessage(task.ID, fmt.Sprintf("重试 %d 个失败分片...", len(failedSegments)))
		var stillFailed []failedSegment
		for _, fs := range failedSegments {
			var lastErr error
			success := false
			for attempt := 0; attempt < 10; attempt++ {
				if task.ctx.Err() != nil {
					break
				}
				if err := c.downloadSegment(task.ctx, fs.url, fs.file, task.Header, fs.seg, fs.key, parsed.MediaSequence); err != nil {
					lastErr = err
					continue
				}
				task.Progress.IncrementDone()
				success = true
				break
			}
			if !success {
				stillFailed = append(stillFailed, fs)
				slog.Error("分片最终失败（共重试13次）", "idx", fs.idx, "url", fs.url, "file", fs.file, "error", lastErr)
			}
		}
		if len(stillFailed) > 0 {
			return fmt.Errorf("%d 个分片下载失败，已重试13次", len(stillFailed))
		}
	}

	return nil
}

// cleanupSegments 删除旧版本遗留的、与当前播放列表不匹配的分片文件。
// keep 为本次任务期望的分片路径；hasInit 为 false 时删除遗留的 init.mp4。
func cleanupSegments(dir string, keep []string, hasInit bool) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	keepSet := make(map[string]bool, len(keep))
	for _, f := range keep {
		keepSet[filepath.Base(f)] = true
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		switch {
		case name == m3u8.InitSegmentName:
			if !hasInit {
				_ = os.Remove(filepath.Join(dir, name))
			}
		case m3u8.IsSegmentFile(name) && !keepSet[name]:
			_ = os.Remove(filepath.Join(dir, name))
		}
	}
}

func (c *DownloadController) parseM3u8(ctx context.Context, m3u8URL string, header http.Header) (*m3u8.M3u8, *url.URL, error) {
	baseURL, _ := url.Parse(m3u8URL)

	req, err := http.NewRequestWithContext(ctx, "GET", m3u8URL, nil)
	if err != nil {
		return nil, nil, err
	}
	if len(header) == 0 {
		header = convertHeaders(c.config.DefaultHeaders)
	}
	for k, v := range header {
		req.Header[k] = v
	}

	resp, err := c.httpClient.Get().Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	parsed, err := m3u8.ParseM3u8Data(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return parsed, baseURL, nil
}

func (c *DownloadController) downloadSegment(ctx context.Context, segURL, filename string, header http.Header, seg *m3u8.Segment, key *m3u8.Key, mediaSeq uint64) error {
	req, err := http.NewRequestWithContext(ctx, "GET", segURL, nil)
	if err != nil {
		return err
	}
	if len(header) == 0 {
		header = convertHeaders(c.config.DefaultHeaders)
	}
	for k, v := range header {
		req.Header[k] = v
	}

	resp, err := c.httpClient.Get().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if key != nil && key.Method == m3u8.CryptMethodAES && len(key.KeyData) > 0 {
		iv, err := m3u8.ParseIV(key.IV, mediaSeq+uint64(seg.KeyIndex))
		if err != nil {
			return fmt.Errorf("parse IV failed: %w", err)
		}
		data, err = m3u8.Decrypt(data, key.KeyData, iv)
		if err != nil {
			return fmt.Errorf("decrypt failed: %w", err)
		}
	}

	return os.WriteFile(filename, data, 0644)
}
