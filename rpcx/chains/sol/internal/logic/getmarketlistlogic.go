package logic

import (
	"context"
	xerror "github.com/simance-ai/smdx/pkg/errors/x_err"
	"github.com/simance-ai/smdx/rpcx/chains/common/tx"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/dao/dbx"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/svc"
	"github.com/simance-ai/smdx/rpcx/chains/sol/sol"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	DirectionASC = "ASC"
)

type GetMarketListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMarketListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMarketListLogic {
	return &GetMarketListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMarketListLogic) GetMarketList(in *sol.GetMarketListRequest) (*sol.MarketListResponse, error) {
	page := in.Page
	size := in.Size
	orderBy := in.OrderBy
	direction := in.Direction

	if page < 0 {
		return &sol.MarketListResponse{}, xerror.NewErrCode(xerror.ErrPageInvalid)
	}

	if size < 0 {
		return &sol.MarketListResponse{}, xerror.NewErrCode(xerror.ErrSizeInvalid)
	}

	if size > 1000 {
		return &sol.MarketListResponse{}, xerror.NewErrCode(xerror.ErrMaxSizeInvalid)
	}

	result, count, err := dbx.FindByPageWithJoin(dbx.Use(l.svcCtx.PgDB).ReadDB(), orderBy, direction, int(page), int(size))
	if err != nil {
		return &sol.MarketListResponse{}, xerror.NewErrCodeMsg(xerror.ErrDBQueryError, err.Error())
	}
	var markets []*tx.Market

	for _, market := range result {
		markets = append(markets, &tx.Market{
			Id:                market.ID,
			Address:           market.Address,
			BaseTokenAddress:  market.BaseTokenAddress,
			QuoteTokenAddress: market.QuoteTokenAddress,

			BaseSymbol:            market.BaseSymbol,
			QuoteSymbol:           market.QuoteSymbol,
			BasePrice:             market.BasePrice,
			QuotePrice:            market.QuotePrice,
			MarketType:            market.MarketType,
			CreatedTimestamp:      market.CreatedTimestamp.Unix(),
			BaseVault:             market.BaseVault,
			QuoteVault:            market.QuoteVault,
			HolderCount:           market.HolderCount,
			Sells_24H:             int64(market.SellCount24h),
			Buys_24H:              int64(market.BuyCount24h),
			Volume_24H:            market.SellVolume24h + market.BuyVolume24h,
			PriceChangePercent1M:  market.PriceChange1m,
			PriceChangePercent5M:  market.PriceChange5m,
			PriceChangePercent1H:  market.PriceChange1h,
			PriceChangePercent6H:  market.PriceChange6h,
			PriceChangePercent24H: market.PriceChange24h,
			Liquidity:             market.Liquidity,
			MarketCap:             market.MarketCap,
			OpenTimestamp:         market.OpenTimestamp.Unix(),
			BaseIcon:              "",
			QuoteIcon:             market.LogoURL,
			InitBaseVault:         market.InitBaseVault,
			InitQuoteVault:        market.InitQuoteVault,
			BaseMintAddress:       market.BaseTokenMintAddress,
			QuoteMintAddress:      market.QuoteTokenMintAddress,
			Logo:                  market.LogoURL,
			PriceChangePercent30M: market.PriceChange24h,
		})
	}

	return &sol.MarketListResponse{
		List:  markets,
		Total: count,
	}, nil
}
