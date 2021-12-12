package core

import (
	"encoding/hex"
	"fmt"
	"time"
)

func ExampleNewClaim() {
	var schemaHash SchemaHash
	expDate := time.Date(2021, 1, 10, 20, 30, 0, 0, time.UTC)
	claim, err := NewClaim(schemaHash,
		WithExpirationDate(expDate),
		WithVersion(42))
	if err != nil {
		panic(err)
	}
	expDateRes, ok := claim.GetExpirationDate()
	fmt.Println(ok)
	fmt.Println(expDateRes.In(time.UTC).Format(time.RFC3339))

	fmt.Println(claim.GetVersion())

	entry := claim.TreeEntry()
	indexHash, valueHash, err := entry.HiHv()
	if err != nil {
		panic(err)
	}
	fmt.Println(hex.EncodeToString(indexHash[:]))
	fmt.Println(hex.EncodeToString(valueHash[:]))
	// Output:
	// true
	// 2021-01-10T20:30:00Z
	// 42
	// a07b32a81b631544f9199f4bf429ad2026baec31ba5e5e707a49cc2c9d243f18
	// 8e6bca4b559d758eca7b6125faea23ed0765cdcb6f85b3fe9477ca4293a6fd05
}
