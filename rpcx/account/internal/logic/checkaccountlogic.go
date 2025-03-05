package logic

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/simance-ai/smdx/rpcx/account/account"
	"github.com/simance-ai/smdx/rpcx/account/internal/model"
	"github.com/simance-ai/smdx/rpcx/account/internal/svc"

	"database/sql"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckAccountLogic {
	return &CheckAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 生成八位邀请码
func GenerateInviteCode() string {
	return uuid.New().String()[:8] // 取前 8 位
}

// 检查有没有Account没有创建
func (l *CheckAccountLogic) CheckAccount(in *account.CheckAccountRequest) (*account.CheckAccountResponse, error) {
	accountData, err := l.svcCtx.AccountModel.FindOne(l.ctx, in.Address)
	if err != nil {
		if err == model.ErrNotFound {
			// Account not found, create a new one
			newAccount := &model.Account{
				Address:        in.Address,
				Chain:          sql.NullString{String: "sol", Valid: true},
				Username:       sql.NullString{String: in.Address, Valid: true},
				CreatedAt:      sql.NullTime{Time: time.Now(), Valid: true},
				UpdatedAt:      sql.NullTime{Time: time.Now(), Valid: true},
				InitialFunding: sql.NullFloat64{Float64: 0, Valid: true},
				Bakance:        sql.NullFloat64{Float64: 0, Valid: true}, // 用户余额 单词错了

				InvitedCode: GenerateInviteCode(),
			}
			_, err = l.svcCtx.AccountModel.Insert(l.ctx, newAccount)
			if err != nil {
				return nil, fmt.Errorf("failed to create account: %v", err)
			}

			return &account.CheckAccountResponse{
				Address: newAccount.Address,
			}, nil
		}
		return nil, err
	}

	return &account.CheckAccountResponse{
		Address: accountData.Address,
	}, nil
}
