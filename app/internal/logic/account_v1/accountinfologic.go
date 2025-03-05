package account_v1

import (
	"context"
	"fmt"
	"github.com/simance-ai/smdx/app/internal/svc"
	"github.com/simance-ai/smdx/app/internal/types"
	"github.com/simance-ai/smdx/rpcx/account/account"
	"github.com/simance-ai/smdx/rpcx/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccountInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 通过token获取个人信息
func NewAccountInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountInfoLogic {
	return &AccountInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AccountInfoLogic) AccountInfo(req *types.AccountInfoRequest) (resp *types.AccountInfoResponse, err error) {
	address, ok := l.ctx.Value("payload").(string)
	if !ok {
		return nil, fmt.Errorf("failed to get address from context")
	}

	accountInfo, err := l.svcCtx.AccountClient.GetAccountInfo(l.ctx, &account.GetAccountInfoRequest{
		Address: address,
	})
	if err != nil {
		return nil, err
	}

	accountOrderInfo, err := l.svcCtx.OrderClient.GetAccountOrderInfo(l.ctx, &order.GetAccountInfoRequest{
		Address: address,
	})
	if err != nil {
		return nil, err
	}

	// 创建响应
	resp = &types.AccountInfoResponse{
		Nickname:           accountInfo.Username,          //昵称
		InitialFunding:     accountInfo.InitialFunding,    //初始资金
		Funding:            accountInfo.Balance,           //资金
		MoneyPerDay:        accountInfo.MoneyPerDay,       //资金每日变化数组
		UnrealizedProfits:  accountInfo.UnrealizedProfits, //未实现盈亏
		TotalProfit:        accountInfo.TotalProfit,       //总盈亏
		Buy:                accountOrderInfo.Buy,          //买入次数
		Sell:               accountOrderInfo.Sell,         //卖出次数
		Role:               accountInfo.Role,
		InviteCode:         accountInfo.InviteCode,
		RebateAddress:      accountInfo.RebateAddress,
		RebateAddressTotal: accountInfo.RebateAddressTotal,
		Balance:            accountInfo.Balance,
		HistoryRebate:      accountInfo.HistoryRebate,
	}
	return resp, nil
}
