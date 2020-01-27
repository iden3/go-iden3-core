package messages

import (
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/utils"
	"github.com/iden3/go-iden3-crypto/babyjub"
)

// SetRoot0Req contains the data to set the SetRootClaim
type SetRoot0Req struct {
	OldRoot           *merkletree.Hash         `json:"oldRoot" binding:"required"`
	NewRoot           *merkletree.Hash         `json:"newRoot" binding:"required"`
	ClaimAuthKOp      *merkletree.Entry        `json:"claimKOp" binding:"required"`
	ProofClaimAuthKOp *proof.ProofClaimGenesis `json:"proofKOp" binding:"required,dive"`
	Date              int64                    `json:"date" binding:"required"`
	Signature         *babyjub.SignatureComp   `json:"signature" binding:"required"` // signature of the Root
}

// ClaimValueReq contains a core.ClaimValue with its signature in Hex
type ClaimValueReq struct {
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
// ProofClaimUserRes is the proof of a claim in the Identity MerkleTree, and the SetRootClaim of that MerkleTree inside the Relay's MerkleTree. Also with the proofs of non revocation of both claims
type ProofClaimUserRes struct {
	ClaimProof                     ProofTreeLeaf
	SetRootClaimProof              ProofTreeLeaf
	ClaimNonRevocationProof        ProofTreeLeaf
	SetRootClaimNonRevocationProof ProofTreeLeaf
	Date                           int64
	Signature                      []byte // signature of the Root of the Relay
}
