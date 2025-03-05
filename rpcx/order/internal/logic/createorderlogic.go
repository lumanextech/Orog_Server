package logic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/simance-ai/smdx/pkg/errors/x_err"
	"github.com/simance-ai/smdx/rpcx/order/internal/model"
	"github.com/simance-ai/smdx/rpcx/order/internal/svc"
	"github.com/simance-ai/smdx/rpcx/order/order"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func generateUniqueOrderHash() string {
	// 使用 UUID 生成唯一 OrderHash
	return uuid.New().String()
}

func (l *CreateOrderLogic) CreateOrder(in *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	// Validate request data
	if in == nil {
		return nil, x_err.NewErrCodeMsg(x_err.In_ERROR, "Request data cannot be empty")
	}

	var newOrder *model.Order
	switch in.Type {
	case 0:
		newOrder = &model.Order{
			Status:          1,
			Message:         sql.NullString{String: "Created successfully", Valid: true},
			OrderHash:       generateUniqueOrderHash(),
			ChainId:         in.ChainName,
			MarketAddress:   in.MarketAddress,
			Side:            in.Side,
			Type:            in.Type,
			Price:           in.Price,
			Amount:          in.Amount,
			Slippage:        float64(in.Slippage),
			FilledAmount:    0,
			RemainingAmount: in.Amount,
			CreatedAt:       time.Now().UTC(),
			UpdatedAt:       time.Now().UTC(),
			AccountAddress:  in.Address,
			PaymentStatus:   0,
			TransactionHash: sql.NullString{},
			CancelReason:    sql.NullString{},
			OpenMev:         false, //只有限价单需要
			LimitOrderType:  0,     // 限价单
			RebateStatus:    0,
		}
	case 1: // Market order
		newOrder = &model.Order{
			Status:          1,
			Message:         sql.NullString{String: "Created successfully", Valid: true},
			OrderHash:       generateUniqueOrderHash(),
			ChainId:         in.ChainName,
			MarketAddress:   in.MarketAddress,
			Side:            in.Side,
			Type:            in.Type,
			Price:           in.Price,
			Amount:          in.Amount,
			Slippage:        float64(in.Slippage),
			FilledAmount:    0,
			RemainingAmount: in.Amount,
			CreatedAt:       time.Now().UTC(),
			UpdatedAt:       time.Now().UTC(),
			AccountAddress:  in.Address,
			PaymentStatus:   0,
			TransactionHash: sql.NullString{String: in.TransactionHash, Valid: true},
			CancelReason:    sql.NullString{},
			OpenMev:         false,
			LimitOrderType:  0,
			RebateStatus:    0,
		}
	}

	// 写入数据库
	_, err := l.svcCtx.OrderModel.Insert(l.ctx, newOrder)
	if err != nil {
		return nil, x_err.NewErrCodeMsg(x_err.DBErr, err.Error())
	}
	// 返回成功响应
	return &order.CreateOrderResponse{
		OrderHash: newOrder.OrderHash,
		Status:    fmt.Sprintf("%d", newOrder.Status),
		Message:   newOrder.Message.String,
	}, nil

}
