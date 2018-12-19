package endpoint

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/relay/config"
	"github.com/iden3/go-iden3/services/adminsrv"
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/services/identitysrv"
	"github.com/iden3/go-iden3/services/namesrv"
	"github.com/iden3/go-iden3/services/rootsrv"

	log "github.com/sirupsen/logrus"
)

var claimservice claimsrv.Service
var rootservice rootsrv.Service

var nameservice namesrv.Service
var idservice identitysrv.Service

var adminservice adminsrv.Service

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func serveServiceApi() *http.Server {
	// start serviceapi
	serviceapi := gin.Default()
	serviceapi.Use(cors.Default())

	serviceapi.GET("/root", handleGetRoot)

	serviceapi.POST("/root/:idaddr", handleCommitNewIDRoot)
	serviceapi.POST("/claim/:idaddr", handlePostClaim)
	serviceapi.GET("/claim/:idaddr/root", handleGetIDRoot)
	serviceapi.GET("/claim/:idaddr/hi/:hi", handleGetClaimByHi)

	serviceapi.POST("/id", handleCreateId)
	serviceapi.GET("/id/:idaddr", handleGetId)
	serviceapi.POST("/id/:idaddr/deploy", handleDeployId)
	serviceapi.POST("/id/:idaddr/forward", handleForwardId)

	serviceapi.POST("/vinculateid", handleVinculateID)
	serviceapi.GET("/identities/resolv/:nameid", handleAssignNameClaimResolv)

	serviceapisrv := &http.Server{Addr: config.C.Server.ServiceApi, Handler: serviceapi}
	go func() {
		log.Info("API server at ", config.C.Server.ServiceApi)
		if err := serviceapisrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("listen: %s\n", err)
		}
	}()
	return serviceapisrv
}

func serveAdminApi(stopch chan interface{}) *http.Server {
	adminapi := gin.Default()
	adminapi.Use(cors.Default())

	adminapi.POST("/stop", func(c *gin.Context) {
		// yeah, use curl -X POST http://<adminserver>/stop
		c.String(http.StatusOK, "got it, shutdowning server")
		stopch <- nil
	})

	adminapi.GET("/info", handleInfo)
	adminapi.GET("/rawdump", handleRawDump)
	adminapi.GET("/claimsdump", handleClaimsDump)
	adminapi.POST("/mimc7", handleMimc7)
	adminapi.POST("/genericClaim", handleAddGenericClaim)

	adminapisrv := &http.Server{Addr: config.C.Server.AdminApi, Handler: adminapi}
	go func() {
		log.Info("ADMIN server at ", config.C.Server.AdminApi)
		if err := adminapisrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("listen: %s\n", err)
		}
	}()
	return adminapisrv
}

func Serve(rs rootsrv.Service, cs claimsrv.Service, ids identitysrv.Service, ns namesrv.Service, as adminsrv.Service) {

	idservice = ids
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
