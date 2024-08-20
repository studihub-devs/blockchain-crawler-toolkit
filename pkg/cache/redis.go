package cache

import (
	"encoding/json"
	"new-token/pkg/log"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-redis/redis_rate"
)

const (
	KEY_NOT_EXIST = redis.Nil
)

// Redis Configuration Struct
type redisConfig struct {
	Host     string
	Port     string
	Password string
	Name     int
	Enable   bool
}

var redisCfg = &redisConfig{}

type RedisCacheStore struct {
	RWMutex sync.RWMutex
	Limiter *redis_rate.Limiter
	Store   *redis.Client
}

// Redis Cache Variable
var RedisCache *RedisCacheStore

// Redis Connect Function
func redisConnect() *RedisCacheStore {
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
	return &RedisCacheStore{
		Limiter: redis_rate.NewLimiter(client),
		Store:   client,
	}
}

// Get method to check redis server connection
func (redisCache *RedisCacheStore) Ping() (string, error) {
	return redisCache.Store.Ping().Result()
}

// SetByKey method to set cache by given key with time to live
func (redisCache *RedisCacheStore) SetByKey(key string, value any, timeToLive time.Duration) error {
	if redisCfg.Enable {
		redisCache.RWMutex.Lock()
		defer redisCache.RWMutex.Unlock()
		byteValue, err := json.Marshal(value)
		if err != nil {
			return err
		}

		return redisCache.Store.Set(key, byteValue, timeToLive).Err()
	}
	return nil
}

// Get method to retrieve the value of a key. If not present, returns false.
func (redisCache *RedisCacheStore) Get(key string) ([]byte, bool, error) {
	if redisCfg.Enable {
		byteValue, err := redisCache.Store.Get(key).Bytes()
		if err == KEY_NOT_EXIST {
			return nil, false, nil
		}
		if err != nil {
			return nil, false, err
		}
		return byteValue, true, nil
	}

	return nil, false, nil
}

// InvalidateByKey  method to delete a key from cahce.
func (redisCache *RedisCacheStore) InvalidateByKey(key string) error {
	if redisCfg.Enable {
		redisCache.RWMutex.Lock()
		defer redisCache.RWMutex.Unlock()

		return redisCache.Store.Del(key).Err()
	}
	return nil
}

// SetByTags method to set cache by given tags with time to live
func (redisCache *RedisCacheStore) SetByTags(key string, value any, timeToLive time.Duration, tags []string) error {
	if redisCfg.Enable {
		redisCache.RWMutex.Lock()
		defer redisCache.RWMutex.Unlock()

		for _, tag := range tags {
			tagSet := NewTagSet()
			tagSetDataByte, found, err := redisCache.Get(tag)
			if err != nil {
				return err
			}

			if found {
				err = json.Unmarshal(tagSetDataByte, &tagSet.Data)
				if err != nil {
					return err
				}
			}
			tagSet.Add(key)
			tagSetDataByte, err = json.Marshal(tagSet.Data)
			if err != nil {
				return err
			}

			err = redisCache.Store.Set(tag, tagSetDataByte, timeToLive).Err()
			if err != nil {
				return err
			}
		}
		valueByte, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return redisCache.Store.Set(key, valueByte, timeToLive).Err()
	}
	return nil
}

// InvalidateByTags method to invalidate cache with given tags.
func (redisCache *RedisCacheStore) InvalidateByTags(tags []string) error {
	if redisCfg.Enable {
		redisCache.RWMutex.Lock()
		defer redisCache.RWMutex.Unlock()

		keys := make([]string, 0)
		for _, tag := range tags {
			tagSet := NewTagSet()
			tagSetByte, found, err := redisCache.Get(tag)
			if err != nil {
				return err
			}
			if found {
				err = json.Unmarshal(tagSetByte, &tagSet.Data)
				if err != nil {
					return err
				}
			}
			keys = append(keys, tagSet.Members()...)
			keys = append(keys, tag)
		}

		for _, k := range keys {
			err := redisCache.Store.Del(k).Err()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Close method clear and then close the cache Store.
func (redisCache *RedisCacheStore) Close() {
	if redisCfg.Enable {
		redisCache.Store.Close()
	}
}
