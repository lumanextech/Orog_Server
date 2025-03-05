package cache

import (
	"time"
)

var DefaultExpireDuration = 1 * time.Hour

var (
	BlockHeightBasePriceRedisZKey = "sol:base_price_height"
	BlockHeightRedisZKey          = "sol:block_height"
)

var (
	MarketLockRedisKey        = "sol:market_lock:%v"
	MarketAddressInfoRedisKey = "sol:market_info:%v"
	//MarketAddressInfoRedisKey               = "sol:market_list:%v"
	MarketKlineInternalEndTimestampRedisKey = "sol:market_kline_internal_end_timestamp:%v:%v:%v"
)
