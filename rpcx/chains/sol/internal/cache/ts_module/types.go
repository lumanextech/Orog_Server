package ts_module

import (
	"time"

	"github.com/simance-ai/smdx/pkg/common"
)

var (
	KlineIntervalDurations = map[string]time.Duration{
		common.MarketKline1s:  time.Second,
		common.MarketKline1m:  time.Minute,
		common.MarketKline5m:  5 * time.Minute,
		common.MarketKline15m: 15 * time.Minute,
		common.MarketKline30m: 30 * time.Minute,
		common.MarketKline1h:  time.Hour,
		common.MarketKline4h:  4 * time.Hour,
		common.MarketKline6h:  6 * time.Hour,
		common.MarketKline12h: 12 * time.Hour,
		common.MarketKline1d:  24 * time.Hour,
	}

	PeriodDurations = map[string]time.Duration{
		common.MarketPeriodTime1m: 1 * time.Minute,
		common.MarketPeriodTime5m: 5 * time.Minute,
		common.MarketPeriodTime1h: 1 * time.Hour,
		common.MarketPeriodTime6h: 6 * time.Hour,
		common.MarketPeriodTime1d: 24 * time.Hour,
	}
)

// RealTimeData
// @TimeRange (pkg/common/market_time.go)
type RealTimeData struct {
	LatestQuotePrice float64 //计数器-最新成交价
	BaseVault        float64 //计数器-池中base总量
	QuoteVault       float64 //计数器-池中quote总量
	SellCount        float64 //计数器-卖出笔数总量
	BuyCount         float64 //计数器-买入笔数总量
	SellVolume       float64 //计数器-卖出成交额
	BuyVolume        float64 //计数器-买入成交额
}

// PeriodTimeData
// @TimeRange (pkg/common/market_time.go)
type PeriodTimeData struct {
	PriceChange float64 //非计数器-价格变化
	SellVolume  float64 //非计数器-卖出成交额
	BuyVolume   float64 //非计数器-买入成交额
	SellCount   float64 //非计数器-卖出笔数
	BuyCount    float64 //非计数器-买入笔数
}
