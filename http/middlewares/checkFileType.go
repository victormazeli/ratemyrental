package middlewares

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"rateMyRentalBackend/http/response"
)

func CheckFileType() gin.HandlerFunc {
	return func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			response.ErrorResponse(http.StatusBadRequest, err.Error(), c)
			c.Abort()
			return
		}

		files := form.File["files"]

		allowedFileTypes := map[string]bool{
			"image/jpeg":      true,
			"image/jpg":       true,
			"image/png":       true,
			"application/pdf": true,
		}

		for _, file := range files {
			f, err := file.Open()
			if err != nil {
				response.ErrorResponse(http.StatusInternalServerError, err.Error(), c)
				c.Abort()
				return
			}

			fileContent, readErr := io.ReadAll(f)
			if err != nil {
				f.Close()
				response.ErrorResponse(http.StatusInternalServerError, readErr.Error(), c)
				c.Abort()
				return
			}

			f.Close()

			fileType := http.DetectContentType(fileContent)
			if !allowedFileTypes[fileType] {
				response.ErrorResponse(http.StatusBadRequest, "Invalid file type. Only JPEG, JPG, PNG, and PDF files are allowed.", c)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
