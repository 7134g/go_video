package task

import (
	"context"
	"dv/internel/serve/api/internal/svc"
	"dv/internel/serve/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShowLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogLogic {
	return &ShowLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShowLogLogic) ShowLog(_ *types.ShowLogRequest) (resp *types.ShowLogResponse, err error) {

	resp = &types.ShowLogResponse{Text: l.svcCtx.LogData.String()}
	return
}
