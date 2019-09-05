package identitysrv

import (
	// "crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/services/claimsrv"
	"github.com/iden3/go-iden3-crypto/babyjub"
)

var PREFIX_MERKLETREE = []byte("merkletree")

type Service struct {
	cs *claimsrv.Service
}

func New(cs *claimsrv.Service) *Service {
	return &Service{
		cs: cs,
	}
}

// CreateIdGenesis initializes the id MerkleTree with the given the kop, kdisable,
// kreenable and kupdateRoots public keys. Where the id is calculated a MerkleTree containing
// that initial data, calculated in the function CalculateIdGenesis()
func (is *Service) CreateIdGenesis(kop *babyjub.PublicKey, kdis, kreen, kupdateRoot common.Address) (*core.ID, *core.ProofClaim, error) {

	id, proofClaims, err := core.CalculateIdGenesisFrom4Keys(kop, kdis, kreen, kupdateRoot)
	if err != nil {
		return nil, nil, err
	}

	// add the claims into the storage merkletree of that identity
	stoUserId := is.cs.MT().Storage().WithPrefix(id.Bytes()).WithPrefix(PREFIX_MERKLETREE)
	userMT, err := merkletree.NewMerkleTree(stoUserId, 140)
	if err != nil {
		return nil, nil, err
	}

	proofClaimsList := []core.ProofClaim{proofClaims.KOp, proofClaims.KDis,
		proofClaims.KReen, proofClaims.KUpdateRoot}
	for _, proofClaim := range proofClaimsList {
		err = userMT.Add(proofClaim.Claim)
		if err != nil {
			return nil, nil, err
		}
	}

	// create new ClaimSetRootKey
	claimSetRootKey, err := core.NewClaimSetRootKey(id, userMT.RootKey())
	if err != nil {
		return nil, nil, err
	}

	// add User's Id Merkle Root into the Relay's Merkle Tree
	err = is.cs.MT().Add(claimSetRootKey.Entry())
	if err != nil {
		return nil, nil, err
	}

	// update Relay's Root in the Smart Contract
	is.cs.RootSrv().SetRoot(*is.cs.MT().RootKey())

	return id, &proofClaims.KOp, nil
}
