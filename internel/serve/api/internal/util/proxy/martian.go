package proxy

// 使用官方包 github.com/google/martian 实现拦截
import (
	"bytes"
	"dv/internel/serve/api/internal/util/model"
	"encoding/json"
	"fmt"
	"github.com/google/martian"
	"github.com/google/martian/priority"
	"github.com/zeromicro/go-zero/core/logx"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	Proxy       *martian.Proxy
	ProxyServer string
	//UserName    string
	//PassWord    string
)

func Martian(taskDB *model.TaskModel) error {
	Proxy = martian.NewProxy()
	group := priority.NewGroup()

	//if UserName != "" {
	//	a := proxyauth.NewModifier()
	//	group.AddRequestModifier(a, 2)
	//	group.AddResponseModifier(a, 2)
	//}

	s := &Skip{taskDB: taskDB}
	group.AddRequestModifier(s, 1)
	group.AddResponseModifier(s, 1)

	Proxy.SetRequestModifier(group)
	Proxy.SetResponseModifier(group)
	// 使用代理发请求时候装载证书
	if ProxyServer != "" {
		fmt.Println("开启代理：", ProxyServer)
		CertReload()
		mc, err := GetMITMConfig()
		if err != nil {
			return err
		}
		Proxy.SetMITM(mc)
	}

	//log.SetLevel(log.Silent)
	listener, err := net.Listen("tcp", ":1080")
	if err != nil {
		return err
	}

	fmt.Println(listener.Addr().String())
	err = Proxy.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

type Skip struct {
	taskDB *model.TaskModel
}

func (r *Skip) ModifyRequest(req *http.Request) error {
	if ProxyServer != "" {
		u, err := url.Parse(ProxyServer)
		if err != nil {
			return err
		}

		Proxy.SetDownstreamProxy(u)
	}

	parts := strings.Split(req.URL.Path, ".")
	if len(parts) > 0 {
		var data string
		ext := parts[len(parts)-1]
		switch ext {
		case model.VideoTypeMp4, model.VideoTypeM3u8:
			v, _ := json.Marshal(req)
			data = string(v)
		default:
			return nil
		}

		findTask, _ := r.taskDB.Exist(data)
		if findTask == nil {
			t := model.Task{
				Name:      fmt.Sprintf("%d", time.Now().UnixMilli()),
				VideoType: ext,
				Type:      model.TypeProxy,
				Data:      data,
			}
			if err := r.taskDB.Insert(&t); err != nil {
				logx.Error(err)
			}
		}

	}

	//ctx := martian.NewContext(req)
	//authCTX := auth.FromContext(ctx)
	//if authCTX.ID() != fmt.Sprintf("%s:%s", UserName, PassWord) {
	//	authCTX.SetError(errors.New("auth error"))
	//	ctx.SkipRoundTrip()
	//}

	return nil
}

func (r *Skip) ModifyResponse(_ *http.Response) error {
	return nil
}

// ExtractRequestToString 提取请求包
func ExtractRequestToString(res *http.Request) string {
	buf := bytes.NewBuffer([]byte{})
	defer buf.Reset()
	err := res.Write(buf)
	if err != nil {
		return ""
	}

	return buf.String()
}
