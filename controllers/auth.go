package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"rateMyRentalBackend/common"
	"rateMyRentalBackend/common/utils"
	"rateMyRentalBackend/config"
	"rateMyRentalBackend/dtos"
	"rateMyRentalBackend/dtos/input"
	"rateMyRentalBackend/models"
	"time"
)

const registration = "registration"
const forgetPassword = "forgot password"
const usedOtp = 1
const unusedOtp = 2
const expiryTime = 10

type AuthController struct {
	Db  *gorm.DB
	Env *config.Env
}

func (a AuthController) Login(c *gin.Context) {
	var loginInput input.LoginInput
	// deserialize JSON to struct
	if err := c.ShouldBindJSON(&loginInput); err == nil {
		var user models.User

		query := &models.User{
			Email: loginInput.Email,
		}
		// check user
		findUser := a.Db.Where(query).First(&user)

		// check db error
		if findUser.Error == nil {
			match, err := utils.ComparePasswordAndHash(loginInput.Password, user.Password)

			if err == nil {
				if match == false {
					common.ErrorResponse(http.StatusBadRequest, "invalid credentials", c)
					return
				} else {
					if user.Status == 0 {
						common.ErrorResponse(http.StatusBadRequest, "user account not verified", c)
						return
					} else {
						newToken := utils.GenerateToken(user.ID, a.Env.JwtKey)

						response := &dtos.AuthResponse{
							User: dtos.UserDTO{
								Email:      user.Email,
								FullName:   user.FullName.String,
								Address:    user.Address.String,
								Status:     user.Status,
								Avatar:     user.Avatar.String,
								Location:   user.Location.String,
								PostalCode: user.PostalCode.String,
								Id:         user.ID,
								CreatedAt:  user.CreatedAt,
								UpdatedAt:  user.UpdatedAt,
								DeletedAt:  user.DeletedAt.Time,
							},
							Token: newToken,
						}
						common.SuccessResponse(http.StatusOK, "Login successful", response, c)
						return
					}
				}
			} else {
				log.Println(err)
				common.ErrorResponse(http.StatusUnauthorized, "invalid credentials", c)
				return
			}
		} else {
			if errors.Is(findUser.Error, gorm.ErrRecordNotFound) == true {
				common.ErrorResponse(http.StatusUnauthorized, "invalid credentials", c)
				return
			} else {
				common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}
		}
	} else {
		common.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}
}
func (a AuthController) Register(c *gin.Context) {
	var registerInput input.RegisterInput
	if err := c.ShouldBind(&registerInput); err == nil {
		var user models.User

		query := &models.User{
			Email: registerInput.Email,
		}
		findUser := a.Db.Where(query).First(&user)

		if findUser.Error == nil {
			common.ErrorResponse(http.StatusBadRequest, "User with email already exist", c)
			return
		} else {
			if errors.Is(findUser.Error, gorm.ErrRecordNotFound) == true {
				encodedHash, err := utils.GenerateFromPassword(registerInput.Password)
				if err == nil {
					registerInput.Password = encodedHash
					newUserPayload := map[string]interface{}{
						"email":       registerInput.Email,
						"password":    registerInput.Password,
						"full_name":   registerInput.FirstName + " " + registerInput.Lastname,
						"address":     registerInput.Address,
						"postal_code": registerInput.PostalCode,
						"location":    registerInput.Location,
					}
					newUser := a.Db.Model(&models.User{}).Create(newUserPayload)
					if newUser.Error == nil {
						otp := utils.GenerateOTP()

						otpPayload := &models.Otp{
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
							utils.SendEmail(registerInput.Email, body, subject)

							common.SuccessResponse(http.StatusOK, "A link to activate your account has been emailed to the address provided.", nil, c)
							return
						} else {
							log.Println("log otp err >>", addOtpRecord.Error)
							common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
							return
						}
					} else {
						log.Println("log newUser err >>", newUser.Error)
						common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
						return
					}
				} else {
					log.Println("log encoding err >>", err)
					common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
					return
				}

			} else {
				log.Println("here >>", findUser.Error)
				common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}
		}
	} else {
		common.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}

}

func (a AuthController) ForgetPassword(c *gin.Context) {
	var forgotPasswordInput input.ForgotPasswordInput
	if err := c.ShouldBind(&forgotPasswordInput); err == nil {
		var user models.User

		query := &models.User{
			Email: forgotPasswordInput.Email,
		}
		findUser := a.Db.Where(query).First(&user)
		// check db error

		if findUser.Error == nil {
			otp := utils.GenerateOTP()

			otpPayload := &models.Otp{
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
				e := utils.SendEmail(user.Email, body, subject)

				if e == nil {
					common.SuccessResponse(http.StatusOK, "If the provided email address is in our database, we will send you an email to reset your password", nil, c)
					return
				} else {
					log.Println(e)
					common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
					return
				}

			} else {
				log.Println(addOtpRecord.Error)
				common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}

		} else {
			if errors.Is(findUser.Error, gorm.ErrRecordNotFound) == true {
				common.SuccessResponse(http.StatusOK, "If the provided email address is in our database, we will send you an email to reset your password", nil, c)
				return

			} else {
				log.Println(findUser.Error)
				common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}

		}

	} else {
		common.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}

}

