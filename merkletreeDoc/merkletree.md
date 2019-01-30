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
// new merkletree of 140 levels of maximum depth using the defined
// storage
mt, err := merkletree.NewMerkleTree(storage, 140)
if err!=nil {
  panic(err)
}
defer mt.Storage().Close()
```

## Add claims
To add claims, first we need to have a claim data struct that fits the `Entrier` interface:
```go
// Data consists of 4 elements of the mimc7 field.
type Data [4]ElemBytes
// An Entry contains Data where the claim will be serialized.
type Entry struct {
  Data Data
  [...]
}
// Entrier is the interface of a generic claim.
type Entrier interface {
  Entry() *Entry
}
```

We can use a new struct, or also use one of the already existing in the `go-iden3/core/claim.go`.

For this example, we will use the `core.ClaimAssignName`.
We add two different claims into the merkletree:
```go
name0 := "alice@iden3.io"
ethAddr0 := common.HexToAddress("0x7b471a1bdbd3b8ac98f3715507449f3a8e1f3b22")
claim0 := core.NewClaimAssignName(name0, ethAddr0)
claimEntry0 := claim0.Entry()

name1 := "bob@iden3.io"
ethAddr1 := common.HexToAddress("0x28f8267fb21e8ce0cdd9888a6e532764eb8d52dd")
claim1 := core.NewClaimAssignName(name1, ethAddr1)
claimEntry1 := claim1.Entry()
```
Once we have the `claim` struct that fits the `Entrier` interface, we can add it to the merkletree:
```go
err = mt.Add(claimEntry0)
if err != nil {
  panic(err)
}
err = mt.Add(claimEntry1)
if err != nil {
  panic(err)
}
```

## Generate merkle proof
Now we can generat the merkle proof of this claim:
```go
mp, err := mt.GenerateProof(claimEntry0.HIndex())
if err != nil {
  panic(err)
}

// We can display the merkleproof:
fmt.Println("merkle proof: ", mp)
// out: 
// merkle proof:  Proof:
//         existence: true
//         depth: 2
//         notempties: 01
//         siblings: 0 a045683a
```

## Check merkle proof
Now from a given merkle proof, we can check that it's data is consistent:
```go
checked := merkletree.VerifyProof(mt.RootKey(), mp,
				  claimEntry0.HIndex(), claimEntry0.HValue())
// checked == true
```

## Get value in position
We can also get the `claim` byte data in a certain position of the merkle tree
(determined by its Hash Index (`HIndex`)):
```go
claimDataInPos, err := mt.GetDataByIndex(claimEntry0.HIndex())
if err!=nil{
  panic(err)
}
```

## Proof of non existence
Also, we can generate a `Proof of non existence`, that is, the merkle proof
that a claim is not in the tree.
For example, we have this `claim2` that is not added in the merkletree:
```go
name2 := "eve@iden3.io"
ethAddr2 := common.HexToAddress("0x29a6a240e2d8f8bf39b5338b9664d414c5d793f4")
claim2 := core.NewClaimAssignName(name2, ethAddr2)
claimEntry2 := claim2.Entry()
```
Now, we can generate the merkle proof of the data in the position of this claim
in the merkletree, and print it to see that it's a non-existence proof:
```go
mp, err = mt.GenerateProof(claimEntry2.HIndex())
if err != nil {
  panic(err)
}

// We can display the merkleproof:
fmt.Println("merkle proof: ", mp)
// out: 
// merkle proof:  Proof:
//         existence: false
//         depth: 2
//         notempties: 01
//         siblings: 0 a045683a
//         node aux: hi: c641b925, ht: eeae8c7e
```
In the `mp` we have the merkleproof that in the position of this `claim2` (that
is determined by its Hash Index (`HIndex`)) there is no data stored (so, it's an
`NodeTypeEmpty` not actually stored in the tree).

We can check this proof by calling the `VerifyProof` function, and in the
parameter where we put the Hash Total (`HtTotal`) we can actually put anything, because we can proof that anything is not there.  We will use the Hash Total of the claim2 for convenience.
```go
checked = merkletree.VerifyProof(mt.RootKey(), mp, claimEntry2.HIndex(), claimEntry2.HValue())
// checked == true
```

## Complete example
The complete example can be found in this directory in the file [`example.go`]( https://github.com/iden3/go-iden3/blob/master/merkletreeDoc/example.go).
