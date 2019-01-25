package core

import (
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/merkletree"
)

// ProofOfClaimPartial is a proof of existence and non-existence of a claim in
// a single tree (only one level).
type ProofOfClaimPartial struct {
	Mtp0 *merkletree.Proof
	Mtp1 *merkletree.Proof
	Root *merkletree.Hash
	Aux  *SetRootAux
}

func (p *ProofOfClaimPartial) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Mtp0": hex.EncodeToString(p.Mtp0.Bytes()),
		"Mtp1": hex.EncodeToString(p.Mtp1.Bytes()),
		"Root": hex.EncodeToString(p.Root[:]),
		"Aux":  p.Aux,
	})
}

// SetRootAux is the auxiliary data to build the set root claim from a root in
// a partial proof of claim.
type SetRootAux struct {
	Version uint32
	Era     uint32
	EthAddr common.Address
}

// ProofOfClaim is a complete proof of a claim that includes all the proofs of
// existence and non-existence for mutliple levels from the leaf of a tree to
// the signed root of possibly another tree whose root.
type ProofOfClaim struct {
	Proofs    []ProofOfClaimPartial
	Leaf      *merkletree.Data
	Date      uint64
	Signature []byte // signature of the Root of the Relay
}

func (p *ProofOfClaim) MarshalJSON() ([]byte, error) {
	leafBytes := p.Leaf.Bytes()
	return json.Marshal(map[string]interface{}{
		"Proofs":    p.Proofs,
		"Leaf":      hex.EncodeToString(leafBytes[:]),
		"Date":      p.Date,
		"Signature": hex.EncodeToString(p.Signature),
	})
}
