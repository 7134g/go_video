package proxy

import (
	"bytes"
	"dv/internel/serve/api/internal/util/model"
	"dv/internel/serve/api/internal/util/table"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/martian"
	"github.com/google/martian/auth"
	"github.com/google/martian/log"
	"github.com/google/martian/mitm"
	"github.com/google/martian/priority"
	"github.com/google/martian/proxyauth"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	monitorAddress = "127.0.0.1:10888" // 监听地址
	taskDB         *model.TaskModel
)

var (
	httpMartian *martian.Proxy // 拦截器全局对象
	certFlag    bool           // 开启自签证书验证
)

var (
	serverProxyUrlParse *url.URL // 解析代理

	serverProxyFlag     bool   // 启用代理
	serverProxy         string // 服务代理地址
	serverProxyUsername string // 用户名
	serverProxyPassword string // 密码
)

func init() {
	log.SetLevel(log.Silent)
}

func OpenCert() {
	certFlag = true
	if err := LoadCert(); err != nil {
		panic(err)
	}
}

func SetServeProxyAddress(address, username, password string) {
	serverProxyFlag = true
	serverProxy = address
	serverProxyUsername = username
	serverProxyPassword = password
}

func SetTaskDb(taskDb *model.TaskModel) {
	taskDB = taskDb
}

func SetMartianAddress(address string) {
	monitorAddress = address
}

func Martian() error {
	httpMartian = martian.NewProxy()
	if certFlag {
		mc, err := mitm.NewConfig(ca, private)
		if err != nil {
			return err
		}
		httpMartian.SetMITM(mc)
	}

	if serverProxyFlag {
		u, err := url.Parse(serverProxy)
		if err != nil {
			return err
		}
		serverProxyUrlParse = u
	}

	group := priority.NewGroup()
	xs := newSkip()
	group.AddRequestModifier(xs, 10)
	group.AddResponseModifier(xs, 10)
	xa := newAuth(proxyauth.NewModifier())
	group.AddRequestModifier(xa, 12)
	group.AddResponseModifier(xa, 12)
	httpMartian.SetRequestModifier(group)
	httpMartian.SetResponseModifier(group)

	fmt.Printf("listen %s, user proxy %s \n", monitorAddress, serverProxy)
	listener, err := net.Listen("tcp", monitorAddress)
	if err != nil {
		return err
	}

	err = httpMartian.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

type skip struct {
}

func newSkip() *skip {
	return &skip{}
}

func (r *skip) ModifyRequest(req *http.Request) error {
	//fmt.Println(req.Method, req.URL.String())
	req.Header.Del("Accept-Encoding")
	parts := strings.Split(req.URL.Path, ".")
	if len(parts) > 0 {
		var header string
		ext := parts[len(parts)-1]
		switch ext {
		case model.VideoTypeMp4, model.VideoTypeM3u8:
			v, _ := json.Marshal(req.Header)
			header = string(v)
		default:
			return nil
		}
		//var name string
		//switch {
		//case req.Header.Get("Origin") != "":
		//	key, _ := getUrlHost(req.Header.Get("Origin"))
		//	value, exist := table.TitleData.Get(key)
		//	if exist {
		//		name = value
		//	}
		//case req.Header.Get("Referer") != "":
		//	key, _ := getUrlHost(req.Header.Get("Referer"))
		//	value, exist := table.TitleData.Get(key)
		//	if exist {
		//		name = value
		//	}
		//}
		findTask, _ := taskDB.Exist(req.URL.String())
		if findTask == nil {
			t := model.Task{
				Name:       fmt.Sprintf("%s", time.Now().Format("2006_01_02_15_04_05")),
				VideoType:  ext,
				Type:       model.TypeProxy,
				Data:       req.URL.String(),
				HeaderJson: header,
			}
			if err := taskDB.Insert(&t); err != nil {
				logx.Error(err)
			}

		}

	}

	return nil
}

func (r *skip) ModifyResponse(res *http.Response) error {
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}

	buf := bytes.NewBuffer(data)
	logx.Debugw(
		"url message",
		logx.Field("method", res.Request.Method),
		logx.Field("url", res.Request.URL.String()))
	title, err := ParseHtmlTitle(buf)
	if err != nil {
		logx.Error(err)
	}
	if title != "" {
		table.TitleData.Set(res.Request.URL.Host, title)
	}

	res.Body = io.NopCloser(bytes.NewBuffer(data))
	return nil
}

type xauth struct {
	pAuth *proxyauth.Modifier
}

func newAuth(pAuth *proxyauth.Modifier) *xauth {
	return &xauth{pAuth: pAuth}
}

func (r *xauth) ModifyRequest(req *http.Request) error {
	if serverProxy == "" {
		return nil
	}

	httpMartian.SetDownstreamProxy(serverProxyUrlParse)

	if serverProxyUsername != "" {
		un := base64.StdEncoding.EncodeToString([]byte(serverProxyUsername))
		pw := base64.StdEncoding.EncodeToString([]byte(serverProxyPassword))
		//req.Header.Set("Proxy-Authorization", fmt.Sprintf("Basic %s:%s", un, pw))
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
	if serverProxy == "" {
		return nil
	}
	return r.pAuth.ModifyResponse(res)
}

func getNumber() int64 {
	currentTime := time.Now().Unix()
	return currentTime - currentTime%10
}

func getUrlHost(u string) (string, error) {
	_url, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	return _url.Host, nil
}
