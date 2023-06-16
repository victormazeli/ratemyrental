package response

import (
	"github.com/gin-gonic/gin"
)

func SuccessResponse(code int, message string, data interface{}, c *gin.Context) {
	c.JSON(code, gin.H{"status": code, "success": true, "message": message, "data": data})
}

func ErrorResponse(code int, message string, c *gin.Context) {
	c.JSON(code, gin.H{"status": code, "error": true, "message": message})
	c.Abort()
}
