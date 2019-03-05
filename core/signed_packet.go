package core

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/utils"
)

// Temporary hardcoded relay address
var relayAddr common.Address

func init() {
	relayAddr = common.HexToAddress("0xe0fbce58cfaa72812103f003adce3f284fe5fc7c")
}

const SIGV01 = "iden3.sig.v0_1"
const IDENASSERTV01 = "iden3.iden_assert.v0_1"
const GENERICSIGV01 = "iden3.gen_sig.v0_1"
const SIGALGV01 = "EK256K1"

// SigHeader is the JSON Web Signature Header of a signed packet
type SigHeader struct {
	Type         string         `json:"typ" binding:"required"`
	Issuer       common.Address `json:"iss" binding:"required"`
	IssuedAtTime int64          `json:"iat" binding:"required"`
	Expiration   int64          `json:"exp" binding:"required"`
	Algorithm    string         `json:"alg" binding:"required"`
}

// SigPayload is the JSON Web Signature Payload of a signed packet
type SigPayload struct {
	Type       string           `json:"type" binding:"required"`
	KSign      *utils.PublicKey `json:"ksign" binding:"required"`
	ProofKSign ProofClaim       `json:"proofKSign" binding:"required"`
	DataRaw    json.RawMessage  `json:"data" binding:"required"`
	Data       interface{}      `json:"-"`
	FormRaw    json.RawMessage  `json:"form" binding:"required"`
	Form       interface{}      `json:"-"`
}

type IdenAssertData struct {
	Challenge string `json:"challenge" binding:"required"`
	Timeout   int64  `json:"timeout" binding:"required"`
	Origin    string `json:"origin" binding:"required"`
}

type IdenAssertForm struct {
	EthName         string
	ProofAssignName ProofClaim
}

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
		var form IdenAssertForm
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

// SignedPacket is a JSON Web Signature unmarshaled packet of a signed packet
type SignedPacket struct {
	Header      SigHeader
	Payload     SigPayload
	SignedBytes []byte
	Signature   *utils.SignatureEthMsg
}

func (sp *SignedPacket) Sign(ks *keystore.KeyStore, addr common.Address) error {
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
	sp.Signature, err = utils.SignEthMsg(ks, accounts.Account{Address: addr}, sp.SignedBytes)
	if err != nil {
		return err
	}
	return nil
}

type Jws struct {
	Jws string `json:"jws" binding:"required"`
}

func (sp *SignedPacket) Marshal() (string, error) {
	if sp.Signature == nil {
		return "", fmt.Errorf("signed packet has not been signed yet")
	}
	sig64 := base64.StdEncoding.EncodeToString(sp.Signature[:])
	return fmt.Sprintf("%v.%v", string(sp.SignedBytes), sig64), nil
}

func (sp *SignedPacket) MarshalJSON() ([]byte, error) {
	str, err := sp.Marshal()
	if err != nil {
		return nil, err
	}
	return json.Marshal(Jws{Jws: str})
}

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
	sp.Signature = &utils.SignatureEthMsg{}
	copy(sp.Signature[:], signature)
	sp.SignedBytes = []byte(s[:strings.LastIndex(s, ".")])
	return nil
}

func (sp *SignedPacket) UnmarshalJSON(bs []byte) error {
	var jws Jws
	if err := json.Unmarshal(bs, &jws); err != nil {
		return err
	}
	return sp.Unmarshal(jws.Jws)
}

func NewSignPacketV01(ks *keystore.KeyStore, idAddr common.Address, kSignPk *ecdsa.PublicKey,
	proofKSign ProofClaim, expireDelta int64, payloadType string,
	data interface{}, form interface{}) (*SignedPacket, error) {
	now := time.Now().Unix()
	header := SigHeader{
		Type:         SIGV01,
		Issuer:       idAddr,
		IssuedAtTime: now,
		Expiration:   now + expireDelta,
		Algorithm:    SIGALGV01,
	}
	payload := SigPayload{
		Type:       payloadType,
		Data:       data,
		KSign:      &utils.PublicKey{PublicKey: *kSignPk},
		ProofKSign: proofKSign,
		Form:       form,
	}
	jws := SignedPacket{Header: header, Payload: payload}
	if err := jws.Sign(ks, crypto.PubkeyToAddress(*kSignPk)); err != nil {
		return nil, err
	}
	return &jws, nil
}

