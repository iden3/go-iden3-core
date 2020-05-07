package zk

import (
	"fmt"
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	zktypes "github.com/iden3/go-circom-prover-verifier/types"
)

func G1ToBigInts(g1 *bn256.G1) [2]*big.Int {
	numBytes := 256 / 8
	bs := g1.Marshal()
	x := new(big.Int).SetBytes(bs[:numBytes])
	y := new(big.Int).SetBytes(bs[numBytes:])
	return [2]*big.Int{x, y}
}

func G2ToBigInts(g2 *bn256.G2) [2][2]*big.Int {
	numBytes := 256 / 8
	bs := g2.Marshal()
	xx := new(big.Int).SetBytes(bs[0*numBytes : 1*numBytes])
	xy := new(big.Int).SetBytes(bs[1*numBytes : 2*numBytes])
	yx := new(big.Int).SetBytes(bs[2*numBytes : 3*numBytes])
	yy := new(big.Int).SetBytes(bs[3*numBytes : 4*numBytes])
	// return [2][2]*big.Int{[2]*big.Int{xy, xx}, [2]*big.Int{yy, yx}}
	return [2][2]*big.Int{[2]*big.Int{xx, xy}, [2]*big.Int{yx, yy}}
}

func ProofToBigInts(proof *zktypes.Proof) (a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int) {
	a = G1ToBigInts(proof.A)
	b = G2ToBigInts(proof.B)
	c = G1ToBigInts(proof.C)
	return a, b, c
}

func PrintProof(proof *zktypes.Proof) {
	proofA, proofB, proofC := ProofToBigInts(proof)
	fmt.Printf(
		`    "a": ["%v",
	    "%v"],
`,
		proofA[0], proofA[1])
	fmt.Printf(
		`    "b": [
           ["%v",
            "%v"],
           ["%v",
            "%v"]],
`,
		proofB[0][0], proofB[0][1], proofB[1][0], proofB[1][1])
	fmt.Printf(
		`    "c": ["%v",
	    "%v"]
`,
		proofC[0], proofC[1])
}
