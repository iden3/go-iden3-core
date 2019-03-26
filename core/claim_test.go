package core

import (
	//"bytes"
	//"encoding/hex"
	"bytes"
	"crypto/ecdsa"

	//"crypto/elliptic"
	//"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	//common3 "github.com/iden3/go-iden3/common"
	//"github.com/iden3/go-iden3/db"
	"github.com/iden3/go-iden3/merkletree"
	//"github.com/iden3/go-iden3/utils"
	"github.com/stretchr/testify/assert"
)

var debug = false

func TestClaimBasic(t *testing.T) {
	// ClaimBasic
	indexSlot := [400 / 8]byte{
		0x29, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a,
		0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a,
		0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a,
		0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a,
		0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a,
		0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a, 0x2a,
		0x2a, 0x2b}
	dataSlot := [496 / 8]byte{
		0x56, 0x58, 0x58, 0x58, 0x58, 0x58, 0x58, 0x58,
		0x58, 0x58, 0x58, 0x58, 0x58, 0x58, 0x58, 0x58,
		0x58, 0x58, 0x58, 0x58, 0x58, 0x58, 0x58, 0x58,
		0x58, 0x58, 0x58, 0x58, 0x58, 0x58, 0x58, 0x58,
		0x58, 0x58, 0x58, 0x58, 0x58, 0x58, 0x58, 0x58,
		0x58, 0x58, 0x58, 0x58, 0x58, 0x58, 0x58, 0x58,
		0x58, 0x58, 0x58, 0x58, 0x58, 0x58, 0x58, 0x58,
		0x58, 0x58, 0x58, 0x58, 0x58, 0x59}
	c0 := NewClaimBasic(indexSlot, dataSlot)
	c0.Version = 1
	e := c0.Entry()
	assert.Equal(t,
		"0x10dae586534145cf6bf10cb0f7dc7f0a343b6fb32c04501e85cace56c9381f73",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x00fa98173f610358b765b4ff593ed7bae0db326d287eec86b1782da9bea45748",
		e.HValue().Hex())
	dataTestOutput(&e.Data)
	assert.Equal(t, ""+
		"0056585858585858585858585858585858585858585858585858585858585858"+
		"0058585858585858585858585858585858585858585858585858585858585859"+
		"00292a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a"+
		"002a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2b000000010000000000000000",
		e.Data.String())
	c1 := NewClaimBasicFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)
}

func TestClaimAssignName(t *testing.T) {
	// ClaimAssignName
	name := "example.iden3.eth"
	ethAddr := common.BytesToAddress([]byte{
		0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39,
		0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39,
		0x39, 0x39, 0x39, 0x3a})
	c0 := NewClaimAssignName(name, ethAddr)
	c0.Version = 1
	e := c0.Entry()
	assert.Equal(t,
		"0x1b90afca6decbbbeb9422a441731f7c9a666b3ce15307d327f43b2c4506fe29c",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x19ff0e5ed571621bbc36bf4c0027dbe097975e11c03841b8c8d773b833b71f84",
		e.HValue().Hex())
	dataTestOutput(&e.Data)
	assert.Equal(t, ""+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"000000000000000000000000393939393939393939393939393939393939393a"+
		"00d67b05d8e2d1ace8f3e84b8451dd2e9da151578c3c6be23e7af11add5a807a"+
		"0000000000000000000000000000000000000000000000010000000000000003",
		e.Data.String())
	c1 := NewClaimAssignNameFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)
}

func TestClaimAuthorizeKSign(t *testing.T) {
	// ClaimAuthorizeKSign
	sign := true
	ay := merkletree.ElemBytes{
		0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05,
		0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05,
		0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05,
		0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x06}
	//ay := [128 / 8]byte{
	//	0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07,
	//	0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x08}
	c0 := NewClaimAuthorizeKSign(sign, ay)
	c0.Version = 1
	e := c0.Entry()
	assert.Equal(t,
		"0x2933b00c626c11d1dacadc6bbfbbdd7d008efb3db3565fbadc9dadfc8ac57eff",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x1541a6b5aa9bf7d9be3d5cb0bcc7cacbca26242016a0feebfc19c90f2224baed",
		e.HValue().Hex())
	dataTestOutput(&e.Data)
	assert.Equal(t, ""+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"0505050505050505050505050505050505050505050505050505050505050506"+
		"0000000000000000000000000000000000000001000000010000000000000001",
		e.Data.String())
	c1 := NewClaimAuthorizeKSignFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)
}

