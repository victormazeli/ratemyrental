package v1

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"rateMyRentalBackend/config"
	"rateMyRentalBackend/controllers"
	"rateMyRentalBackend/middlewares"
)

func PropertyRouter(env *config.Env, db *gorm.DB, group *gin.RouterGroup) {
	p := controllers.PropertyController{
		Db:  db,
		Env: env,
	}
	group.GET("/property/get/:id", middlewares.Auth(env.JwtKey), p.GetProperty)
	group.GET("/property/all", middlewares.Auth(env.JwtKey), p.GetAllProperties)
	group.POST("/property/add", middlewares.Auth(env.JwtKey), middlewares.RateLimiter("/property/add", 5, 60, env), p.AddNewProperty)
	group.POST("/property/upload/image", middlewares.Auth(env.JwtKey), p.UploadImageProperty)
	group.POST("/property/add_favorite", middlewares.Auth(env.JwtKey), p.AddOrRemoveFavoriteProperty)
	group.POST("/property/toggle_publish", middlewares.Auth(env.JwtKey), p.TogglePropertyVisibility)
	group.GET("/property/get_favorite", middlewares.Auth(env.JwtKey), p.GetUserFavoriteProperties)
	group.PUT("/property/edit/:id", middlewares.Auth(env.JwtKey), middlewares.RateLimiter("/property/edit", 5, 60, env), p.UpdatePropertyDetail)
	group.GET("/property/types", p.GetPropertyTypes)
	group.GET("/property/detached_types", p.GetPropertyDetachedTypes)
	group.PUT("/property/update_image", middlewares.Auth(env.JwtKey), p.UpdateSingleImageProperty)
	group.POST("/property/add_rating", middlewares.Auth(env.JwtKey), middlewares.RateLimiter("/property/add_rating", 5, 60, env), p.RateProperty)
	group.GET("/property/user_properties", middlewares.Auth(env.JwtKey), p.GetUserUploadedProperties)

}
