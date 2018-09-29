package commands

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	cfg "github.com/iden3/go-iden3/cmd/relay/config"
	"github.com/iden3/go-iden3/services/identitysrv"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

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
			Name:   "deploy",
			Usage:  "deploy new identity",
			Action: cmdIdDeploy,
		},
	},
}}

func cmdIdDeploy(c *cli.Context) error {

	if err := cfg.MustRead(c); err != nil {
		return err
	}
	ks, acc := cfg.LoadKeyStore()
	client := cfg.LoadWeb3(ks, &acc)
	idservice := cfg.LoadIdService(client)
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

	isDeployed, err := idservice.IsDeployed(idaddr)
	if err != nil {
		return err
	}
	if isDeployed {
		log.Warn("Identity already deployed at ", idaddr.Hex())
		return nil
	}

	_, err = idservice.Deploy(&id)
	if err != nil {
		return err
	}
	return nil
}

func cmdIdInfo(c *cli.Context) error {

	if err := cfg.MustRead(c); err != nil {
		return err
	}
	ks, acc := cfg.LoadKeyStore()
	client := cfg.LoadWeb3(ks, &acc)
	idservice := cfg.LoadIdService(client)
	if len(c.Args()) != 1 {
		return fmt.Errorf("usage: <0xidaddr>")
	}

	idaddr := common.HexToAddress(c.Args()[0])
	info, err := idservice.Info(idaddr)
	if err != nil {
		return err
	}
	infojson, err := json.MarshalIndent(info, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(infojson))

	return nil
}
