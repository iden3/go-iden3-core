package identityagentsrv

import (
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/merkletree"
)

// identityagentsrv is a service that can hold multiple IdentityAgents (one for each Identity)

type Service interface {
	NewIdentity(claimAuthKOp merkletree.Claim, extraGenesisClaims []merkletree.Claim) (core.ID, core.ProofClaim, error)
	LoadIdStorages(prefix []byte) (db.Storage, *merkletree.MerkleTree, error)
}

type RootUpdaterConfig struct {
	Type   string            // "relay", "ethereum"
	Params map[string]string // if type=="relay" -> "url"
}

type ServiceImpl struct {
	storage     db.Storage
	rootUpdater RootUpdaterConfig
}

func New(storage db.Storage) *ServiceImpl {
	return &ServiceImpl{
		storage: storage,
	}
}

type IdStorages struct {
	storage db.Storage
	ecSto   db.Storage
	rcSto   db.Storage
	mt      *merkletree.MerkleTree
}

// LoadPrefixStorage returns the identity storages
func (ia *ServiceImpl) LoadIdStorages(id *core.ID) (*IdStorages, error) {
	idSto := ia.storage.WithPrefix(id.Bytes())
	ecSto := idSto.WithPrefix([]byte("emittedclaims"))
	rcSto := idSto.WithPrefix([]byte("receivedclaims"))
	mt, err := merkletree.NewMerkleTree(idSto, 140)
	return &IdStorages{
		storage: idSto,
		ecSto:   ecSto,
		rcSto:   rcSto,
		mt:      mt,
	}, err
}

// NewIdentity creates a new identity from the given claims
func (ia *ServiceImpl) NewIdentity(claimAuthKOp merkletree.Claim, extraGenesisClaims []merkletree.Claim) (*core.ID, *core.ProofClaim, error) {
	// calculate new ID in a memorydb
	id, proofKOp, err := core.CalculateIdGenesis(claimAuthKOp, extraGenesisClaims)
	if err != nil {
		return nil, nil, err
	}

	// load identity storages
	idStorages, err := ia.LoadIdStorages(id)
	if err != nil {
		return nil, nil, err
	}

	// add claims into the stored MerkleTree & into the EmittedClaimsStorage
	tx, err := idStorages.ecSto.NewTx() // for the moment a simple storage, in the future a storage that allows to query searches
	err = idStorages.mt.Add(claimAuthKOp.Entry())
	if err != nil {
		return nil, nil, err
	}
	tx.Put(claimAuthKOp.Entry().HIndex().Bytes(), claimAuthKOp.Entry().Bytes())
	for _, claim := range extraGenesisClaims {
		err = idStorages.mt.Add(claim.Entry())
		if err != nil {
			return nil, nil, err
		}
		tx.Put(claim.Entry().HIndex().Bytes(), claim.Entry().Bytes())
	}
	tx.Commit()

	// TODO send identity Root to RootUpdater (Relay)
	// this will be implemented when the Connection with RootUpdater is ready

	return id, proofKOp, nil
}
