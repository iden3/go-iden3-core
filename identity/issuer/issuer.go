package issuer

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/iden3/go-iden3-core/components/idenpubonchain"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/core/genesis"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/keystore"
	"github.com/iden3/go-iden3-core/merkletree"

	"github.com/iden3/go-iden3-crypto/babyjub"
)

var (
	ErrIdenPubOnChainNil         = fmt.Errorf("idenPubOnChain is nil")
	ErrIdenStatePendingNotNil    = fmt.Errorf("Update of the published IdenState is pending")
	ErrIdenStateOnChainZero      = fmt.Errorf("No IdenState known to be on chain")
	ErrClaimNotFoundStateOnChain = fmt.Errorf("Claim not found under the on chain identity state")
)

var (
	dbPrefixClaimsTree     = []byte("treeclaims:")
	dbPrefixRevocationTree = []byte("treerevocation:")
	dbPrefixRootsTree      = []byte("treeroots:")
	dbPrefixIdenStateList  = []byte("idenstates:")
	dbKeyConfig            = []byte("config")
	dbKeyKOp               = []byte("kop")
	dbKeyId                = []byte("id")
	dbKeyNonceIdx          = []byte("nonceidx")
	// dbKeyIdenStateOnChain     = []byte("idenstateonchain")
	dbKeyIdenStateDataOnChain = []byte("idenstatedataonchain")
	dbKeyIdenStatePending     = []byte("idenstatepending")
	dbKeyEthTxSetState        = []byte("ethtxsetstate")
	dbKeyEthTxInitState       = []byte("ethtxinitstate")
)

var (
	SigPrefixSetState = []byte("setstate:")
)

// ConfigDefault is a default configuration for the Issuer.
var ConfigDefault = Config{MaxLevelsClaimsTree: 140, MaxLevelsRevocationTree: 140, MaxLevelsRootsTree: 140}

func storeJSON(tx db.Tx, key []byte, v interface{}) error {
	vJSON, err := json.Marshal(v)
	if err != nil {
		return err
	}
	tx.Put(key, vJSON)
	return nil
}

func loadJSON(storage db.Storage, key []byte, v interface{}) error {
	vJSON, err := storage.Get(key)
	if err == db.ErrNotFound {
		return nil
	} else if err != nil {
		return err
	}
	return json.Unmarshal(vJSON, v)
}

// Config allows configuring the creation of an Issuer.
type Config struct {
	MaxLevelsClaimsTree     int
	MaxLevelsRevocationTree int
	MaxLevelsRootsTree      int
}

// IdenStateTreeRoots is the set of the three roots of each Identity Merkle Tree.
type IdenStateTreeRoots struct {
	ClaimsRoot      *merkletree.Hash
	RevocationsRoot *merkletree.Hash
	RootsRoot       *merkletree.Hash
}

// TODO: Add mutex!

// Issuer is an identity that issues claims
type Issuer struct {
	rw            *sync.RWMutex
	storage       db.Storage
	id            *core.ID
	claimsMt      *merkletree.MerkleTree
	revocationsMt *merkletree.MerkleTree
	rootsMt       *merkletree.MerkleTree
	// idenPubOnChain can be nil if the identity doesn't connect to the blockchain.
	idenPubOnChain idenpubonchain.IdenPubOnChainer
	keyStore       *keystore.KeyStore
	kOpComp        *babyjub.PublicKeyComp
	nonceGen       *UniqueNonceGen
	idenStateList  *StorageList
	// _idenStateOnChain     *merkletree.Hash
	// idenStateDataOnChain is the last known identity state checked to be
	// in the Smart Contract.
	_idenStateDataOnChain *proof.IdenStateData
	// idenStatePending is a newly calculated identity state that is being
	// published in the Smart Contract but the transaction to publish it is
	// still pending.
	_idenStatePending *merkletree.Hash
	_ethTxSetState    *types.Transaction
	_ethTxInitState   *types.Transaction
	cfg               Config
}

//
// Persistence setters and getters
//

func (is *Issuer) idenStateDataOnChain() *proof.IdenStateData { return is._idenStateDataOnChain }

func (is *Issuer) setIdenStateDataOnChain(tx db.Tx, v *proof.IdenStateData) error {
	is._idenStateDataOnChain = v
	return storeJSON(tx, dbKeyIdenStateDataOnChain, v)
}

