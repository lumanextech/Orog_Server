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

type MarketTxActivityConsumer struct {
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarketTxActivityConsumer(svcCtx *svc.ServiceContext) *MarketTxActivityConsumer {
	return &MarketTxActivityConsumer{
		svcCtx: svcCtx,
		Logger: logx.WithContext(context.Background()),
	}
}

func (l *MarketTxActivityConsumer) Consume(ctx context.Context, key, value string) error {

	l.Infof("Key Market Tx: %s ", key)

	var mts []types.MarketTxV2
	err := json.Unmarshal([]byte(value), &mts)
	if err != nil {
		return fmt.Errorf("unmarshal MarketTxV2: %w", err)
	}

	if len(mts) == 0 {
		l.Infof("Empty MarketTxV2")
		return nil
	}

	var mt = mts[0]

	channelKey, err := util.GetMarketTxActivityRedisKey(types.ActionSubscribe, types.ChannelMarketTxActivity, mt.Chain, mt.MarketAddress)
	if err != nil {
		l.Error("GetMarketTxActivityRedisKey: ", err)
		return err
	}

	err = l.svcCtx.WebsocketContext.FireJsonMessage(ctx, channelKey, types.ChannelMarketTxActivity, mts)
	if err != nil {
		l.Error("FireJsonMessage: ", err)
		return err
	}

	return nil
}
