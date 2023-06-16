package v1

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"rateMyRentalBackend/config"
	"rateMyRentalBackend/http/controllers"
	middlewares2 "rateMyRentalBackend/http/middlewares"
)

func UserRouter(env *config.Env, db *gorm.DB, group *gin.RouterGroup) {
	uc := controllers.UserController{
		Db:  db,
		Env: env,
	}
	group.GET("/user/me", middlewares2.Auth(env.JwtKey), uc.GetCurrentUser)
	group.POST("/user/switch", middlewares2.Auth(env.JwtKey), uc.SwitchProfile)
	group.GET("/user", uc.GetAllUsers)
	group.GET("/user/:id", uc.GetUserByID)
	group.PUT("/user/update", middlewares2.Auth(env.JwtKey), middlewares2.RateLimiter("/user/update", 5, 60, env), uc.UpdateUserInfo)
}
