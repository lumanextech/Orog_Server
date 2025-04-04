// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package raydium_clmm

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// Update pool status for given vaule
//
// # Arguments
//
// * `ctx`- The context of accounts
// * `status` - The vaule of status
type UpdatePoolStatus struct {
	Status *uint8

	// [0] = [SIGNER] authority
	//
	// [1] = [WRITE] poolState
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewUpdatePoolStatusInstructionBuilder creates a new `UpdatePoolStatus` instruction builder.
func NewUpdatePoolStatusInstructionBuilder() *UpdatePoolStatus {
	nd := &UpdatePoolStatus{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 2),
	}
	return nd
}

// SetStatus sets the "status" parameter.
func (inst *UpdatePoolStatus) SetStatus(status uint8) *UpdatePoolStatus {
	inst.Status = &status
	return inst
}

// SetAuthorityAccount sets the "authority" account.
func (inst *UpdatePoolStatus) SetAuthorityAccount(authority ag_solanago.PublicKey) *UpdatePoolStatus {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(authority).SIGNER()
	return inst
}

// GetAuthorityAccount gets the "authority" account.
func (inst *UpdatePoolStatus) GetAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetPoolStateAccount sets the "poolState" account.
func (inst *UpdatePoolStatus) SetPoolStateAccount(poolState ag_solanago.PublicKey) *UpdatePoolStatus {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(poolState).WRITE()
	return inst
}

// GetPoolStateAccount gets the "poolState" account.
func (inst *UpdatePoolStatus) GetPoolStateAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

func (inst UpdatePoolStatus) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_UpdatePoolStatus,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst UpdatePoolStatus) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *UpdatePoolStatus) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Status == nil {
			return errors.New("Status parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Authority is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.PoolState is not set")
		}
	}
	return nil
}

func (inst *UpdatePoolStatus) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("UpdatePoolStatus")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=1]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("Status", *inst.Status))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=2]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("authority", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("poolState", inst.AccountMetaSlice.Get(1)))
					})
				})
		})
}

func (obj UpdatePoolStatus) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Status` param:
	err = encoder.Encode(obj.Status)
	if err != nil {
		return err
	}
	return nil
}
func (obj *UpdatePoolStatus) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Status`:
	err = decoder.Decode(&obj.Status)
	if err != nil {
		return err
	}
	return nil
}

// NewUpdatePoolStatusInstruction declares a new UpdatePoolStatus instruction with the provided parameters and accounts.
func NewUpdatePoolStatusInstruction(
	// Parameters:
	status uint8,
	// Accounts:
	authority ag_solanago.PublicKey,
	poolState ag_solanago.PublicKey) *UpdatePoolStatus {
	return NewUpdatePoolStatusInstructionBuilder().
		SetStatus(status).
		SetAuthorityAccount(authority).
		SetPoolStateAccount(poolState)
}
