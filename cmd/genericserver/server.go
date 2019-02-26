package genericserver

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewServiceAPI(prefix string) (*gin.Engine, *gin.RouterGroup) {
	api := gin.Default()
	api.Use(cors.Default())

	serviceapi := api.Group(prefix)
	serviceapi.GET("/root", HandleGetRoot)

	return api, serviceapi
}

func NewAdminAPI(prefix string, stopch chan interface{}) (*gin.Engine, *gin.RouterGroup) {
	api := gin.Default()
	api.Use(cors.Default())
	adminapi := api.Group("/api/unstable")

	adminapi.POST("/stop", func(c *gin.Context) {
		// yeah, use curl -X POST http://<adminserver>/stop
		c.String(http.StatusOK, "got it, shutdowning server")
		stopch <- nil
	})

	adminapi.GET("/info", HandleInfo)
	adminapi.GET("/rawdump", HandleRawDump)
	adminapi.POST("/rawimport", HandleRawImport)
	adminapi.GET("/claimsdump", HandleClaimsDump)
	return api, adminapi
}
