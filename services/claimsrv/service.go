package claimsrv

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/rootsrv"
	"github.com/iden3/go-iden3/services/signsrv"
	"github.com/iden3/go-iden3/utils"
)

type Service interface {
	AddAssignNameClaim(assignNameClaim core.AssignNameClaim) error
	ResolvAssignNameClaim(nameid, namespace string) (core.AssignNameClaim, error)
	AddAuthorizeKSignClaim(ethID common.Address, authorizeKSignClaimMsg AuthorizeKSignClaimMsg) error
	AddUserIDClaim(namespace string, ethID common.Address, claimValueMsg ClaimValueMsg) error
	GetIDRoot(ethID common.Address) (merkletree.Hash, []byte, error)
	GetClaimByHi(namespace string, ethID common.Address, hi merkletree.Hash) (ProofOfClaim, error)
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

// AddAssignNameClaim adds AssignNameClaim into the merkletree, updates the root in the smart contract, and returns the merkle proof of the claim in the merkletree
func (cs *ServiceImpl) AddAssignNameClaim(assignNameClaim core.AssignNameClaim) error {
	// get next version of the claim
	version, err := getNextVersion(cs.mt, assignNameClaim.Hi())
	if err != nil {
		return err
	}
	assignNameClaim.BaseIndex.Version = version

	// add AssignNameClaim to the Relay's merkletree
	// err := mt.Add(assignNameClaimMsg.AssignNameClaim)
	err = cs.mt.Add(assignNameClaim)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// update relay's root in smart contract
	cs.rootsrv.SetRoot(cs.mt.Root())

	return nil
}

// ResolvAssignNameClaim returns the AssignNameClaim from the merkletree, given a nameid and a namespace
func (cs *ServiceImpl) ResolvAssignNameClaim(nameid, namespace string) (core.AssignNameClaim, error) {
	// get name and domain
	s := strings.Split(nameid, "@")
	if len(s) != 2 {
		return core.AssignNameClaim{}, fmt.Errorf("Invalid nameid %v", s)
	}
	name := s[0]
	domain := s[1]

	// build the AssignNameClaim Partial with the given data of the Index
	nameHash := merkletree.HashBytes([]byte(name))
	domainHash := merkletree.HashBytes([]byte(domain))
	claimPartial := core.NewAssignNameClaim(namespace, nameHash, domainHash, common.Address{})
	version, err := getNextVersion(cs.mt, claimPartial.Hi())
	if err != nil {
		return core.AssignNameClaim{}, err
	}
	claimPartial.BaseIndex.Version = version - 1

	// get the complete AssignNameClaim in that merkle tree position
	claimInPosBytes, err := cs.mt.GetValueInPos(claimPartial.Hi())
	if err != nil {
		return core.AssignNameClaim{}, err
	}
	if bytes.Equal(claimInPosBytes, merkletree.EmptyNodeValue[:]) {
		return core.AssignNameClaim{}, errors.New("not found")
	}
	assignNameClaim, err := core.ParseAssignNameClaimBytes(claimInPosBytes)
	return assignNameClaim, err
}

// AddAuthorizeKSignClaim adds AuthorizeKSignClaim into the ID's merkletree, and adds the ID's merkle root into the Relay's merkletree inside a SetRootClaim. Returns the merkle proof of both Claims
func (cs *ServiceImpl) AddAuthorizeKSignClaim(ethID common.Address, authorizeKSignClaimMsg AuthorizeKSignClaimMsg) error {
	// TODO this function is almost the same than AddUserIDClaim() function. Maybe delete this one and keep the generic one
	// verify signature of the AuthorizeKSignClaim
	signature, err := common3.HexToBytes(authorizeKSignClaimMsg.Signature)
	if err != nil {
		return err
	}
	msgHash := utils.EthHash(authorizeKSignClaimMsg.AuthorizeKSignClaim.Bytes())
	signature[64] -= 27
	// if !utils.VerifySig(ethID, signature, msgHash[:]) {
	if !utils.VerifySig(authorizeKSignClaimMsg.KSign, signature, msgHash[:]) {
		return errors.New("signature can not be verified")
	}

	// get the user's id storage, using the user id prefix (the idaddress itself)
	stoUserID := cs.mt.Storage().WithPrefix(ethID.Bytes())

	// open the MerkleTree of the user
	userMT, err := merkletree.New(stoUserID, 140)
	if err != nil {
		return err
	}

	// add AuthorizeKSignClaim into the User's ID Merkle Tree
	err = userMT.Add(authorizeKSignClaimMsg.AuthorizeKSignClaim)
	if err != nil {
		return err
	}

	// create new SetRootClaim
	setRootClaim := core.NewSetRootClaim("iden3.io", ethID, userMT.Root())

	// get next version of the claim
	version, err := getNextVersion(cs.mt, setRootClaim.Hi())
	if err != nil {
		return err
	}
	setRootClaim.BaseIndex.Version = version

	// add User's ID Merkle Root into the Relay's Merkle Tree
	err = cs.mt.Add(setRootClaim)
	if err != nil {
		return err
	}

	// update Relay's Root in the Smart Contract
	cs.rootsrv.SetRoot(cs.mt.Root())

	return nil
}

// AddUserIDClaim adds a claim into the ID's merkle tree, and with the ID's root, creates a new SetRootClaim and adds it to the Relay's merkletree
func (cs *ServiceImpl) AddUserIDClaim(namespace string, ethID common.Address, claimValueMsg ClaimValueMsg) error {
	// verify proof of KSign
	proofOfKSign, err := claimValueMsg.ProofOfKSignHex.Unhex()
	if err != nil {
		return err
	}
	if !CheckProofOfClaim(proofOfKSign, 140) {
		return errors.New("ProofOfClaim can not be verified")
	}

	// verify signature with KSign
	signature, err := common3.HexToBytes(claimValueMsg.Signature)
	if err != nil {
		return err
	}

	msgHash := utils.EthHash(claimValueMsg.ClaimValue.Bytes())
	signature[64] -= 27
	ksign := claimValueMsg.KSign
	// if !utils.VerifySig(ethID, signature, msgHash[:]) {
	if !utils.VerifySig(ksign, signature, msgHash[:]) {
		return errors.New("signature can not be verified")
	}

	// get the user's id storage, using the user id prefix (the idaddress itself)
	stoUserID := cs.mt.Storage().WithPrefix(ethID.Bytes())

	// open the MerkleTree of the user
	userMT, err := merkletree.New(stoUserID, 140)
	if err != nil {
		return err
	}

	// add claim in User ID Merkle Tree
	err = userMT.Add(claimValueMsg.ClaimValue)
	if err != nil {
		return err
	}

	// setRootClaim of the user in the Relay Merkle Tree
	// create new SetRootClaim
	setRootClaim := core.NewSetRootClaim(namespace, ethID, userMT.Root())
	version, err := getNextVersion(cs.mt, setRootClaim.Hi())
	if err != nil {
		return err
	}
	setRootClaim.BaseIndex.Version = version

	// add User's ID Merkle Root into the Relay's Merkle Tree
	err = cs.mt.Add(setRootClaim)
	if err != nil {
		return err
	}

	// update Relay Root in Smart Contract
	cs.rootsrv.SetRoot(cs.mt.Root())

	return nil
}

// GetIDRoot returns the root of an ID tree, and the proof of that Root ID tree in the Relay Merkle Tree
func (cs *ServiceImpl) GetIDRoot(ethID common.Address) (merkletree.Hash, []byte, error) {
	// get the user's id storage, using the user id prefix (the idaddress itself)
	stoUserID := cs.mt.Storage().WithPrefix(ethID.Bytes())
	// open the MerkleTree of the user
	userMT, err := merkletree.New(stoUserID, 140)
	if err != nil {
		return merkletree.Hash{}, []byte{}, err
	}
	// build SetRootClaim of the user id
	setRootClaim := core.NewSetRootClaim("iden3.io", ethID, userMT.Root())
	version, err := getNextVersion(cs.mt, setRootClaim.Hi())
	if err != nil {
		return merkletree.Hash{}, []byte{}, err
	}
	setRootClaim.BaseIndex.Version = version - 1
	// get proof of SetRootProof in the Relay tree
	idRootProof, err := cs.mt.GenerateProof(setRootClaim.Hi())
	if err != nil {
		return merkletree.Hash{}, []byte{}, err
	}
	return userMT.Root(), idRootProof, nil
}

// GetClaimByHi given a Hash(index) (Hi) and an ID, returns the Claim in that Hi position inside the ID's merkletree, and the SetRootClaim with the ID's root in the Relay's merkletree
func (cs *ServiceImpl) GetClaimByHi(namespace string, ethID common.Address, hi merkletree.Hash) (ProofOfClaim, error) {
	// get the user's id storage, using the user id prefix (the idaddress itself)
	stoUserID := cs.mt.Storage().WithPrefix(ethID.Bytes())

	// open the MerkleTree of the user
	userMT, err := merkletree.New(stoUserID, 140)
	if err != nil {
		return ProofOfClaim{}, err
	}

	// get the value in the hi position
	valueBytes, err := userMT.GetValueInPos(hi)
	if err != nil {
		return ProofOfClaim{}, err
	}
	value, err := core.ParseValueFromBytes(valueBytes)
	if err != nil {
		return ProofOfClaim{}, err
	}

	// get the proof of the value in the User ID Tree
	idProof, err := userMT.GenerateProof(merkletree.HashBytes(value.Bytes()[:value.IndexLength()]))
	if err != nil {
		return ProofOfClaim{}, err
	}

	claimProof := ProofOfTreeLeaf{
		Leaf:  valueBytes,
		Hi:    merkletree.HashBytes(value.Bytes()[:value.IndexLength()]),
		Proof: idProof,
		Root:  userMT.Root(),
	}

	// build SetRootClaim
	setRootClaim := core.NewSetRootClaim(namespace, ethID, userMT.Root())
	version, err := getNextVersion(cs.mt, setRootClaim.Hi())
	if err != nil {
		return ProofOfClaim{}, err
	}
	setRootClaim.BaseIndex.Version = version - 1

	// get the proof of the SetRootClaim in the Relay Tree
	relayProof, err := cs.mt.GenerateProof(setRootClaim.Hi())
	if err != nil {
		return ProofOfClaim{}, err
	}
	setRootClaimProof := ProofOfTreeLeaf{
		Leaf:  setRootClaim.Bytes(),
		Hi:    setRootClaim.Hi(),
		Proof: relayProof,
		Root:  cs.mt.Root(),
	}

	// get non revocation proofs of the claim
	claimNonRevocationProof, err := getNonRevocationProof(userMT, merkletree.HashBytes(value.Bytes()[:value.IndexLength()]))
	if err != nil {
		return ProofOfClaim{}, err
	}
	setRootClaimNonRevocationProof, err := getNonRevocationProof(cs.mt, setRootClaim.Hi())
	if err != nil {
		return ProofOfClaim{}, err
	}

	// sign root
	sig, err := cs.signer.SignHash(setRootClaimProof.Root)
	if err != nil {
		return ProofOfClaim{}, err
	}
	proofOfClaim := ProofOfClaim{
		claimProof,
		setRootClaimProof,
		claimNonRevocationProof,
		setRootClaimNonRevocationProof,
		sig,
	}
	return proofOfClaim, nil
}
func (cs *ServiceImpl) MT() *merkletree.MerkleTree {
	return cs.mt
}

// GetNonRevocationProof returns the next version Hi (that don't exist in the tree, it's value is Empty) with merkleproof and root
func getNonRevocationProof(mt *merkletree.MerkleTree, hi merkletree.Hash) (ProofOfTreeLeaf, error) {
	var value merkletree.Value

	// get claim value in bytes
	b, err := mt.GetValueInPos(hi)
	if err != nil {
		return ProofOfTreeLeaf{}, err
	}
	if bytes.Equal(b, merkletree.EmptyNodeValue[:]) {
		return ProofOfTreeLeaf{}, errors.New("Hi not found in the merkle tree")
	}

	nextVersion, err := getNextVersion(mt, hi)
	if err != nil {
		return ProofOfTreeLeaf{}, err
	}

	// get claim with version+1 from the merkletree
	nextVersionBytes, err := core.Uint32ToEthBytes(nextVersion)
	if err != nil {
		return ProofOfTreeLeaf{}, err
	}
	copy(b[60:64], nextVersionBytes)

	value, err = core.ParseValueFromBytes(b)
	if err != nil {
		return ProofOfTreeLeaf{}, err
	}
	mp, err := mt.GenerateProof(merkletree.HashBytes(value.Bytes()[:value.IndexLength()]))
	if err != nil {
		return ProofOfTreeLeaf{}, err
	}
	nonRevocationProof := ProofOfTreeLeaf{
		Leaf:  merkletree.EmptyNodeValue[:],
		Hi:    merkletree.HashBytes(value.Bytes()[:value.IndexLength()]),
		Proof: mp,
		Root:  mt.Root(),
	}
	return nonRevocationProof, nil
}

// GetNextVersion returns the next version of a claim, given a Hash(index)
func getNextVersion(mt *merkletree.MerkleTree, hi merkletree.Hash) (uint32, error) {
	var version uint32

	// loop until found a nextversion that don't exist
	for {
		// merkletree.GetValueInPos(hi)
		b, err := mt.GetValueInPos(hi)
		if err != nil {
			return 0, err
		}
		// if value not exist, return version 0
		if bytes.Equal(merkletree.EmptyNodeValue[:], b) {
			break
		}
		// get version bytes
		versionBytes := b[60:64]
		version = core.EthBytesToUint32(versionBytes)
		version++

		// get claim with version+1 from the merkletree
		versionBytes, err = core.Uint32ToEthBytes(version)
		if err != nil {
			return 0, err
		}
		copy(b[60:64], versionBytes)
		value, err := core.ParseValueFromBytes(b)
		if err != nil {
			return 0, err
		}
		hi = merkletree.HashBytes(value.Bytes()[:value.IndexLength()])
	}
	// return version
	return version, nil
}
