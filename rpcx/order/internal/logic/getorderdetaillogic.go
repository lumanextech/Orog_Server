package logic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/simance-ai/smdx/pkg/errors/x_err"
	"github.com/simance-ai/smdx/rpcx/order/internal/svc"
	"github.com/simance-ai/smdx/rpcx/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderDetailLogic {
	return &GetOrderDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Helper function for formatting sql.NullTime
func FormatNullTime(nullTime sql.NullTime) string {
	if nullTime.Valid {
		return nullTime.Time.Format(time.RFC3339)
	}
	return ""
}
func (l *GetOrderDetailLogic) GetOrderDetail(in *order.GetOrderDetailRequest) (*order.GetOrderDetailResponse, error) {
	// Validate request
	if in == nil {
		return nil, x_err.NewErrCodeMsg(x_err.In_ERROR, "Request cannot be empty")
	}
	if in.Id <= 0 {
		return nil, x_err.NewErrCodeMsg(x_err.In_ERROR, "Invalid order ID")
	}

	// Get order from database
	orderInfo, err := l.svcCtx.OrderModel.FindOrderById(l.ctx, in.Id)
	if err != nil {
		return nil, x_err.NewErrCodeMsg(x_err.DBErr, fmt.Sprintf("Failed to get order: %v", err))
	}

	// Convert to response
	return &order.GetOrderDetailResponse{
		Id:              orderInfo.Id,
		OrderHash:       orderInfo.OrderHash,
		Status:          fmt.Sprintf("%d", orderInfo.Status),
		Message:         orderInfo.Message.String,
		ChainName:       orderInfo.ChainId,
		MarketAddress:   orderInfo.MarketAddress,
		Side:            orderInfo.Side,
		Type:            orderInfo.Type,
		Price:           fmt.Sprintf("%f", orderInfo.Price),
		Amount:          fmt.Sprintf("%f", orderInfo.Amount),
		Slippage:        int64(orderInfo.Slippage),
		FilledAmount:    fmt.Sprintf("%f", orderInfo.FilledAmount),
		RemainingAmount: fmt.Sprintf("%f", orderInfo.RemainingAmount),
		CreatedAt:       orderInfo.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       orderInfo.UpdatedAt.Format(time.RFC3339),
		UserId:          orderInfo.AccountAddress,
		PaymentStatus:   fmt.Sprintf("%d", orderInfo.PaymentStatus),
		TransactionHash: orderInfo.TransactionHash.String,
		CancelReason:    orderInfo.CancelReason.String,
	}, nil
}
