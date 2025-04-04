// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package jupter

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// AldrinSwap is the `aldrinSwap` instruction.
type AldrinSwap struct {

	// [0] = [] swapProgram
	//
	// [1] = [] pool
	//
	// [2] = [] poolSigner
	//
	// [3] = [WRITE] poolMint
	//
	// [4] = [WRITE] baseTokenVault
	//
	// [5] = [WRITE] quoteTokenVault
	//
	// [6] = [WRITE] feePoolTokenAccount
	//
	// [7] = [] walletAuthority
	//
	// [8] = [WRITE] userBaseTokenAccount
	//
	// [9] = [WRITE] userQuoteTokenAccount
	//
	// [10] = [] tokenProgram
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewAldrinSwapInstructionBuilder creates a new `AldrinSwap` instruction builder.
func NewAldrinSwapInstructionBuilder() *AldrinSwap {
	nd := &AldrinSwap{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 11),
	}
	return nd
}

// SetSwapProgramAccount sets the "swapProgram" account.
func (inst *AldrinSwap) SetSwapProgramAccount(swapProgram ag_solanago.PublicKey) *AldrinSwap {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(swapProgram)
	return inst
}

// GetSwapProgramAccount gets the "swapProgram" account.
func (inst *AldrinSwap) GetSwapProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetPoolAccount sets the "pool" account.
func (inst *AldrinSwap) SetPoolAccount(pool ag_solanago.PublicKey) *AldrinSwap {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(pool)
	return inst
}

// GetPoolAccount gets the "pool" account.
func (inst *AldrinSwap) GetPoolAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetPoolSignerAccount sets the "poolSigner" account.
func (inst *AldrinSwap) SetPoolSignerAccount(poolSigner ag_solanago.PublicKey) *AldrinSwap {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(poolSigner)
	return inst
}

// GetPoolSignerAccount gets the "poolSigner" account.
func (inst *AldrinSwap) GetPoolSignerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetPoolMintAccount sets the "poolMint" account.
func (inst *AldrinSwap) SetPoolMintAccount(poolMint ag_solanago.PublicKey) *AldrinSwap {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(poolMint).WRITE()
	return inst
}

// GetPoolMintAccount gets the "poolMint" account.
func (inst *AldrinSwap) GetPoolMintAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

// SetBaseTokenVaultAccount sets the "baseTokenVault" account.
func (inst *AldrinSwap) SetBaseTokenVaultAccount(baseTokenVault ag_solanago.PublicKey) *AldrinSwap {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(baseTokenVault).WRITE()
	return inst
}

// GetBaseTokenVaultAccount gets the "baseTokenVault" account.
func (inst *AldrinSwap) GetBaseTokenVaultAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(4)
}

// SetQuoteTokenVaultAccount sets the "quoteTokenVault" account.
func (inst *AldrinSwap) SetQuoteTokenVaultAccount(quoteTokenVault ag_solanago.PublicKey) *AldrinSwap {
	inst.AccountMetaSlice[5] = ag_solanago.Meta(quoteTokenVault).WRITE()
	return inst
}

// GetQuoteTokenVaultAccount gets the "quoteTokenVault" account.
func (inst *AldrinSwap) GetQuoteTokenVaultAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(5)
}

// SetFeePoolTokenAccountAccount sets the "feePoolTokenAccount" account.
func (inst *AldrinSwap) SetFeePoolTokenAccountAccount(feePoolTokenAccount ag_solanago.PublicKey) *AldrinSwap {
	inst.AccountMetaSlice[6] = ag_solanago.Meta(feePoolTokenAccount).WRITE()
	return inst
}

// GetFeePoolTokenAccountAccount gets the "feePoolTokenAccount" account.
func (inst *AldrinSwap) GetFeePoolTokenAccountAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(6)
}

// SetWalletAuthorityAccount sets the "walletAuthority" account.
func (inst *AldrinSwap) SetWalletAuthorityAccount(walletAuthority ag_solanago.PublicKey) *AldrinSwap {
	inst.AccountMetaSlice[7] = ag_solanago.Meta(walletAuthority)
	return inst
}

// GetWalletAuthorityAccount gets the "walletAuthority" account.
func (inst *AldrinSwap) GetWalletAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(7)
}

// SetUserBaseTokenAccountAccount sets the "userBaseTokenAccount" account.
func (inst *AldrinSwap) SetUserBaseTokenAccountAccount(userBaseTokenAccount ag_solanago.PublicKey) *AldrinSwap {
	inst.AccountMetaSlice[8] = ag_solanago.Meta(userBaseTokenAccount).WRITE()
	return inst
}

// GetUserBaseTokenAccountAccount gets the "userBaseTokenAccount" account.
func (inst *AldrinSwap) GetUserBaseTokenAccountAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(8)
}

