package logic

import (
	"context"
	"fmt"

	"github.com/simance-ai/smdx/rpcx/account/account"
	"github.com/simance-ai/smdx/rpcx/account/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowTokenListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowTokenListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowTokenListLogic {
	return &GetFollowTokenListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户关注的代币
func (l *GetFollowTokenListLogic) GetFollowTokenList(in *account.GetFollowTokenListRequest) (*account.GetFollowTokenListResponse, error) {
	response := &account.GetFollowTokenListResponse{}

	userTokens, err := l.svcCtx.UserTokenFollowModel.FindAllByAddressAndStatus(l.ctx, in.Address, "1")
	if err != nil {
		return nil, err
	}

	//打印 userTokens
	fmt.Println(userTokens)
	for _, token := range userTokens {
		response.TokenAddress = append(response.TokenAddress, token.TokenAddress.String)
	}
	return response, nil
}
