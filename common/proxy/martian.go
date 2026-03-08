package proxy

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/martian"
	"github.com/google/martian/auth"
	"github.com/google/martian/log"
	"github.com/google/martian/mitm"
	"github.com/google/martian/priority"
	"github.com/google/martian/proxyauth"
)

type VideoTask struct {
	URL       string
	VideoType string
	Header    map[string][]string
	Title     string
	Timestamp time.Time
}

type TaskHandler func(task VideoTask) error

type Proxy struct {
	config      *Config
	martian     *martian.Proxy
	ca          *x509.Certificate
	privateKey  *rsa.PrivateKey
	upstreamURL *url.URL
}

func init() {
	log.SetLevel(log.Silent)
}

func New(cfg *Config) (*Proxy, error) {
	p := &Proxy{config: cfg}

	if cfg.EnableMITM {
		if err := p.loadCert(); err != nil {
			return nil, err
		}
	}

	if cfg.UpstreamProxy != "" {
		u, err := url.Parse(cfg.UpstreamProxy)
		if err != nil {
			return nil, err
		}
		p.upstreamURL = u
	}

	return p, nil
}

func (p *Proxy) Start() error {
	p.martian = martian.NewProxy()

	if p.config.EnableMITM {
		mc, err := mitm.NewConfig(p.ca, p.privateKey)
		if err != nil {
			return err
		}
		p.martian.SetMITM(mc)
	}

	group := priority.NewGroup()
	xs := &skip{handler: p.config.TaskHandler}
	group.AddRequestModifier(xs, 10)
	group.AddResponseModifier(xs, 10)
	xa := &xauth{proxy: p, pAuth: proxyauth.NewModifier()}
	group.AddRequestModifier(xa, 12)
	group.AddResponseModifier(xa, 12)
	p.martian.SetRequestModifier(group)
	p.martian.SetResponseModifier(group)

	fmt.Printf("listen %s, upstream proxy %s\n", p.config.ListenAddr, p.config.UpstreamProxy)
	listener, err := net.Listen("tcp", p.config.ListenAddr)
	if err != nil {
		return err
	}

	return p.martian.Serve(listener)
}

type skip struct {
	handler TaskHandler
}

func (r *skip) ModifyRequest(req *http.Request) error {
	req.Header.Del("Accept-Encoding")

	parts := strings.Split(req.URL.Path, ".")
	if len(parts) == 0 {
		return nil
	}

	ext := parts[len(parts)-1]
	if ext != ExtMP4 && ext != ExtM3U8 {
		return nil
	}

	if r.handler != nil {
		task := VideoTask{
			URL:       req.URL.String(),
			VideoType: ext,
			Header:    req.Header,
			Timestamp: time.Now(),
		}
		_ = r.handler(task)
	}

	return nil
}

func (r *skip) ModifyResponse(res *http.Response) error {
	if !strings.HasSuffix(res.Request.URL.String(), ".html") {
		return nil
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if len(data) > 0 {
		if title, _ := ParseHtmlTitle(bytes.NewBuffer(data)); title != "" {
			// 可以在这里处理 HTML 标题
		}
	}

	res.Body = io.NopCloser(bytes.NewBuffer(data))
	return nil
}

type xauth struct {
	proxy *Proxy
	pAuth *proxyauth.Modifier
}

func (r *xauth) ModifyRequest(req *http.Request) error {
	if r.proxy.config.UpstreamProxy == "" {
		return nil
	}

	r.proxy.martian.SetDownstreamProxy(r.proxy.upstreamURL)

	if r.proxy.config.ProxyUsername != "" {
		un := base64.StdEncoding.EncodeToString([]byte(r.proxy.config.ProxyUsername))
		pw := base64.StdEncoding.EncodeToString([]byte(r.proxy.config.ProxyPassword))
		ctx := martian.NewContext(req)
		authCTX := auth.FromContext(ctx)
		if authCTX.ID() != fmt.Sprintf("%s:%s", un, pw) {
			authCTX.SetError(errors.New("auth error"))
			ctx.SkipRoundTrip()
		}
	}

	return nil
}

func (r *xauth) ModifyResponse(res *http.Response) error {
	if r.proxy.config.UpstreamProxy == "" {
		return nil
	}
	return r.pAuth.ModifyResponse(res)
}