func (is *Issuer) loadIdenStateDataOnChain() error {
	is._idenStateDataOnChain = &proof.IdenStateData{}
	return loadJSON(is.storage, dbKeyIdenStateDataOnChain, is._idenStateDataOnChain)
}

func (is *Issuer) idenStateOnChain() *merkletree.Hash {
	return is._idenStateDataOnChain.IdenState
}

func (is *Issuer) idenStatePending() *merkletree.Hash {
	return is._idenStatePending
}

func (is *Issuer) setIdenStatePending(tx db.Tx, v *merkletree.Hash) {
	is._idenStatePending = v
	tx.Put(dbKeyIdenStatePending, v[:])
}

func (is *Issuer) loadIdenStatePending() error {
	b, err := is.storage.Get(dbKeyIdenStatePending)
	if err != nil {
		return err
	}
	var v merkletree.Hash
	copy(v[:], b)
	is._idenStatePending = &v
	return nil
}

// func (is *Issuer) ethTxSetState() *types.Transaction { return is._ethTxSetState }

func (is *Issuer) setEthTxSetState(tx db.Tx, v *types.Transaction) error {
	is._ethTxSetState = v
	return storeJSON(tx, dbKeyEthTxSetState, v)
}

func (is *Issuer) loadEthTxSetState() error {
	is._ethTxSetState = &types.Transaction{}
	return loadJSON(is.storage, dbKeyEthTxSetState, is._ethTxSetState)
}

// func (is *Issuer) ethTxInitState() *types.Transaction { return is._ethTxInitState }

func (is *Issuer) setEthTxInitState(tx db.Tx, v *types.Transaction) error {
	is._ethTxInitState = v
	return storeJSON(tx, dbKeyEthTxInitState, v)
}

func (is *Issuer) loadEthTxInitState() error {
	is._ethTxInitState = &types.Transaction{}
	return loadJSON(is.storage, dbKeyEthTxInitState, is._ethTxInitState)
}

// loadMTs loads the three identity merkle trees from the storage using the configuration.
func loadMTs(cfg *Config, storage db.Storage) (*merkletree.MerkleTree, *merkletree.MerkleTree, *merkletree.MerkleTree, error) {
	cltStorage := storage.WithPrefix(dbPrefixClaimsTree)
	retStorage := storage.WithPrefix(dbPrefixRevocationTree)
	rotStorage := storage.WithPrefix(dbPrefixRootsTree)

	clt, err := merkletree.NewMerkleTree(cltStorage, cfg.MaxLevelsClaimsTree)
	if err != nil {
		return nil, nil, nil, err
	}
	ret, err := merkletree.NewMerkleTree(retStorage, cfg.MaxLevelsRevocationTree)
	if err != nil {
		return nil, nil, nil, err
	}
	rot, err := merkletree.NewMerkleTree(rotStorage, cfg.MaxLevelsRootsTree)
	if err != nil {
		return nil, nil, nil, err
	}

	return clt, ret, rot, nil
}

// New creates a new Issuer, creating a new genesis ID and initializes the storages.
func New(cfg Config, kOpComp *babyjub.PublicKeyComp, extraGenesisClaims []merkletree.Entrier, storage db.Storage, keyStore *keystore.KeyStore, idenPubOnChain idenpubonchain.IdenPubOnChainer) (*Issuer, error) {
	clt, ret, rot, err := loadMTs(&cfg, storage)
	if err != nil {
		return nil, err
	}

	tx, err := storage.NewTx()

	if err != nil {
		return nil, err
	}

	// Initialize the UniqueNonceGen to generate revocation nonces for claims.
	nonceGen := NewUniqueNonceGen(NewStorageValue(dbKeyNonceIdx))
	nonceGen.Init(tx)

	// Create the Claim to authorize the Operational Key (kOp)
	kOp, err := kOpComp.Decompress()
	if err != nil {
		return nil, err
	}
	nonce, err := nonceGen.Next(tx)
	if err != nil {
		return nil, err
	}
	claimKOp := claims.NewClaimAuthorizeKSignBabyJub(kOp, nonce)
	id, _, err := genesis.CalculateIdGenesisMT(clt, rot, claimKOp, extraGenesisClaims)
	if err != nil {
		return nil, err
	}

	tx.Put(dbKeyId, id[:])
	tx.Put(dbKeyKOp, kOpComp[:])

	cfgJSON, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	tx.Put(dbKeyConfig, cfgJSON)

	idenStateList := NewStorageList(dbPrefixIdenStateList)

	is := Issuer{
		rw:             &sync.RWMutex{},
		id:             id,
		claimsMt:       clt,
		revocationsMt:  ret,
		rootsMt:        rot,
		idenPubOnChain: idenPubOnChain,
		// idenStateWriter: idenStateWriter,
		keyStore:      keyStore,
		kOpComp:       kOpComp,
		storage:       storage,
		nonceGen:      nonceGen,
		idenStateList: idenStateList,
		cfg:           cfg,
	}

	// Initalize the history of idenStates
	idenState, idenStateTreeRoots := is.state()
	idenStateList.Init(tx)

	if err := idenStateList.Append(tx, idenState[:], &idenStateTreeRoots); err != nil {
		return nil, err
	}

	// Initialize IdenStateDataOnChain and IdenStatePending to zero (writes to storage).
	if err := is.setIdenStateDataOnChain(tx, &proof.IdenStateData{IdenState: &merkletree.HashZero}); err != nil {
		return nil, err
	}
	is.setIdenStatePending(tx, &merkletree.HashZero)

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &is, nil

}

