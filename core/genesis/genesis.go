package genesis

import (
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-merkletree-sql"
	"github.com/iden3/go-merkletree-sql/db/memory"
)

// CalculateIdGenesis calculates the ID given the input claims using memory Merkle Trees.
// Adds the given parameters into an ephemeral MerkleTree to calculate the MerkleRoot.
// ID: base58 ( [ type | root_genesis | checksum ] )
// where checksum: hash( [type | root_genesis ] )
// where the hash function is Poseidon
func CalculateIdGenesis(claimKOp *claims.ClaimKeyBabyJub, extraGenesisClaims []merkletree.Entrier) (*core.ID, error) {
	// add the claims into an ephemeral merkletree to calculate the genesis root to get that identity
	clt, err := merkletree.NewMerkleTree(memory.NewMemoryStorage(), 140)
	if err != nil {
		return nil, err
	}
	rot, err := merkletree.NewMerkleTree(memory.NewMemoryStorage(), 140)
	if err != nil {
		return nil, err
	}

	return CalculateIdGenesisMT(clt, rot, claimKOp, extraGenesisClaims)
}

// CalculateIdGenesisMT calculates the Genesis ID from the given claims using
// the given Claims Merkle Tree and Roots Merkle Tree.
func CalculateIdGenesisMT(clt *merkletree.MerkleTree, rot *merkletree.MerkleTree, claimKOp *claims.ClaimKeyBabyJub, extraGenesisClaims []merkletree.Entrier) (*core.ID, error) {
	err := clt.AddEntry(claimKOp.Entry())
	if err != nil {
		return nil, err
	}

	for _, claim := range extraGenesisClaims {
		if err := clt.AddEntry(claim.Entry()); err != nil {
			return nil, err
		}
	}

	clr := clt.Root()

	if err := claims.AddLeafRootsTree(rot, clr); err != nil {
		return nil, err
	}

	ror := rot.Root()

	idenState := core.IdenState(clr, &merkletree.HashZero, ror)
	id := core.IdGenesisFromIdenState(idenState)

	return id, nil
}
