package merkletree

/*
This is just an example of basic tests for the iden3-merkletree-specification.
The methods and variables names can be different.
Complete specification in https://github.com/iden3/iden3-merkletree-specification
*/

import (
	"strconv"
	"testing"

	common3 "github.com/iden3/go-iden3/common"
	"github.com/stretchr/testify/assert"
)

type testBytesClaim struct {
	data        []byte
	indexLength uint32
}

func (c testBytesClaim) Bytes() (b []byte) {
	return c.data
}
func (c testBytesClaim) IndexLength() uint32 {
	return c.indexLength
}
func (c testBytesClaim) Hi() Hash {
	h := HashBytes(c.Bytes()[:c.IndexLength()])
	return h
}
func newTestBytesClaim(data string, indexLength uint32) testBytesClaim {
	return testBytesClaim{
		data:        []byte(data),
		indexLength: indexLength,
	}
}

// Test to check the iden3-merkletree-specification
func TestIden3MerkletreeSpecification(t *testing.T) {
	h := HashBytes([]byte("test")).Hex()
	assert.Equal(t, "0x9c22ff5f21f0b81b113e63f7db6da94fedef11b2119b4088b89664fb9a3cb658", h)

	h = HashBytes([]byte("authorizeksign")).Hex()
	assert.Equal(t, "0x353f867ef725411de05e3d4b0a01c37cf7ad24bcc213141a05ed7726d7932a1f", h)

	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	// empty tree
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000000", mt.Root().Hex())

	// add claim
	claim := testBytesClaim{
		data:        []byte("this is a test claim"),
		indexLength: 15,
	}
	assert.Nil(t, mt.Add(claim))
	assert.Equal(t, "0x1c4160fe7330f22ef5bd5f4eefc3a818a6dec63a5014600c83fe0ef8495e28ed", mt.Root().Hex())

	// proof with only one claim in the MerkleTree
	proof, err := mt.GenerateProof(claim.Hi())
	assert.Nil(t, err)
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000000", common3.BytesToHex(proof))

	// add a second claim
	claim2 := testBytesClaim{
		data:        []byte("this is a second test claim"),
		indexLength: 15,
	}
	err = mt.Add(claim2)
	assert.Nil(t, err)
	assert.Equal(t, "0xc85f08a5500320b7877bffec8298f5c222c260e6ba86968114d70f8591ccef3e", mt.Root().Hex())

	// proof of the second claim, with two claims in the MerkleTree
	proof2, err := mt.GenerateProof(claim2.Hi())
	assert.Nil(t, err)
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000001feedc5746452611b2d5fc83bbc72ebeb1e284c071e1552a1876ae7e1d5043946", common3.BytesToHex(proof2))

	// proof of emptyLeaf
	claim3 := testBytesClaim{
		data:        []byte("this is a third test claim"),
		indexLength: 15,
	}
	proof3, err := mt.GenerateProof(claim3.Hi())
	assert.Nil(t, err)
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000003c11c2813e3b6ab49fb0a1236bd6b0b150d06a9ddc04fbde23d3cb71f58ee9d7ffeedc5746452611b2d5fc83bbc72ebeb1e284c071e1552a1876ae7e1d5043946", common3.BytesToHex(proof3))

	// getClaimByHi/GetValueInPos
	bytesInHi, err := mt.GetValueInPos(claim2.Hi())
	assert.Nil(t, err)
	assert.Equal(t, claim2.Bytes(), bytesInHi)

	// check proof
	rootBytes, err := common3.HexToBytes("0x7d7c5e8f4b3bf434f3d9d223359c4415e2764dd38de2e025fbf986e976a7ed3d")
	assert.Nil(t, err)
	mp, err := common3.HexToBytes("0x0000000000000000000000000000000000000000000000000000000000000002d45aada6eec346222eaa6b5d3a9260e08c9b62fcf63c72bc05df284de07e6a52")
	assert.Nil(t, err)
	hiBytes, err := common3.HexToBytes("0x786677808ba77bdd9090a969f1ef2cbd1ac5aecd9e654f340500159219106878")
	assert.Nil(t, err)
	htBytes, err := common3.HexToBytes("0x786677808ba77bdd9090a969f1ef2cbd1ac5aecd9e654f340500159219106878")
	assert.Nil(t, err)
	var root, hi, ht Hash
	copy(root[:], rootBytes)
	copy(hi[:], hiBytes)
	copy(ht[:], htBytes)
	verified := CheckProof(root, mp, hi, ht, 140)
	assert.True(t, verified)

	// check proof of empty
	rootBytes, err = common3.HexToBytes("0x8f021d00c39dcd768974ddfe0d21f5d13f7215bea28db1f1cb29842b111332e7")
	assert.Nil(t, err)
	mp, err = common3.HexToBytes("0x0000000000000000000000000000000000000000000000000000000000000004bf8e980d2ed328ae97f65c30c25520aeb53ff837579e392ea1464934c7c1feb9")
	assert.Nil(t, err)
	hiBytes, err = common3.HexToBytes("0xa69792a4cff51f40b7a1f7ae596c6ded4aba241646a47538898f17f2a8dff647")
	assert.Nil(t, err)
	htBytes, err = common3.HexToBytes("0x0000000000000000000000000000000000000000000000000000000000000000")
	assert.Nil(t, err)
	copy(root[:], rootBytes)
	copy(hi[:], hiBytes)
	copy(ht[:], htBytes)
	verified = CheckProof(root, mp, hi, ht, 140)
	assert.True(t, verified)

	// check the proof generated in the previous steps
	verified = CheckProof(mt.Root(), proof2, claim2.Hi(), HashBytes(claim2.Bytes()), 140)
	assert.True(t, verified)
	// check proof of no existence (emptyLeaf), as we are prooving an empty leaf, the Ht is an empty value (0x000...0)
	verified = CheckProof(mt.Root(), proof3, claim3.Hi(), EmptyNodeValue, 140)
	assert.True(t, verified)

	// add claims in different orders
	mt1 := newTestingMerkle(t, 140)
	defer mt1.Storage().Close()

	mt1.Add(newTestBytesClaim("0 this is a test claim", 15))
	mt1.Add(newTestBytesClaim("1 this is a test claim", 15))
	mt1.Add(newTestBytesClaim("2 this is a test claim", 15))
	mt1.Add(newTestBytesClaim("3 this is a test claim", 15))
	mt1.Add(newTestBytesClaim("4 this is a test claim", 15))

	mt2 := newTestingMerkle(t, 140)
	defer mt2.Storage().Close()

	mt2.Add(newTestBytesClaim("2 this is a test claim", 15))
	mt2.Add(newTestBytesClaim("1 this is a test claim", 15))
	mt2.Add(newTestBytesClaim("0 this is a test claim", 15))
	mt2.Add(newTestBytesClaim("3 this is a test claim", 15))
	mt2.Add(newTestBytesClaim("4 this is a test claim", 15))

	assert.Equal(t, mt1.Root().Hex(), mt2.Root().Hex())

	// adding 1000 claims
	mt1000 := newTestingMerkle(t, 140)
	defer mt1000.Storage().Close()

	numToAdd := 1000
	for i := 0; i < numToAdd; i++ {
		claim := newTestBytesClaim(strconv.Itoa(i)+" this is a test claim", 15)
		mt1000.Add(claim)
	}
	assert.Equal(t, "0xf1f6e6380d311dd7742be1aaecc35e9d7218bf11218d9f5bf8f7497b00a830c9", mt1000.Root().Hex())
}
