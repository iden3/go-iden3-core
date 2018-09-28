package config

import (
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
		Address  string
		Password string
	}
	Contracts struct {
		RootCommits   ContractInfo
		Iden3Impl     ContractInfo
		Iden3Deployer ContractInfo
		Iden3Proxy    ContractInfo
	}
	Namespace string
	Jwtsecret string
}

var C Config

func MustRead(c *cli.Context) error {

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")      // adding home directory as first search path
	viper.SetEnvPrefix("relay2i") // so viper.AutomaticEnv will get matching envvars starting with O2M_
	viper.AutomaticEnv()          // read in environment variables that match

	if c.GlobalString("config") != "" {
		viper.SetConfigFile(c.GlobalString("config"))
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	err := viper.Unmarshal(&C)
	if err != nil {
		return err
	}
	return nil
}
