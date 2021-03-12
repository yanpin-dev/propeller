package db

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Options struct {
	Type     string
	Host     string
	Port     uint16
	Username string
	Password string
	Database string
	ChartSet string
}

var defaultOptions = &Options{
	Type:     "mysql",
	Host:     "127.0.0.1",
	Port:     3306,
	Username: "root",
	Password: "root",
	Database: "test",
	ChartSet: "utf8mb4",
}

func NewOptions(v *viper.Viper) (*Options, error) {
	if err := v.UnmarshalKey("datasource", defaultOptions); err != nil {
		return nil, errors.Wrap(err, "unmarshal datasource option error")
	}
	return defaultOptions, nil
}
