package event

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Options struct {
	Source   string
	Exchange string
}

func NewOptions(v *viper.Viper) (*Options, error) {
	o := new(Options)
	if err := v.UnmarshalKey("event.publisher", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal publisher option error")
	}
	return o, nil
}
