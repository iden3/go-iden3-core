package claimsrv

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	// "time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/crypto/babyjub"
	"github.com/iden3/go-iden3/db"
	babykeystore "github.com/iden3/go-iden3/keystore"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/signsrv"
	// "github.com/iden3/go-iden3/utils"
	"github.com/ipfsconsortium/go-ipfsc/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var debug = true

var service *ServiceImpl
var mt *merkletree.MerkleTree
var c config.Config

const relayIdHex = "1N7d2qVEJeqnYAWVi5Cq6PLj6GwxaW6FYcfmY2fps"
const relaySkHex = "4be5471a938bdf3606888472878baace4a6a64e14a153adf9a1333969e4e573c"

var relayID core.ID
var keyStore *babykeystore.KeyStore
var relaySk babyjub.PrivateKey
var relayPkComp *babyjub.PublicKeyComp
var relayPk *babyjub.PublicKey

type RootServiceMock struct {
	mock.Mock
}

func (m *RootServiceMock) Start() {

}
func (m *RootServiceMock) StopAndJoin() {

}

func (m *RootServiceMock) GetRoot(addr common.Address) (merkletree.Hash, error) {
	args := m.Called(addr)
	return args.Get(0).(merkletree.Hash), args.Error(1)
}
func (m *RootServiceMock) SetRoot(hash merkletree.Hash) {
	m.Called(hash)
}

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

func newTestingMerkle(numLevels int) (*merkletree.MerkleTree, error) {
	dir, err := ioutil.TempDir("", "db")
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

	id, err := core.IDFromString(relayIdHex)
	assert.Nil(t, err)

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
	service = New(id, mt, &RootServiceMock{}, *signSrv)

	relayID, err = core.IDFromString(relayIdHex)
	if err != nil {
		panic(err)
	}
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
	assert.Nil(t, err)
	assert.Equal(t, uint32(0), version)
	claim.Version = version
	mt.Add(claim.Entry())
	version, err = GetNextVersion(mt, claim.Entry().HIndex())
	assert.Nil(t, err)
	assert.Equal(t, uint32(1), version)

	claim.Version = version
	mt.Add(claim.Entry())
	version, err = GetNextVersion(mt, claim.Entry().HIndex())
	assert.Nil(t, err)
	assert.Equal(t, uint32(2), version)

	claim.Version = version
	mt.Add(claim.Entry())
	version, err = GetNextVersion(mt, claim.Entry().HIndex())
	assert.Nil(t, err)
	assert.Equal(t, uint32(3), version)

	claim.Version = version
	mt.Add(claim.Entry())
	version, err = GetNextVersion(mt, claim.Entry().HIndex())
	assert.Nil(t, err)
	assert.Equal(t, uint32(4), version)
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
	assert.Nil(t, err)
	version, err := GetNextVersion(mt, claim.Entry().HIndex())
	assert.Nil(t, err)
	assert.Equal(t, uint32(1), version)

	claimProof, err := getNonRevocationProof(mt, *claim.Entry().HIndex())
	assert.Nil(t, err)

	assert.Equal(t,
		"0x030000000000000000000000000000000000000000000000000000000000000015d78037f4d6da9bfd2f3f308944e7b58e008ae10d8499c776e00c95d14f783306d4571fb9634e4bed32e265f91a373a852c476656c5c13b09bc133ac61bc5a6",
		common3.HexEncode(claimProof.Proof))
	assert.Equal(t,
		"0x1584066523678cd35a24da5feafeda02d8a02520698d9916b8f1e4f1d08254e1",
		claimProof.Root.Hex())

	proof, err := merkletree.NewProofFromBytes(claimProof.Proof)
	assert.Nil(t, err)

	// VerifyProof with HIndex and HValue from claimProof.Leaf
	var leafBytes [128]byte
	copy(leafBytes[:], claimProof.Leaf)
	dataNonE := merkletree.NewDataFromBytes(leafBytes)
	claimProofEntry := merkletree.Entry{
		Data: *dataNonE,
	}
	verified := merkletree.VerifyProof(&claimProof.Root, proof, claimProofEntry.HIndex(), claimProofEntry.HValue())
	assert.True(t, verified)
}

