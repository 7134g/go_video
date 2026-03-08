package proxy

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Interceptor interface {
	Intercept(req *http.Request, resp *http.Response) *VideoTask
}

type VideoDetector struct{}

func (v *VideoDetector) IsVideo(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	path := u.Path
	return strings.HasSuffix(path, ".m3u8") || strings.HasSuffix(path, ".mp4")
}

type RequestCapture struct{}

func (r *RequestCapture) Capture(req *http.Request) *VideoTask {
	headers := make(map[string]string)
	for k, v := range req.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
	var body []byte
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(body))
	}
	return &VideoTask{
		URL:     req.URL.String(),
		Method:  req.Method,
		Headers: headers,
		Body:    body,
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
