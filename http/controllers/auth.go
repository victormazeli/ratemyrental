package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
	"log"
	"net/http"
	"rateMyRentalBackend/common/utils"
	"rateMyRentalBackend/config"
	models2 "rateMyRentalBackend/database/models"
	"rateMyRentalBackend/http/request"
	response2 "rateMyRentalBackend/http/response"
	"time"
)

const (
	registration   = "registration"
	forgetPassword = "forgot password"
	usedOtp        = 1
	unusedOtp      = 2
	expiryTime     = 10
)

var httpClient = resty.New()

type AuthController struct {
	Db  *gorm.DB
	Env *config.Env
}

func (a AuthController) Login(c *gin.Context) {
	var loginInput request.LoginInput
	// deserialize JSON to struct
	if err := c.ShouldBindJSON(&loginInput); err == nil {
		var user models2.User

		query := &models2.User{
			Email: loginInput.Email,
		}
		// check user
		findUser := a.Db.Where(query).First(&user)

		// check db error
		if findUser.Error == nil {
			match, err := utils.ComparePasswordAndHash(loginInput.Password, user.Password)

			if err == nil {
				if match == false {
					response2.ErrorResponse(http.StatusBadRequest, "invalid credentials", c)
					return
				} else {
					if user.Status == 0 {
						response2.ErrorResponse(http.StatusBadRequest, "user account not verified", c)
						return
					} else {
						newToken := utils.GenerateToken(user.ID, a.Env.JwtKey)

						response := &response2.AuthResponse{
							Token: newToken,
						}
						response2.SuccessResponse(http.StatusOK, "Login successful", response, c)
						return
					}
				}
			} else {
				log.Println(err)
				response2.ErrorResponse(http.StatusUnauthorized, "invalid credentials", c)
				return
			}
		} else {
			if errors.Is(findUser.Error, gorm.ErrRecordNotFound) == true {
				response2.ErrorResponse(http.StatusUnauthorized, "invalid credentials", c)
				return
			} else {
				response2.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}
		}
	} else {
		response2.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}
}
func (a AuthController) Register(c *gin.Context) {
	var registerInput request.RegisterInput
	if err := c.ShouldBind(&registerInput); err == nil {
		var user models2.User

		query := &models2.User{
			Email: registerInput.Email,
		}
		findUser := a.Db.Where(query).First(&user)

		if findUser.Error == nil {
			response2.ErrorResponse(http.StatusBadRequest, "User with email already exist", c)
			return
		} else {
			if errors.Is(findUser.Error, gorm.ErrRecordNotFound) == true {
				encodedHash, err := utils.GenerateFromPassword(registerInput.Password)
				if err == nil {
					registerInput.Password = encodedHash
					newUserPayload := map[string]interface{}{
						"email":     registerInput.Email,
						"password":  registerInput.Password,
						"full_name": registerInput.FullName,
					}
					newUser := a.Db.Model(&models2.User{}).Create(newUserPayload)
					if newUser.Error == nil {
						otp := utils.GenerateOTP()

						otpPayload := &models2.Otp{
							Email:      registerInput.Email,
							Purpose:    registration,
							Otp:        otp,
							Status:     0,
							ExpiryDate: time.Now().Add(expiryTime * time.Minute),
						}
						addOtpRecord := a.Db.Create(otpPayload)

						if addOtpRecord.Error == nil {
							subject := "Registration"
							body := utils.GenerateOTPEmailTemplate(otp)
							res, _ := utils.SendEmail(a.Env, registerInput.Email, body, subject)
							log.Print("email sent successfully", res)

							response2.SuccessResponse(http.StatusOK, "A link to activate your account has been emailed to the address provided.", nil, c)
							return
						} else {
							log.Println("log otp err >>", addOtpRecord.Error)
							response2.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
							return
						}
					} else {
						log.Println("log newUser err >>", newUser.Error)
						response2.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
						return
					}
				} else {
					log.Println("log encoding err >>", err)
					response2.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
					return
				}

			} else {
				log.Println("here >>", findUser.Error)
				response2.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}
		}
	} else {
		response2.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}

}