func TestGetClaimProof(t *testing.T) {
	initializeEnvironment(t)

	id, err := core.IDFromString("1pnWU7Jdr4yLxp1azs1r1PpvfErxKGRQdcLBZuq3Z")
	assert.Nil(t, err)

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
	userMT, err := NewMerkleTreeUser(id, mt.Storage(), 140)
	assert.Nil(t, err)

	// add claimBasic in User ID Merkle Tree
	err = userMT.Add(claimBasic.Entry())
	assert.Nil(t, err)

	// add claimAuthKSign in User ID Merkle Tree
	err = userMT.Add(claimAuthKSign.Entry())
	assert.Nil(t, err)

	// add claimAuthKSignBabyJub in User ID Merkle Tree
	err = userMT.Add(claimAuthKSignBabyJub.Entry())
	assert.Nil(t, err)

	// setRootClaim of the user in the Relay Merkle Tree
	setRootClaim := core.NewClaimSetRootKey(id, *userMT.RootKey())
	// setRootClaim.BaseIndex.Version++ // TODO autoincrement
	// add User's ID Merkle Root into the Relay's Merkle Tree
	err = mt.Add(setRootClaim.Entry())
	assert.Nil(t, err)

	proofClaim, err := service.GetClaimProofByHi(setRootClaim.Entry().HIndex())
	assert.Nil(t, err)
	p, err := json.Marshal(proofClaim)
	assert.Nil(t, err)
	if debug {
		fmt.Printf("\n\tSetRoot claim proof\n\n")
		fmt.Println(string(p))
	}

	ok, err := core.VerifyProofClaim(relayPk, proofClaim)
	if !ok || err != nil {
		panic(err)
	}

	proofClaimUser, err := service.GetClaimProofUserByHi(id, claimBasic.Entry().HIndex())
	assert.Nil(t, err)
	p, err = json.Marshal(proofClaimUser)
	if debug {
		fmt.Printf("\n\tclaim basic claim proof\n\n")
		fmt.Println(string(p))
	}

	ok, err = core.VerifyProofClaim(relayPk, proofClaimUser)
	assert.Equal(t, ok, true)
	assert.Nil(t, err)

	proofClaimUser2, err := service.GetClaimProofUserByHi(id, claimAuthKSign.Entry().HIndex())
	assert.Nil(t, err)
	p, err = json.Marshal(proofClaimUser2)
	if debug {
		fmt.Printf("\n\tclaim authorize ksign secp256k1 claim proof\n\n")
		fmt.Println("id", id.String())
		fmt.Println("pk", common3.HexEncode(crypto.CompressPubkey(pk)))
		fmt.Println(string(p))
	}

	ok, err = core.VerifyProofClaim(relayPk, proofClaimUser2)
	assert.Equal(t, ok, true)
	assert.Nil(t, err)

	ok, err = core.VerifyProofClaim(relayPk, proofClaimUser)
	assert.Equal(t, ok, true)
	assert.Nil(t, err)

	proofClaimUser3, err := service.GetClaimProofUserByHi(id, claimAuthKSignBabyJub.Entry().HIndex())
	assert.Nil(t, err)
	p, err = json.Marshal(proofClaimUser3)
	if debug {
		fmt.Printf("\n\tclaim authorize ksign babyjub claim proof\n\n")
		fmt.Println("id", id.String())
		fmt.Println("pk", kSignPk)
		fmt.Println(string(p))
	}

	// ClaimAssignName
	// id, err := core.IDFromString("1oqcKzijA2tyUS6tqgGWoA1jLiN1gS5sWRV6JG8XY")
	// assert.Nil(t, err)
	claimAssignName := core.NewClaimAssignName("testName@iden3.eth", id)
	// add assignNameClaim in User ID Merkle Tree
	err = mt.Add(claimAssignName.Entry())
	assert.Nil(t, err)
	proofClaimAssignName, err := service.GetClaimProofByHi(claimAssignName.Entry().HIndex())
	assert.Nil(t, err)
	p, err = json.Marshal(proofClaimAssignName)
	if debug {
		fmt.Printf("\n\tclaim assign name claim proof\n\n")
		fmt.Println(string(p))
	}

	//proofClaim, err := service.GetClaimProofUserByHi(ethAddr, *claim.Entry().HIndex())
	//assert.Nil(t, err)
	//if err != nil {
	//	panic(err)
	//}
	//claimProof := proofClaim.ClaimProof
	//claimNonRevocationProof := proofClaim.ClaimNonRevocationProof
	//setRootClaimProof := proofClaim.SetRootClaimProof
	//setRootClaimNonRevocationProof := proofClaim.SetRootClaimNonRevocationProof

	//assert.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074cfee7c08a98f4b565d124c7e4e28acc52e1bc780e3887db000000048000000006461746161736466", common3.HexEncode(claimProof.Leaf))
	//assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000000", common3.HexEncode(claimProof.Proof))
	//assert.Equal(t, "0x1415376b054a9ab3c7f9bd0ec956b0f403ae98d7e37dcbafdf26b465b23dd970", claimProof.Root.Hex())
	//assert.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c49694030749b9a76a0132a0814192c05c9321efc30c7286f6187f18fc60000005400000000970e8128ab834e8eac17ab8e3812f010678cf7911415376b054a9ab3c7f9bd0ec956b0f403ae98d7e37dcbafdf26b465b23dd970", common3.HexEncode(setRootClaimProof.Leaf))
	//assert.Equal(t, "0x000000000000000000000000000000000000000000000000000000000000000474c3e76aebd3df03ff91325d245e72ea9ad9599777f5d2c5e560b3f049d68309", common3.HexEncode(setRootClaimProof.Proof))
	//assert.Equal(t, "0xf73c98cbaa1d43ada4ed5520300c348985dd47cc283e3cf7186434a07a46886a", setRootClaimProof.Root.Hex())
	//assert.Equal(t, "0x00000000000000000000000000000000000000000000000000000000000000025563046fb69f065953f0fdb0b3033f721457184adfae2824c02932090bf8f281", common3.HexEncode(claimNonRevocationProof.Proof))
	//assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000014367d7c39348c9b7f2d488a7cd2edfbae56d608ec92a1b0a747adda3c4aaf763d74c3e76aebd3df03ff91325d245e72ea9ad9599777f5d2c5e560b3f049d68309", common3.HexEncode(setRootClaimNonRevocationProof.Proof))

	//assert.Equal(t, claimProof.Root.Bytes(), claimNonRevocationProof.Root.Bytes())
	//assert.Equal(t, setRootClaimProof.Root.Bytes(), setRootClaimNonRevocationProof.Root.Bytes())

	//var leafBytes [128]byte
	//copy(leafBytes[:], claimProof.Leaf)
	//entry := merkletree.Entry{Data: *merkletree.BytesToData(leafBytes)}
	//proof, err := merkletree.NewProofFromBytes(claimProof.Proof)
	//assert.Nil(t, err)
	//verified := merkletree.VerifyProof(&claimProof.Root, proof, entry.HIndex(), entry.HValue())
	//assert.True(t, verified)

	//leafBytes = [128]byte{}
	//copy(leafBytes[:], setRootClaimProof.Leaf)
	//entry = merkletree.Entry{Data: *merkletree.BytesToData(leafBytes)}
	//proof, err = merkletree.NewProofFromBytes(setRootClaimProof.Proof)
	//assert.Nil(t, err)
	//verified = merkletree.VerifyProof(&setRootClaimProof.Root, proof, entry.HIndex(), entry.HValue())
	//assert.True(t, verified)

	//leafBytes = [128]byte{}
	//copy(leafBytes[:], claimNonRevocationProof.Leaf)
	//entry = merkletree.Entry{Data: *merkletree.BytesToData(leafBytes)}
	//proof, err = merkletree.NewProofFromBytes(claimNonRevocationProof.Proof)
	//assert.Nil(t, err)
	//verified = merkletree.VerifyProof(&claimNonRevocationProof.Root, proof, entry.HIndex(), entry.HValue())
	//assert.True(t, verified)

	//leafBytes = [128]byte{}
	//copy(leafBytes[:], setRootClaimNonRevocationProof.Leaf)
	//entry = merkletree.Entry{Data: *merkletree.BytesToData(leafBytes)}
	//proof, err = merkletree.NewProofFromBytes(setRootClaimNonRevocationProof.Proof)
	//assert.Nil(t, err)
	//verified = merkletree.VerifyProof(&setRootClaimNonRevocationProof.Root, proof, entry.HIndex(), entry.HValue())
	//assert.True(t, verified)
}

