package genericserver

import (
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	common3 "github.com/iden3/go-iden3/common"
	log "github.com/sirupsen/logrus"
)

func Fail(c *gin.Context, msg string, err error) {
	if err != nil {
		log.WithError(err).Error(msg)
	} else {
		log.Error(msg)
	}
	c.JSON(400, gin.H{
		"error": msg,
	})
	return
}

// Generic
func HandleGetRoot(c *gin.Context) {
	// get the contract data
	contractAddress := common.HexToAddress(C.Contracts.RootCommits.Address)
	root, err := Rootservice.GetRoot(contractAddress)
	if err != nil {
		Fail(c, "error contract.GetRoot(contractAddress)", err)
		return
	}
	c.JSON(200, gin.H{
		"root":         Claimservice.MT().RootKey().Hex(),
		"contractRoot": common3.HexEncode(root[:]),
	})
}

// Admin
func HandleInfo(c *gin.Context) {
	r := Adminservice.Info(common.HexToAddress(C.Contracts.RootCommits.Address))
	c.JSON(200, gin.H{
		"info": r,
	})
}
func HandleRawDump(c *gin.Context) {
	Adminservice.RawDump(c)
}

func HandleRawImport(c *gin.Context) {
	var data map[string]string
	err := c.BindJSON(&data)
	if err != nil {
		Fail(c, "json parsing error", err)
		return
	}

	count, err := Adminservice.RawImport(data)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
	c.String(http.StatusOK, "imported "+strconv.Itoa(count)+" key,value entries")
}

func HandleClaimsDump(c *gin.Context) {
	r := Adminservice.ClaimsDump()
	c.JSON(http.StatusOK, r)
}
