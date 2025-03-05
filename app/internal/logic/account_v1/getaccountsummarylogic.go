package account_v1

import (
	"context"

	"github.com/simance-ai/smdx/app/internal/svc"
	"github.com/simance-ai/smdx/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAccountSummaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAccountSummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAccountSummaryLogic {
	return &GetAccountSummaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAccountSummaryLogic) GetAccountSummary(req *types.GetAccountSummaryRequest) (resp *types.GetAccountSummaryResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
