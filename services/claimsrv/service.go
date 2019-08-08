package claimsrv

import (
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/db"
	"github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/services/rootsrv"
	"github.com/iden3/go-iden3-core/services/signsrv"
	"github.com/iden3/go-iden3-core/utils"
)

var (
	ErrNotFound = errors.New("value not found")
)

type Service interface {
	CommitNewIdRoot(id core.ID, kSignPk *ecdsa.PublicKey, root merkletree.Hash, timestamp int64, signature *utils.SignatureEthMsg) (*core.ClaimSetRootKey, error)
	AddClaimAuthorizeKSignSecp256k1First(id core.ID,
		claimAuthorizeKSignSecp256k1 core.ClaimAuthorizeKSignSecp256k1) error
	AddUserIdClaim(id core.ID, claimValueMsg ClaimValueMsg) error
	AddClaim(claim merkletree.Entrier) error
	GetIdRoot(id core.ID) (merkletree.Hash, []byte, error)
	GetClaimProofUserByHi(id core.ID, hi *merkletree.Hash) (*core.ProofClaim, error)
	GetClaimProofUserByHiOld(id core.ID, hi merkletree.Hash) (*ProofClaimUser, error)
	GetClaimProofByHi(hi *merkletree.Hash) (*core.ProofClaim, error)
	MT() *merkletree.MerkleTree
	RootSrv() rootsrv.Service
}

type ServiceImpl struct {
	id      core.ID
	mt      *merkletree.MerkleTree
	rootsrv rootsrv.Service
	signer  signsrv.Service
}

func New(id core.ID, mt *merkletree.MerkleTree, rootsrv rootsrv.Service,
	signer signsrv.Service) *ServiceImpl {
	return &ServiceImpl{id, mt, rootsrv, signer}
}

// MT returns the merkle tree.
func (cs *ServiceImpl) MT() *merkletree.MerkleTree {
	return cs.mt
}

// RootSrv returns the RootService
func (cs *ServiceImpl) RootSrv() rootsrv.Service {
	return cs.rootsrv
}

// SetNewIdRoot checks that the data is valid and performs a claim in the Relay merkletree setting the new Root of the emiting Id
func (cs *ServiceImpl) CommitNewIdRoot(id core.ID, kSignPk *ecdsa.PublicKey, root merkletree.Hash, timestamp int64, signature *utils.SignatureEthMsg) (*core.ClaimSetRootKey, error) {
	// get the user's id storage, using the user id prefix (the id itself)
	stoUserId := cs.mt.Storage().WithPrefix(id.Bytes())

	// open the MerkleTree of the user
	userMT, err := merkletree.NewMerkleTree(stoUserId, 140)
	if err != nil {
		return &core.ClaimSetRootKey{}, err
	}

	// verify that the KSign is authorized
	if !CheckKSignInIddb(userMT, kSignPk) {
		return &core.ClaimSetRootKey{}, errors.New("can not verify the KSign")
	}
	// in the future the user merkletree will be in the client side, and this step will be a check of the ProofKSign

	// check data timestamp
	verified := utils.VerifyTimestamp(timestamp, 30000) //needs to be from last 30 seconds
	if !verified {
		return &core.ClaimSetRootKey{}, errors.New("timestamp too old")
	}
	// check signature with id
	// whee data signed is id+root+timestamp
	timestampBytes := utils.Uint64ToEthBytes(uint64(timestamp))
	// signature of id+root+timestamp, only valid if is from last X seconds
	var msg []byte
	msg = append(msg, id.Bytes()...)
	msg = append(msg, root.Bytes()...)
	msg = append(msg, timestampBytes...)
	if !utils.VerifySigEthMsg(crypto.PubkeyToAddress(*kSignPk), signature, msg) {
		return &core.ClaimSetRootKey{}, errors.New("signature can not be verified")
	}

	// claimSetRootKey of the user in the Relay Merkle Tree
	// create new ClaimSetRootKey
	claimSetRootKey, err := core.NewClaimSetRootKey(id, root)
	if err != nil {
		return nil, err
	}
	// entry := claimSetRootKey.Entry()
	// version, err := GetNextVersion(cs.mt, entry.HIndex())
	version, err := GetNextVersion(cs.mt, claimSetRootKey.Entry().HIndex())
	if err != nil {
		return &core.ClaimSetRootKey{}, err
	}
	claimSetRootKey.Version = version

	// add User's Id Merkle Root into the Relay's Merkle Tree
	e := claimSetRootKey.Entry()
	err = cs.mt.Add(e)
	if err != nil {
		return &core.ClaimSetRootKey{}, err
	}

	// update Relay Root in Smart Contract
	cs.rootsrv.SetRoot(*cs.mt.RootKey())

	return claimSetRootKey, nil
}

