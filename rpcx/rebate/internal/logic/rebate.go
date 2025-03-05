package logic

import (
	"context"
	"github.com/simance-ai/smdx/rpcx/account/account"
	"github.com/simance-ai/smdx/rpcx/order/order"
	"github.com/simance-ai/smdx/rpcx/rebate/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type RebateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRebateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RebateLogic {
	return &RebateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RebateLogic) Timeout() time.Duration {
	return time.Minute * 10
}

func (l *RebateLogic) Rebate() {
	logx.Info("开始返佣")

	// 获取返佣订单列表
	orderListResp, err := l.svcCtx.OrderClient.GetRebateOrderList(l.ctx, &order.RebateOrderListRequest{})
	if err != nil {
		logx.Error(err.Error())
	}

	// 返佣金额加入账户钱包，生成返佣记录
	_, err = l.svcCtx.AccountClient.Rebate(l.ctx, &account.RebateRequest{
		RebateOrderList: orderListResp.OrderList,
	})
	if err != nil {
		logx.Error("rebate account err: ", err.Error())
	}

	// 标记已返佣订单
	var ids []int64
	for _, rebateOrder := range orderListResp.OrderList {
		ids = append(ids, rebateOrder.Id)
	}
	_, err = l.svcCtx.OrderClient.UpdateRebateOrder(l.ctx, &order.UpdateRebateOrderRequest{
		OrderIdList: ids,
	})
	if err != nil {
		logx.Error("rebate order update err: ", err.Error())
	}
}
