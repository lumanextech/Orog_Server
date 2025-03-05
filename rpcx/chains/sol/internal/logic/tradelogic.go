package logic

import (
	"context"
	"errors"
	xerror2 "github.com/simance-ai/smdx/pkg/errors/x_err"
	raydium_amm2 "github.com/simance-ai/smdx/rpcx/chains/sol/internal/common/contracts/raydium_amm"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/common/types"

	"github.com/deatil/go-cryptobin/cryptobin/crypto"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/simance-ai/smdx/pkg/common"
	"github.com/simance-ai/smdx/rpcx/account/account"
	"github.com/simance-ai/smdx/rpcx/chains/common/trade"
	"gorm.io/gorm"

	"os"

	"github.com/gagliardetto/solana-go"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/text"
	"github.com/shopspring/decimal"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/dao/dbx"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/dao/model"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/svc"
	"github.com/simance-ai/smdx/rpcx/chains/sol/sol"
	"github.com/zeromicro/go-zero/core/logx"
)

func init() {
	raydium_amm2.SetProgramID(types.RaydiumAmmMarket.ProgramID)
}

type TradeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TradeLogic {
	return &TradeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TradeLogic) Trade(in *sol.TradeRequest) (*sol.Response, error) {
	marketId := in.TradeRequest.MarketId
	accountId := in.TradeRequest.AccountId

	accountInfo, err := l.svcCtx.AccountClient.GetAccount(l.ctx, &account.AccountRequest{
		Id: accountId,
	})
	if err != nil {
		return nil, xerror2.NewErrCodeMsg(xerror2.GetAccountError, err.Error())
	}

	marketDB := dbx.Use(l.svcCtx.PgDB).Market

	//find market
	market, err := marketDB.WithContext(l.ctx).Where(marketDB.ID.Eq(marketId)).First()
	switch {
	case err == nil:
		err = l.handleMarketTrade(market, accountInfo, in.TradeRequest)
		switch {
		case err == nil:
		case errors.Is(err, xerror2.NewErrCode(xerror2.InsufficientBalance)):
			return nil, err
		default:
			return nil, xerror2.NewErrCodeMsg(xerror2.MarketTradeError, err.Error())
		}
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, xerror2.NewErrCode(xerror2.NotFoundMarket)
	default:
		return nil, xerror2.NewErrCodeMsg(xerror2.NotFoundMarket, err.Error())
	}

	return &sol.Response{
		Message: "success",
	}, nil
}

func (l *TradeLogic) handleMarketTrade(
	market *model.Market,
	accountInfo *account.AccountResponse,
	request *trade.TradeRequest) error {
	if market == nil || request == nil {
		return errors.New("market or request not found")
	}

	//build instruction for trade
	switch market.MarketType {
	case types.RaydiumAmmMarket.MarketTypeName:
		return l.handleRaydiumAmmTrade(market, accountInfo, request)
	default:
		return errors.New("market type not supported")
	}
}

func (l *TradeLogic) handleRaydiumAmmTrade(
	market *model.Market,
	accountInfo *account.AccountResponse,
	tradeRequest *trade.TradeRequest) error {

	logger := logx.WithContext(l.ctx).WithFields(logx.LogField{
		Key:   "Method",
		Value: "handleRaydiumAmmTrade",
	})

	if market == nil || tradeRequest == nil || accountInfo == nil {
		return errors.New("market or request or accountInfo not found")
	}

	//default gas
	unitLimit := uint32(0)
	unitPrice := uint64(0)
	rentLamport := uint64(2039280)
	rentSpace := uint64(165)
	switch tradeRequest.TradeMode {
	case common.TradeModeNormal:
		// normal mode
		unitLimit = 200000
		unitPrice = 1
	case common.TradeModeFaster:
		// fast mode
		unitLimit = 200000
		unitPrice = 14255778
	case common.TradeModePriority:
		// priority mode
		unitLimit = 200000
		unitPrice = 24255778
	default:
		return errors.New("mode not supported")
	}

	logger.Debug("accountInfo: ", accountInfo)

	if accountInfo.Pk == "" {
		return errors.New("accountInfo.Pk is empty")
	}

	decodePk := crypto.
		FromBase64String(accountInfo.Pk).
		SetKey(l.svcCtx.Config.PkEncode.Key).
		SetIv(l.svcCtx.Config.PkEncode.Iv).
		Aes().
		CBC().
		PKCS7Padding().
		Decrypt().
		ToString()

	if decodePk == "" {
		return errors.New("decodePk is empty")
	}

	player := solana.MustPrivateKeyFromBase58(decodePk)

	logger.Debug("player: ", player.PublicKey().String())

	switch tradeRequest.TradeType {
	case common.TradeBuy:
		//use sol -> token
		wsolAmountIn := uint64(tradeRequest.TradeBaseAmount * 1e9)

		return l.handleRaydiumAmmTradeBuy(logger, player, market,
			tradeRequest, wsolAmountIn, unitLimit, unitPrice, rentLamport, rentSpace)

	case common.TradeSell:
		//use token -> sol
		tokenAmountIn := decimal.NewFromFloat(float64(tradeRequest.TradeQuoteAmount)).
			Mul(decimal.New(1, int32(market.QuoteTokenDecimals))).IntPart()

		return l.handleRaydiumAmmTradeSell(logger, player, market,
			tradeRequest, uint64(tokenAmountIn), unitLimit, unitPrice, rentLamport, rentSpace)

	default:
		return errors.New("trade type not supported")
	}
}

