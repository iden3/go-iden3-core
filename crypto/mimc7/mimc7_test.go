package mimc7

import (
	"math/big"
	"testing"

	"github.com/arnaucube/go-snark/bn128"
	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	b1 := big.NewInt(int64(1))
	b2 := big.NewInt(int64(2))
	b3 := big.NewInt(int64(3))
	arrBigInt := []*big.Int{b1, b2, b3}

	// *big.Int array to RElem array
	rElems, err := BigIntsToRElems(arrBigInt)
	assert.Nil(t, err)

	// RElem array to *big.Int array
	bElems := RElemsToBigInts(rElems)

	assert.Equal(t, arrBigInt, bElems)

	r, ok := new(big.Int).SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)
	assert.True(t, ok)

	// greater or equal than 2**253 -1 number will give error when passing bigInts to RElems, to fit in the R Finite Field
	overR, ok := new(big.Int).SetString("14474011154664524427946373126085988481658748083205070504932198000989141204991", 10)
	assert.True(t, ok)
	_, err = BigIntsToRElems([]*big.Int{b1, overR, b2})
	assert.True(t, err != nil)

	_, err = BigIntsToRElems([]*big.Int{b1, r, b2})
	assert.True(t, err != nil)

	// smaller than 2**253 -1 number will not give error when passing bigInts to RElems, to fit in the R Finite Field
	underR, ok := new(big.Int).SetString("14474011154664524427946373126085988481658748083205070504932198000989141204990", 10)
	assert.True(t, ok)
	_, err = BigIntsToRElems([]*big.Int{b1, underR, b2})
	assert.Nil(t, err)
}

func TestMIMC7Generic(t *testing.T) {
	b1 := big.NewInt(int64(1))
	b2 := big.NewInt(int64(2))
	b3 := big.NewInt(int64(3))

	fqR, err := bn128.NewFqR()
	assert.Nil(t, err)

	bigArray := []*big.Int{b1, b2, b3}
	elementsArray, err := BigIntsToRElems(bigArray)
	assert.Nil(t, err)

	// Generic Hash
	mhg, err := MIMC7HashGeneric(fqR, b1, b2, 91)
	assert.Nil(t, err)
	assert.Equal(t, mhg.String(), "12727861578807605408775670530609274399788562023377579932124539591649752410226")
	hg, err := HashGeneric(elementsArray, fqR, 91)
	assert.Nil(t, err)
	assert.Equal(t, (*big.Int)(hg).String(), "10001192134743444757278983923787274376044444355175924720153500128284360571540")
}

func TestMIMC7(t *testing.T) {
	b1 := big.NewInt(int64(1))
	b2 := big.NewInt(int64(2))
	b3 := big.NewInt(int64(3))

	bigArray := []*big.Int{b1, b2, b3}
	elementsArray, err := BigIntsToRElems(bigArray)
	assert.Nil(t, err)

	// Hash
	mh := MIMC7Hash(b1, b2)
	assert.Nil(t, err)
	assert.Equal(t, mh.String(), "12727861578807605408775670530609274399788562023377579932124539591649752410226")

	h := Hash(elementsArray)
	assert.Nil(t, err)
	assert.Equal(t, (*big.Int)(h).String(), "10001192134743444757278983923787274376044444355175924720153500128284360571540")
}
