package logic

import (
	"context"
	"fmt"

	"github.com/simance-ai/smdx/rpcx/order/internal/svc"
	"github.com/simance-ai/smdx/rpcx/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRebateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRebateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRebateOrderLogic {
	return &UpdateRebateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateRebateOrderLogic) UpdateRebateOrder(in *order.UpdateRebateOrderRequest) (*order.Response, error) {
	// Query the database or data store for the order list
	err := l.svcCtx.OrderModel.UpdateRebateOrders(context.Background(), in.OrderIdList)
	if err != nil {
		return nil, fmt.Errorf("failed to get order list: %w", err)
	}
	return &order.Response{
		Pong: "success",
	}, nil
}
