package checker

import (
	"github.com/yanpin-dev/propeller/pkg/health"
	"github.com/go-redis/redis"
	"strings"
)

// Redis is a interface used to abstract the access of the Version string
type Redis interface {
	GetVersion() (string, error)
}

// Checker is a checker that check a given redis
type RedisChecker struct {
	name   string
	client *redis.Client
}

// NewRedisChecker returns a new redis.Checker
func NewRedisChecker(name string, client *redis.Client) health.Checker {
	return RedisChecker{
		name:   name,
		client: client,
	}
}

// NewDefaultRedisChecker returns a new redis.Checker configured with a custom Redis implementation
func NewDefaultRedisChecker(client *redis.Client) health.Checker {
	return NewRedisChecker("redis", client)
}

// Check obtain the version string from redis info command
func (c RedisChecker) Check() health.Health {
	health := health.NewHealth()
	c.client.Info()

	version, err := c.getVersion()

	if err != nil {
		health.Down().AddInfo("error", err.Error())
		return health
	}

	health.Up().AddInfo("version", version)

	return health
}

func (c RedisChecker) Name() string {
	return c.name
}

func (c RedisChecker) getVersion() (string, error) {
	result, err := c.client.Info("server").Result()
	if err != nil {
		return "", err
	}

	lines := strings.Split(result, "\n")
	for _, line := range lines {
		values := strings.SplitN(line, ":", 2)
		if values[0] == "redis_version" {
			return strings.TrimSpace(values[1]), nil
		}
	}
	return "unknown", nil
}
