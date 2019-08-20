package identityagentsrv

import (
	"bytes"
	"encoding/hex"
	"errors"
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

type IdStorages struct {
	storage db.Storage
	ecSto   db.Storage
	rcSto   db.Storage
	mt      *merkletree.MerkleTree
}

// LoadPrefixStorage returns the identity storages
func (ia *Service) LoadIdStorages(id *core.ID) (*IdStorages, error) {
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
func (ia *Service) CreateIdentity(claimAuthKOp *merkletree.Entry,
	extraGenesisClaims []*merkletree.Entry) (*core.ID, *core.ProofClaim, error) {
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
	err = idStorages.mt.Add(claimAuthKOp)
	if err != nil {
		return nil, nil, err
	}
	tx.Put(claimAuthKOp.HIndex().Bytes(), claimAuthKOp.Bytes())
	for _, claim := range extraGenesisClaims {
		err = idStorages.mt.Add(claim)
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

func (ia *Service) AddClaim(id *core.ID, claim *merkletree.Entry) error {
	return ia.AddClaims(id, []*merkletree.Entry{claim})
}

func (ia *Service) AddClaims(id *core.ID, claims []*merkletree.Entry) error {
	idStorages, err := ia.LoadIdStorages(id)
	if err != nil {
		return err
	}

	tx, err := idStorages.ecSto.NewTx()
	for _, claim := range claims {
		err = idStorages.mt.Add(claim)
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

func (ia *Service) GetAllReceivedClaims(id *core.ID, idStorages *IdStorages) ([]core.ClaimObj, error) {
	var err error
	if idStorages == nil {
		idStorages, err = ia.LoadIdStorages(id)
		if err != nil {
			return []core.ClaimObj{}, err
		}
	}

	// get received claims
	var receivedClaims []core.ClaimObj
	idStorages.rcSto.Iterate(func(key, value []byte) {
		// filter only the key-value with the prefix of id+emittedclaims
		// as the Iterate from the Storage don't filters by prefix
		// in the future do a more efficient way to filter without going through all the keys-values
		var prefix []byte
		prefix = append(prefix[:], id.Bytes()...)
		prefix = append(prefix[:], PREFIX_RECEIVEDCLAIMS...)
		if bytes.Equal(key[:len(prefix)], prefix) {
			rClaim := core.ClaimObj{
				// TODO to be defined the way to store received claims
				Claim: &core.ClaimBasic{},
				Proof: core.ProofClaimPartial{},
			}
			receivedClaims = append(receivedClaims, rClaim)
		}
	})

	return receivedClaims, err
}

func (ia *Service) GetAllEmittedClaims(id *core.ID, idStorages *IdStorages) ([]core.ClaimObj, error) {
	var err error
	if idStorages == nil {
		idStorages, err = ia.LoadIdStorages(id)
		if err != nil {
			return []core.ClaimObj{}, err
		}
	}

	// get emitted claims, and generate fresh proof with current Root
	var emittedClaims []core.ClaimObj
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
			// once RootUpdater is ready, tie the claim proof in
			// the identity merkletree. Also the mtpNonRevokated
			// depends on the RootUpdater
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
			eClaim := core.ClaimObj{
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

func (ia *Service) GetAllClaims(id *core.ID) ([]core.ClaimObj, []core.ClaimObj, error) {
	idStorages, err := ia.LoadIdStorages(id)
	if err != nil {
		return []core.ClaimObj{}, []core.ClaimObj{}, err
	}

	// get received claims
	receivedClaims, err := ia.GetAllReceivedClaims(id, idStorages)
	if err != nil {
		return []core.ClaimObj{}, []core.ClaimObj{}, err
	}

	// get emitted claims, and generate fresh proof with current Root
	emittedClaims, err := ia.GetAllEmittedClaims(id, idStorages)
	return emittedClaims, receivedClaims, err

}

func (ia *Service) GetClaimByHi(id *core.ID, hi *merkletree.Hash) (merkletree.Claim, *core.ProofClaimPartial, error) {
	idStorages, err := ia.LoadIdStorages(id)
	if err != nil {
		return nil, nil, err
	}
	mtp, err := idStorages.mt.GenerateProof(hi, nil)
	// TODO see issues #167 and #169
	// once RootUpdater is ready, tie the claim proof in the identity
	// merkletree. Also the mtpNonRevokated depends on the RootUpdater
	// with the proof of the SetRootClaim in the Relay's merkletree
	proof := core.ProofClaimPartial{
		Mtp0: mtp,
		Mtp1: &merkletree.Proof{},
		Root: idStorages.mt.RootKey(),
		Aux:  nil,
	}
	// var leafBytes [merkletree.ElemBytesLen * merkletree.DataLen]byte
	// copy(leafBytes[:], value[:merkletree.ElemBytesLen*merkletree.DataLen])
	// leafData := merkletree.NewDataFromBytes(leafBytes)
	leafData, err := idStorages.mt.GetDataByIndex(hi)
	if err != nil {
		return nil, nil, err
	}
	entry := merkletree.Entry{
		Data: *leafData,
	}
	claim, err := core.NewClaimFromEntry(&entry)
	if err != nil {
		return nil, nil, err
	}
	return claim, &proof, nil
}

func (ia *Service) GetFullMT(id *core.ID) (map[string]string, error) {
	mt := make(map[string]string)
	idStorages, err := ia.LoadIdStorages(id)
	if err != nil {
		return mt, err
	}

	idStorages.ecSto.Iterate(func(key, value []byte) {
		// filter only the key-value with the prefix of id+emittedclaims
		// as the Iterate from the Storage don't filters by prefix
		// in the future do a more efficient way to filter without going through all the keys-values
		var prefix []byte
		prefix = append(prefix[:], id.Bytes()...)
		prefix = append(prefix[:], merkletree.PREFIX_MERKLETREE...)
		if bytes.Equal(key[:len(prefix)], prefix) {
			mt["0x"+hex.EncodeToString(key[len(prefix):])] = "0x" + hex.EncodeToString(value)
		}
	})
	return mt, nil
}

// GetCurrentRoot is used from wallet to check if it is syncronized with the
// MerkleTree in the IdentityAgent
func (ia *Service) GetCurrentRoot(id *core.ID) (*merkletree.Hash, error) {
	idStorages, err := ia.LoadIdStorages(id)
	if err != nil {
		return &merkletree.Hash{}, err
	}

	return idStorages.mt.RootKey(), nil
}
