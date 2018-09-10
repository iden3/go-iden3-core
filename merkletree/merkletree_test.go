package merkletree

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"
	"time"

	common3 "github.com/iden3/go-iden3/common"
	"github.com/stretchr/testify/assert"
)

type testBase struct {
	Length    [4]byte
	Namespace Hash
	Type      Hash
	Version   uint32
}
type testClaim struct {
	testBase
	extraIndex struct {
		Data []byte
	}
}

func parseTestClaimBytes(b []byte) testClaim {
	var c testClaim
	copy(c.testBase.Length[:], b[0:4])
	copy(c.testBase.Namespace[:], b[4:36])
	copy(c.testBase.Type[:], b[36:68])
	versionBytes := b[68:72]
	c.testBase.Version = common3.BytesToUint32(versionBytes)
	c.extraIndex.Data = b[72:]
	return c
}
func (c testClaim) Bytes() (b []byte) {
	b = append(b, c.testBase.Length[:]...)
	b = append(b, c.testBase.Namespace[:]...)
	b = append(b, c.testBase.Type[:]...)
	versionBytes, _ := common3.Uint32ToBytes(c.testBase.Version)
	b = append(b, versionBytes[:]...)
	b = append(b, c.extraIndex.Data[:]...)
	return b
}
func (c testClaim) IndexLength() uint32 {
	return uint32(len(c.Bytes()))
}
func (c testClaim) hi() Hash {
	h := HashBytes(c.Bytes())
	return h
}
func newTestClaim(namespaceStr, typeStr string, data []byte) testClaim {
	var c testClaim
	c.testBase.Length = [4]byte{0x00, 0x00, 0x00, 0x48}
	c.testBase.Namespace = HashBytes([]byte(namespaceStr))
	c.testBase.Type = HashBytes([]byte(typeStr))
	c.testBase.Version = 0
	c.extraIndex.Data = data
	return c
}

type Fatalable interface {
	Fatal(args ...interface{})
}

func newTestingMerkle(f Fatalable, numLevels int) *MerkleTree {
	dir, err := ioutil.TempDir("", "db")
	if err != nil {
		f.Fatal(err)
		return nil
	}
	sto, err := NewLevelDbStorage(dir)
	if err != nil {
		f.Fatal(err)
		return nil
	}
	mt, err := New(sto, numLevels)
	if err != nil {
		f.Fatal(err)
		return nil
	}
	return mt
}

func TestNewMT(t *testing.T) {

	//create a new MT
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000000", mt.Root().Hex())
}

func TestAddClaim(t *testing.T) {

	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	claim := newTestClaim("iden3.io", "typespec", []byte("c1"))
	assert.Equal(t, "0x939862c94ca9772fc9e2621df47128b1d4041b514e19edc969a92d8f0dae558f", claim.hi().Hex())

	assert.Nil(t, mt.Add(claim))
	assert.Equal(t, "0x9d3c407ff02c813cd474c0a6366b4f7c58bf417a38268f7a0d73a8bca2490b9b", mt.Root().Hex())
}

func TestAddClaims(t *testing.T) {

	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	claim := newTestClaim("iden3.io", "typespec", []byte("c1"))
	assert.Equal(t, "0x939862c94ca9772fc9e2621df47128b1d4041b514e19edc969a92d8f0dae558f", claim.hi().Hex())
	assert.Nil(t, mt.Add(claim))

	assert.Nil(t, mt.Add(newTestClaim("iden3.io2", "typespec2", []byte("c2"))))
	assert.Equal(t, "0xebae8fb483b48ba6c337136535198eb8bcf891daba40ac81e28958c09b9b229b", mt.Root().Hex())

	mt.Add(newTestClaim("iden3.io3", "typespec3", []byte("c3")))
	mt.Add(newTestClaim("iden3.io4", "typespec4", []byte("c4")))
	assert.Equal(t, "0xb4b51aa0c77a8e5ed0a099d7c11c7d2a9219ef241da84f0689da1f40a5f6ac31", mt.Root().Hex())
}

func TestAddClaimsCollision(t *testing.T) {

	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	claim := newTestClaim("iden3.io", "typespec", []byte("c1"))
	assert.Nil(t, mt.Add(claim))

	root1 := mt.Root()
	assert.EqualError(t, mt.Add(claim), ErrNodeAlreadyExists.Error())

	assert.Equal(t, root1.Hex(), mt.Root().Hex())
}

