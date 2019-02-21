package endpoint

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/nameserver/config"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/services/adminsrv"
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/services/namesrv"
	"github.com/iden3/go-iden3/services/rootsrv"

	log "github.com/sirupsen/logrus"
)

var claimservice claimsrv.Service
var rootservice rootsrv.Service

var nameservice namesrv.Service

var adminservice adminsrv.Service

func init() {
	gin.SetMode(gin.ReleaseMode)
}

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

func serveServiceApi() *http.Server {
	// start serviceapi
	api := gin.Default()
	api.Use(cors.Default())

	serviceapi := api.Group("/api/unstable")
	serviceapi.GET("/root", handleGetRoot)

	serviceapi.POST("/names", handleVinculateId)
	serviceapi.GET("/names/:nameid", handleClaimAssignNameResolv)

	serviceapisrv := &http.Server{Addr: config.C.Server.ServiceApi, Handler: api}
	go func() {
		log.Info("API server at ", config.C.Server.ServiceApi)
		if err := serviceapisrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("listen: %s\n", err)
		}
	}()
	return serviceapisrv
}

func serveAdminApi(stopch chan interface{}) *http.Server {
	api := gin.Default()
	api.Use(cors.Default())
	adminapi := api.Group("/api/unstable")

	adminapi.POST("/stop", func(c *gin.Context) {
		// yeah, use curl -X POST http://<adminserver>/stop
		c.String(http.StatusOK, "got it, shutdowning server")
		stopch <- nil
	})

	adminapi.GET("/info", handleInfo)
	adminapi.GET("/rawdump", handleRawDump)
	adminapi.POST("/rawimport", handleRawImport)
	adminapi.GET("/claimsdump", handleClaimsDump)

	adminapisrv := &http.Server{Addr: config.C.Server.AdminApi, Handler: api}
	go func() {
		log.Info("ADMIN server at ", config.C.Server.AdminApi)
		if err := adminapisrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("listen: %s\n", err)
		}
	}()
	return adminapisrv
}

func Serve(rs rootsrv.Service, cs claimsrv.Service, ns namesrv.Service, as adminsrv.Service) {

	claimservice = cs
	rootservice = rs
	nameservice = ns
	adminservice = as

	stopch := make(chan interface{})

	// catch ^C to send the stop signal
	ossig := make(chan os.Signal, 1)
	signal.Notify(ossig, os.Interrupt)
	go func() {
		for sig := range ossig {
			if sig == os.Interrupt {
				stopch <- nil
			}
		}
	}()

	// start servers
	rootservice.Start()
	serviceapisrv := serveServiceApi()
	adminapisrv := serveAdminApi(stopch)

	// wait until shutdown signal
	<-stopch
	log.Info("Shutdown Server ...")

	if err := serviceapisrv.Shutdown(context.Background()); err != nil {
		log.Error("ServiceApi Shutdown:", err)
	} else {
		log.Info("ServiceApi stopped")
	}

	if err := adminapisrv.Shutdown(context.Background()); err != nil {
		log.Error("AdminApi Shutdown:", err)
	} else {
		log.Info("AdminApi stopped")
	}

}
