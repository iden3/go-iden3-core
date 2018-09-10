package config

import (
	"github.com/spf13/viper"
)

type server struct {
	Port  string
	PrivK string
}
type geth struct {
	URL string
}
type contractsAddress struct {
	Identities string
}
type config struct {
	Server           server
	Geth             geth
	ContractsAddress contractsAddress
	Domain           string
	Namespace        string
}

var C config

func MustRead(path, filename string) {
	viper.SetConfigName(filename)
	viper.AddConfigPath(path)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	err := viper.Unmarshal(&C)
	if err != nil {
		panic(err)
	}
}