// Load creates an Issuer by loading a previously created Issuer (with New).
func Load(storage db.Storage, keyStore *keystore.KeyStore, idenPubOnChain idenpubonchain.IdenPubOnChainer) (*Issuer, error) {
	var cfg Config
	cfgJSON, err := storage.Get(dbKeyConfig)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(cfgJSON, &cfg); err != nil {
		return nil, err
	}

	kOpCompBytes, err := storage.Get(dbKeyConfig)
	if err != nil {
		return nil, err
	}
	var kOpComp babyjub.PublicKeyComp
	copy(kOpComp[:], kOpCompBytes)

	var id core.ID
	idBytes, err := storage.Get(dbKeyId)
	if err != nil {
		return nil, err
	}
	copy(id[:], idBytes)

	clt, ret, rot, err := loadMTs(&cfg, storage)
	if err != nil {
		return nil, err
	}

	nonceGen := NewUniqueNonceGen(NewStorageValue(dbKeyNonceIdx))
	idenStateList := NewStorageList(dbPrefixIdenStateList)

	is := Issuer{
		rw:             &sync.RWMutex{},
		id:             &id,
		claimsMt:       clt,
		revocationsMt:  ret,
		rootsMt:        rot,
		idenPubOnChain: idenPubOnChain,
		keyStore:       keyStore,
		kOpComp:        &kOpComp,
		storage:        storage,
		nonceGen:       nonceGen,
		idenStateList:  idenStateList,
		cfg:            cfg,
	}

	if err := is.loadIdenStateDataOnChain(); err != nil {
		return nil, err
	}
	if err := is.loadIdenStatePending(); err != nil {
		return nil, err
	}
	if err := is.loadEthTxInitState(); err != nil {
		return nil, err
	}
	if err := is.loadEthTxSetState(); err != nil {
		return nil, err
	}

	if err := is.SyncIdenStatePublic(); err != nil {
		if err != ErrIdenPubOnChainNil {
			return nil, err
		}
	}
	return &is, nil
}

// state returns the current Identity State and the three merkle tree roots.
func (is *Issuer) state() (*merkletree.Hash, IdenStateTreeRoots) {
	clr, rer, ror := is.claimsMt.RootKey(), is.revocationsMt.RootKey(), is.rootsMt.RootKey()
	idenState := core.IdenState(clr, rer, ror)
	return idenState, IdenStateTreeRoots{
		ClaimsRoot:      clr,
		RevocationsRoot: rer,
		RootsRoot:       ror,
	}
}

// State calculates and returns the current Identity State and the three merkle tree roots.
func (is *Issuer) State() (*merkletree.Hash, IdenStateTreeRoots) {
	is.rw.RLock()
	defer is.rw.RUnlock()
	return is.state()
}

// StateDataOnChain returns the last known IdentityState Data known to be on chain.
func (is *Issuer) StateDataOnChain() *proof.IdenStateData {
	is.rw.RLock()
	defer is.rw.RUnlock()
	return is.idenStateDataOnChain()
}

// ID returns the Issuer ID (Identity ID).
func (is *Issuer) ID() *core.ID {
	return is.id
}

