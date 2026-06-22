package controller

import (
	"context"
	"fmt"
	"go_video/pkg/m3u8"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync/atomic"
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
	task.Progress.SetSegment(0, len(segments))

	// 快照 pool，避免 ApplyConfig 中途替换。
	c.mu.RLock()
	pool := c.downloadPool
	c.mu.RUnlock()

	var consecutiveErrors int32
	errChan := make(chan error, 1)
	group := pool.NewGroup()

	for i, seg := range segments {
		if task.ctx.Err() != nil {
			break
		}
		if atomic.LoadInt32(&consecutiveErrors) >= int32(c.config.MaxConsecutiveErrors) {
			break
		}

		segFile := filepath.Join(dir, fmt.Sprintf("%06d.ts", i))
		if _, err := os.Stat(segFile); err == nil {
			task.Progress.IncrementSegmentDone()
			continue
		}

		segURL := m3u8.ResolveURL(seg.URI, baseURL)
		key := parsed.GetSegmentKey(seg)

		idx := i
		downloadURL := segURL
		file := segFile
		segment := seg
		segKey := key

		group.Submit(downloadURL, func() error {
			if err := c.downloadSegment(task.ctx, downloadURL, file, task.Header, segment, segKey, parsed.MediaSequence); err != nil {
				atomic.AddInt32(&consecutiveErrors, 1)
				select {
				case errChan <- fmt.Errorf("segment %d failed: %w", idx, err):
				default:
				}
				return err
			}
			atomic.StoreInt32(&consecutiveErrors, 0)
			task.Progress.IncrementSegmentDone()
			return nil
		})
	}

	group.Wait()
	BroadcastMessage(task.ID, "下载结束..."+task.Name)

	if task.ctx.Err() != nil {
		return context.Canceled
	}

	if atomic.LoadInt32(&consecutiveErrors) >= int32(c.config.MaxConsecutiveErrors) {
		select {
		case err := <-errChan:
			return fmt.Errorf("max consecutive errors reached: %w", err)
		default:
			return fmt.Errorf("max consecutive errors reached")
		}
	}

	return nil
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
