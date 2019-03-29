package signedpacketsrv

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/discoverysrv"
	"github.com/iden3/go-iden3/services/nameresolversrv"
	"github.com/iden3/go-iden3/utils"
)

type Service struct {
	DiscoverySrv    *discoverysrv.Service
	nameResolverSrv *nameresolversrv.Service
}

func New(discoverySrv *discoverysrv.Service, nameResolverSrv *nameresolversrv.Service) *Service {
	return &Service{DiscoverySrv: discoverySrv, nameResolverSrv: nameResolverSrv}
}

// VerifySignedPacketV01 verifies a SIGV01 signed packet.
func (ss *Service) VerifySignedPacketV01(jws *SignedPacket) error {
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
	claim, err := core.NewClaimFromEntry(&merkletree.Entry{Data: *jws.Payload.ProofKSign.Leaf})
	if err != nil {
		return err
	}
	claimAuthorizeKSign, ok := claim.(*core.ClaimAuthorizeKSignSecp256k1)
	if !ok {
		return fmt.Errorf("Invalid claim type in payload.proofksign.leaf," +
			"expected ClaimAuthorizeKSignSecp256k1")
	}
	if !reflect.DeepEqual(jws.Payload.KSign.PublicKey, *claimAuthorizeKSign.PubKey) {
		return fmt.Errorf("Pub key in payload.proofksign doesn't match payload.ksign")
	}

	// X. check that 1 <= jwsPayload.proofKSign.proofs.length <= 2
	if len(jws.Payload.ProofKSign.Proofs) < 1 {
		return fmt.Errorf("No proofs found in payload.proofKSign")
	} else if len(jws.Payload.ProofKSign.Proofs) > 2 {
		return fmt.Errorf("Authorize KSign claim proofs of depth > 2 not allowed yet")
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

	// 7a. Get the operational key from the signer and in case it's a
	// relay, check if it's trusted.
	signerIdAddr := jws.Payload.ProofKSign.Signer
	signer, err := ss.DiscoverySrv.GetEntity(signerIdAddr)
	if err != nil {
		return fmt.Errorf("Unable to get payload.proofKSign.signer entity data: %v", err)
	}
	if len(jws.Payload.ProofKSign.Proofs) > 1 {
		if !signer.Trusted.Relay {
			return fmt.Errorf("payload.proofKSign.signer is not a trusted relay")
		}
	}

	// 7b. VerifyProofClaim(jwsPayload.proofOfKSign, signerOperational)
	if ok, err := core.VerifyProofClaim(signer.OperationalAddr, &jws.Payload.ProofKSign); !ok {
		return fmt.Errorf("Invalid proofKSign: %v", err)
	}

	return nil
}

// VerifySignedPacket verifies a signed packet.
func (ss *Service) VerifySignedPacket(jws *SignedPacket) error {
	switch jws.Header.Type {
	// 1. Verify jwsHeader.typ is 'iden3.sig.v0_1'
	case SIGV01:
		return ss.VerifySignedPacketV01(jws)
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
	IdAddr   common.Address
}

// VerifyIdenAssertV01 verifies an IDENASSERTV01 payload of a signed packet.
func (ss *Service) VerifyIdenAssertV01(nonceDb *core.NonceDb, origin string,
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
		return &IdenAssertResult{NonceObj: nonceObj, EthName: nil, IdAddr: jws.Header.Issuer}, nil
	}

	// 4. Verify that jwsHeader.iss and jwsPayload.form.ethName are in jwsPayload.form.proofAssignName.leaf
	claim, err := core.NewClaimFromEntry(&merkletree.Entry{Data: *form.ProofAssignName.Leaf})
	if err != nil {
		return nil, fmt.Errorf("Error parsing form.proofAssignNam.leaf: %v", err)
	}
	claimAssignName, ok := claim.(*core.ClaimAssignName)
	if !ok {
		return nil, fmt.Errorf("Invalid claim type in form.proofAssignName.leaf")
	}
	if core.HashName(form.EthName) != claimAssignName.NameHash {
		return nil, fmt.Errorf("Assign Name claim name doesn't match with form.ethName")
	}
	if jws.Header.Issuer != claimAssignName.IdAddr {
		return nil, fmt.Errorf("Assign Name claim idAddr doesn't match with header.iss")
	}

	// 5a. Extract domain from the name
	var domain string
	if idx := strings.LastIndexByte(form.EthName, '@'); idx == -1 {
		return nil, fmt.Errorf("Invalid form.ethName %v, it doesn't containt '@'", form.EthName)
	} else {
		domain = form.EthName[idx+1 : len(form.EthName)]
	}

	// 5b. Resolve name to obtain name server idAddr and verify that it matches the signer idAddr
	if len(form.ProofAssignName.Proofs) != 1 {
		return nil, fmt.Errorf("Assign Name claim cannot be delegated to a child entity tree")
	}
	nameServerIdAddr, err := ss.nameResolverSrv.Resolve(domain)
	if err != nil {
		return nil, fmt.Errorf("Unable to resolve %v: %v", domain, err)
	}
	signerIdAddr := form.ProofAssignName.Signer
	if *nameServerIdAddr != signerIdAddr {
		return nil, fmt.Errorf("Resolved idAddr (%v) doesn't match signer idAddr (%v)",
			common3.HexEncode(nameServerIdAddr[:]), common3.HexEncode(signerIdAddr[:]))
	}

	// 5c. Get the operational key from the signer (name server).
	signer, err := ss.DiscoverySrv.GetEntity(signerIdAddr)
	if err != nil {
		return nil, fmt.Errorf("Unable to get payload.proofKSign.signer entity data: %v", err)
	}

	// 5d. VerifyProofClaim(jwsPayload.form.proofAssignName, signerOperational)
	if ok, err := core.VerifyProofClaim(signer.OperationalAddr, &jws.Payload.ProofKSign); !ok {
		return nil, fmt.Errorf("form.proofAssignName not verified: %v", err)
	}

	return &IdenAssertResult{NonceObj: nonceObj, EthName: &form.EthName, IdAddr: jws.Header.Issuer}, nil
}

// VerifySignedPacketIdenAssert verifies a signed packet and the
// IDENASSERTV01 payload of the signed packet.
func (ss *Service) VerifySignedPacketIdenAssert(jws *SignedPacket, nonceDb *core.NonceDb, origin string) (*IdenAssertResult, error) {
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
func (ss *Service) VerifySignedPacketGeneric(jws *SignedPacket) error {
	if jws.Payload.Type != GENERICSIGV01 {
		return fmt.Errorf("Invalid payload.type: %v", jws.Payload.Type)
	}
	if err := ss.VerifySignedPacket(jws); err != nil {
		return err
	}
	return nil
}