// SyncIdenStatePublic updates the IdenStateOnChain and IdenStatePending from
// the values in the Smart Contract.
func (is *Issuer) SyncIdenStatePublic() error {
	if is.idenPubOnChain == nil {
		return ErrIdenPubOnChainNil
	}
	is.rw.Lock()
	defer is.rw.Unlock()
	idenStateData, err := is.idenPubOnChain.GetState(is.id)
	if err != nil {
		return err
	}
	if is.idenStatePending().Equals(&merkletree.HashZero) {
		// If there's no IdenState pending to be set on chain, the
		// obtained one must be the idenStateOnChain (Zero for genesis
		// / empty in the smart contract).
		if idenStateData.IdenState.Equals(is.idenStateOnChain()) {
			return nil
		}

		return fmt.Errorf("Fatal error: Identity State in the Smart Contract (%v)"+
			" doesn't match the expected OnChain one (%v).",
			idenStateData.IdenState, is.idenStateOnChain())
	}
	// If there's an IdenState pending to be set on chain, the
	// obtained one can be:

	// a. the idenStateOnchan (in this case, we still have an
	// IdenState pending to be set on chain).
	if idenStateData.IdenState.Equals(is.idenStateOnChain()) {
		return nil
	}

	// b. the idenStatePending (in this case, we no longer have an
	// IdenState pending and it becomes the idenStateOnChain, so we update
	// the sync state).
	if idenStateData.IdenState.Equals(is.idenStatePending()) {
		tx, err := is.storage.NewTx()
		if err != nil {
			return err
		}
		is.setIdenStatePending(tx, &merkletree.HashZero)
		if err := is.setIdenStateDataOnChain(tx, idenStateData); err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		return nil
	}

	// c. Neither the idenStatePending nor the idenStateOnchain
	// (unexpected result).
	return fmt.Errorf("Fatal error: Identity State in the Smart Contract (%v)"+
		" doesn't match the Pending one (%v) nor the OnChain one (%v).",
		idenStateData.IdenState, is.idenStatePending(), is.idenStateOnChain())
}

// IssueClaim adds a new claim to the Claims Merkle Tree of the Issuer.  The
// Identity State is not updated.
func (is *Issuer) IssueClaim(claim merkletree.Entrier) error {
	is.rw.Lock()
	defer is.rw.Unlock()
	if is.idenPubOnChain == nil {
		return ErrIdenPubOnChainNil
	}
	err := is.claimsMt.AddClaim(claim)
	if err != nil {
		return err
	}
	return nil
}

// getIdenStateByIdx gets identity state and identity state tree roots of the
// Issuer from the stored list at index idx.
func (is *Issuer) getIdenStateByIdx(tx db.Tx, idx uint32) (*merkletree.Hash, *IdenStateTreeRoots, error) {
	var idenStateTreeRoots IdenStateTreeRoots
	idenStateBytes, err := is.idenStateList.GetByIdx(tx, idx, &idenStateTreeRoots)
	if err != nil {
		return nil, nil, err
	}
	var idenState merkletree.Hash
	copy(idenState[:], idenStateBytes)
	return &idenState, &idenStateTreeRoots, nil
}

// getIdenStateTreeRoots gets the identity state tree roots of the Issuer from
// the stored list by identity state.
func (is *Issuer) getIdenStateTreeRoots(tx db.Tx, idenState *merkletree.Hash) (*IdenStateTreeRoots, error) {
	var idenStateTreeRoots IdenStateTreeRoots
	if err := is.idenStateList.Get(tx, idenState[:], &idenStateTreeRoots); err != nil {
		return nil, err
	}
	return &idenStateTreeRoots, nil
}

