package endpoint

import (
	"net/http"

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
	c.String(http.StatusOK, r)
}

func handleClaimsDump(c *gin.Context) {
	r := adminservice.ClaimsDump()
	c.String(http.StatusOK, r)
}
