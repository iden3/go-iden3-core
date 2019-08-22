package identityagentsrv

import (
	// "bytes"
	// "encoding/hex"
	// "errors"
	"fmt"

	"github.com/dghubble/sling"
	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/services/claimsrv"
	"gopkg.in/go-playground/validator.v9"
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

func (ru *RootUpdaterMock) RootUpdate(setRootReq claimsrv.SetRoot0Req) error {
	return fmt.Errorf("mock mock")
}

func (ru *RootUpdaterMock) GetRootProof() (*core.ProofClaim, error) {
	return nil, fmt.Errorf("mock mock")
}

type RootUpdaterRelay struct {
	RelayUrl string
	RelayId  *core.ID
	UserId   *core.ID
	_client  *sling.Sling
	validate *validator.Validate
}

func NewRootUpdaterRelay(relayUrl string, relayId, userId *core.ID) RootUpdaterRelay {
	if relayUrl[len(relayUrl)-1] != '/' {
		relayUrl += "/"
	}
	client := sling.New().Base(relayUrl)
	return RootUpdaterRelay{RelayUrl: relayUrl, RelayId: relayId, UserId: userId,
		_client: client, validate: validator.New()}
}

func (ru *RootUpdaterRelay) client() *sling.Sling {
	return ru._client.New()
}

func (ru *RootUpdaterRelay) request(s *sling.Sling, res interface{}) error {
	var serverError ServerError
	resp, err := s.Receive(res, &serverError)
	if err == nil {
		defer resp.Body.Close()
		if !(200 <= resp.StatusCode && resp.StatusCode < 300) {
			err = serverError
		} else if res != nil {
			err = ru.validate.Struct(res)
		}
	}
	return err
}

func (ru *RootUpdaterRelay) RootUpdate(setRootReq claimsrv.SetRoot0Req) error {
	var setRootClaim struct {
		SetRootClaim *merkletree.Entry `json:"setRootClaim" validate:"required"`
	}
	path := fmt.Sprintf("ids/%s/setrootclaim", ru.UserId)
	return ru.request(ru.client().Path(path).Post("").BodyJSON(setRootReq), &setRootClaim)
}

func (ru *RootUpdaterRelay) GetRootProof() (*core.ProofClaim, error) {
	var proofClaim struct {
		ProofClaim *core.ProofClaim `json:"proofClaim" validate:"required"`
	}
	path := fmt.Sprintf("ids/%s/setrootclaim", ru.UserId)
	err := ru.request(ru.client().Path(path).Get(""), &proofClaim)
	return proofClaim.ProofClaim, err
}

type RootUpdater interface {
	RootUpdate(setRootMsg claimsrv.SetRoot0Req) error
	GetRootProof() (*core.ProofClaim, error)
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
		// FIXME: If some agent.mt.Add fails, there will be an
		// inconsistency between agent.mt and agent.storage.claims.emitted
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
	err := a.storage.claims.received.Iterate(func(key, value []byte) (bool, error) {
		claim, err := merkletree.NewEntryFromBytes(value)
		if err != nil {
			return false, err
		}
		receivedClaims = append(receivedClaims, claim)
		return true, err
	})
	return receivedClaims, err
}

func (a *Agent) GetAllEmittedClaims() ([]*merkletree.Entry, error) {
	// get emitted claims, and generate fresh proof with current Root
	var emittedClaims []*merkletree.Entry
	err := a.storage.claims.emitted.Iterate(func(key, value []byte) (bool, error) {
		// TODO: Load claim + proof
		claim, err := merkletree.NewEntryFromBytes(value)
		if err != nil {
			return false, err
		}
		emittedClaims = append(emittedClaims, claim)
		return true, err

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

func (a *Agent) ExportMT() ([][2]string, error) {
	mt := [][2]string{}
	err := a.mt.Walk(nil, func(node *merkletree.Node) {
		mt = append(mt, [2]string{common3.HexEncode(node.Key()[:]),
			common3.HexEncode(node.Value())})
	})
	return mt, err
}

// GetCurrentRoot is used from wallet to check if it is syncronized with the
// MerkleTree in the IdentityAgent
func (a *Agent) GetCurrentRoot() *merkletree.Hash {
	return a.mt.RootKey()
}
