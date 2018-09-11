package core

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/stretchr/testify/assert"
)

func TestParseTypeClaimBytes(t *testing.T) {
	// default claim
	claimHex := "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074cfee7c08a98f4b565d124c7e4e28acc52e1bc780e3887db0a02a7d2d5bc66728000000006331"
	claimBytes, err := common3.HexToBytes(claimHex)
	assert.Nil(t, err)
	claimType, err := ParseTypeClaimBytes(claimBytes)
	assert.Nil(t, err)
	assert.Equal(t, "default", claimType)

	// assignNameClaim type
	claimHex = "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074b7ae3d3a2056c54f48763999f3ff99caffaaba3bab58cae9f22abc828264ab70000000009c4b7a6b4af91b44be8d9bb66d41e82589f01974702d3bf1d9b4407a55593c3c3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074101d2fa51f8259df207115af9eaa73f3f4e52e60"
	claimBytes, err = common3.HexToBytes(claimHex)
	assert.Nil(t, err)
	claimType, err = ParseTypeClaimBytes(claimBytes)
	assert.Nil(t, err)
	assert.Equal(t, "assignname", claimType)

	// authorizeKSignClaim type
	claimHex = "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074353f867ef725411de05e3d4b0a01c37cf7ad24bcc213141a05ed7726d7932a1f00000000101d2fa51f8259df207115af9eaa73f3f4e52e602077bb3f0400dd62421c97220536fd6ed2be29228e8db1315e8c6d7525f4bdf4dad9966a2e7371f0a24b1929ed765c0e7a3f2b4665a76a19d58173308bb340623135333534363739333431353335343637393334"
	claimBytes, err = common3.HexToBytes(claimHex)
	assert.Nil(t, err)
	claimType, err = ParseTypeClaimBytes(claimBytes)
	assert.Nil(t, err)
	assert.Equal(t, "authorizeksign", claimType)

	// setRootClaim type
	claimHex = "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c49694030749b9a76a0132a0814192c05c9321efc30c7286f6187f18fc6b6858214fe963e0e00000000101d2fa51f8259df207115af9eaa73f3f4e52e607d2e30055477edf346d81dac92720729aed3ef830cf26b012b38f2f8304ac18c"
	claimBytes, err = common3.HexToBytes(claimHex)
	assert.Nil(t, err)
	claimType, err = ParseTypeClaimBytes(claimBytes)
	assert.Nil(t, err)
	assert.Equal(t, "setroot", claimType)
}
func TestClaimGenerationAndParse(t *testing.T) {
	claim := NewClaimDefault("iden3.io", "default", []byte("c1"))
	assert.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074cfee7c08a98f4b565d124c7e4e28acc52e1bc780e3887db0a02a7d2d5bc66728000000006331", common3.BytesToHex(claim.Bytes()))

	claimHt := claim.Ht()
	assert.Equal(t, "0x54f22c228c99b424f787ad0673782f44cf83d13cb13c74924b0746e120135e4b", common3.BytesToHex(claimHt[:]))

	claimParsed, err := ParseClaimDefaultBytes(claim.Bytes())
	assert.Nil(t, err)
	if !bytes.Equal(claim.Bytes(), claimParsed.Bytes()) {
		t.Errorf("claim and claimParsed not equal")
	}
}

func TestAssignNameClaim(t *testing.T) {
	assignNameClaim := NewAssignNameClaim("iden3.io", merkletree.HashBytes([]byte("john")), merkletree.HashBytes([]byte("iden3.io")), common.HexToAddress("0x101d2fa51f8259df207115af9eaa73f3f4e52e60"))
	assignNameClaimParsed, err := ParseAssignNameClaimBytes(assignNameClaim.Bytes())
	assert.Nil(t, err)
	if !bytes.Equal(assignNameClaimParsed.Bytes(), assignNameClaim.Bytes()) {
		t.Errorf("assignNameClaim and assignNameClaim parsed are not equals")
	}
	assert.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074b7ae3d3a2056c54f48763999f3ff99caffaaba3bab58cae9f22abc828264ab70000000009c4b7a6b4af91b44be8d9bb66d41e82589f01974702d3bf1d9b4407a55593c3c3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074101d2fa51f8259df207115af9eaa73f3f4e52e60", common3.BytesToHex(assignNameClaim.Bytes()))

	assert.Equal(t, "0x635c15786ced9f7b29aa0abdbcb7ab0b18d207a3ba79a8061488308df4cec8c2", assignNameClaim.Hi().Hex())
	assert.Equal(t, "0x85c181e833edd9e12fb1df71c3996bbcb3e26791cf4b39e0d47c9dec6428aa02", assignNameClaim.Ht().Hex())
}

