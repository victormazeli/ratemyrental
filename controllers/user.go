package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"rateMyRentalBackend/common"
	"rateMyRentalBackend/config"
	"rateMyRentalBackend/dtos/input"
	"rateMyRentalBackend/models"
)

type UserController struct {
	Db  *gorm.DB
	Env *config.Env
}

func (u UserController) GetCurrentUser(c *gin.Context) {
	userId, _ := c.Get("user")

	var user models.User

	findUser := u.Db.Model(&models.User{}).Omit("password").Where("id = ?", userId).First(&user)

	if findUser.Error == nil {
		common.SuccessResponse(http.StatusOK, "User fetched successfully", user, c)
		return
	} else {
		if errors.Is(findUser.Error, gorm.ErrRecordNotFound) == true {
			common.ErrorResponse(http.StatusNotFound, "user not found", c)
			return
		} else {
			log.Println(findUser.Error)
			common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
			return
		}
	}

}

func (u UserController) UpdateUserInfo(c *gin.Context) {
	userId, _ := c.Get("user")
	var user models.User
	var userInfoInput input.UserInput

	if err := c.ShouldBindJSON(&userInfoInput); err == nil {
		result := u.Db.Model(&models.User{}).Where("id = ?", userId).Updates(map[string]interface{}{
			"full_name":   userInfoInput.FullName,
			"avatar":      userInfoInput.Avatar,
			"address":     userInfoInput.Address,
			"postal_code": userInfoInput.PostalCode,
			"location":    userInfoInput.Location,
		})

		if result.Error == nil {
			userResult := u.Db.Model(&models.User{}).Where("id = ?", userId).First(&user)
			if userResult.Error == nil {
				common.SuccessResponse(http.StatusOK, "user info updated successfully", user, c)
				return

			} else {
				if errors.Is(userResult.Error, gorm.ErrRecordNotFound) == true {
					log.Println(userResult.Error)
					common.ErrorResponse(http.StatusNotFound, "user not found", c)
					return

				} else {
					log.Println(userResult.Error)
					common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
					return
				}
			}
		} else {
			log.Println(result.Error)
			common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
			return
		}

	} else {
		common.ErrorResponse(http.StatusBadRequest, err.Error(), c)
	}

}

//func (u UserController) DeactivateAccount(c *gin.Context)  {
//
//
//}
