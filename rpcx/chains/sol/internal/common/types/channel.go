package types

import (
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/dao/model"
)

type MarketTxsSummary struct {
	MarketAddress string            `json:"market_address"`
	Txs           []*model.MarketTx `json:"txs"` //一段[]txs
}

type TokenBalanceChange struct {
	Mint          solana.PublicKey   `json:"mint"`
	UiTokenAmount *rpc.UiTokenAmount `json:"uiTokenAmount"`
	IsNegative    bool               `json:"isNegative"`
}
