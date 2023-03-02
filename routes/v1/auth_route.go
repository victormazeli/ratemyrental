package v1

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"rateMyRentalBackend/config"
	"rateMyRentalBackend/controllers"
	"rateMyRentalBackend/middlewares"
)

func AuthRouter(env *config.Env, db *gorm.DB, group *gin.RouterGroup) {
	ac := controllers.AuthController{
		Db:  db,
		Env: env,
	}
	group.POST("/auth/login", middlewares.RateLimiter("login", 5, 60, env), ac.Login)
	group.POST("/auth/register", middlewares.RateLimiter("register", 5, 60, env), ac.Register)
	group.POST("/auth/reset_password", middlewares.RateLimiter("/auth/reset_password", 5, 60, env), ac.ResetPassword)
	group.POST("/auth/forgot_password", middlewares.RateLimiter("/auth/forgot_password", 5, 60, env), ac.ForgetPassword)
	group.POST("/auth/validate_otp", middlewares.RateLimiter("/auth/validate_otp", 5, 60, env), ac.ValidateOtp)
	group.POST("/auth/resend_otp", middlewares.RateLimiter("/auth/resend_otp", 5, 60, env), ac.ResendOtp)
}
