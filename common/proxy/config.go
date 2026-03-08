package proxy

import (
	"regexp"
	"time"
)

const (
	DefaultDomain         = "proxy"
	DefaultOrganization   = "location"
	DefaultValidDuration  = time.Hour * 24 * 365 * 20
	DefaultCAName         = "mitm.crt"
	DefaultPrivateKeyName = "pri.pem"
	DefaultMonitorAddress = "127.0.0.1:10888"
	ExtMP4                = "mp4"
	ExtM3U8               = "m3u8"
)

var regUrl = regexp.MustCompile(`([^/]+)(\.m3u8|\.mp4)$`)

type Config struct {
	ListenAddr     string
	EnableMITM     bool
	CAFile         string
	PrivateKeyFile string
	UpstreamProxy  string
	ProxyUsername  string
	ProxyPassword  string
	TaskHandler    TaskHandler
}

func NewConfig() *Config {
	return &Config{
		ListenAddr:     DefaultMonitorAddress,
		CAFile:         DefaultCAName,
		PrivateKeyFile: DefaultPrivateKeyName,
	}
}
