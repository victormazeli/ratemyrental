package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"math"
	"net/http"
	"rateMyRentalBackend/config"
	models2 "rateMyRentalBackend/database/models"
	request2 "rateMyRentalBackend/http/request"
	"rateMyRentalBackend/http/response"
	"rateMyRentalBackend/http/services"
	"strconv"
)

type PropertyController struct {
	Db              *gorm.DB
	Env             *config.Env
	PropertyService services.PropertyService
}

// AddNewProperty Adds new property
func (p PropertyController) AddNewProperty(c *gin.Context) {
	userId, _ := c.Get("user")

	var propertyInput request2.PropertyInput
	if err := c.ShouldBind(&propertyInput); err != nil {
		response.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}

	// Check if the property already exists
	var property models2.Property
	findProperty := p.Db.Where("property_title = ?", propertyInput.PropertyTitle).First(&property)
	if findProperty.RowsAffected > 0 {
		response.ErrorResponse(http.StatusBadRequest, "property already exists", c)
		return
	}

	// Check if images is empty
	if len(propertyInput.Images) == 0 {
		response.ErrorResponse(http.StatusBadRequest, "Property images can not be empty", c)
		return
	}

	propertyPayload := map[string]interface{}{
		"postal_code":    propertyInput.PostalCode,
		"country":        propertyInput.Country,
		"state":          propertyInput.State,
		"latitude":       propertyInput.Latitude,
		"longitude":      propertyInput.Longitude,
		"property_title": propertyInput.PropertyTitle,
		"description":    propertyInput.Description,
		"is_visible":     1,
		"user_id":        userId,
	}

	// Create the new property
	if err := p.Db.Model(&models2.Property{}).Create(&propertyPayload).Last(&property).Error; err != nil {
		fmt.Print(err)
		response.ErrorResponse(http.StatusInternalServerError, "An error occurred while creating the property", c)
		return
	}

	// Create property images
	var propertyImages []models2.PropertyImage
	for _, image := range propertyInput.Images {
		propertyImages = append(propertyImages, models2.PropertyImage{ImageUrl: image, PropertyID: property.ID})
	}
	if err := p.Db.Create(&propertyImages).Error; err != nil {
		response.ErrorResponse(http.StatusInternalServerError, "An error occurred while creating property images", c)
		return
	}

	response.SuccessResponse(http.StatusOK, "Property added successfully", property, c)
}

// GetProperty Get a single property
func (p PropertyController) GetProperty(c *gin.Context) {
	propertyId := c.Param("id")
	id, err := strconv.ParseInt(propertyId, 10, 64) // Use base 10 and int64 for id

	if err != nil {
		response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		return
	}

	var property models2.Property
	findProperty := p.Db.Preload("PropertyImages").First(&property, id)

	if findProperty.Error != nil {
		if errors.Is(findProperty.Error, gorm.ErrRecordNotFound) {
			response.ErrorResponse(http.StatusNotFound, "property not found", c)
		} else {
			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		}
		return
	}

	response.SuccessResponse(http.StatusOK, "property fetched successfully", property, c)
}

