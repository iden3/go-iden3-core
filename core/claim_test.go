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
	claimHex := "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074cfee7c08a98f4b565d124c7e4e28acc52e1bc780e3887db000000040000000006331"
	claimBytes, err := common3.HexToBytes(claimHex)
	assert.Nil(t, err)
	claimType, err := ParseTypeClaimBytes(claimBytes)
	assert.Nil(t, err)
	assert.Equal(t, "default", claimType)

	// assignNameClaim type
	claimHex = "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074b7ae3d3a2056c54f48763999f3ff99caffaaba3bab58cae900000080000000009c4b7a6b4af91b44be8d9bb66d41e82589f01974702d3bf1d9b4407a55593c3c3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074101d2fa51f8259df207115af9eaa73f3f4e52e60"
	claimBytes, err = common3.HexToBytes(claimHex)
	assert.Nil(t, err)
	claimType, err = ParseTypeClaimBytes(claimBytes)
	assert.Nil(t, err)
	assert.Equal(t, "assignname", claimType)

	// authorizeKSignClaim type
	claimHex = "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074353f867ef725411de05e3d4b0a01c37cf7ad24bcc213141a0000005400000000101d2fa51f8259df207115af9eaa73f3f4e52e602077bb3f0400dd62421c97220536fd6ed2be29228e8db1315e8c6d7525f4bdf4dad9966a2e7371f0a24b1929ed765c0e7a3f2b4665a76a19d58173308bb34062000000005b816b9e000000005b816b9e"
	claimBytes, err = common3.HexToBytes(claimHex)
	assert.Nil(t, err)
	claimType, err = ParseTypeClaimBytes(claimBytes)
	assert.Nil(t, err)
	assert.Equal(t, "authorizeksign", claimType)

	// setRootClaim type
	claimHex = "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c49694030749b9a76a0132a0814192c05c9321efc30c7286f6187f18fc60000005400000000101d2fa51f8259df207115af9eaa73f3f4e52e607d2e30055477edf346d81dac92720729aed3ef830cf26b012b38f2f8304ac18c"
	claimBytes, err = common3.HexToBytes(claimHex)
	assert.Nil(t, err)
	claimType, err = ParseTypeClaimBytes(claimBytes)
	assert.Nil(t, err)
	assert.Equal(t, "setroot", claimType)
}
func TestClaimGenerationAndParse(t *testing.T) {
	claim := NewClaimDefault("iden3.io", "default", []byte("c1"), []byte{})
	assert.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074cfee7c08a98f4b565d124c7e4e28acc52e1bc780e3887db000000042000000006331", common3.BytesToHex(claim.Bytes()))

	claimHt := claim.Ht()
	assert.Equal(t, "0x0fce11cbd33e15d137a3a1953cda71aa81898ee8b917c21615073b59cd4dca8c", common3.BytesToHex(claimHt[:]))

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
	assert.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074b7ae3d3a2056c54f48763999f3ff99caffaaba3bab58cae900000080000000009c4b7a6b4af91b44be8d9bb66d41e82589f01974702d3bf1d9b4407a55593c3c3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074101d2fa51f8259df207115af9eaa73f3f4e52e60", common3.BytesToHex(assignNameClaim.Bytes()))

	assert.Equal(t, "0xfa7c93d75fc4b617a8ab43afd1def1a06ffc0c227ce25edc93b235e4df084cf0", assignNameClaim.Hi().Hex())
	assert.Equal(t, "0xa09dafa7ebe716f4c35a316c4c332d47b2bb3c101402f9ba4467ed676719d45a", assignNameClaim.Ht().Hex())
}

