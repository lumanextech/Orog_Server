package logic

import (
	"context"
	xerror2 "github.com/simance-ai/smdx/pkg/errors/x_err"
	"github.com/simance-ai/smdx/rpcx/chains/common/tx"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/dao/dbx"

	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/svc"
	"github.com/simance-ai/smdx/rpcx/chains/sol/sol"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFollowMarketListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFollowMarketListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFollowMarketListLogic {
	return &GetFollowMarketListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFollowMarketListLogic) GetFollowMarketList(in *sol.GetMarketListRequest) (*sol.RealTimeMarketListResponse, error) {
	page := in.Page
	size := in.Size
	orderBy := in.OrderBy
	direction := in.Direction
	tokenAddressList := in.TokenAddressList
	if page < 0 {
		return &sol.RealTimeMarketListResponse{}, xerror2.NewErrCode(xerror2.ErrPageInvalid)
	}

	if size < 0 {
		return &sol.RealTimeMarketListResponse{}, xerror2.NewErrCode(xerror2.ErrSizeInvalid)
	}

	if size > 1000 {
		return &sol.RealTimeMarketListResponse{}, xerror2.NewErrCode(xerror2.ErrMaxSizeInvalid)
	}

	marketDB := dbx.Use(l.svcCtx.PgDB).ReadDB().MarketRealTimeDatum
	_, exist := marketDB.GetFieldByName(orderBy)
	if !exist {
		orderBy = marketDB.Liquidity.ColumnName().String()
	}

	result, count, err := dbx.FindByPageWithJoinFollow(dbx.Use(l.svcCtx.PgDB).ReadDB(), orderBy, direction, int(page), int(size), tokenAddressList)
	if err != nil {
		return &sol.RealTimeMarketListResponse{}, xerror2.NewErrCodeMsg(xerror2.ErrDBQueryError, err.Error())
	}

	var markets []*tx.RealTimeMarket
	for _, market := range result {
		markets = append(markets, &tx.RealTimeMarket{
			Address:               market.Address,                 // 代币地址
			PoolAddress:           market.Address,                 // 池地址
			Buys:                  market.Buys,                    // 买入数量
			Chain:                 "sol",                          // 所属链
			CreatorClose:          false,                          // 是否关闭创建者权限
			CreatorTokenStatus:    "",                             // 创建者代币状态
			HolderCount:           market.HolderCount,             // 持有人数量
			HotLevel:              0,                              // 热度等级
			Id:                    market.ID,                      // 唯一标识符
			InitialLiquidity:      0,                              // 初始流动性
			IsShowAlert:           false,                          // 是否显示警告
			Launchpad:             "",                             // 启动平台
			Liquidity:             market.Liquidity,               // 流动性
			Logo:                  market.LogoURL,                 // 代币 logo 链接
			MarketCap:             market.MarketCap,               // 市值
			OpenTimestamp:         market.OpenTimestamp.Unix(),    // 开放时间戳
			PoolCreationTimestamp: market.CreatedTimestamp.Unix(), // 池子创建时间戳
			BasePrice:             market.BasePrice,               // 主网币价格
			Price:                 0,                              // 价格
			PriceChangePercent1H:  market.PriceChange1h,           // 1 小时价格变动百分比
			PriceChangePercent1M:  market.PriceChange1m,           // 1 分钟价格变动百分比
			PriceChangePercent5M:  market.PriceChange5m,           // 5 分钟价格变动百分比
			PriceChangePercent30M: market.PriceChange30m,
			PriceChangePercent6H:  market.PriceChange6h,  // 6 小时价格变动百分比
			PriceChangePercent24H: market.PriceChange24h, // 24 小时价格变动百分比

			Sells:  market.Sells,       // 卖出数量
			Swaps:  market.Swaps,       // 交换数量
			Symbol: market.QuoteSymbol, // 代币符号 QuoteSymbol/BaseSymbol

			TwitterUsername: market.Twitter, // Twitter 用户名
			Volume:          market.Volume,  // 交易量

			BaseMintAddress:   market.BaseTokenMintAddress,
			QuoteMintAddress:  market.QuoteTokenMintAddress,
			BaseTokenAddress:  market.BaseTokenAddress,
			QuoteTokenAddress: market.QuoteTokenAddress,

			Website:                "", // 官网链接
			Telegram:               "", // Telegram 链接
			Top_10HolderRate:       0,  // 前 10 持有人占比
			RenouncedFreezeAccount: 0,  // 放弃冻结账户权限
			RenouncedMint:          0,  // 放弃铸造权限
		})
	}
	return &sol.RealTimeMarketListResponse{
		List:  markets,
		Total: count,
	}, nil
}
