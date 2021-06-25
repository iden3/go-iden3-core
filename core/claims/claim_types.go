package claims

import (
	"errors"
	"fmt"
	"strings"

	"github.com/iden3/go-iden3-core/common"
	cryptoUtils "github.com/iden3/go-iden3-crypto/utils"
	"github.com/iden3/go-merkletree"
)

var (
	// ClaimTypeBasic is a simple claim type that can be used for anything.
	ClaimTypeBasic       = NewClaimTypeNum(0)
	ClaimTypeStringBasic = "Basic"

	// ClaimTypeKeyBabyJub is a claim type to autorize a babyjub public key for signing.
	ClaimTypeKeyBabyJub       = NewClaimTypeNum(1)
	ClaimTypeStringKeyBabyJub = "KeyBabyJub"

	// ClaimTypeOtherIden is a simple claim type that can be used for
	// anything with a recipient subject.
	ClaimTypeOtherIden       = NewClaimTypeNum(2)
	ClaimTypeStringOtherIden = "OtherIden"

// 	// ClaimTypeSetRootKey is a claim type of the root key of a merkle tree that goes into the relay.
// 	ClaimTypeSetRootKey = NewClaimTypeNum(2)
// 	// ClaimTypeAssignName is a claim type to assign a name to an ID
// 	ClaimTypeAssignName = NewClaimTypeNum(3)
// 	// ClaimTypeAuthorizeKSignSecp256k1 is a claim type to autorize a secp256k1 public key for signing.
// 	ClaimTypeAuthorizeKSignSecp256k1 = NewClaimTypeNum(4)
// 	// ClaimTypeLinkObjectIdentity is a claim type to link an object (represented by a hash) to an identity.
// 	ClaimTypeLinkObjectIdentity = NewClaimTypeNum(5)
// 	// ClaimTypeAuthorizeService is a claim type to authorize a Service for the identity that performs the claim
// 	ClaimTypeAuthorizeService = NewClaimTypeNum(6)
// 	// ClaimTypeNonce is a claim used to increment the tree nonce to modify the root hash
// 	ClaimTypeNonce = NewClaimTypeNum(7)
// 	// ClaimTypeEthId is a claim type to autorize an Eth Address to be used as Id inside Ethereum
// 	ClaimTypeEthId = NewClaimTypeNum(8)
// 	// ClaimTypeAuthEthKey is a claim type to authorize an Eth Address directly from a private key, allowing to specify if is used as KDisable (revoke), KReenable (recover), etc
// 	ClaimTypeAuthEthKey = NewClaimTypeNum(9)
)

func (ct ClaimType) MarshalText() ([]byte, error) {
	var str string
	switch ct {
	case ClaimTypeBasic:
		str = fmt.Sprintf("str:%v", ClaimTypeStringBasic)
	case ClaimTypeKeyBabyJub:
		str = fmt.Sprintf("str:%v", ClaimTypeStringKeyBabyJub)
	case ClaimTypeOtherIden:
		str = fmt.Sprintf("str:%v", ClaimTypeStringOtherIden)
	default:
		str = fmt.Sprintf("hex:%v", common.Hex(ct[:]))
	}
	return []byte(str), nil
}

func (ct *ClaimType) UnmarshalText(b []byte) error {
	str := string(b)
	if strings.HasPrefix(str, "str:") {
		str := strings.TrimPrefix(str, "str:")
		switch str {
		case ClaimTypeStringBasic:
			*ct = ClaimTypeBasic
		case ClaimTypeStringKeyBabyJub:
			*ct = ClaimTypeKeyBabyJub
		case ClaimTypeStringOtherIden:
			*ct = ClaimTypeOtherIden
		default:
			return fmt.Errorf("Unknown ClaimType str:%v", str)
		}
	} else if strings.HasPrefix(str, "hex:") {
		str := strings.TrimPrefix(str, "hex:")
		if err := common.HexDecodeInto(ct[:], []byte(str)); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("invalid ClaimType prefix")
	}
	return nil
}

// NewClaimFromEntry deserializes a valid claim type into a Claim.
func NewClaimFromEntry(e *merkletree.Entry) (merkletree.Entrier, error) {
	for _, elemBytes := range e.Data {
		bigint := elemBytes.BigInt()
		ok := cryptoUtils.CheckBigIntInField(bigint)
		if !ok {
			return nil, errors.New("Elements not in the Finite Field over R")
		}
	}
	var metadata Metadata
	metadata.Unmarshal(e)
	switch metadata.Type() {
	case ClaimTypeBasic:
		c := NewClaimBasicFromEntry(e)
		return c, nil
	// case *ClaimTypeAssignName:
	// 	c := NewClaimAssignNameFromEntry(e)
	// 	return c, nil
	case ClaimTypeKeyBabyJub:
		c := NewClaimKeyBabyJubFromEntry(e)
		return c, nil
	case ClaimTypeOtherIden:
		c := NewClaimOtherIdenFromEntry(e)
		return c, nil
	// case *ClaimTypeSetRootKey:
	// 	c := NewClaimSetRootKeyFromEntry(e)
	// 	return c, nil
	// case *ClaimTypeAuthorizeKSignSecp256k1:
	// 	return NewClaimAuthorizeKSignSecp256k1FromEntry(e)
	// case *ClaimTypeLinkObjectIdentity:
	// 	c := NewClaimLinkObjectIdentityFromEntry(e)
	// 	return c, nil
	// case *ClaimTypeAuthorizeService:
	// 	c := NewClaimAuthorizeServiceFromEntry(e)
	// 	return c, nil
	// case *ClaimTypeEthId:
	// 	c := NewClaimEthIdFromEntry(e)
	// 	return c, nil
	// case *ClaimTypeAuthEthKey:
	// 	c := NewClaimAuthEthKeyFromEntry(e)
	// 	return c, nil
	default:
		return nil, ErrInvalidClaimType
	}
}

var (
	ClaimHeaderBasic = ClaimHeader{
		Type:       ClaimTypeBasic,
		Subject:    ClaimSubjectSelf,
		Expiration: false,
		Version:    false}
	ClaimHeaderKeyBabyJub = ClaimHeader{
		Type:       ClaimTypeKeyBabyJub,
		Subject:    ClaimSubjectSelf,
		Expiration: false,
		Version:    false}
	ClaimHeaderOtherIden = ClaimHeader{
		Type:       ClaimTypeOtherIden,
		Subject:    ClaimSubjectOtherIden,
		SubjectPos: ClaimSubjectPosIndex,
		Expiration: false,
		Version:    false}
)

func checkHeader(header *ClaimHeader) error {
	switch header.Type {
	case ClaimTypeBasic:
		if *header != ClaimHeaderBasic {
			return fmt.Errorf("claim header for ClaimType %v is different than expected",
				ClaimTypeStringBasic)
		}
	case ClaimTypeKeyBabyJub:
		if *header != ClaimHeaderKeyBabyJub {
			return fmt.Errorf("claim header for ClaimType %v is different than expected",
				ClaimTypeStringKeyBabyJub)
		}
	case ClaimTypeOtherIden:
		if *header != ClaimHeaderOtherIden {
			return fmt.Errorf("claim header for ClaimType %v is different than expected",
				ClaimTypeStringOtherIden)
		}
	default:
	}
	return nil
}
