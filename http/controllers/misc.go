package controllers

import (
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"rateMyRentalBackend/config"
	"rateMyRentalBackend/http/response"
)

type MiscController struct {
	Env *config.Env
	DB  *gorm.DB
}

func (m MiscController) UploadFile(c *gin.Context) {
	cld, err := cloudinary.NewFromParams(m.Env.CloudinaryName, m.Env.CloudinaryApiKey, m.Env.CloudinarySecret)
	if err != nil {
		log.Fatalf("failed to initialize cloudinary: %v", err)
	}

	formData, err := c.MultipartForm()
	if err != nil {
		response.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}

	files := formData.File["files"]
	if files == nil {
		response.ErrorResponse(http.StatusBadRequest, "Please select the input file", c)
		return
	}

	var urls []string

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			response.ErrorResponse(http.StatusInternalServerError, err.Error(), c)
			return
		}
		defer src.Close()

		tmpFileName := fmt.Sprintf("tmp_%s", file.Filename)
		dst, err := os.Create(tmpFileName)
		if err != nil {
			response.ErrorResponse(http.StatusInternalServerError, err.Error(), c)
			return
		}
		defer func() {
			dst.Close()
			os.Remove(tmpFileName)
		}()

		_, err = io.Copy(dst, src)
		if err != nil {
			response.ErrorResponse(http.StatusInternalServerError, err.Error(), c)
			return
		}

		resp, errUpload := cld.Upload.Upload(c, tmpFileName, uploader.UploadParams{})
		if err != nil {
			response.ErrorResponse(http.StatusUnprocessableEntity, errUpload.Error(), c)
			return
		}

		urls = append(urls, resp.SecureURL)
	}

	response.SuccessResponse(http.StatusOK, "upload successful", urls, c)
}
