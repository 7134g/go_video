package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// HasExactlyOneHttp 检查字符串中是否只包含 1 个 http 或 https
func HasExactlyOneHttp(input string) bool {
	// 定义正则：h开头，接 tt，接 p 或 ps，接 ://
	// \b 可以确保匹配单词边界，防止类似 "shhttps://" 的干扰
	re := regexp.MustCompile(`https?://`)

	// FindAllString 返回所有匹配的子串，-1 表示匹配所有
	matches := re.FindAllString(input, -1)

	// 如果匹配到的切片长度为 1，则返回 true
	return len(matches) == 1
}

func GetVideo(u *url.URL) (string, bool) {

	path := u.Path
	if strings.HasSuffix(path, ".m3u8") {
		return "m3u8", true
	}
	if strings.HasSuffix(path, ".mp4") {
		return "mp4", true

	}

	return "", false
}

func Capture(req *http.Request) *VideoTask {

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

func extractTitleFromHTML(body []byte) string {
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

var WebTree = map[string]map[string]*webContent{}

func addWeb(tabId string, u string, body []byte, header http.Header) {
	var dat = map[string]*webContent{}
	if v, ok := WebTree[tabId]; ok {
		dat = v
	}

	dat[u] = &webContent{
		u:      u,
		body:   body,
		header: header.Clone(),
	}

	WebTree[tabId] = dat
}

func search(tabId string) string {
	var dat = map[string]*webContent{}
	if v, ok := WebTree[tabId]; ok {
		dat = v
	}

	for _, v := range dat {
		title := extractTitleFromHTML(v.body)
		if title != "" {
			return title
		}
	}

	return ""
}

func SearchTitle(rawURL string) string {
	var tabId string
	for s, m := range WebTree {
		for _, v := range m {
			if v.u == rawURL {
				tabId = s
			}
		}
	}

	if tabId == "" {
		return ""
	}

	return search(tabId)
}
