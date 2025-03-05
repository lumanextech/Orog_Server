package logic

import (
	"context"
	"errors"
	"github.com/simance-ai/smdx/rpcx/account/account"
	"github.com/simance-ai/smdx/rpcx/account/internal/model"
	"github.com/simance-ai/smdx/rpcx/account/internal/svc"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type RebateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRebateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RebateLogic {
	return &RebateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RebateLogic) Rebate(in *account.RebateRequest) (*account.Response, error) {
	logx.Info("start RebateLogic")
	ctx := context.Background()
	// 获取返佣角色的返佣比例
	rebateRates := make(map[int64]struct {
		Level1 float64
		Level2 float64
	})
	rebateRoles, err := l.svcCtx.RebateRoleModel.FindMany(ctx)
	if err != nil {
		logx.Info("find role err: ", err.Error())
		return &account.Response{}, err
	}
	for _, role := range rebateRoles {
		rebateRates[role.Id] = struct {
			Level1 float64
			Level2 float64
		}{
			Level1: role.FirstRebateRatio,
			Level2: role.SecondRebateRatio,
		}
	}

	// 获取返佣的一级二级地址
	var userAddress []string
	for _, order := range in.RebateOrderList {
		userAddress = append(userAddress, order.UserAddress)
	}
	if len(userAddress) == 0 {
		return &account.Response{}, errors.New("no user address")
	}
	accountData, err := l.svcCtx.AccountModel.FindRebateAccount(ctx, userAddress)
	if err != nil {
		logx.Info("find rebate account err: ", err.Error())
		return &account.Response{}, err
	}
	if len(userAddress) == 0 {
		return &account.Response{}, errors.New("no user rebate address")
	}

	// 封装所有返佣用户
	rebateAccountMap := make(map[string]model.RebateAccount)
	for _, account := range accountData {
		rebateAccountMap[account.Address] = account
	}

	// 返佣记录表: map[用户地址]累计返佣金额
	rebateRecords := make(map[string]float64)

	for _, order := range in.RebateOrderList {
		rebateAccount := rebateAccountMap[order.UserAddress]
		// 一级返佣邀请人
		inviterLevel1 := rebateAccountMap[rebateAccount.InvitedAddress.String]
		// 二级返佣邀请人
		inviterLevel2, ok := rebateAccountMap[inviterLevel1.InvitedAddress.String]

		// 获取一级返佣比例0
		if rateLevel1, ok := rebateRates[inviterLevel1.RoleId.Int64]; ok {
			rebateAmount1 := order.FilledValue * rateLevel1.Level1
			rebateRecords[inviterLevel1.Address] += rebateAmount1
			logx.Infof("记录一级返佣: 用户 %s, 返佣金额 %.2f", inviterLevel1.Address, rebateAmount1)
		} else {
			logx.Infof("一级邀请人 %s 角色 %s 不在返佣列表", inviterLevel1.Address, inviterLevel1.RoleId)
		}

		// 获取二级返佣比例（如果二级邀请人存在）
		if ok {
			if rateLevel2, ok := rebateRates[inviterLevel2.RoleId.Int64]; ok {
				rebateAmount2 := order.FilledValue * rateLevel2.Level2
				rebateRecords[inviterLevel2.Address] += rebateAmount2
				logx.Infof("记录二级返佣: 用户 %s, 返佣金额 %.2f", inviterLevel2.Address, rebateAmount2)
			} else {
				logx.Infof("二级邀请人 %s 角色 %s 不在返佣列表", inviterLevel2.Address, inviterLevel2.RoleId)
			}
		}
	}

	if len(rebateRecords) == 0 {
		return &account.Response{}, errors.New("no rebate record")
	}

	var data []model.RebateDetail
	for address, amount := range rebateRecords {
		data = append(data, model.RebateDetail{
			Address:      address,
			RebateAmount: amount,
			CreatedAt:    time.Now(), // 当前时间
		})
	}

	// 更新用户钱包
	err = l.svcCtx.AccountModel.UpdateRebateBalance(ctx, rebateRecords)
	if err != nil {
		logx.Errorf("批量更新返佣记录失败: %v", err)
		return &account.Response{}, err
	}

	// 插入返佣记录
	_, err = l.svcCtx.RebateDetailModel.InsertMany(ctx, data)
	if err != nil {
		logx.Errorf("批量插入返佣记录失败: %v", err)
		return &account.Response{}, err
	}

	logx.Info("返佣计算完成")
	return &account.Response{}, nil
}