func TestAddClaimsDifferentOrders(t *testing.T) {

	mt1 := newTestingMerkle(t, 140)
	defer mt1.Storage().Close()

	mt1.Add(newTestClaim("iden3.io", "typespec", []byte("c1")))
	mt1.Add(newTestClaim("iden3.io2", "typespec2", []byte("c2")))
	mt1.Add(newTestClaim("iden3.io3", "typespec3", []byte("c3")))
	mt1.Add(newTestClaim("iden3.io4", "typespec4", []byte("c4")))
	mt1.Add(newTestClaim("iden3.io5", "typespec5", []byte("c5")))

	mt2 := newTestingMerkle(t, 140)
	defer mt2.Storage().Close()

	mt2.Add(newTestClaim("iden3.io3", "typespec3", []byte("c3")))
	mt2.Add(newTestClaim("iden3.io2", "typespec2", []byte("c2")))
	mt2.Add(newTestClaim("iden3.io", "typespec", []byte("c1")))
	mt2.Add(newTestClaim("iden3.io4", "typespec4", []byte("c4")))
	mt2.Add(newTestClaim("iden3.io5", "typespec5", []byte("c5")))

	assert.Equal(t, mt1.Root().Hex(), mt2.Root().Hex())
}
func TestBenchmarkAddingClaims(t *testing.T) {

	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	start := time.Now()
	numToAdd := 1000
	for i := 0; i < numToAdd; i++ {
		claim := newTestClaim("iden3.io"+strconv.Itoa(i), "typespec"+strconv.Itoa(i), []byte("c"+strconv.Itoa(i)))
		mt.Add(claim)
	}
	fmt.Print("time elapsed adding " + strconv.Itoa(numToAdd) + " claims: ")
	fmt.Println(time.Since(start))
}
func BenchmarkAddingClaims(b *testing.B) {

	mt := newTestingMerkle(b, 140)
	defer mt.Storage().Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		claim := newTestClaim("iden3.io"+strconv.Itoa(i), "typespec"+strconv.Itoa(i), []byte("c"+strconv.Itoa(i)))
		if err := mt.Add(claim); err != nil {
			b.Fatal(err)
		}
	}
}

func TestGenerateProof(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	mt.Add(newTestClaim("iden3.io_3", "typespec_3", []byte("c3")))
	mt.Add(newTestClaim("iden3.io_2", "typespec_2", []byte("c2")))

	claim1 := newTestClaim("iden3.io_1", "typespec_1", []byte("c1"))
	assert.Nil(t, mt.Add(claim1))

	mp, err := mt.GenerateProof(parseTestClaimBytes(claim1.Bytes()))
	assert.Nil(t, err)

	mpHexExpected := "0000000000000000000000000000000000000000000000000000000000000002beb0fd6dcf18d37fe51cf34beacd4c524d9c039ef9da2a27ccd3e7edf662c39c"
	assert.Equal(t, mpHexExpected, hex.EncodeToString(mp))
}

func TestCheckProof(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	claim1 := newTestClaim("iden3.io_1", "typespec_1", []byte("c1"))
	assert.Nil(t, mt.Add(claim1))

	claim3 := newTestClaim("iden3.io_3", "typespec_3", []byte("c3"))
	assert.Nil(t, mt.Add(claim3))

	mp, err := mt.GenerateProof(parseTestClaimBytes(claim1.Bytes()))
	assert.Nil(t, err)
	verified := CheckProof(mt.Root(), mp, parseTestClaimBytes(claim1.Bytes()), mt.NumLevels())
	assert.True(t, verified)

}

func TestGetClaimInPos(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	for i := 0; i < 50; i++ {
		claim := newTestClaim("iden3.io"+strconv.Itoa(i), "typespec"+strconv.Itoa(i), []byte("c"+strconv.Itoa(i)))
		mt.Add(claim)

	}
	claim1 := newTestClaim("iden3.io_x", "typespec_x", []byte("cx"))
	assert.Nil(t, mt.Add(claim1))

	claim := parseTestClaimBytes(claim1.Bytes())
	claimInPosBytes, err := mt.GetValueInPos(claim.hi())
	assert.Nil(t, err)
	assert.Equal(t, claim1.Bytes(), claimInPosBytes)
}

type vt struct {
	v      []byte
	idxlen uint32
}

func (v vt) IndexLength() uint32 {
	return v.idxlen
}
func (v vt) Bytes() []byte {
	return v.v
}

func TestVector4(t *testing.T) {
	mt := newTestingMerkle(t, 4)
	defer mt.Storage().Close()

	zeros := make([]byte, 32, 32)
	assert.Nil(t, mt.Add(vt{zeros, uint32(1)}))
	v := vt{zeros, uint32(2)}
	assert.Nil(t, mt.Add(v))
	proof, _ := mt.GenerateProof(v)
	assert.True(t, CheckProof(mt.Root(), proof, v, mt.NumLevels()))
	assert.Equal(t, 4, mt.NumLevels())
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", hex.EncodeToString(v.Bytes()))
	assert.Equal(t, "0x8a18703d002272f39b71c3a586a265a92a841bbb2f553b93c7d7951ec892769d", mt.Root().Hex())
	assert.Equal(t, "000000000000000000000000000000000000000000000000000000000000000268a11c64aff5642515fe2cb8b46166fdf5b1661b81b2e7fe30fa5d6f9039b573", hex.EncodeToString(proof))
}

