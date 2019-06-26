package commands

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/cmd/genericserver"
	"github.com/iden3/go-iden3/eth"
	"github.com/iden3/go-iden3/services/counterfactualsrv"
	"github.com/iden3/go-iden3/services/identitysrv"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func loadIdService() (eth.Client, identitysrv.Service, counterfactualsrv.Service) {
	ks, acc := genericserver.LoadKeyStore()
	ksBaby, pkc := genericserver.LoadKeyStoreBabyJub()
	pk, err := pkc.Decompress()
	if err != nil {
		panic(err)
	}
	client := genericserver.LoadWeb3(ks, &acc)
	client2 := genericserver.LoadEthClient2(ks, &acc)
	storage := genericserver.LoadStorage()
	mt := genericserver.LoadMerkele(storage)
	proofClaims := genericserver.LoadGenesis(mt)
	kUpdateMtp := proofClaims.KUpdateRoot.Proofs[0].Mtp0.Bytes()

	rootService := genericserver.LoadRootsService(client2, kUpdateMtp)
	claimService := genericserver.LoadClaimService(mt, rootService, ksBaby, pk)
	return client, genericserver.LoadIdentityService(claimService), genericserver.LoadCounterfactualService(client, claimService, storage)
}

var IdCommands = []cli.Command{{
	Name:    "id",
	Aliases: []string{},
	Usage:   "operate with identities",
	Subcommands: []cli.Command{
		{
			Name:   "info",
			Usage:  "show information about identity",
			Action: cmdIdInfo,
		},
		{
			Name:   "list",
			Usage:  "list identities",
			Action: cmdIdList,
		},
		{
			Name:   "add",
			Usage:  "add new identity to db",
			Action: cmdIdAdd,
		},
		{
			Name:   "deploy",
			Usage:  "deploy new identity",
			Action: cmdIdDeploy,
		},
	},
}}

// TODO this will be create a new ID, not a new counterfactual address
func cmdIdAdd(c *cli.Context) error {
	/*
		if err := genericserver.MustRead(c); err != nil {
			return err
		}

		_, idservice, counterfactualservice := loadIdService()

		if len(c.Args()) != 3 {
			return fmt.Errorf("usage: <0xoperational> <0xrecovery> <0xrevocation>")
		}

		id := counterfactualsrv.Counterfactual{
			Operational: common.HexToAddress(c.Args()[0]),
			Relayer:     common.HexToAddress(genericserver.C.KeyStore.Address),
			Recoverer:   common.HexToAddress(c.Args()[1]),
			Revokator:   common.HexToAddress(c.Args()[2]),
			Impl:        *counterfactualservice.ImplAddr(),
		}

		idaddr, err := counterfactualservice.AddressOf(&id)
		if err != nil {
			return err
		}

		if _, err = counterfactualservice.Add(&id); err != nil {
			return err
		}

		log.Info("New identity stored: ", idaddr.Hex())
	*/
	return nil
}

func cmdIdDeploy(c *cli.Context) error {

	if err := genericserver.MustRead(c); err != nil {
		return err
	}

	client, _, counterfactualservice := loadIdService()

	if len(c.Args()) != 1 {
		return fmt.Errorf("usage: <0xethAddr>")
	}
	ethAddr := common.HexToAddress(c.Args()[0])
	id, err := counterfactualservice.Get(ethAddr)
	if err != nil {
		return err
	}

	isDeployed, err := counterfactualservice.IsDeployed(ethAddr)
	if err != nil {
		return err
	}
	if isDeployed {
		log.Warn("Counterfactual already deployed at ", ethAddr.Hex())
		return nil
	}

	addr, tx, err := counterfactualservice.Deploy(id)
	if err != nil {
		_, err = client.WaitReceipt(tx.Hash())
	}
	if err != nil {
		log.Error(err)
	}
	fmt.Println(addr)
	return nil
}

type idInfo struct {
	IdAddr  common.Address
	LocalDb *counterfactualsrv.Counterfactual
	Onchain *counterfactualsrv.Info
}

func cmdIdInfo(c *cli.Context) error {

	if err := genericserver.MustRead(c); err != nil {
		return err
	}

	if len(c.Args()) != 1 {
		return fmt.Errorf("usage: <0xidaddr>")
	}

	_, _, counterfactualservice := loadIdService()

	var idi idInfo

	idi.IdAddr = common.HexToAddress(c.Args()[0])
	info, err := counterfactualservice.Info(idi.IdAddr)
	if err != nil {
		fmt.Println(err)
	} else {
		idi.Onchain = info
	}
	id, err := counterfactualservice.Get(idi.IdAddr)
	if err != nil {
		fmt.Println(err)
	} else {
		idi.LocalDb = id
	}
	text, err := json.MarshalIndent(idi, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(text))

	return nil
}
func cmdIdList(c *cli.Context) error {

	if err := genericserver.MustRead(c); err != nil {
		return err
	}

	_, _, counterfactualservice := loadIdService()
	addrs, err := counterfactualservice.List(1024)
	if err != nil {
		return err
	}
	addrsjson, err := json.MarshalIndent(addrs, "", "\t")
	if err != nil {
		return err
	} else {
		fmt.Println(string(addrsjson))
	}

	return nil
}
