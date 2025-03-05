package order_v1

import (
	"context"
	"fmt"

	"github.com/simance-ai/smdx/app/internal/svc"
	"github.com/simance-ai/smdx/app/internal/types"
	"github.com/simance-ai/smdx/rpcx/order/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderRequest) (resp *types.CreateOrderResponse, err error) {
	if l.svcCtx == nil {
		return nil, fmt.Errorf("service context is not initialized")
	}
	if l.svcCtx.OrderClient == nil {
		return nil, fmt.Errorf("order client is not initialized")
	}

	address, ok := l.ctx.Value("payload").(string)
	if !ok {
		return nil, fmt.Errorf("failed to get address from context")
	}
	// Check if the account exists, and create it if it doesn't
	rpcResp, err := l.svcCtx.OrderClient.CreateOrder(l.ctx, &order.CreateOrderRequest{
		Address:         address,             // User address
		ChainName:       req.Chain,           // Chain name
		MarketAddress:   req.MarketAddress,   // Market address
		Side:            req.Side,            // Buy/Sell direction
		Type:            req.OrderType,       // Order type (market/limit)
		Price:           req.Px,              // Price
		Amount:          req.Sz,              // Quantity
		Slippage:        req.Slippage,        // Slippage
		TransactionHash: req.TransactionHash, // Transaction hash
	})
	if err != nil {
		return nil, err
	}

	// Return the response
	return &types.CreateOrderResponse{
		OrderId: rpcResp.OrderHash, // Convert Order ID to string
	}, nil
}
