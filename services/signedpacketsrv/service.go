package signedpacketsrv

import (
	// "encoding/hex"
	"fmt"
	"reflect"
	"strings"
	"time"

	// "github.com/ethereum/go-ethereum/crypto"
	common3 "github.com/iden3/go-iden3-core/common"
	"github.com/iden3/go-iden3-core/core"
	babykeystore "github.com/iden3/go-iden3-core/keystore"
	// "github.com/iden3/go-iden3-core/merkletree"
	"github.com/iden3/go-iden3-core/services/discoverysrv"
	"github.com/iden3/go-iden3-core/services/ethsrv"
	"github.com/iden3/go-iden3-core/services/nameresolversrv"

	// "github.com/iden3/go-iden3-core/utils"
	"github.com/iden3/go-iden3-crypto/babyjub"
)

type SignedPacketVerifier struct {
	discoverySrv    *discoverysrv.Service
	nameResolverSrv *nameresolversrv.Service
	ethSrv          ethsrv.Service
}

func NewSignedPacketVerifier(discoverySrv *discoverysrv.Service,
	nameResolverSrv *nameresolversrv.Service, ethSrv ethsrv.Service) *SignedPacketVerifier {
	return &SignedPacketVerifier{discoverySrv: discoverySrv, nameResolverSrv: nameResolverSrv, ethSrv: ethSrv}
}

