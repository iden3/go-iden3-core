package endpoint

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	// "github.com/ethereum/go-ethereum/common"

	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/genericserver"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/services/adminsrv"
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/services/namesrv"
	"github.com/iden3/go-iden3/services/rootsrv"
	"github.com/iden3/go-iden3/services/signedpacketsrv"

	log "github.com/sirupsen/logrus"
)

var claimService claimsrv.Service
var rootService rootsrv.Service
var nameService namesrv.Service
var signedPacketVerifier signedpacketsrv.SignedPacketVerifier
var adminService adminsrv.Service

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func handleGetRoot(c *gin.Context) {
	// get the contract data
	root, err := rootService.GetRoot(&genericserver.C.Id)
	if err != nil {
		genericserver.Fail(c, "error contract.GetRoot(contractAddress)", err)
		return
	}
	c.JSON(200, gin.H{
		"root":         claimService.MT().RootKey().Hex(),
		"contractRoot": common3.HexEncode(root[:]),
	})
}

func serveServiceApi() *http.Server {
	api, serviceapi := genericserver.NewServiceAPI("/api/unstable")

	serviceapi.POST("/names", handleVinculateId)
	serviceapi.GET("/names/:name", handleClaimAssignNameResolv)

	serviceapisrv := &http.Server{Addr: genericserver.C.Server.ServiceApi, Handler: api}
	go func() {
		if err := genericserver.ListenAndServe(serviceapisrv, "Service"); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return serviceapisrv
}

func serveAdminApi(stopch chan interface{}) *http.Server {
	api, _ := genericserver.NewAdminAPI("/api/unstable", stopch)
	adminapisrv := &http.Server{Addr: genericserver.C.Server.AdminApi, Handler: api}
	go func() {
		if err := genericserver.ListenAndServe(adminapisrv, "Admin"); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return adminapisrv
}

func Serve(rs rootsrv.Service, cs claimsrv.Service, ns namesrv.Service,
	spv signedpacketsrv.SignedPacketVerifier, as adminsrv.Service) {

	claimService = cs
	rootService = rs
	nameService = ns
	signedPacketVerifier = spv
	adminService = as

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
	rootService.Start()
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
