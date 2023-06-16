package misc

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"rateMyRentalBackend/config"
	"rateMyRentalBackend/http/controllers"
	"rateMyRentalBackend/http/middlewares"
)

func MiscellaneousRouter(env *config.Env, db *gorm.DB, group *gin.RouterGroup) {
	misc := controllers.MiscController{
		Env: env,
		DB:  db,
	}
	group.POST("/media/upload", middlewares.CheckFileType(), misc.UploadFile)
}
