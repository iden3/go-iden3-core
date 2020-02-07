package genesis

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-crypto/babyjub"
)

// CalculateIdGenesis calculates the ID given the input claims using memory Merkle Trees.
// Adds the given parameters into an ephemeral MerkleTree to calculate the MerkleRoot.
// ID: base58 ( [ type | root_genesis | checksum ] )
// where checksum: hash( [type | root_genesis ] )
// where the hash function is Poseidon
func CalculateIdGenesis(claimKOp *claims.ClaimAuthorizeKSignBabyJub, extraGenesisClaims []merkletree.Entrier) (*core.ID, *proof.ProofClaim, error) {
	// add the claims into an ephemeral merkletree to calculate the genesis root to get that identity
	clt, err := merkletree.NewMerkleTree(db.NewMemoryStorage(), 140)
	if err != nil {
		return nil, nil, err
	}
	rot, err := merkletree.NewMerkleTree(db.NewMemoryStorage(), 140)
	if err != nil {
		return nil, nil, err
	}

	return CalculateIdGenesisMT(clt, rot, claimKOp, extraGenesisClaims)
}

// CalculateIdGenesisMT calculates the Genesis ID from the given claims using
// the given Claims Merkle Tree and Roots Merkle Tree.
func CalculateIdGenesisMT(clt *merkletree.MerkleTree, rot *merkletree.MerkleTree, claimKOp *claims.ClaimAuthorizeKSignBabyJub, extraGenesisClaims []merkletree.Entrier) (*core.ID, *proof.ProofClaim, error) {
	err := clt.AddClaim(claimKOp)
	if err != nil {
		return nil, nil, err
	}

	for _, claim := range extraGenesisClaims {
		if err := clt.AddClaim(claim); err != nil {
			return nil, nil, err
		}
	}

	clr := clt.RootKey()

	if err := claims.AddLeafRootsTree(rot, clr); err != nil {
		return nil, nil, err
	}

	ror := rot.RootKey()

	proofClaimKOp, err := proof.GetClaimProofByHi(clt, claimKOp.Entry().HIndex())
	if err != nil {
		return nil, nil, err
	}

	idenState := core.IdenState(clr, &merkletree.HashZero, ror)
	id := core.IdGenesisFromIdenState(idenState)

	return id, proofClaimKOp, nil
}

// DEPRECATED
type GenesisProofClaims struct {
	KOp         proof.ProofClaim
	KDis        proof.ProofClaim
	KReen       proof.ProofClaim
	KUpdateRoot proof.ProofClaim
}

// DEPRECATED
func CalculateIdGenesisFrom4Keys(kop *babyjub.PublicKey, kdis, kreen, kupdateRoot common.Address) (*core.ID, *GenesisProofClaims, error) {
	// add the claims into an ephemeral merkletree to calculate the genesis root to get that identity
	mt, err := merkletree.NewMerkleTree(db.NewMemoryStorage(), 140)
	if err != nil {
		return nil, nil, err
	}

	claimKOp := claims.NewClaimAuthorizeKSignBabyJub(kop, 0)
	err = mt.AddClaim(claimKOp)
	if err != nil {
		return nil, nil, err
	}

	claimKDis := claims.NewClaimAuthEthKey(kdis, claims.EthKeyTypeDisable)
	err = mt.AddClaim(claimKDis)
	if err != nil {
		return nil, nil, err
	}
	claimKReen := claims.NewClaimAuthEthKey(kreen, claims.EthKeyTypeReenable)
	err = mt.AddClaim(claimKReen)
	if err != nil {
		return nil, nil, err
	}
	claimKUpdateRoot := claims.NewClaimAuthEthKey(kupdateRoot, claims.EthKeyTypeUpdateRoot)
	err = mt.AddClaim(claimKUpdateRoot)
	if err != nil {
		return nil, nil, err
	}

	idGenesis := mt.RootKey()

	proofClaimKOp, err := proof.GetClaimProofByHi(mt, claimKOp.Entry().HIndex())
	if err != nil {
		return nil, nil, err
	}
	proofClaimKDis, err := proof.GetClaimProofByHi(mt, claimKDis.Entry().HIndex())
	if err != nil {
		return nil, nil, err
	}
	proofClaimKReen, err := proof.GetClaimProofByHi(mt, claimKReen.Entry().HIndex())
	if err != nil {
		return nil, nil, err
	}
	proofClaimKUpdateRoot, err := proof.GetClaimProofByHi(mt, claimKUpdateRoot.Entry().HIndex())
	if err != nil {
		return nil, nil, err
	}

	var idGenesisBytes [27]byte
	copy(idGenesisBytes[:], idGenesis.Bytes()[len(idGenesis.Bytes())-27:])
	id := core.NewID(core.TypeBJP0, idGenesisBytes)
	return &id, &GenesisProofClaims{
		KOp:         *proofClaimKOp,
		KDis:        *proofClaimKDis,
		KReen:       *proofClaimKReen,
		KUpdateRoot: *proofClaimKUpdateRoot,
	}, nil
}
