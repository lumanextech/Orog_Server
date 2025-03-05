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

type TairTs struct {
	Client *tair.TairClusterClient
	logx.Logger
}

func NewTairTs(client *tair.TairClusterClient) Ts {
	return &TairTs{Client: client, Logger: logx.WithContext(context.Background())}
}

func (t *TairTs) CreateTsKeyIfNotExist(ctx context.Context, marketAddress string) error {
	pipe := t.Client.TairPipeline()

	//Creates an skey in a specified pkey. If the pkey does not exist, it is automatically created.
	//If an skey with the same name already exists, the skey cannot be created.
	opts := new(tair.ExTsAttributeArgs).New()
	opts.DataEt(DefaultTsDataRetention.Milliseconds())
	opts.Labels([]string{"market", marketAddress})

	_ = pipe.TsSCreate(ctx, fmt.Sprintf(TimeSeriesQuotePriceKey, DefaultPkName), marketAddress, opts)
	_ = pipe.TsSCreate(ctx, fmt.Sprintf(TimeSeriesBuyVolumeKey, DefaultPkName), marketAddress, opts)
	_ = pipe.TsSCreate(ctx, fmt.Sprintf(TimeSeriesSellVolumeKey, DefaultPkName), marketAddress, opts)

	_, err := pipe.Exec(ctx)
	if err != nil {
		t.Debugf("pipe.CreateTsKeyIfNotExist.Exec: %v", err)
	}

	return nil
}

func (t *TairTs) AddTsKeyPointsByTxs(ctx context.Context, marketAddress string, txs []*model.MarketTx) error {
	if len(txs) <= 0 {
		return nil
	}

	pipe := t.Client.TairPipeline()
	keyQuotePriceSlices := make([]*tair.ExTsDataPoint, 0)
	keySellVolumeSlices := make([]*tair.ExTsDataPoint, 0)
	keyBuyVolumeSlices := make([]*tair.ExTsDataPoint, 0)
	for _, txAt := range txs {
		ts := fmt.Sprintf("%d", txAt.CreatedTimestamp.UnixMilli())

		// Add data point
		epQuotePrice := new(tair.ExTsDataPoint)
		epQuotePrice.SetSKey(marketAddress)
		epQuotePrice.SetTs(ts)
		epQuotePrice.SetValue(txAt.QuotePrice)
		keyQuotePriceSlices = append(keyQuotePriceSlices, epQuotePrice)

		pipe.ExTsIncr(ctx, fmt.Sprintf(TimeSeriesBaseVaultIncrByKey, DefaultPkName), marketAddress, ts, txAt.BaseAmount)
		pipe.ExTsIncr(ctx, fmt.Sprintf(TimeSeriesQuoteVaultIncrByKey, DefaultPkName), marketAddress, ts, txAt.QuoteAmount)

		switch txAt.TxType {
		case common.TxSell:
			epSellVolume := new(tair.ExTsDataPoint)
			epSellVolume.SetSKey(marketAddress)
			epSellVolume.SetTs(ts)
			epSellVolume.SetValue(txAt.Volume)
			keySellVolumeSlices = append(keySellVolumeSlices, epSellVolume)

			pipe.ExTsIncr(ctx, fmt.Sprintf(TimeSeriesSellVolumeIncrByKey, DefaultPkName), marketAddress, ts, txAt.Volume)
			pipe.ExTsIncr(ctx, fmt.Sprintf(TimeSeriesSellCountIncrByKey, DefaultPkName), marketAddress, ts, 1)
		case common.TxBuy:
			epBuyVolume := new(tair.ExTsDataPoint)
			epBuyVolume.SetSKey(marketAddress)
			epBuyVolume.SetTs(ts)
			epBuyVolume.SetValue(txAt.Volume)
			keyBuyVolumeSlices = append(keyBuyVolumeSlices, epBuyVolume)

			pipe.ExTsIncr(ctx, fmt.Sprintf(TimeSeriesBuyVolumeIncrByKey, DefaultPkName), marketAddress, ts, txAt.Volume)
			pipe.ExTsIncr(ctx, fmt.Sprintf(TimeSeriesBuyCountIncrByKey, DefaultPkName), marketAddress, ts, 1)
		case common.TxUnknown:
			continue
		}
	}

	t.cmdMADD(ctx, pipe, fmt.Sprintf(TimeSeriesQuotePriceKey, DefaultPkName), keyQuotePriceSlices)
	t.cmdMADD(ctx, pipe, fmt.Sprintf(TimeSeriesSellVolumeKey, DefaultPkName), keySellVolumeSlices)
	t.cmdMADD(ctx, pipe, fmt.Sprintf(TimeSeriesBuyVolumeKey, DefaultPkName), keyBuyVolumeSlices)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("pipe.AddTsKeyPointsByTxs.Exec: %w %v", err, txs[0].CreatedTimestamp)
	}

	return nil
}

