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
	"github.com/iden3/go-iden3/crypto/babyjub"
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
		"0x1d4d6c81f3cd8bd286affa0d5ac3b677d86fea34ba88d450081d703bcf712e6a",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x03c4686d099ffd137b83ba22b57dc954ac1e6c0e2b1e0ef972a936992b8788ff",
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
	// genesis := common.BytesToAddress([]byte{
	//         0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39,
	//         0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39,
	//         0x39, 0x39, 0x39, 0x3a})
	id, err := IDFromString("1pnWU7Jdr4yLxp1azs1r1PpvfErxKGRQdcLBZuq3Z")
	assert.Nil(t, err)
	c0 := NewClaimAssignName(name, id)
	c0.Version = 1
	e := c0.Entry()
	assert.Equal(t,
		"0x106d1a898d4503f4cb20be6ce9aeb2ac1e65d522579805e3633408a4b9ffcb53",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x25867e06233f276f39e298775245bad077eb0852b4eaac8dbf646a95bd3f8625",
		e.HValue().Hex())
	dataTestOutput(&e.Data)
	assert.Equal(t, ""+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"0000041c980d8faa54be797337fa55dbe62a7675e0c83ce5383b78a04b26b9f4"+
		"00d67b05d8e2d1ace8f3e84b8451dd2e9da151578c3c6be23e7af11add5a807a"+
		"0000000000000000000000000000000000000000000000010000000000000003",
		e.Data.String())
	c1 := NewClaimAssignNameFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)
}

func TestClaimAuthorizeKSignBabyJub(t *testing.T) {
	// ClaimAuthorizeKSignBabyJub
	var k babyjub.PrivateKey
	hex.Decode(k[:], []byte("28156abe7fe2fd433dc9df969286b96666489bac508612d0e16593e944c4f69f"))
	pk := k.Public()

	c0 := NewClaimAuthorizeKSignBabyJub(pk)
	c0.Version = 1
	e := c0.Entry()
	assert.Equal(t,
		"0x04f41fdac3240e7b68905df19a2394e4a4f1fb7eaeb310e39e1bb0b225b7763f",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x06d4571fb9634e4bed32e265f91a373a852c476656c5c13b09bc133ac61bc5a6",
		e.HValue().Hex())
	dataTestOutput(&e.Data)
	assert.Equal(t, ""+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"2b05184c7195b259c95169348434f3a7228fbcfb187d3b07649f3791330cf05c"+
		"0000000000000000000000000000000000000001000000010000000000000001",
		e.Data.String())
	c1 := NewClaimAuthorizeKSignBabyJubFromEntry(e)
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
		"0x25aacb66cedd3be6248f68d61e8648ba163333070a4da17d35c424b798248440",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x06d4571fb9634e4bed32e265f91a373a852c476656c5c13b09bc133ac61bc5a6",
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
	id, err := IDFromString("1pnWU7Jdr4yLxp1azs1r1PpvfErxKGRQdcLBZuq3Z")
	assert.Nil(t, err)

	rootKey := merkletree.Hash(merkletree.ElemBytes{
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0c})
	c0 := NewClaimSetRootKey(id, rootKey)
	c0.Version = 1
	c0.Era = 1
	e := c0.Entry()
	assert.Equal(t,
		"0x12bf59ff4171debe81321c04a52298e62650ca8514e9a7a8a64c23cb55eeaa2e",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x01705b25f2cf7cda34d836f09e9b0dd1777bdc16752657cd9d1ae5f6286525ba",
		e.HValue().Hex())
	dataTestOutput(&e.Data)
	assert.Equal(t, ""+
		"0000000000000000000000000000000000000000000000000000000000000000"+
		"0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0c"+
		"0000041c980d8faa54be797337fa55dbe62a7675e0c83ce5383b78a04b26b9f4"+
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
	const objectType = ObjectTypeAddress
	var indexType uint16
	id, err := IDFromString("1pnWU7Jdr4yLxp1azs1r1PpvfErxKGRQdcLBZuq3Z")
	assert.Nil(t, err)

	objectHash := []byte{
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b,
		0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0b, 0x0c}

	auxData := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x09,
		0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x01, 0x02}

	claim := NewClaimLinkObjectIdentity(objectType, indexType, id, objectHash, auxData)
	claim.Version = 1
	entry := claim.Entry()
	assert.Equal(t,
		"0x2dc73c37e603a15f8f028aa5c3f668d1210c86008577188ce279ead60a9afec4",
		entry.HIndex().Hex())
	assert.Equal(t,
		"0x0f55d2c10514bb5be610006cc9a1ff18aa4bb248856b41de516ee6d027b9463c",
		entry.HValue().Hex())
	dataTestOutput(&entry.Data)
	assert.Equal(t, ""+
		"000102030405060708090a0b0c0d0e0f01020304050607090a0b0c0d0e0f0102"+
		"000b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0c"+
		"0000041c980d8faa54be797337fa55dbe62a7675e0c83ce5383b78a04b26b9f4"+
		"0000000000000000000000000000000000000001000000010000000000000005",
		entry.Data.String())
	c1 := NewClaimLinkObjectIdentityFromEntry(entry)
	c2, err := NewClaimFromEntry(entry)
	assert.Nil(t, err)
	assert.Equal(t, claim, c1)
	assert.Equal(t, claim, c2)
}

