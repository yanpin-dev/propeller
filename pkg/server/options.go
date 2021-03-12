package server

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Options struct {
	Host string
	Port uint16
}

var defaultOptions = &Options{
	Host: "0.0.0.0",
	Port: 8080,
}

func NewOptions(v *viper.Viper) (*Options, error) {
	if err := v.UnmarshalKey("server", defaultOptions); err != nil {
		return nil, errors.Wrap(err, "unmarshal server option error")
	}
	return defaultOptions, nil
}