/*
func TestAssignNameClaim(t *testing.T) {
	initializeEnvironment()
	testPrivK, err := crypto.HexToECDSA(testPrivKHex)
	assert.Nil(t, err)

	mt.Add(core.NewGenericClaim("c1", "c1", []byte("c1")))
	mt.Add(core.NewGenericClaim("c2", "c2", []byte("c2")))
	mt.Add(core.NewGenericClaim("c3", "c3", []byte("c3")))

	nameHash := merkletree.HashBytes([]byte("johndoe"))
	domainHash := merkletree.HashBytes([]byte(c.Domain))
	ethAddr := crypto.PubkeyToAddress(testPrivK.PublicKey)
	assignNameClaim := core.NewAssignNameClaim(c.Namespace, nameHash, domainHash, ethAddr)
	// signature, err := utils.Sign(assignNameClaim.Ht(), testPrivK)
	// assert.Nil(t, err)
	// signatureHex := common3.HexEncode(signature)
	// assignNameClaimMsg := AssignNameClaimMsg{
	// 	assignNameClaim,
	// 	signatureHex,
	// }
	privK, err := crypto.HexToECDSA(c.Server.PrivK)
	assert.Nil(t, err)
	root, mp, _, err := AddAssignNameClaim(mt, assignNameClaim, c.ContractsAddress.Identities, privK)
	assert.Nil(t, err)
	mtRoot := mt.Root()
	if !bytes.Equal(root[:], mtRoot[:]) {
		t.Errorf("root != mt.Root")
	}
	expectedRootHex := "0x05175b7c17ea772423da35f9ccd0bb0017355a135e60ba28e541f26e1185b31e"
	if mt.Root().Hex() != expectedRootHex {
		t.Errorf("mt.Root: " + mt.Root().Hex() + " , expected root: " + expectedRootHex)
	}
	expectedMPHex := "0x000000000000000000000000000000000000000000000000000000000000000311a689079d0478b829d23ae5fb3e65ab15ad1abc364eea2965abf1c324e72e817370e48c8a338794dd181314bbd080e4263a802803686bcc2c2d3f554e3a50de"
	if common3.HexEncode(mp) != expectedMPHex {
		t.Errorf("mp: " + common3.HexEncode(mp) + " , expected mp: " + expectedMPHex)
	}
}

func TestResolvAssignNameClaim(t *testing.T) {
	nameHash := merkletree.HashBytes([]byte("johndoe"))
	domainHash := merkletree.HashBytes([]byte(c.Domain))
	testPrivK, err := crypto.HexToECDSA(testPrivKHex)
	ethAddr := crypto.PubkeyToAddress(testPrivK.PublicKey)
	originalAssignNameClaim := core.NewAssignNameClaim(c.Namespace, nameHash, domainHash, ethAddr)
	assignNameClaim, err := ResolvAssignNameClaim(mt, "johndoe@iden3.io", c.Namespace)
	if err != nil {
		t.Errorf(err.Error())
	}
	if !bytes.Equal(assignNameClaim.Bytes(), originalAssignNameClaim.Bytes()) {
		t.Errorf("resolved AssignNameClaim != original AssignNameClaim")
	}
}

func TestNewAuthorizeKSignClaim(t *testing.T) {
	testPrivK, err := crypto.HexToECDSA(testPrivKHex)
	if err != nil {
		t.Errorf(err.Error())
	}
	testAddr := crypto.PubkeyToAddress(testPrivK.PublicKey)

	authorizeKSignClaim := core.NewAuthorizeKSignClaim("iden3.io", testAddr, "app1", "appauthz", 1535208350, 1535208350)
	msgHash := utils.EthHash(authorizeKSignClaim.Bytes())
	signature, err := utils.Sign(msgHash, testPrivK)
	assert.Nil(t, err)
	signatureHex := common3.HexEncode(signature)
	authorizeKSignClaimMsg := AuthorizeKSignClaimMsg{
		authorizeKSignClaim,
		signatureHex,
	}
	claimProof, idRootProof, err := AddAuthorizeKSignClaim(mt, testAddr, authorizeKSignClaimMsg, c.ContractsAddress.Identities)
	assert.Nil(t, err)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, "0x771e1ef9fab9bdf7f55ba7c24112b9c4b9d7e55cd94f57efd0fd4ef174565b66", mt.Root().Hex())

	// check userIDRoot
	stoUserID := mt.Storage().WithPrefix(testAddr.Bytes())
	userMT, err := merkletree.New(stoUserID, 140)
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, "0x8112699ee0bb1a6307dce979a72d77549fdcf1d59648b424c5d65d5080d4b3fa", userMT.Root().Hex())

	expectedClaimProof := "0x0000000000000000000000000000000000000000000000000000000000000000"
	assert.Equal(t, expectedClaimProof, common3.HexEncode(claimProof))
	expectedIdRootProof := "0x000000000000000000000000000000000000000000000000000000000000000730c5c5fe05516470d1963cde3ecc1b93b73b2b4d09e37a4151685d6af5260705d827465cbe023bbcfa073720ce38ab510064b1743310cca89b00fb807ca3b37e7370e48c8a338794dd181314bbd080e4263a802803686bcc2c2d3f554e3a50de"
	assert.Equal(t, expectedIdRootProof, common3.HexEncode(idRootProof))

}

func TestMultipleAuthorizeKSignClaim(t *testing.T) {
	privKHex := "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
	testPrivK, err := crypto.HexToECDSA(privKHex)
	assert.Nil(t, err)
	testAddr := crypto.PubkeyToAddress(testPrivK.PublicKey)

	authorizeKSignClaim := core.NewAuthorizeKSignClaim("iden3.io", testAddr, "app1", "appauthz", 1535208355, 1535208355)
	msgHash := utils.EthHash(authorizeKSignClaim.Bytes())
	signature, err := utils.Sign(msgHash, testPrivK)
	assert.Nil(t, err)
	signatureHex := common3.HexEncode(signature)
	authorizeKSignClaimMsg := AuthorizeKSignClaimMsg{
		authorizeKSignClaim,
		signatureHex,
	}
	claimProof, idRootProof, err := AddAuthorizeKSignClaim(mt, testAddr, authorizeKSignClaimMsg, c.ContractsAddress.Identities)
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, "0xab8da27ef1d44f3853242f095892280390d60932f2dfdd6a9988a67f6cec35ec", mt.Root().Hex())

	stoUserID := mt.Storage().WithPrefix(testAddr.Bytes())
	userMT, err := merkletree.New(stoUserID, 140)
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, "0xbdb2b31ecb9c674995f29a9bdb74065172a85e0e135c56274f8e17137451c684", userMT.Root().Hex())

	expectedClaimProof := "0x0000000000000000000000000000000000000000000000000000000000000000"
	assert.Equal(t, expectedClaimProof, common3.HexEncode(claimProof))
	expectedIdRootProof := "0x0000000000000000000000000000000000000000000000000000000000000007585169e90e5f14f720529326b75be5fe9c4fbe0e78874c8db3c2c0fe879c87062fd3493fd39f4bd7a626383d2617bf4ead5e47941cdbe4e941edcb0bb8626b357370e48c8a338794dd181314bbd080e4263a802803686bcc2c2d3f554e3a50de"
	assert.Equal(t, expectedIdRootProof, common3.HexEncode(idRootProof))

	privKHex2 := "a247c1a3ab5c894d68575fad9f9a519895732ba7b8b0c22afce255338ae8c345"
	testPrivK2, err := crypto.HexToECDSA(privKHex2)
	assert.Nil(t, err)
	testAddr2 := crypto.PubkeyToAddress(testPrivK2.PublicKey)
	authorizeKSignClaim2 := core.NewAuthorizeKSignClaim("iden3.io", testAddr2, "app1", "appauthz", 1535208355, 1535208355)
	msgHash = utils.EthHash(authorizeKSignClaim2.Bytes())
	signature2, err := utils.Sign(msgHash, testPrivK2)
	assert.Nil(t, err)
	signatureHex2 := common3.HexEncode(signature2)
	authorizeKSignClaimMsg2 := AuthorizeKSignClaimMsg{
		authorizeKSignClaim2,
		signatureHex2,
	}
	claimProof2, idRootProof2, err := AddAuthorizeKSignClaim(mt, testAddr2, authorizeKSignClaimMsg2, c.ContractsAddress.Identities)
	assert.Nil(t, err)

	assert.Equal(t, "0xf6c57457fd9ebcd6c21acd511a41303f63e59e74c7c47d98fd0813a9bf39b392", mt.Root().Hex())

	stoUserID2 := mt.Storage().WithPrefix(testAddr2.Bytes())
	userMT2, err := merkletree.New(stoUserID2, 140)
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, "0xfea5cdf67c17737bf9b148a6dc26449c1672b59d37116b916253f0abce72f160", userMT2.Root().Hex())
	expectedClaimProof = "0x0000000000000000000000000000000000000000000000000000000000000000"
	assert.Equal(t, expectedClaimProof, common3.HexEncode(claimProof2))
	expectedIdRootProof = "0x000000000000000000000000000000000000000000000000000000000000001713bc31bd2a88624073b508ade2ce7e8a2207c53b12f0dbdfc4547362d6376e1312610bb2a7c84995083296c0b3eada2d57184d2b4f02adb907a649d7748c614ad25b5563e50227d3c4ff6b9161f5381a292a998ae7d53ec74960ece6a04f5fb07370e48c8a338794dd181314bbd080e4263a802803686bcc2c2d3f554e3a50de"
	assert.Equal(t, expectedIdRootProof, common3.HexEncode(idRootProof2))
}

func TestNewUserIDClaim(t *testing.T) {
	privKHex := "da7079f082a1ced80c5dee3bf00752fd67f75321a637e5d5073ce1489af062d8"
	testPrivK, err := crypto.HexToECDSA(privKHex)
	assert.Nil(t, err)
	testAddr := crypto.PubkeyToAddress(testPrivK.PublicKey)

	claim := core.NewGenericClaim("iden3.io_3", "default", []byte("data2"))
	signature, err := utils.Sign(claim.Ht(), testPrivK)
	assert.Nil(t, err)
	signatureHex := common3.HexEncode(signature)
	claimValueMsg := ClaimValueMsg{
		claim,
		signatureHex,
	}
	claimProof, idRootProof, err := AddUserIDClaim(mt, "iden3.io", testAddr, claimValueMsg, c.ContractsAddress.Identities)
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, "0x964e3bb814386a83eb85ccca2f09812bdb9582afa30fe1e454c5f4dfcb6bd70e", mt.Root().Hex())

	// check userIDRoot
	stoUserID := mt.Storage().WithPrefix(testAddr.Bytes())
	userMT, err := merkletree.New(stoUserID, 140)
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, "0xcda67faf66cf1261e9653c2528883f7e7c6fa4ea9ddef2a3b669817e0b2d1bbc", userMT.Root().Hex())

	expectedClaimProof := "0x000000000000000000000000000000000000000000000000000000000000000257a42f22a7e9b3acf712f7bb8a4e684f965f8c3ee2dc0fc88129c8b624754fcd"
	assert.Equal(t, expectedClaimProof, common3.HexEncode(claimProof))
	expectedIdRootProof := "0x0000000000000000000000000000000000000000000000000000000000000107f3e6294d5cb4ef3ff284318ddce1f539111c3310e04075420b89dac28d1003b15def58d649018d988ff4d4c7cf9cbc4ab7590d58fa06e76b28f802212e2b5083f9e894a02f51799114c844c03d5859069afb4c7287a5403c6c4fba577467bed57370e48c8a338794dd181314bbd080e4263a802803686bcc2c2d3f554e3a50de"
	assert.Equal(t, expectedIdRootProof, common3.HexEncode(idRootProof))

}
func TestGetIDRoot(t *testing.T) {
	privKHex := "da7079f082a1ced80c5dee3bf00752fd67f75321a637e5d5073ce1489af062d8"
	testPrivK, err := crypto.HexToECDSA(privKHex)
	assert.Nil(t, err)
	testAddr := crypto.PubkeyToAddress(testPrivK.PublicKey)

	idRoot, idRootProof, err := GetIDRoot(mt, testAddr)
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.Equal(t, "0xcda67faf66cf1261e9653c2528883f7e7c6fa4ea9ddef2a3b669817e0b2d1bbc", idRoot.Hex())
	expectedProof := "0x0000000000000000000000000000000000000000000000000000000000000007ab9ed10e59863ed65028fda65d43dc320388afd2ff6510e0677d04acf376e89f4f7c6e940a2392179ceb7120d4a3127bd7918a3c0f7bf1726958523214fc73247370e48c8a338794dd181314bbd080e4263a802803686bcc2c2d3f554e3a50de"
	assert.Equal(t, expectedProof, common3.HexEncode(idRootProof))
}

func TestGetClaimByHiThatDontExist(t *testing.T) {
	privKHex := "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
	testPrivK, err := crypto.HexToECDSA(privKHex)
	assert.Nil(t, err)
	testAddr := crypto.PubkeyToAddress(testPrivK.PublicKey)

	hiHex := "0x784adb4a490b9c0521c11298f384bf847881711f1a522a40129d76e3cfc68c9a"
	hiBytes, err := common3.HexDecode(hiHex)
	assert.Nil(t, err)
	hi := merkletree.Hash{}
	copy(hi[:], hiBytes)
	_, _, _, _, err = GetClaimByHi(mt, "namespace.io", testAddr, hi)
	assert.NotNil(t, err)
}

func TestAddClaimAndGetClaimByHi(t *testing.T) {
	privKHex := "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
	testPrivK, err := crypto.HexToECDSA(privKHex)
	assert.Nil(t, err)
	testAddr := crypto.PubkeyToAddress(testPrivK.PublicKey)

	claim := core.NewGenericClaim("namespace.io", "default", []byte("dataasdf"))
	signature, err := utils.Sign(claim.Ht(), testPrivK)
	assert.Nil(t, err)
	signatureHex := common3.HexEncode(signature)
	claimValueMsg := ClaimValueMsg{
		claim,
		signatureHex,
	}
	claimProof1, idRootProof1, err := AddUserIDClaim(mt, "namespace.io", testAddr, claimValueMsg, c.ContractsAddress.Identities)
	assert.Nil(t, err)
	hi := claim.Hi()
	claimProof, setRootClaimProof, claimNonRevocationProof, setRootClaimNonRevocationProof, err := GetClaimByHi(mt, "namespace.io", testAddr, hi)
	if err != nil {
		panic(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, "0xa92591b1ee18783f95fbf358517afed09d888b9db8286c0d19e2419036941d68cfee7c08a98f4b565d124c7e4e28acc52e1bc780e3887db0a02a7d2d5bc66728000000006461746161736466", common3.HexEncode(claimProof.Leaf))
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000002546f8feb74144a5ee688f26ee5c5c202051386d6682164b1746d7481c4c5fda0", common3.HexEncode(claimProof.Proof))
	assert.Equal(t, "0x174798396a958603a3c6b2f60b21a4735000429be4d5dded269b93ba37945898", claimProof.Root.Hex())
	assert.Equal(t, "0x000000000000000000000000000000000000000000000000000000000000000325030b375e7fb70ce357852c717818479d67f15003b30048798c61d8a3e381fc7e57e8df413edef8ca83461bccf69e18815802e3815765b7384185aca868a7f6", common3.HexEncode(setRootClaimProof.Proof))
	assert.Equal(t, "0x7b71af6e80b3db0c67ee967e46808fd42a0f87b82c6068ced1007297261320f5", setRootClaimProof.Root.Hex())

	assert.Equal(t, claimProof1, claimProof.Proof)
	assert.Equal(t, idRootProof1, setRootClaimProof.Proof)
	verified := merkletree.CheckProof(claimProof.Root, claimProof.Proof, claimProof.Hi, merkletree.HashBytes(claimProof.Leaf), 140)
	assert.True(t, verified)
	assert.Equal(t, mt.Root().Bytes(), setRootClaimProof.Root.Bytes())
	verified = merkletree.CheckProof(setRootClaimProof.Root, setRootClaimProof.Proof, setRootClaimProof.Hi, merkletree.HashBytes(setRootClaimProof.Leaf), mt.NumLevels())
	assert.True(t, verified)
	verified = merkletree.CheckProof(claimNonRevocationProof.Root, claimNonRevocationProof.Proof, claimNonRevocationProof.Hi, merkletree.EmptyNodeValue, 140)
	assert.True(t, verified)
	verified = merkletree.CheckProof(setRootClaimNonRevocationProof.Root, setRootClaimNonRevocationProof.Proof, setRootClaimNonRevocationProof.Hi, merkletree.EmptyNodeValue, 140)
	assert.True(t, verified)
}
*/
