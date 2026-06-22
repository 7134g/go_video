package controller

import (
	"net"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
)

// httpClientHolder 持有一个根据配置可重建的 *http.Client。
// 用 atomic.Pointer 让下载 goroutine 在不加锁的情况下读到最新的 Client。
type httpClientHolder struct {
	cli atomic.Pointer[http.Client]
}

func (h *httpClientHolder) Init() {
	h.cli.Store(http.DefaultClient)
}

// SetProxy 在每次 ApplyConfig 时被调用，重新装配 Transport。
// 传空字符串则不走上游代理。
func (h *httpClientHolder) SetProxy(vpnAddress string) {
	tr := &http.Transport{
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   16,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   15 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}
	if vpnAddress != "" {
		proxyURL, err := url.Parse("http://" + vpnAddress)
		if err == nil {
			tr.Proxy = http.ProxyURL(proxyURL)
		}
	}
	client := &http.Client{
		Transport: tr,
		// 整体超时不设 — 视频段可能很大；用每段 ctx + ResponseHeaderTimeout 兜底。
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return nil },
	}
	h.cli.Store(client)
}

func (h *httpClientHolder) Get() *http.Client {
	return h.cli.Load()
}
