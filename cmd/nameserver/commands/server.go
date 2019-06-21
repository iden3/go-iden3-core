package commands

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"

	// "github.com/ethereum/go-ethereum/accounts"
	// "github.com/ethereum/go-ethereum/accounts/keystore"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/iden3/go-iden3/cmd/genericserver"
	"github.com/iden3/go-iden3/cmd/nameserver/endpoint"
	"github.com/iden3/go-iden3/crypto/babyjub"
	babykeystore "github.com/iden3/go-iden3/keystore"
	"github.com/iden3/go-iden3/services/claimsrv"
	"github.com/iden3/go-iden3/services/discoverysrv"
	"github.com/iden3/go-iden3/services/nameresolversrv"
	"github.com/iden3/go-iden3/services/namesrv"
	"github.com/iden3/go-iden3/services/signedpacketsrv"
	"github.com/iden3/go-iden3/services/signsrv"
)

var ServerCommands = []cli.Command{
	{
		Name:    "init",
		Aliases: []string{},
		Usage:   "create keys and identity for the server",
		Action:  genericserver.CmdNewIdentity,
	},
	{
		Name:    "start",
		Aliases: []string{},
		Usage:   "start the server",
		Action:  cmdStart,
	},
	{
		Name:    "stop",
		Aliases: []string{},
		Usage:   "stops the server",
		Action:  cmdStop,
	},
	{
		Name:    "info",
		Aliases: []string{},
		Usage:   "server status",
		Action:  cmdInfo,
	},
}

func cmdStart(c *cli.Context) error {

	if err := genericserver.MustRead(c); err != nil {
		return err
	}

	ks, acc := genericserver.LoadKeyStore()
	ksBaby, pkc := genericserver.LoadKeyStoreBabyJub()
	defer ksBaby.Close()
	pk, err := pkc.Decompress()
	if err != nil {
		return err
	}
	client := genericserver.LoadWeb3(ks, &acc)
	storage := genericserver.LoadStorage()
	defer storage.Close()
	mt := genericserver.LoadMerkele(storage)

	rootService := genericserver.LoadRootsService(client)
	claimService := genericserver.LoadClaimService(mt, rootService, ksBaby, pk)
	nameService := LoadNameService(claimService, ksBaby, *pk, genericserver.C.Domain, genericserver.C.Namespace)
	adminService := genericserver.LoadAdminService(mt, rootService, claimService)
	nameResolveService, err := nameresolversrv.New(genericserver.C.Names.Path)
	if err != nil {
		return err
	}
	discoveryService, err := discoverysrv.New(genericserver.C.Entitites.Path)
	if err != nil {
		return err
	}
	signedPacketVerifier := signedpacketsrv.NewSignedPacketVerifier(discoveryService, nameResolveService)

	// Check for founds
	balance, err := client.BalanceAt(acc.Address)
	if err != nil {
		log.Panic(err)
	}
	log.WithFields(log.Fields{
		"balance": balance.String(),
		"address": acc.Address.Hex(),
	}).Info("Account balance retrieved")
	if balance.Int64() < 3000000 {
		log.Panic("Not enough funds in the nameserver address")
	}

	endpoint.Serve(rootService, claimService, nameService, *signedPacketVerifier, adminService)

	rootService.StopAndJoin()

	return nil
}

func postAdminApi(command string) (string, error) {

	hostport := strings.Split(genericserver.C.Server.AdminApi, ":")
	if hostport[0] == "0.0.0.0" {
		hostport[0] = "127.0.0.1"
	}
	url := "http://" + hostport[0] + ":" + hostport[1] + "/" + command

	var body bytes.Buffer
	resp, err := http.Post(url, "text/plain", &body)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func cmdStop(c *cli.Context) error {
	if err := genericserver.MustRead(c); err != nil {
		return err
	}
	output, err := postAdminApi("stop")
	if err == nil {
		log.Info("Server response: ", output)
	}
	return err
}

func cmdInfo(c *cli.Context) error {
	if err := genericserver.MustRead(c); err != nil {
		return err
	}
	output, err := postAdminApi("info")
	if err == nil {
		log.Info("Server response: ", output)
	}
	return err
}

func LoadNameService(claimservice claimsrv.Service, ks *babykeystore.KeyStore, pk babyjub.PublicKey, domain string, namespace string) namesrv.Service {
	signer := signsrv.New(ks, pk)
	return namesrv.New(claimservice, *signer, domain)
}
