package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

const (
	DefaultConfigFile = "config/config.yml"
	EnvCfgFile        = "CFG_FILE"
	Prefix            = "propeller"

	EnvProvider = "CFG_REMOTE_PROVIDER"
	EnvEndpoint = "CFG_REMOTE_ENDPOINT"
	EnvPath     = "CFG_REMOTE_PATH"
)

func NewViper() (*viper.Viper, error) {
	var (
		err error
		v   = viper.New()
	)

	provider := os.Getenv(EnvProvider)
	if provider != "" {
		endpoint := os.Getenv(EnvEndpoint)
		path := os.Getenv(EnvPath)
		v.AddRemoteProvider(provider, endpoint, path)
		v.ReadRemoteConfig()
	} else {
		file := os.Getenv(EnvCfgFile)
		if file == "" {
			file = DefaultConfigFile
		}
		v.AddConfigPath(".")
		v.SetConfigFile(file)
		if err := v.ReadInConfig(); err == nil {
			fmt.Printf("use config file -> %s\n", v.ConfigFileUsed())
		} else {
			return nil, err
		}
	}

	return v.Sub(Prefix), err
}
