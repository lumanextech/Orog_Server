package logic

import (
	"context"

	"github.com/simance-ai/smdx/rpcx/ws/internal/svc"
	"github.com/simance-ai/smdx/rpcx/ws/internal/types"
	"github.com/simance-ai/smdx/rpcx/ws/internal/util"
	"github.com/simance-ai/smdx/rpcx/ws/ws"

	"github.com/zeromicro/go-zero/core/logx"
)

type FireMarketKlineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFireMarketKlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FireMarketKlineLogic {
	return &FireMarketKlineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FireMarketKline 发送K线数据
func (l *FireMarketKlineLogic) FireMarketKline(in *ws.FireMarketKlineRequest) (*ws.FireMarketKlineResponse, error) {

	logger := l.WithFields(logx.LogField{
		Key:   "Method",
		Value: "FireMarketKline",
	})

	for _, k := range in.Klines {
		channelKey, err := util.GetMarketKlineRedisKey(types.ActionSubscribe, types.ChannelMarketKline, k.Chain, k.MarketAddress, k.Interval)
		if err != nil {
			logger.Error("GetMarketKlineRedisKey: ", err)
			continue
		}

		err = l.svcCtx.WebsocketContext.FireJsonMessage(l.ctx, channelKey, types.ChannelMarketKline, k)
		if err != nil {
			logger.Error("FireJsonMessage: ", err)
			continue
		}
	}
	return &ws.FireMarketKlineResponse{
		Message: "ok",
	}, nil
}
