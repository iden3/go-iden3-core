package identityagentsrv

import (
	"bytes"
	"encoding/hex"
	// "errors"
	"fmt"

	"github.com/dghubble/sling"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/services/claimsrv"
)

// identityagentsrv is a service that can hold multiple IdentityAgents (one for each Identity)

var PREFIX_EMITTEDCLAIMS = []byte("emittedclaims")
var PREFIX_RECEIVEDCLAIMS = []byte("receivedclaims")

// TODO: Move this to a generic place
type ServerError struct {
	Err string `json:"error"`
}

// TODO: Move this to a generic place
func (e ServerError) Error() string {
	return fmt.Sprintf("server: %v", e.Err)
}

type RootUpdaterMock struct{}

func (ru *RootUpdaterMock) RootUpdate(setRootMsg claimsrv.SetRootMsg) (*core.ProofClaim, error) {
	return nil, fmt.Errorf("mock mock")
}

type RootUpdaterRelay struct {
	UrlRelay string
}

func NewRootUpdaterRelay(url string) RootUpdaterRelay {
	return RootUpdaterRelay{UrlRelay: url}
}

func (ru *RootUpdaterRelay) RootUpdate(setRootMsg claimsrv.SetRootMsg) (*core.ProofClaim, error) {
	url := fmt.Sprintf("%s/ids/%s/root", ru.UrlRelay, setRootMsg.Id)
	proofClaim, serverError := struct {
		ProofClaim core.ProofClaim `json:"proofClaim"`
	}{}, ServerError{}
	_, err := sling.New().Base(url).BodyJSON(setRootMsg).Receive(&proofClaim, &serverError)
	if err == nil {
		err = serverError
	}
	return &proofClaim.ProofClaim, err
}

type RootUpdater interface {
	RootUpdate(setRootMsg claimsrv.SetRootMsg) (*core.ProofClaim, error)
}

type Service struct {
	storage     db.Storage
	rootUpdater RootUpdater
}

func New(storage db.Storage, rootUpdater RootUpdater) *Service {
	return &Service{
		storage:     storage,
		rootUpdater: rootUpdater,
	}
}

// NewIdentity creates a new identity from the given claims
func (ia *Service) CreateIdentity(claimAuthKOp *merkletree.Entry,
	extraGenesisClaims []*merkletree.Entry) (*core.ID, *core.ProofClaim, error) {
	// calculate new ID in a memorydb
	id, proofKOp, err := core.CalculateIdGenesis(claimAuthKOp, extraGenesisClaims)
	if err != nil {
		return nil, nil, err
	}

	// load identity storages
	agent, err := ia.NewAgent(id)
	if err != nil {
		return nil, nil, err
	}

	// add claims into the stored MerkleTree & into the EmittedClaimsStorage
	// for the moment a simple storage, in the future a storage that allows to query searches
	tx, err := agent.storage.claims.emitted.NewTx()
	err = agent.mt.Add(claimAuthKOp)
	if err != nil {
		return nil, nil, err
	}
	tx.Put(claimAuthKOp.HIndex().Bytes(), claimAuthKOp.Bytes())
	for _, claim := range extraGenesisClaims {
		err = agent.mt.Add(claim)
		if err != nil {
			return nil, nil, err
		}
		tx.Put(claim.HIndex().Bytes(), claim.Bytes())
	}
	err = tx.Commit()
	if err != nil {
		return nil, nil, err
	}

	// TODO send identity Root to RootUpdater (Relay)
	// this will be implemented when the Connection with RootUpdater is ready

	return id, proofKOp, nil
}

type IdStorage struct {
	base   db.Storage
	claims struct {
		emitted  db.Storage
		received db.Storage
	}
}

type Agent struct {
	storage *IdStorage
	mt      *merkletree.MerkleTree
	id      *core.ID
}

// LoadPrefixStorage returns the identity storages
func (a *Agent) loadStorage(base db.Storage) error {
	emittedClaims := base.WithPrefix(PREFIX_EMITTEDCLAIMS)
	receivedClaims := base.WithPrefix(PREFIX_RECEIVEDCLAIMS)
	a.storage = &IdStorage{
		base: base,
		claims: struct {
			emitted  db.Storage
			received db.Storage
		}{
			emitted:  emittedClaims,
			received: receivedClaims,
		},
	}
	mt, err := merkletree.NewMerkleTree(base, 140)
	a.mt = mt
	return err
}

func (s *Service) NewAgent(id *core.ID) (*Agent, error) {
	agent := &Agent{id: id}
	err := agent.loadStorage(s.storage.WithPrefix(agent.id.Bytes()))
	return agent, err
}

func (a *Agent) AddClaim(claim *merkletree.Entry) error {
	return a.AddClaims([]*merkletree.Entry{claim})
}