func NewSignGenericSigV01(ks *keystore.KeyStore, idAddr common.Address, kSignPk *ecdsa.PublicKey,
	proofKSign ProofClaim, expireDelta int64, form interface{}) (*SignedPacket, error) {
	return NewSignPacketV01(ks, idAddr, kSignPk, proofKSign, expireDelta,
		GENERICSIGV01, nil, form)

}

func NewSignIdenAssertV01(requestIdenAssert *RequestIdenAssert, ethName string,
	proofAssignName *ProofClaim, ks *keystore.KeyStore, idAddr common.Address,
	kSignPk *ecdsa.PublicKey, proofKSign ProofClaim, expireDelta int64) (*SignedPacket, error) {
	return NewSignPacketV01(ks, idAddr, kSignPk, proofKSign, expireDelta,
		IDENASSERTV01, requestIdenAssert.Body.Data,
		IdenAssertForm{EthName: ethName, ProofAssignName: *proofAssignName})
}

func VerifySignedPacketV01(jws *SignedPacket) error {
	// 2. Verify jwsHeader.alg is 'ES255'
	if jws.Header.Algorithm != SIGALGV01 {
		return fmt.Errorf("Unsupported alg: %v", jws.Header.Algorithm)
	}

	// 3. Verify that jwsHeader.iat <= now() < jwsHeader.exp
	now := time.Now().Unix()
	// Moving iat 2 minutes in the past to accomodate time shifts in time synchronization.
	if !((jws.Header.IssuedAtTime-120 <= now) && (now < jws.Header.Expiration)) {
		return fmt.Errorf("Signature not valid for current date (iat:%v, now:%v, exp:%v)",
			jws.Header.IssuedAtTime, now, jws.Header.Expiration)
	}

	// 4. Verify that jwsPayload.ksign is in jwsPayload.proofKSign.leaf
	claim, err := NewClaimFromEntry(&merkletree.Entry{Data: *jws.Payload.ProofKSign.Leaf})
	if err != nil {
		return err
	}
	claimAuthorizeKSign, ok := claim.(*ClaimAuthorizeKSignSecp256k1)
	if !ok {
		return fmt.Errorf("Invalid claim type in payload.proofksign.leaf," +
			"expected ClaimAuthorizeKSignSecp256k1")
	}
	if !reflect.DeepEqual(jws.Payload.KSign.PublicKey, *claimAuthorizeKSign.PubKey) {
		return fmt.Errorf("Pub key in payload.proofksign doesn't match payload.ksign")
	}

	// 5. Verify that jwsHeader.iss is in jwsPayload.proofKSign.
	if jws.Payload.ProofKSign.Proofs[0].Aux == nil {
		return fmt.Errorf("payload.proofksign.proofs[0].aux is nil")
	}
	if jws.Header.Issuer != jws.Payload.ProofKSign.Proofs[0].Aux.IdAddr {
		return fmt.Errorf("header.iss doesn't match with idaddr in proofksign set root claim")
	}

	// 6. Verify that signature of JWS(jwsHeader, jwsPayload) is by jwsPayload.ksign
	//
	// As verifying a signature is cheaper than verifying a merkle tree
	// proof, first we verify signature with ksign, and then we verify the
	// merkle tree proofs.
	if !utils.VerifySigEthMsg(crypto.PubkeyToAddress(jws.Payload.KSign.PublicKey),
		jws.Signature, jws.SignedBytes) {
		return fmt.Errorf("JWS signature doesn't match with pub key in payload.ksign")
	}

	// 7. VerifyProofClaim(jwsPayload.proofOfKSign, relayPk)
	if ok, err := VerifyProofClaim(relayAddr, &jws.Payload.ProofKSign); !ok {
		return fmt.Errorf("Invalid proofKSign: %v", err)
	}

	return nil
}

func VerifySignedPacket(jws *SignedPacket) error {
	switch jws.Header.Type {
	// 1. Verify jwsHeader.typ is 'iden3.sig.v0_1'
	case SIGV01:
		return VerifySignedPacketV01(jws)
	default:
		return fmt.Errorf("Unsupported signature packet typ: %v", jws.Header.Type)
	}
}

