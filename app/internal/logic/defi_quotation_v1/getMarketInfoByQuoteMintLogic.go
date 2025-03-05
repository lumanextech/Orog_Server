package defi_quotation_v1

import (
	"context"
	"github.com/simance-ai/smdx/app/internal/svc"
	"github.com/simance-ai/smdx/app/internal/types"
	"github.com/simance-ai/smdx/pkg/common"
	"github.com/simance-ai/smdx/pkg/errors/api_err"
	"github.com/simance-ai/smdx/rpcx/chains/sol/sol"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetMarketInfoByQuoteMintLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMarketInfoByQuoteMintLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMarketInfoByQuoteMintLogic {
	return &GetMarketInfoByQuoteMintLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMarketInfoByQuoteMintLogic) GetMarketInfoByQuoteMint(req *types.GetMarketByMintAddress) (resp *types.MarketRank, err error) {
	switch req.Chain {
	case common.SolChainId:
		mintResp, err := l.svcCtx.SolClient.GetMarketInfoByQuoteMint(l.ctx, &sol.GetMarketInfoByQuoteMintRequest{QuoteMint: req.QuoteMint})
		if err != nil {
			return nil, api_err.NewErrorWithCodeAndMsg(api_err.ErrCodeInternalErrorCode, err.Error())
		}

		resp = &types.MarketRank{
			Chain:                 common.SolChainId,
			Id:                    mintResp.Market.Id,
			Address:               mintResp.Market.Address,
			QuoteMintAddress:      mintResp.Market.QuoteMintAddress,
			QuoteTokenAddress:     mintResp.Market.QuoteTokenAddress,
			BaseTokenAddress:      mintResp.Market.BaseTokenAddress,
			BaseMintAddress:       mintResp.Market.BaseMintAddress,
			Symbol:                mintResp.Market.QuoteSymbol,
			Logo:                  mintResp.Market.Logo,
			Price:                 mintResp.Market.QuotePrice,
			BasePrice:             mintResp.Market.BasePrice,
			Swaps:                 mintResp.Market.Sells_24H + mintResp.Market.Buys_24H,
			Volume:                mintResp.Market.Volume_24H,
			PoolCreationTimestamp: mintResp.Market.CreatedTimestamp,
			OpenTimestamp:         mintResp.Market.OpenTimestamp,

			PriceChangePercent1M:  mintResp.Market.PriceChangePercent1M,
			PriceChangePercent5M:  mintResp.Market.PriceChangePercent5M,
			PriceChangePercent30M: mintResp.Market.PriceChangePercent30M,
			PriceChangePercent1H:  mintResp.Market.PriceChangePercent1H,
			PriceChangePercent6H:  mintResp.Market.PriceChangePercent6H,
			PriceChangePercent24H: mintResp.Market.PriceChangePercent24H,

			Liquidity:       0,
			MarketCap:       0,
			HotLevel:        0,
			HolderCount:     0,
			TwitterUsername: "",
			Website:         "",
			Telegram:        "",

			Buys:                   0,
			Sells:                  0,
			InitialLiquidity:       0,
			IsShowAlert:            false,
			Top10HolderRate:        0,
			RenouncedMint:          0,
			RenouncedFreezeAccount: 0,
			Launchpad:              "",
			CreatorTokenStatus:     "",
			CreatorClose:           false,
		}

		return resp, nil

	default:
		return nil, api_err.ErrCodeInvalidChainNotSupport
	}

}
