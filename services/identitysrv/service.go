package identitysrv

import (
	// "crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/crypto/babyjub"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/claimsrv"
)

type Service interface {
	CreateIdGenesis(kop *babyjub.PublicKey, kdis, kreen, kupdateRoot common.Address) (*core.ID, *core.ProofClaim, error)
}

type ServiceImpl struct {
	cs claimsrv.Service
}

func New(cs claimsrv.Service) *ServiceImpl {
	return &ServiceImpl{
		cs: cs,
	}
}

// CreateIdGenesis initializes the id MerkleTree with the given the kop, kdisable,
// kreenable and kupdateRoots public keys. Where the id is calculated a MerkleTree containing
// that initial data, calculated in the function CalculateIdGenesis()
func (is *ServiceImpl) CreateIdGenesis(kop *babyjub.PublicKey, kdis, kreen, kupdateRoot common.Address) (*core.ID, *core.ProofClaim, error) {

	id, claims, err := core.CalculateIdGenesis(kop, kdis, kreen, kupdateRoot)
	if err != nil {
		return nil, nil, err
	}

	// add the claims into the storage merkletree of that identity
	stoUserId := is.cs.MT().Storage().WithPrefix(id.Bytes())
	userMT, err := merkletree.NewMerkleTree(stoUserId, 140)
	if err != nil {
		return nil, nil, err
	}

	for _, claim := range claims {
		err = userMT.Add(claim.Entry())
		if err != nil {
			return nil, nil, err
		}
	}

	// create new ClaimSetRootKey
	claimSetRootKey := core.NewClaimSetRootKey(*id, *userMT.RootKey())

	// add User's Id Merkle Root into the Relay's Merkle Tree
	err = is.cs.MT().Add(claimSetRootKey.Entry())
	if err != nil {
		return nil, nil, err
	}

	// update Relay's Root in the Smart Contract
	is.cs.RootSrv().SetRoot(*is.cs.MT().RootKey())

	proofClaimKop, err := is.cs.GetClaimProofUserByHi(*id, claims[0].Entry().HIndex())
	if err != nil {
		return nil, nil, err
	}

	return id, proofClaimKop, nil
}