//func (p PropertyController) UploadImageProperty(c *gin.Context) {
//	propertyId := c.PostForm("property_id")
//
//	parsedPropertyId, err := strconv.ParseInt(propertyId, 0, 8)
//
//	if err == nil {
//		formData, er := c.MultipartForm()
//		var filePaths []string
//		var propertyImages []models.PropertyImage
//
//		if er == nil {
//			files := formData.File["files"]
//
//			if files == nil {
//				response.ErrorResponse(http.StatusBadRequest, "Please select the input file", c)
//				return
//			}
//
//			if propertyId == "" {
//				response.ErrorResponse(http.StatusBadRequest, "property_id is missing in request", c)
//				return
//			}
//
//			for _, file := range files {
//				fileExt := filepath.Ext(file.Filename)
//				originalFileName := strings.TrimSuffix(filepath.Base(file.Filename), filepath.Ext(file.Filename))
//				filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", time.Now().Unix()) + fileExt
//				if e := c.SaveUploadedFile(file, "./public/image/"+filename); e != nil {
//					response.ErrorResponse(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", e.Error()), c)
//					return
//				}
//				filePath := p.Env.BaseUrl + "/static/image/" + filename
//
//				filePaths = append(filePaths, filePath)
//				// attach to property with property id,
//				propertyImages = append(propertyImages, models.PropertyImage{ImageUrl: filePath, PropertyID: uint(parsedPropertyId)})
//			}
//			// store image url  in database
//			result := p.Db.Create(&propertyImages)
//
//			if result.Error == nil {
//				response.SuccessResponse(http.StatusOK, "upload successful", filePaths, c)
//				return
//			} else {
//				response.ErrorResponse(http.StatusInternalServerError, "property_id is missing in request", c)
//				return
//			}
//		} else {
//			response.ErrorResponse(http.StatusBadRequest, er.Error(), c)
//			return
//		}
//
//	} else {
//		response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
//		return
//	}
//}

func (p PropertyController) GetAllProperties(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")

	visibleProperties := map[string]interface{}{"properties.is_visible": 1}

	fetchProperties, err := p.PropertyService.Pagination(page, limit, visibleProperties)
	log.Print(err)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrorResponse(http.StatusNotFound, "no properties found", c)
		} else {
			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		}
		return
	}

	response.SuccessResponse(http.StatusOK, "properties fetched successfully", fetchProperties, c)
}

//func (p PropertyController) RateProperty(c *gin.Context) {
//	// to be implemented
//	var rateInput request2.RateInput
//	var property models.Property
//	var fiveRatingResponse int64
//	var fourRatingResponse int64
//	var threeRatingResponse int64
//	var twoRatingResponse int64
//	var oneRatingResponse int64
//
//	if err := c.ShouldBind(&rateInput); err == nil {
//		findProperty := p.Db.Model(&models.Property{}).Where("id = ?", rateInput.PropertyID).First(&property)
//		if findProperty.Error == nil {
//			// store rating into db
//			ratePayload := &models.Rating{
//				PropertyID: rateInput.PropertyID,
//				Score:      rateInput.Score,
//				Feature:    rateInput.Feature,
//			}
//			result := p.Db.Model(&models.Rating{}).Create(ratePayload)
//			if result.Error == nil {
//				p.Db.Model(&models.Rating{}).Where("score = ?", 1).Count(&oneRatingResponse)
//				p.Db.Model(&models.Rating{}).Where("score = ?", 2).Count(&twoRatingResponse)
//				p.Db.Model(&models.Rating{}).Where("score = ?", 3).Count(&threeRatingResponse)
//				p.Db.Model(&models.Rating{}).Where("score = ?", 4).Count(&fourRatingResponse)
//				p.Db.Model(&models.Rating{}).Where("score = ?", 5).Count(&fiveRatingResponse)
//
//				totalSum := fiveRatingResponse + fourRatingResponse + threeRatingResponse + twoRatingResponse + oneRatingResponse
//				totalRatingSum := (5 * fiveRatingResponse) + (4 * fourRatingResponse) + (3 * threeRatingResponse) + (2 * twoRatingResponse) + (1 * oneRatingResponse)
//
//				score := float64(totalRatingSum) / float64(totalSum)
//
//				averageScore := math.Ceil(score)
//				// update property
//				updateResult := p.Db.Model(&models.Property{}).Where("id = ?", rateInput.PropertyID).Update("average_rating", averageScore)
//
//				if updateResult.Error == nil {
//					response.SuccessResponse(http.StatusOK, "Rating added successfully", nil, c)
//					return
//				} else {
//					response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
//					return
//				}
//
//			} else {
//				response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
//				return
//			}
//
//		} else {
//			if errors.Is(findProperty.Error, gorm.ErrRecordNotFound) == true {
//				response.ErrorResponse(http.StatusNotFound, "property not found", c)
//				return
//
//			} else {
//				response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
//				return
//			}
//		}
//	} else {
//		response.ErrorResponse(http.StatusBadRequest, err.Error(), c)
//		return
//	}
//
//}

