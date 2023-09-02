package ws

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"rateMyRentalBackend/config"
	middlewares2 "rateMyRentalBackend/http/middlewares"
)

func WebsocketRouter(env *config.Env, db *gorm.DB, group *gin.RouterGroup) {
	w := ChatController{
		Env: env,
		DB:  db,
	}
	group.GET("/ws", w.UpgradeWebsocket)
	group.POST("/channel/join", w.JoinChannel)
	group.GET("/chat/user/messages", middlewares2.Auth(env.JwtKey), w.GetUsersChat)
	group.POST("/chat/messages", middlewares2.Auth(env.JwtKey), w.SendMessage)
	group.GET("/chat/messages", middlewares2.Auth(env.JwtKey), w.GetMessages)
}