func TestVector140(t *testing.T) {
	mt := newTestingMerkle(t, 140)
	defer mt.Storage().Close()

	zeros := make([]byte, 32, 32)
	for i := 1; i < len(zeros)-1; i++ {
		v := vt{zeros, uint32(i)}
		assert.Nil(t, mt.Add(v))
		proof, err := mt.GenerateProof(v)
		assert.Nil(t, err)
		assert.True(t, CheckProof(mt.Root(), proof, v, mt.NumLevels()))
		if i == len(zeros)-2 {
			assert.Equal(t, 140, mt.NumLevels())
			assert.Equal(t, uint32(30), v.IndexLength())
			assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", hex.EncodeToString(v.Bytes()))
			assert.Equal(t, "0xf5ba7c6348183578db9218f9259ab2a7112034b1828583ee933da64b7f62200f", mt.Root().Hex())
			assert.Equal(t, "000000000000000000000000000000000000000000000000000000000000001f4267ed627de2aa903e6cf65af0779e870961e47f055a61ffbbfd015711338ff602e699e9c8bcdcdd7b26d899201d982ddd6a19c3617e2d2b2a9297696f0cd04ccc7bf435ab5eb2e08a788ad08a5dd3eb2546b727c40a51d1a4f59a6bac60525f047f00ce6f67cb17f9218bfb92dcaa7f7d2ce21eeb5435d7bab2e8fbc5e5f86b02d992bd0ba4d8dcfdc7c8f71b73eff00b794a3a7907fc3e6e3ddb6d9b6bf93b", hex.EncodeToString(proof))
		}
	}
}

/*** this version is compartible with the following solidity code


contract MerkleVerifier {

    function checkProof(bytes32 root, bytes proof, bytes value, uint256 indexlen, uint numlevels)
    public returns (bool){
        uint256 hi;
        assembly {
            hi := keccak256(add(value,32),indexlen)
        }

        uint256 emptiesmap;
        assembly {
            emptiesmap := mload(add(proof, 32))
        }

        uint256 nextSibling = 64;
        bytes32 nodehash = keccak256(value);

        for (uint256 level =  numlevels - 2 ; int256(level) >= 0; level--) {

            uint256 bitmask= 1 << level;
            bytes32 sibling;

            if (emptiesmap&bitmask>0) {
                assembly {
                    sibling := mload(add(proof, nextSibling))
                }
                nextSibling+=32;
            } else {
                sibling = 0x0;
            }

            if (hi&bitmask>0) {
                nodehash=keccak256(sibling,nodehash);
            } else {
                nodehash=keccak256(nodehash,sibling);
            }
        }
        return nodehash == root;
    }


    function testLevel4() {
        uint  depth = 4;
        uint8 indexlen = 2;
        bytes memory value = hex"0000000000000000000000000000000000000000000000000000000000000000";
        bytes32 root=0x8a18703d002272f39b71c3a586a265a92a841bbb2f553b93c7d7951ec892769d;
        bytes memory proof =hex'000000000000000000000000000000000000000000000000000000000000000268a11c64aff5642515fe2cb8b46166fdf5b1661b81b2e7fe30fa5d6f9039b573';
        require(checkProof(root,proof,value,indexlen,depth));
    }

    function testLevel140() {
        uint  depth = 140;
        uint8 indexlen = 30;
        bytes memory value = hex"0000000000000000000000000000000000000000000000000000000000000000";
        bytes32 root=0xf5ba7c6348183578db9218f9259ab2a7112034b1828583ee933da64b7f62200f;
        bytes memory proof = hex"000000000000000000000000000000000000000000000000000000000000001f4267ed627de2aa903e6cf65af0779e870961e47f055a61ffbbfd015711338ff602e699e9c8bcdcdd7b26d899201d982ddd6a19c3617e2d2b2a9297696f0cd04ccc7bf435ab5eb2e08a788ad08a5dd3eb2546b727c40a51d1a4f59a6bac60525f047f00ce6f67cb17f9218bfb92dcaa7f7d2ce21eeb5435d7bab2e8fbc5e5f86b02d992bd0ba4d8dcfdc7c8f71b73eff00b794a3a7907fc3e6e3ddb6d9b6bf93b";
        require(checkProof(root,proof,value,indexlen,depth));
    }

}

**/
