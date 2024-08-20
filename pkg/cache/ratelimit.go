package cache

import (
	"net/http"
	"new-token/pkg/router"
	"new-token/pkg/server"
	"time"
)

var (
	_rateLimitPerSecond = server.Config.GetInt64("RATE_LIMIT_PER_SECOND")
)

func RateLimitByIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Get user ip address
		IPAddress := r.Header.Get("X-Real-Ip")
		if IPAddress == "" {
			IPAddress = r.Header.Get("X-Forwarded-For")
		}
		if IPAddress == "" {
			IPAddress = r.RemoteAddr
		}

		//Check ip rate limit
		_, _, isUnderLimit := RedisCache.Limiter.Allow(IPAddress, _rateLimitPerSecond, time.Second)
		if !isUnderLimit {
			router.ResponseTooManyRequests(w)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
