package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

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
		RootCommits struct {
			JsonABI string
			Address string
		}
		Iden3Impl struct {
			JsonABI string
			Address string
		}
		Iden3Deployer struct {
			JsonABI string
			Address string
		}
		Iden3Proxy struct {
			JsonABI string
		}
	}
	Storage struct {
		Path string
	}
	Domain    string
	Namespace string
}

var C Config

func MustRead(path, filename string) {
	viper.SetConfigName(filename)
	viper.AddConfigPath(path)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Panic(err)
	}
	err := viper.Unmarshal(&C)
	if err != nil {
		log.Panic(err)
	}
}
