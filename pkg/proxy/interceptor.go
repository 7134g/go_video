package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Interceptor interface {
	Intercept(req *http.Request, resp *http.Response) *VideoTask
}

type VideoDetector struct{}

func (v *VideoDetector) GetVideo(rawURL string) (string, bool) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", false
	}
	path := u.Path
	if strings.HasSuffix(path, ".m3u8") {
		return "m3u8", true
	}
	if strings.HasSuffix(path, ".mp4") {
		return "mp4", true

	}

	return "", false
}

type RequestCapture struct{}

func (r *RequestCapture) Capture(req *http.Request) *VideoTask {

	headers, _ := json.Marshal(req.Header)

	var body []byte
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(body))
	}
	now := time.Now()
	return &VideoTask{
		URL:      req.URL.String(),
		Method:   req.Method,
		Headers:  string(headers),
		Body:     body,
		Title:    fmt.Sprintf("%d", now.UnixMilli()),
		Type:     "",
		CreateAt: now,
	}
}

type TitleExtractor struct{}

func (t *TitleExtractor) Extract(resp *http.Response) string {
	if resp == nil || !strings.HasSuffix(resp.Request.URL.Path, ".html") {
		return ""
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	content := string(body)
	start := strings.Index(content, "<title>")
	if start == -1 {
		return ""
	}
	start += 7
	end := strings.Index(content[start:], "</title>")
	if end == -1 {
		return ""
	}
	return content[start : start+end]
}

type webContent struct {
	u      string
	body   []byte
	header http.Header
}

var WebTree map[string]*webContent
