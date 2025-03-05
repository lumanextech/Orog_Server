package logic

import (
	"context"

	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/svc"
	"github.com/simance-ai/smdx/rpcx/chains/sol/sol"

	"github.com/zeromicro/go-zero/core/logx"
)

type StartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartLogic {
	return &StartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StartLogic) Start(in *sol.StartRequest) (*sol.Response, error) {
	// todo: add your logic here and delete this line

	return &sol.Response{}, nil
}
