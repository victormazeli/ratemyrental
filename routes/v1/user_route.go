package v1

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"rateMyRentalBackend/config"
	"rateMyRentalBackend/controllers"
	"rateMyRentalBackend/middlewares"
)

func UserRouter(env *config.Env, db *gorm.DB, group *gin.RouterGroup) {
	uc := controllers.UserController{
		Db:  db,
		Env: env,
	}
	group.GET("/user/me", middlewares.Auth(env.JwtKey), uc.GetCurrentUser)
	group.PUT("/user/update", middlewares.Auth(env.JwtKey), middlewares.RateLimiter("/user/update", 5, 60, env), uc.UpdateUserInfo)
}
