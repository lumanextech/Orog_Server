package common

const (
	//-1TxUnknown 0created, 1buy, 2sell, 3add_liquidity, 4remove_liquidity
	TxUnknown         int32 = -1
	TxCreated               = 0
	TxBuy                   = 1
	TxSell                  = 2
	TxAddLiquidity          = 3
	TxRemoveLiquidity       = 4
)

// IsSellType
// 买入交易-池中变化 quote减少 base增加
// 卖出交易-池中变化 quote增加 base减少
func IsSellType(txType int32) bool {
	return txType == TxSell
}

func IsUnknownType(txType int32) bool {
	return txType == TxUnknown
}

func IsLiquidityType(txType int32) bool {
	return txType == TxAddLiquidity ||
		txType == TxCreated ||
		txType == TxRemoveLiquidity
}