func (p PropertyController) RateProperty(c *gin.Context) {
	var rateInput request2.RateInput

	if err := c.ShouldBind(&rateInput); err != nil {
		response.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}

	var property models2.Property
	if err := p.Db.First(&property, rateInput.PropertyID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrorResponse(http.StatusNotFound, "property not found", c)
		} else {
			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		}
		return
	}

	ratePayload := &models2.Rating{
		PropertyID: rateInput.PropertyID,
		Score:      rateInput.Score,
		Feature:    rateInput.Feature,
	}

	if err := p.Db.Create(ratePayload).Error; err != nil {
		response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		return
	}

	var ratingCount []struct {
		Score int
		Count int64
	}

	if err := p.Db.Model(&models2.Rating{}).Select("score, COUNT(*) as count").Group("score").Scan(&ratingCount).Error; err != nil {
		response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		return
	}

	var totalSum int64
	var totalRatingSum int64
	for _, r := range ratingCount {
		totalSum += r.Count
		totalRatingSum += int64(r.Score) * r.Count
	}

	score := float64(totalRatingSum) / float64(totalSum)
	averageScore := math.Round(score)

	if err := p.Db.Model(&models2.Property{}).Where("id = ?", rateInput.PropertyID).Update("average_rating", averageScore).Error; err != nil {
		response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		return
	}

	response.SuccessResponse(http.StatusOK, "Rating added successfully", nil, c)
}

//func (p PropertyController) UpdateSingleImageProperty(c *gin.Context) {
//	// to be implemented
//	// add the image id
//	// check if image exist
//	//propertyId := c.PostForm("property_id")
//	imageId := c.PostForm("image_id")
//	//var parsedPropertyId int64
//	//var parsedimageId int64
//	var propertyImage models.PropertyImage
//
//	formData, er := c.FormFile("file")
//
//	if er == nil {
//
//		if formData == nil {
//			response.ErrorResponse(http.StatusBadRequest, "Please select the input file", c)
//			return
//		}
//
//		//if propertyId == "" {
//		//	parsedPropertyId, _ = strconv.ParseInt(propertyId, 0, 8)
//		//	common.ErrorResponse(http.StatusBadRequest, "property_id is missing in request", c)
//		//	return
//		//}
//
//		if imageId == "" {
//			//parsedimageId, _ = strconv.ParseInt(imageId, 0, 8)
//			response.ErrorResponse(http.StatusBadRequest, "image_id is missing in request", c)
//			return
//		}
//
//		// check if image exist
//		imageExistResult := p.Db.Model(&models.PropertyImage{}).Where("id = ?", imageId).First(&propertyImage)
//
//		if imageExistResult.Error == nil {
//			fileExt := filepath.Ext(formData.Filename)
//			originalFileName := strings.TrimSuffix(filepath.Base(formData.Filename), filepath.Ext(formData.Filename))
//			filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", time.Now().Unix()) + fileExt
//			if e := c.SaveUploadedFile(formData, "./public/image/"+filename); e != nil {
//				response.ErrorResponse(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", e.Error()), c)
//				return
//			}
//			filePath := p.Env.BaseUrl + "/static/image/" + filename
//
//			result := p.Db.Model(&models.PropertyImage{}).Where("id = ?", imageId).Update("image_url", filePath)
//
//			if result.Error == nil {
//				response.SuccessResponse(http.StatusOK, "upload successful", filePath, c)
//				return
//			} else {
//				response.ErrorResponse(http.StatusInternalServerError, "property_id is missing in request", c)
//				return
//			}
//
//		} else {
//			if errors.Is(imageExistResult.Error, gorm.ErrRecordNotFound) == true {
//				response.ErrorResponse(http.StatusNotFound, "image not found", c)
//				return
//
//			} else {
//				response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
//				return
//			}
//
//		}
//
//	} else {
//		response.ErrorResponse(http.StatusBadRequest, er.Error(), c)
//		return
//	}
//
//}

