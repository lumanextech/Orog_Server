package logic

import (
	"context"
	"fmt"
	"github.com/simance-ai/smdx/pkg/errors/x_err"

	"github.com/simance-ai/smdx/rpcx/order/internal/svc"
	"github.com/simance-ai/smdx/rpcx/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAccountOrderInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAccountOrderInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAccountOrderInfoLogic {
	return &GetAccountOrderInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAccountOrderInfoLogic) GetAccountOrderInfo(in *order.GetAccountInfoRequest) (*order.GetAccountInfoResponse, error) {

	buy, sell, err := l.svcCtx.OrderModel.FindOrderTradeNumber(context.Background(), in.Address)
	if err != nil {
		return nil, x_err.NewErrCodeMsg(x_err.DBErr, fmt.Sprintf("Failed to get account trade total: %v", err))
	}

	return &order.GetAccountInfoResponse{
		Funding:           0,
		MoneyPerDay:       make([]float64, 0),
		UnrealizedProfits: 0,
		TotalProfit:       0,
		Buy:               buy,
		Sell:              sell,
	}, nil
}
