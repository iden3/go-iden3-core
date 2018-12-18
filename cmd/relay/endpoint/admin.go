package endpoint

import (
	"github.com/gin-gonic/gin"
)

func handleRawDump(c *gin.Context) {
	out := adminservice.RawDump()
	c.JSON(200, gin.H{
		"dump": out,
	})
}

func handleClaimsDump(c *gin.Context) {
	out := adminservice.ClaimsDump()
	c.JSON(200, gin.H{
		"claimsdump": out,
	})
}
