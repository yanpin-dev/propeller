package health

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewOptions(t *testing.T) {
	v := viper.New()
	v.Set("health.port", 8082)
	o, err := NewOptions(v)
	assert.Nil(t, err)
	assert.Equal(t, "127.0.0.1", o.Host)
	assert.Equal(t, uint16(8082), o.Port)
}