func TestClaimAuthorizeKSignSecp256k1(t *testing.T) {
	// ClaimAuthorizeKSignSecp256k1
	secKeyHex := "79156abe7fe2fd433dc9df969286b96666489bac508612d0e16593e944c4f69f"
	secKey, err := crypto.HexToECDSA(secKeyHex)
	if err != nil {
		panic(err)
	}
	pubKey := secKey.Public().(*ecdsa.PublicKey)
	assert.Equal(t,
		"036d94c84a7096c572b83d44df576e1ffb3573123f62099f8d4fa19de806bd4d59",
		hex.EncodeToString(crypto.CompressPubkey(pubKey)))
	c0 := NewClaimAuthorizeKSignSecp256k1(pubKey)
	c0.Version = 1
	e := c0.Entry()
	assert.Equal(t,
		"0x2f04e25bb5c9e01946aca5c80307d73bace3c3f1496175451db4495dd5a6c423",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x1541a6b5aa9bf7d9be3d5cb0bcc7cacbca26242016a0feebfc19c90f2224baed",
		e.HValue().Hex())
	dataTestOutput(&e.Data)
	assert.Equal(t, ""+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"00036d94c84a7096c572b83d44df576e1ffb3573123f62099f8d4fa19de806bd"+
		"0000000000000000000000000000000000004d59000000010000000000000004",
		e.Data.String())
	c1, err := NewClaimAuthorizeKSignSecp256k1FromEntry(e)
	if err != nil {
		panic(err)
	}
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)
}

func TestClaimSetRootKey(t *testing.T) {
	// ClaimSetRootKey
	ethAddr := common.BytesToAddress([]byte{
		0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39,
		0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39,
		0x39, 0x39, 0x39, 0x3a})
	rootKey := merkletree.Hash(merkletree.ElemBytes{
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0c})
	c0 := NewClaimSetRootKey(ethAddr, rootKey)
	c0.Version = 1
	c0.Era = 1
	e := c0.Entry()
	assert.Equal(t,
		"0x0f92abb8209409736929512c2018467a497ed35f409bb90579c62b9a4e0b2aa8",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x0ad7edbf562757b1ad2282c44e2c248f95e9e6b09ba0d32809aa724fbf148e0c",
		e.HValue().Hex())
	dataTestOutput(&e.Data)
	assert.Equal(t, ""+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0c"+
		"000000000000000000000000393939393939393939393939393939393939393a"+
		"0000000000000000000000000000000000000001000000010000000000000002",
		e.Data.String())
	c1 := NewClaimSetRootKeyFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)
}

func TestClaimLinkObjectIdentity(t *testing.T) {
	// ClaimLinkObjectIdentity
	var hashType uint32 = 1
	var objectType uint32 = 1
	var indexType uint16
	idAddr := common.BytesToAddress([]byte{
		0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39,
		0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39,
		0x39, 0x39, 0x39, 0x3a})
	objectHash := merkletree.Hash{
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0c}
	claim := NewClaimLinkObjectIdentity(hashType, objectType, indexType, idAddr, objectHash)
	claim.Version = 1
	entry := claim.Entry()
	assert.Equal(t,
		"0x2180103fca164a1f4afe33b610f8573076a24de6d7ca3d20fcf62540d3c73338",
		entry.HIndex().Hex())
	assert.Equal(t,
		"0x1bc5fac04887744f1831a8d881e4d68330e9f5f60a1c85265896df294329c7cc",
		entry.HValue().Hex())
	dataTestOutput(&entry.Data)
	assert.Equal(t, ""+
		"0000000000000000000000000000000000000000000000000000000000000001"+
		"000b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0c"+
		"000000000000000000000000393939393939393939393939393939393939393a"+
		"0000000000000000000000000000000000000001000000010000000000000005",
		entry.Data.String())
	c1 := NewClaimLinkObjectIdentityFromEntry(entry)
	c2, err := NewClaimFromEntry(entry)
	assert.Nil(t, err)
	assert.Equal(t, claim, c1)
	assert.Equal(t, claim, c2)
}

func dataTestOutput(d *merkletree.Data) {
	if !debug {
		return
	}
	s := bytes.NewBufferString("")
	fmt.Fprintf(s, "\t\t\"%v\"+\n", hex.EncodeToString(d[0][:]))
	fmt.Fprintf(s, "\t\t\"%v\"+\n", hex.EncodeToString(d[1][:]))
	fmt.Fprintf(s, "\t\t\"%v\"+\n", hex.EncodeToString(d[2][:]))
	fmt.Fprintf(s, "\t\t\"%v\",", hex.EncodeToString(d[3][:]))
	fmt.Println(s.String())
}