func TestAuthorizeKSign(t *testing.T) {
	authorizeKSignClaim := NewAuthorizeKSignClaim("iden3.io", common.HexToAddress("0x101d2fa51f8259df207115af9eaa73f3f4e52e60"), "appToAuthName", "authz", 1535208350, 1535208350)
	authorizeKSignClaimParsed, err := ParseAuthorizeKSignClaimBytes(authorizeKSignClaim.Bytes())
	assert.Nil(t, err)
	if !bytes.Equal(authorizeKSignClaimParsed.Bytes(), authorizeKSignClaim.Bytes()) {
		t.Errorf("ksignClaim and ksignClaim parsed are not equals")
	}
	assert.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074353f867ef725411de05e3d4b0a01c37cf7ad24bcc213141a05ed7726d7932a1f00000000101d2fa51f8259df207115af9eaa73f3f4e52e602077bb3f0400dd62421c97220536fd6ed2be29228e8db1315e8c6d7525f4bdf4dad9966a2e7371f0a24b1929ed765c0e7a3f2b4665a76a19d58173308bb34062000000005b816b9e000000005b816b9e", common3.BytesToHex(authorizeKSignClaim.Bytes()))
	assert.Equal(t, "0xf94b1fbc765c2925f59fc266861d7e585eda11f804340769b62b41e1d7df9e89", authorizeKSignClaim.Hi().Hex())
	assert.Equal(t, "0x8282b432d409ef113ad443bb8ec95e5accdeb63bba0bcd10d1bc9d1155952fe0", authorizeKSignClaim.Ht().Hex())

}
func TestSetRootClaim(t *testing.T) {
	setRootClaim := NewSetRootClaim("iden3.io", common.HexToAddress("0x101d2fa51f8259df207115af9eaa73f3f4e52e60"), merkletree.HashBytes([]byte("root of the MT")))
	setRootClaimParsed, err := ParseSetRootClaimBytes(setRootClaim.Bytes())
	assert.Nil(t, err)
	if !bytes.Equal(setRootClaimParsed.Bytes(), setRootClaim.Bytes()) {
		t.Errorf("setRootClaim and setRootClaim parsed are not equals")
	}
	assert.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c49694030749b9a76a0132a0814192c05c9321efc30c7286f6187f18fc6b6858214fe963e0e00000000101d2fa51f8259df207115af9eaa73f3f4e52e607d2e30055477edf346d81dac92720729aed3ef830cf26b012b38f2f8304ac18c", common3.BytesToHex(setRootClaim.Bytes()))

	assert.Equal(t, "0x988f9b94f3904b1b2e18f821068e517e7c2ba9fbf401086c79ecc55c84fb4475", setRootClaim.Hi().Hex())
	assert.Equal(t, "0x9dbb985810faf851379c8c66f00d62fc9f346b8c3269774820bdd833b02a1cea", setRootClaim.Ht().Hex())
}

