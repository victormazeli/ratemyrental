package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"rateMyRentalBackend/config"
	"rateMyRentalBackend/http/response"
	"strings"
	"time"
)

func RateLimiter(path string, ratelimit int64, ttl int, env *config.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		var key string
		clientIp := c.ClientIP()
		bearerToken := c.Request.Header.Get("Authorization")
		if len(strings.Split(bearerToken, " ")) == 2 {
			t := strings.Split(bearerToken, " ")[1]
			token = t
		} else {
			token = "0"
		}
		if path == "login" || path == "register" || token == "0" {
			key = fmt.Sprintf("%s - %s", path, clientIp)
		} else {
			key = fmt.Sprintf("%s - %s - %s", path, token, clientIp)
		}
		// Create a redis client.
		option, err := redis.ParseURL(env.RedisUrl)
		if err != nil {
			log.Fatal(err)
			return
		}
		client := redis.NewClient(option)

		count, e := client.Incr(c, key).Result()
		if e != nil {
			log.Fatal(e)
		}

		if count == 1 {
			client.Expire(c, key, time.Duration(ttl)*time.Second)
		}

		if count > ratelimit {
			response.ErrorResponse(http.StatusTooManyRequests, "too many request", c)
			return
		}
		c.Next()
	}
}
