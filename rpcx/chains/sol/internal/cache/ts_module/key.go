package ts_module

import "time"

var (
	LabelsTagMarket        = "market"
	DefaultPkName          = "pk"
	DefaultTsDataRetention = 24 * time.Hour
)

var (
	TimeSeriesQuotePriceKey = "{sol:ts:time}:%v:quote_price_ts" //quote_price_ts 	(Compaction policy: last, retention: 1d)
	TimeSeriesBuyVolumeKey  = "{sol:ts:time}:%v:buy_volume_ts"  //buy_volume_ts 	(Compaction policy: sum,retention: 1d)
	TimeSeriesSellVolumeKey = "{sol:ts:time}:%v:sell_volume_ts" //sell_volume_ts 	(Compaction policy: sum,retention: 1d)

	TimeSeriesQuoteVaultIncrByKey = "{sol:ts:incr}:%v:quote_vault" //quote_vault
	TimeSeriesBaseVaultIncrByKey  = "{sol:ts:incr}:%v:base_vault"  //base_vault
	TimeSeriesSellVolumeIncrByKey = "{sol:ts:incr}:%v:sell_volume" //sell_volume
	TimeSeriesBuyVolumeIncrByKey  = "{sol:ts:incr}:%v:buy_volume"  //buy_volume
	TimeSeriesSellCountIncrByKey  = "{sol:ts:incr}:%v:sell_count"  //sell_count
	TimeSeriesBuyCountIncrByKey   = "{sol:ts:incr}:%v:buy_count"   //buy_count
)
