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

func TestForwardingInterop(t *testing.T) {

	// address 0xee602447b5a75cf4f25367f5d199b860844d10c4
	// pvk     8A85AAA2A8CE0D24F66D3EAA7F9F501F34992BACA0FF942A8EDF7ECE6B91F713

	dir, err := ioutil.TempDir("", "db")
	assert.Nil(t, err)

	stobase, err := merkletree.NewLevelDbStorage(dir)
	assert.Nil(t, err)
	defer stobase.Close()

	sto := stobase.WithPrefix([]byte{1})

	mt, err := merkletree.New(sto, 140)
	assert.Nil(t, err)

	// create ksignclaim ----------------------------------------------

	ksignClaim := NewOperationalKSignClaim(
		"iden3.io",
		common.HexToAddress("0xee602447b5a75cf4f25367f5d199b860844d10c4"),
		631152000,  // 1990
		2524608000, // 2050
	)

	assert.Nil(t, mt.Add(ksignClaim))

	kroot := mt.Root()
	kproof, err := mt.GenerateProof(ksignClaim.Hi())
	assert.Nil(t, err)
	assert.True(t, merkletree.CheckProof(kroot, kproof, ksignClaim.Hi(), ksignClaim.Ht(), 140))

	assert.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c4969403074353f867ef725411de05e3d4b0a01c37cf7ad24bcc213141a0000005400000000ee602447b5a75cf4f25367f5d199b860844d10c40000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000259e9d8000000000967a7600", common3.BytesToHex(ksignClaim.Bytes()))
	assert.Equal(t, uint32(84), ksignClaim.BaseIndex.IndexLength)
	assert.Equal(t, 84, int(ksignClaim.IndexLength()))
	assert.Equal(t, "0x68be938284f64944bd8ebc172792687f680fb8db13e383227c8c668820b40078", ksignClaim.Hi().Hex())
	assert.Equal(t, "0xd440292f476cde9f575d8fed36fbf096f5d2986fb11190ce0ee5d0e448bf1113", ksignClaim.Ht().Hex())
	assert.Equal(t, "0x0c7fbb73b49a62b75c44cc0b8559a67af866bcd942fa3bc1e7888d43e2f186f2", kroot.Hex())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000000", common3.BytesToHex(kproof))

	ksignClaim.BaseIndex.Version = 1
	kproofneg, err := mt.GenerateProof(ksignClaim.Hi())
	assert.Nil(t, err)
	assert.Equal(t, "0xeab0608b8891dcca4f421c69244b17f208fbed899b540d01115ca7d907cbf6a5", ksignClaim.Hi().Hex())
	assert.True(t, merkletree.CheckProof(kroot, kproofneg, ksignClaim.Hi(), merkletree.EmptyNodeValue, 140))
	assert.Equal(t, "0x00000000000000000000000000000000000000000000000000000000000000017a0ec823c79c6d1756a29edbf52eb228a69c5435ead519eb96cdb2412927b865", common3.BytesToHex(kproofneg))

	// create setrootclaim ----------------------------------------------

	sto = stobase.WithPrefix([]byte{2})

	mt, err = merkletree.New(sto, 140)
	assert.Nil(t, err)

	setRootClaim := NewSetRootClaim(
		"iden3.io",
		common.HexToAddress("0xd79ae0a65e7dd29db1eac700368e693de09610b8"),
		kroot,
	)

	assert.Nil(t, mt.Add(setRootClaim))

	rroot := mt.Root()
	rproof, err := mt.GenerateProof(setRootClaim.Hi())
	assert.Nil(t, err)

	assert.True(t, merkletree.CheckProof(rroot, rproof, setRootClaim.Hi(), setRootClaim.Ht(), 140))
	assert.Equal(t, uint32(84), setRootClaim.BaseIndex.IndexLength)
	assert.Equal(t, 84, int(setRootClaim.IndexLength()))
	assert.Equal(t, "0x3cfc3a1edbf691316fec9b75970fbfb2b0e8d8edfc6ec7628db77c49694030749b9a76a0132a0814192c05c9321efc30c7286f6187f18fc60000005400000000d79ae0a65e7dd29db1eac700368e693de09610b80c7fbb73b49a62b75c44cc0b8559a67af866bcd942fa3bc1e7888d43e2f186f2", common3.BytesToHex(setRootClaim.Bytes()))
	assert.Equal(t, "0x497d8626567f90e3e14de025961133ca7e4959a686c75a062d4d4db750d607b0", setRootClaim.Hi().Hex())
	assert.Equal(t, "0xb4f391a7eb28eb66adf447fb16da9d25408806a4f9154ffb7c6b13bb1f2bfd79", setRootClaim.Ht().Hex())
	assert.Equal(t, "0xa392bc7458973721c1266b2ac65db038a87bb6ad2e822c2509298803e9941119", rroot.Hex())
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000000", common3.BytesToHex(rproof))

	setRootClaim.BaseIndex.Version++
	rproofneg, err := mt.GenerateProof(setRootClaim.Hi())
	assert.Nil(t, err)
	assert.True(t, merkletree.CheckProof(rroot, rproofneg, setRootClaim.Hi(), merkletree.EmptyNodeValue, 140))
	assert.Equal(t, "0x00000000000000000000000000000000000000000000000000000000000000016602097464f2c4a8f7854f1c29a7671a85d5aa670dbbe04a65f9d9c50a70626d", common3.BytesToHex(rproofneg))
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
