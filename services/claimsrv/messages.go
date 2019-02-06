package claimsrv

import (
	"crypto/ecdsa"

	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/utils"
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

// ClaimAuthorizeKSignMsg contains a core.AuthorizeKSignClaim with its signature in Hex
type ClaimAuthorizeKSignMsg struct {
	ClaimAuthorizeKSign core.ClaimAuthorizeKSign `binding:"required"`
	Signature           *utils.SignatureEthMsg   `binding:"required"`
	KSignPk             *utils.PublicKey         `binding:"required"`
}

// ClaimAuthorizeKSignSecp256k1Msg contains a core.ClaimAuthorizeKSignP256 with its signature in Hex
type ClaimAuthorizeKSignSecp256k1Msg struct {
	ClaimAuthorizeKSignSecp256k1 core.ClaimAuthorizeKSignSecp256k1 `binding:"required"`
	Signature                    *utils.SignatureEthMsg            `binding:"required"`
	KSignP256                    *ecdsa.PublicKey                  `binding:"required"`
}

// SetRootMsg contains the data to set the SetRootClaim with its signature in Hex
type SetRootMsg struct {
	Root      string                 `binding:"required"`
	IdAddr    string                 `binding:"required"`
	KSignPk   *utils.PublicKey       `binding:"required"`
	Timestamp uint64                 `binding:"required"`
	Signature *utils.SignatureEthMsg `binding:"required"`
}

// ClaimValueMsg contains a core.ClaimValue with its signature in Hex
type ClaimValueMsg struct {
	ClaimValue merkletree.Entry       `binding:"required"`
	Signature  *utils.SignatureEthMsg `binding:"required"`
	KSignPk    *utils.PublicKey       `binding:"required"`
}

// TODO: Remove in next refactor
// ProofOfTreeLeaf contains all the parameters needed to proof that a Leaf is in a merkletree with a given Root
type ProofOfTreeLeaf struct {
	Leaf  []byte
	Proof []byte
	Root  merkletree.Hash
}

// TODO: Remove in next refactor
// ProofOfTreeLeafHex is the same data structure than ProofOfTreeLeaf but in Hexadecimal string representation
type ProofOfTreeLeafHex struct {
	Leaf  string
	Proof string
	Root  string
}

// TODO: Remove in next refactor
func (plh *ProofOfTreeLeafHex) Unhex() ProofOfTreeLeaf {
	var r ProofOfTreeLeaf
	r.Leaf, _ = common3.HexDecode(plh.Leaf)
	r.Proof, _ = common3.HexDecode(plh.Proof)
	rootBytes, _ := common3.HexDecode(plh.Root)
	copy(r.Root[:], rootBytes[:32])
	return r
}

// TODO: Remove in next refactor
// Hex returns a ProofOfTreeLeafHex data structure
func (pl *ProofOfTreeLeaf) Hex() ProofOfTreeLeafHex {
	r := ProofOfTreeLeafHex{
		common3.HexEncode(pl.Leaf),
		common3.HexEncode(pl.Proof),
		pl.Root.Hex(),
	}
	return r
}

// TODO: Remove in next refactor
// ProofOfClaimUser is the proof of a claim in the Identity MerkleTree, and the SetRootClaim of that MerkleTree inside the Relay's MerkleTree. Also with the proofs of non revocation of both claims
type ProofOfClaimUser struct {
	ClaimProof                     ProofOfTreeLeaf
	SetRootClaimProof              ProofOfTreeLeaf
	ClaimNonRevocationProof        ProofOfTreeLeaf
	SetRootClaimNonRevocationProof ProofOfTreeLeaf
	Date                           uint64
	Signature                      []byte // signature of the Root of the Relay
}

// TODO: Remove in next refactor
type ProofOfClaimUserHex struct {
	ClaimProof                     ProofOfTreeLeafHex
	SetRootClaimProof              ProofOfTreeLeafHex
	ClaimNonRevocationProof        ProofOfTreeLeafHex
	SetRootClaimNonRevocationProof ProofOfTreeLeafHex
	Date                           uint64
	Signature                      string // signature of the Root of the Relay
}

// TODO: Remove in next refactor
func (pc *ProofOfClaimUser) Hex() ProofOfClaimUserHex {
	r := ProofOfClaimUserHex{
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
func (pch *ProofOfClaimUserHex) Unhex() (ProofOfClaimUser, error) {
	sigBytes, err := common3.HexDecode(pch.Signature)
	if err != nil {
		return ProofOfClaimUser{}, err
	}
	r := ProofOfClaimUser{
		pch.ClaimProof.Unhex(),
		pch.SetRootClaimProof.Unhex(),
		pch.ClaimNonRevocationProof.Unhex(),
		pch.SetRootClaimNonRevocationProof.Unhex(),
		pch.Date,
		sigBytes,
	}
	return r, nil
}
