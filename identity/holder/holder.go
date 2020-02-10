package holder

import (
	"fmt"

	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/identity/issuer"
)

// Holder is an identity that holds claims.  It is an extension of an Issuer.
type Holder struct {
	*issuer.Issuer
}

// New creates a new Holder, calling the internal Issuer.New().
func New() (*Holder, error) {
	return nil, fmt.Errorf("TODO")
}

// New creates a Holder by loading a previously created Holder (with New, and calling the internal Issuer.Load().
func Load() (*Holder, error) {
	return nil, fmt.Errorf("TODO")
}

// HolderGetCredentialValidity gets a Credential of Validity from a Credential
// of Existence.  This requires a request to the Issuer IdenStatePubOffChain.
func (h *Holder) HolderGetCredentialValidity(credentialExistence *proof.CredentialExistence) (*proof.CredentialValidity, error) {
	return nil, fmt.Errorf("TODO")
}

// HolderImportCredentialExistence imports a received Credential of Existence into the ClaimsDB.
func (h *Holder) HolderImportCredentialExistence(credentialExistence *proof.ProofClaim) error {
	return fmt.Errorf("TODO")
}
