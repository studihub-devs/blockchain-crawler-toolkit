package ms

import (
	"new-token/pkg/log"
	"sync"

	"github.com/go-redis/redis"
	"github.com/go-redis/redis_rate"
)

type redisConfig struct {
	Host     string
	Port     string
	Password string
	Name     int
	Enable   bool
}

var redisCfg = &redisConfig{}

type RedisMS struct {
	RWMutex sync.RWMutex
	Limiter *redis_rate.Limiter
	Store   *redis.Client
}

// Redis MS Variable
var Redis *RedisMS

// Redis Connect Function
func redisConnect() *RedisMS {
	// Initialize Connection
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Host + ":" + redisCfg.Port,
		Password: redisCfg.Password,
		DB:       redisCfg.Name,
	})

	// Test Connection
	_, err := client.Ping().Result()
	if err != nil {
		log.Println(log.LogLevelFatal, "redis-connect", err.Error())
	} else {
		log.Println(log.LogLevelInfo, "redis-connect", "Connect redis: Successfully connected")
	}

	// Return Connection
	return &RedisMS{
		Limiter: redis_rate.NewLimiter(client),
		Store:   client,
	}
}

func (redisMS *RedisMS) Ping() (string, error) {
	return redisMS.Store.Ping().Result()
}

// Close method clear and then close the cache Store.
func (redisMS *RedisMS) Close() {
	if redisCfg.Enable {
		redisMS.Store.Close()
	}
}
