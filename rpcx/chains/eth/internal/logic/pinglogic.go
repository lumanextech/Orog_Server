package logic

import (
	"context"

	"github.com/simance-ai/smdx/rpcx/chains/eth/eth"
	"github.com/simance-ai/smdx/rpcx/chains/eth/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *eth.Request) (*eth.Response, error) {
	// todo: add your logic here and delete this line

	return &eth.Response{}, nil
}
