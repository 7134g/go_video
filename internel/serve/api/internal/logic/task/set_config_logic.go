package task

import (
	"context"
	"github.com/jinzhu/copier"

	"dv/internel/serve/api/internal/svc"
	"dv/internel/serve/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetConfigLogic {
	return &SetConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetConfigLogic) SetConfig(req *types.SetConfigRequest) (resp *types.SetConfigResponse, err error) {

	_ = copier.Copy(l.svcCtx.Config.TaskControlConfig, req)

	return
}