// TODO: Update to new claim spec.
//func TestForwardingInterop(t *testing.T) {
//
//	// address 0xee602447b5a75cf4f25367f5d199b860844d10c4
//	// pvk     8A85AAA2A8CE0D24F66D3EAA7F9F501F34992BACA0FF942A8EDF7ECE6B91F713
//
//	mt, err := merkletree.New(db.NewMemoryStorage(), 140)
//	assert.Nil(t, err)
//
//	// create ksignclaim ----------------------------------------------
//
//	ksignClaim := NewOperationalKSignClaim(common.HexToAddress("0xee602447b5a75cf4f25367f5d199b860844d10c4"))
//
//	assert.Nil(t, mt.Add(ksignClaim))
//
//	kroot := mt.Root()
//	kproof, err := mt.GenerateProof(ksignClaim.Hi())
//	assert.Nil(t, err)
//	assert.True(t, merkletree.CheckProof(kroot, kproof, ksignClaim.Hi(), ksignClaim.Ht(), 140))
//
//	assert.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074353f867ef725411de05e3d4b0a01c37cf7ad24bcc213141a0000005400000000ee602447b5a75cf4f25367f5d199b860844d10c4000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffff", common3.HexEncode(ksignClaim.Bytes()))
//	assert.Equal(t, uint32(84), ksignClaim.BaseIndex.IndexLength)
//	assert.Equal(t, 84, int(ksignClaim.IndexLength()))
//	assert.Equal(t, "0x68be938284f64944bd8ebc172792687f680fb8db13e383227c8c668820b40078", ksignClaim.Hi().Hex())
//	assert.Equal(t, "0xdd675b18734a480868ed7b258ec2306a8e676690a81d53bcda7490c31368edd2", ksignClaim.Ht().Hex())
//	assert.Equal(t, "0x93bf43768a1e034e583832a9ee992c37374047be910aa1e80258fc2f27d46628", kroot.Hex())
//	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000000", common3.HexEncode(kproof))
//
//	ksignClaim.BaseIndex.Version = 1
//	kproofneg, err := mt.GenerateProof(ksignClaim.Hi())
//	assert.Nil(t, err)
//	assert.Equal(t, "0xeab0608b8891dcca4f421c69244b17f208fbed899b540d01115ca7d907cbf6a5", ksignClaim.Hi().Hex())
//	assert.True(t, merkletree.CheckProof(kroot, kproofneg, ksignClaim.Hi(), merkletree.EmptyNodeValue, 140))
//	assert.Equal(t, "0x000000000000000000000000000000000000000000000000000000000000000103aab4f597fe23598cc10f1af68192195a7538d3d6fc83cf49e5cfd53eaac527", common3.HexEncode(kproofneg))
//
//	// create setrootclaim ----------------------------------------------
//
//	mt, err = merkletree.New(db.NewMemoryStorage(), 140)
//	assert.Nil(t, err)
//
//	setRootClaim := NewSetRootClaim(
//		common.HexToAddress("0xd79ae0a65e7dd29db1eac700368e693de09610b8"),
//		kroot,
//	)
//
//	assert.Nil(t, mt.Add(setRootClaim))
//
//	rroot := mt.Root()
//	rproof, err := mt.GenerateProof(setRootClaim.Hi())
//	assert.Nil(t, err)
//
//	assert.True(t, merkletree.CheckProof(rroot, rproof, setRootClaim.Hi(), setRootClaim.Ht(), 140))
//	assert.Equal(t, uint32(84), setRootClaim.BaseIndex.IndexLength)
//	assert.Equal(t, 84, int(setRootClaim.IndexLength()))
//	assert.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c49694030749b9a76a0132a0814192c05c9321efc30c7286f6187f18fc60000005400000000d79ae0a65e7dd29db1eac700368e693de09610b893bf43768a1e034e583832a9ee992c37374047be910aa1e80258fc2f27d46628", common3.HexEncode(setRootClaim.Bytes()))
//	assert.Equal(t, "0x497d8626567f90e3e14de025961133ca7e4959a686c75a062d4d4db750d607b0", setRootClaim.Hi().Hex())
//	assert.Equal(t, "0x6da033d96fdde2c687a48a4902823f9f8e91b31e3d73c57f3858e8a9650f9c39", setRootClaim.Ht().Hex())
//	assert.Equal(t, "0xab63a4a3c5fe879e1b55315b945ac7f1ac1ac4b059e7301964b99b6813b514c7", rroot.Hex())
//	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000000", common3.HexEncode(rproof))
//
//	setRootClaim.BaseIndex.Version++
//	rproofneg, err := mt.GenerateProof(setRootClaim.Hi())
//	assert.Nil(t, err)
//	assert.True(t, merkletree.CheckProof(rroot, rproofneg, setRootClaim.Hi(), merkletree.EmptyNodeValue, 140))
//	assert.Equal(t, "0x00000000000000000000000000000000000000000000000000000000000000016f33cf71ff7bdbc492f9c3bd63b15577e6cedc70afd09051e1dfe2f04340c073", common3.HexEncode(rproofneg))
//}
