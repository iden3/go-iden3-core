package identityagentsrv

import (
	"bytes"
	"errors"

	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/merkletree"
)

// identityagentsrv is a service that can hold multiple IdentityAgents (one for each Identity)

var PREFIX_EMITTEDCLAIMS = []byte("emittedclaims")
var PREFIX_RECEIVEDCLAIMS = []byte("receivedclaims")

type Service interface {
	LoadIdStorages(id *core.ID) (db.Storage, *merkletree.MerkleTree, error)
	NewIdentity(claimAuthKOp merkletree.Claim, extraGenesisClaims []merkletree.Claim) (core.ID, core.ProofClaim, error)
	AddClaim(id core.ID, claim merkletree.Claim) error
	AddClaims(id core.ID, claims []merkletree.Claim) error
	GetAllReceivedClaims(id *core.ID, idStorages *IdStorages) ([]ClaimObj, error)
	GetAllEmittedClaims(id *core.ID, idStorages *IdStorages) ([]ClaimObj, error)
	GetAllClaims(id *core.ID) ([]ClaimObj, []ClaimObj, error)
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
	ecSto := idSto.WithPrefix(PREFIX_EMITTEDCLAIMS)
	rcSto := idSto.WithPrefix(PREFIX_RECEIVEDCLAIMS)
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
	err = tx.Commit()
	if err != nil {
		return nil, nil, err
	}

	// TODO send identity Root to RootUpdater (Relay)
	// this will be implemented when the Connection with RootUpdater is ready

	return id, proofKOp, nil
}

func (ia *ServiceImpl) AddClaim(id *core.ID, claim merkletree.Claim) error {
	idStorages, err := ia.LoadIdStorages(id)
	if err != nil {
		return err
	}

	err = idStorages.mt.Add(claim.Entry())
	if err != nil {
		return err
	}
	tx, err := idStorages.ecSto.NewTx()
	tx.Put(claim.Entry().HIndex().Bytes(), claim.Entry().Bytes())
	err = tx.Commit()
	if err != nil {
		return err
	}

	// TODO send identity Root to RootUpdater (Relay)
	// this will be implemented when the Connection with RootUpdater is ready

	return nil
}

