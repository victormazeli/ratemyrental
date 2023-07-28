package v1

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"rateMyRentalBackend/config"
	"rateMyRentalBackend/http/controllers"
	middlewares2 "rateMyRentalBackend/http/middlewares"
	"rateMyRentalBackend/http/services"
)

func PropertyRouter(env *config.Env, db *gorm.DB, group *gin.RouterGroup) {
	ps := services.PropertyService{
		Db:  db,
		Env: env,
	}
	p := controllers.PropertyController{
		Db:              db,
		Env:             env,
		PropertyService: ps,
	}
	group.GET("/property/get/:id", middlewares2.Auth(env.JwtKey), p.GetProperty)
	group.GET("/property/all", middlewares2.Auth(env.JwtKey), p.GetAllProperties)
	group.POST("/property/add", middlewares2.Auth(env.JwtKey), middlewares2.RateLimiter("/property/add", 5, 60, env), p.AddNewProperty)
	//group.POST("/property/upload/image", middlewares2.Auth(env.JwtKey), p.UploadImageProperty)
	group.POST("/property/add_favorite", middlewares2.Auth(env.JwtKey), p.AddOrRemoveFavoriteProperty)
	group.POST("/property/toggle_publish", middlewares2.Auth(env.JwtKey), p.TogglePropertyVisibility)
	group.GET("/property/get_favorite", middlewares2.Auth(env.JwtKey), p.GetUserFavoriteProperties)
	group.PUT("/property/edit/:id", middlewares2.Auth(env.JwtKey), middlewares2.RateLimiter("/property/edit", 5, 60, env), p.UpdatePropertyDetail)
	group.GET("/property/types", p.GetPropertyTypes)
	group.GET("/property/detached_types", p.GetPropertyDetachedTypes)
	group.PUT("/property/update_image", middlewares2.Auth(env.JwtKey), p.UpdateSingleImageProperty)
	group.POST("/property/add_rating", middlewares2.Auth(env.JwtKey), middlewares2.RateLimiter("/property/add_rating", 5, 60, env), p.RateProperty)
	group.GET("/property/user_properties", middlewares2.Auth(env.JwtKey), p.GetUserUploadedProperties)
	group.GET("/property/recommendations", middlewares2.Auth(env.JwtKey), p.PropertyRecommendations)

}
