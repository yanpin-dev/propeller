package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

const (
	DefaultConfigFile = "config/config.yml"
	EnvKey            = "CFG_FILE"
	Prefix            = "propeller"
)

func NewViper() (*viper.Viper, error) {
	var (
		err error
		v   = viper.New()
	)

	file := os.Getenv(EnvKey)
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

	return v.Sub(Prefix), err
}
