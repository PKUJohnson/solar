package toolkit

import (
	"fmt"

	"strconv"
	"time"

	std "github.com/PKUJohnson/solar/std"
	"gopkg.in/redis.v5"
)

// InitRedis initializes the redis connection.
// Usage:
// client := InitRedis(config)
func InitRedis(config std.ConfigRedis) *redis.Client {
	if !config.Enabled {
		return nil
	}

	hostAddress := fmt.Sprintf("%s:%d", config.Host, config.Port)
	_client := redis.NewClient(&redis.Options{
		Addr:     hostAddress,
		Password: config.Auth, // no password set
		DB:       0,           // use default DB
		PoolSize: 100,
	})
	if _, err := _client.Ping().Result(); err != nil {
		std.LogError(std.LogFields{"args": []interface{}{config}, "err": err}, "Failed to connect Redis.")
		_client = nil
	}

	return _client
}

func InitRedisRing(config ...std.ConfigRedis) *redis.Ring {

	servers := map[string]string{}
	var pass string
	var timeoutRead, timeoutWrite, timeoutConnect, timeoutIdle time.Duration

	for i, v := range config {
		servers[strconv.Itoa(i)] = fmt.Sprintf("%s:%d", v.Host, v.Port)
		if i == 0 {
			pass = v.Auth
			timeoutIdle = time.Second * time.Duration(v.IdleTimeout)
		}
	}
	client := redis.NewRing(&redis.RingOptions{
		Addrs:        servers,
		Password:     pass,
		IdleTimeout:  timeoutIdle,
		ReadTimeout:  timeoutRead,
		WriteTimeout: timeoutWrite,
		DialTimeout:  timeoutConnect,
		PoolSize:     200,
	})

	if _, err := client.Ping().Result(); err != nil {
		std.LogError(std.LogFields{"args": []interface{}{config}, "err": err}, "Failed to connect Redis.")
		client = nil
	}

	return client
}