// SetUserQuoteTokenAccountAccount sets the "userQuoteTokenAccount" account.
func (inst *AldrinSwap) SetUserQuoteTokenAccountAccount(userQuoteTokenAccount ag_solanago.PublicKey) *AldrinSwap {
	inst.AccountMetaSlice[9] = ag_solanago.Meta(userQuoteTokenAccount).WRITE()
	return inst
}

// GetUserQuoteTokenAccountAccount gets the "userQuoteTokenAccount" account.
func (inst *AldrinSwap) GetUserQuoteTokenAccountAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(9)
}

// SetTokenProgramAccount sets the "tokenProgram" account.
func (inst *AldrinSwap) SetTokenProgramAccount(tokenProgram ag_solanago.PublicKey) *AldrinSwap {
	inst.AccountMetaSlice[10] = ag_solanago.Meta(tokenProgram)
	return inst
}

// GetTokenProgramAccount gets the "tokenProgram" account.
func (inst *AldrinSwap) GetTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(10)
}

func (inst AldrinSwap) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_AldrinSwap,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst AldrinSwap) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *AldrinSwap) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.SwapProgram is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Pool is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.PoolSigner is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.PoolMint is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.BaseTokenVault is not set")
		}
		if inst.AccountMetaSlice[5] == nil {
			return errors.New("accounts.QuoteTokenVault is not set")
		}
		if inst.AccountMetaSlice[6] == nil {
			return errors.New("accounts.FeePoolTokenAccount is not set")
		}
		if inst.AccountMetaSlice[7] == nil {
			return errors.New("accounts.WalletAuthority is not set")
		}
		if inst.AccountMetaSlice[8] == nil {
			return errors.New("accounts.UserBaseTokenAccount is not set")
		}
		if inst.AccountMetaSlice[9] == nil {
			return errors.New("accounts.UserQuoteTokenAccount is not set")
		}
		if inst.AccountMetaSlice[10] == nil {
			return errors.New("accounts.TokenProgram is not set")
		}
	}
	return nil
}

func (inst *AldrinSwap) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("AldrinSwap")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=0]").ParentFunc(func(paramsBranch ag_treeout.Branches) {})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=11]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("    swapProgram", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("           pool", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("     poolSigner", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("       poolMint", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(ag_format.Meta(" baseTokenVault", inst.AccountMetaSlice.Get(4)))
						accountsBranch.Child(ag_format.Meta("quoteTokenVault", inst.AccountMetaSlice.Get(5)))
						accountsBranch.Child(ag_format.Meta("   feePoolToken", inst.AccountMetaSlice.Get(6)))
						accountsBranch.Child(ag_format.Meta("walletAuthority", inst.AccountMetaSlice.Get(7)))
						accountsBranch.Child(ag_format.Meta("  userBaseToken", inst.AccountMetaSlice.Get(8)))
						accountsBranch.Child(ag_format.Meta(" userQuoteToken", inst.AccountMetaSlice.Get(9)))
						accountsBranch.Child(ag_format.Meta("   tokenProgram", inst.AccountMetaSlice.Get(10)))
					})
				})
		})
}

func (obj AldrinSwap) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	return nil
}
func (obj *AldrinSwap) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	return nil
}

// NewAldrinSwapInstruction declares a new AldrinSwap instruction with the provided parameters and accounts.
func NewAldrinSwapInstruction(
	// Accounts:
	swapProgram ag_solanago.PublicKey,
	pool ag_solanago.PublicKey,
	poolSigner ag_solanago.PublicKey,
	poolMint ag_solanago.PublicKey,
	baseTokenVault ag_solanago.PublicKey,
	quoteTokenVault ag_solanago.PublicKey,
	feePoolTokenAccount ag_solanago.PublicKey,
	walletAuthority ag_solanago.PublicKey,
	userBaseTokenAccount ag_solanago.PublicKey,
	userQuoteTokenAccount ag_solanago.PublicKey,
	tokenProgram ag_solanago.PublicKey) *AldrinSwap {
	return NewAldrinSwapInstructionBuilder().
		SetSwapProgramAccount(swapProgram).
		SetPoolAccount(pool).
		SetPoolSignerAccount(poolSigner).
		SetPoolMintAccount(poolMint).
		SetBaseTokenVaultAccount(baseTokenVault).
		SetQuoteTokenVaultAccount(quoteTokenVault).
		SetFeePoolTokenAccountAccount(feePoolTokenAccount).
		SetWalletAuthorityAccount(walletAuthority).
		SetUserBaseTokenAccountAccount(userBaseTokenAccount).
		SetUserQuoteTokenAccountAccount(userQuoteTokenAccount).
		SetTokenProgramAccount(tokenProgram)
}
