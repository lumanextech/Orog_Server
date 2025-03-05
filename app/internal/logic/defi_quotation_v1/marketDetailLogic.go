package defi_quotation_v1

import (
	"context"
	"fmt"
	"github.com/simance-ai/smdx/pkg/common"
	"github.com/simance-ai/smdx/pkg/errors/api_err"
	"github.com/simance-ai/smdx/rpcx/account/account"
	"github.com/simance-ai/smdx/rpcx/chains/sol/sol"

	"github.com/simance-ai/smdx/app/internal/svc"
	"github.com/simance-ai/smdx/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarketDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMarketDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarketDetailLogic {
	return &MarketDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarketDetailLogic) MarketDetail(req *types.GetMarketDetailRequest) (resp *types.GetMarketDetailResponse, err error) {
	chain := req.Chain
	marketAddress := req.MarketAddress
	// 从上下文拿出address
	var tokenAddressList []string
	address, ok := l.ctx.Value("payload").(string)
	if ok {
		//调用accountrpc的GetFollowTokenList方法，获取用户的followtoken
		followTokenList, err := l.svcCtx.AccountClient.GetFollowTokenList(l.ctx, &account.GetFollowTokenListRequest{
			Chain:   chain,
			Address: address,
		})
		if err == nil {
			tokenAddressList = followTokenList.GetTokenAddress()
		} else {
			fmt.Printf("get token address wrong: %v\n", err)
		}
	}

	switch chain {
	case common.SolChainId:
		marketResult, err := l.svcCtx.SolClient.GetMarketDetail(l.ctx, &sol.GetMarketListRequest{
			MarketAddress:    marketAddress,
			TokenAddressList: tokenAddressList,
		})
		if err != nil {
			return nil, api_err.NewErrorWithCodeAndMsg(api_err.ErrCodeInternalErrorCode, err.Error())
		}
		follow := IsAddressInList(marketResult.Address, tokenAddressList)
		resp = &types.GetMarketDetailResponse{
			Address:           marketResult.Address,
			Symbol:            marketResult.Symbol,
			Name:              marketResult.Name,
			Decimals:          marketResult.Decimals,
			Logo:              marketResult.Logo,
			BiggestPoolAddr:   marketResult.BiggestPoolAddress,
			OpenTimestamp:     marketResult.OpenTimestamp,
			HolderCount:       marketResult.HolderCount,
			CirculatingSupply: marketResult.CirculatingSupply,
			TotalSupply:       marketResult.TotalSupply,
			MaxSupply:         marketResult.MaxSupply,
			Liquidity:         marketResult.Liquidity,
			CreationTimestamp: marketResult.CreationTimestamp,
			Follow:            follow,
			BasePrice:         marketResult.BasePrice,
			Pool: types.PoolDetail{
				Address:             marketResult.Pool.Address,
				QuoteMintAddress:    marketResult.Pool.QuoteMintAddress,
				QuoteAddress:        marketResult.Pool.QuoteAddress,
				QuoteSymbol:         marketResult.Pool.QuoteSymbol,
				Liquidity:           marketResult.Pool.Liquidity,
				BaseReserve:         marketResult.Pool.BaseReserve,
				QuoteReserve:        marketResult.Pool.QuoteReserve,
				InitialLiquidity:    marketResult.Pool.InitialLiquidity,
				InitialBaseReserve:  marketResult.Pool.InitialBaseReserve,
				InitialQuoteReserve: marketResult.Pool.InitialQuoteReserve,
				CreationTimestamp:   marketResult.Pool.CreationTimestamp,
			},
			Dev: types.DevDetail{
				Address:             marketResult.Dev.Address,
				CreatorAddress:      marketResult.Dev.CreatorAddress,
				CreatorTokenBalance: marketResult.Dev.CreatorTokenBalance,
				CreatorTokenStatus:  marketResult.Dev.CreatorTokenStatus,
				Top10HolderRate:     marketResult.Dev.Top_10HolderRate,
				Telegram:            marketResult.Dev.Telegram,
				TwitterUsername:     marketResult.Dev.TwitterUsername,
				Website:             marketResult.Dev.Website,
			},
			Price: types.PriceDetail{
				Address:       marketResult.PriceInfo.Address,
				Price:         marketResult.PriceInfo.Price,
				Price1m:       marketResult.PriceInfo.Price_1M,
				Price5m:       marketResult.PriceInfo.Price_5M,
				Price1h:       marketResult.PriceInfo.Price_1H,
				Price6h:       marketResult.PriceInfo.Price_6H,
				Price24h:      marketResult.PriceInfo.Price_24H,
				Buys1m:        marketResult.PriceInfo.Buys_1M,
				Buys5m:        marketResult.PriceInfo.Buys_5M,
				Buys1h:        marketResult.PriceInfo.Buys_1H,
				Buys6h:        marketResult.PriceInfo.Buys_6H,
				Buys24h:       marketResult.PriceInfo.Buys_24H,
				Sells1m:       marketResult.PriceInfo.Sells_1M,
				Sells5m:       marketResult.PriceInfo.Sells_5M,
				Sells1h:       marketResult.PriceInfo.Sells_1H,
				Sells6h:       marketResult.PriceInfo.Sells_6H,
				Sells24h:      marketResult.PriceInfo.Sells_24H,
				Volume1m:      marketResult.PriceInfo.Volume_1M,
				Volume5m:      marketResult.PriceInfo.Volume_5M,
				Volume1h:      marketResult.PriceInfo.Volume_1H,
				Volume6h:      marketResult.PriceInfo.Volume_6H,
				Volume24h:     marketResult.PriceInfo.Volume_24H,
				BuyVolume1m:   marketResult.PriceInfo.BuyVolume_1M,
				BuyVolume5m:   marketResult.PriceInfo.BuyVolume_5M,
				BuyVolume1h:   marketResult.PriceInfo.BuyVolume_1H,
				BuyVolume6h:   marketResult.PriceInfo.BuyVolume_6H,
				BuyVolume24h:  marketResult.PriceInfo.BuyVolume_24H,
				SellVolume1m:  marketResult.PriceInfo.SellVolume_1M,
				SellVolume5m:  marketResult.PriceInfo.SellVolume_5M,
				SellVolume1h:  marketResult.PriceInfo.SellVolume_1H,
				SellVolume6h:  marketResult.PriceInfo.SellVolume_6H,
				SellVolume24h: marketResult.PriceInfo.SellVolume_24H,
				MarketCap:     marketResult.PriceInfo.MarketCap,
				Volume:        marketResult.PriceInfo.Volume,
				Swaps:         marketResult.PriceInfo.Swaps,
				Sells:         marketResult.PriceInfo.Sells,
				Buys:          marketResult.PriceInfo.Buys,
			},
		}
		return resp, nil

	default:
		return nil, api_err.ErrCodeInvalidChainNotSupport
	}
}
func IsAddressInList(address string, tokenAddressList []string) bool {
	for _, tokenAddress := range tokenAddressList {
		if tokenAddress == address {
			return true
		}
	}
	return false
}
