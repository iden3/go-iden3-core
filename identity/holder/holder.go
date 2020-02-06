package holder

import (
	"fmt"

	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/identity/issuer"
)

type Holder struct {
	*issuer.Issuer
}

func New() (*Holder, error) {
	return nil, fmt.Errorf("TODO")
}

func Load() (*Holder, error) {
	return nil, fmt.Errorf("TODO")
}

func (h *Holder) HolderGetCredentialValidity(credentialExistence *proof.CredentialExistence) (*proof.CredentialValidity, error) {
	return nil, fmt.Errorf("TODO")
}

func (h *Holder) HolderImportCredentialExistence(credentialExistence *proof.ProofClaim) error {
	return fmt.Errorf("TODO")
}
