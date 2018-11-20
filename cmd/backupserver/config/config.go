package config

import (
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

type Config struct {
	Server struct {
		ServiceApi string
	}
	Storage struct {
		Path string
	}
	Mongodb struct {
		Database string
		Url      string
	}
	PoW struct {
		Difficulty int
	}
}

var C Config

func MustRead(c *cli.Context) error {

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")             // adding home directory as first search path
	viper.SetEnvPrefix("backupserver2i") // so viper.AutomaticEnv will get matching envvars starting with O2M_
	viper.AutomaticEnv()                 // read in environment variables that match

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