func (p PropertyController) UpdateSingleImageProperty(c *gin.Context) {
	propertyImageId := c.Param("id")
	id, err := strconv.ParseInt(propertyImageId, 10, 64) // Use base 10 and int64 for id

	if err != nil {
		response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		return
	}

	var propertyImageUpload request2.PropertyImagesUpload

	if err := c.ShouldBind(&propertyImageUpload); err != nil {
		response.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}

	updateImage := p.Db.Model(&models2.PropertyImage{}).Where("id = ?", id).Updates(propertyImageUpload)

	if updateImage.Error != nil {
		response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		return
	}

	response.SuccessResponse(http.StatusOK, "image updated successfully", nil, c)
	return

}

//func (p PropertyController) SearchProperties(c *gin.Context) {
//
//}

func (p PropertyController) GetPropertyTypes(c *gin.Context) {
	var propertyTypes []models2.PropertyType
	if err := p.Db.Find(&propertyTypes).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrorResponse(http.StatusNotFound, "property types not found", c)
		} else {
			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		}
		return
	}

	response.SuccessResponse(http.StatusOK, "property types fetched successfully", propertyTypes, c)
}

func (p PropertyController) GetPropertyDetachedTypes(c *gin.Context) {
	var propertyDetachedTypes []models2.PropertyDetachedType
	if err := p.Db.Find(&propertyDetachedTypes).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrorResponse(http.StatusNotFound, "property detached types not found", c)
		} else {
			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		}
		return
	}
	response.SuccessResponse(http.StatusOK, "property detached types fetched successfully", propertyDetachedTypes, c)

}

// UpdatePropertyDetail update property details
//func (p PropertyController) UpdatePropertyDetail(c *gin.Context) {
//	propertyId := c.Param("id")
//	var property models.Property
//	var updateInput request2.PropertyUpdateInput
//
//	if err := c.ShouldBind(&updateInput); err == nil {
//		findProperty := p.Db.Model(&models.Property{}).Where("id = ?", propertyId).First(&property)
//
//		if findProperty.Error == nil {
//			// if property found; proceed to update property
//			//propertyUpdateData := map[string]interface{}{
//			//	"postal_code":            updateInput.PostalCode,
//			//	"country":                updateInput.Country,
//			//	"state":                  updateInput.State,
//			//	"location":               updateInput.Location,
//			//	"property_title":         updateInput.PropertyTitle,
//			//	"description":            updateInput.Description,
//			//	"floors":                 updateInput.Floors,
//			//	"number_of_rooms":        updateInput.NumberOfRooms,
//			//	"bed_rooms":              updateInput.BedRooms,
//			//	"bath_rooms":             updateInput.BathRooms,
//			//	"cloak_rooms":            updateInput.CloakRooms,
//			//	"utility_rooms":          updateInput.UtilityRooms,
//			//	"conservatory":           updateInput.Conservatory,
//			//	"entrance_hall":          updateInput.EntranceHall,
//			//	"front_yard":             updateInput.FrontYard,
//			//	"mud_room":               updateInput.MudRoom,
//			//	"furnished_room":         updateInput.FurnishedRoom,
//			//	"garden":                 updateInput.Garden,
//			//	"garage":                 updateInput.Garage,
//			//	"ensuite":                updateInput.Ensuite,
//			//	"character_feature":      updateInput.CharacterFeature,
//			//	"epc_ratings":            updateInput.EpcRatings,
//			//	"pets_allowed":           updateInput.PetsAllowed,
//			//	"smoking_allowed":        updateInput.SmokingAllowed,
//			//	"dss_allowed":            updateInput.DssAllowed,
//			//	"sharers_allowed":        updateInput.SharersAllowed,
//			//	"property_type":          updateInput.PropertyType,
//			//	"property_detached_type": updateInput.PropertyDetachedType,
//			//}
//			updateProperty := p.Db.Model(&models.Property{}).Where("id = ?", propertyId).Updates(updateInput)
//
//			if updateProperty.Error == nil {
//				response.SuccessResponse(http.StatusOK, "properties updated successfully", nil, c)
//				return
//			} else {
//				response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
//				return
//			}
//
//		} else {
//			if errors.Is(findProperty.Error, gorm.ErrRecordNotFound) == true {
//				response.ErrorResponse(http.StatusNotFound, "property not found", c)
//				return
//			} else {
//				response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
//				return
//			}
//
//		}
//
//	} else {
//		response.ErrorResponse(http.StatusBadRequest, err.Error(), c)
//	}
//
//}

