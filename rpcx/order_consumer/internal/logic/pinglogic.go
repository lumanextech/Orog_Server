package logic

import (
	"context"

	"github.com/simance-ai/smdx/rpcx/order_consumer/internal/svc"
	"github.com/simance-ai/smdx/rpcx/order_consumer/order_consumer"

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

func (l *PingLogic) Ping(in *order_consumer.Request) (*order_consumer.Response, error) {
	// todo: add your logic here and delete this line

	return &order_consumer.Response{}, nil
}
