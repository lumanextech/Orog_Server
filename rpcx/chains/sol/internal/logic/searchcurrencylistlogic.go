package logic

import (
	"context"
	"fmt"

	"github.com/simance-ai/smdx/pkg/common"
	xerror "github.com/simance-ai/smdx/pkg/errors/x_err"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/dao/dbx"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/svc"
	"github.com/simance-ai/smdx/rpcx/chains/sol/sol"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchCurrencyListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchCurrencyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchCurrencyListLogic {
	return &SearchCurrencyListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchCurrencyListLogic) SearchCurrencyList(in *sol.GetMarketListRequest) (*sol.SearchCurrencyListResponse, error) {
	page := in.Page
	size := in.Size
	query := in.SearchValue
	tokenAddressList := in.TokenAddressList
	tokenAddressSet := make(map[string]struct{})
	for _, addr := range tokenAddressList {
		tokenAddressSet[addr] = struct{}{}
		fmt.Println("tokenAddressSet:", addr)
	}

	if page < 0 {
		return &sol.SearchCurrencyListResponse{}, xerror.NewErrCode(xerror.ErrPageInvalid)
	}

	if size < 0 {
		return &sol.SearchCurrencyListResponse{}, xerror.NewErrCode(xerror.ErrSizeInvalid)
	}

	if size > 20 {
		return &sol.SearchCurrencyListResponse{}, xerror.NewErrCode(xerror.ErrMaxSizeInvalid)
	}

	result, count, err := dbx.SearchCurrencyByPage(dbx.Use(l.svcCtx.PgDB).ReadDB(), query, int(page), int(size))
	if err != nil {
		return &sol.SearchCurrencyListResponse{}, xerror.NewErrCodeMsg(xerror.ErrDBQueryError, err.Error())
	}

	var markets []*sol.Currency
	for _, market := range result {
		_, found := tokenAddressSet[market.QuoteTokenMintAddress]
		markets = append(markets, &sol.Currency{
			Address:          market.Address, // 代币地址
			Volume_24:        market.BuyVolume24h + market.SellVolume24h,
			Logo:             market.LogoURL,               // 代币 logo 链接
			Liquidity:        market.Liquidity,             // 交易量
			Symbol:           market.QuoteSymbol,           // 代币符号 QuoteSymbol/BaseSymbol
			QuoteMintAddress: market.QuoteTokenMintAddress, // 池地址
			Chain:            common.SolChainId,            // 所属链
			Follow:           found,
		})
	}

	return &sol.SearchCurrencyListResponse{
		List:  markets,
		Total: count,
	}, nil
}