func (a AuthController) ForgetPassword(c *gin.Context) {
	var forgotPasswordInput request.ForgotPasswordInput
	if err := c.ShouldBind(&forgotPasswordInput); err == nil {
		var user models2.User

		query := &models2.User{
			Email: forgotPasswordInput.Email,
		}
		findUser := a.Db.Where(query).First(&user)
		// check db error

		if findUser.Error == nil {
			otp := utils.GenerateOTP()

			otpPayload := &models2.Otp{
				Email:      user.Email,
				Purpose:    forgetPassword,
				Otp:        otp,
				Status:     0,
				ExpiryDate: time.Now().Add(expiryTime * time.Minute),
			}
			addOtpRecord := a.Db.Create(otpPayload)

			if addOtpRecord.Error == nil {
				subject := "Password Reset"
				body := utils.GenerateForgetPasswordEmailTemplate(otp)
				res, _ := utils.SendEmail(a.Env, user.Email, body, subject)

				log.Println("email sent successfully", res)

				response2.SuccessResponse(http.StatusOK, "Reset password initiated successfully", nil, c)

			} else {
				log.Println(addOtpRecord.Error)
				response2.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}

		} else {
			if errors.Is(findUser.Error, gorm.ErrRecordNotFound) == true {
				response2.ErrorResponse(http.StatusNotFound, "user not found", c)
				return

			} else {
				log.Println(findUser.Error)
				response2.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}

		}

	} else {
		response2.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}

}

func (a AuthController) ResetPassword(c *gin.Context) {
	var resetPasswordInput request.ResetPassword

	if err := c.ShouldBind(&resetPasswordInput); err == nil {
		var otpData models2.Otp

		query := &models2.Otp{
			Otp:     resetPasswordInput.Otp,
			Purpose: forgetPassword,
			Status:  0,
		}
		findUserOtp := a.Db.Where(query).First(&otpData)

		// check db errors

		if findUserOtp.Error == nil {
			if otpData.ExpiryDate.Unix() > time.Now().Unix() {
				encodedHash, err := utils.GenerateFromPassword(resetPasswordInput.NewPassword)
				if err == nil {
					newPassword := encodedHash

					updateUser := a.Db.Model(&models2.User{}).Where("email = ?", otpData.Email).Update("password", newPassword)
					if updateUser.Error == nil {
						updateOtpTable := a.Db.Model(&otpData).Update("status", usedOtp)

						if updateOtpTable.Error == nil {
							response2.SuccessResponse(http.StatusOK, "Password reset successful", nil, c)
							return
						} else {
							log.Println(updateOtpTable.Error)
							response2.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
							return
						}

					} else {
						log.Println(updateUser.Error)
						response2.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
						return
					}

				} else {
					log.Println(err)
					response2.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
					return
				}
			} else {
				response2.ErrorResponse(http.StatusBadRequest, "otp expired", c)
				return
			}

		} else {
			if errors.Is(findUserOtp.Error, gorm.ErrRecordNotFound) == true {
				response2.ErrorResponse(http.StatusBadRequest, "invalid otp", c)
				return
			}
			log.Println(findUserOtp.Error)
			response2.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
			return
		}

	} else {
		response2.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}
}

func (a AuthController) ValidateOtp(c *gin.Context) {
	var validateOtpInput request.ValidateOtp

	if err := c.ShouldBind(&validateOtpInput); err == nil {
		var otpData models2.Otp

		findUserOtp := a.Db.Where(map[string]interface{}{"otp": validateOtpInput.Otp, "status": 0}).First(&otpData)

		if findUserOtp.Error == nil {
			// determine if otp is for registration then activate user account
			if otpData.Purpose == registration {
				if otpData.ExpiryDate.Unix() > time.Now().Local().Unix() {
					updateOtpTable := a.Db.Model(&otpData).Update("status", usedOtp)

					if updateOtpTable.Error == nil {
						updateUser := a.Db.Model(&models2.User{}).Where("email = ?", otpData.Email).Update("status", 1)
						if updateUser.Error == nil {
							response2.SuccessResponse(http.StatusOK, "Otp validation successful", nil, c)
							return
						} else {
							response2.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
							return
						}
					} else {
						log.Println(updateOtpTable.Error)
						response2.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
						return
					}
				} else {
					response2.ErrorResponse(http.StatusBadRequest, "otp expired", c)
					return
				}
			} else {
				if otpData.ExpiryDate.Unix() > time.Now().Local().Unix() {
					updateOtpTable := a.Db.Model(&otpData).Update("status", usedOtp)

					if updateOtpTable.Error == nil {
						response2.SuccessResponse(http.StatusOK, "Otp validation successful", nil, c)
						return
					} else {
						log.Println(updateOtpTable.Error)
						response2.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
						return
					}
				} else {
					response2.ErrorResponse(http.StatusBadRequest, "otp expired", c)
					return
				}
			}
		} else {
			if errors.Is(findUserOtp.Error, gorm.ErrRecordNotFound) == true {
				response2.ErrorResponse(http.StatusBadRequest, "invalid otp", c)
				return
			} else {
				response2.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}
		}

	} else {
		response2.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}

}

