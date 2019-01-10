package core

import (
	//"bytes"
	//"encoding/hex"
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	//common3 "github.com/iden3/go-iden3/common"
	//"github.com/iden3/go-iden3/db"
	"github.com/iden3/go-iden3/merkletree"
	//"github.com/iden3/go-iden3/utils"
	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(t, "0x0d1770cf7af29da78eb31086bfa35a5945f39a8c4fa35edee71ac12a75b4a30b", e.HIndex().Hex())
	assert.Equal(t, "0x14869ce50566e440424a2571816b117d88a2e5e3d10a0abb7f89a89032b9e07f", e.HValue().Hex())
	assert.Equal(t, ""+
		"0056585858585858585858585858585858585858585858585858585858585858"+
		"0058585858585858585858585858585858585858585858585858585858585859"+
		"00292a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a"+
		"002a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2b000000015735944c6eb8f12d",
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
	ethID := common.BytesToAddress([]byte{
		0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39,
		0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39,
		0x39, 0x39, 0x39, 0x3a})
	c0 := NewClaimAssignName(name, ethID)
	c0.Version = 1
	e := c0.Entry()
	assert.Equal(t, "0x23966b07b31bad5aebd8af6c72c7650f8ab45886e442f427da6c1bce73dbd2bb", e.HIndex().Hex())
	assert.Equal(t, "0x279689e54ed1540614ba9ca682a01e83eb8b6aa3abf85b1f659fd537a75c5d6a", e.HValue().Hex())
	assert.Equal(t, ""+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"000000000000000000000000393939393939393939393939393939393939393a"+
		"00d67b05d8e2d1ace8f3e84b8451dd2e9da151578c3c6be23e7af11add5a807a"+
		"000000000000000000000000000000000000000000000001f60d928459d792ed",
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
	ax := [128 / 8]byte{
		0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05,
		0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x06}
	ay := [128 / 8]byte{
		0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07,
		0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x08}
	c0 := NewClaimAuthorizeKSign(sign, ax, ay)
	c0.Version = 1
	e := c0.Entry()
	assert.Equal(t, "0x18f1032141d6a2abda87e2cf053edcffb5be55ba0dc4c5a9073805c4aa7aee54", e.HIndex().Hex())
	assert.Equal(t, "0x15331daa10ae035babcaabb76a80198bc449d32240ebb7f456ff2b03cd69bca4", e.HValue().Hex())
	assert.Equal(t, ""+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"0000000000000000000000000000000007070707070707070707070707070708"+
		"0000000505050505050505050505050505050601000000015714c3724876e56d",
		e.Data.String())
	c1 := NewClaimAuthorizeKSignFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)
}

func TestClaimSetRootKey(t *testing.T) {
	// ClaimSetRootKey
	ethID := common.BytesToAddress([]byte{
		0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39,
		0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39,
		0x39, 0x39, 0x39, 0x3a})
	rootKey := merkletree.Hash(merkletree.ElemBytes{
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0c})
	c0 := NewClaimSetRootKey(ethID, rootKey)
	c0.Version = 1
	c0.Era = 1
	e := c0.Entry()
	assert.Equal(t, "0x0a2d38687ea5f987637ded13030b22d1657be60bdb35add74bb53c8d5d126f8f", e.HIndex().Hex())
	assert.Equal(t, "0x2e27903d404fcab9363967a4ffe7da6a615f9ce6f55c43661a0297a040d336a4", e.HValue().Hex())
	assert.Equal(t, ""+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0c"+
		"000000000000000000000000393939393939393939393939393939393939393a"+
		"000000000000000000000000000000000000000100000001b111df93ad32c22c",
		e.Data.String())
	c1 := NewClaimSetRootKeyFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)
}

func dataTestOutput(d *merkletree.Data) string {
	s := bytes.NewBufferString("")
	fmt.Fprintf(s, "\"%v\"+\n", hex.EncodeToString(d[0][:]))
	fmt.Fprintf(s, "\"%v\"+\n", hex.EncodeToString(d[1][:]))
	fmt.Fprintf(s, "\"%v\"+\n", hex.EncodeToString(d[2][:]))
	fmt.Fprintf(s, "\"%v\",", hex.EncodeToString(d[3][:]))
	return s.String()
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
//	assert.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074353f867ef725411de05e3d4b0a01c37cf7ad24bcc213141a0000005400000000ee602447b5a75cf4f25367f5d199b860844d10c4000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000ffffffffffffffff", common3.BytesToHex(ksignClaim.Bytes()))
//	assert.Equal(t, uint32(84), ksignClaim.BaseIndex.IndexLength)
//	assert.Equal(t, 84, int(ksignClaim.IndexLength()))
//	assert.Equal(t, "0x68be938284f64944bd8ebc172792687f680fb8db13e383227c8c668820b40078", ksignClaim.Hi().Hex())
//	assert.Equal(t, "0xdd675b18734a480868ed7b258ec2306a8e676690a81d53bcda7490c31368edd2", ksignClaim.Ht().Hex())
//	assert.Equal(t, "0x93bf43768a1e034e583832a9ee992c37374047be910aa1e80258fc2f27d46628", kroot.Hex())
//	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000000", common3.BytesToHex(kproof))
//
//	ksignClaim.BaseIndex.Version = 1
//	kproofneg, err := mt.GenerateProof(ksignClaim.Hi())
//	assert.Nil(t, err)
//	assert.Equal(t, "0xeab0608b8891dcca4f421c69244b17f208fbed899b540d01115ca7d907cbf6a5", ksignClaim.Hi().Hex())
//	assert.True(t, merkletree.CheckProof(kroot, kproofneg, ksignClaim.Hi(), merkletree.EmptyNodeValue, 140))
//	assert.Equal(t, "0x000000000000000000000000000000000000000000000000000000000000000103aab4f597fe23598cc10f1af68192195a7538d3d6fc83cf49e5cfd53eaac527", common3.BytesToHex(kproofneg))
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
//	assert.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c49694030749b9a76a0132a0814192c05c9321efc30c7286f6187f18fc60000005400000000d79ae0a65e7dd29db1eac700368e693de09610b893bf43768a1e034e583832a9ee992c37374047be910aa1e80258fc2f27d46628", common3.BytesToHex(setRootClaim.Bytes()))
//	assert.Equal(t, "0x497d8626567f90e3e14de025961133ca7e4959a686c75a062d4d4db750d607b0", setRootClaim.Hi().Hex())
//	assert.Equal(t, "0x6da033d96fdde2c687a48a4902823f9f8e91b31e3d73c57f3858e8a9650f9c39", setRootClaim.Ht().Hex())
//	assert.Equal(t, "0xab63a4a3c5fe879e1b55315b945ac7f1ac1ac4b059e7301964b99b6813b514c7", rroot.Hex())
//	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000000", common3.BytesToHex(rproof))
//
//	setRootClaim.BaseIndex.Version++
//	rproofneg, err := mt.GenerateProof(setRootClaim.Hi())
//	assert.Nil(t, err)
//	assert.True(t, merkletree.CheckProof(rroot, rproofneg, setRootClaim.Hi(), merkletree.EmptyNodeValue, 140))
//	assert.Equal(t, "0x00000000000000000000000000000000000000000000000000000000000000016f33cf71ff7bdbc492f9c3bd63b15577e6cedc70afd09051e1dfe2f04340c073", common3.BytesToHex(rproofneg))
//}
