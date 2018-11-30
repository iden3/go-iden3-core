# Merkletree usage


## Import
Import packages:
```go
import (
  "github.com/iden3/go-iden3/db"
  "github.com/iden3/go-iden3/merkletree"
  "github.com/iden3/go-iden3/core"
  common3 "github.com/iden3/go-iden3/common"
)
```

## New Merkletree
Define new tree:
```go
// first we create the storage, where will be placed the leveldb database
storage, err := db.NewLevelDbStorage("./path", false)
if err!=nil {
  panic(err)
}
mt, err := merkletree.New(storage, 140) // new merkletree of 140 levels of depth using the defined storage
if err!=nil {
  panic(err)
}
defer mt.Storage().Close()
```

## Add claims
To add claims, first we need to have a claim data struct that fits the `Value` interface.
Value interface:
```go
// Value is the interface of a generic claim, a key value object stored in the leveldb
type Value interface {
	IndexLength() uint32 // returns the index length value
	Bytes() []byte // returns the value in byte array representation
}
```

We can use a new struct, or also use one of the already existing in the `go-iden3/core/claim.go`.

For this example, we will use the `core.GenericClaim`, as it have implemenented more useful methods arround it. We add two different claims into the merkletree:
```go
indexData := []byte("this is a first claim, this is the data that we put in the claim, that will affect it's position in the merkletree")
data := []byte("data that we put in the claim, that will not affect it's position in the merkletree")
claim0 := core.NewGenericClaim("namespace", "default", indexData, data)

indexData = []byte("this is a second claim")
data = []byte("data that we put in the claim, that will not affect it's position in the merkletree")
claim1 := core.NewGenericClaim("namespace", "default", indexData, data)
```
Once we have the `claim` struct that fits the `Value` interface, we can add it to the merkletree:
```go
err := mt.Add(claim)
if err!=nil {
  panic(err)
}
```

## Generate merkle proof
Now we can generat the merkle proof of this claim:
```go
mp, err := mt.GenerateProof(claim0.Hi())
if err!=nil {
  panic(err)
}

// If we display the merkleproof in hex:
mpHex := common3.BytesToHex(mp)
fmt.Println(mpHex)
// out: 0x00000000000000000000000000000000000000000000000000000000000000023035e951da1f81bea095e46ba26d9a4c29ed69aeb6678cc47247219d1c089250
```

## Check merkle proof
Now from a given merkle proof, we can check that it's data is consistent:
```go
checked := CheckProof(mt.Root(), mp, claim0.Hi(), claim0.Ht(), mt.NumLevels())
// checked == true
```

## Get value in position
We can also get the `claim` byte data in a certain position of the merkle tree (determined by its Hash_index (`Hi`)):
```go
claimInPosBytes, err := mt.GetValueInPos(claim0.Hi())
if err!=nil{
  panic(err)
}
```

## Proof of non existence
Also, we can generate a `Proof of non existence`, that is the merkle proof that a claim is not in the tree.
For example, we have this `claim2` that is not added in the merkletree:
useful methods arround it. We add two different claims into the merkletree:
```go
indexData := []byte("this claim will not be stored")
data := []byte("")
claim2 := core.NewGenericClaim("namespace", "default", indexData, data)
```
Now, we can generate the merkle proof of the data in the position of this claim in the merkletree:
```go
mp, err := mt.GenerateProof(claim2.Hi())
if err!=nil {
  panic(err)
}
```
In the `mp` we have the merkleproof that in the position of this `claim2` (that is determined by its `Hash_index` (`Hi`)) there is no data stored (so, it's an `EmptyNodeValue`, that is represented by an empty array of 32 bytes).

We can check this proof by calling the `CheckProof` function, but, in the parameter where we usually put the `Hash_total` (`Ht`), his time we put a `merkletree.EmptyNodeValue`, as we are checking that in that position (`Hi`) there is no data stored.
```go
checked := merkletree.CheckProof(mt.Root(), mp, claim0.Hi(), merkletree.EmptyNodeValue, mt.NumLevels())
// checked == true
```

## Complete example
The complete example can be found in this directory in the file [`example.go`]( https://github.com/iden3/go-iden3/blob/master/merkletreeDoc/example.go).
