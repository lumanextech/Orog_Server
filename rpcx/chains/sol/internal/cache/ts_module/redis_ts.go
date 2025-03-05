package ts_module

import (
	"context"
	"fmt"
	"time"

	"github.com/alibaba/tair-go/tair"
	"github.com/redis/go-redis/v9"
	"github.com/simance-ai/smdx/pkg/common"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/dao/model"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	defaultRetentionPolicy = int((24 * time.Hour).Milliseconds())
)

type RedisTs struct {
	Client *tair.TairClusterClient
	logx.Logger
}

func NewRedisTs(client *tair.TairClusterClient) Ts {
	return &RedisTs{Client: client, Logger: logx.WithContext(context.Background())}
}

// CreateTsKeyIfNotExist implements Ts.
func (r *RedisTs) CreateTsKeyIfNotExist(ctx context.Context, marketAddress string) error {

	pipe := r.Client.Pipeline()
	//create key
	quotePriceType := "quote_price_ts"
	marketQuotePriceKey := fmt.Sprintf(TimeSeriesQuotePriceKey, marketAddress)
	if r.Client.TSInfo(ctx, marketQuotePriceKey).Err() != nil {
		_ = pipe.TSCreateWithArgs(ctx, marketQuotePriceKey, &redis.TSOptions{
			Retention:       defaultRetentionPolicy,
			DuplicatePolicy: "LAST",
			Labels:          map[string]string{"market": marketAddress, "type": quotePriceType},
		})
	}

	sellVolumeType := "sell_volume_ts"
	sellVolumeKey := fmt.Sprintf(TimeSeriesSellVolumeKey, marketAddress)
	if r.Client.TSInfo(ctx, sellVolumeKey).Err() != nil {
		_ = pipe.TSCreateWithArgs(ctx, sellVolumeKey, &redis.TSOptions{
			Retention:       defaultRetentionPolicy,
			DuplicatePolicy: "SUM",
			Labels:          map[string]string{"market": marketAddress, "type": sellVolumeType},
		}).Err()
	}

	buyVolumeType := "buy_volume_ts"
	buyVolumeKey := fmt.Sprintf(TimeSeriesBuyVolumeKey, marketAddress)
	if r.Client.TSInfo(ctx, buyVolumeKey).Err() != nil {
		_ = pipe.TSCreateWithArgs(ctx, buyVolumeKey, &redis.TSOptions{
			Retention:       defaultRetentionPolicy,
			DuplicatePolicy: "SUM",
			Labels:          map[string]string{"market": marketAddress, "type": buyVolumeType},
		}).Err()
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		r.Debugf("pipe.CreateTsKeyIfNotExist.Exec: %v", err)
	}

	return nil
}

