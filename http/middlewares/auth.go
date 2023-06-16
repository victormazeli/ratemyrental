package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rateMyRentalBackend/common/utils"
	"rateMyRentalBackend/http/response"
	"strings"
)

func Auth(jwtkey string) gin.HandlerFunc {
	return func(context *gin.Context) {
		bearerToken := context.Request.Header.Get("Authorization")
		if len(strings.Split(bearerToken, " ")) == 2 {
			token := strings.Split(bearerToken, " ")[1]
			sub, err := utils.ValidateToken(token, jwtkey)
			if err != nil {
				response.ErrorResponse(http.StatusUnauthorized, err.Error(), context)
				return
			} else {
				context.Set("user", sub)
				context.Next()
			}
		} else {
			response.ErrorResponse(http.StatusUnauthorized, "invalid token", context)
			return
		}
	}
}
