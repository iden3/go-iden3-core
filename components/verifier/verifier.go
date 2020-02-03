package verifier

import (
	"fmt"
)

type Verifier struct {
}

func New() (*Verifier, error) {
	return nil, fmt.Errorf("TODO")
}

func (v *Verifier) VerifyCredentialExistence() error {
	return fmt.Errorf("TODO")
}

func (v *Verifier) VerifyCredentialValidity() error {
	return fmt.Errorf("TODO")
}
