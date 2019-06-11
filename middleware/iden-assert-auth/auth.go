package auth

import (
	"fmt"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/services/signedpacketsrv"
)

const identityKey = "id"

// User represents an authenticated identity with an optionally assigned name.
type User struct {
	Id      core.ID
	EthName *string
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
// claims: "id" and "ethName", which are extracted from the idenAssert
// signed packet.
func NewAuthMiddleware(domain string, nonceDb *core.NonceDb, key []byte,
	signedPacketVerifier *signedpacketsrv.SignedPacketVerifier) (*jwt.GinJWTMiddleware, error) {
	// The JWT middleware
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "iden3.iden_assert.v0_1",
		Key:         key,
		Timeout:     24 * time.Hour,
		MaxRefresh:  0,
		IdentityKey: identityKey,
		// Add id JWT claim to the token
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if user, ok := data.(*User); ok {
				return jwt.MapClaims{
					"id":      user.Id,
					"ethName": user.EthName,
				}
			}
			return jwt.MapClaims{}
		},
		// Generate identity (id) from JWT claims
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)

			id, err := core.IDFromString(claims["id"].(string))
			if err != nil {
				panic(err)
			}

			var ethName *string
			switch v := claims["ethName"].(type) {
			case string:
				ethName = &v
			}
			return &User{
				Id:      id,
				EthName: ethName,
			}
		},
		// handler to validate login
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var idenAssert signedpacketsrv.SignedPacket
			if err := c.BindJSON(&idenAssert); err != nil {
				return nil, fmt.Errorf("invalid JWS signed packet: %v", err)
			}
			if res, err := signedPacketVerifier.VerifySignedPacketIdenAssert(&idenAssert, nonceDb,
				domain); err != nil {
				return nil, fmt.Errorf("failed verification of JWS signed packet: %v ", err)
			} else {
				return &User{Id: res.Id, EthName: res.EthName}, nil
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
func AddAuthMiddleware(r *gin.RouterGroup, domain string, nonceDb *core.NonceDb,
	key []byte, signedPacketVerifier *signedpacketsrv.SignedPacketVerifier) (*gin.RouterGroup, error) {
	authMiddleware, err := NewAuthMiddleware(domain, nonceDb, key, signedPacketVerifier)
	if err != nil {
		return nil, fmt.Errorf("JWT auth middleware error: %v", err)
	}

	r.GET("/login", func(c *gin.Context) {
		req := signedpacketsrv.NewRequestIdenAssert(nonceDb, domain, 60)
		c.JSON(200, gin.H{
			"sigReq": req,
		})
	})
	r.POST("/login", authMiddleware.LoginHandler)

	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	return auth, nil
}
