package logic

import (
	"context"
	"fmt"

	"github.com/simance-ai/smdx/rpcx/account/account"
	"github.com/simance-ai/smdx/rpcx/account/internal/model"
	"github.com/simance-ai/smdx/rpcx/account/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAccountInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAccountInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAccountInfoLogic {
	return &GetAccountInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAccountInfoLogic) GetAccountInfo(in *account.GetAccountInfoRequest) (*account.GetAccountInfoResponse, error) {
	// 删除缓存
	l.svcCtx.AccountModel.DelAccountAddressCache(in.Address)

	// Retrieve account information by address
	accountInfo, err := l.svcCtx.AccountModel.FindOne(context.Background(), in.Address)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, fmt.Errorf("account not found for address: %s", in.Address)
		}
		return nil, err
	}

	inviteAccount, err := l.svcCtx.AccountModel.FindInviteAccount(context.Background(), in.Address)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, fmt.Errorf("account invite not found for address: %s", in.Address)
		}
		return nil, err
	}

	amount, err := l.svcCtx.RebateDetailModel.FindAddressRebateAmount(context.Background(), in.Address)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, fmt.Errorf("account rebate not found for address: %s", in.Address)
		}
		return nil, err
	}

	role, err := l.svcCtx.RebateRoleModel.FindOne(context.Background(), accountInfo.RoleId.Int64)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, fmt.Errorf("account role not found for address: %s", in.Address)
		}
		return nil, err
	}

	// Map the data to the response
	resp := &account.GetAccountInfoResponse{
		Username:           accountInfo.Username.String,
		InitialFunding:     accountInfo.InitialFunding.Float64,
		Funding:            0, //  未填充
		MoneyPerDay:        make([]float64, 0),
		UnrealizedProfits:  0,
		TotalProfit:        0,
		Buy:                0,
		Sell:               0,
		Role:               role.Role,
		InviteCode:         accountInfo.InvitedCode,
		RebateAddress:      inviteAccount,
		RebateAddressTotal: int64(len(inviteAccount)),
		Balance:            accountInfo.Bakance.Float64,
		HistoryRebate:      amount,
		CreateAt:           accountInfo.CreatedAt.Time.Unix(),
	}

	return resp, nil
}
