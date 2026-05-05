package proxy

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"go_video/pkg/m3u8"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/martian"
	"github.com/google/martian/auth"
	"github.com/google/martian/mitm"
	"github.com/klauspost/compress/zstd"
)

type Server struct {
	proxy     *martian.Proxy
	collector *Collector
	listener  net.Listener

	Stop chan bool
}

func NewServer(vpnAddress string) (*Server, error) {
	ca, err := LoadCA()
	if err != nil {
		return nil, err
	}
	key, err := LoadKey()
	if err != nil {
		return nil, err
	}

	agent := martian.NewProxy()
	mc, err := mitm.NewConfig(ca, key)
	if err != nil {
		return nil, err
	}
	agent.SetMITM(mc)

	if vpnAddress != "" {
		address := fmt.Sprintf("http://%s", vpnAddress)
		fmt.Println("被动代理启动vpn: " + address)
		proxyUrl, err := url.Parse(address)
		if err != nil {
			return nil, err
		}
		agent.SetDownstreamProxy(proxyUrl)
	}

	s := &Server{
		proxy:     agent,
		collector: NewCollector(),
	}
	s.Stop = make(chan bool)
	agent.SetRequestModifier(s)
	agent.SetResponseModifier(s)
	return s, nil
}

func (s *Server) ModifyRequest(req *http.Request) error {
	return nil
}

func (s *Server) ModifyResponse(res *http.Response) error {
	ctx := martian.NewContext(res.Request)
	actx := auth.FromContext(ctx)
	tabId := res.Request.Header.Get("X-Tab-Id")

	u := res.Request.URL
	if strings.Contains(u.Host, "localhost") || strings.Contains(u.Host, "127.0.0.1") {
		return nil
	}

	var body []byte
	if res.Body != nil {
		body, _ = io.ReadAll(res.Body)
		res.Body = io.NopCloser(bytes.NewReader(body))
	}

	encoding := res.Header.Get("Content-Encoding")
	switch {
	case strings.Contains(res.Request.URL.String(), ".css"):
	case strings.Contains(res.Request.URL.String(), ".js"):
	case strings.Contains(res.Request.Header.Get("Sec-Fetch-Dest"), "image"):
	case strings.Contains(res.Request.Header.Get("Content-Type"), "image"):
	case strings.Contains(res.Request.Header.Get("Content-Type"), "jpeg"):
	case strings.Contains(encoding, "gzip") && len(body) > 0:
		reader, err := gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			return err
		}
		defer reader.Close()
		body, _ = io.ReadAll(reader)
		addWeb(tabId, u.String(), body, res.Request.Header)
	case strings.Contains(encoding, "zstd") && len(body) > 0:
		reader, err := zstd.NewReader(bytes.NewReader(body))
		if err != nil {
			return err
		}
		body, _ = io.ReadAll(reader)
		addWeb(tabId, u.String(), body, res.Request.Header)
	default:
		addWeb(tabId, u.String(), body, res.Request.Header)
	}

	var isVideo bool
	// 判断url类型
	videoType, ok := GetVideo(u)
	if ok {
		switch videoType {
		case "mp4":
			isVideo = true
		case "m3u8":
			_, err := m3u8.ParseM3u8Data(bytes.NewReader(body))
			if err != nil {
				fmt.Println("解析失败: ", u.String(), err, string(body))
			} else {
				isVideo = true
			}
		}
	}

	// 判断是否是视频请求
	var isVideoUrl bool
	if HasExactlyOneHttp(u.String()) {
		isVideoUrl = true
	}

	if isVideo && isVideoUrl {
		title := search(tabId)
		s.collector.Collect(res.Request, title, videoType)
	}

	if actx.Error() != nil {
		res.StatusCode = 403
		res.Status = http.StatusText(403)
	}

	return nil
}

func (s *Server) Listen(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.listener = l
	return s.proxy.Serve(l)
}

func (s *Server) Close() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	s.Stop <- true
	return nil
}

func (s *Server) Tasks() <-chan *VideoTask {
	return s.collector.Tasks()
}
