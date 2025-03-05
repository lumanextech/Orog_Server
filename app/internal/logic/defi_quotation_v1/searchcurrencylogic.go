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

type SearchCurrencyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchCurrencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchCurrencyLogic {
	return &SearchCurrencyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchCurrencyLogic) SearchCurrency(req *types.GetSearchCurrencyRequest) (resp *types.GetSearchCurrencyResponse, err error) {
	chain := req.Chain
	page := req.Page
	size := req.Size
	searchValue := req.Query
	var tokenAddressList []string
	if page == 0 || size == 0 {
		page = 1
		size = 10
	}
	// 从上下文拿出address
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
		marketListResult, err := l.svcCtx.SolClient.SearchCurrencyList(l.ctx, &sol.GetMarketListRequest{
			Page:             page,
			Size:             size,
			SearchValue:      searchValue,
			TokenAddressList: tokenAddressList,
		})
		if err != nil {
			return nil, api_err.NewErrorWithCodeAndMsg(api_err.ErrCodeInternalErrorCode, err.Error())
		}

		marketListResp := make([]*types.Currency, 0)
		for _, market := range marketListResult.List {
			marketListResp = append(marketListResp, &types.Currency{
				Follow:           market.Follow,
				QuoteMintAddress: market.QuoteMintAddress,
				Liquidity:        market.Liquidity,
				Chain:            market.Chain,
				Address:          market.Address,
				Symbol:           market.Symbol,
				Logo:             market.Logo,
				Volume24H:        market.Volume_24,
			})
		}
		resp = &types.GetSearchCurrencyResponse{
			Total: marketListResult.Total,
			List:  marketListResp,
		}
		return resp, nil
	default:
		return nil, api_err.ErrCodeInvalidChainNotSupport
	}

	return
}
