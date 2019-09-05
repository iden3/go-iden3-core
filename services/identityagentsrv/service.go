package identityagentsrv

import (
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

var PREFIX_CLAIMSEMITTED = []byte("claimsemitted")
var PREFIX_CLAIMSRECEIVED = []byte("claimsreceived")
var PREFIX_CLAIMSGENESIS = []byte("claimsgenesis")
var PREFIX_MERKLETREE = []byte("merkletree")

// TODO: Move this to a generic place
type ServerError struct {
	Err string `json:"error"`
}

// TODO: Move this to a generic place
func (e ServerError) Error() string {
	return fmt.Sprintf("server: %v", e.Err)
}

type RootUpdaterMock struct{}

func (ru *RootUpdaterMock) RootUpdate(id *core.ID, setRootReq claimsrv.SetRoot0Req) error {
	return fmt.Errorf("mock mock")
}

func (ru *RootUpdaterMock) GetRootProof(id *core.ID) (*core.ProofClaim, error) {
	return nil, fmt.Errorf("mock mock")
}

func (ru *RootUpdaterMock) ClaimAuthService() *merkletree.Entry { return nil }

type RootUpdaterRelay struct {
	RelayUrl string
	RelayId  *core.ID
	_client  *sling.Sling
	validate *validator.Validate
}

func NewRootUpdaterRelay(relayUrl string, relayId *core.ID) RootUpdaterRelay {
	if relayUrl[len(relayUrl)-1] != '/' {
		relayUrl += "/"
	}
	client := sling.New().Base(relayUrl)
	return RootUpdaterRelay{RelayUrl: relayUrl, RelayId: relayId,
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

func (ru *RootUpdaterRelay) RootUpdate(id *core.ID, setRootReq claimsrv.SetRoot0Req) error {
	var setRootClaim struct {
		SetRootClaim *merkletree.Entry `json:"setRootClaim" validate:"required"`
	}
	path := fmt.Sprintf("ids/%s/setrootclaim", id)
	return ru.request(ru.client().Path(path).Post("").BodyJSON(setRootReq), &setRootClaim)
}

func (ru *RootUpdaterRelay) GetRootProof(id *core.ID) (*core.ProofClaim, error) {
	var proofClaim struct {
		ProofClaim *core.ProofClaim `json:"proofClaim" validate:"required"`
	}
	path := fmt.Sprintf("ids/%s/setrootclaim", id)
	err := ru.request(ru.client().Path(path).Get(""), &proofClaim)
	return proofClaim.ProofClaim, err
}

func (ru *RootUpdaterRelay) ClaimAuthService() *merkletree.Entry {
	return core.NewClaimAuthorizeService(core.ServiceTypeRelay, ru.RelayId.String(), "", ru.RelayUrl).Entry()
}

type RootUpdater interface {
	RootUpdate(id *core.ID, setRootReq claimsrv.SetRoot0Req) error
	GetRootProof(id *core.ID) (*core.ProofClaim, error)
	ClaimAuthService() *merkletree.Entry
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

// CreateIdentity creates a new identity from the given claims
func (ia *Service) CreateIdentity(claimAuthKOp *merkletree.Entry,
	extraGenesisClaims []*merkletree.Entry) (*core.ID, *core.ProofClaim, error) {
	claimAuthService := ia.rootUpdater.ClaimAuthService()
	if claimAuthService != nil {
		extraGenesisClaims = append(extraGenesisClaims)
	}
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
	tx0, err := agent.storage.claims.genesis.NewTx()
	if err != nil {
		return nil, nil, err
	}
	tx1, err := agent.storage.claims.emitted.NewTx()
	if err != nil {
		return nil, nil, err
	}

	if err = agent.mt.Add(claimAuthKOp); err != nil {
		return nil, nil, err
	}
	tx0.Put(claimAuthKOp.HIndex().Bytes(), claimAuthKOp.Bytes())
	tx1.Put(claimAuthKOp.HIndex().Bytes(), claimAuthKOp.Bytes())
	for _, claim := range extraGenesisClaims {
		// FIXME: If some agent.mt.Add fails, there will be an
		// inconsistency between agent.mt and agent.storage.claims.emitted
		err = agent.mt.Add(claim)
		if err != nil {
			return nil, nil, err
		}
		tx0.Put(claim.HIndex().Bytes(), claim.Bytes())
		tx1.Put(claim.HIndex().Bytes(), claim.Bytes())
	}
	if err = tx0.Commit(); err != nil {
		return nil, nil, err
	}
	if err = tx1.Commit(); err != nil {
		return nil, nil, err
	}

	// TODO send identity Root to RootUpdater (Relay)
	// this will be implemented when the Connection with RootUpdater is ready

	return id, proofKOp, nil
}

type IdStorage struct {
	claims struct {
		emitted  db.Storage
		received db.Storage
		genesis  db.Storage
	}
}

type Agent struct {
	rootUpdater RootUpdater
	storage     *IdStorage
	mt          *merkletree.MerkleTree
	id          *core.ID
}

// LoadPrefixStorage returns the identity storages
func (a *Agent) loadStorage(base db.Storage) error {
	emittedClaims := base.WithPrefix(PREFIX_CLAIMSEMITTED)
	receivedClaims := base.WithPrefix(PREFIX_CLAIMSRECEIVED)
	genesisClaims := base.WithPrefix(PREFIX_CLAIMSGENESIS)
	mtStorage := base.WithPrefix(PREFIX_MERKLETREE)
	a.storage = &IdStorage{
		claims: struct {
			emitted  db.Storage
			received db.Storage
			genesis  db.Storage
		}{
			emitted:  emittedClaims,
			received: receivedClaims,
			genesis:  genesisClaims,
		},
	}
	mt, err := merkletree.NewMerkleTree(mtStorage, 140)
	a.mt = mt
	return err
}

func (s *Service) NewAgent(id *core.ID) (*Agent, error) {
	agent := &Agent{id: id, rootUpdater: s.rootUpdater}
	err := agent.loadStorage(s.storage.WithPrefix(agent.id.Bytes()))
	return agent, err
}

// RootUpdate checks the signature and send the
func (a *Agent) RootUpdate(setRootReq claimsrv.SetRoot0Req) error {
	ok, err := claimsrv.CheckSetRootParams(a.id, setRootReq)
	if err != nil || !ok {
		return fmt.Errorf("SetRoot params verification not passed, " + err.Error())
	}

	return a.rootUpdater.RootUpdate(a.id, setRootReq)
}

func (a *Agent) GetRootProof(id *core.ID) (*core.ProofClaim, error) {
	return a.rootUpdater.GetRootProof(a.id)
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
	return nil
}

func (a *Agent) ClaimsGenesis() ([]*merkletree.Entry, error) {
	var genesisClaims []*merkletree.Entry
	err := a.storage.claims.genesis.Iterate(func(key, value []byte) (bool, error) {
		claim, err := merkletree.NewEntryFromBytes(value)
		if err != nil {
			return false, err
		}
		genesisClaims = append(genesisClaims, claim)
		return true, err
	})
	return genesisClaims, err
}

func (a *Agent) ClaimsReceived() ([]*merkletree.Entry, error) {
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

func (a *Agent) ClaimsEmitted() ([]*merkletree.Entry, error) {
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

type CurrentRoot struct {
	Local     *merkletree.Hash `json:"local"`
	Published *merkletree.Hash `json:"published"`
}

// GetCurrentRoot is used from wallet to check if it is syncronized with the
// MerkleTree in the IdentityAgent
func (a *Agent) GetCurrentRoot() (*CurrentRoot, error) {
	proofClaim, err := a.rootUpdater.GetRootProof(a.id)
	if err != nil {
		return nil, err
	}

	claim, err := core.NewClaimFromEntry(proofClaim.Claim)
	if err != nil {
		return nil, fmt.Errorf("Error parsing proofClaim.leaf: %v", err)
	}
	claimSetRootKey, ok := claim.(*core.ClaimSetRootKey)
	if !ok {
		return nil, fmt.Errorf("Error casting claim type for claim parsed from proofClaim.leaf")
	}

	return &CurrentRoot{
		Local:     a.mt.RootKey(),
		Published: &claimSetRootKey.RootKey,
	}, nil
}
