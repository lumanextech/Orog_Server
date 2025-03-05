package defi_quotation_v1

import (
	"context"

	"github.com/simance-ai/smdx/app/internal/svc"
	"github.com/simance-ai/smdx/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMarketInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMarketInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMarketInfoLogic {
	return &GetMarketInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMarketInfoLogic) GetMarketInfo(req *types.GetMemeMarketRequest) (resp *types.GetMemeMarketResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
