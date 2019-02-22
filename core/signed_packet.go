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
	Data       interface{}      `json:"data" binding:"required"`
	KSign      *utils.PublicKey `json:"ksign" binding:"required"`
	ProofKSign ProofClaim       `json:"proofKSign" binding:"required"`
	Form       interface{}      `json:"form" binding:"required"`
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
		return fmt.Errorf("Invalid JWT: it doesn't contain 3 dot separated fields")
	}
	jwsHeader64, jwsPayload64, signature64 := fields[0], fields[1], fields[2]
	jwsHeader, err := base64.StdEncoding.DecodeString(jwsHeader64)
	jwsPayload, err := base64.StdEncoding.DecodeString(jwsPayload64)
	if err := json.Unmarshal(jwsHeader, &sp.Header); err != nil {
		return err
	}
	if err := json.Unmarshal(jwsPayload, &sp.Payload); err != nil {
		return err
	}
	signature, err := base64.StdEncoding.DecodeString(signature64)
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
	sp.Unmarshal(jws.Jws)
	return nil
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
		Type:       IDENASSERTV01,
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
		return err
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
