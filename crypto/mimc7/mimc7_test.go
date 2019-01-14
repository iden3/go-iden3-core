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

	// greater or equal than R will give error when passing bigInts to RElems, to fit in the R Finite Field
	overR, ok := new(big.Int).SetString("21888242871839275222246405745257275088548364400416034343698204186575808495618", 10)
	assert.True(t, ok)
	_, err = BigIntsToRElems([]*big.Int{b1, overR, b2})
	assert.True(t, err != nil)

	_, err = BigIntsToRElems([]*big.Int{b1, r, b2})
	assert.True(t, err != nil)

	// smaller than R will not give error when passing bigInts to RElems, to fit in the R Finite Field
	underR, ok := new(big.Int).SetString("21888242871839275222246405745257275088548364400416034343698204186575808495616", 10)
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
	assert.Equal(t, "10594780656576967754230020536574539122676596303354946869887184401991294982664", mhg.String())
	hg, err := HashGeneric(elementsArray, fqR, 91)
	assert.Nil(t, err)
	assert.Equal(t, "6464402164086696096195815557694604139393321133243036833927490113253119343397", (*big.Int)(hg).String())
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
	assert.Equal(t, "10594780656576967754230020536574539122676596303354946869887184401991294982664", mh.String())

	h := Hash(elementsArray)
	assert.Nil(t, err)
	assert.Equal(t, "6464402164086696096195815557694604139393321133243036833927490113253119343397", (*big.Int)(h).String())
}
