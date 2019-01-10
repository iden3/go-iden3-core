package claimsrv

import (
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/rootsrv"
	"github.com/iden3/go-iden3/services/signsrv"
	"github.com/iden3/go-iden3/utils"
	log "github.com/sirupsen/logrus"
)

var ErrNotFound = errors.New("value not found")

type Service interface {
	CommitNewIDRoot(idaddr common.Address, kSign common.Address, root merkletree.Hash, timestamp uint64, signature []byte) (*core.ClaimSetRootKey, error)
	AddClaimAssignName(claimAssignName core.ClaimAssignName) error
	AddClaimAuthorizeKSign(ethID common.Address, claimAuthorizeKSignMsg ClaimAuthorizeKSignMsg) error
	AddClaimAuthorizeKSignFirst(ethID common.Address, claimAuthorizeKSign core.ClaimAuthorizeKSign) error
	AddUserIDClaim(ethID common.Address, claimValueMsg ClaimValueMsg) error
	AddDirectClaim(claim core.ClaimBasic) error
	GetIDRoot(ethID common.Address) (merkletree.Hash, []byte, error)
	GetClaimByHi(ethID common.Address, hi merkletree.Hash) (ProofOfClaim, error)
	GetRelayClaimByHi(hi merkletree.Hash) (ProofOfRelayClaim, error)
	MT() *merkletree.MerkleTree
}

type ServiceImpl struct {
	mt      *merkletree.MerkleTree
	rootsrv rootsrv.Service
	signer  signsrv.Service
}

func New(mt *merkletree.MerkleTree, rootsrv rootsrv.Service, signer signsrv.Service) *ServiceImpl {
	return &ServiceImpl{mt, rootsrv, signer}
}

// SetNewIDRoot checks that the data is valid and performs a claim in the Relay merkletree setting the new Root of the emmiting ID
func (cs *ServiceImpl) CommitNewIDRoot(idaddr common.Address, kSign common.Address, root merkletree.Hash, timestamp uint64, signature []byte) (*core.ClaimSetRootKey, error) {
	// get the user's id storage, using the user id prefix (the idaddress itself)
	stoUserID := cs.mt.Storage().WithPrefix(idaddr.Bytes())

	// open the MerkleTree of the user
	userMT, err := merkletree.NewMerkleTree(stoUserID, 140)
	if err != nil {
		return &core.ClaimSetRootKey{}, err
	}

	// verify that the KSign is authorized
	if !CheckKSignInIDdb(userMT, kSign) {
		return &core.ClaimSetRootKey{}, errors.New("can not verify the KSign")
	}
	// in the future the user merkletree will be in the client side, and this step will be a check of the ProofOfKSign

	// check data timestamp
	verified := utils.VerifyTimestamp(timestamp, 30000) //needs to be from last 30 seconds
	if !verified {
		return &core.ClaimSetRootKey{}, errors.New("timestamp too old")
	}
	// check signature with idaddr
	// whee data signed is idaddr+root+timestamp
	timestampBytes, err := utils.Uint64ToEthBytes(timestamp)
	if err != nil {
		return &core.ClaimSetRootKey{}, err
	}
	// signature of idaddr+root+timestamp, only valid if is from last X seconds
	var msg []byte
	msg = append(msg, idaddr.Bytes()...)
	msg = append(msg, root.Bytes()...)
	msg = append(msg, timestampBytes...)
	msgHash := utils.EthHash(msg)
	signature[64] -= 27
	if !utils.VerifySig(kSign, signature, msgHash[:]) {
		return &core.ClaimSetRootKey{}, errors.New("signature can not be verified")
	}

	// claimSetRootKey of the user in the Relay Merkle Tree
	// create new ClaimSetRootKey
	claimSetRootKey := core.NewClaimSetRootKey(idaddr, root)
	// entry := claimSetRootKey.Entry()
	// version, err := GetNextVersion(cs.mt, entry.HIndex())
	version, err := GetNextVersion(cs.mt, claimSetRootKey.Entry().HIndex())
	if err != nil {
		return &core.ClaimSetRootKey{}, err
	}
	claimSetRootKey.Version = version

	// add User's ID Merkle Root into the Relay's Merkle Tree
	e := claimSetRootKey.Entry()
	err = cs.mt.Add(e)
	if err != nil {
		return &core.ClaimSetRootKey{}, err
	}

	// update Relay Root in Smart Contract
	cs.rootsrv.SetRoot(*cs.mt.RootKey())

	return claimSetRootKey, nil
}

