package common

// ChainID chain id
const (
	SolChainId = "sol"
	EthChainId = "eth"
	BscChainId = "bsc"
)

func CheckIsSupportChain(chainId string) bool {
	switch chainId {
	case SolChainId, EthChainId, BscChainId:
		return true
	default:
		return false
	}
}
