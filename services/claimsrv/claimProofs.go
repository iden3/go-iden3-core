package claimsrv

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/utils"
)

// CheckProofOfClaim checks the Merkle Proof of the Claim, the SetRootClaim, and the non revocation proof of both claims
func CheckProofOfClaim(relayAddr common.Address, pc ProofOfClaim, numLevels int) bool {
	hiClaim := core.HiFromClaimBytes(pc.ClaimProof.Leaf)
	vClaimProof := merkletree.CheckProof(pc.ClaimProof.Root, pc.ClaimProof.Proof,
		hiClaim, merkletree.HashBytes(pc.ClaimProof.Leaf), numLevels)

	hiSetRootClaim := core.HiFromClaimBytes(pc.SetRootClaimProof.Leaf)
	vSetRootClaimProof := merkletree.CheckProof(pc.SetRootClaimProof.Root, pc.SetRootClaimProof.Proof,
		hiSetRootClaim, merkletree.HashBytes(pc.SetRootClaimProof.Leaf), numLevels)

	hiNonRevocationClaim := core.HiFromClaimBytes(pc.ClaimNonRevocationProof.Leaf)
	vClaimNonRevocationProof := merkletree.CheckProof(pc.ClaimNonRevocationProof.Root, pc.ClaimNonRevocationProof.Proof,
		hiNonRevocationClaim, merkletree.EmptyNodeValue, numLevels)

	hiNonRevocationSetRootClaim := core.HiFromClaimBytes(pc.SetRootClaimNonRevocationProof.Leaf)
	vSetRootClaimNonRevocationProof := merkletree.CheckProof(pc.SetRootClaimNonRevocationProof.Root, pc.SetRootClaimNonRevocationProof.Proof,
		hiNonRevocationSetRootClaim, merkletree.EmptyNodeValue, numLevels)

	// additional, check caducity of the pc.Date

	// check signature of the ProofOfClaim.SetRootClaimProof.Root with the identity of the Relay
	// checking this Root and the four Merkle Proofs, we check the full ProofOfClaim
	dateBytes, err := core.Uint64ToEthBytes(pc.Date)
	if err != nil {
		return false
	}
	rootdate := pc.SetRootClaimProof.Root[:]
	rootdate = append(rootdate, dateBytes...)
	rootdateHash := merkletree.HashBytes(rootdate)
	if !utils.VerifySig(relayAddr, pc.Signature, rootdateHash[:]) {
		return false
	}

	if vClaimProof && vSetRootClaimProof && vClaimNonRevocationProof && vSetRootClaimNonRevocationProof {
		return true
	}
	return false
}