// AddClaimAssignName adds ClaimAssignName into the merkletree, updates the root in the smart contract, and returns the merkle proof of the claim in the merkletree
func (cs *ServiceImpl) AddClaimAssignName(claimAssignName core.ClaimAssignName) error {
	// get next version of the claim
	entry := claimAssignName.Entry()
	version, err := GetNextVersion(cs.mt, entry.HIndex())
	if err != nil {
		return err
	}
	claimAssignName.Version = version

	// add ClaimAssignName to the Relay's merkletree
	e := claimAssignName.Entry()
	err = cs.mt.Add(e)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// update relay's root in smart contract
	cs.rootsrv.SetRoot(*cs.mt.RootKey())

	return nil
}

// AddClaimAuthorizeKSign adds ClaimAuthorizeKSign into the ID's merkletree, and adds the ID's merkle root into the Relay's merkletree inside a ClaimSetRootKey. Returns the merkle proof of both Claims
func (cs *ServiceImpl) AddClaimAuthorizeKSign(ethID common.Address, claimAuthorizeKSignMsg ClaimAuthorizeKSignMsg) error {

	// get the user's id storage, using the user id prefix (the idaddress itself)
	stoUserID := cs.mt.Storage().WithPrefix(ethID.Bytes())

	// open the MerkleTree of the user
	userMT, err := merkletree.NewMerkleTree(stoUserID, 140)
	if err != nil {
		return err
	}

	// verify that the KSign is authorized
	if !CheckKSignInIDdb(userMT, claimAuthorizeKSignMsg.KSign) {
		return errors.New("can not verify the KSign")
	}

	// verify signature of the ClaimAuthorizeKSign
	signature, err := common3.HexToBytes(claimAuthorizeKSignMsg.Signature)
	if err != nil {
		return err
	}
	msgHash := utils.EthHash(claimAuthorizeKSignMsg.ClaimAuthorizeKSign.Entry().Bytes())
	signature[64] -= 27
	if !utils.VerifySig(claimAuthorizeKSignMsg.KSign, signature, msgHash[:]) {
		return errors.New("signature can not be verified")
	}

	// add ClaimAuthorizeKSign into the User's ID Merkle Tree
	err = userMT.Add(claimAuthorizeKSignMsg.ClaimAuthorizeKSign.Entry())
	if err != nil {
		return err
	}

	// create new ClaimSetRootKey
	claimSetRootKey := core.NewClaimSetRootKey(ethID, *userMT.RootKey())

	// get next version of the claim
	version, err := GetNextVersion(cs.mt, claimSetRootKey.Entry().HIndex())
	if err != nil {
		return err
	}
	claimSetRootKey.Version = version

	// add User's ID Merkle Root into the Relay's Merkle Tree
	err = cs.mt.Add(claimSetRootKey.Entry())
	if err != nil {
		return err
	}

	// update Relay's Root in the Smart Contract
	cs.rootsrv.SetRoot(*cs.mt.RootKey())

	return nil
}

