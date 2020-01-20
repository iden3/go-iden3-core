package claimsrv

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	// "time"

	// "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/core"
	"github.com/iden3/go-iden3-core/db"
	babykeystore "github.com/iden3/go-iden3-core/keystore"
	"github.com/iden3/go-iden3-core/merkletree"
	rootsrvmock "github.com/iden3/go-iden3-core/services/rootsrv/mock"
	"github.com/iden3/go-iden3-core/services/signsrv"
	"github.com/iden3/go-iden3-crypto/babyjub"

	// "github.com/iden3/go-iden3-core/utils"
	"github.com/ipfsconsortium/go-ipfsc/config"
	// "github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var debug = true

var service *Service
var rootService *rootsrvmock.RootServiceMock
var mt *merkletree.MerkleTree
var c config.Config

const relayIdHex = "113kyY52PSBr9oUqosmYkCavjjrQFuiuAw47FpZeUf"
const relaySkHex = "4be5471a938bdf3606888472878baace4a6a64e14a153adf9a1333969e4e573c"

var relayID core.ID
var keyStore *babykeystore.KeyStore
var relaySk babyjub.PrivateKey
var relayPkComp *babyjub.PublicKeyComp
var relayPk *babyjub.PublicKey

// type SignServiceMock struct {
// 	mock.Mock
// }
//
// func (m *SignServiceMock) SignEthMsg(msg []byte) (*utils.SignatureEthMsg, error) {
// 	h := utils.EthHash(msg)
// 	sig, err := crypto.Sign(h[:], relaySecKey)
// 	if err != nil {
// 		return nil, err
// 	}
// 	sig[64] += 27
// 	sigEthMsg := &utils.SignatureEthMsg{}
// 	copy(sigEthMsg[:], sig)
// 	return sigEthMsg, nil
// }
//
// func (m *SignServiceMock) SignEthMsgDate(msg []byte) (*utils.SignatureEthMsg, int64, error) {
// 	dateInt64 := time.Now().Unix()
// 	dateBytes := utils.Uint64ToEthBytes(uint64(dateInt64))
// 	sig, err := m.SignEthMsg(append(msg, dateBytes...))
// 	return sig, dateInt64, err
// }
//
// func (m *SignServiceMock) PublicKey() *ecdsa.PublicKey {
// 	return relayPubKey
// }

var rmDirs []string

func newTestingMerkle(numLevels int) (*merkletree.MerkleTree, error) {
	dir, err := ioutil.TempDir("", "db")
	rmDirs = append(rmDirs, dir)
	if err != nil {
		return &merkletree.MerkleTree{}, err
	}
	sto, err := db.NewLevelDbStorage(dir, false)
	if err != nil {
		return &merkletree.MerkleTree{}, err
	}

	mt, err := merkletree.NewMerkleTree(sto, numLevels)
	return mt, err
}

/*
func loadConfig() {
	c.Server.Port = "5000"
	c.Server.PrivK = "da7079f082a1ced80c5dee3bf00752fd67f75321a637e5d5073ce1489af062d8"
	c.Geth.URL = ""
	c.ContractsAddress.Identities = "0x101d2fa51f8259df207115af9eaa73f3f4e52e60"
	c.Domain = "iden3.io"
	c.Namespace = "iden3.io"
}
*/
func initializeEnvironment(t *testing.T) {

	// MerkleTree leveldb
	var err error
	mt, err = newTestingMerkle(140)
	if err != nil {
		t.Error(err)
	}

	pass := []byte("my passphrase")
	storage := babykeystore.MemStorage([]byte{})
	keyStore, err := babykeystore.NewKeyStore(&storage, babykeystore.LightKeyStoreParams)
	if err != nil {
		panic(err)
	}

	if _, err := hex.Decode(relaySk[:], []byte(relaySkHex)); err != nil {
		panic(err)
	}
	if relayPkComp, err = keyStore.ImportKey(relaySk, pass); err != nil {
		panic(err)
	}
	if err := keyStore.UnlockKey(relayPkComp, pass); err != nil {
		panic(err)
	}
	if relayPk, err = relayPkComp.Decompress(); err != nil {
		panic(err)
	}

	signSrv := signsrv.New(keyStore, *relayPk)

	relayID, err = core.IDFromString(relayIdHex)
	if err != nil {
		panic(err)
	}
	rootService = rootsrvmock.New()
	service = New(&relayID, mt, rootService, *signSrv)
}

