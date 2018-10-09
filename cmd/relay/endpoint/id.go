package endpoint

import (
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	cfg "github.com/iden3/go-iden3/cmd/relay/config"
	"github.com/iden3/go-iden3/services/identitysrv"
)

type handlePostIdReq struct {
	Operational common.Address `json:"operational"`
	Recoverer   common.Address `json:"recoverer"`
	Revokator   common.Address `json:"revokator"`
}

type handlePostIdRes struct {
	IDAddr common.Address `json:"idaddr"`
}

type handleDeployIdRes struct {
	IDAddr common.Address `json:"idaddr"`
	Tx     string         `json:"tx"`
}

func handleCreateId(c *gin.Context) {

	if idservice.ImplAddr() == nil {
		fail(c, "idservice.ImplAddr()==nil", fmt.Errorf("Implementation not set"))
		return
	}

	var idreq handlePostIdReq
	if err := c.BindJSON(&idreq); err != nil {
		fail(c, "cannot parse json body", err)
		return
	}

	id := &identitysrv.Identity{
		Operational: idreq.Operational,
		Relayer:     common.HexToAddress(cfg.C.KeyStore.Address),
		Recoverer:   idreq.Recoverer,
		Revokator:   idreq.Revokator,
		Impl:        *idservice.ImplAddr(),
	}

	addr, err := idservice.AddressOf(id)
	if err != nil {
		fail(c, "failed generating identity address ", err)
		return
	}

	if err := idservice.Add(id); err != nil {
		fail(c, "failed adding identity ", err)
		return
	}

	c.JSON(http.StatusOK, handlePostIdRes{addr})
}

func handleDeployId(c *gin.Context) {

	idaddr := common.HexToAddress(c.Param("idaddr"))
	id, err := idservice.Get(idaddr)
	if err != nil {
		fail(c, "cannot retrieve idaddr", err)
		return
	}

	isDeployed, err := idservice.IsDeployed(idaddr)
	if err != nil {
		fail(c, "cannot retrieve deployment status", err)
		return
	}

	if isDeployed {
		fail(c, "already deployed", fmt.Errorf("already deployed"))
		return
	}

	addr, tx, err := idservice.Deploy(id)
	if err != nil {
		fail(c, "error deploying", err)
		return
	}
	c.JSON(http.StatusOK, handleDeployIdRes{addr, tx.Hash().Hex()})
}

type handleGetIdRes struct {
	IDAddr  common.Address
	LocalDb *identitysrv.Identity
	Onchain *identitysrv.Info
}

func handleGetId(c *gin.Context) {
	var idi handleGetIdRes

	idi.IDAddr = common.HexToAddress(c.Param("idaddr"))

	if info, err := idservice.Info(idi.IDAddr); err == nil {
		idi.Onchain = info
	}
	if id, err := idservice.Get(idi.IDAddr); err == nil {
		idi.LocalDb = id
	}
	c.JSON(http.StatusOK, idi)
}