func (a AuthController) ResetPassword(c *gin.Context) {
	var resetPasswordInput input.ResetPassword

	if err := c.ShouldBind(&resetPasswordInput); err == nil {
		var otpData models.Otp

		query := &models.Otp{
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

					updateUser := a.Db.Model(&models.User{}).Where("email = ?", otpData.Email).Update("password", newPassword)
					if updateUser.Error == nil {
						updateOtpTable := a.Db.Model(&otpData).Update("status", usedOtp)

						if updateOtpTable.Error == nil {
							common.SuccessResponse(http.StatusOK, "Password reset successful", nil, c)
							return
						} else {
							log.Println(updateOtpTable.Error)
							common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
							return
						}

					} else {
						log.Println(updateUser.Error)
						common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
						return
					}

				} else {
					log.Println(err)
					common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
					return
				}
			} else {
				common.ErrorResponse(http.StatusBadRequest, "otp expired", c)
				return
			}

		} else {
			if errors.Is(findUserOtp.Error, gorm.ErrRecordNotFound) == true {
				common.ErrorResponse(http.StatusBadRequest, "invalid otp", c)
				return
			}
			log.Println(findUserOtp.Error)
			common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
			return
		}

	} else {
		common.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}
}

func (a AuthController) ValidateOtp(c *gin.Context) {
	var validateOtpInput input.ValidateOtp

	if err := c.ShouldBind(&validateOtpInput); err == nil {
		var otpData models.Otp

		findUserOtp := a.Db.Where(map[string]interface{}{"otp": validateOtpInput.Otp, "status": 0}).First(&otpData)

		if findUserOtp.Error == nil {
			// determine if otp is for registration then activate user account
			if otpData.Purpose == registration {
				if otpData.ExpiryDate.Unix() > time.Now().Local().Unix() {
					updateOtpTable := a.Db.Model(&otpData).Update("status", usedOtp)

					if updateOtpTable.Error == nil {
						updateUser := a.Db.Model(&models.User{}).Where("email = ?", otpData.Email).Update("status", 1)
						if updateUser.Error == nil {
							common.SuccessResponse(http.StatusOK, "Otp validation successful", nil, c)
							return
						} else {
							common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
							return
						}
					} else {
						log.Println(updateOtpTable.Error)
						common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
						return
					}
				} else {
					common.ErrorResponse(http.StatusBadRequest, "otp expired", c)
					return
				}
			} else {
				if otpData.ExpiryDate.Unix() > time.Now().Local().Unix() {
					updateOtpTable := a.Db.Model(&otpData).Update("status", usedOtp)

					if updateOtpTable.Error == nil {
						common.SuccessResponse(http.StatusOK, "Otp validation successful", nil, c)
						return
					} else {
						log.Println(updateOtpTable.Error)
						common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
						return
					}
				} else {
					common.ErrorResponse(http.StatusBadRequest, "otp expired", c)
					return
				}
			}
		} else {
			if errors.Is(findUserOtp.Error, gorm.ErrRecordNotFound) == true {
				common.ErrorResponse(http.StatusBadRequest, "invalid otp", c)
				return
			} else {
				common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}
		}

	} else {
		common.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}

}

func (a AuthController) ResendOtp(c *gin.Context) {
	var resendOtpInput input.ResendOtpInput

	if err := c.ShouldBind(&resendOtpInput); err == nil {
		var purpose string
		var resetPasswordTruthiness = true

		if resendOtpInput.IsPasswordReset == &resetPasswordTruthiness {
			purpose = forgetPassword
		} else {
			purpose = registration
		}

		var user models.User

		query := &models.User{
			Email: resendOtpInput.Email,
		}
		findUser := a.Db.Where(query).First(&user)

		if findUser.Error == nil {
			otp := utils.GenerateOTP()

			otpPayload := &models.Otp{
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
					er := utils.SendEmail(user.Email, body, subject)
					if er == nil {
						common.SuccessResponse(http.StatusOK, "Otp sent successfully", nil, c)
						return
					} else {
						common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
						return
					}
				} else {
					subject := "Registration"
					body := utils.GenerateOTPEmailTemplate(otp)
					e := utils.SendEmail(user.Email, body, subject)
					if e == nil {
						common.SuccessResponse(http.StatusOK, "Otp sent successfully", nil, c)
						return
					} else {
						common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
						return
					}
				}

			} else {
				log.Println(addOtpRecord.Error)
				common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}

		} else {
			if errors.Is(findUser.Error, gorm.ErrRecordNotFound) == true {
				common.ErrorResponse(http.StatusBadRequest, "user with email not found", c)
				return
			} else {
				common.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}
		}
	} else {
		common.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}

}
