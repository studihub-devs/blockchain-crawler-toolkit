package cache

import (
	"new-token/pkg/server"
	"strings"
)

// Initialize Function in Cache Package
func init() {
	// Remote Cache Configuration Value
	switch strings.ToLower(server.Config.GetString("REMOTE_CACHE_DRIVER")) {
	case "redis":
		redisCfg.Host = server.Config.GetString("REMOTE_CACHE_HOST")
		redisCfg.Port = server.Config.GetString("REMOTE_CACHE_PORT")
		redisCfg.Password = server.Config.GetString("REMOTE_CACHE_PASSWORD")
		redisCfg.Name = server.Config.GetInt("REMOTE_CACHE_NAME")
		redisCfg.Enable = server.Config.GetBool("REMOTE_CACHE_ENABLE")

		if len(redisCfg.Host) != 0 && len(redisCfg.Port) != 0 {
			RedisCache = redisConnect()
		}
	}
}
