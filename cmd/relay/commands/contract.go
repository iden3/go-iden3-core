package commands

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	cfg "github.com/iden3/go-iden3/cmd/relay/config"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var ContractCommands = []cli.Command{{
	Name:  "contract",
	Usage: "operate with contracts",
	Subcommands: []cli.Command{
		{
			Name:   "info",
			Usage:  "show information about contracts",
			Action: cmdContractInfo,
		},
		{
			Name:   "deploy",
			Usage:  "deploy contract",
			Action: cmdContractDeploy,
		},
	},
}}

func contractInfo() map[string]cfg.ContractInfo {
	var info map[string]cfg.ContractInfo = make(map[string]cfg.ContractInfo)
	info["rootcommits"] = cfg.C.Contracts.RootCommits
	info["iden3impl"] = cfg.C.Contracts.Iden3Impl
	info["iden3deployer"] = cfg.C.Contracts.Iden3Deployer
	return info
}

func cmdContractInfo(c *cli.Context) error {

	if err := cfg.MustRead(c); err != nil {
		return err
	}
	ks, acc := cfg.LoadKeyStore()
	client := cfg.LoadWeb3(ks, &acc)

	info := func(name string, info cfg.ContractInfo) {
		if len(info.Address) > 0 {
			code, err := client.CodeAt(common.HexToAddress(info.Address))
			if err != nil {
				log.Panic(err)
			}
			if len(code) > 0 {
				log.Info(name, ": code set at ", info.Address)
			} else {
				log.Info(name, ": code NOT set at ", info.Address)
			}
		} else {
			log.Info(name, ": address not set")
		}
	}

	for k, v := range contractInfo() {
		info(k, v)
	}

	return nil
}

func cmdContractDeploy(c *cli.Context) error {

	if err := cfg.MustRead(c); err != nil {
		return err
	}

	ks, acc := cfg.LoadKeyStore()
	client := cfg.LoadWeb3(ks, &acc)

	if len(c.Args()) != 1 {
		return fmt.Errorf("should specify contract")
	}

	contractid := c.Args()[0]
	info, ok := contractInfo()[contractid]
	if !ok {
		return fmt.Errorf("contract %v does not exist", contractid)
	}
	contract := cfg.LoadContract(client, info.JsonABI, nil)

	_, _, err := contract.DeploySync()
	if err != nil {
		return err
	}

	log.Info("Contract deployed at ", contract.Address().Hex())

	return nil
}
