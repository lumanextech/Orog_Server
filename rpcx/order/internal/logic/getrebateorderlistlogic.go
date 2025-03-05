package logic

import (
	"context"
	"fmt"
	"github.com/simance-ai/smdx/rpcx/order/internal/svc"
	"github.com/simance-ai/smdx/rpcx/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRebateOrderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRebateOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRebateOrderListLogic {
	return &GetRebateOrderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRebateOrderListLogic) GetRebateOrderList(in *order.RebateOrderListRequest) (*order.RebateOrderListResponse, error) {
	orders, err := l.svcCtx.OrderModel.FindRebateOrders(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get order list: %w", err)
	}

	// Construct the response
	orderList := make([]*order.RebateOrder, len(orders))
	for i, o := range orders {
		orderList[i] = &order.RebateOrder{
			Id:              o.Id,
			OrderHash:       o.OrderHash,
			ChainName:       o.ChainId,
			MarketAddress:   o.MarketAddress,
			UserAddress:     o.AccountAddress,
			FilledValue:     o.Amount * o.Price,
			TransactionHash: o.TransactionHash.String,
		}
	}

	return &order.RebateOrderListResponse{OrderList: orderList}, nil
}