// AddTsKeyPointsByTxs implements Ts.
func (r *RedisTs) AddTsKeyPointsByTxs(ctx context.Context, marketAddress string, txs []*model.MarketTx) error {
	pipe := r.Client.Pipeline()
	keySlices := make([][]interface{}, 0) //two keys

	baseVault := 0.0
	quoteVault := 0.0
	sellVolume := 0.0
	buyVolume := 0.0
	sellCount := 0
	buyCount := 0
	for _, txAt := range txs {

		ts := txAt.CreatedTimestamp.UnixMilli()

		// Add data point
		keySlices = append(keySlices,
			[]interface{}{
				fmt.Sprintf(TimeSeriesQuotePriceKey, marketAddress),
				ts,
				txAt.QuotePrice})

		baseVault += txAt.BaseAmount
		quoteVault += txAt.QuoteAmount

		switch txAt.TxType {
		case common.TxSell:
			keySlices = append(keySlices,
				[]interface{}{
					fmt.Sprintf(TimeSeriesSellVolumeKey, marketAddress),
					ts,
					txAt.Volume})

			sellVolume += txAt.Volume
			sellCount++
		case common.TxBuy:
			keySlices = append(keySlices,
				[]interface{}{
					fmt.Sprintf(TimeSeriesBuyVolumeKey, marketAddress),
					ts,
					txAt.Volume})

			buyVolume += txAt.Volume
			buyCount++
		case common.TxUnknown:
			r.Errorf("invalid txType: %d", txAt.TxType)
			continue
		}
	}

	pipe.TSIncrBy(ctx, fmt.Sprintf(TimeSeriesSellVolumeIncrByKey, marketAddress), sellVolume)
	pipe.TSIncrBy(ctx, fmt.Sprintf(TimeSeriesSellCountIncrByKey, marketAddress), float64(sellCount))
	pipe.TSIncrBy(ctx, fmt.Sprintf(TimeSeriesBuyVolumeIncrByKey, marketAddress), buyVolume)
	pipe.TSIncrBy(ctx, fmt.Sprintf(TimeSeriesBuyCountIncrByKey, marketAddress), float64(buyCount))

	pipe.TSIncrBy(ctx, fmt.Sprintf(TimeSeriesBaseVaultIncrByKey, marketAddress), baseVault)
	pipe.TSIncrBy(ctx, fmt.Sprintf(TimeSeriesQuoteVaultIncrByKey, marketAddress), quoteVault)

	pipe.TSMAdd(ctx, keySlices)

	_, err := pipe.Exec(ctx)
	if err != nil {
		r.Errorf("pipe.AddTsKeyPointsByTxs.Exec: %v", err)
	}

	return nil
}

