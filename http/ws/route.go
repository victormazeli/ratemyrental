package ws

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"rateMyRentalBackend/config"
)

func WebsocketRouter(env *config.Env, db *gorm.DB, group *gin.RouterGroup) {
	w := ChatController{
		Env: env,
		DB:  db,
	}
	group.GET("/ws", w.UpgradeWebsocket)
	group.POST("/channel/join", w.JoinChannel)
	group.POST("/chat/messages", w.SendMessage)
	group.GET("/chat/messages", w.GetMessages)
}
