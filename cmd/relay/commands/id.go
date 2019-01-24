package commands

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	cfg "github.com/iden3/go-iden3/cmd/relay/config"
	"github.com/iden3/go-iden3/eth"
	"github.com/iden3/go-iden3/services/identitysrv"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func loadIdService() (eth.Client, identitysrv.Service) {
	ks, acc := cfg.LoadKeyStore()
	client := cfg.LoadWeb3(ks, &acc)
	storage := cfg.LoadStorage()
	mt := cfg.LoadMerkele(storage)
	rootservice := cfg.LoadRootsService(client)
	claimservice := cfg.LoadClaimService(mt, rootservice, ks, acc)
	return client, cfg.LoadIdService(client, claimservice, storage)
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

func cmdIdAdd(c *cli.Context) error {

	if err := cfg.MustRead(c); err != nil {
		return err
	}

	_, idservice := loadIdService()

	if len(c.Args()) != 3 {
		return fmt.Errorf("usage: <0xoperational> <0xrecovery> <0xrevocation>")
	}

	id := identitysrv.Identity{
		Operational: common.HexToAddress(c.Args()[0]),
		Relayer:     common.HexToAddress(cfg.C.KeyStore.Address),
		Recoverer:   common.HexToAddress(c.Args()[1]),
		Revokator:   common.HexToAddress(c.Args()[2]),
		Impl:        *idservice.ImplAddr(),
	}

	idaddr, err := idservice.AddressOf(&id)
	if err != nil {
		return err
	}

	if _, err = idservice.Add(&id); err != nil {
		return err
	}

	log.Info("New identity stored: ", idaddr.Hex())

	return nil
}

func cmdIdDeploy(c *cli.Context) error {

	if err := cfg.MustRead(c); err != nil {
		return err
	}

	client, idservice := loadIdService()

	if len(c.Args()) != 1 {
		return fmt.Errorf("usage: <0xidaddr>")
	}
	idaddr := common.HexToAddress(c.Args()[0])
	id, err := idservice.Get(idaddr)
	if err != nil {
		return err
	}

	isDeployed, err := idservice.IsDeployed(idaddr)
	if err != nil {
		return err
	}
	if isDeployed {
		log.Warn("Identity already deployed at ", idaddr.Hex())
		return nil
	}

	addr, tx, err := idservice.Deploy(id)
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
	IDAddr  common.Address
	LocalDb *identitysrv.Identity
	Onchain *identitysrv.Info
}

func cmdIdInfo(c *cli.Context) error {

	if err := cfg.MustRead(c); err != nil {
		return err
	}

	if len(c.Args()) != 1 {
		return fmt.Errorf("usage: <0xidaddr>")
	}

	_, idservice := loadIdService()

	var idi idInfo

	idi.IDAddr = common.HexToAddress(c.Args()[0])
	info, err := idservice.Info(idi.IDAddr)
	if err != nil {
		fmt.Println(err)
	} else {
		idi.Onchain = info
	}
	id, err := idservice.Get(idi.IDAddr)
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

	if err := cfg.MustRead(c); err != nil {
		return err
	}

	_, idservice := loadIdService()
	addrs, err := idservice.List(1024)
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
