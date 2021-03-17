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
	EnvType     = "CFG_TYPE"
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
		cfgType := os.Getenv(EnvType)
		v.SetConfigType(cfgType)
		v.AddRemoteProvider(provider, endpoint, path)
		if err := v.ReadRemoteConfig(); err != nil {
			return nil, err
		}
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
