package claimsrv

import "github.com/iden3/go-iden3/merkletree"

// CheckProofOfClaim checks the Merkle Proof of the Claim, the SetRootClaim, and the non revocation proof of both claims
func CheckProofOfClaim(pc ProofOfClaim, numLevels int) bool {
	vClaimProof := merkletree.CheckProof(pc.ClaimProof.Root, pc.ClaimProof.Proof,
		pc.ClaimProof.Hi, merkletree.HashBytes(pc.ClaimProof.Leaf), numLevels)

	vSetRootClaimProof := merkletree.CheckProof(pc.SetRootClaimProof.Root, pc.SetRootClaimProof.Proof,
		pc.SetRootClaimProof.Hi, merkletree.HashBytes(pc.SetRootClaimProof.Leaf), numLevels)

	vClaimNonRevocationProof := merkletree.CheckProof(pc.ClaimNonRevocationProof.Root, pc.ClaimNonRevocationProof.Proof,
		pc.ClaimNonRevocationProof.Hi, merkletree.EmptyNodeValue, numLevels)

	vSetRootClaimNonRevocationProof := merkletree.CheckProof(pc.SetRootClaimNonRevocationProof.Root, pc.SetRootClaimNonRevocationProof.Proof,
		pc.SetRootClaimNonRevocationProof.Hi, merkletree.EmptyNodeValue, numLevels)

	if vClaimProof && vSetRootClaimProof && vClaimNonRevocationProof && vSetRootClaimNonRevocationProof {
		return true
	}
	return false
}
