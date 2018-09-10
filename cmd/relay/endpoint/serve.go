package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/relay/config"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/syndtr/goleveldb/leveldb"
)

var mt *merkletree.MerkleTree
var dbNameResolver *leveldb.DB

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Add("Origin", "*")
		c.Writer.Header().Add("X-Requested-With", "*")
		c.Next()
	}
}

func Serve(mtree *merkletree.MerkleTree) {

	mt = mtree

	r := gin.Default()
	r.Use(corsMiddleware())
	r.GET("/root", handleGetRoot)
	r.POST("/claim/:idaddr", handlePostClaim)
	r.GET("/claim/:idaddr/root", handleGetIDRoot)
	r.GET("/claim/:idaddr/hi/:hi", handleGetClaimByHi)
	r.Run(":" + config.C.Server.Port)
}
