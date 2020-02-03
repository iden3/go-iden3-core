package holder

import (
	"fmt"

	"github.com/iden3/go-iden3-core/components/identity/issuer"
	"github.com/iden3/go-iden3-core/core/proof"
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
