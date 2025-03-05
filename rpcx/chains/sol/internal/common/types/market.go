package types

import (
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/shopspring/decimal"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/dao/model"
)

type MarketProgram struct {
	MarketTypeName     string
	ProgramID          solana.PublicKey
	AuthorityProgramID solana.PublicKey
}

type MarketType struct {
	MP                        *MarketProgram
	MainCompiledInstruction   solana.CompiledInstruction
	InnerCompiledInstructions []solana.CompiledInstruction
}

type MarketPoolInfo struct {
	MarketAddress solana.PublicKey
	MarketInfo    *model.Market
	Pool0         solana.PublicKey // pool0 base_mint_address(wsol)
	Pool1         solana.PublicKey
	Maker         solana.PublicKey // maker_address(signer)
}

type MarketPoolTokenBalance struct {
	PreMap  map[string]rpc.TokenBalance
	PostMap map[string]rpc.TokenBalance
}

func (m *MarketPoolTokenBalance) GetTokenBalanceChange(poolAccount solana.PublicKey) (*TokenBalanceChange, error) {
	preBalance, ok := m.PreMap[poolAccount.String()]
	if !ok {
		return nil, fmt.Errorf("preBalance not found, poolAccount: %s", poolAccount.String())
	}

	postBalance, ok := m.PostMap[poolAccount.String()]
	if !ok {
		return nil, fmt.Errorf("postBalance not found, poolAccount: %s", poolAccount.String())
	}

	preTokenAmount := preBalance.UiTokenAmount
	postTokenAmount := postBalance.UiTokenAmount
	preTokenAmountDecimal, err := decimal.NewFromString(preTokenAmount.Amount)
	if err != nil {
		return nil, err
	}
	postTokenAmountDecimal, err := decimal.NewFromString(postTokenAmount.Amount)
	if err != nil {
		return nil, err
	}
	preTokenUiAmountDecimal, err := decimal.NewFromString(preTokenAmount.UiAmountString)
	if err != nil {
		return nil, err
	}
	postTokenUiAmountDecimal, err := decimal.NewFromString(postTokenAmount.UiAmountString)
	if err != nil {
		return nil, err
	}
	changeAmountDecimal := postTokenAmountDecimal.Sub(preTokenAmountDecimal)
	changeUiAmountDecimal := postTokenUiAmountDecimal.Sub(preTokenUiAmountDecimal)
	uiAmountFloat := changeUiAmountDecimal.InexactFloat64()
	return &TokenBalanceChange{
		Mint: preBalance.Mint,
		UiTokenAmount: &rpc.UiTokenAmount{
			Amount:         changeAmountDecimal.String(),
			Decimals:       preBalance.UiTokenAmount.Decimals,
			UiAmount:       &uiAmountFloat,
			UiAmountString: changeUiAmountDecimal.String(),
		},
		IsNegative: changeAmountDecimal.IsNegative(),
	}, nil
}
