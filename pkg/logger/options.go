package logger

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Options struct {
	Use         string
	Environment string
	LogLevel    string
}

var defaultOptions = &Options{
	Use:         "zapLogger",
	Environment: "prod",
	LogLevel:    "info",
}

func NewOptions(v *viper.Viper) (*Options, error) {
	if err := v.UnmarshalKey("logger", defaultOptions); err != nil {
		return nil, errors.Wrap(err, "unmarshal logger option error")
	}
	return defaultOptions, nil
}
