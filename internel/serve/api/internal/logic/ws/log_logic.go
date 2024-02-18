package ws

import (
	"context"
	"dv/internel/serve/api/internal/svc"
	"dv/internel/serve/api/internal/util/ws_conn"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"net/http"
)

type ShowLog struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLog {
	return &ShowLog{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShowLog) ShowLogWsRoute(w http.ResponseWriter, r *http.Request) error {
	client, err := ws_conn.InitClient(l.svcCtx.Hub, w, r)
	if err != nil {
		return err
	}

	client.SetReadHandle(func(message []byte) []byte {
		switch string(message) {
		case "get":
			return l.svcCtx.LogData.Bytes()
		}
		return nil
	})

	threading.GoSafe(func() {
		for {
			select {
			case <-client.Ctx.Done():
				return
			case data := <-l.svcCtx.LogData.Cache:
				client.Write(data)
			}
		}
	})

	ws_conn.Run(client)
	return nil
}