// TODO
// AddClaimAuthorizeKSignSecp256k1 adds ClaimAuthorizeKSignSecp256k1 into the Id's merkletree, and adds the Id's merkle root into the Relay's merkletree inside a ClaimSetRootKey. Returns the merkle proof of both Claims
//func (cs *ServiceImpl) AddClaimAuthorizeKSignSecp256k1(id common.Address, claimAuthorizeKSignMsg ClaimAuthorizeKSignMsg) error {
//
//	// get the user's id storage, using the user id prefix (the id itself)
//	stoUserId := cs.mt.Storage().WithPrefix(id.Bytes())
//
//	// open the MerkleTree of the user
//	userMT, err := merkletree.NewMerkleTree(stoUserId, 140)
//	if err != nil {
//		return err
//	}
//
//	// verify that the KSign is authorized
//	if !CheckKSignInIddb(userMT, claimAuthorizeKSignMsg.KSign) {
//		return errors.New("can not verify the KSign")
//	}
//
//	// verify signature of the ClaimAuthorizeKSign
//	signature, err := common3.HexDecode(claimAuthorizeKSignMsg.Signature)
//	if err != nil {
//		return err
//	}
//	if !utils.VerifySigEthMsg(claimAuthorizeKSignMsg.KSign, signature,
//		claimAuthorizeKSignMsg.ClaimAuthorizeKSign.Entry().Bytes()) {
//		return errors.New("signature can not be verified")
//	}
//
//	// add ClaimAuthorizeKSign into the User's Id Merkle Tree
//	err = userMT.Add(claimAuthorizeKSignMsg.ClaimAuthorizeKSign.Entry())
//	if err != nil {
//		return err
//	}
//
//	// create new ClaimSetRootKey
//	claimSetRootKey := core.NewClaimSetRootKey(id, *userMT.RootKey())
//
//	// get next version of the claim
//	version, err := GetNextVersion(cs.mt, claimSetRootKey.Entry().HIndex())
//	if err != nil {
//		return err
//	}
//	claimSetRootKey.Version = version
//
//	// add User's Id Merkle Root into the Relay's Merkle Tree
//	err = cs.mt.Add(claimSetRootKey.Entry())
//	if err != nil {
//		return err
//	}
//
//	// update Relay's Root in the Smart Contract
//	cs.rootsrv.SetRoot(*cs.mt.RootKey())
//
//	return nil
//}

// AddClaimAuthorizeKSignSecp256k1First adds ClaimAuthorizeKSignSecp256k1 into
// the Id's merkletree, and adds the Id's merkle root into the Relay's
// merkletree inside a ClaimSetRootKey. Returns the merkle proof of both Claims
func (cs *ServiceImpl) AddClaimAuthorizeKSignSecp256k1First(id core.ID,
	claimAuthorizeKSignSecp256k1 core.ClaimAuthorizeKSignSecp256k1) error {

	// get the user's id storage, using the user id prefix (the id itself)
	stoUserId := cs.mt.Storage().WithPrefix(id.Bytes())

	// open the MerkleTree of the user
	userMT, err := merkletree.NewMerkleTree(stoUserId, 140)
	if err != nil {
		return err
	}

	// add ClaimAuthorizeKSign into the User's Id Merkle Tree
	err = userMT.Add(claimAuthorizeKSignSecp256k1.Entry())
	if err != nil {
		return err
	}

	// create new ClaimSetRootKey
	claimSetRootKey, err := core.NewClaimSetRootKey(id, *userMT.RootKey())
	if err != nil {
		return err
	}

	// get next version of the claim
	version, err := GetNextVersion(cs.mt, claimSetRootKey.Entry().HIndex())
	if err != nil {
		return err
	}
	claimSetRootKey.Version = version

	// add User's Id Merkle Root into the Relay's Merkle Tree
	err = cs.mt.Add(claimSetRootKey.Entry())
	if err != nil {
		return err
	}

	// update Relay's Root in the Smart Contract
	cs.rootsrv.SetRoot(*cs.mt.RootKey())

	return nil
}