// PublishState calculates the current Issuer identity state, and if it's
// different than the last one, it publishes in in the blockchain.
func (is *Issuer) PublishState() error {
	is.rw.Lock()
	defer is.rw.Unlock()
	if is.idenPubOnChain == nil {
		return ErrIdenPubOnChainNil
	}
	if !is.idenStatePending().Equals(&merkletree.HashZero) {
		return ErrIdenStatePendingNotNil
	}
	idenState, idenStateTreeRoots := is.state()

	tx, err := is.storage.NewTx()
	if err != nil {
		return err
	}

	idenStateListLen, err := is.idenStateList.Length(tx)
	if err != nil {
		return err
	}
	idenStateLast, _, err := is.getIdenStateByIdx(tx, idenStateListLen-1)
	if err != nil {
		return err
	}

	if idenState.Equals(idenStateLast) {
		// IdenState hasn't changed, there's no need to do anything!
		return nil
	}

	if err := is.idenStateList.Append(tx, idenState[:], &idenStateTreeRoots); err != nil {
		return err
	}

	// Sign [minor] identity transition from last state to new (current) state.
	sig, err := is.SignBinary(SigPrefixSetState, append(idenStateLast[:], idenState[:]...))
	if err != nil {
		return err
	}

	if is.idenStateOnChain().Equals(&merkletree.HashZero) {
		// Identity State not present in the Smart Contract. First time
		// publishing it.
		ethTx, err := is.idenPubOnChain.InitState(is.id, idenStateLast, idenState, nil, nil, sig)
		if err != nil {
			return err
		}

		if err := is.setEthTxInitState(tx, ethTx); err != nil {
			return err
		}
	} else {
		// Identity State already present in the Smart Contract.
		// Update it.
		ethTx, err := is.idenPubOnChain.SetState(is.id, idenState, nil, nil, sig)
		if err != nil {
			return err
		}

		if err := is.setEthTxSetState(tx, ethTx); err != nil {
			return err
		}
	}

	is.setIdenStatePending(tx, idenState)

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// RevokeClaim revokes an already issued claim.
func (is *Issuer) RevokeClaim(claim merkletree.Entrier) error {
	if is.idenPubOnChain == nil {
		return ErrIdenPubOnChainNil
	}
	is.rw.Lock()
	defer is.rw.Unlock()
	data, err := is.claimsMt.GetDataByIndex(claim.Entry().HIndex())
	if err != nil {
		return err
	}
	nonce := claims.GetRevocationNonce(&merkletree.Entry{Data: *data})

	if err := claims.AddLeafRevocationsTree(is.revocationsMt, nonce, 0xffffffff); err != nil {
		return err
	}
	return nil
}

// UpdateClaim allows updating the value of an already issued claim.
func (is *Issuer) UpdateClaim(hIndex *merkletree.Hash, value []merkletree.ElemBytes) error {
	if is.idenPubOnChain == nil {
		return ErrIdenPubOnChainNil
	}
	return fmt.Errorf("TODO")
}

// Sign signs a message by the kOp of the issuer.
func (is *Issuer) Sign(string) (string, error) {
	return "", fmt.Errorf("TODO")
}

// Sign signs a binary message by the kOp of the issuer.
func (is *Issuer) SignBinary(prefix, msg []byte) (*babyjub.SignatureComp, error) {
	return is.keyStore.SignRaw(is.kOpComp, append(prefix, msg...))
}

func generateExistenceMTProof(mt *merkletree.MerkleTree, hi, root *merkletree.Hash) (*merkletree.Proof, error) {
	mtp, err := mt.GenerateProof(hi, root)
	if err != nil {
		return nil, err
	}
	if !mtp.Existence {
		return nil, ErrClaimNotFoundStateOnChain
	}
	return mtp, nil
}

// GenCredentialExistence generates an existence credential (claim + proof of
// existence) of an issued claim.  The result contains all data necessary to
// validate the credential against the Identity State found in the blockchain.
// For now, there are no genesis credentials.
func (is *Issuer) GenCredentialExistence(claim merkletree.Entrier) (*proof.CredentialExistence, error) {
	tx, err := is.storage.NewTx()
	if err != nil {
		return nil, err
	}
	is.rw.RLock()
	defer is.rw.RUnlock()
	idenStateData := is.idenStateDataOnChain()
	if idenStateData.IdenState.Equals(&merkletree.HashZero) {
		return nil, ErrIdenStateOnChainZero
	}
	idenStateTreeRoots, err := is.getIdenStateTreeRoots(tx, idenStateData.IdenState)
	if err != nil {
		return nil, err
	}
	mtpExist, err := generateExistenceMTProof(is.claimsMt, claim.Entry().HIndex(), idenStateTreeRoots.ClaimsRoot)
	if err != nil {
		return nil, err
	}
	return &proof.CredentialExistence{
		Id:              is.id,
		IdenStateData:   *idenStateData,
		MtpClaim:        mtpExist,
		Claim:           claim.Entry(),
		RevocationsRoot: idenStateTreeRoots.RevocationsRoot,
		RootsRoot:       idenStateTreeRoots.RootsRoot,
		IdPub:           "http://TODO",
	}, nil
}
