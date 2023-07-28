package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"rateMyRentalBackend/config"
	"rateMyRentalBackend/database/models"
	"rateMyRentalBackend/http/request"
	"rateMyRentalBackend/http/response"
	"strconv"
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
		response.SuccessResponse(http.StatusOK, "User fetched successfully", user, c)
		return
	} else {
		if errors.Is(findUser.Error, gorm.ErrRecordNotFound) == true {
			response.ErrorResponse(http.StatusNotFound, "user not found", c)
			return
		} else {
			log.Println(findUser.Error)
			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
			return
		}
	}

}

func (u UserController) GetUserByID(c *gin.Context) {

	userId := c.Param("id")
	id, err := strconv.ParseInt(userId, 0, 8)

	if err != nil {
		response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		return
	}

	var user models.User

	findUser := u.Db.Model(&models.User{}).Omit("password").Where("id = ?", id).First(&user)

	if findUser.Error == nil {
		response.SuccessResponse(http.StatusOK, "User fetched successfully", user, c)
		return
	} else {
		if errors.Is(findUser.Error, gorm.ErrRecordNotFound) == true {
			response.ErrorResponse(http.StatusNotFound, "user not found", c)
			return
		} else {
			log.Println(findUser.Error)
			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
			return
		}
	}

}

func (u UserController) UpdateUserInfo(c *gin.Context) {
	userId, _ := c.Get("user")
	var user models.User
	var userInfoInput request.UserInput

	if err := c.ShouldBindJSON(&userInfoInput); err == nil {
		result := u.Db.Model(&models.User{}).Where("id = ?", userId).Updates(map[string]interface{}{
			"full_name": userInfoInput.FullName,
			"avatar":    userInfoInput.Avatar,
			//"address":     userInfoInput.Address,
			"postal_code": userInfoInput.PostalCode,
			"latitude":    userInfoInput.Latitude,
			"longitude":   userInfoInput.Longitude,
			"city":        userInfoInput.City,
			"country":     userInfoInput.Country,
		})

		if result.Error == nil {
			userResult := u.Db.Model(&models.User{}).Where("id = ?", userId).First(&user)
			if userResult.Error == nil {
				response.SuccessResponse(http.StatusOK, "user info updated successfully", user, c)
				return

			} else {
				if errors.Is(userResult.Error, gorm.ErrRecordNotFound) == true {
					log.Println(userResult.Error)
					response.ErrorResponse(http.StatusNotFound, "user not found", c)
					return

				} else {
					log.Println(userResult.Error)
					response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
					return
				}
			}
		} else {
			log.Println(result.Error)
			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
			return
		}

	} else {
		response.ErrorResponse(http.StatusBadRequest, err.Error(), c)
	}

}

func (u UserController) SwitchProfile(c *gin.Context) {
	userId, _ := c.Get("user")
	var userTypeSwitch request.UserSwitchType
	var user models.User
	if err := c.ShouldBind(&userTypeSwitch); err != nil {
		response.ErrorResponse(http.StatusBadRequest, err.Error(), c)

	}

	result := u.Db.Model(&models.User{}).Where("id = ?", userId).Updates(map[string]interface{}{
		"user_type": 2,
	})

	if result.Error == nil {
		response.SuccessResponse(http.StatusOK, "user info updated successfully", user, c)
		return

	} else {
		log.Println(result.Error)
		response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		return
	}
	//userResult := u.Db.Model(&models.User{}).Where("id = ?", userId).First(&user)
	//if userResult.Error == nil {
	//	if user.UserType == 1 {
	//
	//		result := u.Db.Model(&models.User{}).Where("id = ?", userId).Updates(map[string]interface{}{
	//			"user_type": 2,
	//		})
	//
	//		if result.Error == nil {
	//			response.SuccessResponse(http.StatusOK, "user info updated successfully", user, c)
	//			return
	//
	//		} else {
	//			log.Println(result.Error)
	//			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
	//			return
	//		}
	//
	//	} else {
	//
	//		result := u.Db.Model(&models.User{}).Where("id = ?", userId).Updates(map[string]interface{}{
	//			"user_type": 1,
	//		})
	//
	//		if result.Error == nil {
	//			response.SuccessResponse(http.StatusOK, "user info updated successfully", user, c)
	//			return
	//
	//		} else {
	//			log.Println(result.Error)
	//			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
	//			return
	//		}
	//	}
	//
	//} else {
	//	if errors.Is(userResult.Error, gorm.ErrRecordNotFound) == true {
	//		log.Println(userResult.Error)
	//		response.ErrorResponse(http.StatusNotFound, "user not found", c)
	//		return
	//
	//	} else {
	//		log.Println(userResult.Error)
	//		response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
	//		return
	//	}
	//}

}

func (u UserController) GetAllUsers(c *gin.Context) {
	var users []models.User

	err := u.Db.Model(&models.User{}).Omit("password").Find(&users).Error

	if err != nil {
		response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
	}

	response.SuccessResponse(http.StatusOK, "users fetched successfully", users, c)

}