func TestGetNextVersion(t *testing.T) {
	initializeEnvironment(t)

	indexData := []byte("c1")
	data := []byte{}
	var indexSlot [400 / 8]byte
	var dataSlot [496 / 8]byte
	// copy(indexSlot[:], indexData[:400/8])
	// copy(dataSlot[:], data[:496/8])
	copy(indexSlot[:], indexData[:])
	copy(dataSlot[:], data[:])
	claim := core.NewClaimBasic(indexSlot, dataSlot)

	version, err := GetNextVersion(mt, claim.Entry().HIndex())
	require.Nil(t, err)
	require.Equal(t, uint32(0), version)
	claim.Version = version
	mt.Add(claim.Entry())
	version, err = GetNextVersion(mt, claim.Entry().HIndex())
	require.Nil(t, err)
	require.Equal(t, uint32(1), version)

	claim.Version = version
	mt.Add(claim.Entry())
	version, err = GetNextVersion(mt, claim.Entry().HIndex())
	require.Nil(t, err)
	require.Equal(t, uint32(2), version)

	claim.Version = version
	mt.Add(claim.Entry())
	version, err = GetNextVersion(mt, claim.Entry().HIndex())
	require.Nil(t, err)
	require.Equal(t, uint32(3), version)

	claim.Version = version
	mt.Add(claim.Entry())
	version, err = GetNextVersion(mt, claim.Entry().HIndex())
	require.Nil(t, err)
	require.Equal(t, uint32(4), version)
}

func TestGetNonRevocationProof(t *testing.T) {
	initializeEnvironment(t)
	indexData := []byte("c1")
	data := []byte{}
	var indexSlot [400 / 8]byte
	var dataSlot [496 / 8]byte
	copy(indexSlot[:], indexData[:])
	copy(dataSlot[:], data[:])
	claim := core.NewClaimBasic(indexSlot, dataSlot)

	err := mt.Add(claim.Entry())
	require.Nil(t, err)
	version, err := GetNextVersion(mt, claim.Entry().HIndex())
	require.Nil(t, err)
	require.Equal(t, uint32(1), version)

	claimProof, err := getNonRevocationProof(mt, claim.Entry().HIndex())
	require.Nil(t, err)

	require.Equal(t,
		"0x03000000000000000000000000000000000000000000000000000000000000002d3258815847168fd948039e4c3295028c06300755badb5c7a7a19076e175328021a76d5f2cdcf354ab66eff7b4dee40f02501545def7bb66b3502ae68e1b781",
		common3.HexEncode(claimProof.Proof))
	require.Equal(t,
		"0x1ba11fb41fcdfff061db7bfb1d89a5e27eca53ea2b7c471bdf69360151dcf4f4",
		claimProof.Root.Hex())

	proof, err := merkletree.NewProofFromBytes(claimProof.Proof)
	require.Nil(t, err)

	// VerifyProof with HIndex and HValue from claimProof.Leaf
	var leafBytes [128]byte
	copy(leafBytes[:], claimProof.Leaf)
	dataNonE := merkletree.NewDataFromBytes(leafBytes)
	claimProofEntry := merkletree.Entry{
		Data: *dataNonE,
	}
	verified := merkletree.VerifyProof(&claimProof.Root, proof, claimProofEntry.HIndex(), claimProofEntry.HValue())
	require.True(t, verified)
}