func (l *TradeLogic) handleRaydiumAmmTradeBuy(logger logx.Logger,
	player solana.PrivateKey,
	market *model.Market,
	tradeRequest *trade.TradeRequest,
	wsolAmountIn uint64,
	unitLimit uint32, unitPrice uint64, rentLamport uint64, rentSpace uint64) error {

	logger = logx.WithContext(l.ctx).WithFields(logx.LogField{
		Key:   "Method2",
		Value: "handleRaydiumAmmTradeBuy",
	})

	//check balance
	balance, err := l.svcCtx.SolClient.GetBalance(l.ctx, player.PublicKey(), rpc.CommitmentConfirmed)
	if err != nil {
		return errors.New("raydium amm GetBalance error: " + err.Error())
	}

	logger.Debug("wSol balance: ", balance.Value, wsolAmountIn)

	if balance.Value < wsolAmountIn {
		return xerror2.NewErrCode(xerror2.InsufficientBalance)
	}

	minAmountOut := uint64(0)
	ammMarketAccount := solana.MustPublicKeyFromBase58(market.Address)
	quoteTokenMintAccount := solana.MustPublicKeyFromBase58(market.QuoteTokenMintAddress)
	baseTokenAccount := solana.MustPublicKeyFromBase58(market.BaseTokenAddress)
	quoteTokenAccount := solana.MustPublicKeyFromBase58(market.QuoteTokenAddress)

	//build computeUnits instruction
	//Compute Budget
	computeBudgetUnitLimit := computebudget.NewSetComputeUnitLimitInstructionBuilder(). //gas limit, 默认情况下，计算预算是200000个计算单元
												SetUnits(unitLimit).
												Build()

	computeBudgetUnitPrice := computebudget.NewSetComputeUnitPriceInstructionBuilder(). //gas price, 基本费用（5,000个Lamport）+ prioritization fee
												SetMicroLamports(unitPrice).
												Build()

	//find wSol account
	newWsolAccount, _, err := solana.FindAssociatedTokenAddress(player.PublicKey(), solana.WrappedSol)
	if err != nil {
		return err
	}

	createWsolAccountInstruction, err := associatedtokenaccount.
		NewCreateInstruction(
			player.PublicKey(),
			player.PublicKey(),
			solana.WrappedSol).ValidateAndBuild()
	if err != nil {
		return errors.New("raydium amm newWsolAccount error: " + err.Error())
	}

	transferInstruction, err := system.NewTransferInstruction(uint64(wsolAmountIn),
		player.PublicKey(), newWsolAccount).ValidateAndBuild()
	if err != nil {
		return errors.New("raydium amm tokenTransfer error: " + err.Error())
	}

	syncNativeInstruction, err := token.NewSyncNativeInstruction(newWsolAccount).ValidateAndBuild()
	if err != nil {
		return errors.New("raydium amm tokenSyncNative error: " + err.Error())
	}

	newQuoteTokenAccount, _, err := solana.FindAssociatedTokenAddress(player.PublicKey(), quoteTokenMintAccount)
	if err != nil {
		return errors.New("raydium amm quoteTokenMintAccount FindAssociatedTokenAddress error: " + err.Error())
	}

	var createTokenAccountInstruction *associatedtokenaccount.Instruction
	_, err = l.svcCtx.SolClient.GetAccountInfo(l.ctx, newQuoteTokenAccount)
	switch {
	case err == nil:
	case errors.Is(err, rpc.ErrNotFound):
		createTokenAccountInstruction, err = associatedtokenaccount.
			NewCreateInstruction(
				player.PublicKey(),
				player.PublicKey(),
				quoteTokenMintAccount).ValidateAndBuild()
		if err != nil {
			return errors.New("raydium amm newQuoteTokenAccount error: " + err.Error())
		}
	default:
		return errors.New("raydium amm quoteTokenMintAccount GetAccountInfo error: " + err.Error())
	}

	//build instruction for raydium swapBaseIn
	instruction := raydium_amm2.NewSwapBaseInInstructionBuilder()
	instruction.SetTokenProgramAccount(token.ProgramID)
	instruction.SetAmmAccount(ammMarketAccount)
	instruction.SetAmmAuthorityAccount(types.RaydiumAmmMarket.AuthorityProgramID)
	instruction.SetAmmOpenOrdersAccount(ammMarketAccount)
	instruction.SetAmmTargetOrdersAccount(ammMarketAccount)

	// 因表结构修改暂时替换成下面的set
	//if market.BaseTokenIsPcOrToken0 {
	//	instruction.SetPoolPcTokenAccountAccount(baseTokenAccount)
	//	instruction.SetPoolCoinTokenAccountAccount(quoteTokenAccount)
	//} else {
	//	instruction.SetPoolPcTokenAccountAccount(quoteTokenAccount)
	//	instruction.SetPoolCoinTokenAccountAccount(baseTokenAccount)
	//}
	instruction.SetPoolPcTokenAccountAccount(baseTokenAccount)
	instruction.SetPoolCoinTokenAccountAccount(quoteTokenAccount)

	instruction.SetSerumProgramAccount(ammMarketAccount)
	instruction.SetSerumMarketAccount(ammMarketAccount)
	instruction.SetSerumBidsAccount(ammMarketAccount)
	instruction.SetSerumAsksAccount(ammMarketAccount)
	instruction.SetSerumEventQueueAccount(ammMarketAccount)
	instruction.SetSerumCoinVaultAccountAccount(ammMarketAccount)
	instruction.SetSerumPcVaultAccountAccount(ammMarketAccount)
	instruction.SetSerumVaultSignerAccount(ammMarketAccount)

	instruction.SetAmountIn(wsolAmountIn)
	instruction.SetMinimumAmountOut(minAmountOut)

	instruction.SetUerSourceTokenAccountAccount(newWsolAccount)
	instruction.SetUerDestinationTokenAccountAccount(newQuoteTokenAccount)
	instruction.SetUserSourceOwnerAccount(player.PublicKey()) //payer
	ammSwapBaseInInstruction, err := instruction.ValidateAndBuild()
	if err != nil {
		return errors.New("raydium amm swapBaseIn error: " + err.Error())
	}

	closeAccountInstruction, err := token.NewCloseAccountInstruction(newWsolAccount, player.PublicKey(), player.PublicKey(), nil).ValidateAndBuild()
	if err != nil {
		return errors.New("raydium amm closeAccount error: " + err.Error())
	}

	recent, err := l.svcCtx.SolClient.GetRecentBlockhash(l.ctx, rpc.CommitmentFinalized)
	if err != nil {
		return errors.New("raydium amm GetRecentBlockhash error: " + err.Error())
	}

	tx := new(solana.Transaction)
	if createTokenAccountInstruction == nil {
		tx, err = solana.NewTransaction(
			[]solana.Instruction{
				computeBudgetUnitLimit,
				computeBudgetUnitPrice,
				createWsolAccountInstruction,
				transferInstruction,
				syncNativeInstruction,
				//createTokenAccountInstruction,
				ammSwapBaseInInstruction,
				closeAccountInstruction,
			},
			recent.Value.Blockhash,
			solana.TransactionPayer(player.PublicKey()),
			solana.TransactionAddressTables(
				map[solana.PublicKey]solana.PublicKeySlice{},
			),
		)
		if err != nil {
			return errors.New("raydium amm newTransaction error: " + err.Error())
		}
	} else {
		tx, err = solana.NewTransaction(
			[]solana.Instruction{
				computeBudgetUnitLimit,
				computeBudgetUnitPrice,
				createWsolAccountInstruction,
				transferInstruction,
				syncNativeInstruction,
				createTokenAccountInstruction,
				ammSwapBaseInInstruction,
				closeAccountInstruction,
			},
			recent.Value.Blockhash,
			solana.TransactionPayer(player.PublicKey()),
			solana.TransactionAddressTables(
				map[solana.PublicKey]solana.PublicKeySlice{},
			),
		)
		if err != nil {
			return errors.New("raydium amm newTransaction error: " + err.Error())
		}
	}

	tx.EncodeTree(text.NewTreeEncoder(os.Stdout, "Raydiumv4 Transaction"))

	//Sign Transaction
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(player.PublicKey()) {
			return &player
		}
		if key.Equals(newWsolAccount) {
			return &player
		}
		if key.Equals(newQuoteTokenAccount) {
			return &player
		}
		return nil
	})
	if err != nil {
		return errors.New("raydium amm sign error: " + err.Error())
	}

	signatureTx, err := l.svcCtx.SolClient.SendTransaction(l.ctx, tx)
	if err != nil {
		return errors.New("raydium amm sendTransaction error: " + err.Error())
	}

	logger.Info("raydium amm sendTransaction success", " txHash: ", signatureTx.String())

	return nil
}

