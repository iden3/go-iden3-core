package endpoint

import (
	"encoding/base64"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/genericserver"
	common3 "github.com/iden3/go-iden3/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/merkletree"
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/services/notificationsrv"
	"github.com/iden3/go-iden3/utils"
)

// IdData struct representing user data that claim server will manage afterwards.
type IdData struct {
	IdAddr      common.Address `json:"idAddr"`
	NotifSrvUrl string         `json:"notifSrvUrl"`
}

type IdDataB64 IdData

// UnmarshalText retrieve data from an array of bytes.
func (d *IdDataB64) UnmarshalText(text []byte) error {
	idDataJSON, err := base64.URLEncoding.WithPadding(base64.NoPadding).
		DecodeString(string(text))
	if err != nil {
		return err
	}
	var idData IdData
	if err := json.Unmarshal(idDataJSON, &idData); err != nil {
		return err
	}
	*d = IdDataB64(idData)
	return nil
}

// claimData struct representing data needed in order to be accepted by handlePostClaim function.
type claimData struct {
	IdData IdDataB64 `json:"idData" binding:"required"`
	Cert   string    `json:"data" binding:"required"`
}

// handlePostClaim handles the request to add a claim to a user tree.
func handlePostClaim(c *gin.Context) {
	var m claimData
	if err := c.BindJSON(&m); err != nil {
		genericserver.Fail(c, "cannot parse json body", err)
		return
	}

	hash := utils.HashBytes([]byte(m.Cert))
	hashType := core.HashTypeKeccak256
	objectType := core.ObjectTypeCertificate
	indexObject := uint16(0)
	claim := core.NewClaimLinkObjectIdentity(hashType, objectType, indexObject,
		m.IdData.IdAddr, hash[:])

	// If necessary store the claim with a version higher than an existing
	// claim to invalidate the later.
	version, err := claimsrv.GetNextVersion(genericserver.Claimservice.MT(), claim.Entry().HIndex())
	if err != nil {
		genericserver.Fail(c, "error on GetNextVersion", err)
		return
	}
	claim.Version = version

	// Add claim to claim server merke tree.
	err = genericserver.Claimservice.AddClaim(claim)
	if err != nil {
		genericserver.Fail(c, "error on AddLinkObjectClaim", err)
		return
	}

	// return claim with proofs.
	proofClaim, err := genericserver.Claimservice.GetClaimProofByHi(claim.Entry().HIndex())
	if err != nil {
		genericserver.Fail(c, "error on GetClaimProofByHi", err)
		return
	}

	// Send proofClaim to notification server.
	service := notificationsrv.New(m.IdData.NotifSrvUrl, &genericserver.SignedPacketService)

	// Send packet.
	notification := notificationsrv.NewMsgProofClaim(proofClaim)
	err = service.SendNotification(notification, m.IdData.IdAddr)
	if err != nil {
		genericserver.Fail(c, "error at sending notification", err)
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
	return
}

// handleGetClaimProofByHi handles the request to query the claim proof of a
// server claim (by hIndex).
func handleGetClaimProofByHi(c *gin.Context) {
	hihex := c.Param("hi")
	hiBytes, err := common3.HexDecode(hihex)
	if err != nil {
		genericserver.Fail(c, "error on HexDecode of Hi", err)
		return
	}
	hi := &merkletree.Hash{}
	copy(hi[:], hiBytes)
	proofOfClaim, err := genericserver.Claimservice.GetClaimProofByHi(hi)
	if err != nil {
		genericserver.Fail(c, "error on GetClaimProofByHi", err)
		return
	}
	c.JSON(200, gin.H{
		"proofClaim": proofOfClaim,
	})
	return
}