// AddClaimAuthorizeKSign adds ClaimAuthorizeKSign into the ID's merkletree, and adds the ID's merkle root into the Relay's merkletree inside a ClaimSetRootKey. Returns the merkle proof of both Claims
func (cs *ServiceImpl) AddClaimAuthorizeKSignFirst(ethID common.Address, claimAuthorizeKSign core.ClaimAuthorizeKSign) error {

	// get the user's id storage, using the user id prefix (the idaddress itself)
	stoUserID := cs.mt.Storage().WithPrefix(ethID.Bytes())

	// open the MerkleTree of the user
	userMT, err := merkletree.NewMerkleTree(stoUserID, 140)
	if err != nil {
		return err
	}

	// add ClaimAuthorizeKSign into the User's ID Merkle Tree
	err = userMT.Add(claimAuthorizeKSign.Entry())
	if err != nil {
		return err
	}

	// create new ClaimSetRootKey
	claimSetRootKey := core.NewClaimSetRootKey(ethID, *userMT.RootKey())

	// get next version of the claim
	version, err := GetNextVersion(cs.mt, claimSetRootKey.Entry().HIndex())
	if err != nil {
		return err
	}
	claimSetRootKey.Version = version

	// add User's ID Merkle Root into the Relay's Merkle Tree
	err = cs.mt.Add(claimSetRootKey.Entry())
	if err != nil {
		return err
	}

	// update Relay's Root in the Smart Contract
	cs.rootsrv.SetRoot(*cs.mt.RootKey())

	return nil
}

// AddUserIDClaim adds a claim into the ID's merkle tree, and with the ID's root, creates a new ClaimSetRootKey and adds it to the Relay's merkletree
func (cs *ServiceImpl) AddUserIDClaim(ethID common.Address, claimValueMsg ClaimValueMsg) error {
	// get the user's id storage, using the user id prefix (the idaddress itself)
	stoUserID := cs.mt.Storage().WithPrefix(ethID.Bytes())

	// open the MerkleTree of the user
	userMT, err := merkletree.NewMerkleTree(stoUserID, 140)
	if err != nil {
		return err
	}

	// verify that the KSign is authorized
	if !CheckKSignInIDdb(userMT, claimValueMsg.KSign) {
		return errors.New("can not verify the KSign")
	}

	// verify signature with KSign
	signature, err := common3.HexToBytes(claimValueMsg.Signature)
	if err != nil {
		return err
	}

	msgHash := utils.EthHash(claimValueMsg.ClaimValue.Bytes())
	signature[64] -= 27
	ksign := claimValueMsg.KSign
	if !utils.VerifySig(ksign, signature, msgHash[:]) {
		return errors.New("signature can not be verified")
	}

	// add claim in User ID Merkle Tree
	err = userMT.Add(&claimValueMsg.ClaimValue)
	if err != nil {
		return err
	}

	// claimSetRootKey of the user in the Relay Merkle Tree
	// create new ClaimSetRootKey
	claimSetRootKey := core.NewClaimSetRootKey(ethID, *userMT.RootKey())
	version, err := GetNextVersion(cs.mt, claimSetRootKey.Entry().HIndex())
	if err != nil {
		return err
	}
	claimSetRootKey.Version = version

	// add User's ID Merkle Root into the Relay's Merkle Tree
	err = cs.mt.Add(claimSetRootKey.Entry())
	if err != nil {
		return err
	}

	// update Relay Root in Smart Contract
	cs.rootsrv.SetRoot(*cs.mt.RootKey())

	return nil
}

// AddDirectClaim adds a claim directly to the Relay merkletree
func (cs *ServiceImpl) AddDirectClaim(claim core.ClaimBasic) error {
	err := cs.mt.Add(claim.Entry())
	if err != nil {
		return err
	}
	cs.rootsrv.SetRoot(*cs.mt.RootKey())
	return nil
}

// GetIDRoot returns the root of an ID tree, and the proof of that Root ID tree in the Relay Merkle Tree
func (cs *ServiceImpl) GetIDRoot(ethID common.Address) (merkletree.Hash, []byte, error) {
	// get the user's id storage, using the user id prefix (the idaddress itself)
	stoUserID := cs.mt.Storage().WithPrefix(ethID.Bytes())

	// open the MerkleTree of the user
	userMT, err := merkletree.NewMerkleTree(stoUserID, 140)
	if err != nil {
		return merkletree.Hash{}, []byte{}, err
	}

	// build ClaimSetRootKey of the user id
	claimSetRootKey := core.NewClaimSetRootKey(ethID, *userMT.RootKey())
	version, err := GetNextVersion(cs.mt, claimSetRootKey.Entry().HIndex())
	if err != nil {
		return merkletree.Hash{}, []byte{}, err
	}
	claimSetRootKey.Version = version - 1

	// get proof of SetRootProof in the Relay tree
	idRootProof, err := cs.mt.GenerateProof(claimSetRootKey.Entry().HIndex())
	if err != nil {
		return merkletree.Hash{}, []byte{}, err
	}
	return *userMT.RootKey(), idRootProof.Bytes(), nil
}

