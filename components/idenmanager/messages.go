package idenmanager

import (
	"crypto/ecdsa"

	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/utils"
	"github.com/iden3/go-iden3-crypto/babyjub"
)

// BytesSignedMsg contains the value and its signature in Hex representation
type BytesSignedMsg struct {
	ValueHex  string                 `json:"valueHex" binding:"required"` // claim.Bytes() in a hex format
	Signature *utils.SignatureEthMsg `json:"signatureHex" binding:"required"`
	KSignPk   *utils.PublicKey       `json:"kSignPK" binding:"required"`
}

// ClaimBasicMsg contains a core.ClaimBasic with its signature in Hex
type ClaimBasicMsg struct {
	ClaimBasic core.ClaimBasic        `binding:"required"`
	Signature  *utils.SignatureEthMsg `binding:"required"`
}

// ClaimAssignNameMsg contains a core.ClaimAssignName with its signature in Hex
type ClaimAssignNameMsg struct {
	ClaimAssignName core.ClaimAssignName   `binding:"required"`
	Signature       *utils.SignatureEthMsg `binding:"required"`
}

// ClaimAuthorizeKSignSecp256k1Msg contains a core.ClaimAuthorizeKSignP256 with its signature in Hex
type ClaimAuthorizeKSignSecp256k1Msg struct {
	ClaimAuthorizeKSignSecp256k1 core.ClaimAuthorizeKSignSecp256k1 `binding:"required"`
	Signature                    *utils.SignatureEthMsg            `binding:"required"`
	KSignP256                    *ecdsa.PublicKey                  `binding:"required"`
}

// SetRootMsg contains the data to set the SetRootClaim with its signature in Hex
type SetRootMsg struct {
	Root      *merkletree.Hash       `binding:"required"`
	Id        *core.ID               `binding:"required"`
	KSignPk   *utils.PublicKey       `binding:"required"`
	Timestamp int64                  `binding:"required"`
	Signature *utils.SignatureEthMsg `binding:"required"`
}

// SetRoot0Req contains the data to set the SetRootClaim
type SetRoot0Req struct {
	OldRoot           *merkletree.Hash        `json:"oldRoot" binding:"required"`
	NewRoot           *merkletree.Hash        `json:"newRoot" binding:"required"`
	ClaimAuthKOp      *merkletree.Entry       `json:"claimKOp" binding:"required"`
	ProofClaimAuthKOp *core.ProofClaimGenesis `json:"proofKOp" binding:"required,dive"`
	Date              int64                   `json:"date" binding:"required"`
	Signature         *babyjub.SignatureComp  `json:"signature" binding:"required"` // signature of the Root
}

// ClaimValueMsg contains a core.ClaimValue with its signature in Hex
type ClaimValueMsg struct {
	ClaimValue merkletree.Entry       `binding:"required"`
	Signature  *utils.SignatureEthMsg `binding:"required"`
	KSignPk    *utils.PublicKey       `binding:"required"`
}

// TODO: Remove in next refactor
// ProofTreeLeaf contains all the parameters needed to proof that a Leaf is in a merkletree with a given Root
type ProofTreeLeaf struct {
	Leaf  []byte
	Proof []byte
	Root  merkletree.Hash
}

// TODO: Remove in next refactor
// ProofTreeLeafHex is the same data structure than ProofTreeLeaf but in Hexadecimal string representation
type ProofTreeLeafHex struct {
	Leaf  string
	Proof string
	Root  string
}

// TODO: Remove in next refactor
func (plh *ProofTreeLeafHex) Unhex() ProofTreeLeaf {
	var r ProofTreeLeaf
	r.Leaf, _ = common3.HexDecode(plh.Leaf)
	r.Proof, _ = common3.HexDecode(plh.Proof)
	rootBytes, _ := common3.HexDecode(plh.Root)
	copy(r.Root[:], rootBytes[:32])
	return r
}

// TODO: Remove in next refactor
// Hex returns a ProofTreeLeafHex data structure
func (pl *ProofTreeLeaf) Hex() ProofTreeLeafHex {
	r := ProofTreeLeafHex{
		common3.HexEncode(pl.Leaf),
		common3.HexEncode(pl.Proof),
		pl.Root.Hex(),
	}
	return r
}

// TODO: Remove in next refactor
// ProofClaimUser is the proof of a claim in the Identity MerkleTree, and the SetRootClaim of that MerkleTree inside the Relay's MerkleTree. Also with the proofs of non revocation of both claims
type ProofClaimUser struct {
	ClaimProof                     ProofTreeLeaf
	SetRootClaimProof              ProofTreeLeaf
	ClaimNonRevocationProof        ProofTreeLeaf
	SetRootClaimNonRevocationProof ProofTreeLeaf
	Date                           int64
	Signature                      []byte // signature of the Root of the Relay
}

// TODO: Remove in next refactor
type ProofClaimUserHex struct {
	ClaimProof                     ProofTreeLeafHex
	SetRootClaimProof              ProofTreeLeafHex
	ClaimNonRevocationProof        ProofTreeLeafHex
	SetRootClaimNonRevocationProof ProofTreeLeafHex
	Date                           int64
	Signature                      string // signature of the Root of the Relay
}

// TODO: Remove in next refactor
func (pc *ProofClaimUser) Hex() ProofClaimUserHex {
	r := ProofClaimUserHex{
		pc.ClaimProof.Hex(),
		pc.SetRootClaimProof.Hex(),
		pc.ClaimNonRevocationProof.Hex(),
		pc.SetRootClaimNonRevocationProof.Hex(),
		pc.Date,
		common3.HexEncode(pc.Signature),
	}
	return r
}

// TODO: Remove in next refactor
func (pch *ProofClaimUserHex) Unhex() (ProofClaimUser, error) {
	sigBytes, err := common3.HexDecode(pch.Signature)
	if err != nil {
		return ProofClaimUser{}, err
	}
	r := ProofClaimUser{
		pch.ClaimProof.Unhex(),
		pch.SetRootClaimProof.Unhex(),
		pch.ClaimNonRevocationProof.Unhex(),
		pch.SetRootClaimNonRevocationProof.Unhex(),
		pch.Date,
		sigBytes,
	}
	return r, nil
}