func TestAuthorizeKSign(t *testing.T) {
	authorizeKSignClaim := NewAuthorizeKSignClaim("iden3.io", common.HexToAddress("0x101d2fa51f8259df207115af9eaa73f3f4e52e60"), "appToAuthName", "authz", 1535208350, 1535208350)
	authorizeKSignClaimParsed, err := ParseAuthorizeKSignClaimBytes(authorizeKSignClaim.Bytes())
	assert.Nil(t, err)
	if !bytes.Equal(authorizeKSignClaimParsed.Bytes(), authorizeKSignClaim.Bytes()) {
		t.Errorf("ksignClaim and ksignClaim parsed are not equals")
	}
	assert.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074353f867ef725411de05e3d4b0a01c37cf7ad24bcc213141a0000005400000000101d2fa51f8259df207115af9eaa73f3f4e52e602077bb3f0400dd62421c97220536fd6ed2be29228e8db1315e8c6d7525f4bdf4dad9966a2e7371f0a24b1929ed765c0e7a3f2b4665a76a19d58173308bb34062000000005b816b9e000000005b816b9e", common3.BytesToHex(authorizeKSignClaim.Bytes()))
	assert.Equal(t, "0xb98902d35fe0861daaeb78ada21e60e1d7c009b6e56d85127e892aeb4ed37ef2", authorizeKSignClaim.Hi().Hex())
	assert.Equal(t, "0x9a1d4978ced5adfd4c4de9ee1bb4f850e0db426855737a5ecf0749d150620422", authorizeKSignClaim.Ht().Hex())

}
func TestSetRootClaim(t *testing.T) {
	setRootClaim := NewSetRootClaim("iden3.io", common.HexToAddress("0x101d2fa51f8259df207115af9eaa73f3f4e52e60"), merkletree.HashBytes([]byte("root of the MT")))
	setRootClaimParsed, err := ParseSetRootClaimBytes(setRootClaim.Bytes())
	assert.Nil(t, err)
	if !bytes.Equal(setRootClaimParsed.Bytes(), setRootClaim.Bytes()) {
		t.Errorf("setRootClaim and setRootClaim parsed are not equals")
	}
	assert.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c49694030749b9a76a0132a0814192c05c9321efc30c7286f6187f18fc60000005400000000101d2fa51f8259df207115af9eaa73f3f4e52e607d2e30055477edf346d81dac92720729aed3ef830cf26b012b38f2f8304ac18c", common3.BytesToHex(setRootClaim.Bytes()))

	assert.Equal(t, "0x7e76f2faee2b80b1d794271fe06b6fa52b06bcabc48441473ae28aa281b19965", setRootClaim.Hi().Hex())
	assert.Equal(t, "0xaf49c9214bb28b886e197df5b5e38c0ea54b3334039f419375e0fb4d1f70e44c", setRootClaim.Ht().Hex())
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
	proof, err := mt.GenerateProof(ksignClaim.Hi())
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, merkletree.CheckProof(root, proof, ksignClaim.Hi(), ksignClaim.Ht(), 140))

	assert.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074353f867ef725411de05e3d4b0a01c37cf7ad24bcc213141a0000005400000000ee602447b5a75cf4f25367f5d199b860844d10c4d6f028ca0e8edb4a8c9757ca4fdccab25fa1e0317da1188108f7d2dee14902fbdad9966a2e7371f0a24b1929ed765c0e7a3f2b4665a76a19d58173308bb3406200000000259e9d8000000000967a7600", common3.BytesToHex(ksignClaim.Bytes()))
	assert.Equal(t, uint32(84), ksignClaim.BaseIndex.IndexLength)
	assert.Equal(t, 84, int(ksignClaim.IndexLength()))
	assert.Equal(t, uint32(0x54), ksignClaim.IndexLength())
	assert.Equal(t, "0x68be938284f64944bd8ebc172792687f680fb8db13e383227c8c668820b40078", ksignClaim.Hi().Hex())
	assert.Equal(t, "0x63b43ece0a9f5f63a4333143563896d6d4e8b0ce8acd8dd2e6f7aaec52a007bb", ksignClaim.Ht().Hex())
	assert.Equal(t, "0x532abdf4d17d806893915c6d04ebd669ea02f127bd0f48b897dabbac75764ed6", root.Hex())
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
	proof, err := mt.GenerateProof(setRootClaim.Hi())
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, merkletree.CheckProof(root, proof, setRootClaim.Hi(), setRootClaim.Ht(), 140))
	assert.Equal(t, uint32(84), setRootClaim.BaseIndex.IndexLength)
	assert.Equal(t, 84, int(setRootClaim.IndexLength()))
	assert.Equal(t, uint32(0x54), setRootClaim.IndexLength())
	assert.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c49694030749b9a76a0132a0814192c05c9321efc30c7286f6187f18fc60000005400000000d79ae0a65e7dd29db1eac700368e693de09610b8562c7589149679a8dce7c53c16475eb572ea4b75d23539132d3093b483b8f1a3", common3.BytesToHex(setRootClaim.Bytes()))
	assert.Equal(t, "0x497d8626567f90e3e14de025961133ca7e4959a686c75a062d4d4db750d607b0", setRootClaim.Hi().Hex())
	assert.Equal(t, "0x4920867cc5963579b7c919dbb8a1cf3164fdde0c9f06b2af3d6613b3346c7f9e", setRootClaim.Ht().Hex())
	assert.Equal(t, "0x585d9dcd51abce33f55f7be8ba04719aef308a2f9e4280593eaef981672be24c", root.Hex())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000000", common3.BytesToHex(proof))
}

// hexToBytes converts from a hex string into an array of bytes
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
