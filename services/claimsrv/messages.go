package claimsrv

import (
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/merkletree"
)

// BytesSignedMsg contains the value and its signature in Hex representation
type BytesSignedMsg struct {
	ValueHex     string `json:"valueHex"` // claim.Bytes() in a hex format
	SignatureHex string `json:"signatureHex"`
}

// ClaimDefaultMsg contains a core.ClaimDefault with its signature in Hex
type ClaimDefaultMsg struct {
	ClaimDefault core.ClaimDefault
	Signature    string
}

// AssignNameClaimMsg contains a core.AssignNameClaim with its signature in Hex
type AssignNameClaimMsg struct {
	AssignNameClaim core.AssignNameClaim
	Signature       string
}

// AuthorizeKSignClaimMsg contains a core.AuthorizeKSignClaim with its signature in Hex
type AuthorizeKSignClaimMsg struct {
	AuthorizeKSignClaim core.AuthorizeKSignClaim
	Signature           string
}

// SetRootClaimMsg contains a core.SetRootClaim with its signature in Hex
type SetRootClaimMsg struct {
	SetRootClaim core.SetRootClaim
	Signature    string
}

// ClaimValueMsg contains a core.ClaimValue with its signature in Hex
type ClaimValueMsg struct {
	ClaimValue merkletree.Value
	Signature  string
}

// ProofOfTreeLeaf contains all the parameters needed to proof that a Leaf is in a merkletree with a given Root
type ProofOfTreeLeaf struct {
	Leaf  []byte
	Hi    merkletree.Hash
	Proof []byte
	Root  merkletree.Hash
}

// ProofOfTreeLeafHex is the same data structure than ProofOfTreeLeaf but in Hexadecimal string representation
type ProofOfTreeLeafHex struct {
	Leaf  string
	Hi    string
	Proof string
	Root  string
}

func (plh *ProofOfTreeLeafHex) Unhex() ProofOfTreeLeaf {
	var r ProofOfTreeLeaf
	r.Leaf, _ = common3.HexToBytes(plh.Leaf)
	hiBytes, _ := common3.HexToBytes(plh.Hi)
	copy(r.Hi[:], hiBytes[:32])
	r.Proof, _ = common3.HexToBytes(plh.Proof)
	rootBytes, _ := common3.HexToBytes(plh.Root)
	copy(r.Root[:], rootBytes[:32])
	return r
}

// Hex returns a ProofOfTreeLeafHex data structure
func (pl *ProofOfTreeLeaf) Hex() ProofOfTreeLeafHex {
	r := ProofOfTreeLeafHex{
		common3.BytesToHex(pl.Leaf),
		pl.Hi.Hex(),
		common3.BytesToHex(pl.Proof),
		pl.Root.Hex(),
	}
	return r
}

type ProofOfClaim struct {
	ClaimProof                     ProofOfTreeLeaf
	SetRootClaimProof              ProofOfTreeLeaf
	ClaimNonRevocationProof        ProofOfTreeLeaf
	SetRootClaimNonRevocationProof ProofOfTreeLeaf
	Signature                      []byte // signature of the Root of the Relay
}
type ProofOfClaimHex struct {
	ClaimProof                     ProofOfTreeLeafHex
	SetRootClaimProof              ProofOfTreeLeafHex
	ClaimNonRevocationProof        ProofOfTreeLeafHex
	SetRootClaimNonRevocationProof ProofOfTreeLeafHex
	Signature                      string // signature of the Root of the Relay
}

func (pc *ProofOfClaim) Hex() ProofOfClaimHex {
	r := ProofOfClaimHex{
		pc.ClaimProof.Hex(),
		pc.SetRootClaimProof.Hex(),
		pc.ClaimNonRevocationProof.Hex(),
		pc.SetRootClaimNonRevocationProof.Hex(),
		common3.BytesToHex(pc.Signature),
	}
	return r
}
func (pch *ProofOfClaimHex) Unhex() (ProofOfClaim, error) {
	sigBytes, err := common3.HexToBytes(pch.Signature)
	if err != nil {
		return ProofOfClaim{}, err
	}
	r := ProofOfClaim{
		pch.ClaimProof.Unhex(),
		pch.SetRootClaimProof.Unhex(),
		pch.ClaimNonRevocationProof.Unhex(),
		pch.SetRootClaimNonRevocationProof.Unhex(),
		sigBytes,
	}
	return r, nil
}
