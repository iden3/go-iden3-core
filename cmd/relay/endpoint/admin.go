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

type addGenericClaimMsg struct {
	Namespace string
	IndexData string
	Data      string
}

func handleAddGenericClaim(c *gin.Context) {
	var m addGenericClaimMsg
	c.BindJSON(&m)

	proofOfClaim, err := adminservice.AddGenericClaim([]byte(m.IndexData), []byte(m.Data))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
	c.JSON(http.StatusOK, proofOfClaim.Hex())
}
