package issuer

import (
	"encoding/json"
	"fmt"

	"github.com/iden3/go-iden3-core/components/idenstatereader"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/core/claims"
	"github.com/iden3/go-iden3-core/core/genesis"
	"github.com/iden3/go-iden3-core/core/proof"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/keystore"
	"github.com/iden3/go-iden3-core/merkletree"

	//"github.com/iden3/go-iden3-core/services/idenstatewriter"
	"github.com/iden3/go-iden3-crypto/babyjub"
)

var (
	ErrIdenStateWriterNil = fmt.Errorf("IdenStateWriter is nil")
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
)

// Config allows configuring the creation of an Issuer.
type Config struct {
	MaxLevelsClaimsTree     int
	MaxLevelsRevocationTree int
	MaxLevelsRootsTree      int
}

// ConfigDefault is a default configuration for the Issuer.
var ConfigDefault = Config{MaxLevelsClaimsTree: 140, MaxLevelsRevocationTree: 140, MaxLevelsRootsTree: 140}

// IdenStateTreeRoots is the set of the three roots of each Identity Merkle Tree.
type IdenStateTreeRoots struct {
	ClaimsRoot      *merkletree.Hash
	RevocationsRoot *merkletree.Hash
	RootsRoot       *merkletree.Hash
}

// Issuer is an identity that issues claims
type Issuer struct {
	id              *core.ID
	claimsMt        *merkletree.MerkleTree
	revocationsMt   *merkletree.MerkleTree
	rootsMt         *merkletree.MerkleTree
	idenStateReader idenstatereader.IdenStateReader
	// idenStateWriter idenstatewriter.IdenStateWriter
	keyStore        *keystore.KeyStore
	kOpComp         *babyjub.PublicKeyComp
	storage         db.Storage
	nonceGen        *UniqueNonceGen
	idenStateList   *StorageList
	idenStatePublic *merkletree.Hash
	cfg             Config
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
func New(cfg Config, kOpComp *babyjub.PublicKeyComp, extraGenesisClaims []merkletree.Entrier, storage db.Storage, keyStore *keystore.KeyStore, idenStateReader idenstatereader.IdenStateReader) (*Issuer, error) {
	clt, ret, rot, err := loadMTs(&cfg, storage)
	if err != nil {
		return nil, err
	}

	tx, err := storage.NewTx()

	if err != nil {
		return nil, err
	}

	nonceGen := NewUniqueNonceGen(NewStorageValue(dbKeyNonceIdx))
	nonceGen.Init(tx)

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

	issuer := &Issuer{
		id:              id,
		claimsMt:        clt,
		revocationsMt:   ret,
		rootsMt:         rot,
		idenStateReader: idenStateReader,
		// idenStateWriter: idenStateWriter,
		keyStore:        keyStore,
		kOpComp:         kOpComp,
		storage:         storage,
		nonceGen:        nonceGen,
		idenStateList:   idenStateList,
		idenStatePublic: nil,
		cfg:             cfg,
	}

	idenState, idenStateTreeRoots := issuer.state()
	idenStateList.Init(tx)

	if err := idenStateList.Append(tx, idenState[:], &idenStateTreeRoots); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return issuer, nil

}

// Load creates an Issuer by loading a previously created Issuer (with New).
func Load(storage db.Storage, keyStore *keystore.KeyStore, idenStateReader idenstatereader.IdenStateReader) (*Issuer, error) {
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

	return &Issuer{
		id:            &id,
		claimsMt:      clt,
		revocationsMt: ret,
		rootsMt:       rot,
		// idenStateWriter: idenStateWriter,
		idenStateReader: idenStateReader,
		keyStore:        keyStore,
		kOpComp:         &kOpComp,
		storage:         storage,
		nonceGen:        nonceGen,
		idenStateList:   idenStateList,
		cfg:             cfg,
	}, nil
}

// state returns the current Identity State and the three merkle tree roots.
func (i *Issuer) state() (*merkletree.Hash, IdenStateTreeRoots) {
	clr, rer, ror := i.claimsMt.RootKey(), i.revocationsMt.RootKey(), i.rootsMt.RootKey()
	idenState := core.IdenState(clr, rer, ror)
	return idenState, IdenStateTreeRoots{
		ClaimsRoot:      clr,
		RevocationsRoot: rer,
		RootsRoot:       ror,
	}
}

// ID returns the Issuer ID (Identity ID).
func (i *Issuer) ID() *core.ID {
	return i.id
}

func (i *Issuer) SyncIdenStatePublic() error {
	// i.idenStateReader.GetState(i.id)
	return fmt.Errorf("TODO")
}

// GenCredentialExistence generates an existence credential (claim + proof of
// existence) of an issued claim.  The result contains all data necessary to
// validate the credential against the Identity State found in the blockchain.
func (i *Issuer) GenCredentialExistence(claim merkletree.Entrier) (*proof.CredentialExistence, error) {
	// idenState, err := m.idenStateWriter.GetIdenState(i.id)
	// if err != nil {
	// 	return nil, err
	// }
	// var clr *merkletree.Hash // TODO
	// clt, err := m.claimsMt.Snapshot(clr)
	// if err != nil {
	// 	return nil, err
	// }
	// proofClaim, err := proof.GetClaimProofByHi(mt, hi)
	// if err != nil {
	// 	return nil, err
	// }
	// proofClaim.ID = i.id
	// proofClaim.BlockN = idenState.BlockN
	// proofClaim.BlockTimestamp = idenState.BlockTs
	// return proofClaim, nil
	return nil, fmt.Errorf("TODO")
}

// IssueClaim adds a new claim to the Claims Merkle Tree of the Issuer.  The
// Identity State is not updated.
func (i *Issuer) IssueClaim(claim merkletree.Entrier) error {
	if i.idenStateReader == nil {
		return ErrIdenStateWriterNil
	}
	err := i.claimsMt.AddClaim(claim)
	if err != nil {
		return err
	}
	return nil
}

// getIdenStateByIdx gets identity state and identity state tree roots of the
// Issuer from the stored list at index idx.
func (i *Issuer) getIdenStateByIdx(tx db.Tx, idx uint32) (*merkletree.Hash, *IdenStateTreeRoots, error) {
	var idenStateTreeRoots IdenStateTreeRoots
	idenStateBytes, err := i.idenStateList.GetByIdx(tx, idx, &idenStateTreeRoots)
	if err != nil {
		return nil, nil, err
	}
	var idenState merkletree.Hash
	copy(idenState[:], idenStateBytes)
	return &idenState, &idenStateTreeRoots, nil
}

// getIdenStateTreeRoots gets the identity state tree roots of the Issuer from
// the stored list by identity state.
// func (i *Issuer) getIdenStateTreeRoots(tx db.Tx, idenState *merkletree.Hash) (*IdenStateTreeRoots, error) {
// 	var idenStateTreeRoots IdenStateTreeRoots
// 	if err := i.idenStateList.Get(tx, idenState[:], &idenStateTreeRoots); err != nil {
// 		return nil, err
// 	}
// 	return &idenStateTreeRoots, nil
// }

// PublishState calculates the current Issuer identity state, and if it's
// different than the last one, it publishes in in the blockchain.
func (i *Issuer) PublishState() error {
	if i.idenStateReader == nil {
		return ErrIdenStateWriterNil
	}
	idenState, idenStateTreeRoots := i.state()

	tx, err := i.storage.NewTx()
	if err != nil {
		return err
	}
	idenStateListLen, err := i.idenStateList.Length(tx)
	if err != nil {
		return err
	}
	idenStateLast, _, err := i.getIdenStateByIdx(tx, idenStateListLen-1)
	if err != nil {
		return err
	}

	if idenState == idenStateLast {
		// IdenState hasn't changed, there's no need to do anything!
		return nil
	}

	if err := i.idenStateList.Append(tx, idenState[:], &idenStateTreeRoots); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return fmt.Errorf("TODO")
}

// RevokeClaim revokes an already issued claim.
func (i *Issuer) RevokeClaim(claim merkletree.Entrier) error {
	if i.idenStateReader == nil {
		return ErrIdenStateWriterNil
	}
	data, err := i.claimsMt.GetDataByIndex(claim.Entry().HIndex())
	if err != nil {
		return err
	}
	nonce := claims.GetRevocationNonce(&merkletree.Entry{Data: *data})

	if err := claims.AddLeafRevocationsTree(i.revocationsMt, nonce, 0xffffffff); err != nil {
		return err
	}
	return nil
}

// UpdateClaim allows updating the value of an already issued claim.
func (i *Issuer) UpdateClaim(hIndex *merkletree.Hash, value []merkletree.ElemBytes) error {
	if i.idenStateReader == nil {
		return ErrIdenStateWriterNil
	}
	return fmt.Errorf("TODO")
}

// Sign signs a message by the kOp of the issuer.
func (i *Issuer) Sign(string) (string, error) {
	return "", fmt.Errorf("TODO")
}

// Sign signs a binary message by the kOp of the issuer.
func (i *Issuer) SignBinary(string) (string, error) {
	return "", fmt.Errorf("TODO")
}
