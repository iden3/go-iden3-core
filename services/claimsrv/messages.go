package claimsrv

import (
	"github.com/ethereum/go-ethereum/common"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/merkletree"
)

// BytesSignedMsg contains the value and its signature in Hex representation
type BytesSignedMsg struct {
	ValueHex     string         `json:"valueHex"` // claim.Bytes() in a hex format
	SignatureHex string         `json:"signatureHex"`
	KSign        common.Address `json:"ksign"`
}

// ClaimBasicMsg contains a core.ClaimBasic with its signature in Hex
type ClaimBasicMsg struct {
	ClaimBasic core.ClaimBasic
	Signature  string
}

// ClaimAssignNameMsg contains a core.ClaimAssignName with its signature in Hex
type ClaimAssignNameMsg struct {
	ClaimAssignName core.ClaimAssignName
	Signature       string
}

// ClaimAuthorizeKSignMsg contains a core.AuthorizeKSignClaim with its signature in Hex
type ClaimAuthorizeKSignMsg struct {
	ClaimAuthorizeKSign core.ClaimAuthorizeKSign
	Signature           string
	KSign               common.Address
}

// SetRootMsg contains the data to set the SetRootClaim with its signature in Hex
type SetRootMsg struct {
	Root      string
	IdAddr    string
	KSign     string
	Timestamp uint64
	Signature string
}

// ClaimValueMsg contains a core.ClaimValue with its signature in Hex
type ClaimValueMsg struct {
	ClaimValue merkletree.Entry
	Signature  string
	KSign      common.Address
}

// ProofOfTreeLeaf contains all the parameters needed to proof that a Leaf is in a merkletree with a given Root
type ProofOfTreeLeaf struct {
	Leaf  []byte
	Proof []byte
	Root  merkletree.Hash
}

// ProofOfTreeLeafHex is the same data structure than ProofOfTreeLeaf but in Hexadecimal string representation
type ProofOfTreeLeafHex struct {
	Leaf  string
	Proof string
	Root  string
}

func (plh *ProofOfTreeLeafHex) Unhex() ProofOfTreeLeaf {
	var r ProofOfTreeLeaf
	r.Leaf, _ = common3.HexToBytes(plh.Leaf)
	r.Proof, _ = common3.HexToBytes(plh.Proof)
	rootBytes, _ := common3.HexToBytes(plh.Root)
	copy(r.Root[:], rootBytes[:32])
	return r
}

// Hex returns a ProofOfTreeLeafHex data structure
func (pl *ProofOfTreeLeaf) Hex() ProofOfTreeLeafHex {
	r := ProofOfTreeLeafHex{
		common3.BytesToHex(pl.Leaf),
		common3.BytesToHex(pl.Proof),
		pl.Root.Hex(),
	}
	return r
}

// ProofOfClaim is the proof of a claim in the Identity MerkleTree, and the SetRootClaim of that MerkleTree inside the Relay's MerkleTree. Also with the proofs of non revocation of both claims
type ProofOfClaim struct {
	ClaimProof                     ProofOfTreeLeaf
	SetRootClaimProof              ProofOfTreeLeaf
	ClaimNonRevocationProof        ProofOfTreeLeaf
	SetRootClaimNonRevocationProof ProofOfTreeLeaf
	Date                           uint64
	Signature                      []byte // signature of the Root of the Relay
}
type ProofOfClaimHex struct {
	ClaimProof                     ProofOfTreeLeafHex
	SetRootClaimProof              ProofOfTreeLeafHex
	ClaimNonRevocationProof        ProofOfTreeLeafHex
	SetRootClaimNonRevocationProof ProofOfTreeLeafHex
	Date                           uint64
	Signature                      string // signature of the Root of the Relay
}

func (pc *ProofOfClaim) Hex() ProofOfClaimHex {
	r := ProofOfClaimHex{
		pc.ClaimProof.Hex(),
		pc.SetRootClaimProof.Hex(),
		pc.ClaimNonRevocationProof.Hex(),
		pc.SetRootClaimNonRevocationProof.Hex(),
		pc.Date,
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
		pch.Date,
		sigBytes,
	}
	return r, nil
}

// ProofOfRelayClaim is the proof of a claim in the Relay's MerkleTree, and the proof of non revocation of the claim
type ProofOfRelayClaim struct {
	ClaimProof              ProofOfTreeLeaf
	ClaimNonRevocationProof ProofOfTreeLeaf
	Date                    uint64
	Signature               []byte // signature of the Root of the Relay
}

type ProofOfRelayClaimHex struct {
	ClaimProof              ProofOfTreeLeafHex
	ClaimNonRevocationProof ProofOfTreeLeafHex
	Date                    uint64
	Signature               string // signature of the Root of the Relay
}

func (pc *ProofOfRelayClaim) Hex() ProofOfRelayClaimHex {
	r := ProofOfRelayClaimHex{
		pc.ClaimProof.Hex(),
		pc.ClaimNonRevocationProof.Hex(),
		pc.Date,
		common3.BytesToHex(pc.Signature),
	}
	return r
}
func (pch *ProofOfRelayClaimHex) Unhex() (ProofOfRelayClaim, error) {
	sigBytes, err := common3.HexToBytes(pch.Signature)
	if err != nil {
		return ProofOfRelayClaim{}, err
	}
	r := ProofOfRelayClaim{
		pch.ClaimProof.Unhex(),
		pch.ClaimNonRevocationProof.Unhex(),
		pch.Date,
		sigBytes,
	}
	return r, nil
}
