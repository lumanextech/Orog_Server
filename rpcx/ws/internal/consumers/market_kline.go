package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/simance-ai/smdx/rpcx/ws/internal/svc"
	"github.com/simance-ai/smdx/rpcx/ws/internal/types"
	"github.com/simance-ai/smdx/rpcx/ws/internal/util"
	"github.com/zeromicro/go-zero/core/logx"
)

// MarketKlineConsumer
// 消费marketKline数据 并且推送到websocket的client
type MarketKlineConsumer struct {
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarketKlineConsumer(svcCtx *svc.ServiceContext) *MarketKlineConsumer {
	return &MarketKlineConsumer{
		svcCtx: svcCtx,
		Logger: logx.WithContext(context.Background()),
	}
}

func (l *MarketKlineConsumer) Consume(ctx context.Context, key, value string) error {

	l.Infof("Key Market Kline: %s", key)

	var mKLines []types.MarketKlineV2
	err := json.Unmarshal([]byte(value), &mKLines)
	if err != nil {
		return fmt.Errorf("unmarshal market_kline: %w", err)
	}

	for _, mKline := range mKLines {
		channelKey, err := util.GetMarketKlineRedisKey(types.ActionSubscribe, types.ChannelMarketKline,
			mKline.Chain, mKline.MarketAddress, mKline.Interval)
		if err != nil {
			l.Error("GetMarketKlineRedisKey: ", err)
			return err
		}

		err = l.svcCtx.WebsocketContext.FireJsonMessage(ctx, channelKey, types.ChannelMarketKline, mKline)
		if err != nil {
			l.Error("FireJsonMessage: ", err)
			return err
		}
	}

	return nil
}