func TestClaimAuthorizeService(t *testing.T) {
	// ClaimAuthorizeService
	ethAddr := common.BytesToAddress([]byte{
		0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39,
		0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39, 0x39,
		0x39, 0x39, 0x39, 0x3a})
	pubKstr := "af048ddcc131d526699d928e8b8548c5c85fb7d407fc408bb543e4e58f305347f67942a7e56d7dc90bbcecca865f2fbde3118c91516594262f62857136f71dbc"
	c0 := NewClaimAuthorizeService(ServiceTypeRelay, ethAddr.Hex(), pubKstr, "relay.iden3.io")
	e := c0.Entry()
	assert.Equal(t,
		"0x0ee7fb1c970abca8667607eca3974704783f8812bc7f745c1c7ee49a2faf7927",
		e.HIndex().Hex())
	assert.Equal(t,
		"0x2ac15c5d5a255d7d92d84580c9a19b2e6beed42cfd26978b448d6b4abfa6d017",
		e.HValue().Hex())
	dataTestOutput(&e.Data)
	assert.Equal(t, ""+
		"00f3b1c89978c483ef94f9ecff889cbef9db68036f3b2dc251e72b7960b8529d"+
		"00f28abb0b5b73fdcc8eed8e707f33d8dd9b50b3e2c6e1957a585903ae3b729a"+
		"00f54d900e54dfb5d19c0e19e5e3abca0d744fee18b72cb8b9cc05f655495983"+
		"0000000000000000000000000000000000000000000000000000000000000006",
		e.Data.String())
	c1 := NewClaimAuthorizeServiceFromEntry(e)
	c2, err := NewClaimFromEntry(e)
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)
	assert.Equal(t, c0.ServiceType, ServiceTypeRelay)
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

func TestClaimEthId(t *testing.T) {
	ethId := common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
	identityFactoryAddr := common.HexToAddress("0x66D0c2F85F1B717168cbB508AfD1c46e07227130")

	c0 := NewClaimEthId(ethId, identityFactoryAddr)

	c1 := NewClaimEthIdFromEntry(c0.Entry())
	c2, err := NewClaimFromEntry(c0.Entry())
	assert.Nil(t, err)
	assert.Equal(t, c0, c1)
	assert.Equal(t, c0, c2)

	assert.Equal(t, c0.Address, ethId)
	assert.Equal(t, c0.IdentityFactory, identityFactoryAddr)
	assert.Equal(t, c0.Address, c1.Address)
	assert.Equal(t, c0.IdentityFactory, c1.IdentityFactory)

	assert.Equal(t, c0.Entry().Bytes(), c1.Entry().Bytes())
	assert.Equal(t, c0.Entry().Bytes(), c2.Entry().Bytes())
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
