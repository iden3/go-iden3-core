package endpoint

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/relay/config"
	common3 "github.com/iden3/go-iden3/common"
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
