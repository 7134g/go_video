package task

import (
	"context"

	"dv/internel/serve/api/internal/svc"
	"dv/internel/serve/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(_ *types.TaskListRequest) (resp *types.TaskListResponse, err error) {

	list, err := l.svcCtx.TaskModel.List()
	if err != nil {
		return nil, err
	}

	resp = &types.TaskListResponse{
		List: list,
	}

	return
}
