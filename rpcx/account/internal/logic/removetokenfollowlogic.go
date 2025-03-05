package logic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/simance-ai/smdx/rpcx/account/account"
	"github.com/simance-ai/smdx/rpcx/account/internal/model"
	"github.com/simance-ai/smdx/rpcx/account/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveTokenFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveTokenFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveTokenFollowLogic {
	return &RemoveTokenFollowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 取消关注token
func (l *RemoveTokenFollowLogic) RemoveTokenFollow(in *account.AddTokenFollowRequest) (*account.AddTokenFollowResponse, error) {
	// 检查是否存在关注条目
	existing, err := l.svcCtx.UserTokenFollowModel.FindOneByAddressAndTokenAddress(l.ctx, in.Address, in.TokenAddress)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, fmt.Errorf("no follow entry found for address: %s", in.Address)
		}
		return nil, err
	}

	// 验证链和代币地址是否匹配
	if existing.Chain.String != in.Chain || existing.TokenAddress.String != in.TokenAddress {
		return nil, fmt.Errorf("no matching follow entry found for chain: %s and token address: %s", in.Chain, in.TokenAddress)
	}

	// 更新状态以表明取消关注
	existing.Status = sql.NullString{String: "0", Valid: true}
	// 更新取消关注时间
	existing.UnfollowedAt = sql.NullTime{Time: time.Now(), Valid: true}
	err = l.svcCtx.UserTokenFollowModel.Update(l.ctx, existing)
	if err != nil {
		return nil, err
	}

	// 返回成功响应
	return &account.AddTokenFollowResponse{
		Message: "Token unfollowed successfully",
	}, nil
}
