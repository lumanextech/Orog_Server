package defi_quotation_v1

import (
	"context"
	"github.com/simance-ai/smdx/pkg/errors/api_err"

	"github.com/simance-ai/smdx/pkg/common"
	"github.com/simance-ai/smdx/rpcx/chains/sol/solclient"

	"github.com/simance-ai/smdx/app/internal/svc"
	"github.com/simance-ai/smdx/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarketKlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMarketKlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarketKlineLogic {
	return &MarketKlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarketKlineLogic) MarketKline(req *types.GetMarketKlineRequest) (resp *types.MarketKlineListResponse, err error) {
	marketAddress := req.MarketAddress
	chain := req.Chain
	internal := req.Internal
	fromTimestamp := req.From
	toTimestamp := req.To

	//check params
	if marketAddress == "" || chain == "" || fromTimestamp <= 0 || toTimestamp <= 0 {
		return nil, api_err.NewErrorWithCodeAndMsg(api_err.ErrCodeParamErrorCode, "invalid params")
	}

	//get market kline
	switch chain {
	case common.SolChainId:

		result, err := l.svcCtx.SolClient.GetMarketKlineList(l.ctx, &solclient.GetMarketKlineRequest{
			MarketAddress: marketAddress,
			EndTimestamp:  toTimestamp,
			Interval:      internal,
		})
		if err != nil {
			return nil, api_err.NewErrorWithCodeAndMsg(api_err.ErrCodeInternalErrorCode, err.Error())
		}

		items := make([]*types.MarketKlineResponse, 0, len(result.List))
		for _, item := range result.List {
			items = append(items, &types.MarketKlineResponse{
				O:         item.O,
				C:         item.C,
				H:         item.H,
				L:         item.L,
				V:         item.V,
				Timestamp: item.Timestamp,
			})
		}
		resp = &types.MarketKlineListResponse{
			Items: items,
		}
		return resp, nil
	default:
		return nil, api_err.ErrCodeInvalidChainNotSupport
	}
}
