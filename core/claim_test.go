package core

import (
	//"bytes"
	//"encoding/hex"
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
	indexSlot := [400 / 8]byte{42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42, 42}
	dataSlot := [496 / 8]byte{88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88, 88}
	c0 := NewClaimBasic(indexSlot, dataSlot)
	e := c0.ToEntry()
	assert.Equal(t, e.Data.String(),
		"00585858585858585858585858585858585858585858585858585858585858580058585858585858585858585858585858585858585858585858585858585858002a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a002a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a0000000080db04d364e4c1aa")
	c1 := NewClaimBasicFromEntry(&e)
	c2, err := NewClaimFromEntry(&e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, &c0, c2)
}

func TestClaimAssignName(t *testing.T) {
	// ClaimAssignName
	name := "example.iden3.eth"
	ethID := common.BytesToAddress([]byte{71, 71, 71, 71, 71, 71, 71, 71, 71, 71, 71, 71, 71, 71, 71, 71, 71, 71, 71, 71})
	c0 := NewClaimAssignName(name, ethID)
	e := c0.ToEntry()
	assert.Equal(t, e.Data.String(),
		"0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000474747474747474747474747474747474747474700805add1af17a3ee26b3c8c5751a19d2edd51844be8f3e8acd1e2d8057bd6480000000000000000000000000000000000000000000000005fd72e7912afa9ff")
	c1 := NewClaimAssignNameFromEntry(&e)
	c2, err := NewClaimFromEntry(&e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, &c0, c2)
}

func TestClaimAuthorizeKSign(t *testing.T) {
	// ClaimAuthorizeKSign
	sign := true
	ax := [128 / 8]byte{25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25}
	ay := [128 / 8]byte{77, 77, 77, 77, 77, 77, 77, 77, 77, 77, 77, 77, 77, 77, 77, 77}
	c0 := NewClaimAuthorizeKSign(sign, ax, ay)
	e := c0.ToEntry()
	assert.Equal(t, e.Data.String(),
		"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d000000191919191919191919191919191919190100000000e75c97ad68ba5cc8")
	c1 := NewClaimAuthorizeKSignFromEntry(&e)
	c2, err := NewClaimFromEntry(&e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, &c0, c2)
}

func TestClaimSetRootKey(t *testing.T) {
	// ClaimSetRootKey
	ethID := common.BytesToAddress([]byte{57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57, 57})
	rootKey := merkletree.Hash(merkletree.ElemBytes{00, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11})
	c0 := NewClaimSetRootKey(ethID, rootKey)
	e := c0.ToEntry()
	assert.Equal(t, e.Data.String(),
		"0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b0b0b0b0b0b0b0b0b0b0b0b0b0b0b000000000000000000000000003939393939393939393939393939393900000000000000000000000000000000000000000000000000000000e400a1345fb8a750")
	c1 := NewClaimSetRootKeyFromEntry(&e)
	c2, err := NewClaimFromEntry(&e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, &c0, c2)
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
