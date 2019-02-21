package endpoint

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	common3 "github.com/iden3/go-iden3/common"
)

func handleInfo(c *gin.Context) {
	r := adminservice.Info()
	c.JSON(200, gin.H{
		"info": r,
	})
}

func handleRawDump(c *gin.Context) {
	r := adminservice.RawDump()
	c.JSON(http.StatusOK, r)
}

func handleRawImport(c *gin.Context) {
	var data map[string]string
	c.BindJSON(&data)

	count, err := adminservice.RawImport(data)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
	c.String(http.StatusOK, "imported "+strconv.Itoa(count)+" key,value entries")
}

func handleClaimsDump(c *gin.Context) {
	r := adminservice.ClaimsDump()
	c.JSON(http.StatusOK, r)
}

type addClaimBasicMsg struct {
	IndexData string
	Data      string
}

func handleAddClaimBasic(c *gin.Context) {
	var m addClaimBasicMsg
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
	proofOfClaim, err := adminservice.AddClaimBasic(indexSlot, dataSlot)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
	c.JSON(http.StatusOK, proofOfClaim)
}