func (a *Agent) AddClaims(claims []*merkletree.Entry) error {
	tx, err := a.storage.claims.emitted.NewTx()
	for _, claim := range claims {
		err = a.mt.Add(claim)
		if err != nil {
			return err
		}
		tx.Put(claim.HIndex().Bytes(), claim.Bytes())
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	// TODO send identity Root to RootUpdater (Relay)
	// this will be implemented when the Connection with RootUpdater is ready

	// TODO: @arnaucube: what is this?
	// cBytes := claims[0].Entry().Bytes()

	// var leafBytes [merkletree.ElemBytesLen * merkletree.DataLen]byte
	// copy(leafBytes[:], cBytes[:merkletree.ElemBytesLen*merkletree.DataLen])
	// leafData := merkletree.NewDataFromBytes(leafBytes)
	// // leafDataBytes := leafData.Bytes()

	// // assert.Equal(t, cBytes, leafDataBytes[:])
	// // assert.Equal(t, cBytes, leafBytes[:])

	// entry := merkletree.Entry{
	// 	Data: *leafData,
	// }
	// for _, elemBytes := range entry.Data {
	// 	if _, err := merkletree.ElemBytesToRElem(elemBytes); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func (a *Agent) GetAllReceivedClaims() ([]*merkletree.Entry, error) {
	var receivedClaims []*merkletree.Entry
	err := a.storage.claims.received.Iterate(func(key, value []byte) {
		// filter only the key-value with the prefix of id+emittedclaims
		// as the Iterate from the Storage don't filters by prefix
		// in the future do a more efficient way to filter without going through all the keys-values
		var prefix []byte
		prefix = append(prefix[:], a.id.Bytes()...)
		prefix = append(prefix[:], PREFIX_RECEIVEDCLAIMS...)
		if bytes.Equal(key[:len(prefix)], prefix) {
			// TODO to be defined the way to store received claims
			claim := merkletree.Entry{}
			receivedClaims = append(receivedClaims, &claim)
		}
	})
	return receivedClaims, err
}

func (a *Agent) GetAllEmittedClaims() ([]*merkletree.Entry, error) {
	// get emitted claims, and generate fresh proof with current Root
	var emittedClaims []*merkletree.Entry
	err := a.storage.claims.emitted.Iterate(func(key, value []byte) {
		// filter only the key-value with the prefix of id+emittedclaims
		// as the Iterate from the Storage don't filters by prefix
		// in the future do a more efficient way to filter without going through all the keys-values
		var prefix []byte
		prefix = append(prefix[:], a.id.Bytes()...)
		prefix = append(prefix[:], PREFIX_EMITTEDCLAIMS...)
		if bytes.Equal(key[:len(prefix)], prefix) {
			// where key is the hi, value is the leaf
			// var hiBytes [32]byte
			// copy(hiBytes[:], key[:32])
			// hi := merkletree.Hash(hiBytes)

			// mtp, err := idStorages.mt.GenerateProof(&hi, nil)

			// TODO see issues #167 and #169
			// once RootUpdater is ready, tie the claim proof in
			// the identity merkletree. Also the mtpNonRevokated
			// depends on the RootUpdater
			// with the proof of the SetRootClaim in the Relay's merkletree
			// proof := core.ProofClaimPartial{
			// 	Mtp0: mtp,
			// 	Mtp1: &merkletree.Proof{},
			// 	Root: idStorages.mt.RootKey(),
			// 	Aux:  nil,
			// }

			var data [merkletree.ElemBytesLen * merkletree.DataLen]byte
			copy(data[:], value)
			claim := merkletree.NewEntryFromBytes(data)
			// leafData := merkletree.NewDataFromBytes(leafBytes)
			// entry := merkletree.Entry{
			// 	Data: *leafData,
			// }
			// c, err := core.NewClaimFromEntry(&entry)
			// if err != nil {
			// 	iterErr = errors.New(err.Error())
			// 	return
			// }
			// eClaim := core.ClaimObj{
			// 	Claim: c,
			// 	Proof: proof,
			// }
			emittedClaims = append(emittedClaims, claim)
		}
	})
	return emittedClaims, err

}

func (a *Agent) GetClaimByHi(hi *merkletree.Hash) (*merkletree.Entry, *core.ProofClaimPartial, error) {
	leafData, err := a.mt.GetDataByIndex(hi)
	if err != nil {
		return nil, nil, err
	}
	entry := &merkletree.Entry{
		Data: *leafData,
	}

	mtp, err := a.mt.GenerateProof(hi, nil)
	if err != nil {
		return nil, nil, err
	}
	// TODO see issues #167 and #169
	// once RootUpdater is ready, tie the claim proof in the identity
	// merkletree. Also the mtpNonRevokated depends on the RootUpdater
	// with the proof of the SetRootClaim in the Relay's merkletree
	proof := &core.ProofClaimPartial{
		Mtp0: mtp,
		Mtp1: &merkletree.Proof{},
		// FIXME: There is a data race here!
		Root: a.mt.RootKey(),
		Aux:  nil,
	}
	return entry, proof, nil
}

func (a *Agent) GetFullMT() (map[string]string, error) {
	mt := make(map[string]string)

	// FIXME: This is not a full Merkle Tree, but a list of emitted claims :S
	// FIXME: Ok, this is abusing the no prefix filter in the db iteration... Please, use an iterator over the mt!
	a.storage.claims.emitted.Iterate(func(key, value []byte) {
		// filter only the key-value with the prefix of id+emittedclaims
		// as the Iterate from the Storage don't filters by prefix
		// in the future do a more efficient way to filter without going through all the keys-values
		var prefix []byte
		prefix = append(prefix[:], a.id.Bytes()...)
		prefix = append(prefix[:], merkletree.PREFIX_MERKLETREE...)
		if bytes.Equal(key[:len(prefix)], prefix) {
			mt["0x"+hex.EncodeToString(key[len(prefix):])] = "0x" + hex.EncodeToString(value)
		}
	})
	return mt, nil
}

// GetCurrentRoot is used from wallet to check if it is syncronized with the
// MerkleTree in the IdentityAgent
func (a *Agent) GetCurrentRoot() *merkletree.Hash {
	return a.mt.RootKey()
}
