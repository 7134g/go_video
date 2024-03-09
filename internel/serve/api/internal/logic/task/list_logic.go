package task

import (
	"context"
	"dv/internel/serve/api/internal/util/model"
	"dv/internel/serve/api/internal/util/table"
	"github.com/jinzhu/copier"

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

func (l *ListLogic) List(req *types.TaskListRequest) (resp *types.TaskListResponse, err error) {

	total, err := l.svcCtx.TaskModel.Count(req)
	if err != nil {
		return nil, err
	}

	list, err := l.svcCtx.TaskModel.List(req)
	if err != nil {
		return nil, err
	}

	resp = &types.TaskListResponse{
		Total: total,
		List:  nil,
	}
	for _, task := range list {
		var data types.TaskInfo
		_ = copier.Copy(&data, &task)
		if model.StatusSuccess.Eq(task.Status) {
			data.Score = 10000
		} else {
			data.Score, _ = table.DownloadTaskScore.Get(task.ID)
		}
		resp.List = append(resp.List, data)
	}

	return
}