// AddUserIdClaim adds a claim into the Id's merkle tree, and with the Id's root, creates a new ClaimSetRootKey and adds it to the Relay's merkletree
func (cs *ServiceImpl) AddUserIdClaim(id core.ID, claimValueMsg ClaimValueMsg) error {
	// get the user's id storage, using the user id prefix (the id itself)
	stoUserId := cs.mt.Storage().WithPrefix(id.Bytes())

	// open the MerkleTree of the user
	userMT, err := merkletree.NewMerkleTree(stoUserId, 140)
	if err != nil {
		return err
	}

	// verify that the KSign is authorized
	if !CheckKSignInIddb(userMT, &claimValueMsg.KSignPk.PublicKey) {
		return errors.New("can not verify the KSign")
	}

	// verify signature with KSign
	if !utils.VerifySigEthMsg(crypto.PubkeyToAddress(claimValueMsg.KSignPk.PublicKey),
		claimValueMsg.Signature, claimValueMsg.ClaimValue.Bytes()) {
		return errors.New("signature can not be verified")
	}

	// add claim in User Id Merkle Tree
	err = userMT.Add(&claimValueMsg.ClaimValue)
	if err != nil {
		return err
	}

	// claimSetRootKey of the user in the Relay Merkle Tree
	// create new ClaimSetRootKey
	claimSetRootKey, err := core.NewClaimSetRootKey(id, *userMT.RootKey())
	if err != nil {
		return err
	}
	version, err := GetNextVersion(cs.mt, claimSetRootKey.Entry().HIndex())
	if err != nil {
		return err
	}
	claimSetRootKey.Version = version

	// add User's Id Merkle Root into the Relay's Merkle Tree
	err = cs.mt.Add(claimSetRootKey.Entry())
	if err != nil {
		return err
	}

	// update Relay Root in Smart Contract
	cs.rootsrv.SetRoot(*cs.mt.RootKey())

	return nil
}

// AddClaim adds a claim directly to the Relay merkletree
func (cs *ServiceImpl) AddClaim(claim merkletree.Entrier) error {
	err := cs.mt.Add(claim.Entry())
	if err != nil {
		return err
	}
	cs.rootsrv.SetRoot(*cs.mt.RootKey())
	return nil
}

// GetIdRoot returns the root of an Id tree, and the proof of that Root Id tree in the Relay Merkle Tree
func (cs *ServiceImpl) GetIdRoot(id core.ID) (merkletree.Hash, []byte, error) {
	// get the user's id storage, using the user id prefix (the id itself)
	stoUserId := cs.mt.Storage().WithPrefix(id.Bytes())

	// open the MerkleTree of the user
	userMT, err := merkletree.NewMerkleTree(stoUserId, 140)
	if err != nil {
		return merkletree.Hash{}, []byte{}, err
	}

	// build ClaimSetRootKey of the user id
	claimSetRootKey, err := core.NewClaimSetRootKey(id, *userMT.RootKey())
	if err != nil {
		return merkletree.Hash{}, []byte{}, err
	}

	version, err := GetNextVersion(cs.mt, claimSetRootKey.Entry().HIndex())
	if err != nil {
		return merkletree.Hash{}, []byte{}, err
	}
	claimSetRootKey.Version = version - 1

	// get proof of SetRootProof in the Relay tree
	idRootProof, err := cs.mt.GenerateProof(claimSetRootKey.Entry().HIndex(), nil)
	if err != nil {
		return merkletree.Hash{}, []byte{}, err
	}
	return *userMT.RootKey(), idRootProof.Bytes(), nil
}