func (ia *ServiceImpl) AddClaims(id *core.ID, claims []merkletree.Claim) error {
	idStorages, err := ia.LoadIdStorages(id)
	if err != nil {
		return err
	}

	tx, err := idStorages.ecSto.NewTx()
	for _, claim := range claims {
		err = idStorages.mt.Add(claim.Entry())
		if err != nil {
			return err
		}
		tx.Put(claim.Entry().HIndex().Bytes(), claim.Entry().Bytes())
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	// TODO send identity Root to RootUpdater (Relay)
	// this will be implemented when the Connection with RootUpdater is ready

	cBytes := claims[0].Entry().Bytes()

	var leafBytes [merkletree.ElemBytesLen * merkletree.DataLen]byte
	copy(leafBytes[:], cBytes[:merkletree.ElemBytesLen*merkletree.DataLen])
	leafData := merkletree.NewDataFromBytes(leafBytes)
	// leafDataBytes := leafData.Bytes()

	// assert.Equal(t, cBytes, leafDataBytes[:])
	// assert.Equal(t, cBytes, leafBytes[:])

	entry := merkletree.Entry{
		Data: *leafData,
	}
	for _, elemBytes := range entry.Data {
		if _, err := merkletree.ElemBytesToRElem(elemBytes); err != nil {
			return err
		}
	}

	return nil
}

type ClaimObj struct {
	Claim merkletree.Claim
	Proof core.ProofClaimPartial // TODO once the RootUpdater is ready we can use here the proof part of the Relay (or direct from blockchain)
}

func (ia *ServiceImpl) GetAllReceivedClaims(id *core.ID, idStorages *IdStorages) ([]ClaimObj, error) {
	var err error
	if idStorages == nil {
		idStorages, err = ia.LoadIdStorages(id)
		if err != nil {
			return []ClaimObj{}, err
		}
	}

	// get received claims
	var receivedClaims []ClaimObj
	idStorages.rcSto.Iterate(func(key, value []byte) {
		// filter only the key-value with the prefix of id+emittedclaims
		// as the Iterate from the Storage don't filters by prefix
		// in the future do a more efficient way to filter without going through all the keys-values
		var prefix []byte
		prefix = append(prefix[:], id.Bytes()...)
		prefix = append(prefix[:], PREFIX_RECEIVEDCLAIMS...)
		if bytes.Equal(key[:len(prefix)], prefix) {
			rClaim := ClaimObj{
				// TODO to be defined the way to store received claims
				Claim: &core.ClaimBasic{},
				Proof: core.ProofClaimPartial{},
			}
			receivedClaims = append(receivedClaims, rClaim)
		}
	})

	return receivedClaims, err
}

func (ia *ServiceImpl) GetAllEmittedClaims(id *core.ID, idStorages *IdStorages) ([]ClaimObj, error) {
	var err error
	if idStorages == nil {
		idStorages, err = ia.LoadIdStorages(id)
		if err != nil {
			return []ClaimObj{}, err
		}
	}

	// get emitted claims, and generate fresh proof with current Root
	var emittedClaims []ClaimObj
	var iterErr error
	idStorages.ecSto.Iterate(func(key, value []byte) {
		// filter only the key-value with the prefix of id+emittedclaims
		// as the Iterate from the Storage don't filters by prefix
		// in the future do a more efficient way to filter without going through all the keys-values
		var prefix []byte
		prefix = append(prefix[:], id.Bytes()...)
		prefix = append(prefix[:], PREFIX_EMITTEDCLAIMS...)
		if bytes.Equal(key[:len(prefix)], prefix) {
			// where key is the hi, value is the leaf
			var hi_b [32]byte
			copy(hi_b[:], key[:32])
			hi := merkletree.Hash(hi_b)

			mtp, err := idStorages.mt.GenerateProof(&hi, nil)

			// TODO see issues #167 and #169
			// once RootUpdater is ready, tie the claim proof in the identity merkletree. Also the mtpNonRevokated depends on the RootUpdater
			// with the proof of the SetRootClaim in the Relay's merkletree
			proof := core.ProofClaimPartial{
				Mtp0: mtp,
				Mtp1: &merkletree.Proof{},
				Root: idStorages.mt.RootKey(),
				Aux:  nil,
			}

			var leafBytes [merkletree.ElemBytesLen * merkletree.DataLen]byte
			copy(leafBytes[:], value[:merkletree.ElemBytesLen*merkletree.DataLen])
			leafData := merkletree.NewDataFromBytes(leafBytes)
			entry := merkletree.Entry{
				Data: *leafData,
			}
			c, err := core.NewClaimFromEntry(&entry)
			if err != nil {
				iterErr = errors.New(err.Error())
				return
			}
			eClaim := ClaimObj{
				Claim: c,
				Proof: proof,
			}
			emittedClaims = append(emittedClaims, eClaim)
		}
	})
	if iterErr != nil {
		return emittedClaims, iterErr
	}
	return emittedClaims, err

}

func (ia *ServiceImpl) GetAllClaims(id *core.ID) ([]ClaimObj, []ClaimObj, error) {
	idStorages, err := ia.LoadIdStorages(id)
	if err != nil {
		return []ClaimObj{}, []ClaimObj{}, err
	}

	// get received claims
	receivedClaims, err := ia.GetAllReceivedClaims(id, idStorages)
	if err != nil {
		return []ClaimObj{}, []ClaimObj{}, err
	}

	// get emitted claims, and generate fresh proof with current Root
	emittedClaims, err := ia.GetAllEmittedClaims(id, idStorages)
	return emittedClaims, receivedClaims, err

}
