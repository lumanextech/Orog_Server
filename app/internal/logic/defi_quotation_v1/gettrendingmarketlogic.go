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

// GetTrendingMarketLogic 获取市场排名Rank信息
type GetTrendingMarketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTrendingMarketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTrendingMarketLogic {
	return &GetTrendingMarketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTrendingMarketLogic) GetTrendingMarket(req *types.GetTrendingMarketRequest) (resp *types.GetTrendingMarketResponse, err error) {
	chain := req.Chain
	page := req.Page
	size := req.Size
	orderBy := req.OrderBy
	filters := req.Filters
	direction := req.Direction
	period := req.Period

	switch chain {
	case common.SolChainId:
		marketListResult, err := l.svcCtx.SolClient.GetTrendingMarket(l.ctx, &sol.GetMarketListRequest{
			Page:      page,
			Size:      size,
			OrderBy:   orderBy,
			Direction: direction,
			Period:    period,
			Filters:   filters,
		})
		if err != nil {
			return nil, api_err.NewErrorWithCodeAndMsg(api_err.ErrCodeInternalErrorCode, err.Error())
		}

		marketListResp := make([]*types.MarketRank, 0)
		for _, market := range marketListResult.List {
			marketListResp = append(marketListResp, &types.MarketRank{
				Id:                     market.Id,
				Chain:                  common.SolChainId,
				Address:                market.Address,
				QuoteMintAddress:       market.QuoteMintAddress,
				QuoteTokenAddress:      market.QuoteTokenAddress,
				BaseTokenAddress:       market.BaseTokenAddress,
				BaseMintAddress:        market.BaseMintAddress,
				Symbol:                 market.Symbol,
				Logo:                   market.Logo,
				Price:                  market.Price,
				BasePrice:              market.BasePrice,
				Swaps:                  market.Swaps,
				Volume:                 market.Volume,
				Liquidity:              market.Liquidity,
				MarketCap:              market.MarketCap,
				HotLevel:               int64(market.HotLevel),
				PoolCreationTimestamp:  market.PoolCreationTimestamp,
				HolderCount:            market.HolderCount,
				TwitterUsername:        market.TwitterUsername,
				Website:                market.Website,
				Telegram:               market.Telegram,
				OpenTimestamp:          market.OpenTimestamp,
				PriceChangePercent1M:   market.PriceChangePercent1M,
				PriceChangePercent5M:   market.PriceChangePercent5M,
				PriceChangePercent30M:  market.PriceChangePercent30M,
				PriceChangePercent1H:   market.PriceChangePercent1H,
				PriceChangePercent6H:   market.PriceChangePercent6H,
				PriceChangePercent24H:  market.PriceChangePercent24H,
				Buys:                   market.Buys,
				Sells:                  market.Sells,
				InitialLiquidity:       market.InitialLiquidity,
				IsShowAlert:            market.IsShowAlert,
				Top10HolderRate:        market.Top_10HolderRate,
				RenouncedMint:          market.RenouncedMint,
				RenouncedFreezeAccount: market.RenouncedFreezeAccount,
				Launchpad:              "",
				CreatorTokenStatus:     "",
				CreatorClose:           false,
			})
		}
		resp = &types.GetTrendingMarketResponse{
			Total: marketListResult.Total,
			List:  marketListResp,
		}
		return resp, nil
	default:
		return nil, api_err.ErrCodeInvalidChainNotSupport
	}

	return
}
