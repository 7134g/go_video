package logic

import (
	"context"
	"io"
	"os"

	"dv/internel/serve/api/internal/svc"
	"dv/internel/serve/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCertFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	writer io.Writer
}

func NewGetCertFileLogic(ctx context.Context, svcCtx *svc.ServiceContext, writer io.Writer) *GetCertFileLogic {
	return &GetCertFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		writer: writer,
	}
}

func (l *GetCertFileLogic) GetCertFile(req *types.GetCertRequest) (resp *types.GetCertResponse, err error) {
	body, err := os.ReadFile(req.File)
	if err != nil {
		return nil, err
	}

	n, err := l.writer.Write(body)
	if err != nil {
		return nil, err
	}

	if n < len(body) {
		return nil, io.ErrClosedPipe
	}

	return nil, nil
}
