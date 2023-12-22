package task

import (
	"context"
	"dv/internel/serve/api/internal/model"
	"dv/internel/serve/api/internal/svc"
	"dv/internel/serve/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RunLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRunLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RunLogic {
	return &RunLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RunLogic) Run(req *types.TaskRunRequest) (resp *types.TaskRunResponse, err error) {
	if req.Stop {
		l.svcCtx.TaskControl.Stop()
		return &types.TaskRunResponse{Message: "正在停止中"}, err
	}

	task := make([]model.Task, 0)
	l.svcCtx.TaskModel.DB.Where("status != ?", model.StatsuSuccess).Find(&task)

	if l.svcCtx.TaskControl.GetStatus() {
		return &types.TaskRunResponse{Message: "正在执行中"}, err
	}
	go l.svcCtx.TaskControl.Run(task)

	return
}