func (t *TairTs) GetTsKLineTimeData(ctx context.Context, marketAddress string, now time.Time) (map[string]*model.MarketKline1, error) {
	pipe := t.Client.TairPipeline()
	resultsFirst := make(map[string]*tair.ExTsSKeyCmd)
	resultsLast := make(map[string]*tair.ExTsSKeyCmd)
	resultMax := make(map[string]*tair.ExTsSKeyCmd)
	resultMin := make(map[string]*tair.ExTsSKeyCmd)
	resultSellVolume := make(map[string]*tair.ExTsSKeyCmd)
	resultBuyVolume := make(map[string]*tair.ExTsSKeyCmd)
	for key, durations := range KlineIntervalDurations {
		fromTimestamp := fmt.Sprintf("%d", now.Truncate(durations).UnixMilli())
		toTimestamp := fmt.Sprintf("%d", now.UnixMilli())
		bucketDuration := durations.Milliseconds()

		optsFirst := new(tair.ExTsAggregationArgs).New()
		optsFirst.Reverse()
		optsFirst.MaxCount(1)
		optsFirst.First(bucketDuration)
		resultsFirst[key] = pipe.ExTsRangeArgs(ctx, fmt.Sprintf(TimeSeriesQuotePriceKey, DefaultPkName), marketAddress, fromTimestamp, toTimestamp, optsFirst)

		optsLast := new(tair.ExTsAggregationArgs).New()
		optsLast.Reverse()
		optsLast.MaxCount(1)
		optsLast.Last(bucketDuration)
		resultsLast[key] = pipe.ExTsRangeArgs(ctx, fmt.Sprintf(TimeSeriesQuotePriceKey, DefaultPkName), marketAddress, fromTimestamp, toTimestamp, optsLast)

		optsMax := new(tair.ExTsAggregationArgs).New()
		optsMax.Reverse()
		optsMax.MaxCount(1)
		optsMax.Max(bucketDuration)
		resultMax[key] = pipe.ExTsRangeArgs(ctx, fmt.Sprintf(TimeSeriesQuotePriceKey, DefaultPkName), marketAddress, fromTimestamp, toTimestamp, optsMax)

		optsMin := new(tair.ExTsAggregationArgs).New()
		optsMin.Reverse()
		optsMin.MaxCount(1)
		optsMin.Min(bucketDuration)
		resultMin[key] = pipe.ExTsRangeArgs(ctx, fmt.Sprintf(TimeSeriesQuotePriceKey, DefaultPkName), marketAddress, fromTimestamp, toTimestamp, optsMin)

		optsSellVolume := new(tair.ExTsAggregationArgs).New()
		optsSellVolume.Reverse()
		optsSellVolume.MaxCount(1)
		optsSellVolume.Sum(bucketDuration)
		resultSellVolume[key] = pipe.ExTsRangeArgs(ctx, fmt.Sprintf(TimeSeriesSellVolumeKey, DefaultPkName), marketAddress, fromTimestamp, toTimestamp, optsSellVolume)

		optsBuyVolume := new(tair.ExTsAggregationArgs).New()
		optsBuyVolume.Reverse()
		optsBuyVolume.MaxCount(1)
		optsBuyVolume.Sum(bucketDuration)
		resultBuyVolume[key] = pipe.ExTsRangeArgs(ctx, fmt.Sprintf(TimeSeriesBuyVolumeKey, DefaultPkName), marketAddress, fromTimestamp, toTimestamp, optsBuyVolume)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		t.Errorf("pipe.GetTsKLineTimeData.Exec: %w", err) //ignore: key not exist error
	}

	result := make(map[string]*model.MarketKline1)
	for key := range KlineIntervalDurations {

		firstResult, ok := resultsFirst[key]
		if !ok {
			t.Debugf("pipe.GetTsKLineTimeData resultsFirst not found")
			continue
		}

		lastResult, ok := resultsLast[key]
		if !ok {
			t.Debugf("pipe.GetTsKLineTimeData resultsLast not found")
			continue
		}

		maxResult, ok := resultMax[key]
		if !ok {
			t.Debugf("pipe.GetTsKLineTimeData resultMax not found")
			continue
		}

		minResult, ok := resultMin[key]
		if !ok {
			t.Debugf("pipe.GetTsKLineTimeData resultMin not found")
			continue
		}

		sellVolumeResult, ok := resultSellVolume[key]
		if !ok {
			t.Debugf("pipe.GetTsKLineTimeData resultSellVolume not found")
			continue
		}

		buyVolumeResult, ok := resultBuyVolume[key]
		if !ok {
			t.Debugf("pipe.GetTsKLineTimeData resultBuyVolume not found")
			continue
		}

		if len(firstResult.DataPoints()) <= 0 {
			t.Debugf("pipe.GetTsKLineTimeData firstResult.Val() empty: %v %v %v", key, marketAddress, now)
			continue
		}

		timestamp := firstResult.DataPoints()[0].Ts()
		sellVolume := 0.0
		buyVolume := 0.0
		openPrice := 0.0
		closePrice := 0.0
		highPrice := 0.0
		lowPrice := 0.0
		if len(sellVolumeResult.DataPoints()) > 0 {
			sellVolume = sellVolumeResult.DataPoints()[0].Value()
		}
		if len(buyVolumeResult.DataPoints()) > 0 {
			buyVolume = buyVolumeResult.DataPoints()[0].Value()
		}
		if len(firstResult.DataPoints()) > 0 {
			openPrice = firstResult.DataPoints()[0].Value()
		}
		if len(lastResult.DataPoints()) > 0 {
			closePrice = lastResult.DataPoints()[0].Value()
		}
		if len(maxResult.DataPoints()) > 0 {
			highPrice = maxResult.DataPoints()[0].Value()
		}
		if len(minResult.DataPoints()) > 0 {
			lowPrice = minResult.DataPoints()[0].Value()
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

func (t *TairTs) GetTsPeriodTimeData(ctx context.Context, marketAddress string, now time.Time) (map[string]*PeriodTimeData, error) {
	pipe := t.Client.TairPipeline()

	priceChangeResultMap := make(map[string]*tair.ExTsSKeyCmd)
	sellVolumeResultMap := make(map[string]*tair.ExTsSKeyCmd)
	buyVolumeResultMap := make(map[string]*tair.ExTsSKeyCmd)
	sellCountResultMap := make(map[string]*tair.ExTsSKeyCmd)
	buyCountResultMap := make(map[string]*tair.ExTsSKeyCmd)
	for key, duration := range PeriodDurations {

		fromTimestamp := fmt.Sprintf("%d", now.Add(-duration).UnixMilli())
		endTimestamp := fmt.Sprintf("%d", now.UnixMilli())
		bucketDuration := duration.Milliseconds()

		optsPriceChange := new(tair.ExTsAggregationArgs).New()
		optsPriceChange.Reverse()
		priceChangeResultMap[key] = pipe.ExTsRangeArgs(ctx,
			fmt.Sprintf(TimeSeriesQuotePriceKey, DefaultPkName),
			marketAddress,
			fromTimestamp,
			endTimestamp,
			optsPriceChange)

		optsSellVolume := new(tair.ExTsAggregationArgs).New()
		optsSellVolume.Reverse()
		optsSellVolume.MaxCount(1)
		optsSellVolume.Sum(bucketDuration)
		sellVolumeResultMap[key] = pipe.ExTsRangeArgs(ctx,
			fmt.Sprintf(TimeSeriesSellVolumeKey, DefaultPkName),
			marketAddress,
			fromTimestamp, endTimestamp, optsSellVolume)

		optsBuyVolume := new(tair.ExTsAggregationArgs).New()
		optsBuyVolume.Reverse()
		optsBuyVolume.MaxCount(1)
		optsBuyVolume.Sum(bucketDuration)
		buyVolumeResultMap[key] = pipe.ExTsRangeArgs(ctx,
			fmt.Sprintf(TimeSeriesBuyVolumeKey, DefaultPkName),
			marketAddress,
			fromTimestamp, endTimestamp, optsBuyVolume)

		optsSellCount := new(tair.ExTsAggregationArgs).New()
		optsSellCount.Reverse()
		optsSellCount.MaxCount(1)
		optsSellCount.Count(bucketDuration)
		sellCountResultMap[key] = pipe.ExTsRangeArgs(ctx,
			fmt.Sprintf(TimeSeriesSellVolumeKey, DefaultPkName),
			marketAddress,
			fromTimestamp,
			endTimestamp,
			optsSellCount)

		optsBuyCount := new(tair.ExTsAggregationArgs).New()
		optsBuyCount.Reverse()
		optsBuyCount.MaxCount(1)
		optsBuyCount.Count(bucketDuration)
		buyCountResultMap[key] = pipe.ExTsRangeArgs(ctx,
			fmt.Sprintf(TimeSeriesBuyVolumeKey, DefaultPkName),
			marketAddress,
			fromTimestamp, endTimestamp, optsBuyCount)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		t.Debugf("pipe.GetTsPeriodTimeData.Exec: %w", err) //ignore: key not exist error
	}

	result := make(map[string]*PeriodTimeData)
	for key := range PeriodDurations {
		priceChangeResult, ok := priceChangeResultMap[key]
		if !ok {
			t.Debugf("pipe.GetTsPeriodTimeData priceChangeResultMap not found: %v %v", key, marketAddress)
			continue
		}

		beforePrice := 1.0
		afterPrice := 1.0
		if len(priceChangeResult.DataPoints()) > 1 {
			beforePrice = priceChangeResult.DataPoints()[0].Value()
			afterPrice = priceChangeResult.DataPoints()[len(result)-1].Value()
		} else {
			t.Errorf("pipe.GetTsPeriodTimeData priceChangeResult.DataPoints() len <= 1: %v %v", key, marketAddress)
		}

		sellVolumeResult, ok := sellVolumeResultMap[key]
		if !ok {
			t.Debugf("pipe.GetTsPeriodTimeData sellVolumeResultMap not found: %v %v", key, marketAddress)
			continue
		}

		sellVolume := 0.0
		if len(sellVolumeResult.DataPoints()) > 0 {
			sellVolume = sellVolumeResult.DataPoints()[0].Value()
		}

		buyVolumeResult, ok := buyVolumeResultMap[key]
		if !ok {
			t.Debugf("pipe.GetTsPeriodTimeData buyVolumeResultMap not found: %v %v", key, marketAddress)
			continue
		}

		buyVolume := 0.0
		if len(buyVolumeResult.DataPoints()) > 0 {
			buyVolume = buyVolumeResult.DataPoints()[0].Value()
		}

		sellCountResult, ok := sellCountResultMap[key]
		if !ok {
			t.Debugf("pipe.GetTsPeriodTimeData sellCountResultMap not found: %v %v", key, marketAddress)
			continue
		}

		sellCount := 0.0
		if len(sellCountResult.DataPoints()) > 0 {
			sellCount = sellCountResult.DataPoints()[0].Value()
		}

		buyCountResult, ok := buyCountResultMap[key]
		if !ok {
			t.Debugf("pipe.GetTsPeriodTimeData buyCountResultMap not found: %v %v", key, marketAddress)
			continue
		}

		buyCount := 0.0
		if len(buyCountResult.DataPoints()) > 0 {
			buyCount = buyCountResult.DataPoints()[0].Value()
		}

		result[key] = &PeriodTimeData{
			PriceChange: t.calPriceChangePercent(afterPrice, beforePrice),
			SellVolume:  sellVolume,
			BuyVolume:   buyVolume,
			SellCount:   sellCount,
			BuyCount:    buyCount,
		}
	}

	return result, nil
}

func (t *TairTs) GetTsRealTimeData(ctx context.Context, marketAddress string, now time.Time) (*RealTimeData, error) {
	pipe := t.Client.TairPipeline()

	latestQuotePriceValueCmd := pipe.ExTsGet(ctx, fmt.Sprintf(TimeSeriesQuotePriceKey, DefaultPkName), marketAddress)
	baseVaultValueCmd := pipe.ExTsGet(ctx, fmt.Sprintf(TimeSeriesBaseVaultIncrByKey, DefaultPkName), marketAddress)
	quoteVaultValueCmd := pipe.ExTsGet(ctx, fmt.Sprintf(TimeSeriesQuoteVaultIncrByKey, DefaultPkName), marketAddress)
	sellCountCmd := pipe.ExTsGet(ctx, fmt.Sprintf(TimeSeriesSellCountIncrByKey, DefaultPkName), marketAddress)
	buyCountCmd := pipe.ExTsGet(ctx, fmt.Sprintf(TimeSeriesBuyCountIncrByKey, DefaultPkName), marketAddress)
	sellVolumeCmd := pipe.ExTsGet(ctx, fmt.Sprintf(TimeSeriesSellVolumeIncrByKey, DefaultPkName), marketAddress)
	buyVolumeCmd := pipe.ExTsGet(ctx, fmt.Sprintf(TimeSeriesBuyVolumeIncrByKey, DefaultPkName), marketAddress)

	_, err := pipe.Exec(ctx)
	if err != nil {
		logx.Errorf("pipe.GetTsRealTimeData.Exec:  %v", err) //ignore: key not exist error
	}

	latestQuotePriceResult, _ := latestQuotePriceValueCmd.Result()
	baseVaultResult, _ := baseVaultValueCmd.Result()
	quoteVaultResult, _ := quoteVaultValueCmd.Result()
	sellCountResult, _ := sellCountCmd.Result()
	buyCountResult, _ := buyCountCmd.Result()
	sellVolumeResult, _ := sellVolumeCmd.Result()
	buyVolumeResult, _ := buyVolumeCmd.Result()

	quotePrice := latestQuotePriceResult.Value()
	baseVault := baseVaultResult.Value()
	quoteVault := quoteVaultResult.Value()
	sellCount := sellCountResult.Value()
	buyCount := buyCountResult.Value()
	sellVolume := sellVolumeResult.Value()
	buyVolume := buyVolumeResult.Value()

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

func (t *TairTs) calPriceChangePercent(nowPrice, lastPrice float64) float64 {
	if lastPrice <= 0 {
		return 0
	}
	return (nowPrice - lastPrice) / lastPrice
}

func (t *TairTs) cmdMADD(ctx context.Context, pipe tair.TairPipeline, pKey string, sKeys []*tair.ExTsDataPoint) *redis.StringSliceCmd {
	a := make([]interface{}, 2)
	a[0] = "EXTS.S.MADD"
	a[1] = len(sKeys)
	args := tair.ExTsAttributeArgs{}.New()
	a = append(a, args.JoinArgs(pKey, sKeys)...)
	cmd := redis.NewStringSliceCmd(ctx, a...)
	_ = pipe.Process(ctx, cmd)
	return cmd
}