func (l *TradeLogic) handleRaydiumAmmTradeSell(logger logx.Logger,
	player solana.PrivateKey,
	market *model.Market,
	tradeRequest *trade.TradeRequest,
	quoteAmountIn uint64,
	unitLimit uint32, unitPrice uint64, rentLamport uint64, rentSpace uint64) error {

	logger = logx.WithContext(l.ctx).WithFields(logx.LogField{
		Key:   "Method2",
		Value: "handleRaydiumAmmTradeSell",
	})

	minAmountOut := uint64(0)
	ammMarketAccount := solana.MustPublicKeyFromBase58(market.Address)
	quoteTokenMintAccount := solana.MustPublicKeyFromBase58(market.QuoteTokenMintAddress)
	baseTokenAccount := solana.MustPublicKeyFromBase58(market.BaseTokenAddress)
	quoteTokenAccount := solana.MustPublicKeyFromBase58(market.QuoteTokenAddress)

	//build computeUnits instruction
	//Compute Budget
	computeBudgetUnitLimit := computebudget.NewSetComputeUnitLimitInstructionBuilder(). //gas limit, 默认情况下，计算预算是200000个计算单元
												SetUnits(unitLimit).
												Build()

	computeBudgetUnitPrice := computebudget.NewSetComputeUnitPriceInstructionBuilder(). //gas price, 基本费用（5,000个Lamport）+ prioritization fee
												SetMicroLamports(unitPrice).
												Build()

	//find quote token account
	newQuoteAccount, _, err := solana.FindAssociatedTokenAddress(player.PublicKey(), quoteTokenMintAccount)
	if err != nil {
		return errors.New("newQuoteAccount error: " + err.Error())
	}

	//check balance
	newQuoteAccountBalanceResult, err := l.svcCtx.SolClient.GetTokenAccountBalance(l.ctx, newQuoteAccount, rpc.CommitmentConfirmed)
	if err != nil {
		return errors.New("GetBalance error: " + err.Error())
	}

	newQuoteAccountBalance, err := decimal.NewFromString(newQuoteAccountBalanceResult.Value.Amount)
	if err != nil {
		return errors.New("GetBalance error: " + err.Error())
	}
	logger.Debug("newQuoteAccountBalance: ", newQuoteAccountBalance.String(), " quoteAmountIn: ", quoteAmountIn)

	if newQuoteAccountBalance.LessThan(decimal.NewFromInt(int64(quoteAmountIn))) {
		return xerror2.NewErrCode(xerror2.InsufficientBalance)
	}

	newWsolAccount, _, err := solana.FindAssociatedTokenAddress(player.PublicKey(), solana.WrappedSol)
	if err != nil {
		return errors.New("raydium amm WrappedSol FindAssociatedTokenAddress error: " + err.Error())
	}

	var createWsolAccountInstruction *associatedtokenaccount.Instruction
	_, err = l.svcCtx.SolClient.GetAccountInfo(l.ctx, newWsolAccount)
	switch {
	case err == nil:
	case errors.Is(err, rpc.ErrNotFound):
		createWsolAccountInstruction, err = associatedtokenaccount.
			NewCreateInstruction(
				player.PublicKey(),
				player.PublicKey(),
				solana.WrappedSol).ValidateAndBuild()
		if err != nil {
			return errors.New("newWsolAccount error: " + err.Error())
		}
	default:
		return errors.New("quoteTokenMintAccount GetAccountInfo error: " + err.Error())
	}

	//build instruction for raydium swapBaseIn
	instruction := raydium_amm2.NewSwapBaseInInstructionBuilder()
	instruction.SetTokenProgramAccount(token.ProgramID)
	instruction.SetAmmAccount(ammMarketAccount)
	instruction.SetAmmAuthorityAccount(types.RaydiumAmmMarket.AuthorityProgramID)
	instruction.SetAmmOpenOrdersAccount(ammMarketAccount)
	instruction.SetAmmTargetOrdersAccount(ammMarketAccount)

	// 因表结构修改暂时替换成下面的set
	//if market.BaseTokenIsPcOrToken0 {
	//	instruction.SetPoolPcTokenAccountAccount(baseTokenAccount)
	//	instruction.SetPoolCoinTokenAccountAccount(quoteTokenAccount)
	//} else {
	//	instruction.SetPoolPcTokenAccountAccount(quoteTokenAccount)
	//	instruction.SetPoolCoinTokenAccountAccount(baseTokenAccount)
	//}
	instruction.SetPoolPcTokenAccountAccount(baseTokenAccount)
	instruction.SetPoolCoinTokenAccountAccount(quoteTokenAccount)

	instruction.SetSerumProgramAccount(ammMarketAccount)
	instruction.SetSerumMarketAccount(ammMarketAccount)
	instruction.SetSerumBidsAccount(ammMarketAccount)
	instruction.SetSerumAsksAccount(ammMarketAccount)
	instruction.SetSerumEventQueueAccount(ammMarketAccount)
	instruction.SetSerumCoinVaultAccountAccount(ammMarketAccount)
	instruction.SetSerumPcVaultAccountAccount(ammMarketAccount)
	instruction.SetSerumVaultSignerAccount(ammMarketAccount)

	instruction.SetAmountIn(quoteAmountIn)
	instruction.SetMinimumAmountOut(minAmountOut)

	instruction.SetUerSourceTokenAccountAccount(newQuoteAccount)
	instruction.SetUerDestinationTokenAccountAccount(newWsolAccount)
	instruction.SetUserSourceOwnerAccount(player.PublicKey()) //payer
	ammSwapBaseInInstruction, err := instruction.ValidateAndBuild()
	if err != nil {
		return errors.New("swapBaseIn error: " + err.Error())
	}

	closeAccountInstruction, err := token.NewCloseAccountInstruction(newWsolAccount, player.PublicKey(), player.PublicKey(), nil).ValidateAndBuild()
	if err != nil {
		return errors.New("closeAccount error: " + err.Error())
	}

	recent, err := l.svcCtx.SolClient.GetRecentBlockhash(l.ctx, rpc.CommitmentFinalized)
	if err != nil {
		return errors.New("GetRecentBlockhash error: " + err.Error())
	}

	tx := new(solana.Transaction)
	if createWsolAccountInstruction == nil {
		tx, err = solana.NewTransaction(
			[]solana.Instruction{
				computeBudgetUnitLimit,
				computeBudgetUnitPrice,
				//createWsolAccountInstruction,
				ammSwapBaseInInstruction,
				closeAccountInstruction,
			},
			recent.Value.Blockhash,
			solana.TransactionPayer(player.PublicKey()),
			solana.TransactionAddressTables(
				map[solana.PublicKey]solana.PublicKeySlice{},
			),
		)
		if err != nil {
			return errors.New("raydium amm newTransaction error: " + err.Error())
		}
	} else {
		tx, err = solana.NewTransaction(
			[]solana.Instruction{
				computeBudgetUnitLimit,
				computeBudgetUnitPrice,
				createWsolAccountInstruction,
				ammSwapBaseInInstruction,
				closeAccountInstruction,
			},
			recent.Value.Blockhash,
			solana.TransactionPayer(player.PublicKey()),
			solana.TransactionAddressTables(
				map[solana.PublicKey]solana.PublicKeySlice{},
			),
		)
		if err != nil {
			return errors.New("newTransaction error: " + err.Error())
		}
	}

	tx.EncodeTree(text.NewTreeEncoder(os.Stdout, "Raydiumv4 Transaction"))

	//Sign Transaction
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(player.PublicKey()) {
			return &player
		}
		if key.Equals(newWsolAccount) {
			return &player
		}
		if key.Equals(newQuoteAccount) {
			return &player
		}
		return nil
	})
	if err != nil {
		return errors.New("sign error: " + err.Error())
	}

	signatureTx, err := l.svcCtx.SolClient.SendTransaction(l.ctx, tx)
	if err != nil {
		return errors.New("sendTransaction error: " + err.Error())
	}

	logger.Info("sendTransaction success", " txHash: ", signatureTx.String())

	return nil
}
