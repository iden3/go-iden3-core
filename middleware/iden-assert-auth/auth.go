package auth

import (
	"fmt"
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
)

const identityKey = "id"

// User represents an authenticated identity with an assigned name.
type User struct {
	IdAddr  common.Address
	EthName string
}

// GetUser extracts the User from the gin.Context for an authenticated
// endpoint.
func GetUser(c *gin.Context) *User {
	user, _ := c.Get(identityKey)
	return user.(*User)
}

// NewAuthMiddleware creates a JWT middleware struct that contains an endpoint
// for authenticatin a user through an idenAssert signed packet as well as a
// gin middleware to validate a JWT session token.  The JWT contains two
// claims: "idAddr" and "ethName", which are extracted from the idenAssert
// signed packet.
func NewAuthMiddleware(domain string, nonceDb *core.NonceDb, key []byte) (*jwt.GinJWTMiddleware, error) {
	// The JWT middleware
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "iden3.iden_assert.v0_1",
		Key:         key,
		Timeout:     24 * time.Hour,
		MaxRefresh:  0,
		IdentityKey: identityKey,
		// Add idAddr JWT claim to the token
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if user, ok := data.(*User); ok {
				return jwt.MapClaims{
					"idAddr":  user.IdAddr,
					"ethName": user.EthName,
				}
			}
			return jwt.MapClaims{}
		},
		// Generate identity (idAddr) from JWT claims
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			var idAddr common.Address
			if err := common3.HexDecodeInto(idAddr[:],
				[]byte(claims["idAddr"].(string))); err != nil {
				panic(err)
			}
			return &User{
				IdAddr:  idAddr,
				EthName: claims["ethName"].(string),
			}
		},
		// handler to validate login
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var idenAssert core.SignedPacket
			if err := c.BindJSON(&idenAssert); err != nil {
				return nil, fmt.Errorf("invalid JWS signed packet: %v", err)
			}
			if res, err := core.VerifySignedPacketIdenAssert(&idenAssert, nonceDb,
				domain); err != nil {
				return nil, fmt.Errorf("failed verification of JWS signed packet: %v ", err)
			} else {
				return &User{IdAddr: res.IdAddr, EthName: res.EthName}, nil
			}
		},
		// handler for failed authentication (when Authenticator returns error)
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup: "header: Authorization",
		TimeFunc:    time.Now,
	})
}

// AddAuthMiddleware adds the login endpoints to r and returns a router with
// prefix "/auth" that contains the JWT authorization validation middleware.
// The login endpoints are:
//    - GET "/login": returns a JSON with sigReq field containing a
//    RequestIdenAssert.
//    - POST "/login": requires a JSON with jwt field containing a JWS
//    signedPacket.
func AddAuthMiddleware(r *gin.Engine, domain string, nonceDb *core.NonceDb, key []byte) (*gin.RouterGroup, error) {
	authMiddleware, err := NewAuthMiddleware(domain, nonceDb, key)
	if err != nil {
		return nil, fmt.Errorf("JWT auth middleware error: %v", err)
	}

	r.GET("/login", func(c *gin.Context) {
		req := core.NewRequestIdenAssert(nonceDb, domain, 60)
		c.JSON(200, gin.H{
			"sigReq": req,
		})
	})
	r.POST("/login", authMiddleware.LoginHandler)

	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	return auth, nil
}
