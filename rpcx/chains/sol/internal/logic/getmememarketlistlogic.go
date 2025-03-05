package logic

import (
	"context"
	"github.com/simance-ai/smdx/pkg/common"
	xerror "github.com/simance-ai/smdx/pkg/errors/x_err"
	"github.com/simance-ai/smdx/rpcx/chains/common/tx"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/dao/dbx"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/svc"
	"github.com/simance-ai/smdx/rpcx/chains/sol/sol"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetMemeMarketListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMemeMarketListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMemeMarketListLogic {
	return &GetMemeMarketListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMemeMarketListLogic) GetMemeMarketList(in *sol.GetMarketListRequest) (*sol.RealTimeMarketListResponse, error) {
	page := in.Page
	size := in.Size
	orderBy := in.OrderBy
	direction := in.Direction
	pumpOrRaydium := in.PumpOrRaydium

	if page < 0 {
		return &sol.RealTimeMarketListResponse{}, xerror.NewErrCode(xerror.ErrPageInvalid)
	}

	if size < 0 {
		return &sol.RealTimeMarketListResponse{}, xerror.NewErrCode(xerror.ErrSizeInvalid)
	}

	if size > 1000 {
		return &sol.RealTimeMarketListResponse{}, xerror.NewErrCode(xerror.ErrMaxSizeInvalid)
	}

	marketDB := dbx.Use(l.svcCtx.PgDB).ReadDB().MarketRealTimeDatum
	_, exist := marketDB.GetFieldByName(orderBy)
	if !exist {
		orderBy = marketDB.Liquidity.ColumnName().String()
	}

	result, count, err := dbx.FindByPageWithJoinWithMarketType(dbx.Use(l.svcCtx.PgDB).ReadDB(), orderBy, direction, int(page), int(size), pumpOrRaydium)
	if err != nil {
		return &sol.RealTimeMarketListResponse{}, xerror.NewErrCodeMsg(xerror.ErrDBQueryError, err.Error())
	}

	var markets []*tx.RealTimeMarket
	for _, market := range result {
		markets = append(markets, &tx.RealTimeMarket{
			Address:               market.Address, // 代币地址
			PoolAddress:           market.Address, // 池地址
			BaseMintAddress:       market.BaseTokenMintAddress,
			QuoteMintAddress:      market.QuoteTokenMintAddress,
			PriceChangePercent30M: market.PriceChange30m,
			BaseTokenAddress:      market.BaseTokenAddress,
			QuoteTokenAddress:     market.QuoteTokenAddress,

			Buys:                  market.Buys,                    // 买入数量
			Chain:                 common.SolChainId,              // 所属链
			HolderCount:           market.HolderCount,             // 持有人数量
			Id:                    market.ID,                      // 唯一标识符
			Liquidity:             market.Liquidity,               // 流动性
			Logo:                  market.LogoURL,                 // 代币 logo 链接
			MarketCap:             market.MarketCap,               // 市值
			OpenTimestamp:         market.OpenTimestamp.Unix(),    // 开放时间戳
			PoolCreationTimestamp: market.CreatedTimestamp.Unix(), // 池子创建时间戳
			BasePrice:             market.BasePrice,               // 主网币价格
			PriceChangePercent1H:  market.PriceChange1h,           // 1 小时价格变动百分比
			PriceChangePercent1M:  market.PriceChange1m,           // 1 分钟价格变动百分比
			PriceChangePercent5M:  market.PriceChange5m,           // 5 分钟价格变动百分比
			PriceChangePercent6H:  market.PriceChange6h,           // 6 小时价格变动百分比
			PriceChangePercent24H: market.PriceChange24h,          // 24 小时价格变动百分比
			Sells:                 market.Sells,                   // 卖出数量
			Swaps:                 market.Swaps,                   // 交换数量
			Symbol:                market.QuoteSymbol,             // 代币符号 QuoteSymbol/BaseSymbol
			TwitterUsername:       market.Twitter,                 // Twitter 用户名
			Volume:                market.Volume,                  // 交易量
			Price:                 market.QuotePrice,              // 价格
			InitialLiquidity:      0,                              // 初始流动性

			CreatorClose:           false, // 是否关闭创建者权限
			CreatorTokenStatus:     "",    // 创建者代币状态
			HotLevel:               0,     // 热度等级
			IsShowAlert:            false, // 是否显示警告
			Launchpad:              "",    // 启动平台
			RenouncedFreezeAccount: 0,     // 放弃冻结账户权限
			RenouncedMint:          0,     // 放弃铸造权限
			Telegram:               "",    // Telegram 链接
			Top_10HolderRate:       0,     // 前 10 持有人占比
			Website:                "",    // 官网链接

		})
	}
	return &sol.RealTimeMarketListResponse{
		List:  markets,
		Total: count,
	}, nil
}