// GetTsKLineTimeData implements Ts.
func (r *RedisTs) GetTsKLineTimeData(ctx context.Context, marketAddress string, now time.Time) (map[string]*model.MarketKline1, error) {
	pipe := r.Client.Pipeline()
	resultsFirst := make(map[string]*redis.TSTimestampValueSliceCmd)
	resultsLast := make(map[string]*redis.TSTimestampValueSliceCmd)
	resultMax := make(map[string]*redis.TSTimestampValueSliceCmd)
	resultMin := make(map[string]*redis.TSTimestampValueSliceCmd)
	resultSellVolume := make(map[string]*redis.TSTimestampValueSliceCmd)
	resultBuyVolume := make(map[string]*redis.TSTimestampValueSliceCmd)
	for key, durations := range KlineIntervalDurations {
		fromTimestamp := int(now.Truncate(durations).UnixMilli())
		toTimestamp := int(now.UnixMilli())
		bucketDuration := int(durations.Milliseconds())
		resultsFirst[key] = pipe.TSRevRangeWithArgs(ctx,
			fmt.Sprintf(TimeSeriesQuotePriceKey, marketAddress), fromTimestamp, toTimestamp, &redis.TSRevRangeOptions{
				Empty:          true,
				Count:          1,
				Aggregator:     redis.First,
				BucketDuration: bucketDuration,
			})

		resultsLast[key] = pipe.TSRevRangeWithArgs(ctx,
			fmt.Sprintf(TimeSeriesQuotePriceKey, marketAddress), fromTimestamp, toTimestamp, &redis.TSRevRangeOptions{
				Empty:          true,
				Count:          1,
				Aggregator:     redis.Last,
				BucketDuration: bucketDuration,
			})

		resultMax[key] = pipe.TSRevRangeWithArgs(ctx,
			fmt.Sprintf(TimeSeriesQuotePriceKey, marketAddress), fromTimestamp, toTimestamp, &redis.TSRevRangeOptions{
				Empty:          true,
				Count:          1,
				Aggregator:     redis.Max,
				BucketDuration: bucketDuration,
			})

		resultMin[key] = pipe.TSRevRangeWithArgs(ctx,
			fmt.Sprintf(TimeSeriesQuotePriceKey, marketAddress), fromTimestamp, toTimestamp, &redis.TSRevRangeOptions{
				Empty:          true,
				Count:          1,
				Aggregator:     redis.Min,
				BucketDuration: bucketDuration,
			})

		resultSellVolume[key] = pipe.TSRevRangeWithArgs(ctx,
			fmt.Sprintf(TimeSeriesSellVolumeKey, marketAddress), fromTimestamp, toTimestamp, &redis.TSRevRangeOptions{
				Empty:          true,
				Count:          1,
				Aggregator:     redis.Sum,
				BucketDuration: bucketDuration,
			})

		resultBuyVolume[key] = pipe.TSRevRangeWithArgs(ctx,
			fmt.Sprintf(TimeSeriesBuyVolumeKey, marketAddress), fromTimestamp, toTimestamp, &redis.TSRevRangeOptions{
				Empty:          true,
				Count:          1,
				Aggregator:     redis.Sum,
				BucketDuration: bucketDuration,
			})
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		r.Errorf("pipe.GetTsKLineTimeData.Exec: %v", err) //ignore: key not exist error
	}

	result := make(map[string]*model.MarketKline1)
	for key := range KlineIntervalDurations {

		firstResult, ok := resultsFirst[key]
		if !ok {
			r.Debugf("pipe.GetTsKLineTimeData resultsFirst not found")
			continue
		}

		lastResult, ok := resultsLast[key]
		if !ok {
			r.Debugf("pipe.GetTsKLineTimeData resultsLast not found")
			continue
		}

		maxResult, ok := resultMax[key]
		if !ok {
			r.Debugf("pipe.GetTsKLineTimeData resultMax not found")
			continue
		}

		minResult, ok := resultMin[key]
		if !ok {
			r.Debugf("pipe.GetTsKLineTimeData resultMin not found")
			continue
		}

		sellVolumeResult, ok := resultSellVolume[key]
		if !ok {
			r.Debugf("pipe.GetTsKLineTimeData resultSellVolume not found")
			continue
		}

		buyVolumeResult, ok := resultBuyVolume[key]
		if !ok {
			r.Debugf("pipe.GetTsKLineTimeData resultBuyVolume not found")
			continue
		}

		if len(firstResult.Val()) <= 0 {
			r.Debugf("pipe.GetTsKLineTimeData firstResult.Val() empty: %v %v %v", key, marketAddress, now)
			continue
		}

		timestamp := firstResult.Val()[0].Timestamp
		sellVolume := 0.0
		buyVolume := 0.0
		openPrice := 0.0
		closePrice := 0.0
		highPrice := 0.0
		lowPrice := 0.0
		if len(sellVolumeResult.Val()) > 0 {
			sellVolume = sellVolumeResult.Val()[0].Value
		}
		if len(buyVolumeResult.Val()) > 0 {
			buyVolume = buyVolumeResult.Val()[0].Value
		}
		if len(firstResult.Val()) > 0 {
			openPrice = firstResult.Val()[0].Value
		}
		if len(lastResult.Val()) > 0 {
			closePrice = lastResult.Val()[0].Value
		}
		if len(maxResult.Val()) > 0 {
			highPrice = maxResult.Val()[0].Value
		}
		if len(minResult.Val()) > 0 {
			lowPrice = minResult.Val()[0].Value
		}

		result[key] = &model.MarketKline1{
			MarketAddress: marketAddress,
			O:             openPrice,
			H:             closePrice,
			L:             highPrice,
			C:             lowPrice,
			V:             sellVolume + buyVolume,
			Timestamp:     time.Unix(timestamp, 0),
		}
	}

	return result, nil
}