// GetClaimByHi given a Hash(index) (Hi) and an ID, returns the Claim in that Hi position inside the ID's merkletree, and the ClaimSetRootKey with the ID's root in the Relay's merkletree
func (cs *ServiceImpl) GetClaimByHi(ethID common.Address, hi merkletree.Hash) (ProofOfClaim, error) {
	// get the user's id storage, using the user id prefix (the idaddress itself)
	stoUserID := cs.mt.Storage().WithPrefix(ethID.Bytes())

	// open the MerkleTree of the user
	userMT, err := merkletree.NewMerkleTree(stoUserID, 140)
	if err != nil {
		return ProofOfClaim{}, err
	}

	// get the value in the hi position
	// valueBytes, err := userMT.GetValueInPos(hi)
	leafData, err := userMT.GetDataByIndex(&hi)
	if err != nil {
		return ProofOfClaim{}, err
	}
	// if bytes.Equal(valueBytes, merkletree.EmptyNodeValue[:]) {
	//         return ProofOfClaim{}, ErrNotFound
	// }

	// value, err := core.ParseValueFromBytes(valueBytes)
	// if err != nil {
	//         return ProofOfClaim{}, err
	// }

	// get the proof of the value in the User ID Tree
	// idProof, err := userMT.GenerateProof(merkletree.HashBytes(value.Bytes()[:value.IndexLength()]))
	idProof, err := userMT.GenerateProof(&hi)
	if err != nil {
		return ProofOfClaim{}, err
	}

	leafBytes := leafData.Bytes()
	claimProof := ProofOfTreeLeaf{
		Leaf:  leafBytes[:],
		Proof: idProof.Bytes(),
		Root:  *userMT.RootKey(),
	}

	// build ClaimSetRootKey
	claimSetRootKey := core.NewClaimSetRootKey(ethID, *userMT.RootKey())
	version, err := GetNextVersion(cs.mt, claimSetRootKey.Entry().HIndex())
	if err != nil {
		return ProofOfClaim{}, err
	}
	claimSetRootKey.Version = version - 1

	// get the proof of the ClaimSetRootKey in the Relay Tree
	relayProof, err := cs.mt.GenerateProof(claimSetRootKey.Entry().HIndex())
	if err != nil {
		return ProofOfClaim{}, err
	}
	claimSetRootKeyProof := ProofOfTreeLeaf{
		Leaf:  claimSetRootKey.Entry().Bytes(),
		Proof: relayProof.Bytes(),
		Root:  *cs.mt.RootKey(),
	}

	// get non revocation proofs of the claim
	claimNonRevocationProof, err := getNonRevocationProof(userMT, hi)
	if err != nil {
		return ProofOfClaim{}, err
	}
	claimSetRootKeyNonRevocationProof, err := getNonRevocationProof(cs.mt, *claimSetRootKey.Entry().HIndex())
	if err != nil {
		return ProofOfClaim{}, err
	}

	// sign root + date
	dateUint64 := uint64(time.Now().Unix())
	dateBytes, err := utils.Uint64ToEthBytes(dateUint64)
	if err != nil {
		return ProofOfClaim{}, err
	}
	rootdate := claimSetRootKeyProof.Root[:]
	rootdate = append(rootdate, dateBytes...)
	rootdateHash := merkletree.HashBytes(rootdate)
	sig, err := cs.signer.SignHash(rootdateHash)
	// sig[64] += 27
	if err != nil {
		return ProofOfClaim{}, err
	}

	proofOfClaim := ProofOfClaim{
		claimProof,
		claimSetRootKeyProof,
		claimNonRevocationProof,
		claimSetRootKeyNonRevocationProof,
		dateUint64,
		sig,
	}
	return proofOfClaim, nil
}

