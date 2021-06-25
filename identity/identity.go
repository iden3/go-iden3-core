package identity

import (
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/merkletree_old"
	"github.com/iden3/go-merkletree"
)

// Issuer is an interface of an Identity that is only capable of issuing
// claims.  The identity can be set up without access to the IdenStates Smart
// Contract, in which case it will be a Genesis Only Identity and Identity
// update functions should fail.
type Issuer interface {
	ID() *core.ID
	GenCredentialExistence(claim merkletree_old.Entrier) (*proof.CredentialExistence, error)
	IssueClaim(claim merkletree_old.Entrier) error
	PublishState() error
	RevokeClaim(claim merkletree_old.Entrier) error
	UpdateClaim(hIndex *merkletree.Hash, value []merkletree_old.ElemBytes) error
	Sign(string) (string, error)
	SignBinary(string) (string, error)
}

// Holder is an interface of an Identity that is only capable of holding
// received claims.  Usually this interface is never used because the minimum
// Identity should be able to act as a Genesis Only Identity.  This interface
// is defined to be used in the IssuerHolder interface.
type Holder interface {
	HolderGetCredentialValidity(credentialExistence *proof.CredentialExistence) (*proof.CredentialValidity, error)
	HolderImportCredentialExistence(credentialExistence *proof.CredentialExistence) error
}

// IssuerHolder is an interface of an Identity capable of issuing claims and
// holding claims.  It combines the Issuer and Holder interfaces.
type IssuerHolder interface {
	Issuer
	Holder
}
