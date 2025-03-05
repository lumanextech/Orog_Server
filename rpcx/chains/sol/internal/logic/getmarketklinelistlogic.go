package logic

import (
	"context"
	"fmt"
	"github.com/simance-ai/smdx/pkg/common"
	"github.com/simance-ai/smdx/rpcx/chains/common/tx"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/svc"
	"github.com/simance-ai/smdx/rpcx/chains/sol/sol"
	"time"

	xerror "github.com/simance-ai/smdx/pkg/errors/x_err"
	"github.com/zeromicro/go-zero/core/logx"
)

const DefaultLimit = 1000

type GetMarketKlineListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMarketKlineListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMarketKlineListLogic {
	return &GetMarketKlineListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMarketKlineListLogic) GetMarketKlineList(in *sol.GetMarketKlineRequest) (*sol.MarketKlineListResponse, error) {
	marketAddress := in.MarketAddress
	interval := in.Interval
	endTimestamp := in.EndTimestamp //当前时间戳

	var endTime = time.Unix(endTimestamp, 0)
	var truncateStartTime time.Time //向下取整到最近的时间
	var truncateEndTime time.Time
	result := make([]*tx.MarketKline, 0)
	switch interval {
	case common.MarketKline1m:

		truncateStartTime = endTime.Truncate(1 * time.Minute)
		truncateEndTime = truncateStartTime.Add(1 * time.Minute)
		marketList, err := l.svcCtx.CacheDB.GetMarketKline1m(l.ctx, marketAddress, truncateStartTime)
		if err != nil {
			return nil, xerror.NewErrCodeMsg(xerror.ErrRedisGetError, err.Error())
		}

		for _, market := range marketList {
			result = append(result, &tx.MarketKline{
				Chain:         common.SolChainId,
				MarketAddress: market.MarketAddress,
				Timestamp:     market.Timestamp.Unix(),
				O:             market.O,
				H:             market.H,
				L:             market.L,
				C:             market.C,
				V:             market.V,
			})
		}
	case common.MarketKline5m:
		truncateStartTime = endTime.Truncate(5 * time.Minute)
		truncateEndTime = truncateStartTime.Add(5 * time.Minute)

		marketList, err := l.svcCtx.CacheDB.GetMarketKline5m(l.ctx, marketAddress, truncateStartTime)
		if err != nil {
			return nil, fmt.Errorf("GetMarketKline5m: %w", err)
		}
		for _, market := range marketList {
			result = append(result, &tx.MarketKline{
				Chain:         common.SolChainId,
				MarketAddress: market.MarketAddress,
				Timestamp:     market.Timestamp.Unix(),
				O:             market.O,
				H:             market.H,
				L:             market.L,
				C:             market.C,
				V:             market.V,
			})
		}
	case common.MarketKline15m:
		truncateStartTime = endTime.Truncate(15 * time.Minute)
		truncateEndTime = truncateStartTime.Add(15 * time.Minute)

		marketList, err := l.svcCtx.CacheDB.GetMarketKline15m(l.ctx, marketAddress, truncateStartTime)
		if err != nil {
			return nil, xerror.NewErrCodeMsg(xerror.ErrRedisGetError, err.Error())
		}
		for _, market := range marketList {
			result = append(result, &tx.MarketKline{
				Chain:         common.SolChainId,
				MarketAddress: market.MarketAddress,
				Timestamp:     market.Timestamp.Unix(),
				O:             market.O,
				H:             market.H,
				L:             market.L,
				C:             market.C,
				V:             market.V,
			})
		}
	case common.MarketKline30m:
		truncateStartTime = endTime.Truncate(30 * time.Minute)
		truncateEndTime = truncateStartTime.Add(30 * time.Minute)

		marketList, err := l.svcCtx.CacheDB.GetMarketKline30m(l.ctx, marketAddress, truncateStartTime)
		if err != nil {
			return nil, xerror.NewErrCodeMsg(xerror.ErrRedisGetError, err.Error())
		}

		for _, market := range marketList {
			result = append(result, &tx.MarketKline{
				Chain:         common.SolChainId,
				MarketAddress: market.MarketAddress,
				Timestamp:     market.Timestamp.Unix(),
				O:             market.O,
				H:             market.H,
				L:             market.L,
				C:             market.C,
				V:             market.V,
			})
		}
	case common.MarketKline1h:
		truncateStartTime = endTime.Truncate(1 * time.Hour)
		truncateEndTime = truncateStartTime.Add(1 * time.Hour)

		marketList, err := l.svcCtx.CacheDB.GetMarketKline1h(l.ctx, marketAddress, truncateStartTime)
		if err != nil {
			return nil, xerror.NewErrCodeMsg(xerror.ErrRedisGetError, err.Error())
		}

		for _, market := range marketList {
			result = append(result, &tx.MarketKline{
				Chain:         common.SolChainId,
				MarketAddress: market.MarketAddress,
				Timestamp:     market.Timestamp.Unix(),
				O:             market.O,
				H:             market.H,
				L:             market.L,
				C:             market.C,
				V:             market.V,
			})
		}
	case common.MarketKline4h:
		truncateStartTime = endTime.Truncate(4 * time.Hour)
		truncateEndTime = truncateStartTime.Add(4 * time.Hour)

		marketList, err := l.svcCtx.CacheDB.GetMarketKline4h(l.ctx, marketAddress, truncateStartTime)
		if err != nil {
			return nil, xerror.NewErrCodeMsg(xerror.ErrRedisGetError, err.Error())
		}

		for _, market := range marketList {
			result = append(result, &tx.MarketKline{
				Chain:         common.SolChainId,
				MarketAddress: market.MarketAddress,
				Timestamp:     market.Timestamp.Unix(),
				O:             market.O,
				H:             market.H,
				L:             market.L,
				C:             market.C,
				V:             market.V,
			})
		}
	case common.MarketKline6h:
		truncateStartTime = endTime.Truncate(6 * time.Hour)
		truncateEndTime = truncateStartTime.Add(6 * time.Hour)

		marketList, err := l.svcCtx.CacheDB.GetMarketKline6h(l.ctx, marketAddress, truncateStartTime)
		if err != nil {
			return nil, xerror.NewErrCodeMsg(xerror.ErrRedisGetError, err.Error())
		}

		for _, market := range marketList {
			result = append(result, &tx.MarketKline{
				Chain:         common.SolChainId,
				MarketAddress: market.MarketAddress,
				Timestamp:     market.Timestamp.Unix(),
				O:             market.O,
				H:             market.H,
				L:             market.L,
				C:             market.C,
				V:             market.V,
			})
		}
	case common.MarketKline12h:
		truncateStartTime = endTime.Truncate(12 * time.Hour)
		truncateEndTime = truncateStartTime.Add(12 * time.Hour)

		marketList, err := l.svcCtx.CacheDB.GetMarketKline12h(l.ctx, marketAddress, truncateStartTime)
		if err != nil {
			return nil, xerror.NewErrCodeMsg(xerror.ErrRedisGetError, err.Error())
		}

		for _, market := range marketList {
			result = append(result, &tx.MarketKline{
				Chain:         common.SolChainId,
				MarketAddress: market.MarketAddress,
				Timestamp:     market.Timestamp.Unix(),
				O:             market.O,
				H:             market.H,
				L:             market.L,
				C:             market.C,
				V:             market.V,
			})
		}
	case common.MarketKline1d:
		truncateStartTime = endTime.Truncate(24 * time.Hour)
		truncateEndTime = truncateStartTime.Add(24 * time.Hour)

		marketList, err := l.svcCtx.CacheDB.GetMarketKline1d(l.ctx, marketAddress, truncateStartTime)
		if err != nil {
			return nil, xerror.NewErrCodeMsg(xerror.ErrRedisGetError, err.Error())
		}

		for _, market := range marketList {
			result = append(result, &tx.MarketKline{
				Chain:         common.SolChainId,
				MarketAddress: market.MarketAddress,
				Timestamp:     market.Timestamp.Unix(),
				O:             market.O,
				H:             market.H,
				L:             market.L,
				C:             market.C,
				V:             market.V,
			})
		}
	default:
		return nil, xerror.NewErrCodeMsg(xerror.ErrInvalidParam, "invalid interval")
	}

	//当前那一根的蜡烛 需要通过market_tx表查找
	var resultLatest []*tx.MarketKline
	err := l.svcCtx.ReadPgDB.WithContext(l.ctx).
		Raw(`
        SELECT market_address,
               FIRST(quote_price, created_timestamp) AS o,
               MAX(quote_price) AS h,
               MIN(quote_price) AS l,
               LAST(quote_price, created_timestamp) AS c,
               SUM(volume) AS v
        FROM market_tx
        WHERE market_address = ?
          AND created_timestamp >= ?
          AND created_timestamp <= ?
        GROUP BY market_address;
    `, marketAddress, truncateStartTime, truncateEndTime).
		Scan(&resultLatest).Error
	if err != nil {
		return nil, fmt.Errorf("db.Where: %w", err)
	}

	for _, resL := range resultLatest {
		resL.Timestamp = truncateStartTime.Unix()
	}

	result = append(result, resultLatest...)

	return &sol.MarketKlineListResponse{List: result}, nil

}