// VerifySignedPacketV02 verifies a SIGV02 signed packet.
func (ss *SignedPacketVerifier) VerifySignedPacketV02(jws *SignedPacket) error {
	// 2. Verify jwsHeader.alg is 'ED256BJ'
	if jws.Header.Algorithm != SIGALGV02 {
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
	claim, err := core.NewClaimFromEntry(jws.Payload.ProofKSign.Claim)
	if err != nil {
		return err
	}
	claimAuthorizeKSign, ok := claim.(*core.ClaimAuthorizeKSignBabyJub)
	if !ok {
		return fmt.Errorf("Invalid claim type in payload.proofksign.leaf," +
			"expected ClaimAuthorizeKSignBabyJub")
	}
	claimAuthorizeKSignPkComp := babyjub.PublicKeyComp(
		babyjub.PackPoint(claimAuthorizeKSign.Ay, claimAuthorizeKSign.Sign))
	if !reflect.DeepEqual(jws.Payload.KSign.Compress(), claimAuthorizeKSignPkComp) {
		return fmt.Errorf("Pub key in payload.proofksign doesn't match payload.ksign")
	}

	if !jws.Header.Issuer.Equals(jws.Payload.ProofKSign.ID) {
		return fmt.Errorf("header.iss doesn't match with id in proofksign set root claim")
	}

	// 6. Verify that signature of JWS(jwsHeader, jwsPayload) is by jwsPayload.ksign
	//
	// As verifying a signature is cheaper than verifying a merkle tree
	// proof, first we verify signature with ksign, and then we verify the
	// merkle tree proofs.
	kSignComp := jws.Payload.KSign.Compress()
	if ok, err := babykeystore.VerifySignatureRaw(&kSignComp, jws.Signature, jws.SignedBytes); !ok {
		return fmt.Errorf("JWS signature doesn't match with pub key in payload.ksign: %v", err)
	}

	// 7b. VerifyProofClaim(jwsPayload.proofOfKSign)
	if ok, err := ss.ethSrv.VerifyProofClaim(&jws.Payload.ProofKSign); !ok {
		return fmt.Errorf("Invalid proofKSign: %v", err)
	}

	return nil
}

// VerifySignedPacket verifies a signed packet.
func (ss *SignedPacketVerifier) VerifySignedPacket(jws *SignedPacket) error {
	switch jws.Header.Type {
	// 1. Verify jwsHeader.typ is 'iden3.sig.v0_1'
	case SIGV01:
		// return ss.VerifySignedPacketV01(jws)
		return fmt.Errorf("Deprecated signature packet typ: %v", jws.Header.Type)
	// 1. Verify jwsHeader.typ is 'iden3.sig.v0_2'
	case SIGV02:
		return ss.VerifySignedPacketV02(jws)
	default:
		return fmt.Errorf("Unsupported signature packet typ: %v", jws.Header.Type)
	}
}

// IdenAssertResult is the result of a successfull verification of an
// IDENASSERTV01 payload from a signed packet.  EthName will be nil if no name
// ownership was proved (the form field of the signed packet was nil).
type IdenAssertResult struct {
	NonceObj *core.NonceObj
	EthName  *string
	Id       core.ID
}

// VerifyIdenAssertV01 verifies an IDENASSERTV01 payload of a signed packet.
func (ss *SignedPacketVerifier) VerifyIdenAssertV01(nonceDb *core.NonceDb, origin string,
	jws *SignedPacket) (*IdenAssertResult, error) {
	data, ok := jws.Payload.Data.(IdenAssertData)
	if !ok {
		return nil, fmt.Errorf("Invalid payload.data")
	}
	form, ok := jws.Payload.Form.(*IdenAssertForm)
	if !ok {
		return nil, fmt.Errorf("Invalid payload.form")
	}

	// 2. Verify jwsPayload.data.origin is origin
	if data.Origin != origin {
		return nil, fmt.Errorf("Invalid origin: expected %v, but got %v", origin, data.Origin)
	}

	// 3. Verify jwsPayload.data.challenge is in nonceDB and hasn't expired, delete it
	nonceObj, ok := nonceDb.SearchAndDelete(data.Challenge)
	if !ok {
		return nil, fmt.Errorf("Invalid nonce")
	}

	if form == (*IdenAssertForm)(nil) {
		return &IdenAssertResult{NonceObj: nonceObj, EthName: nil, Id: jws.Header.Issuer}, nil
	}

	// 4. Verify that jwsHeader.iss and jwsPayload.form.ethName are in jwsPayload.form.proofAssignName.leaf
	claim, err := core.NewClaimFromEntry(form.ProofAssignName.Claim)
	if err != nil {
		return nil, fmt.Errorf("Error parsing form.proofAssignNam.leaf: %v", err)
	}
	claimAssignName, ok := claim.(*core.ClaimAssignName)
	if !ok {
		return nil, fmt.Errorf("Invalid claim type in form.proofAssignName.leaf")
	}
	if core.HashString(form.EthName) != claimAssignName.NameHash {
		return nil, fmt.Errorf("Assign Name claim name doesn't match with form.ethName")
	}
	if jws.Header.Issuer != claimAssignName.Id {
		return nil, fmt.Errorf("Assign Name claim id doesn't match with header.iss")
	}

	// 5a. Extract domain from the name
	var domain string
	if idx := strings.LastIndexByte(form.EthName, '@'); idx == -1 {
		return nil, fmt.Errorf("Invalid form.ethName %v, it doesn't containt '@'", form.EthName)
	} else {
		domain = form.EthName[idx+1 : len(form.EthName)]
	}

	// 5b. Resolve name to obtain name server id and verify that it matches the signer id
	if form.ProofAssignName.RelayAux != nil {
		return nil, fmt.Errorf("Assign Name claim cannot be delegated to a child entity tree")
	}
	nameServerId, err := ss.nameResolverSrv.Resolve(domain)
	if err != nil {
		return nil, fmt.Errorf("Unable to resolve %v: %v", domain, err)
	}
	signerId := form.ProofAssignName.ID
	if !nameServerId.Equals(signerId) {
		return nil, fmt.Errorf("Resolved id (%v) doesn't match signer id (%v)",
			common3.HexEncode(nameServerId[:]), common3.HexEncode(signerId[:]))
	}

	// 5d. VerifyProofClaim(jwsPayload.form.proofAssignName, signerOperational)
	if ok, err := ss.ethSrv.VerifyProofClaim(&jws.Payload.ProofKSign); !ok {
		return nil, fmt.Errorf("form.proofAssignName not verified: %v", err)
	}

	return &IdenAssertResult{NonceObj: nonceObj, EthName: &form.EthName, Id: jws.Header.Issuer}, nil
}

// VerifySignedPacketIdenAssert verifies a signed packet and the
// IDENASSERTV01 payload of the signed packet.
func (ss *SignedPacketVerifier) VerifySignedPacketIdenAssert(jws *SignedPacket, nonceDb *core.NonceDb, origin string) (*IdenAssertResult, error) {
	if jws.Payload.Type != IDENASSERTV01 {
		return nil, fmt.Errorf("Invalid payload.type: %v", jws.Payload.Type)
	}
	if err := ss.VerifySignedPacket(jws); err != nil {
		return nil, err
	}
	return ss.VerifyIdenAssertV01(nonceDb, origin, jws)
}

// VerifySignedPacketGeneric verifies a signed packet and checks that the
// payload type is GENERICSIGV01.
func (ss *SignedPacketVerifier) VerifySignedPacketGeneric(jws *SignedPacket) error {
	if jws.Payload.Type != GENERICSIGV01 {
		return fmt.Errorf("Invalid payload.type: %v", jws.Payload.Type)
	}
	if err := ss.VerifySignedPacket(jws); err != nil {
		return err
	}
	return nil
}
