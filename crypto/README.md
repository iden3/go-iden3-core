# crypto [![GoDoc](https://godoc.org/github.com/iden3/go-iden3/crypto?status.svg)](https://godoc.org/github.com/iden3/go-iden3/crypto)
iden3 crypto Go package

## MIMC7
MIMC7 Hash function

Usage:
```go
package main

import (
	"math/big"
	"github.com/iden3/go-iden3/crypto/mimc7"
)

func mimc7Example() {
	// for this example, define an array of big ints to hash
	b1 := big.NewInt(int64(1))
	b2 := big.NewInt(int64(2))
	b3 := big.NewInt(int64(3))
	bigArr := []*big.Int{b1, b2, b3}
	arr, err := mimc7.BigIntsToRElems(bigArr)

	// mimc7 hash
	h := mimc7.Hash(arr)
	fmt.Println((*big.Int)(h).String())
	// h == 10001192134743444757278983923787274376044444355175924720153500128284360571540
}
```
