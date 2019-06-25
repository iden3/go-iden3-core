package genericserver

import (
	"strings"

	// common3 "github.com/iden3/go-iden3/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-iden3/core"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

type ContractInfo struct {
	JsonABI string
	Address string
}

type Config struct {
	Server struct {
		ServiceApi string
		AdminApi   string
	}
	Web3 struct {
		Url string
	}
	KeyStore struct {
		Path     string
		Password string
	}
	KeyStoreBaby struct {
		Path     string
		Password string
	}
	Keys struct {
		Ethereum struct {
			KDisRaw        string         `mapstructure:"kdis"`
			KDis           common.Address `mapstructure:"-"`
			KReenRaw       string         `mapstructure:"kreen"`
			KReen          common.Address `mapstructure:"-"`
			KUpdateRootRaw string         `mapstructure:"kupdateroot"`
			KUpdateRoot    common.Address `mapstructure:"-"`
		}
		BabyJub struct {
			KOpRaw string            `mapstructure:"kop"`
			KOp    babyjub.PublicKey `mapstructure:"-"`
		}
	}
	IdRaw     string  `mapstructure:"id"`
	Id        core.ID `mapstructure:"-"`
	Contracts struct {
		RootCommits   ContractInfo
		Iden3Impl     ContractInfo
		Iden3Deployer ContractInfo
		Iden3Proxy    ContractInfo
	}
	Storage struct {
		Path string
	}
	Domain    string
	Namespace string
	Names     struct {
		Path string
	}
	Entitites struct {
		Path string
	}
}

var C Config

func MustRead(c *cli.Context) error {

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")    // adding home directory as first search path
	viper.SetEnvPrefix("iden3") // so viper.AutomaticEnv will get matching envvars starting with O2M_
	viper.AutomaticEnv()        // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if c.GlobalString("config") != "" {
		viper.SetConfigFile(c.GlobalString("config"))
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&C); err != nil {
		return err
	}
	var err error
	if C.IdRaw != "" {
		if C.Id, err = core.IDFromString(C.IdRaw); err != nil {
			return err
		}
	}
	if C.Keys.BabyJub.KOpRaw != "" {
		if err := C.Keys.BabyJub.KOp.UnmarshalText([]byte(C.Keys.BabyJub.KOpRaw)); err != nil {
			return err
		}
	}
	C.Keys.Ethereum.KDis = common.HexToAddress(C.Keys.Ethereum.KDisRaw)
	C.Keys.Ethereum.KReen = common.HexToAddress(C.Keys.Ethereum.KReenRaw)
	C.Keys.Ethereum.KUpdateRoot = common.HexToAddress(C.Keys.Ethereum.KUpdateRootRaw)
	return nil
}
