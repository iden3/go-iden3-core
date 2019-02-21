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
	//common3 "github.com/iden3/go-iden3/common"
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
const SIGALGV01 = "EK256K1"

// JwsHeader is the JSON Web Signature Header
type JwsHeader struct {
	Type         string         `json:"type" binding:"required"`
	Issuer       common.Address `json:"iss" binding:"required"`
	IssuedAtTime int64          `json:"iat" binding:"required"`
	Expiration   int64          `json:"exp" binding:"required"`
	Algorithm    string         `json:"alg" binding:"required"`
}

// JwsPayload is the JSON Web Signature Payload
type JwsPayload struct {
	Type       string           `json:"type" binding:"required"`
	Data       interface{}      `json:"data" binding:"required"`
	KSign      *utils.PublicKey `json:"ksign" binding:"required"`
	ProofKSign ProofClaim       `json:"proofKSign" binding:"required"`
	Form       interface{}      `json:"form" binding:"required"`
}

// Jws is a JSON Web Signature unmarshaled packet
type Jws struct {
	Header    JwsHeader
	Payload   JwsPayload
	Signature utils.SignatureEthMsg
}

func (j *Jws) MarshalSign(ks *keystore.KeyStore, addr common.Address) (string, error) {
	headerJSON, err := json.Marshal(j.Header)
	if err != nil {
		return "", err
	}
	payloadJSON, err := json.Marshal(j.Payload)
	if err != nil {
		return "", err
	}
	dataToSign := fmt.Sprintf("%v.%v", base64.StdEncoding.EncodeToString([]byte(headerJSON)),
		base64.StdEncoding.EncodeToString([]byte(payloadJSON)))
	sig, err := utils.SignEthMsg(ks, accounts.Account{Address: addr}, []byte(dataToSign))
	if err != nil {
		return "", err
	}
	sig64 := base64.StdEncoding.EncodeToString(sig[:])
	return fmt.Sprintf("%v.%v", dataToSign, sig64), nil
}

func (j *Jws) Unmarshal(s string) error {
	fields := strings.Split(s, ".")
	if len(fields) != 3 {
		return fmt.Errorf("Invalid JWT: it doesn't contain 3 dot separated fields")
	}
	jwsHeader64, jwsPayload64, signature64 := fields[0], fields[1], fields[2]
	jwsHeader, err := base64.StdEncoding.DecodeString(jwsHeader64)
	jwsPayload, err := base64.StdEncoding.DecodeString(jwsPayload64)
	if err := json.Unmarshal(jwsHeader, &j.Header); err != nil {
		return err
	}
	if err := json.Unmarshal(jwsPayload, &j.Payload); err != nil {
		return err
	}
	signature, err := base64.StdEncoding.DecodeString(signature64)
	if err != nil {
		return err
	}
	copy(j.Signature[:], signature)
	return nil
}

func SignPacketV01(ks *keystore.KeyStore, idAddr common.Address, kSignPk *ecdsa.PublicKey, proofKSign ProofClaim,
	expireDelta int64, payloadType string, data interface{}, form interface{}) (string, error) {
	now := time.Now().Unix()
	header := JwsHeader{
		Type:         SIGV01,
		Issuer:       idAddr,
		IssuedAtTime: now,
		Expiration:   now + expireDelta,
		Algorithm:    SIGALGV01,
	}
	payload := JwsPayload{
		Type:       IDENASSERTV01,
		Data:       data,
		KSign:      &utils.PublicKey{PublicKey: *kSignPk},
		ProofKSign: proofKSign,
		Form:       form,
	}
	jws := Jws{Header: header, Payload: payload}
	signedPacket, err := jws.MarshalSign(ks, crypto.PubkeyToAddress(*kSignPk))
	if err != nil {
		return "", err
	}
	return signedPacket, nil
}

func VerifySignedPacketV01(signedPacket string, jws Jws) error {
	// 2. Verify jwsHeader.alg is 'ES255'
	if jws.Header.Algorithm != SIGALGV01 {
		return fmt.Errorf("Unsupported alg: %v", jws.Header.Algorithm)
	}

	// 3. Verify that jwsHeader.iat <= now() < jwsHeader.exp
	now := time.Now().Unix()
	if !((jws.Header.IssuedAtTime <= now) && (now < jws.Header.Expiration)) {
		return fmt.Errorf("Signature not valid for current date")
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
	msg := signedPacket[:strings.LastIndex(signedPacket, ".")]
	if !utils.VerifySigEthMsg(crypto.PubkeyToAddress(jws.Payload.KSign.PublicKey),
		&jws.Signature, []byte(msg)) {
		return fmt.Errorf("JWS signature doesn't match with pub key in payload.ksign")
	}

	// 7. VerifyProofClaim(jwsPayload.proofOfKSign, relayPk)
	if ok, err := VerifyProofClaim(relayAddr, &jws.Payload.ProofKSign); !ok {
		return err
	}

	return nil
}

func VerifySignedPacket(signedPacket string) error {
	var jws Jws
	if err := jws.Unmarshal(signedPacket); err != nil {
		return err
	}
	switch jws.Header.Type {
	// 1. Verify jwsHeader.typ is 'iden3.sig.v0_1'
	case SIGV01:
		return VerifySignedPacketV01(signedPacket, jws)
	default:
		return fmt.Errorf("Unsupported signature packet typ: %v", jws.Header.Type)
	}
}
