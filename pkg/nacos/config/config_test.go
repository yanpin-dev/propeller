package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtractNacosConfig(t *testing.T) {
	cfg, err := extractNacosConfig(&testRemoteProvider{
		provider:      "nacos",
		endpoint:      "http://console.nacos.io:8848/nacos?namespace=test&group=DEFAULT_GROUP",
		path:          "",
		secretKeyring: "",
	})

	assert.Nil(t, err)
	assert.Equal(t, "http", cfg.Scheme)
	assert.Equal(t, "console.nacos.io", cfg.IpAddr)
	assert.Equal(t, uint64(8848), cfg.Port)
	assert.Equal(t, "/nacos", cfg.ContextPath)
	assert.Equal(t, "test", cfg.namespace)
	assert.Equal(t, "DEFAULT_GROUP", cfg.group)
}

type testRemoteProvider struct {
	provider      string
	endpoint      string
	path          string
	secretKeyring string
}

func (p *testRemoteProvider) Provider() string {
	return p.provider
}
func (p *testRemoteProvider) Endpoint() string {
	return p.endpoint
}
func (p *testRemoteProvider) Path() string {
	return p.path
}
func (p *testRemoteProvider) SecretKeyring() string {
	return p.secretKeyring
}
