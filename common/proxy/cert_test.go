package proxy

import (
	"testing"
	"time"
)

func TestGenMITM(t *testing.T) {
	if err := GenMITM(DefaultCAName, DefaultPrivateKeyName); err != nil {
		t.Fatal(err)
	}
}

func TestNewCert(t *testing.T) {
	cd, pd, err := newCert("proxy", "location", time.Hour*24*365*10)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(cd))
	t.Log(string(pd))
}

func TestLoadCert(t *testing.T) {
	cfg := NewConfig()
	cfg.EnableMITM = true

	p, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if p.ca == nil || p.privateKey == nil {
		t.Fatal("cert not loaded")
	}
}
