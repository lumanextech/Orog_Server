package account_v1

import (
	"context"
	"github.com/simance-ai/smdx/app/internal/svc"
	"github.com/simance-ai/smdx/app/internal/types"
	"github.com/simance-ai/smdx/rpcx/account/account"
	"github.com/simance-ai/smdx/rpcx/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAccountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAccountLogic {
	return &GetAccountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAccountLogic) GetAccount(req *types.GetAccountRequest) (resp *types.GetAccountResponse, err error) {
	accountInfo, err := l.svcCtx.AccountClient.GetAccountInfo(context.Background(), &account.GetAccountInfoRequest{
		Address: req.Address,
	})
	if err != nil {
		return nil, err
	}

	accountOrderInfo, err := l.svcCtx.OrderClient.GetAccountOrderInfo(context.Background(), &order.GetAccountInfoRequest{
		Address: req.Address,
	})
	if err != nil {
		return nil, err
	}

	// 创建响应
	resp = &types.GetAccountResponse{
		Nickname:          accountInfo.Username,          //昵称
		InitialFunding:    accountInfo.InitialFunding,    //初始资金
		Funding:           accountInfo.Balance,           //资金
		MoneyPerDay:       accountInfo.MoneyPerDay,       //资金每日变化数组
		UnrealizedProfits: accountInfo.UnrealizedProfits, //未实现盈亏
		TotalProfit:       accountInfo.TotalProfit,       //总盈亏
		Buy:               accountOrderInfo.Buy,          //买入次数
		Sell:              accountOrderInfo.Sell,         //卖出次数
		Role:              accountInfo.Role,
	}
	return resp, nil
}
