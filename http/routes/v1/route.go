package v1

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"rateMyRentalBackend/config"
)

func Setup(env *config.Env, db *gorm.DB, routerV1 *gin.RouterGroup) {
	//v1 := routerV1.Group("")
	// All auth route
	AuthRouter(env, db, routerV1)
	UserRouter(env, db, routerV1)
	PropertyRouter(env, db, routerV1)

}
