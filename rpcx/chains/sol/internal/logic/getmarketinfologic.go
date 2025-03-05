package logic

import (
	"context"

	xerror "github.com/simance-ai/smdx/pkg/errors/x_err"
	"github.com/simance-ai/smdx/rpcx/chains/common/tx"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/svc"
	"github.com/simance-ai/smdx/rpcx/chains/sol/sol"

	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/dao/model"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetMarketInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMarketInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMarketInfoLogic {
	return &GetMarketInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMarketInfoLogic) GetMarketInfo(in *sol.GetMarketInfoRequest) (*sol.MarketInfoResponse, error) {
	//check marketAddress is valid
	if in.MarketAddress == "" {
		return nil, xerror.NewErrCode(xerror.ErrCodeInvalidMarketAddress)
	}

	//from redisCache
	result := new(sol.MarketInfoResponse)
	marketInfo := new(model.Market)

	marketInfo, err := l.svcCtx.CacheDB.GetMarketInfo(l.ctx, in.MarketAddress)
	if err != nil {
		return nil, xerror.NewErrCodeMsg(xerror.ErrCodeRedisGet, err.Error())
	}

	result = &sol.MarketInfoResponse{
		Market: &tx.Market{
			Id:                    marketInfo.ID,
			Address:               marketInfo.Address,
			BaseTokenAddress:      marketInfo.BaseTokenAddress,
			QuoteTokenAddress:     marketInfo.QuoteTokenAddress,
			BaseSymbol:            marketInfo.BaseSymbol,
			QuoteSymbol:           marketInfo.QuoteSymbol,
			MarketType:            marketInfo.MarketType,
			CreatedTimestamp:      marketInfo.CreatedTimestamp.Unix(),
			BasePrice:             0,
			QuotePrice:            0,
			BaseVault:             0,
			QuoteVault:            0,
			HolderCount:           0,
			Sells_24H:             0,
			Buys_24H:              0,
			Volume_24H:            0,
			PriceChangePercent1M:  0,
			PriceChangePercent5M:  0,
			PriceChangePercent1H:  0,
			PriceChangePercent6H:  0,
			PriceChangePercent24H: 0,
			Liquidity:             0,
			MarketCap:             0,
			PriceChangePercent30M: 0,

			OpenTimestamp:    marketInfo.OpenTimestamp.Unix(),
			BaseIcon:         "",
			QuoteIcon:        marketInfo.LogoURL,
			InitBaseVault:    marketInfo.InitBaseVault,
			InitQuoteVault:   marketInfo.InitQuoteVault,
			BaseMintAddress:  marketInfo.BaseTokenMintAddress,
			QuoteMintAddress: marketInfo.QuoteTokenMintAddress,
			Logo:             marketInfo.LogoURL,
		},
	}

	return result, nil
}
