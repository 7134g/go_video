package controller

import (
	"context"
	"fmt"
	m3u9 "go_video/pkg/m3u8"
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
	dir := filepath.Join(c.config.DownloadDir, task.Name)
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
		streamURL := m3u9.ResolveURL(parsed.MasterPlaylist[idx].URI, baseURL)
		parsed, baseURL, err = c.parseM3u8(task.ctx, streamURL, task.Header)
		if err != nil {
			return err
		}
	}

	// 预下载加密密钥
	header := task.Header
	if len(header) == 0 {
		header = convertHeaders(c.config.DefaultHeaders)
	}
	keyCache := m3u9.NewKeyCache()
	if parsed.HasEncryption() {
		for _, key := range parsed.Keys {
			if key.Method == m3u9.CryptMethodAES && key.URI != "" {
				keyData, err := m3u9.DownloadKey(key.URI, baseURL, header)
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

	var consecutiveErrors int32
	errChan := make(chan error, 1)
	group := c.downloadPool.NewGroup()

	for i, seg := range segments {
		if task.ctx.Err() != nil {
			break
		}
		if atomic.LoadInt32(&consecutiveErrors) >= int32(c.config.MaxConsecutiveErrors) {
			break
		}

		segFile := filepath.Join(dir, fmt.Sprintf("%06d.ts", i))
		if _, err := os.Stat(segFile); err == nil {
			task.Progress.SetSegment(i+1, len(segments))
			//fmt.Printf("已存在 ---> %s\n", segFile)
			continue
		}

		segURL := m3u9.ResolveURL(seg.URI, baseURL)
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
			task.Progress.SetSegment(idx+1, len(segments))
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

func (c *DownloadController) parseM3u8(ctx context.Context, m3u8URL string, header http.Header) (*m3u9.M3u8, *url.URL, error) {
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

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}
	if c.config.HttpProxyAddress != "" {
		proxyURL, _ := url.Parse("http://" + c.config.HttpProxyAddress)
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	parsed, err := m3u9.ParseM3u8Data(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return parsed, baseURL, nil
}

func (c *DownloadController) downloadSegment(ctx context.Context, segURL, filename string, header http.Header, seg *m3u9.Segment, key *m3u9.Key, mediaSeq uint64) error {
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

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}
	if c.config.HttpProxyAddress != "" {
		proxyURL, _ := url.Parse("http://" + c.config.HttpProxyAddress)
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 解密
	if key != nil && key.Method == m3u9.CryptMethodAES && len(key.KeyData) > 0 {
		iv, err := m3u9.ParseIV(key.IV, mediaSeq+uint64(seg.KeyIndex))
		if err != nil {
			return fmt.Errorf("parse IV failed: %w", err)
		}
		data, err = m3u9.Decrypt(data, key.KeyData, iv)
		if err != nil {
			return fmt.Errorf("decrypt failed: %w", err)
		}
	}

	return os.WriteFile(filename, data, 0644)
}
