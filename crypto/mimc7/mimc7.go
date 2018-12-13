package mimc7

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/arnaucube/go-snark/bn128"
	"github.com/arnaucube/go-snark/fields"
	"github.com/ethereum/go-ethereum/crypto"
)

// MaxFieldVal is the maximum value that a element of a field can have, fitting inside an R Finite Field. This value is equal to 2**253 - 1
const MaxFieldValHex = "0x1fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
const SEED = "iden3_mimc"

// RElem is a big.Int of maximum 253 bits
type RElem *big.Int

var constants = generateConstantsData()

type constantsData struct {
	maxFieldVal *big.Int
	seedHash    *big.Int
	fqR         fields.Fq
	nRounds     int
	cts         []*big.Int
}

func generateConstantsData() constantsData {
	var constants constantsData
	maxFieldBytes, err := hex.DecodeString(MaxFieldValHex[2:])
	if err != nil {
		panic(err)
	}
	constants.maxFieldVal = new(big.Int).SetBytes(maxFieldBytes)
	constants.seedHash = new(big.Int).SetBytes(crypto.Keccak256([]byte(SEED)))
	fqR, err := bn128.NewFqR()
	if err != nil {
		panic(err)
	}
	constants.fqR = fqR
	constants.nRounds = 91
	cts, err := getConstants(constants.fqR, SEED, constants.nRounds)
	if err != nil {
		panic(err)
	}
	constants.cts = cts
	return constants
}

// BigIntToRElem checks if given big.Int fits in a Field R element, and returns the RElem type
func BigIntToRElem(a *big.Int) (RElem, error) {
	if a.Cmp(constants.maxFieldVal) != -1 {
		return RElem(a), errors.New("Given big.Int don't fits in the Finite Field over R")
	}
	return RElem(a), nil
}

//BigIntsToRElems converts from array of *big.Int to array of RElem
func BigIntsToRElems(arr []*big.Int) ([]RElem, error) {
	var o []RElem
	for i, a := range arr {
		e, err := BigIntToRElem(a)
		if err != nil {
			return o, fmt.Errorf("element in position %v don't fits in Finite Field over R", i)
		}
		o = append(o, e)
	}
	return o, nil
}

// RElemsToBigInts converts from array of RElem to array of *big.Int
func RElemsToBigInts(arr []RElem) []*big.Int {
	var o []*big.Int
	for _, a := range arr {
		o = append(o, a)
	}
	return o
}

func getConstants(fqR fields.Fq, seed string, nRounds int) ([]*big.Int, error) {
	var cts []*big.Int
	cts = append(cts, big.NewInt(int64(0)))
	c := new(big.Int).SetBytes(crypto.Keccak256([]byte(SEED)))
	for i := 1; i < nRounds; i++ {
		c = new(big.Int).SetBytes(crypto.Keccak256(c.Bytes()))

		n := fqR.Affine(c)
		cts = append(cts, n)
	}
	return cts, nil
}

// MIMC7HashGeneric performs the MIMC7 hash over a RElem, in a generic way, where it can be specified the Finite Field over R, and the number of rounds
func MIMC7HashGeneric(fqR fields.Fq, xIn, k *big.Int, nRounds int) (*big.Int, error) {
	cts, err := getConstants(fqR, SEED, nRounds)
	if err != nil {
		return &big.Int{}, err
	}
	var r *big.Int
	for i := 0; i < nRounds; i++ {
		var t *big.Int
		if i == 0 {
			t = fqR.Add(xIn, k)
		} else {
			t = fqR.Add(fqR.Add(r, k), cts[i])
		}
		t2 := fqR.Square(t)
		t4 := fqR.Square(t2)
		r = fqR.Mul(fqR.Mul(t4, t2), t)
	}
	return fqR.Affine(fqR.Add(r, k)), nil
}

// HashGeneric performs the MIMC7 hash over a RElem array, in a generic way, where it can be specified the Finite Field over R, and the number of rounds
func HashGeneric(arrEl []RElem, fqR fields.Fq, nRounds int) (RElem, error) {
	arr := RElemsToBigInts(arrEl)
	r := fqR.Zero()
	var err error
	for i := 0; i < len(arr); i++ {
		r, err = MIMC7HashGeneric(fqR, r, arr[i], nRounds)
		if err != nil {
			return r, err
		}
	}
	return RElem(r), nil
}

// MIMC7Hash performs the MIMC7 hash over a RElem, using the Finite Field over R and the number of rounds setted in the `constants` variable
func MIMC7Hash(xIn, k *big.Int) *big.Int {
	var r *big.Int
	for i := 0; i < constants.nRounds; i++ {
		var t *big.Int
		if i == 0 {
			t = constants.fqR.Add(xIn, k)
		} else {
			t = constants.fqR.Add(constants.fqR.Add(r, k), constants.cts[i])
		}
		t2 := constants.fqR.Square(t)
		t4 := constants.fqR.Square(t2)
		r = constants.fqR.Mul(constants.fqR.Mul(t4, t2), t)
	}
	return constants.fqR.Affine(constants.fqR.Add(r, k))
}

// Hash performs the MIMC7 hash over a RElem array
func Hash(arrEl []RElem) RElem {
	arr := RElemsToBigInts(arrEl)
	r := constants.fqR.Zero()
	for i := 0; i < len(arr); i++ {
		r = MIMC7Hash(r, arr[i])
	}
	return RElem(r)
}
