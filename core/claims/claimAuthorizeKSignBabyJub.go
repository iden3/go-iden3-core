package claims

import (
	"github.com/iden3/go-iden3-crypto/babyjub"
)

// NewClaimAuthorizeKSignBabyJub returns a ClaimAuthorizeKSignBabyJub with the
// given elliptic public key parameters.
func NewClaimAuthorizeKSignBabyJub(pk *babyjub.PublicKey) *ClaimKeyBabyJub {
	return &ClaimKeyBabyJub{
		metadata: NewMetadata(ClaimHeaderAuthorizeKSignBabyJub),
		Sign:     babyjub.PointCoordSign(pk.X),
		Ay:       pk.Y,
	}
}