// GetTsPeriodTimeData implements Ts.
func (r *RedisTs) GetTsPeriodTimeData(ctx context.Context, marketAddress string, now time.Time) (map[string]*PeriodTimeData, error) {
	pipe := r.Client.Pipeline()

	priceChangeResultMap := make(map[string]*redis.TSTimestampValueSliceCmd)
	sellVolumeResultMap := make(map[string]*redis.TSTimestampValueSliceCmd)
	buyVolumeResultMap := make(map[string]*redis.TSTimestampValueSliceCmd)
	sellCountResultMap := make(map[string]*redis.TSTimestampValueSliceCmd)
	buyCountResultMap := make(map[string]*redis.TSTimestampValueSliceCmd)
	for key, duration := range PeriodDurations {

		fromTimestamp := int(now.Add(-duration).UnixMilli())
		endTimestamp := int(now.UnixMilli())
		bucketDurationMilli := int(duration.Milliseconds())
		priceChangeResultMap[key] = pipe.TSRevRange(ctx,
			fmt.Sprintf(TimeSeriesQuotePriceKey, marketAddress),
			fromTimestamp,
			endTimestamp)

		sellVolumeResultMap[key] = pipe.TSRevRangeWithArgs(ctx,
			fmt.Sprintf(TimeSeriesSellVolumeKey, marketAddress), fromTimestamp, endTimestamp, &redis.TSRevRangeOptions{
				Empty:          true,
				Count:          1,
				Aggregator:     redis.Sum,
				BucketDuration: bucketDurationMilli,
			})

		buyVolumeResultMap[key] = pipe.TSRevRangeWithArgs(ctx,
			fmt.Sprintf(TimeSeriesBuyVolumeKey, marketAddress), fromTimestamp, endTimestamp, &redis.TSRevRangeOptions{
				Empty:          true,
				Count:          1,
				Aggregator:     redis.Sum,
				BucketDuration: bucketDurationMilli,
			})

		sellCountResultMap[key] = pipe.TSRevRangeWithArgs(ctx,
			fmt.Sprintf(TimeSeriesSellVolumeKey, marketAddress), fromTimestamp, endTimestamp, &redis.TSRevRangeOptions{
				Empty:          true,
				Count:          1,
				Aggregator:     redis.Count,
				BucketDuration: bucketDurationMilli,
			})

		buyCountResultMap[key] = pipe.TSRevRangeWithArgs(ctx,
			fmt.Sprintf(TimeSeriesBuyVolumeKey, marketAddress), fromTimestamp, endTimestamp, &redis.TSRevRangeOptions{
				Empty:          true,
				Count:          1,
				Aggregator:     redis.Count,
				BucketDuration: bucketDurationMilli,
			})
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		r.Debugf("pipe.GetTsPeriodTimeData.Exec: %v", err) //ignore: key not exist error
	}

	result := make(map[string]*PeriodTimeData)
	for key := range PeriodDurations {
		priceChangeResult, ok := priceChangeResultMap[key]
		if !ok {
			r.Debugf("pipe.GetTsPeriodTimeData priceChangeResultMap not found: %v %v", key, marketAddress)
			continue
		}

		beforePrice := 1.0
		afterPrice := 1.0
		if len(priceChangeResult.Val()) > 1 {
			beforePrice = priceChangeResult.Val()[0].Value
			afterPrice = priceChangeResult.Val()[len(priceChangeResult.Val())-1].Value
		} else {
			r.Debugf("pipe.GetTsPeriodTimeData priceChangeResult.Val() len <= 1: %v %v", key, marketAddress)
		}

		sellVolumeResult, ok := sellVolumeResultMap[key]
		if !ok {
			r.Debugf("pipe.GetTsPeriodTimeData sellVolumeResultMap not found: %v %v", key, marketAddress)
			continue
		}

		sellVolume := 0.0
		if len(sellVolumeResult.Val()) > 0 {
			sellVolume = sellVolumeResult.Val()[0].Value
		}

		buyVolumeResult, ok := buyVolumeResultMap[key]
		if !ok {
			r.Debugf("pipe.GetTsPeriodTimeData buyVolumeResultMap not found: %v %v", key, marketAddress)
			continue
		}

		buyVolume := 0.0
		if len(buyVolumeResult.Val()) > 0 {
			buyVolume = buyVolumeResult.Val()[0].Value
		}

		sellCountResult, ok := sellCountResultMap[key]
		if !ok {
			r.Debugf("pipe.GetTsPeriodTimeData sellCountResultMap not found: %v %v", key, marketAddress)
			continue
		}

		sellCount := 0.0
		if len(sellCountResult.Val()) > 0 {
			sellCount = sellCountResult.Val()[0].Value
		}

		buyCountResult, ok := buyCountResultMap[key]
		if !ok {
			r.Debugf("pipe.GetTsPeriodTimeData buyCountResultMap not found: %v %v", key, marketAddress)
			continue
		}

		buyCount := 0.0
		if len(buyCountResult.Val()) > 0 {
			buyCount = buyCountResult.Val()[0].Value
		}

		result[key] = &PeriodTimeData{
			PriceChange: r.calPriceChangePercent(afterPrice, beforePrice),
			SellVolume:  sellVolume,
			BuyVolume:   buyVolume,
			SellCount:   sellCount,
			BuyCount:    buyCount,
		}
	}

	return result, nil
}

