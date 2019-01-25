package core

import (
	"encoding/hex"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/utils"
)

// ProofOfClaimPartial is a proof of existence and non-existence of a claim in
// a single tree (only one level).
type ProofOfClaimPartial struct {
	Mtp0 *merkletree.Proof `json:"mtp0" binding:"required"`
	Mtp1 *merkletree.Proof `json:"mtp1" binding:"required"`
	Root *merkletree.Hash  `json:"root" binding:"required"`
	Aux  *SetRootAux       `json:"aux" binding:"required"`
}

func (p *ProofOfClaimPartial) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"mtp0": hex.EncodeToString(p.Mtp0.Bytes()),
		"mtp1": hex.EncodeToString(p.Mtp1.Bytes()),
		"root": hex.EncodeToString(p.Root[:]),
		"aux":  p.Aux,
	})
}

// SetRootAux is the auxiliary data to build the set root claim from a root in
// a partial proof of claim.
type SetRootAux struct {
	Version uint32         `json:"version" binding:"required"`
	Era     uint32         `json:"era" binding:"required"`
	EthAddr common.Address `json:"ethAddr" binding:"required"`
}

// ProofOfClaim is a complete proof of a claim that includes all the proofs of
// existence and non-existence for mutliple levels from the leaf of a tree to
// the signed root of possibly another tree whose root.
type ProofOfClaim struct {
	Proofs    []ProofOfClaimPartial  `json:"proofs" binding:"required"`
	Leaf      *merkletree.Data       `json:"leaf" binding:"required"`
	Date      uint64                 `json:"date" binding:"required"`
	Signature *utils.SignatureEthMsg `json:"signature" binding:"required"` // signature of the Root of the Relay
}

func (p *ProofOfClaim) MarshalJSON() ([]byte, error) {
	leafBytes := p.Leaf.Bytes()
	return json.Marshal(map[string]interface{}{
		"proofs":    p.Proofs,
		"leaf":      hex.EncodeToString(leafBytes[:]),
		"date":      p.Date,
		"signature": p.Signature,
	})
}
