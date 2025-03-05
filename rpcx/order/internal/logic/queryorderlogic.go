package logic

import (
	"context"
	"fmt"

	"github.com/simance-ai/smdx/rpcx/order/internal/svc"
	"github.com/simance-ai/smdx/rpcx/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryOrderLogic {
	return &QueryOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryOrderLogic) QueryOrder(in *order.QueryOrderRequest) (*order.QueryOrderResponse, error) {
	// Query the database or data store for the order details
	orderDetails, err := l.svcCtx.OrderModel.FindOrderById(l.ctx, in.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to query order: %w", err)
	}

	// Construct the response
	response := &order.QueryOrderResponse{
		Id:     orderDetails.Id,
		Status: orderDetails.Status,
	}
	// Return the response
	return response, nil
}
