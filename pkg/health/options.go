package health

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Options struct {
	Host         string
	Port         uint16
	ReadTimeout  int64
	WriteTimeout int64
}

var defaultOptions = &Options{
	Host:         "127.0.0.1",
	Port:         8081,
	ReadTimeout:  10,
	WriteTimeout: 10,
}

func NewOptions(v *viper.Viper) (*Options, error) {
	if err := v.UnmarshalKey("health", defaultOptions); err != nil {
		return nil, errors.Wrap(err, "unmarshal server option error")
	}
	return defaultOptions, nil
}
