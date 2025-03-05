package ts_module

import (
	"context"
	"time"

	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/dao/model"
)

type Ts interface {
	CreateTsKeyIfNotExist(ctx context.Context, marketAddress string) error

	AddTsKeyPointsByTxs(ctx context.Context, marketAddress string, txs []*model.MarketTx) error

	GetTsPeriodTimeData(ctx context.Context, marketAddress string, now time.Time) (map[string]*PeriodTimeData, error)

	GetTsRealTimeData(ctx context.Context, marketAddress string, now time.Time) (*RealTimeData, error)

	GetTsKLineTimeData(ctx context.Context, marketAddress string, now time.Time) (map[string]*model.MarketKline1, error)
}