// GetTsRealTimeData implements Ts.
func (r *RedisTs) GetTsRealTimeData(ctx context.Context, marketAddress string, now time.Time) (*RealTimeData, error) {
	pipe := r.Client.Pipeline()

	latestQuotePriceValueCmd := pipe.TSGet(ctx, fmt.Sprintf(TimeSeriesQuotePriceKey, marketAddress))
	baseVaultValueCmd := pipe.TSGet(ctx, fmt.Sprintf(TimeSeriesBaseVaultIncrByKey, marketAddress))
	quoteVaultValueCmd := pipe.TSGet(ctx, fmt.Sprintf(TimeSeriesQuoteVaultIncrByKey, marketAddress))
	sellCountCmd := pipe.TSGet(ctx, fmt.Sprintf(TimeSeriesSellCountIncrByKey, marketAddress))
	buyCountCmd := pipe.TSGet(ctx, fmt.Sprintf(TimeSeriesBuyCountIncrByKey, marketAddress))
	sellVolumeCmd := pipe.TSGet(ctx, fmt.Sprintf(TimeSeriesSellVolumeIncrByKey, marketAddress))
	buyVolumeCmd := pipe.TSGet(ctx, fmt.Sprintf(TimeSeriesBuyVolumeIncrByKey, marketAddress))

	_, err := pipe.Exec(ctx)
	if err != nil {
		r.Debugf("pipe.GetTsRealTimeData.Exec:  %v", err) //ignore: key not exist error
	}

	latestQuotePriceResult, _ := latestQuotePriceValueCmd.Result()
	baseVaultResult, _ := baseVaultValueCmd.Result()
	quoteVaultResult, _ := quoteVaultValueCmd.Result()
	sellCountResult, _ := sellCountCmd.Result()
	buyCountResult, _ := buyCountCmd.Result()
	sellVolumeResult, _ := sellVolumeCmd.Result()
	buyVolumeResult, _ := buyVolumeCmd.Result()

	quotePrice := latestQuotePriceResult.Value
	baseVault := baseVaultResult.Value
	quoteVault := quoteVaultResult.Value
	sellCount := sellCountResult.Value
	buyCount := buyCountResult.Value
	sellVolume := sellVolumeResult.Value
	buyVolume := buyVolumeResult.Value

	return &RealTimeData{
		LatestQuotePrice: quotePrice,
		BaseVault:        baseVault,
		QuoteVault:       quoteVault,
		SellCount:        sellCount,
		BuyCount:         buyCount,
		SellVolume:       sellVolume,
		BuyVolume:        buyVolume,
	}, nil
}

func (r *RedisTs) calPriceChangePercent(nowPrice, lastPrice float64) float64 {
	if lastPrice <= 0 {
		return 0
	}
	return (nowPrice - lastPrice) / lastPrice
}