type IdenAssertResult struct {
	NonceObj *NonceObj
	EthName  string
	IdAddr   common.Address
}

func VerifyIdenAssertV01(nonceDb *NonceDb, origin string, jws *SignedPacket) (*IdenAssertResult, error) {
	data, ok := jws.Payload.Data.(IdenAssertData)
	if !ok {
		return nil, fmt.Errorf("Invalid payload.data")
	}
	form, ok := jws.Payload.Form.(IdenAssertForm)
	if !ok {
		return nil, fmt.Errorf("Invalid payload.form")
	}

	// 2. Verify jwsPayload.data.origin is origin
	if data.Origin != origin {
		return nil, fmt.Errorf("Invalid origin: expected %v, but got %v", origin, data.Origin)
	}

	// check that jwsPayload.proofKSign.proofs.length <= 2
	if len(jws.Payload.ProofKSign.Proofs) > 2 {
		return nil, fmt.Errorf("Authorize KSign claim proofs of depth > 2 not allowed yet")
	}

	// 3. Verify jwsPayload.data.challenge is in nonceDB and hasn't expired, delete it
	nonceObj, ok := nonceDb.SearchAndDelete(data.Challenge)
	if !ok {
		return nil, fmt.Errorf("Invalid nonce")
	}

	// 4. Verify that jwsHeader.iss and jwsPayload.form.ethName are in jwsPayload.proofAssignName.leaf
	claim, err := NewClaimFromEntry(&merkletree.Entry{Data: *form.ProofAssignName.Leaf})
	if err != nil {
		return nil, fmt.Errorf("Error parsing form.proofAssignNam.leaf: %v", err)
	}
	claimAssignName, ok := claim.(*ClaimAssignName)
	if !ok {
		return nil, fmt.Errorf("Invalid claim type in form.proofAssignName.leaf")
	}
	if HashName(form.EthName) != claimAssignName.NameHash {
		return nil, fmt.Errorf("Assign Name claim name doesn't match with form.ethName")
	}
	if jws.Header.Issuer != claimAssignName.IdAddr {
		return nil, fmt.Errorf("Assign Name claim idAddr doesn't match with header.iss")
	}

	// 5. VerifyProofClaim(jwsPayload.form.proofAssignName, relayPk)
	if ok, err := VerifyProofClaim(relayAddr, &form.ProofAssignName); !ok {
		return nil, fmt.Errorf("form.proofAssignName not verified: %v", err)
	}

	return &IdenAssertResult{NonceObj: nonceObj, EthName: form.EthName, IdAddr: jws.Header.Issuer}, nil
}

type RequestIdenAssertHeader struct {
	Type string `json:"typ" binding:"required"`
}

type RequestIdenAssertBody struct {
	Type string         `json:"type" binding:"required"`
	Data IdenAssertData `json:"data" binding:"required"`
}

type RequestIdenAssert struct {
	Header RequestIdenAssertHeader `json:"header" binding:"required"`
	Body   RequestIdenAssertBody   `json:"body" binding:"required"`
}

func NewRequestIdenAssert(nonceDb *NonceDb, origin string, expireDelta int64) *RequestIdenAssert {
	nonceObj := nonceDb.New(expireDelta, nil)
	return &RequestIdenAssert{
		Header: RequestIdenAssertHeader{Type: SIGV01},
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

func VerifySignedPacketIdenAssert(jws *SignedPacket, nonceDb *NonceDb, origin string) (*IdenAssertResult, error) {
	if jws.Payload.Type != IDENASSERTV01 {
		return nil, fmt.Errorf("Invalid payload.type: %v", jws.Payload.Type)
	}
	if err := VerifySignedPacket(jws); err != nil {
		return nil, err
	}
	return VerifyIdenAssertV01(nonceDb, origin, jws)
}

func VerifySignedPacketGeneric(jws *SignedPacket) error {
	if jws.Payload.Type != GENERICSIGV01 {
		return fmt.Errorf("Invalid payload.type: %v", jws.Payload.Type)
	}
	if err := VerifySignedPacket(jws); err != nil {
		return err
	}
	return nil
}
