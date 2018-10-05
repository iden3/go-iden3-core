package claimsrv

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/utils"
)

// CheckProofOfClaim checks the Merkle Proof of the Claim, the SetRootClaim, and the non revocation proof of both claims
func CheckProofOfClaim(relayAddr common.Address, pc ProofOfClaim, numLevels int) bool {
	vClaimProof := merkletree.CheckProof(pc.ClaimProof.Root, pc.ClaimProof.Proof,
		pc.ClaimProof.Hi, merkletree.HashBytes(pc.ClaimProof.Leaf), numLevels)

	vSetRootClaimProof := merkletree.CheckProof(pc.SetRootClaimProof.Root, pc.SetRootClaimProof.Proof,
		pc.SetRootClaimProof.Hi, merkletree.HashBytes(pc.SetRootClaimProof.Leaf), numLevels)

	vClaimNonRevocationProof := merkletree.CheckProof(pc.ClaimNonRevocationProof.Root, pc.ClaimNonRevocationProof.Proof,
		pc.ClaimNonRevocationProof.Hi, merkletree.EmptyNodeValue, numLevels)

	vSetRootClaimNonRevocationProof := merkletree.CheckProof(pc.SetRootClaimNonRevocationProof.Root, pc.SetRootClaimNonRevocationProof.Proof,
		pc.SetRootClaimNonRevocationProof.Hi, merkletree.EmptyNodeValue, numLevels)

	// check signature of the ProofOfClaim.SetRootClaimProof.Root with the identity of the Relay
	// checking this Root and the four Merkle Proofs, we check the full ProofOfClaim
	if !utils.VerifySig(relayAddr, pc.Signature, pc.SetRootClaimProof.Root[:]) {
		return false
	}

	if vClaimProof && vSetRootClaimProof && vClaimNonRevocationProof && vSetRootClaimNonRevocationProof {
		return true
	}
	return false
}
