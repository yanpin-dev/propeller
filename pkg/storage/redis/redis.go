package redis

import (
	"fmt"
	"github.com/go-redis/redis"
)

func NewClient(options *Options) (*redis.Client, error) {
	var stubClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", options.Host, options.Port),
		Password: options.Password,
		DB:       int(options.Database),
		//MaxRetries: 10,
		//DialTimeout:  config.DialTimeout,
		//ReadTimeout:  config.ReadTimeout,
		//WriteTimeout: config.WriteTimeout,
		//PoolSize:     config.PoolSize,
		//MinIdleConns: config.MinIdleConns,
		//IdleTimeout:  config.IdleTimeout,
	})
	if err := stubClient.Ping().Err(); err != nil {
		return nil, err
	}
	return stubClient, nil
}