func (p PropertyController) UpdatePropertyDetail(c *gin.Context) {
	propertyId := c.Param("id")
	var property models2.Property
	var updateInput request2.PropertyUpdateInput

	if err := c.ShouldBind(&updateInput); err != nil {
		response.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}

	if err := p.Db.First(&property, propertyId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrorResponse(http.StatusNotFound, "property not found", c)
		} else {
			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		}
		return
	}

	if err := p.Db.Model(&property).Updates(updateInput).Error; err != nil {
		response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		return
	}

	response.SuccessResponse(http.StatusOK, "property updated successfully", nil, c)
}

//func (p PropertyController) TogglePropertyVisibility(c *gin.Context) {
//	var propertyVisibility request2.PropertyVisibility
//	var property models.Property
//
//	if err := c.ShouldBind(&propertyVisibility); err == nil {
//		// check if property exist
//		checkPropertyExist := p.Db.Model(&models.Property{}).Where("id = ?", propertyVisibility.PropertyID).First(&property)
//
//		if checkPropertyExist.Error == nil {
//			// check if property is already visible; if visible, make it invisible
//			checkPropertyIsVisible := p.Db.Model(&models.Property{}).Where("id = ?", propertyVisibility.PropertyID).Where("is_visible", 1).First(&property)
//			if checkPropertyIsVisible.Error == nil {
//				makePropertyInvisible := p.Db.Model(&models.Property{}).Where("id = ?", propertyVisibility.PropertyID).Update("is_visible", 0)
//				if makePropertyInvisible.Error == nil {
//					response.SuccessResponse(http.StatusOK, "property visibility updated successfully", nil, c)
//					return
//				} else {
//					response.ErrorResponse(http.StatusInternalServerError, "An error occurred..", c)
//					return
//				}
//
//			} else {
//				if errors.Is(checkPropertyIsVisible.Error, gorm.ErrRecordNotFound) {
//					makePropertyVisible := p.Db.Model(&models.Property{}).Where("id = ?", propertyVisibility.PropertyID).Update("is_visible", 1)
//					if makePropertyVisible.Error == nil {
//						response.SuccessResponse(http.StatusOK, "property visibility updated successfully", nil, c)
//						return
//					} else {
//						response.ErrorResponse(http.StatusInternalServerError, "An error occurred.", c)
//						return
//					}
//				} else {
//					response.ErrorResponse(http.StatusInternalServerError, "An error occurred!!", c)
//					return
//				}
//			}
//
//		} else {
//			if errors.Is(checkPropertyExist.Error, gorm.ErrRecordNotFound) == true {
//
//				response.ErrorResponse(http.StatusNotFound, "property not found", c)
//				return
//
//			} else {
//				response.ErrorResponse(http.StatusInternalServerError, "An error occurred!", c)
//				return
//			}
//
//		}
//
//	} else {
//		response.ErrorResponse(http.StatusBadRequest, err.Error(), c)
//		return
//	}
//
//}

