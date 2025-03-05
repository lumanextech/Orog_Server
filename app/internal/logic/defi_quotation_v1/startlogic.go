package defi_quotation_v1

import (
	"context"
	"github.com/simance-ai/smdx/app/internal/svc"
	"github.com/simance-ai/smdx/app/internal/types"
	"github.com/simance-ai/smdx/pkg/common"
	"github.com/simance-ai/smdx/pkg/errors/api_err"
	"github.com/simance-ai/smdx/rpcx/chains/sol/solclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type StartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartLogic {
	return &StartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StartLogic) Start(req *types.StartRequest) (resp *types.Response, err error) {
	chain := req.Chain
	height := req.Height

	switch chain {
	case common.SolChainId:
		_, err = l.svcCtx.SolClient.Start(l.ctx, &solclient.StartRequest{
			Height: height,
		})
		if err != nil {
			return nil, api_err.NewErrorWithCodeAndMsg(api_err.ErrCodeInternalErrorCode, err.Error())
		}
	default:
		return nil, api_err.ErrCodeInvalidChainNotSupport
	}

	return &types.Response{Message: "success"}, nil
}
