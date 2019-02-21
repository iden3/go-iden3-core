package endpoint

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/claimserver/config"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"

	"github.com/iden3/go-iden3/merkletree"
)

func handleGetRoot(c *gin.Context) {
	// get the contract data
	contractAddress := common.HexToAddress(config.C.Contracts.RootCommits.Address)
	root, err := rootservice.GetRoot(contractAddress)
	if err != nil {
		fail(c, "error contract.GetRoot(contractAddress)", err)
		return
	}
	c.JSON(200, gin.H{
		"root":         claimservice.MT().RootKey().Hex(),
		"contractRoot": common3.HexEncode(root[:]),
	})
}

type claimData struct {
	IndexData string
	Data      string
}

// handlePostClaim handles the request to add a claim to a user tree.
func handlePostClaim(c *gin.Context) {
	var m claimData
	c.BindJSON(&m)

	indexData, err := common3.HexDecode(m.IndexData)
	if err != nil {
		fail(c, "error on handlePostClaim", err)
		return
	}
	data, err := common3.HexDecode(m.Data)
	if err != nil {
		fail(c, "error on handlePostClaim", err)
		return
	}

	if len(indexData) < 400/8 {
		fail(c, "error on handlePostClaim", errors.New("indexData smaller than 400/8"))
		return
	}
	if len(data) < 496/8 {
		fail(c, "error on handlePostClaim", errors.New("data smaller than 496/8"))
		return
	}

	var indexSlot [400 / 8]byte
	var dataSlot [496 / 8]byte
	copy(indexSlot[:], indexData[:400/8])
	copy(dataSlot[:], data[:496/8])
	claim := core.NewClaimBasic(indexSlot, dataSlot)
	err = claimservice.AddDirectClaim(*claim)
	if err != nil {
		fail(c, "error on AddDirectClaim", err)
		return
	}

	// return claim with proofs
	proofOfClaim, err := claimservice.GetClaimProofByHi(claim.Entry().HIndex())
	if err != nil {
		fail(c, "error on GetClaimByHi", err)
		return
	}

	c.JSON(200, gin.H{
		"proofClaim": proofOfClaim,
	})
	return
}

// handleGetClaimProofByHi handles the request to query the claim proof of a
// server claim (by hIndex).
func handleGetClaimProofByHi(c *gin.Context) {
	hihex := c.Param("hi")
	hiBytes, err := common3.HexDecode(hihex)
	if err != nil {
		fail(c, "error on HexDecode of Hi", err)
		return
	}
	hi := &merkletree.Hash{}
	copy(hi[:], hiBytes)
	proofOfClaim, err := claimservice.GetClaimProofByHi(hi)
	if err != nil {
		fail(c, "error on GetClaimProofByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"proofClaim": proofOfClaim,
	})
	return
}
