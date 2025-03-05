package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/simance-ai/smdx/rpcx/order/internal/svc"
	"github.com/simance-ai/smdx/rpcx/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderListLogic {
	return &GetOrderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOrderListLogic) GetOrderList(in *order.GetOrderListRequest) (*order.OrderListResponse, error) {
	// Extract request parameters
	pageSize := in.PageSize
	page := in.Page
	userId := in.UserId
	status := in.Status
	chainName := in.ChainName

	// Query the database or data store for the order list
	orders, totalCount, err := l.svcCtx.OrderModel.FindOrders(l.ctx, pageSize, page, userId, status, chainName)
	if err != nil {
		return nil, fmt.Errorf("failed to get order list: %w", err)
	}

	// Construct the response
	orderList := make([]*order.Orders, len(orders))
	for i, o := range orders {
		orderList[i] = &order.Orders{
			Id:              o.Id,
			OrderHash:       o.OrderHash,
			Status:          fmt.Sprintf("%d", o.Status),
			Message:         o.Message.String,
			ChainName:       o.ChainId,
			MarketAddress:   o.MarketAddress,
			Side:            o.Side,
			Type:            o.Type,
			Price:           fmt.Sprintf("%f", o.Price),
			Amount:          fmt.Sprintf("%f", o.Amount),
			Slippage:        int64(o.Slippage),
			FilledAmount:    fmt.Sprintf("%f", o.FilledAmount),
			RemainingAmount: fmt.Sprintf("%f", o.RemainingAmount),
			CreatedAt:       o.CreatedAt.Format(time.RFC3339),
			UpdatedAt:       o.UpdatedAt.Format(time.RFC3339),
			UserId:          o.AccountAddress,
			PaymentStatus:   fmt.Sprintf("%d", o.PaymentStatus),
			TransactionHash: o.TransactionHash.String,
			CancelReason:    o.CancelReason.String,
		}
	}

	response := &order.OrderListResponse{
		OrderList:   orderList,
		TotalCount:  totalCount,
		PageSize:    pageSize,
		CurrentPage: page,
	}

	// Return the response
	return response, nil
}
