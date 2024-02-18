package ws

import (
	"dv/internel/serve/api/internal/logic/ws"
	"dv/internel/serve/api/internal/svc"
	xhttp "github.com/zeromicro/x/http"
	"net/http"
)

func LogHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := ws.NewShowLogLogic(r.Context(), svcCtx)
		if err := l.ShowLogWsRoute(w, r); err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		}
	}
}
