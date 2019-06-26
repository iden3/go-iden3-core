package signedpacketsrv

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/services/signsrv"

	// "github.com/iden3/go-iden3/utils"
	"github.com/iden3/go-iden3-crypto/babyjub"
)

// SIGV01 is the JWS type of an iden3 signed packet.
const SIGV01 = "iden3.sig.v0_1" // V01 uses SIGALGV01 = "EK256K1"
const SIGV02 = "iden3.sig.v0_2" // V02 uses SIGALGV02 = "ED256BJ"

// IDENASSERTV01 is the signed packet payload type for an identity assertion.
const IDENASSERTV01 = "iden3.iden_assert.v0_1"

// GENERICSIGV01 is the signed packet payload type for a generic signature that
// contains an empty data field and a string key to string value mapping as
// form.
const GENERICSIGV01 = "iden3.gen_sig.v0_1"

// MSGV01 is the signed packet payload type of a signed message.  The signature
// is only used to guarantee the origin of the message.  The signature doesn't
// say anything about the message content.  This signed packet type is used for
// notifications.
const MSGV01 = "iden3.msg.v0_1"

// SIGALGV01 is the JWS algorithm used in SIGV01.  It's ECDSA with secp256k1
// and keccak.
const SIGALGV01 = "EK256K1"
const SIGALGV02 = "ED256BJ"

// SigHeader is the JSON Web Signature Header of a signed packet.
type SigHeader struct {
	Type         string  `json:"typ" binding:"required"`
	Issuer       core.ID `json:"iss" binding:"required"`
	IssuedAtTime int64   `json:"iat" binding:"required"`
	Expiration   int64   `json:"exp" binding:"required"`
	Algorithm    string  `json:"alg" binding:"required"`
}

// SigPayload is the JSON Web Signature Payload of a signed packet.
type SigPayload struct {
	Type       string             `json:"type" binding:"required"`
	KSign      *babyjub.PublicKey `json:"ksign" binding:"required"`
	ProofKSign core.ProofClaim    `json:"proofKSign" binding:"required"`
	DataRaw    json.RawMessage    `json:"data" binding:"required"`
	Data       interface{}        `json:"-"`
	FormRaw    json.RawMessage    `json:"form" binding:"required"`
	Form       interface{}        `json:"-"`
}

// MarshalJSON marshals the signed packet payload into JSON.
func (p SigPayload) MarshalJSON() ([]byte, error) {
	var err error
	p.DataRaw, err = json.Marshal(p.Data)
	if err != nil {
		return nil, err
	}
	p.FormRaw, err = json.Marshal(p.Form)
	if err != nil {
		return nil, err
	}
	type SigPayloadRaw SigPayload
	return json.Marshal(SigPayloadRaw(p))
}

// UnmarshalJSON unmarshals the signed packet payload from a JSON.
func (p *SigPayload) UnmarshalJSON(bs []byte) error {
	type SigPayloadRaw SigPayload
	var pRaw SigPayloadRaw
	if err := json.Unmarshal(bs, &pRaw); err != nil {
		return err
	}
	switch pRaw.Type {
	case IDENASSERTV01:
		var data IdenAssertData
		if err := json.Unmarshal(pRaw.DataRaw, &data); err != nil {
			return err
		}
		pRaw.Data = data
		var form *IdenAssertForm
		if err := json.Unmarshal(pRaw.FormRaw, &form); err != nil {
			return err
		}
		pRaw.Form = form
	case GENERICSIGV01:
		pRaw.Data = nil
		var form map[string]string
		if err := json.Unmarshal(pRaw.FormRaw, &form); err != nil {
			return err
		}
		pRaw.Form = form
	default:
		return fmt.Errorf("unknown signed packet type: %v", pRaw.Type)
	}
	*p = SigPayload(pRaw)
	return nil
}

// Jws represents a JWS (JSON Web Signature) sent over the network.
type Jws struct {
	Jws string `json:"jws" binding:"required"`
}

// SignedPacket is a JSON Web Signature unmarshaled packet of a signed packet.
type SignedPacket struct {
	Header      SigHeader
	Payload     SigPayload
	SignedBytes []byte
	Signature   *babyjub.SignatureComp
}

// Sign signs the signed packet with the key corresponding to addr.
func (sp *SignedPacket) Sign(signer signsrv.Service) error {
	headerJSON, err := json.Marshal(sp.Header)
	if err != nil {
		return err
	}
	payloadJSON, err := json.Marshal(sp.Payload)
	if err != nil {
		return err
	}
	sp.SignedBytes = []byte(fmt.Sprintf("%v.%v", base64.StdEncoding.EncodeToString([]byte(headerJSON)),
		base64.StdEncoding.EncodeToString([]byte(payloadJSON))))
	sp.Signature, err = signer.SignEthMsg(sp.SignedBytes)
	if err != nil {
		return err
	}
	return nil
}

// Marshal serializes a signed packet (that has been signed) into a string,
// encoding it as JWS.
func (sp *SignedPacket) Marshal() (string, error) {
	if sp.Signature == nil {
		return "", fmt.Errorf("signed packet has not been signed yet")
	}
	sig64 := base64.StdEncoding.EncodeToString(sp.Signature[:])
	return fmt.Sprintf("%v.%v", string(sp.SignedBytes), sig64), nil
}

