package defi_quotation_v1

import (
	"context"

	"github.com/simance-ai/smdx/app/internal/svc"
	"github.com/simance-ai/smdx/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetXTopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetXTopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetXTopLogic {
	return &GetXTopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetXTopLogic) GetXTop(req *types.GetXTopRequest) (resp *types.GetXTopResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