func (a AuthController) ResendOtp(c *gin.Context) {
	var resendOtpInput request.ResendOtpInput

	if err := c.ShouldBind(&resendOtpInput); err == nil {
		var purpose string
		var resetPasswordTruthiness = true

		if resendOtpInput.IsPasswordReset == &resetPasswordTruthiness {
			purpose = forgetPassword
		} else {
			purpose = registration
		}

		var user models2.User

		query := &models2.User{
			Email: resendOtpInput.Email,
		}
		findUser := a.Db.Where(query).First(&user)

		if findUser.Error == nil {
			otp := utils.GenerateOTP()

			otpPayload := &models2.Otp{
				Email:      user.Email,
				Purpose:    purpose,
				Otp:        otp,
				Status:     0,
				ExpiryDate: time.Now().Add(expiryTime * time.Minute),
			}
			addOtpRecord := a.Db.Create(otpPayload)

			if addOtpRecord.Error == nil {
				if resendOtpInput.IsPasswordReset == &resetPasswordTruthiness {
					subject := "Password Reset"
					body := utils.GenerateForgetPasswordEmailTemplate(otp)
					res, _ := utils.SendEmail(a.Env, user.Email, body, subject)
					log.Println("email sent", res)
					response2.SuccessResponse(http.StatusOK, "Otp sent successfully", nil, c)
				} else {
					subject := "Registration"
					body := utils.GenerateOTPEmailTemplate(otp)
					res, _ := utils.SendEmail(a.Env, user.Email, body, subject)
					log.Println("email sent", res)
					response2.SuccessResponse(http.StatusOK, "Otp sent successfully", nil, c)
				}

			} else {
				log.Println(addOtpRecord.Error)
				response2.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}

		} else {
			if errors.Is(findUser.Error, gorm.ErrRecordNotFound) == true {
				response2.ErrorResponse(http.StatusBadRequest, "user with email not found", c)
				return
			} else {
				response2.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}
		}
	} else {
		response2.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}

}

func (a AuthController) GoogleLogin(c *gin.Context) {
	var payload request.GoogleAuthModel
	var reqBody request.GoogleAuth

	if err := c.ShouldBind(&reqBody); err != nil {
		response2.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}

	url := "https://www.googleapis.com/oauth2/v3/userinfo"

	res, err := httpClient.R().SetHeader("Accept", "application/json").SetAuthToken(reqBody.AccessToken).Get(url)

	if err != nil {
		response2.ErrorResponse(http.StatusInternalServerError, err.Error(), c)
		return
	}
	if res.StatusCode() == http.StatusOK {
		err := json.Unmarshal(res.Body(), &payload)
		if err != nil {
			response2.ErrorResponse(http.StatusInternalServerError, err.Error(), c)
			return
		}

		user := &models2.User{}

		query := &models2.User{
			Email: *payload.Email,
		}
		findUser := a.Db.Where(query).First(&user)

		if findUser.Error != nil {
			// create user
			newUserPayload := map[string]interface{}{
				"email":     *payload.Email,
				"password":  "",
				"full_name": *payload.Name,
			}
			if insertError := a.Db.Model(&models2.User{}).Create(newUserPayload).Error; insertError != nil {
				response2.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}

			token := utils.GenerateToken(user.ID, a.Env.JwtKey)

			//var appTokens []string
			//
			//appTokens = append(appTokens, token)
			//
			//redis.CacheService{}.SetAppToken(ctx, "appToken", appTokens)

			response := &response2.AuthResponse{
				Token: token,
			}
			response2.SuccessResponse(http.StatusOK, "Login successful", response, c)
			return
		}

		token := utils.GenerateToken(user.ID, a.Env.JwtKey)

		//var appTokens []string
		//
		//appTokens = append(appTokens, token)
		//
		//redis.CacheService{}.SetAppToken(ctx, "appToken", appTokens)

		response := &response2.AuthResponse{
			Token: token,
		}
		response2.SuccessResponse(http.StatusOK, "Login successful", response, c)
		return

	} else {
		response2.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		return
	}

}
