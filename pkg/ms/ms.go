package ms

import (
	"new-token/pkg/server"
	"strings"
)

// Initialize Function in Cache Package
func init() {
	// Remote Cache Configuration Value
	switch strings.ToLower(server.Config.GetString("MS_DRIVER")) {
	case "redis":
		redisCfg.Host = server.Config.GetString("MS_HOST")
		redisCfg.Port = server.Config.GetString("MS_PORT")
		redisCfg.Password = server.Config.GetString("MS_PASSWORD")
		redisCfg.Name = server.Config.GetInt("MS_NAME")

		if len(redisCfg.Host) != 0 && len(redisCfg.Port) != 0 {
			Redis = redisConnect()
		}
	}
}