func TestKSignClaimInterop(t *testing.T) {

	// address 0xee602447b5a75cf4f25367f5d199b860844d10c4
	// pvk     8A85AAA2A8CE0D24F66D3EAA7F9F501F34992BACA0FF942A8EDF7ECE6B91F713

	dir, err := ioutil.TempDir("", "db")
	if err != nil {
		t.Fatal(err)
	}
	sto, err := merkletree.NewLevelDbStorage(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer sto.Close()

	mt, err := merkletree.New(sto, 140)
	if err != nil {
		t.Fatal(err)
	}

	ksignClaim := NewAuthorizeKSignClaim(
		"iden3.io",
		common.HexToAddress("0xee602447b5a75cf4f25367f5d199b860844d10c4"),
		"app", "authz",
		631152000,  // 1990
		2524608000, // 2050
	)

	if err = mt.Add(ksignClaim); err != nil {
		t.Fatal(err)
	}
	root := mt.Root()
	proof, err := mt.GenerateProof(ksignClaim)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, merkletree.CheckProof(root, proof, ksignClaim, 140))

	assert.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074353f867ef725411de05e3d4b0a01c37cf7ad24bcc213141a05ed7726d7932a1f00000000ee602447b5a75cf4f25367f5d199b860844d10c4d6f028ca0e8edb4a8c9757ca4fdccab25fa1e0317da1188108f7d2dee14902fbdad9966a2e7371f0a24b1929ed765c0e7a3f2b4665a76a19d58173308bb3406200000000259e9d8000000000967a7600", common3.BytesToHex(ksignClaim.Bytes()))
	assert.Equal(t, uint32(0x58), ksignClaim.IndexLength())
	assert.Equal(t, "0xefaf444c30354019722a8da1b5a1eca8fd4ff454aff3bfd477c8eb4ce05e75f0", ksignClaim.Hi().Hex())
	assert.Equal(t, "0xc98ce0dbbf4cd1fc05f2093b2ebb8b2fc4699cb4cde2b8e4c0a37f957c72e64f", ksignClaim.Ht().Hex())
	assert.Equal(t, "0x562c7589149679a8dce7c53c16475eb572ea4b75d23539132d3093b483b8f1a3", root.Hex())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000000", common3.BytesToHex(proof))
}

func TestSetRootClaimInterop(t *testing.T) {

	dir, err := ioutil.TempDir("", "db")
	if err != nil {
		t.Fatal(err)
	}
	sto, err := merkletree.NewLevelDbStorage(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer sto.Close()

	mt, err := merkletree.New(sto, 140)
	if err != nil {
		t.Fatal(err)
	}

	setRootClaim := NewSetRootClaim(
		"iden3.io",
		common.HexToAddress("0xd79ae0a65e7dd29db1eac700368e693de09610b8"),
		hexToHash("0x562c7589149679a8dce7c53c16475eb572ea4b75d23539132d3093b483b8f1a3"),
	)

	if err = mt.Add(setRootClaim); err != nil {
		t.Fatal(err)
	}
	root := mt.Root()
	proof, err := mt.GenerateProof(setRootClaim)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, merkletree.CheckProof(root, proof, setRootClaim, 140))
	assert.Equal(t, uint32(0x58), setRootClaim.IndexLength())
	assert.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c49694030749b9a76a0132a0814192c05c9321efc30c7286f6187f18fc6b6858214fe963e0e00000000d79ae0a65e7dd29db1eac700368e693de09610b8562c7589149679a8dce7c53c16475eb572ea4b75d23539132d3093b483b8f1a3", common3.BytesToHex(setRootClaim.Bytes()))
	assert.Equal(t, "0xaaad7b30f89608e270551f207688ea8f112bb3416e4ca07018d9e80bb05f26a8", setRootClaim.Hi().Hex())
	assert.Equal(t, "0xfa9cd70ad96d731f5d24d38baeba7a6a8d89c6910bdc430793b870a39d2f81d7", setRootClaim.Ht().Hex())
	assert.Equal(t, "0x1dce20a20a0f93a139de6069dcfb16b91f0a7d3a540eee0a57d1fa78c2f401c3", root.Hex())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000000", common3.BytesToHex(proof))
}

// HexToBytes converts from a hex string into an array of bytes
func hexToHash(hexstr string) merkletree.Hash {
	var h merkletree.Hash
	b, err := hex.DecodeString(hexstr[2:])
	if err != nil {
		panic(err)
	}
	if len(b) != len(h) {
		panic("Invalid hash length")
	}
	copy(h[:], b[:])
	return h
}