// TODO: Remove this
// GetClaimProofUserByHiOld given a Hash(index) (Hi) and an Id, returns the Claim in that Hi position inside the Id's merkletree, and the ClaimSetRootKey with the Id's root in the Relay's merkletree
func (cs *ServiceImpl) GetClaimProofUserByHiOld(id core.ID, hi merkletree.Hash) (*ProofClaimUser, error) {
	// get the user's id storage, using the user id prefix (the id itself)
	stoUserId := cs.mt.Storage().WithPrefix(id.Bytes())

	// open the MerkleTree of the user
	userMT, err := merkletree.NewMerkleTree(stoUserId, 140)
	if err != nil {
		return nil, err
	}

	// get the value in the hi position
	// valueBytes, err := userMT.GetValueInPos(hi)
	leafData, err := userMT.GetDataByIndex(&hi)
	if err != nil {
		return nil, err
	}
	// if bytes.Equal(valueBytes, merkletree.EmptyNodeValue[:]) {
	//         return nil, ErrNotFound
	// }

	// value, err := core.ParseValueFromBytes(valueBytes)
	// if err != nil {
	//         return nil, err
	// }

	// get the proof of the value in the User Id Tree
	// idProof, err := userMT.GenerateProof(merkletree.HashBytes(value.Bytes()[:value.IndexLength()]), nil)
	idProof, err := userMT.GenerateProof(&hi, nil)
	if err != nil {
		return nil, err
	}

	leafBytes := leafData.Bytes()
	claimProof := ProofTreeLeaf{
		Leaf:  leafBytes[:],
		Proof: idProof.Bytes(),
		Root:  *userMT.RootKey(),
	}

	// build ClaimSetRootKey
	claimSetRootKey, err := core.NewClaimSetRootKey(id, *userMT.RootKey())
	if err != nil {
		return nil, err
	}
	version, err := GetNextVersion(cs.mt, claimSetRootKey.Entry().HIndex())
	if err != nil {
		return nil, err
	}
	claimSetRootKey.Version = version - 1

	// get the proof of the ClaimSetRootKey in the Relay Tree
	relayProof, err := cs.mt.GenerateProof(claimSetRootKey.Entry().HIndex(), nil)
	if err != nil {
		return nil, err
	}
	claimSetRootKeyProof := ProofTreeLeaf{
		Leaf:  claimSetRootKey.Entry().Bytes(),
		Proof: relayProof.Bytes(),
		Root:  *cs.mt.RootKey(),
	}

	// get non revocation proofs of the claim
	claimNonRevocationProof, err := getNonRevocationProof(userMT, hi)
	if err != nil {
		return nil, err
	}
	claimSetRootKeyNonRevocationProof, err := getNonRevocationProof(cs.mt, *claimSetRootKey.Entry().HIndex())
	if err != nil {
		return nil, err
	}

	// sign root + date
	sig, date, err := cs.signer.SignEthMsgDate(claimSetRootKeyProof.Root[:])
	if err != nil {
		return nil, err
	}

	proofClaim := ProofClaimUser{
		claimProof,
		claimNonRevocationProof,
		claimSetRootKeyProof,
		claimSetRootKeyNonRevocationProof,
		date,
		sig[:],
	}
	return &proofClaim, nil
}

// GetClaimProofUserByHi given a Hash(index) (Hi) and an id, returns the Claim
// in that Hi position inside the User merkletree, it's proof of existence and
// of non-revocation, and the proof of existence and of non-revocation for the
// set root claim in the relay tree, all in the form of a ProofClaim.
func (cs *ServiceImpl) GetClaimProofUserByHi(id core.ID,
	hi *merkletree.Hash) (*core.ProofClaim, error) {
	// open the MerkleTree of the user
	userMT, err := NewMerkleTreeUser(id, cs.mt.Storage(), 140)
	if err != nil {
		return nil, err
	}

	// get the value in the hi position
	leafData, err := userMT.GetDataByIndex(hi)
	if err != nil {
		return nil, err
	}

	// get the MT proof of existence of the claim and the non-existence of
	// the claim's next version in the User Tree
	mtpExistUser, err := userMT.GenerateProof(hi, nil)
	if err != nil {
		return nil, err
	}
	mtpNonExistUser, err := core.GetNonRevocationMTProof(userMT, leafData, hi)
	if err != nil {
		return nil, err
	}

	// build ClaimSetRootKey
	claimSetRootKey, err := core.NewClaimSetRootKey(id, *userMT.RootKey())
	if err != nil {
		return nil, err
	}
	// TODO in a future iteration: make an efficient implementation of GetNextVersion
	version, err := GetNextVersion(cs.mt, claimSetRootKey.Entry().HIndex())
	if err != nil {
		return nil, err
	}
	claimSetRootKey.Version = version - 1

	// Call GetClaimProofByHi to generate a Proof for the top level tree
	proofClaim, err := cs.GetClaimProofByHi(claimSetRootKey.Entry().HIndex())
	if err != nil {
		return nil, err
	}

	// Generate the partial claim proof for the user claim and add it to the ProofClaim
	proofClaimUserPartial := core.ProofClaimPartial{
		Mtp0: mtpExistUser,
		Mtp1: mtpNonExistUser,
		Root: userMT.RootKey(),
		Aux: &core.SetRootAux{
			Version: claimSetRootKey.Version,
			Era:     0, // NOTE: For the login milestone we don't support Era
			Id:      id,
		},
	}
	proofClaim.Proofs = []core.ProofClaimPartial{proofClaimUserPartial, proofClaim.Proofs[0]}
	proofClaim.Leaf = leafData

	return proofClaim, nil
}

