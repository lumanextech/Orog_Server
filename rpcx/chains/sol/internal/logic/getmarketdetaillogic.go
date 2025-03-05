package logic

import (
	"context"
	xerror "github.com/simance-ai/smdx/pkg/errors/x_err"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/dao/dbx"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/dao/model"

	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/svc"
	"github.com/simance-ai/smdx/rpcx/chains/sol/sol"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMarketDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMarketDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMarketDetailLogic {
	return &GetMarketDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMarketDetailLogic) GetMarketDetail(in *sol.GetMarketListRequest) (*sol.MarketDetailResponse, error) {
	if in.MarketAddress == "" {
		return nil, xerror.NewErrCode(xerror.ErrCodeInvalidQuoteMintAddress)
	}

	marketResult := new(model.Market)
	marketDB := dbx.Use(l.svcCtx.PgDB).Market
	marketResult, err := marketDB.WithContext(l.ctx).ReadDB().Where(marketDB.Address.Eq(in.MarketAddress)).First()
	if err != nil {
		return nil, xerror.NewErrCode(xerror.ErrDBQueryError)
	}

	marketRealTimeResult := new(model.MarketRealTimeDatum)
	marketRealTimeDB := dbx.Use(l.svcCtx.PgDB).MarketRealTimeDatum
	marketRealTimeResult, err = marketRealTimeDB.WithContext(l.ctx).ReadDB().Where(marketRealTimeDB.Address.Eq(in.MarketAddress)).First()
	if err != nil {
		return nil, xerror.NewErrCode(xerror.ErrDBQueryError)
	}

	marketAuditResult := new(model.MarketAuditMedium)
	marketAuditDB := dbx.Use(l.svcCtx.PgDB).MarketAuditMedium
	marketAuditResult, err = marketAuditDB.WithContext(l.ctx).ReadDB().Where(marketAuditDB.Address.Eq(in.MarketAddress)).First()
	if err != nil {
		return nil, xerror.NewErrCode(xerror.ErrDBQueryError)
	}

	return &sol.MarketDetailResponse{
		Address:            marketResult.Address,
		Symbol:             marketResult.QuoteSymbol,
		Name:               marketAuditResult.Name,
		Decimals:           int64(marketResult.QuoteTokenDecimals),
		Logo:               marketResult.LogoURL,
		BiggestPoolAddress: "",
		OpenTimestamp:      marketResult.OpenTimestamp.Unix(),
		HolderCount:        marketRealTimeResult.HolderCount,
		CirculatingSupply:  0,
		TotalSupply:        marketAuditResult.QuoteMaxSupply,
		MaxSupply:          0,
		Liquidity:          marketRealTimeResult.Liquidity,
		CreationTimestamp:  marketResult.CreatedTimestamp.Unix(),
		Price:              marketRealTimeResult.QuotePrice,
		BasePrice:          marketRealTimeResult.BasePrice,
		Follow:             false,
		Pool: &sol.Pool{
			Address:             marketResult.Address,
			QuoteMintAddress:    marketResult.QuoteTokenMintAddress,
			QuoteAddress:        marketResult.QuoteTokenAddress,
			QuoteSymbol:         marketResult.QuoteSymbol,
			Liquidity:           marketRealTimeResult.Liquidity,
			BaseReserve:         marketRealTimeResult.BaseVault,
			QuoteReserve:        marketRealTimeResult.QuoteVault,
			InitialLiquidity:    marketResult.InitLiquidity,
			InitialBaseReserve:  marketResult.InitBaseVault,
			InitialQuoteReserve: marketResult.InitQuoteVault,
			CreationTimestamp:   marketResult.CreatedTimestamp.Unix(),
		},
		Dev: &sol.Developer{
			Address:             marketResult.Address,
			CreatorAddress:      marketResult.DevAddress,
			CreatorTokenBalance: 0,
			CreatorTokenStatus:  marketAuditResult.CreatorTokenStatus,
			Top_10HolderRate:    0,
			Telegram:            marketAuditResult.Telegram,
			TwitterUsername:     marketAuditResult.Twitter,
			Website:             marketAuditResult.OfficialWebsite,
		},
		PriceInfo: &sol.PriceInfo{
			Address:        marketResult.Address,
			Price:          marketRealTimeResult.QuotePrice,
			Price_1M:       marketRealTimeResult.PriceChange1m,
			Price_5M:       marketRealTimeResult.PriceChange5m,
			Price_1H:       marketRealTimeResult.PriceChange1h,
			Price_6H:       marketRealTimeResult.PriceChange6h,
			Price_24H:      marketRealTimeResult.PriceChange24h,
			Buys_1M:        marketRealTimeResult.BuyCount1m,
			Buys_5M:        marketRealTimeResult.BuyCount5m,
			Buys_1H:        marketRealTimeResult.BuyCount1h,
			Buys_6H:        marketRealTimeResult.BuyCount6h,
			Buys_24H:       marketRealTimeResult.BuyCount24h,
			Sells_1M:       marketRealTimeResult.SellCount1m,
			Sells_5M:       marketRealTimeResult.SellCount5m,
			Sells_1H:       marketRealTimeResult.SellCount1h,
			Sells_6H:       marketRealTimeResult.SellCount6h,
			Sells_24H:      marketRealTimeResult.SellCount24h,
			Volume_1M:      marketRealTimeResult.BuyVolume1m + marketRealTimeResult.SellVolume1m,
			Volume_5M:      marketRealTimeResult.BuyVolume5m + marketRealTimeResult.SellVolume5m,
			Volume_1H:      marketRealTimeResult.BuyVolume1h + marketRealTimeResult.SellVolume1h,
			Volume_6H:      marketRealTimeResult.BuyVolume6h + marketRealTimeResult.SellVolume6h,
			Volume_24H:     marketRealTimeResult.BuyVolume24h + marketRealTimeResult.SellVolume24h,
			BuyVolume_1M:   marketRealTimeResult.BuyVolume1m,
			BuyVolume_5M:   marketRealTimeResult.BuyVolume5m,
			BuyVolume_1H:   marketRealTimeResult.BuyVolume1h,
			BuyVolume_6H:   marketRealTimeResult.BuyVolume6h,
			BuyVolume_24H:  marketRealTimeResult.BuyVolume24h,
			SellVolume_1M:  marketRealTimeResult.SellVolume1m,
			SellVolume_5M:  marketRealTimeResult.SellVolume5m,
			SellVolume_1H:  marketRealTimeResult.SellVolume1h,
			SellVolume_6H:  marketRealTimeResult.SellVolume6h,
			SellVolume_24H: marketRealTimeResult.SellVolume24h,
			Volume:         marketRealTimeResult.Volume,
			MarketCap:      marketRealTimeResult.MarketCap,
			Swaps:          marketRealTimeResult.Swaps,
			Sells:          marketRealTimeResult.Sells,
			Buys:           marketRealTimeResult.Buys,
		},
	}, nil
}
