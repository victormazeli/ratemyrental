package ws

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"log"
	"net/http"
	"rateMyRentalBackend/config"
	models2 "rateMyRentalBackend/database/models"
	"rateMyRentalBackend/http/response"
	"strconv"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	hub = HubInstance
)

var chatMessages []struct {
	RecipientName string
	Content       string
	CreatedAt     time.Time
	ChannelID     uint
}

type ChatController struct {
	Env *config.Env
	DB  *gorm.DB
}

func (ch ChatController) JoinChannel(c *gin.Context) {
	var channel models2.Channel

	if err := c.ShouldBindJSON(&channel); err != nil {
		response.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}

	// Get the channel ID based on the sender and recipient IDs
	channelKey := fmt.Sprintf("%v:%v", channel.SenderID, channel.RecipientID)
	channelID, ok := hub.channels[channelKey]
	if !ok {
		// Check if the channel exists in the database
		err := ch.DB.Where("(sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)", channel.SenderID, channel.RecipientID, channel.RecipientID, channel.SenderID).First(&channel).Error
		if err != nil {
			if err := ch.DB.Create(&channel).Error; err != nil {
				response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}
		}

		// Store the channel ID in the hub and return it to the client
		hub.channels[channelKey] = channel.ID
		channelID = channel.ID
	}

	resData := make(map[string]interface{})
	resData["channel_id"] = channelID
	response.SuccessResponse(http.StatusOK, "channel joined successfully", resData, c)

}

func (ch ChatController) SendMessage(c *gin.Context) {
	var message models2.ChatMessage
	if err := c.ShouldBindJSON(&message); err != nil {
		response.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}

	// Set the channel ID in the WebSocketMessage struct
	wsMessage := WebSocketMessage{
		Sender:    message.Sender,
		Recipient: message.Recipient,
		Content:   message.Content,
		ChannelID: message.ChannelID,
	}

	if err := ch.DB.Create(&message).Error; err != nil {
		response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		return
	}

	hub.broadcast <- wsMessage

	response.SuccessResponse(http.StatusOK, "Message sent successfully", message.ChannelID, c)
}

func (ch ChatController) GetMessages(c *gin.Context) {
	channelID := c.Query("channelID")

	if channelID == "" {
		response.ErrorResponse(http.StatusBadRequest, "channelID is missing in query", c)
		return
	}

	var messages []models2.ChatMessage
	if err := ch.DB.Where("channel_id = ?", channelID).Find(&messages).Error; err != nil {
		response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		return
	}
	response.SuccessResponse(http.StatusOK, "Messages retrieved successfully", messages, c)
}

func (ch ChatController) GetUsersChat(c *gin.Context) {
	userId, _ := c.Get("user")
	log.Print(userId)
	if err := ch.DB.Table("chat_messages").
		Select("users.full_name AS recipient_name, chat_messages.content, chat_messages.created_at, chat_messages.channel_id").
		Joins("LEFT JOIN users ON chat_messages.recipient = users.id").
		Where("(chat_messages.sender = ? OR chat_messages.recipient = ?)", userId, userId).
		Where("chat_messages.created_at = (SELECT MAX(created_at) FROM chat_messages WHERE chat_messages.channel_id = chat_messages.channel_id)").
		Find(&chatMessages).Error; err != nil {
		response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		return
	}

	response.SuccessResponse(http.StatusOK, "Messages retrieved successfully", chatMessages, c)
}

//func (ch ChatController) GetMessages(c *gin.Context) {
//	sender := c.Query("sender")
//	recipient := c.Query("recipient")
//
//	// Get the channel ID based on the sender and recipient IDs
//	channelKey := fmt.Sprintf("%s:%s", sender, recipient)
//	channelID, ok := hub.channels[channelKey]
//	if !ok {
//		response.ErrorResponse(http.StatusBadRequest, "Channel not found", c)
//		return
//	}
//
//	var messages []models.ChatMessage
//	if err := ch.DB.Where("channel_id = ? AND ((sender = ? AND recipient = ?) OR (sender = ? AND recipient = ?))", channelID, sender, recipient, recipient, sender).Find(&messages).Error; err != nil {
//		response.ErrorResponse(http.StatusInternalServerError, err.Error(), c)
//		return
//	}
//	response.SuccessResponse(http.StatusOK, "Message retrieved successfully", messages, c)
//
//}

func (ch ChatController) UpgradeWebsocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade WebSocket connection:", err)
		return
	}

	// Get the channel ID from the connection's query parameters
	channelID, err := strconv.ParseUint(c.Query("channel_id"), 10, 32)
	if err != nil {
		log.Println("Failed to parse channel ID:", err)
		conn.Close()
		return
	}

	// Add the connection to the corresponding channel's client list
	if _, ok := hub.clients[uint(channelID)]; !ok {
		hub.clients[uint(channelID)] = make(map[*websocket.Conn]bool)
	}
	hub.clients[uint(channelID)][conn] = true

	go func() {
		defer func() {
			// Remove the connection from the corresponding channel's client list
			for channelID, clients := range hub.clients {
				if _, ok := clients[conn]; ok {
					delete(clients, conn)
					conn.Close()

					// If no clients are left in the channel, delete the channel
					if len(clients) == 0 {
						delete(hub.clients, channelID)
						delete(hub.channels, strconv.FormatUint(uint64(channelID), 10))
					}
				}
			}
		}()

		for {
			var message WebSocketMessage
			err := conn.ReadJSON(&message)
			if err != nil {
				log.Println("Error reading WebSocket message:", err)
				break
			}

			hub.broadcast <- message
		}
	}()

}
