package proxy

import (
	"strings"
	"sync"
	"testing"
)

func TestWebIndex_Concurrent(t *testing.T) {
	// 重现旧代码会触发 "fatal: concurrent map writes" 的场景。
	idx := newWebIndex(8)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tabId := "tab" + string(rune('A'+i%4))
			url := "https://example.com/" + string(rune('a'+i%26))
			body := []byte("<html><title>title-" + string(rune('A'+i%26)) + "</title></html>")
			idx.record(tabId, url, body)
			_ = idx.titleByTab(tabId)
			_ = idx.titleByURL(url)
		}(i)
	}
	wg.Wait()
}

func TestWebIndex_LRUEviction(t *testing.T) {
	idx := newWebIndex(2)
	idx.record("a", "u1", []byte("<title>A</title>"))
	idx.record("b", "u2", []byte("<title>B</title>"))
	idx.record("c", "u3", []byte("<title>C</title>")) // 应淘汰 a

	if got := idx.titleByTab("a"); got != "" {
		t.Errorf("tab a should be evicted, got %q", got)
	}
	if got := idx.titleByTab("c"); got != "C" {
		t.Errorf("tab c title = %q, want C", got)
	}
}

func TestWebIndex_OnlyStoresTitle(t *testing.T) {
	idx := newWebIndex(4)
	bigBody := strings.Repeat("x", 1<<20) + "<title>real</title>"
	idx.record("t", "u", []byte(bigBody))
	if got := idx.titleByTab("t"); got != "real" {
		t.Errorf("title = %q, want real", got)
	}
	// 间接验证：tabEntry 只存了 lastTitle / urls，body 没被引用
	idx.mu.Lock()
	entry := idx.tabs["t"]
	idx.mu.Unlock()
	if len(entry.lastTitle) >= 1<<20 {
		t.Errorf("body leaked into index: %d bytes", len(entry.lastTitle))
	}
}

func TestHasExactlyOneHttp(t *testing.T) {
	cases := map[string]bool{
		"https://a.com/x":               true,
		"http://a.com/?x=https://b.com": false, // 嵌套 URL
		"https://a.com/x?y=foo":         true,
		"no-protocol":                   false,
		"https://a.com#https://b.com":   false,
	}
	for in, want := range cases {
		if got := HasExactlyOneHttp(in); got != want {
			t.Errorf("HasExactlyOneHttp(%q) = %v, want %v", in, got, want)
		}
	}
}
