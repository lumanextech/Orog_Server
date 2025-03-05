package logic

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/simance-ai/smdx/rpcx/account/account"
	"github.com/simance-ai/smdx/rpcx/account/internal/model"
	"github.com/simance-ai/smdx/rpcx/account/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserBaseInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserBaseInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserBaseInfoLogic {
	return &UpdateUserBaseInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新用户基础信息
func (l *UpdateUserBaseInfoLogic) UpdateUserBaseInfo(in *account.UpdateUserBaseInfoRequest) (*account.UpdateUserBaseInfoResponse, error) {
	// 获取用户记录
	user, err := l.svcCtx.AccountModel.FindOne(l.ctx, in.Address)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, fmt.Errorf("user not found: %s", in.Address)
		}
		return nil, err
	}

	// 更新用户信息
	user.Username = sql.NullString{String: in.Username, Valid: true}

	// 保存更改
	err = l.svcCtx.AccountModel.Update(l.ctx, user)
	if err != nil {
		return nil, err
	}

	// 返回成功响应
	return &account.UpdateUserBaseInfoResponse{
		Message: "User information updated successfully",
	}, nil
}
