package proxy

import (
	"net"
	"net/http"

	"github.com/google/martian"
	"github.com/google/martian/mitm"
)

type Server struct {
	proxy     *martian.Proxy
	collector *Collector
	detector  *VideoDetector
	capture   *RequestCapture
}

func NewServer(caFile, keyFile string) (*Server, error) {
	ca, err := LoadCA(caFile)
	if err != nil {
		return nil, err
	}
	key, err := LoadKey(keyFile)
	if err != nil {
		return nil, err
	}

	proxy := martian.NewProxy()
	mc, err := mitm.NewConfig(ca, key)
	if err != nil {
		return nil, err
	}
	proxy.SetMITM(mc)

	s := &Server{
		proxy:     proxy,
		collector: NewCollector(),
		detector:  &VideoDetector{},
		capture:   &RequestCapture{},
	}

	proxy.SetRequestModifier(s)
	return s, nil
}

func (s *Server) ModifyRequest(req *http.Request) error {
	if s.detector.IsVideo(req.URL.String()) {
		task := s.capture.Capture(req)
		s.collector.Collect(task)
	}
	return nil
}

func (s *Server) Listen(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.proxy.Serve(l)
}

func (s *Server) Tasks() <-chan *VideoTask {
	return s.collector.Tasks()
}
