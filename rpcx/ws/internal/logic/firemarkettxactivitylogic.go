package logic

import (
	"context"
	xerror2 "github.com/simance-ai/smdx/pkg/errors/x_err"
	"github.com/simance-ai/smdx/rpcx/ws/internal/types"
	"github.com/simance-ai/smdx/rpcx/ws/internal/util"

	"github.com/simance-ai/smdx/rpcx/ws/internal/svc"
	"github.com/simance-ai/smdx/rpcx/ws/ws"

	"github.com/zeromicro/go-zero/core/logx"
)

type FireMarketTxActivityLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFireMarketTxActivityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FireMarketTxActivityLogic {
	return &FireMarketTxActivityLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FireMarketTxActivityLogic) FireMarketTxActivity(in *ws.FireMarketTxActivityRequest) (*ws.FireMarketTxActivityResponse, error) {

	txs := in.Txs
	marketAddress := in.MarketAddress
	chain := in.Chain
	channelKey, err := util.GetMarketTxActivityRedisKey(types.ActionSubscribe, types.ChannelMarketTxActivity, chain, marketAddress)
	if err != nil {
		return nil, xerror2.NewErrCodeMsg(xerror2.GetChannelKeyErr, err.Error())
	}

	err = l.svcCtx.WebsocketContext.FireJsonMessage(l.ctx, channelKey, types.ChannelMarketTxActivity, txs)
	if err != nil {
		return nil, xerror2.NewErrCodeMsg(xerror2.FireJsonMessageErr, err.Error())
	}

	return &ws.FireMarketTxActivityResponse{
		Message: "ok",
	}, nil
}
