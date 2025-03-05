package account_v1

import (
	"context"
	"fmt"

	"github.com/simance-ai/smdx/app/internal/svc"
	"github.com/simance-ai/smdx/app/internal/types"
	"github.com/simance-ai/smdx/rpcx/account/account"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAccountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAccountLogic {
	return &UpdateAccountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 更新用户信息
func (l *UpdateAccountLogic) UpdateAccount(req *types.UpdateAccountRequest) (resp *types.UpdateAccountResponse, err error) {
	// 从上下文拿出address
	address, ok := l.ctx.Value("payload").(string)
	if !ok {
		return nil, fmt.Errorf("failed to get address from context")
	}
	updateBackInfo, err := l.svcCtx.AccountClient.UpdateUserBaseInfo(l.ctx, &account.UpdateUserBaseInfoRequest{
		Address:  address,
		Username: req.Name,
		Chain:    "sol",
	})
	if err != nil {
		return nil, err
	}
	return &types.UpdateAccountResponse{
		Message: updateBackInfo.Message,
	}, nil
}