// MarshalJSON marshals a signed packet into a Jws JSON.
func (sp *SignedPacket) MarshalJSON() ([]byte, error) {
	str, err := sp.Marshal()
	if err != nil {
		return nil, err
	}
	return json.Marshal(Jws{Jws: str})
}

// Unmarshal deserializes a signed packet (encoded as JWS) from a string.
func (sp *SignedPacket) Unmarshal(s string) error {
	fields := strings.Split(s, ".")
	if len(fields) != 3 {
		return fmt.Errorf("Invalid JWS: it doesn't contain 3 dot separated fields")
	}
	jwsHeader64, jwsPayload64, signature64 := fields[0], fields[1], fields[2]
	jwsHeader, err := common3.Base64Decode(jwsHeader64)
	if err != nil {
		return err
	}
	jwsPayload, err := common3.Base64Decode(jwsPayload64)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(jwsHeader, &sp.Header); err != nil {
		return err
	}
	if err := json.Unmarshal(jwsPayload, &sp.Payload); err != nil {
		return err
	}
	signature, err := common3.Base64Decode(signature64)
	if err != nil {
		return err
	}
	sp.Signature = &babyjub.SignatureComp{}
	copy(sp.Signature[:], signature)
	sp.SignedBytes = []byte(s[:strings.LastIndex(s, ".")])
	return nil
}

// UnmarshalJSON unmarshals a signed packet from a Jws JSON.
func (sp *SignedPacket) UnmarshalJSON(bs []byte) error {
	var jws Jws
	if err := json.Unmarshal(bs, &jws); err != nil {
		return err
	}
	return sp.Unmarshal(jws.Jws)
}

type SignedPacketSigner struct {
	signer     signsrv.Service
	id         core.ID
	proofKSign core.ProofClaim
}

func NewSignedPacketSigner(signer signsrv.Service, proofKSign core.ProofClaim,
	id core.ID) *SignedPacketSigner {
	return &SignedPacketSigner{
		id:         id,
		signer:     signer,
		proofKSign: proofKSign,
	}
}

func (sps *SignedPacketSigner) SetProofKSign(proofKSign core.ProofClaim) {
	sps.proofKSign = proofKSign
}

// NewSignPacketV02 generates and signs a SIGV02 signed packet.
func (sps *SignedPacketSigner) NewSignPacketV02(expireDelta int64,
	payloadType string, data interface{}, form interface{}) (*SignedPacket, error) {
	now := time.Now().Unix()
	header := SigHeader{
		Type:         SIGV02,
		Issuer:       sps.id,
		IssuedAtTime: now,
		Expiration:   now + expireDelta,
		Algorithm:    SIGALGV02,
	}
	payload := SigPayload{
		Type:       payloadType,
		Data:       data,
		KSign:      sps.signer.PublicKey(),
		ProofKSign: sps.proofKSign,
		Form:       form,
	}
	jws := SignedPacket{Header: header, Payload: payload}
	if err := jws.Sign(sps.signer); err != nil {
		return nil, err
	}
	return &jws, nil
}

// NewSignGenericSigV01 generates and signs a signed packet with payload type GENERICSIGV01.
func (sps *SignedPacketSigner) NewSignGenericSigV01(expireDelta int64,
	form interface{}) (*SignedPacket, error) {
	return sps.NewSignPacketV02(expireDelta, GENERICSIGV01, nil, form)

}

type MsgForm struct {
	Type string      `json:"type" binding:"required"`
	Data interface{} `json:"data" binding:"required"`
}

func (sps *SignedPacketSigner) NewSignMsgV01(expireDelta int64, msgType string,
	msg interface{}) (*SignedPacket, error) {
	return sps.NewSignPacketV02(expireDelta, MSGV01, nil, MsgForm{Type: msgType, Data: msg})

}

// IdenAssertData contains the data field of a signed packet of type
// iden3.iden_assert.v0_1
type IdenAssertData struct {
	Challenge string `json:"challenge" binding:"required"`
	Timeout   int64  `json:"timeout" binding:"required"`
	Origin    string `json:"origin" binding:"required"`
}

// IdenAssertForm contains the form field of a signed packet of type
// iden3.iden_assert.v0_1
type IdenAssertForm struct {
	EthName         string           `json:"ethName" binding:"required"`
	ProofAssignName *core.ProofClaim `json:"proofAssignName" binding:"required"`
}

// NewSignIdenAssertV01 generates and signs a signed packet with payload type
// IDENASSERTV01.  idenAssertForm may be nil if proving the ownership of a name
// is not desired.
func (sps *SignedPacketSigner) NewSignIdenAssertV01(requestIdenAssert *RequestIdenAssert,
	idenAssertForm *IdenAssertForm, expireDelta int64) (*SignedPacket, error) {
	return sps.NewSignPacketV02(expireDelta, IDENASSERTV01,
		requestIdenAssert.Body.Data, idenAssertForm)
}
