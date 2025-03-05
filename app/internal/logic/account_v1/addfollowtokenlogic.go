package account_v1

import (
	"context"
	"fmt"

	"github.com/simance-ai/smdx/app/internal/svc"
	"github.com/simance-ai/smdx/app/internal/types"
	"github.com/simance-ai/smdx/rpcx/account/account"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFollowTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddFollowTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFollowTokenLogic {
	return &AddFollowTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddFollowTokenLogic) AddFollowToken(req *types.AddFollowTokenRequest) (resp *types.AddFollowTokenResponse, err error) {
	// 从上下文拿出address
	address, ok := l.ctx.Value("payload").(string)
	if !ok {
		return nil, fmt.Errorf("failed to get address from context")
	}

	// 查询该地址是否已经创建过account，没有则创建account
	rpcResp, err := l.svcCtx.AccountClient.AddTokenFollow(l.ctx, &account.AddTokenFollowRequest{
		Address:      address,
		Chain:        req.Chain,
		TokenAddress: req.Token,
	})
	if err != nil {
		return nil, err
	}
	// 创建响应
	resp = &types.AddFollowTokenResponse{
		Message: rpcResp.Message,
	}
	return resp, nil
}
