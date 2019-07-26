package signedpacketsrv

import (
	"github.com/iden3/go-iden3-core/core"
)

// RequestIdenAssertBody is the header request of a RequestIdenAssert.
type RequestIdenAssertHeader struct {
	Type string `json:"typ" binding:"required"`
}

// RequestIdenAssertBody is the body request of a RequestIdenAssert.
type RequestIdenAssertBody struct {
	Type string         `json:"type" binding:"required"`
	Data IdenAssertData `json:"data" binding:"required"`
}

// RequestIdenAssert is a request for a signed packet with payload type
// IDENASSERTV01.
type RequestIdenAssert struct {
	Header RequestIdenAssertHeader `json:"header" binding:"required"`
	Body   RequestIdenAssertBody   `json:"body" binding:"required"`
}

// NewRequestIdenAssert generates a signing request for a signed packet with
// payload type IDENASSERTV01.
func NewRequestIdenAssert(nonceDb *core.NonceDb, origin string, expireDelta int64) *RequestIdenAssert {
	nonceObj := nonceDb.New(expireDelta, nil)
	return &RequestIdenAssert{
		Header: RequestIdenAssertHeader{Type: SIGV02},
		Body: RequestIdenAssertBody{
			Type: IDENASSERTV01,
			Data: IdenAssertData{
				Challenge: nonceObj.Nonce,
				Timeout:   nonceObj.Expiration,
				Origin:    origin,
			},
		},
	}
}