// GetClaimProofByHi given a Hash(index) (Hi), returns the Claim in that Hi
// position inside the Relay merkletree, and it's proof of existence and of
// non-revocated, all in the form of a ProofClaim.  The result is signed (with
// a timestamp) by the service.
func (cs *ServiceImpl) GetClaimProofByHi(hi *merkletree.Hash) (*core.ProofClaim, error) {
	mt, err := cs.mt.Snapshot(cs.mt.RootKey())
	if err != nil {
		return nil, err
	}
	proofClaim, err := core.GetClaimProofByHi(mt, hi)
	if err != nil {
		return nil, err
	}

	sig, date, err := cs.signer.SignEthMsgDate(proofClaim.Proofs[0].Root[:])
	if err != nil {
		return nil, err
	}
	proofClaim.Signer, proofClaim.Signature, proofClaim.Date = cs.id, sig, date

	return proofClaim, nil
}

// getNonRevocationProof returns the next version Hi (that don't exist in the tree, it's value is Empty) with merkleproof and root
func getNonRevocationProof(mt *merkletree.MerkleTree, hi merkletree.Hash) (ProofTreeLeaf, error) {
	// var value merkletree.Value

	// get claim value in bytes
	leafData, err := mt.GetDataByIndex(&hi)
	if err != nil {
		return ProofTreeLeaf{}, err
	}

	claimType, _ := core.GetClaimTypeVersionFromData(leafData)
	nextVersion, err := GetNextVersion(mt, &hi)
	if err != nil {
		return ProofTreeLeaf{}, err
	}

	core.SetClaimTypeVersionInData(leafData, claimType, nextVersion)

	entry := merkletree.Entry{
		Data: *leafData,
	}
	mp, err := mt.GenerateProof(entry.HIndex(), nil)
	if err != nil {
		return ProofTreeLeaf{}, err
	}
	leafBytes := entry.Bytes()
	nonRevocationProof := ProofTreeLeaf{
		Leaf:  leafBytes[:],
		Proof: mp.Bytes(),
		Root:  *mt.RootKey(),
	}
	return nonRevocationProof, nil
}

// GetNextVersion returns the next version of a claim, given a Hash(index)
func GetNextVersion(mt *merkletree.MerkleTree, hi *merkletree.Hash) (uint32, error) {
	var claimType core.ClaimType
	var version uint32

	// loop until we find a nextversion that don't exist
	for {
		leafData, err := mt.GetDataByIndex(hi)
		if err == merkletree.ErrEntryIndexNotFound {
			return version, nil
		} else if err != nil {
			return 0, err
		}
		claimType, version = core.GetClaimTypeVersionFromData(leafData)
		version++

		core.SetClaimTypeVersionInData(leafData, claimType, version)

		entry := merkletree.Entry{
			Data: *leafData,
		}
		hi = entry.HIndex()
	}
}

// NewMerkleTreeUser creates a new user merkle tree by using an storage with
// the user addres prefix.
func NewMerkleTreeUser(id core.ID, storage db.Storage, levels int) (*merkletree.MerkleTree, error) {
	stoUserId := storage.WithPrefix(id.Bytes())
	if userMT, err := merkletree.NewMerkleTree(stoUserId, levels); err != nil {
		return nil, err
	} else {
		return userMT, nil
	}
}
