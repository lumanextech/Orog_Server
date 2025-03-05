package logic

import (
	"context"

	"github.com/simance-ai/smdx/rpcx/account/account"
	"github.com/simance-ai/smdx/rpcx/account/internal/model"
	"github.com/simance-ai/smdx/rpcx/account/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsTokenFollowedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsTokenFollowedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsTokenFollowedLogic {
	return &IsTokenFollowedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户是否收藏代币
func (l *IsTokenFollowedLogic) IsTokenFollowed(in *account.IsTokenFollowedRequest) (*account.IsTokenFollowedResponse, error) {
	// Query the database to check if the token is followed
	_, err := l.svcCtx.UserTokenFollowModel.FindOneByAddressAndTokenAddress(l.ctx, in.Address, in.TokenAddress)
	if err != nil {
		if err == model.ErrNotFound {
			// Token is not followed
			return &account.IsTokenFollowedResponse{Followed: false}, nil
		}
		return nil, err
	}

	// Token is followed
	return &account.IsTokenFollowedResponse{Followed: true}, nil
}
