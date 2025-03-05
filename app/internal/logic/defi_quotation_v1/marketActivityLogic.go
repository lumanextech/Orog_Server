package defi_quotation_v1

import (
	"context"
	"github.com/simance-ai/smdx/pkg/common"
	"github.com/simance-ai/smdx/pkg/errors/api_err"
	"github.com/simance-ai/smdx/rpcx/chains/sol/sol"

	"github.com/simance-ai/smdx/app/internal/svc"
	"github.com/simance-ai/smdx/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarketActivityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMarketActivityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarketActivityLogic {
	return &MarketActivityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarketActivityLogic) MarketActivity(req *types.GetMarketActivityRequest) (resp *types.GetMarketActivityResponse, err error) {
	chain := req.Chain
	page := req.Page
	size := req.Size
	orderBy := req.OrderBy
	address := req.MarketAddress
	direction := req.Direction

	// 设置分页的默认值
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	switch chain {
	case common.SolChainId:
		marketListResult, err := l.svcCtx.SolClient.MarketActivityList(l.ctx, &sol.GetMarketListRequest{
			Page:          page,
			Size:          size,
			OrderBy:       orderBy,
			Direction:     direction,
			MarketAddress: address,
		})
		if err != nil {
			return nil, api_err.NewErrorWithCodeAndMsg(api_err.ErrCodeInternalErrorCode, err.Error())
		}

		marketListResp := make([]*types.Activity, 0)
		for _, market := range marketListResult.List {
			marketListResp = append(marketListResp, &types.Activity{
				MakerAddress:        market.Maker,
				BaseAmount:          market.BaseAmount,
				QuoteAmount:         market.QuoteAmount,
				Volume:              market.AmountUsd,
				Timestamp:           market.Timestamp,
				SwapType:            market.Type,
				TxHash:              market.TxHash,
				QuotePrice:          market.PriceUsd,
				MakerTags:           market.MakerTags,
				MakerTwitterName:    market.MakerTwitterName,
				MakerTwitterUser:    market.MakerTwitterUsername,
				MakerName:           market.MakerName,
				MakerAvatar:         market.MakerAvatar,
				MakerENS:            market.MakerEns,
				MakerTokenTags:      market.MakerTokenTags,
				TokenAddress:        market.TokenAddress,
				QuoteAddress:        market.QuoteAddress,
				TotalTrade:          int(market.TotalTrade),
				IsFollowing:         int(market.IsFollowing),
				IsOpenOrClose:       int(market.IsOpenOrClose),
				BuyCostUSD:          market.BuyCostUsd,
				Balance:             market.Balance,
				Cost:                market.Cost,
				HistoryBoughtAmount: market.HistoryBoughtAmount,
				HistorySoldIncome:   market.HistorySoldIncome,
				HistorySoldAmount:   market.HistorySoldAmount,
				UnrealizedProfit:    market.UnrealizedProfit,
				RealizedProfit:      market.RealizedProfit,
			})

		}
		resp = &types.GetMarketActivityResponse{
			Total: marketListResult.Total,
			List:  marketListResp,
		}
		return resp, nil
	default:
		return nil, api_err.ErrCodeInvalidChainNotSupport
	}
}
