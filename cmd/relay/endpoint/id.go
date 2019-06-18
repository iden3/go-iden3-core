package endpoint

import (
	"net/http"

	// "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/genericserver"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/crypto/babyjub"
	// "github.com/iden3/go-iden3/utils"
)

type handleIdGenesis struct {
	KOp   			*babyjub.PublicKey `json:"operationalPk" binding:"required"`
	KDis  			common.Address     `json:"kdisable" binding:"required"`
	KReen 			common.Address     `json:"kreenable" binding:"required"`
	KUpdateRoot common.Address     `json:"kupdateRoot" binding:"required"`
}

// handlePostIdRes is the response of a creation of a new user tree in the relay.
type handlePostIdRes struct {
	Id         core.ID          `json:"id"`
	ProofClaim *core.ProofClaim `json:"proofClaim"`
}

// handleCreateIdGenesis creates the identity creating a new MerkleTree that contains
// the initial keys of that identity. The Merkle Root of that tree will be the
// identity address
func handleCreateIdGenesis(c *gin.Context) {
	var idgen handleIdGenesis
	if err := c.BindJSON(&idgen); err != nil {
		genericserver.Fail(c, "cannot parse json body", err)
		return
	}

	id, proofKOp, err := genericserver.Idservice.CreateIdGenesis(idgen.KOp, idgen.KDis, idgen.KReen, idgen.KUpdateRoot)
	if err != nil {
		genericserver.Fail(c, "failed generating identity address ", err)
		return
	}

	c.JSON(http.StatusOK, handlePostIdRes{Id: *id, ProofClaim: proofKOp})
}
