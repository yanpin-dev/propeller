package config

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/yanpin-dev/propeller/pkg/nacos/config"
	"os"
)

const (
	DefaultConfigFile = "config/config.yml"
	EnvCfgFile        = "CFG_FILE"
	Prefix            = "propeller"

	EnvNacosEnabled  = "NACOS_ENABLED"
	EnvNacosEndpoint = "NACOS_ENDPOINT"
	EnvNacosPath     = "NACOS_PATH"
)

func NewViper() (*viper.Viper, error) {
	var (
		err error
		v   = viper.New()
	)

	nacosEnabled := os.Getenv(EnvNacosEnabled)
	if nacosEnabled == "true" {
		endpoint := os.Getenv(EnvNacosEndpoint)
		path := os.Getenv(EnvNacosPath)
		v.AddRemoteProvider(config.ProviderName, endpoint, path)
	} else {
		file := os.Getenv(EnvCfgFile)
		if file == "" {
			file = DefaultConfigFile
		}
		v.AddConfigPath(".")
		v.SetConfigFile(file)
	}

	if err := v.ReadInConfig(); err == nil {
		fmt.Printf("use config file -> %s\n", v.ConfigFileUsed())
	} else {
		return nil, err
	}

	return v.Sub(Prefix), err
}