func TestGetClaimProof(t *testing.T) {
	initializeEnvironment(t)

	id, err := core.IDFromString("11AVZrKNJVqDJoyKrdyaAgEynyBEjksV5z2NjZoPxf")
	require.Nil(t, err)

	// Basic Claim
	indexData := []byte("index01")
	data := []byte("data01")
	var indexSlot [400 / 8]byte
	var dataSlot [496 / 8]byte
	copy(indexSlot[:], indexData[:])
	copy(dataSlot[:], data[:])
	claimBasic := core.NewClaimBasic(indexSlot, dataSlot)

	// KSign Claim
	sk, err := crypto.HexToECDSA("0b8bdda435a144fc12764c0afe4ac9e2c4d544bf5692d2a6353ec2075dc1fcb4")
	if err != nil {
		panic(err)
	}
	pk := sk.Public().(*ecdsa.PublicKey)
	claimAuthKSign := core.NewClaimAuthorizeKSignSecp256k1(pk)

	var kSignSk babyjub.PrivateKey
	if _, err := hex.Decode(kSignSk[:], []byte("9b3260823e7b07dd26ef357ccfed23c10bcef1c85940baa3d02bbf29461bbbbe")); err != nil {
		panic(err)
	}
	kSignPk := kSignSk.Public()
	claimAuthKSignBabyJub := core.NewClaimAuthorizeKSignBabyJub(kSignPk)

	// open the MerkleTree of the user
	userMT, err := newTestingMerkle(140)
	require.Nil(t, err)

	// add claimBasic in User ID Merkle Tree
	err = userMT.Add(claimBasic.Entry())
	require.Nil(t, err)

	// add claimAuthKSign in User ID Merkle Tree
	err = userMT.Add(claimAuthKSign.Entry())
	require.Nil(t, err)

	// add claimAuthKSignBabyJub in User ID Merkle Tree
	err = userMT.Add(claimAuthKSignBabyJub.Entry())
	require.Nil(t, err)

	// setRootClaim of the user in the Relay Merkle Tree
	setRootClaim, err := core.NewClaimSetRootKey(&id, userMT.RootKey())
	require.Nil(t, err)
	// setRootClaim.BaseIndex.Version++ // TODO autoincrement
	// add User's ID Merkle Root into the Relay's Merkle Tree
	err = mt.Add(setRootClaim.Entry())
	require.Nil(t, err)

	rootService.On("GetRoot", &relayID).Return(
		&core.RootData{BlockN: 123, BlockTimestamp: 456, Root: mt.RootKey()}, nil).Once()
	proofClaim, err := service.GetClaimProofByHiBlockchain(setRootClaim.Entry().HIndex())
	require.Nil(t, err)
	p, err := json.Marshal(proofClaim)
	require.Nil(t, err)
	if debug {
		fmt.Printf("\n\tSetRoot claim proof\n\n")
		fmt.Println(string(p))
	}

	ok, err := proofClaim.Verify(proofClaim.Proof.Root)
	if !ok || err != nil {
		panic(err)
	}

	// ClaimAssignName
	// id, err := core.IDFromString("1oqcKzijA2tyUS6tqgGWoA1jLiN1gS5sWRV6JG8XY")
	// require.Nil(t, err)
	claimAssignName := core.NewClaimAssignName("testName@iden3.eth", id)
	// add assignNameClaim in User ID Merkle Tree
	err = mt.Add(claimAssignName.Entry())
	require.Nil(t, err)
	fmt.Printf("> A %+v\n", mt.RootKey().String())
	rootService.On("GetRoot", &relayID).Return(
		&core.RootData{BlockN: 123, BlockTimestamp: 456, Root: mt.RootKey()}, nil).Once()
	fmt.Printf("> R %+v\n", mt.RootKey().String())
	proofClaimAssignName, err := service.GetClaimProofByHiBlockchain(claimAssignName.Entry().HIndex())
	require.Nil(t, err)
	p, err = json.Marshal(proofClaimAssignName)
	if debug {
		fmt.Printf("\n\tclaim assign name claim proof\n\n")
		fmt.Println(string(p))
	}

	//proofClaim, err := service.GetClaimProofUserByHi(ethAddr, *claim.Entry().HIndex())
	//require.Nil(t, err)
	//if err != nil {
	//	panic(err)
	//}
	//claimProof := proofClaim.ClaimProof
	//claimNonRevocationProof := proofClaim.ClaimNonRevocationProof
	//setRootClaimProof := proofClaim.SetRootClaimProof
	//setRootClaimNonRevocationProof := proofClaim.SetRootClaimNonRevocationProof

	//require.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074cfee7c08a98f4b565d124c7e4e28acc52e1bc780e3887db000000048000000006461746161736466", common3.HexEncode(claimProof.Leaf))
	//require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000000", common3.HexEncode(claimProof.Proof))
	//require.Equal(t, "0x1415376b054a9ab3c7f9bd0ec956b0f403ae98d7e37dcbafdf26b465b23dd970", claimProof.Root.Hex())
	//require.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c49694030749b9a76a0132a0814192c05c9321efc30c7286f6187f18fc60000005400000000970e8128ab834e8eac17ab8e3812f010678cf7911415376b054a9ab3c7f9bd0ec956b0f403ae98d7e37dcbafdf26b465b23dd970", common3.HexEncode(setRootClaimProof.Leaf))
	//require.Equal(t, "0x000000000000000000000000000000000000000000000000000000000000000474c3e76aebd3df03ff91325d245e72ea9ad9599777f5d2c5e560b3f049d68309", common3.HexEncode(setRootClaimProof.Proof))
	//require.Equal(t, "0xf73c98cbaa1d43ada4ed5520300c348985dd47cc283e3cf7186434a07a46886a", setRootClaimProof.Root.Hex())
	//require.Equal(t, "0x00000000000000000000000000000000000000000000000000000000000000025563046fb69f065953f0fdb0b3033f721457184adfae2824c02932090bf8f281", common3.HexEncode(claimNonRevocationProof.Proof))
	//require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000014367d7c39348c9b7f2d488a7cd2edfbae56d608ec92a1b0a747adda3c4aaf763d74c3e76aebd3df03ff91325d245e72ea9ad9599777f5d2c5e560b3f049d68309", common3.HexEncode(setRootClaimNonRevocationProof.Proof))

	//require.Equal(t, claimProof.Root.Bytes(), claimNonRevocationProof.Root.Bytes())
	//require.Equal(t, setRootClaimProof.Root.Bytes(), setRootClaimNonRevocationProof.Root.Bytes())

	//var leafBytes [128]byte
	//copy(leafBytes[:], claimProof.Leaf)
	//entry := merkletree.Entry{Data: *merkletree.BytesToData(leafBytes)}
	//proof, err := merkletree.NewProofFromBytes(claimProof.Proof)
	//require.Nil(t, err)
	//verified := merkletree.VerifyProof(&claimProof.Root, proof, entry.HIndex(), entry.HValue())
	//require.True(t, verified)

	//leafBytes = [128]byte{}
	//copy(leafBytes[:], setRootClaimProof.Leaf)
	//entry = merkletree.Entry{Data: *merkletree.BytesToData(leafBytes)}
	//proof, err = merkletree.NewProofFromBytes(setRootClaimProof.Proof)
	//require.Nil(t, err)
	//verified = merkletree.VerifyProof(&setRootClaimProof.Root, proof, entry.HIndex(), entry.HValue())
	//require.True(t, verified)

	//leafBytes = [128]byte{}
	//copy(leafBytes[:], claimNonRevocationProof.Leaf)
	//entry = merkletree.Entry{Data: *merkletree.BytesToData(leafBytes)}
	//proof, err = merkletree.NewProofFromBytes(claimNonRevocationProof.Proof)
	//require.Nil(t, err)
	//verified = merkletree.VerifyProof(&claimNonRevocationProof.Root, proof, entry.HIndex(), entry.HValue())
	//require.True(t, verified)

	//leafBytes = [128]byte{}
	//copy(leafBytes[:], setRootClaimNonRevocationProof.Leaf)
	//entry = merkletree.Entry{Data: *merkletree.BytesToData(leafBytes)}
	//proof, err = merkletree.NewProofFromBytes(setRootClaimNonRevocationProof.Proof)
	//require.Nil(t, err)
	//verified = merkletree.VerifyProof(&setRootClaimNonRevocationProof.Root, proof, entry.HIndex(), entry.HValue())
	//require.True(t, verified)
}

func TestMain(m *testing.M) {
	result := m.Run()
	for _, dir := range rmDirs {
		os.RemoveAll(dir)
	}
	os.Exit(result)
}
