package identitysrv

import (
	"crypto/ecdsa"

	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/claimsrv"
)

type Service interface {
	CreateIdGenesis(kop, krec, krev *ecdsa.PublicKey) (*core.ID, *core.ProofClaim, error)
}

type ServiceImpl struct {
	cs claimsrv.Service
}

func New(cs claimsrv.Service) *ServiceImpl {
	return &ServiceImpl{
		cs: cs,
	}
}

// CreateIdGenesis initializes the idAddress MerkleTree with the given the kop, krec,
// krev public keys. Where the idAddress is calculated a MerkleTree containing
// that initial data, calculated in the function CalculateIdGenesis()
func (is *ServiceImpl) CreateIdGenesis(kop, krec, krev *ecdsa.PublicKey) (*core.ID, *core.ProofClaim, error) {

	idAddr, err := core.CalculateIdGenesis(kop, krec, krev)
	if err != nil {
		return nil, nil, err
	}

	// add the claims into the storage merkletree of that identity
	stoUserId := is.cs.MT().Storage().WithPrefix(idAddr.Bytes())
	userMT, err := merkletree.NewMerkleTree(stoUserId, 140)
	if err != nil {
		return nil, nil, err
	}

	// generate the Authorize KSign Claims for the given public Keys
	claims := core.GenerateArrayClaimAuthorizeKSignFromPublicKeys(kop, krec, krev)

	for _, claim := range claims {
		err = userMT.Add(claim.Entry())
		if err != nil {
			return nil, nil, err
		}
	}

	// create new ClaimSetRootKey
	claimSetRootKey := core.NewClaimSetRootKey(*idAddr, *userMT.RootKey())

	// add User's Id Merkle Root into the Relay's Merkle Tree
	err = is.cs.MT().Add(claimSetRootKey.Entry())
	if err != nil {
		return nil, nil, err
	}

	// update Relay's Root in the Smart Contract
	is.cs.RootSrv().SetRoot(*is.cs.MT().RootKey())

	proofClaimKop, err := is.cs.GetClaimProofUserByHi(*idAddr, claims[0].Entry().HIndex())
	if err != nil {
		return nil, nil, err
	}

	return idAddr, proofClaimKop, nil
}
