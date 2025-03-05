package common

const (
	// MarketKline1s 1s 1s 1m 5m 15m 30m 1h 4h 6h 12h 1d
	MarketKline1s  = "1s"
	MarketKline1m  = "1m"
	MarketKline5m  = "5m"
	MarketKline15m = "15m"
	MarketKline30m = "30m"
	MarketKline1h  = "1h"
	MarketKline4h  = "4h"
	MarketKline6h  = "6h"
	MarketKline12h = "12h"
	MarketKline1d  = "1d"
)

const (
	// MarketPeriodTime1m 1m 5m 1h 1d
	MarketPeriodTime1m = "1m"
	MarketPeriodTime5m = "5m"
	MarketPeriodTime1h = "1h"
	MarketPeriodTime6h = "6h"
	MarketPeriodTime1d = "1d"
)

func CheckIsSupportMarketKline(interval string) bool {
	switch interval {
	case MarketKline1s, MarketKline1m, MarketKline5m, MarketKline15m, MarketKline30m, MarketKline1h, MarketKline4h, MarketKline1d:
		return true
	case MarketKline6h, MarketKline12h:
		return true
	default:
		return false
	}
}
