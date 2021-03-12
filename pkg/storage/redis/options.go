package redis

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Options struct {
	Host     string
	Port     uint16
	Database uint8
	Password string
}

var defaultOptions = &Options{
	Host:     "127.0.0.1",
	Port:     6379,
	Database: 0,
	Password: "",
}

func NewOptions(v *viper.Viper) (*Options, error) {
	if err := v.UnmarshalKey("redis", defaultOptions); err != nil {
		return nil, errors.Wrap(err, "unmarshal redis option error")
	}
	return defaultOptions, nil
}
