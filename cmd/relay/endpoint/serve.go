package endpoint

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/relay/config"
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/services/rootsrv"

	log "github.com/sirupsen/logrus"
)

var claimservice claimsrv.Service
var rootservice rootsrv.Service

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Add("Origin", "*")
		c.Writer.Header().Add("X-Requested-With", "*")
		c.Next()
	}
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func serveServiceApi() *http.Server{
	// start serviceapi
	serviceapi := gin.Default()
	serviceapi.Use(corsMiddleware())

	serviceapi.GET("/root", handleGetRoot)
	serviceapi.POST("/claim/:idaddr", handlePostClaim)
	serviceapi.GET("/claim/:idaddr/root", handleGetIDRoot)
	serviceapi.GET("/claim/:idaddr/hi/:hi", handleGetClaimByHi)
	serviceapisrv := &http.Server{Addr: config.C.Server.ServiceApi, Handler: serviceapi}
	go func() {
		log.Info("API server at ", config.C.Server.ServiceApi)
		if err := serviceapisrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("listen: %s\n", err)
		}
	}()
	return serviceapisrv
}

func serveAdminApi(stopch chan interface{}) *http.Server{
	adminapi := gin.Default()
	adminapi.Use(corsMiddleware())

	adminapi.POST("/stop", func(c *gin.Context) {
		// yeah, use curl -X POST http://<adminserver>/stop
		c.String(http.StatusOK,"got it, shutdowning server")
		stopch <- nil
	})

	adminapi.POST("/info", func(c *gin.Context) {
		// yeah, use curl -X POST http://<adminserver>/info
		c.String(http.StatusOK,"ping? pong!")
	})

	adminapisrv := &http.Server{Addr: config.C.Server.AdminApi, Handler: adminapi}
	go func() {
		log.Info("ADMIN server at ", config.C.Server.AdminApi)
		if err := adminapisrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("listen: %s\n", err)
		}
	}()
	return adminapisrv
}


func Serve(rs rootsrv.Service, cs claimsrv.Service) {

	claimservice = cs
	rootservice = rs

<<<<<<< HEAD
	r := gin.Default()
	// r.Use(corsMiddleware())
	r.Use(cors.Default())
	r.GET("/root", handleGetRoot)
	r.POST("/claim/:idaddr", handlePostClaim)
	r.GET("/claim/:idaddr/root", handleGetIDRoot)
	r.GET("/claim/:idaddr/hi/:hi", handleGetClaimByHi)
	r.Run(config.C.Server.ServiceApi)
=======
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
	serviceapisrv:=serveServiceApi()
	adminapisrv :=serveAdminApi(stopch)

	// wait until shutdown signal
	<-stopch
	log.Info("Shutdown Server ...")

	rootservice.StopAndJoin()

	if err := serviceapisrv.Shutdown(context.Background()); err != nil {
		log.Error("ServiceAapi Shutdown:", err)
	}

	if err := adminapisrv.Shutdown(context.Background()); err != nil {
		log.Error("AdminApi Shutdown:", err)
	}

>>>>>>> 4739c4bade9163aead6c46fe820b22b1cb4aca48
}
