package proxy

import (
	"bytes"
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
)

type Server struct {
	proxy     *martian.Proxy
	collector *Collector
	listener  net.Listener

	Stop chan bool
}

func NewServer(proxyAddress string) (*Server, error) {
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

	if proxyAddress != "" {
		address := fmt.Sprintf("http://%s", proxyAddress)
		fmt.Println("local proxy on :" + address)
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
	//fmt.Println("收到请求:", req.URL.String(), "TabID:", req.Header.Get("X-Tab-Id"))
	//if videoType, ok := s.detector.GetVideo(req.URL.String()); ok {
	//	task := s.capture.Capture(req)
	//	task.Type = videoType
	//	s.collector.Collect(task)
	//}
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

	if isHtml(res.Request) {
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