func (p PropertyController) TogglePropertyVisibility(c *gin.Context) {
	var propertyVisibility request2.PropertyVisibility
	var property models2.Property

	if err := c.ShouldBind(&propertyVisibility); err != nil {
		response.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}

	if err := p.Db.First(&property, propertyVisibility.PropertyID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrorResponse(http.StatusNotFound, "property not found", c)
		} else {
			response.ErrorResponse(http.StatusInternalServerError, "An error occurred!", c)
		}
		return
	}

	// Toggle property visibility
	if property.IsVisible == 1 {
		property.IsVisible = 0
	} else {
		property.IsVisible = 1
	}

	if err := p.Db.Save(&property).Error; err != nil {
		response.ErrorResponse(http.StatusInternalServerError, "An error occurred!", c)
		return
	}

	response.SuccessResponse(http.StatusOK, "property visibility updated successfully", nil, c)
}

//func (p PropertyController) AddOrRemoveFavoriteProperty(c *gin.Context) {
//	userId, _ := c.Get("user")
//	var addFavorite request2.FavoritePropertyInput
//	var favoriteProperty models.FavoriteProperty
//	deletedAt := time.Now()
//	if err := c.ShouldBind(&addFavorite); err == nil {
//		// check if already property is already added
//		findProperty := p.Db.Model(&models.FavoriteProperty{}).Where("property_id = ?", addFavorite.PropertyID).Where("user_id = ?", userId).Where("deleted_at IS NULL").First(&favoriteProperty)
//		// if no error and property found in favorite, remove property from favorite
//		if findProperty.Error == nil {
//			// remove property
//			result := p.Db.Model(&models.FavoriteProperty{}).Where("property_id = ?", addFavorite.PropertyID).Where("user_id = ?", userId).Update("deleted_at", deletedAt)
//
//			if result.Error == nil {
//				response.SuccessResponse(http.StatusOK, "property removed from favorite", nil, c)
//				return
//			} else {
//				response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
//				return
//			}
//		} else {
//			// check if property already exist update its deleted status to null
//			if errors.Is(findProperty.Error, gorm.ErrRecordNotFound) {
//				checkProperty := p.Db.Model(&models.FavoriteProperty{}).Where("property_id = ?", addFavorite.PropertyID).Where("user_id = ?", userId).Not("deleted_at", nil).First(&favoriteProperty)
//				if checkProperty.Error == nil {
//					UpdateResult := p.Db.Model(&models.FavoriteProperty{}).Where("property_id = ?", addFavorite.PropertyID).Where("user_id = ?", userId).Updates(map[string]interface{}{"deleted_at": nil})
//					if UpdateResult.Error == nil {
//						response.SuccessResponse(http.StatusOK, "property added to favorite", nil, c)
//						return
//					} else {
//						response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
//						return
//					}
//				} else {
//					// else if error is equals to not found, add property to favorite
//					if errors.Is(checkProperty.Error, gorm.ErrRecordNotFound) == true {
//						addProperty := map[string]interface{}{
//							"property_id": addFavorite.PropertyID,
//							"user_id":     userId,
//						}
//						addPropertyResult := p.Db.Model(&models.FavoriteProperty{}).Create(addProperty)
//						if addPropertyResult.Error == nil {
//							response.SuccessResponse(http.StatusOK, "property added to favorite", nil, c)
//							return
//
//						} else {
//							response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
//							return
//						}
//
//					} else {
//						response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
//						return
//					}
//				}
//
//			} else {
//				response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
//				return
//			}
//
//		}
//
//	} else {
//		response.ErrorResponse(http.StatusBadRequest, err.Error(), c)
//		return
//	}
//
//}

