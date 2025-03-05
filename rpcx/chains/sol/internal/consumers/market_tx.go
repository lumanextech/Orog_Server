package consumers

import (
	"context"
	"encoding/json"

	"github.com/simance-ai/smdx/pkg/common"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/dao/model"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type MarketTxLogic struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewMarketTxLogic(svcCtx *svc.ServiceContext) *MarketTxLogic {

	ctx := context.Background()
	logger := logx.WithContext(ctx).WithFields(logx.LogField{
		Key:   "Module",
		Value: "MarketTxLogic",
	})

	return &MarketTxLogic{
		svcCtx: svcCtx,
		ctx:    ctx,
		Logger: logger,
	}
}

// Consume  ConsumerHandler
func (l *MarketTxLogic) Consume(ctx context.Context, key, value string) error {

	logger := logx.WithContext(ctx).WithFields(logx.LogField{
		Key:   "Method",
		Value: "MarketTxConsume",
	})

	marketTxs := make([]*model.MarketTx, 0)
	err := json.Unmarshal([]byte(value), &marketTxs)
	if err != nil {
		logger.Error("json.Unmarshal marketTx: ", err)
		return nil
	}

	for _, marketTx := range marketTxs {
		if marketTx.CreatedTimestamp.IsZero() {
			logger.Error("marketTx.CreatedTimestamp nil empty")
			return nil
		}

		if marketTx.QuotePrice <= 0 && !common.IsLiquidityType(marketTx.TxType) {
			logger.Errorf("marketTx.QuotePrice not a liquidity type and price is empty, hash: %v key: %v", marketTx.TxHash, key)
			return nil
		}

		if marketTx.MarketAddress == "" {
			logger.Error("marketTx.MarketAddress nil empty")
			return nil
		}

		if marketTx.TxType == common.TxUnknown {
			logger.Error("marketTx.TxType Unknown")
			return nil
		}
	}

	return nil
}
