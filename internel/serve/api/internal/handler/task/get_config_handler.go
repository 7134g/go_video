package task

import (
	xhttp "github.com/zeromicro/x/http"
	"net/http"

	"dv/internel/serve/api/internal/logic/task"
	"dv/internel/serve/api/internal/svc"
	"dv/internel/serve/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetConfigRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := task.NewGetConfigLogic(r.Context(), svcCtx)
		resp, err := l.GetConfig(&req)
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