// GetRelayClaimByHi given a Hash(index) (Hi), returns the Claim in that Hi position inside the Relay merkletree, and it's proof of non revocated
func (cs *ServiceImpl) GetRelayClaimByHi(hi merkletree.Hash) (ProofOfRelayClaim, error) {
	// get the value in the hi position
	leafData, err := cs.mt.GetDataByIndex(&hi)
	if err != nil {
		return ProofOfRelayClaim{}, err
	}

	// get the proof of the ClaimSetRootKey in the Relay Tree
	relayProof, err := cs.mt.GenerateProof(&hi)
	if err != nil {
		return ProofOfRelayClaim{}, err
	}

	leafBytes := leafData.Bytes()
	claimProof := ProofOfTreeLeaf{
		Leaf:  leafBytes[:],
		Proof: relayProof.Bytes(),
		Root:  *cs.mt.RootKey(),
	}

	claimNonRevocationProof, err := getNonRevocationProof(cs.mt, hi)
	if err != nil {
		return ProofOfRelayClaim{}, err
	}

	// sign root + date
	dateUint64 := uint64(time.Now().Unix())
	dateBytes, err := utils.Uint64ToEthBytes(dateUint64)
	if err != nil {
		return ProofOfRelayClaim{}, err
	}
	rootdate := claimProof.Root[:]
	rootdate = append(rootdate, dateBytes...)
	rootdateHash := merkletree.HashBytes(rootdate)
	sig, err := cs.signer.SignHash(rootdateHash)
	if err != nil {
		return ProofOfRelayClaim{}, err
	}

	proofOfRelayClaim := ProofOfRelayClaim{
		claimProof,
		claimNonRevocationProof,
		dateUint64,
		sig,
	}
	return proofOfRelayClaim, nil
}

func (cs *ServiceImpl) MT() *merkletree.MerkleTree {
	return cs.mt
}

// getNonRevocationProof returns the next version Hi (that don't exist in the tree, it's value is Empty) with merkleproof and root
func getNonRevocationProof(mt *merkletree.MerkleTree, hi merkletree.Hash) (ProofOfTreeLeaf, error) {
	// var value merkletree.Value

	// get claim value in bytes
	leafData, err := mt.GetDataByIndex(&hi)
	if err != nil {
		return ProofOfTreeLeaf{}, err
	}

	claimType, _ := core.GetClaimTypeVersionFromData(leafData)
	nextVersion, err := GetNextVersion(mt, &hi)
	if err != nil {
		return ProofOfTreeLeaf{}, err
	}

	core.SetClaimTypeVersionInData(leafData, claimType, nextVersion)

	entry := merkletree.Entry{
		Data: *leafData,
	}
	mp, err := mt.GenerateProof(entry.HIndex())
	if err != nil {
		return ProofOfTreeLeaf{}, err
	}
	leafBytes := entry.Bytes()
	nonRevocationProof := ProofOfTreeLeaf{
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

	// loop until found a nextversion that don't exist
	for {
		leafData, err := mt.GetDataByIndex(hi)
		if err == merkletree.ErrEntryIndexNotFound {
			return version, nil
		}
		if err != nil {
			return 0, err
		}
		// if value not exist, return version 0
		// if bytes.Equal(merkletree.EmptyNodeValue[:], b) {
		//         break
		// }
		// get version bytes
		// versionBytes := b[60:64]
		// version = utils.EthBytesToUint32(versionBytes)
		claimType, version = core.GetClaimTypeVersionFromData(leafData)
		version++

		//
		// // get claim with version+1 from the merkletree
		// versionBytes, err = utils.Uint32ToEthBytes(version)
		// if err != nil {
		//         return 0, err
		// }
		// copy(b[60:64], versionBytes)
		// value, err := core.ParseValueFromBytes(b)
		// if err != nil {
		//         return 0, err
		// }

		core.SetClaimTypeVersionInData(leafData, claimType, version)

		entry := merkletree.Entry{
			Data: *leafData,
		}
		hi = entry.HIndex()
		// hi = merkletree.HashBytes(value.Bytes()[:value.IndexLength()])
	}
}
