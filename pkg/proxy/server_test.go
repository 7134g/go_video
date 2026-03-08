package proxy

import "testing"

func TestGenCA(t *testing.T) {
	GenCA()
}

func TestServer(t *testing.T) {
	s, err := NewServer()
	if err != nil {
		panic(err)
	}

	//if err := InstallCert(); err != nil {
	//	panic(err)
	//}
	_ = s.Listen("127.0.0.1:8888")
}
