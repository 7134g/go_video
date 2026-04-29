package proxy

import (
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/google/martian"
	"github.com/google/martian/mitm"
)

type Server struct {
	proxy     *martian.Proxy
	collector *Collector
	detector  *VideoDetector
	capture   *RequestCapture
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

	proxyUrl, err := url.Parse(fmt.Sprintf("http://%s", proxyAddress))
	if err != nil {
		return nil, err
	}
	agent.SetDownstreamProxy(proxyUrl)

	s := &Server{
		proxy:     agent,
		collector: NewCollector(),
		detector:  &VideoDetector{},
		capture:   &RequestCapture{},
	}
	s.Stop = make(chan bool)
	agent.SetRequestModifier(s)
	return s, nil
}

func (s *Server) ModifyRequest(req *http.Request) error {
	//fmt.Println(req.URL.String())
	if videoType, ok := s.detector.GetVideo(req.URL.String()); ok {
		task := s.capture.Capture(req)
		task.Type = videoType
		s.collector.Collect(task)
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
