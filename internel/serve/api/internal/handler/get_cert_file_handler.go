package handler

import (
	"net/http"

	"dv/internel/serve/api/internal/logic"
	"dv/internel/serve/api/internal/svc"
	"dv/internel/serve/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetCertFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetCertRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetCertFileLogic(r.Context(), svcCtx, w)
		_, err := l.GetCertFile(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("Content-Disposition", "attachment;filename=mitm.crt")
			httpx.Ok(w)
		}
	}
}