func (p PropertyController) AddOrRemoveFavoriteProperty(c *gin.Context) {
	userId, _ := c.Get("user")
	var addFavorite request2.FavoritePropertyInput
	if err := c.ShouldBind(&addFavorite); err != nil {
		response.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}

	var favoriteProperty models2.FavoriteProperty
	result := p.Db.Where("property_id = ?", addFavorite.PropertyID).Where("user_id = ?", userId).First(&favoriteProperty)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Property not found in favorites, so add it
			newFavorite := models2.FavoriteProperty{
				PropertyID: addFavorite.PropertyID,
				UserID:     userId.(uint),
			}
			if err := p.Db.Create(&newFavorite).Error; err != nil {
				response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}
			response.SuccessResponse(http.StatusOK, "property added to favorite", nil, c)
		} else {
			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		}
		return
	}

	// Property found in favorites, so remove it by soft-deleting (setting deleted_at)
	if err := p.Db.Delete(&favoriteProperty, addFavorite.PropertyID).Error; err != nil {
		response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		return
	}
	response.SuccessResponse(http.StatusOK, "property removed from favorite", nil, c)
}

func (p PropertyController) GetUserFavoriteProperties(c *gin.Context) {
	userId, _ := c.Get("user")

	var properties []*models2.Property
	results := p.Db.Joins("left join favorite_properties on favorite_properties.property_id = properties.id").
		Preload("PropertyImages").
		Where("favorite_properties.user_id = ?", userId).
		Find(&properties)

	if results.Error != nil {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			response.ErrorResponse(http.StatusNotFound, "no properties found", c)
		} else {
			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		}
		return
	}

	response.SuccessResponse(http.StatusOK, "properties fetched successfully", properties, c)
}

func (p PropertyController) GetUserUploadedProperties(c *gin.Context) {
	userId, _ := c.Get("user")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	fetchProperties, err := p.PropertyService.Pagination(page, limit, map[string]interface{}{"user_id": userId})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrorResponse(http.StatusNotFound, "no properties found", c)
		} else {
			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		}
		return
	}
	response.SuccessResponse(http.StatusOK, "properties fetched successfully", fetchProperties, c)
}

func (p PropertyController) PropertyRecommendations(c *gin.Context) {
	userId, _ := c.Get("user")
	var user models2.User
	var properties []models2.Property

	if err := p.Db.Preload("PropertyImages").Find(&properties).Limit(300).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrorResponse(http.StatusNotFound, "no properties found", c)
		} else {
			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		}
		return
	}

	if err := p.Db.Model(&models2.User{}).Where("id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrorResponse(http.StatusNotFound, "user not found", c)
		} else {
			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		}
		return
	}
	recommendation := Recommend(properties, user)

	response.SuccessResponse(http.StatusOK, "recommendation fetched successfully", recommendation, c)

}

func Recommend(destinations []models2.Property, currDest models2.User) []models2.Property {
	var similarPlaces []models2.Property

	for _, d := range destinations {
		sim := ComputeSimilarity(d, currDest)

		fmt.Printf("Distance to %v: %f\n", d.Address, sim)

		if sim > 0.4 {
			similarPlaces = append(similarPlaces, d)
		}

	}

	return similarPlaces

}

func ComputeSimilarity(item1 models2.Property, item2 models2.User) float64 {
	// Compute the distance between the two items using their geoPoints
	lat1 := item1.Latitude
	lon1 := item1.Longitude
	lat2 := item2.Latitude
	lon2 := item2.Longitude
	dist := haversine(lat1, lon1, lat2, lon2)

	// Compute the similarity score as the inverse of the distance
	return 1 / (1 + dist)
}

// haversine formula implementation
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const r = 6371e3 // Earth radius in meters
	phi1 := toRadians(lat1)
	phi2 := toRadians(lat2)
	deltaPhi := toRadians(lat2 - lat1)
	deltaLambda := toRadians(lon2 - lon1)

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return r * c
}

func toRadians(deg float64) float64 {
	return deg * math.Pi / 180
}
