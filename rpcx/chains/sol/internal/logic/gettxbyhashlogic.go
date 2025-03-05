package logic

import (
	"context"

	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/svc"
	"github.com/simance-ai/smdx/rpcx/chains/sol/sol"

	"github.com/zeromicro/go-zero/core/logx"

	chainscommontx "github.com/simance-ai/smdx/rpcx/chains/common/tx"
)

type GetTxByHashLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTxByHashLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTxByHashLogic {
	return &GetTxByHashLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTxByHashLogic) GetTxByHash(in *sol.GetTxByHashRequest) (*sol.Tx, error) {

	return &sol.Tx{
		Result: &chainscommontx.Tx{
			MarketAddress:    "xx",
			TxHash:           "",
			BaseAmount:       0,
			QuoteAmount:      0,
			BaseSymbol:       "",
			QuoteSymbol:      "",
			BasePrice:        0,
			QuotePrice:       0,
			CreatedTimestamp: 0,
			BlockHeight:      0,
			TxIndex:          0,
			SenderAddress:    "",
			RecipientAddress: "",
		},
	}, nil
}
