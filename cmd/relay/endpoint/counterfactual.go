package endpoint

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3/cmd/genericserver"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/services/counterfactualsrv"
	"github.com/iden3/go-iden3/utils"
)

// handlePostCounterfactualReq is the request used to create a new user tree in the relay.
type handlePostCounterfactualReq struct {
	Id core.ID
	//Operational   common.Address     `json:"operational"`
	OperationalPk *utils.PublicKey `json:"operationalpk" binding:"required"`
	Recoverer     common.Address   `json:"recoverer"`
	Revokator     common.Address   `json:"revokator"`
}

// handlePostCounterfactualRes is the response of a creation of a new counterfactual
type handlePostCounterfactualRes struct {
	Id         core.ID          `json:"id"`
	EthAddr    common.Address   `json:"ethAddr"`
	ProofClaim *core.ProofClaim `json:"proofClaim"`
}

// handleDeployCounterfactualRes is the response of a deploy of the user contract in the blockchain.
type handleDeployCounterfactualRes struct {
	IdAddr common.Address `json:"idAddr"`
	Tx     string         `json:"tx"`
}

type handleForwardCounterfactualReq struct {
	CounterfactualAddr common.Address   `json:"counterfactualAddr"`
	KSignPk            *utils.PublicKey `json:"ksignpk" binding:"required"`
	To                 common.Address   `json:"to"`
	Data               string           `json:"data"`
	Value              string           `json:"value"`
	Gas                uint64           `json:"gas"` // gaslimit
	Sig                string           `json:"sig"`
}

type handleForwardCounterfactualRes struct {
	Tx common.Hash `json:"tx"`
}

// handleCreateCounterfactual handles the creation of a new user tree from the user keys.
func handleCreateCounterfactual(c *gin.Context) {

	if genericserver.Counterfactualservice.ImplAddr() == nil {
		genericserver.Fail(c, "counterfactualservice.ImplAddr()==nil", fmt.Errorf("Implementation not set"))
		return
	}

	var idreq handlePostCounterfactualReq
	if err := c.BindJSON(&idreq); err != nil {
		genericserver.Fail(c, "cannot parse json body", err)
		return
	}

	operational := crypto.PubkeyToAddress(idreq.OperationalPk.PublicKey)
	counterfactual := &counterfactualsrv.Counterfactual{
		Operational:   operational,
		OperationalPk: idreq.OperationalPk,
		Relayer:       common.HexToAddress(genericserver.C.KeyStore.Address),
		Recoverer:     idreq.Recoverer,
		Revokator:     idreq.Revokator,
		Impl:          *genericserver.Counterfactualservice.ImplAddr(),
	}

	ethAddr, err := genericserver.Counterfactualservice.AddressOf(counterfactual)
	if err != nil {
		genericserver.Fail(c, "failed generating identity address ", err)
		return
	}

	if proofClaim, err := genericserver.Counterfactualservice.Add(idreq.Id, counterfactual); err != nil {
		genericserver.Fail(c, "failed adding identity ", err)
		return
	} else {
		c.JSON(http.StatusOK, handlePostCounterfactualRes{Id: idreq.Id, EthAddr: ethAddr, ProofClaim: proofClaim})
	}
}

// handleDeployCounterfactual handles the deploying of the user contract in the blockchain.
func handleDeployCounterfactual(c *gin.Context) {

	idAddr := common.HexToAddress(c.Param("idaddr"))
	id, err := genericserver.Counterfactualservice.Get(idAddr)
	if err != nil {
		genericserver.Fail(c, "cannot retrieve id", err)
		return
	}

	isDeployed, err := genericserver.Counterfactualservice.IsDeployed(idAddr)
	if err != nil {
		genericserver.Fail(c, "cannot retrieve deployment status", err)
		return
	}

	if isDeployed {
		genericserver.Fail(c, "already deployed", fmt.Errorf("already deployed"))
		return
	}

	addr, tx, err := genericserver.Counterfactualservice.Deploy(id)
	if err != nil {
		genericserver.Fail(c, "error deploying", err)
		return
	}
	c.JSON(http.StatusOK, handleDeployCounterfactualRes{addr, tx.Hash().Hex()})
}

type handleGetCounterfactualRes struct {
	EthAddr common.Address
	LocalDb *counterfactualsrv.Counterfactual
	Onchain *counterfactualsrv.Info
}

func handleGetCounterfactual(c *gin.Context) {
	var counterfRes handleGetCounterfactualRes

	counterfRes.EthAddr = common.HexToAddress(c.Param("ethaddr"))

	if info, err := genericserver.Counterfactualservice.Info(counterfRes.EthAddr); err == nil {
		counterfRes.Onchain = info
	}
	if id, err := genericserver.Counterfactualservice.Get(counterfRes.EthAddr); err == nil {
		counterfRes.LocalDb = id
	}
	c.JSON(http.StatusOK, counterfRes)
}

func decodeBigIntParamOrFail(c *gin.Context, param, bivalue string) *big.Int {
	value := new(big.Int)
	value, ok := value.SetString(bivalue, 10)
	if !ok {
		genericserver.Fail(c, "bad "+param+" parameter", fmt.Errorf("bad"+param+" paremeter"))
		return nil
	}
	return value
}

func decodeHexParamOrFail(c *gin.Context, param, hexvalue string) []byte {
	if !strings.HasPrefix(hexvalue, "0x") {
		genericserver.Fail(c, "bad "+param+" parameter", fmt.Errorf("bad "+param+" paremeter"))
		return nil
	}
	if hexvalue == "0x0" {
		return []byte{}
	}
	data, err := hex.DecodeString(hexvalue[2:])
	if err != nil {
		genericserver.Fail(c, "bad data parameter", err)
		return nil
	}
	return data
}

func handleForwardCounterfactual(c *gin.Context) {

	if genericserver.Counterfactualservice.ImplAddr() == nil {
		genericserver.Fail(c, "idservice.ImplAddr()==nil", fmt.Errorf("Implementation not set"))
		return
	}

	var req handleForwardCounterfactualReq
	if err := c.BindJSON(&req); err != nil {
		genericserver.Fail(c, "cannot parse json body", err)
		return
	}

	astxt, _ := json.MarshalIndent(req, "", "   ")
	fmt.Println(string(astxt))

	// idaddr := common.HexToAddress(c.Param("id"))
	id, err := core.IDFromString(c.Param("id"))
	if err := c.BindJSON(&req); err != nil {
		genericserver.Fail(c, "cannot parse id", err)
		return
	}

	var data, sig []byte
	var value *big.Int

	if data = decodeHexParamOrFail(c, "data", req.Data); data == nil {
		return
	}
	if sig = decodeHexParamOrFail(c, "sig", req.Sig); sig == nil {
		return
	}

	if value = decodeBigIntParamOrFail(c, "value", req.Value); value == nil {
		return
	}

	tx, err := genericserver.Counterfactualservice.Forward(id, req.CounterfactualAddr,
		&req.KSignPk.PublicKey,
		req.To, data, value, req.Gas, sig)

	if err != nil {
		genericserver.Fail(c, "failed to forward", err)
		return
	}

	c.JSON(http.StatusOK, handleForwardCounterfactualRes{tx})
}
