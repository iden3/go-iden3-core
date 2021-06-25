package proof

import (

	// common3 "github.com/iden3/go-iden3-core/common"
	"fmt"

	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-merkletree"
	// "github.com/iden3/go-iden3-crypto/babyjub"
)

type IdenStateData struct {
	BlockTs   int64
	BlockN    uint64
	IdenState *merkletree.Hash
}

type CredentialExistence struct {
	Id                  *core.ID
	IdenStateData       IdenStateData
	MtpClaim            *merkletree.Proof
	Claim               *merkletree.Entry
	RevocationsTreeRoot *merkletree.Hash
	RootsTreeRoot       *merkletree.Hash
	IdenPubUrl          string
}

func (c CredentialExistence) String() string {
	type alias CredentialExistence
	return fmt.Sprintf("%+v", alias(c))
}

type CredentialValidity struct {
	CredentialExistence CredentialExistence
	IdenStateData       IdenStateData
	MtpNotNonce         *merkletree.Proof
	ClaimsTreeRoot      *merkletree.Hash
	RootsTreeRoot       *merkletree.Hash
}

func (c CredentialValidity) String() string {
	type alias CredentialValidity
	return fmt.Sprintf("%+v", alias(c))
}
