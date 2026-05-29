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
	"sync"
	"time"
)

// HasExactlyOneHttp 仅当字符串中恰好出现一次 http(s):// 时返回 true。
// 用于过滤埋点/跳转链接里嵌套了第二个 URL 的伪视频地址。
func HasExactlyOneHttp(input string) bool {
	re := regexp.MustCompile(`https?://`)
	return len(re.FindAllString(input, -1)) == 1
}

func GetVideo(u *url.URL) (string, bool) {
	switch {
	case strings.HasSuffix(u.Path, ".m3u8"):
		return "m3u8", true
	case strings.HasSuffix(u.Path, ".mp4"):
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
	start += len("<title>")
	end := strings.Index(content[start:], "</title>")
	if end == -1 {
		return ""
	}
	return strings.TrimSpace(content[start : start+end])
}

// webIndex 用于按浏览器 Tab 维度缓存"看到过的 URL"和"该 Tab 内出现过的 HTML 标题"。
// 视频任务命名来自这里：捕获到视频 URL 时取该 Tab 最近的标题作为任务名。
// 设计上：
//   - 只保留 HTML 中提取到的标题，不存原始 body（避免内存膨胀）
//   - 通过 maxTabs LRU 限容，避免长跑代理占用无界内存
//   - 全部访问串行化（sync.Mutex），消除原全局 map 的并发崩溃
type webIndex struct {
	mu      sync.Mutex
	tabs    map[string]*tabEntry
	order   []string // 进入顺序，便于 LRU 淘汰
	maxTabs int
}

type tabEntry struct {
	urls      map[string]struct{}
	lastTitle string
}

func newWebIndex(maxTabs int) *webIndex {
	return &webIndex{
		tabs:    make(map[string]*tabEntry),
		maxTabs: maxTabs,
	}
}

func (w *webIndex) record(tabId, rawURL string, body []byte) {
	if tabId == "" {
		return
	}
	title := extractTitleFromHTML(body)

	w.mu.Lock()
	defer w.mu.Unlock()

	t, ok := w.tabs[tabId]
	if !ok {
		if len(w.tabs) >= w.maxTabs {
			oldest := w.order[0]
			w.order = w.order[1:]
			delete(w.tabs, oldest)
		}
		t = &tabEntry{urls: make(map[string]struct{})}
		w.tabs[tabId] = t
		w.order = append(w.order, tabId)
	}
	t.urls[rawURL] = struct{}{}
	if title != "" {
		t.lastTitle = title
	}
}

func (w *webIndex) titleByTab(tabId string) string {
	w.mu.Lock()
	defer w.mu.Unlock()
	if t, ok := w.tabs[tabId]; ok {
		return t.lastTitle
	}
	return ""
}

func (w *webIndex) titleByURL(rawURL string) string {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, t := range w.tabs {
		if _, ok := t.urls[rawURL]; ok {
			return t.lastTitle
		}
	}
	return ""
}

// 包级单例：上限 64 个 Tab 足够桌面浏览场景。
var defaultIndex = newWebIndex(64)

func addWeb(tabId, rawURL string, body []byte) { defaultIndex.record(tabId, rawURL, body) }
func search(tabId string) string               { return defaultIndex.titleByTab(tabId) }
func SearchTitle(rawURL string) string         { return defaultIndex.titleByURL(rawURL) }
