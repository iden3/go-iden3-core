package endpoint

import (
	"math/big"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

func handleMimc7(c *gin.Context) {
	var elements []*big.Int
	c.BindJSON(&elements)

	r, err := adminservice.Mimc7(elements)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
	c.String(http.StatusOK, r.String())
}

type addClaimBasicMsg struct {
	Namespace string
	IndexData string
	Data      string
}

func handleAddClaimBasic(c *gin.Context) {
	var m addClaimBasicMsg
	c.BindJSON(&m)

	if len(m.IndexData) != 400/8 {
		c.String(http.StatusBadRequest, "indexData smaller than 400/8")
		return
	}
	if len(m.Data) != 496/8 {
		c.String(http.StatusBadRequest, "data smaller than 496/8")
		return
	}

	var indexSlot [400 / 8]byte
	var dataSlot [496 / 8]byte
	copy(indexSlot[:], m.IndexData[:400/8])
	copy(dataSlot[:], m.Data[:496/8])
	proofOfClaim, err := adminservice.AddClaimBasic(indexSlot, dataSlot)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
	c.JSON(http.StatusOK, proofOfClaim)
}
