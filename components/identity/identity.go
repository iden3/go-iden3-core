package identity

import (
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/merkletree"
)

type Issuer interface {
	ID() *core.ID
	GenCredentialExistence(claim merkletree.Entrier) (*proof.CredentialExistence, error)
	IssueClaim(claim merkletree.Entrier) error
	PublishState() error
	RevokeClaim(claim merkletree.Entrier) error
	UpdateClaim(hIndex *merkletree.Hash, value []merkletree.ElemBytes) error
	Sign(string) (string, error)
	SignBinary(string) (string, error)
}

type Holder interface {
	HolderGetCredentialValidity(credentialExistence *proof.CredentialExistence) (*proof.CredentialValidity, error)
	HolderImportCredentialExistence(credentialExistence *proof.CredentialExistence) error
}

type IssuerHolder interface {
	Issuer
	Holder
}
