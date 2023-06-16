package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math"
	"net/http"
	"path/filepath"
	"rateMyRentalBackend/common"
	"rateMyRentalBackend/config"
	request2 "rateMyRentalBackend/http/request"
	"rateMyRentalBackend/http/response"
	"rateMyRentalBackend/models"
	"strconv"
	"strings"
	"time"
)

type PropertyController struct {
	Db  *gorm.DB
	Env *config.Env
}

// AddNewProperty Adds new property
func (p PropertyController) AddNewProperty(c *gin.Context) {
	userId, _ := c.Get("user")
	var propertyInput request2.PropertyInput
	var property models.Property
	var propertyImages []models.PropertyImage

	if err := c.ShouldBind(&propertyInput); err == nil {
		// check if property already exist
		findProperty := p.Db.Model(&models.Property{}).Where("property_title = ?", propertyInput.PropertyTitle).First(&property)
		// check db errors
		if findProperty.Error == nil {
			response.ErrorResponse(http.StatusBadRequest, "property already exist", c)
			return
		} else {
			if errors.Is(findProperty.Error, gorm.ErrRecordNotFound) == true {
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

				addProperty := p.Db.Model(&models.Property{}).Create(&propertyPayload).Last(&property)

				if addProperty.Error == nil {

					for _, image := range propertyInput.Images {
						propertyImages = append(propertyImages, models.PropertyImage{ImageUrl: image, PropertyID: property.ID})

					}

					result := p.Db.Create(&propertyImages)

					if result.Error == nil {
						response.SuccessResponse(http.StatusOK, "property added successful", property, c)

					} else {
						response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
					}

				} else {
					response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				}

			} else {
				response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
			}
		}

	} else {
		response.ErrorResponse(http.StatusBadRequest, err.Error(), c)
	}

}

// GetProperty Get a single property
func (p PropertyController) GetProperty(c *gin.Context) {
	var property models.Property
	propertyId := c.Param("id")
	id, err := strconv.ParseInt(propertyId, 0, 8)

	if err == nil {
		findProperty := p.Db.Preload("PropertyImages").Find(&property, id)

		if findProperty.Error == nil {
			response.SuccessResponse(http.StatusOK, "property fetched successfully", property, c)
			return

		} else {
			if errors.Is(findProperty.Error, gorm.ErrRecordNotFound) == true {
				response.ErrorResponse(http.StatusNotFound, "property not found", c)
				return
			} else {
				response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}
		}

	} else {
		response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		return
	}

}

func (p PropertyController) UploadImageProperty(c *gin.Context) {
	propertyId := c.PostForm("property_id")

	parsedPropertyId, err := strconv.ParseInt(propertyId, 0, 8)

	if err == nil {
		formData, er := c.MultipartForm()
		var filePaths []string
		var propertyImages []models.PropertyImage

		if er == nil {
			files := formData.File["files"]

			if files == nil {
				response.ErrorResponse(http.StatusBadRequest, "Please select the input file", c)
				return
			}

			if propertyId == "" {
				response.ErrorResponse(http.StatusBadRequest, "property_id is missing in request", c)
				return
			}

			for _, file := range files {
				fileExt := filepath.Ext(file.Filename)
				originalFileName := strings.TrimSuffix(filepath.Base(file.Filename), filepath.Ext(file.Filename))
				filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", time.Now().Unix()) + fileExt
				if e := c.SaveUploadedFile(file, "./public/image/"+filename); e != nil {
					response.ErrorResponse(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", e.Error()), c)
					return
				}
				filePath := p.Env.BaseUrl + "/static/image/" + filename

				filePaths = append(filePaths, filePath)
				// attach to property with property id,
				propertyImages = append(propertyImages, models.PropertyImage{ImageUrl: filePath, PropertyID: uint(parsedPropertyId)})
			}
			// store image url  in database
			result := p.Db.Create(&propertyImages)

			if result.Error == nil {
				response.SuccessResponse(http.StatusOK, "upload successful", filePaths, c)
				return
			} else {
				response.ErrorResponse(http.StatusInternalServerError, "property_id is missing in request", c)
				return
			}
		} else {
			response.ErrorResponse(http.StatusBadRequest, er.Error(), c)
			return
		}

	} else {
		response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
		return
	}
}

func (p PropertyController) GetAllProperties(c *gin.Context) {
	//userId, _ := c.Get("user")

	page := c.DefaultQuery("page", "1")

	limit := c.DefaultQuery("limit", "10")

	//var user models.User

	//findUser := p.Db.Model(&models.User{}).Where("id = ?", userId).First(&user)

	preloadQuery := "PropertyImages"

	fetchProperties, err := common.Pagination("properties", page, limit, "created_at desc", map[string]interface{}{"is_visible": 1, "deleted_at": nil}, &preloadQuery, p.Db)

	if err == nil {
		response.SuccessResponse(http.StatusOK, "properties fetched successfully", fetchProperties, c)
		return
	} else {
		if errors.Is(err, gorm.ErrRecordNotFound) == true {

			response.ErrorResponse(http.StatusNotFound, "no properties found", c)
			return

		} else {

			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
			return
		}
	}

	//if findUser.Error == nil {

	//} else {
	//	if errors.Is(findUser.Error, gorm.ErrRecordNotFound) == true {
	//
	//		response.ErrorResponse(http.StatusNotFound, "user not found", c)
	//		return
	//
	//	} else {
	//
	//		response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
	//		return
	//	}
	//}

}

func (p PropertyController) RateProperty(c *gin.Context) {
	// to be implemented
	var rateInput request2.RateInput
	var property models.Property
	var fiveRatingResponse int64
	var fourRatingResponse int64
	var threeRatingResponse int64
	var twoRatingResponse int64
	var oneRatingResponse int64

	if err := c.ShouldBind(&rateInput); err == nil {
		findProperty := p.Db.Model(&models.Property{}).Where("id = ?", rateInput.PropertyID).First(&property)
		if findProperty.Error == nil {
			// store rating into db
			ratePayload := &models.Rating{
				PropertyID: rateInput.PropertyID,
				Score:      rateInput.Score,
				Feature:    rateInput.Feature,
			}
			result := p.Db.Model(&models.Rating{}).Create(ratePayload)
			if result.Error == nil {
				p.Db.Model(&models.Rating{}).Where("score = ?", 1).Count(&oneRatingResponse)
				p.Db.Model(&models.Rating{}).Where("score = ?", 2).Count(&twoRatingResponse)
				p.Db.Model(&models.Rating{}).Where("score = ?", 3).Count(&threeRatingResponse)
				p.Db.Model(&models.Rating{}).Where("score = ?", 4).Count(&fourRatingResponse)
				p.Db.Model(&models.Rating{}).Where("score = ?", 5).Count(&fiveRatingResponse)

				totalSum := fiveRatingResponse + fourRatingResponse + threeRatingResponse + twoRatingResponse + oneRatingResponse
				totalRatingSum := (5 * fiveRatingResponse) + (4 * fourRatingResponse) + (3 * threeRatingResponse) + (2 * twoRatingResponse) + (1 * oneRatingResponse)

				score := float64(totalRatingSum) / float64(totalSum)

				averageScore := math.Ceil(score)
				// update property
				updateResult := p.Db.Model(&models.Property{}).Where("id = ?", rateInput.PropertyID).Update("average_rating", averageScore)

				if updateResult.Error == nil {
					response.SuccessResponse(http.StatusOK, "Rating added successfully", nil, c)
					return
				} else {
					response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
					return
				}

			} else {
				response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}

		} else {
			if errors.Is(findProperty.Error, gorm.ErrRecordNotFound) == true {
				response.ErrorResponse(http.StatusNotFound, "property not found", c)
				return

			} else {
				response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}
		}
	} else {
		response.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}

}

func (p PropertyController) UpdateSingleImageProperty(c *gin.Context) {
	// to be implemented
	// add the image id
	// check if image exist
	//propertyId := c.PostForm("property_id")
	imageId := c.PostForm("image_id")
	//var parsedPropertyId int64
	//var parsedimageId int64
	var propertyImage models.PropertyImage

	formData, er := c.FormFile("file")

	if er == nil {

		if formData == nil {
			response.ErrorResponse(http.StatusBadRequest, "Please select the input file", c)
			return
		}

		//if propertyId == "" {
		//	parsedPropertyId, _ = strconv.ParseInt(propertyId, 0, 8)
		//	common.ErrorResponse(http.StatusBadRequest, "property_id is missing in request", c)
		//	return
		//}

		if imageId == "" {
			//parsedimageId, _ = strconv.ParseInt(imageId, 0, 8)
			response.ErrorResponse(http.StatusBadRequest, "image_id is missing in request", c)
			return
		}

		// check if image exist
		imageExistResult := p.Db.Model(&models.PropertyImage{}).Where("id = ?", imageId).First(&propertyImage)

		if imageExistResult.Error == nil {
			fileExt := filepath.Ext(formData.Filename)
			originalFileName := strings.TrimSuffix(filepath.Base(formData.Filename), filepath.Ext(formData.Filename))
			filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", time.Now().Unix()) + fileExt
			if e := c.SaveUploadedFile(formData, "./public/image/"+filename); e != nil {
				response.ErrorResponse(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", e.Error()), c)
				return
			}
			filePath := p.Env.BaseUrl + "/static/image/" + filename

			result := p.Db.Model(&models.PropertyImage{}).Where("id = ?", imageId).Update("image_url", filePath)

			if result.Error == nil {
				response.SuccessResponse(http.StatusOK, "upload successful", filePath, c)
				return
			} else {
				response.ErrorResponse(http.StatusInternalServerError, "property_id is missing in request", c)
				return
			}

		} else {
			if errors.Is(imageExistResult.Error, gorm.ErrRecordNotFound) == true {
				response.ErrorResponse(http.StatusNotFound, "image not found", c)
				return

			} else {
				response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}

		}

	} else {
		response.ErrorResponse(http.StatusBadRequest, er.Error(), c)
		return
	}

}

func (p PropertyController) SearchProperties(c *gin.Context) {

}

func (p PropertyController) GetPropertyTypes(c *gin.Context) {
	// to be implemented
	var propertyTypes []models.PropertyType
	results := p.Db.Model(&models.PropertyType{}).Find(&propertyTypes)
	if results.Error == nil {
		response.SuccessResponse(http.StatusOK, "property types fetched successfully", propertyTypes, c)
		return
	} else {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) == true {
			response.ErrorResponse(http.StatusNotFound, "property types not found", c)
			return
		} else {
			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
			return
		}

	}

}

func (p PropertyController) GetPropertyDetachedTypes(c *gin.Context) {
	// to be implemented
	var propertyDetachedTypes []models.PropertyDetachedType
	results := p.Db.Model(&models.PropertyDetachedType{}).Find(&propertyDetachedTypes)
	if results.Error == nil {
		response.SuccessResponse(http.StatusOK, "property detached types fetched successfully", propertyDetachedTypes, c)
		return
	} else {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) == true {
			response.ErrorResponse(http.StatusNotFound, "property detached types not found", c)
			return
		} else {
			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
			return
		}

	}

}

// UpdatePropertyDetail update property details
func (p PropertyController) UpdatePropertyDetail(c *gin.Context) {
	propertyId := c.Param("id")
	var property models.Property
	var updateInput request2.PropertyUpdateInput

	if err := c.ShouldBind(&updateInput); err == nil {
		findProperty := p.Db.Model(&models.Property{}).Where("id = ?", propertyId).First(&property)

		if findProperty.Error == nil {
			// if property found; proceed to update property
			//propertyUpdateData := map[string]interface{}{
			//	"postal_code":            updateInput.PostalCode,
			//	"country":                updateInput.Country,
			//	"state":                  updateInput.State,
			//	"location":               updateInput.Location,
			//	"property_title":         updateInput.PropertyTitle,
			//	"description":            updateInput.Description,
			//	"floors":                 updateInput.Floors,
			//	"number_of_rooms":        updateInput.NumberOfRooms,
			//	"bed_rooms":              updateInput.BedRooms,
			//	"bath_rooms":             updateInput.BathRooms,
			//	"cloak_rooms":            updateInput.CloakRooms,
			//	"utility_rooms":          updateInput.UtilityRooms,
			//	"conservatory":           updateInput.Conservatory,
			//	"entrance_hall":          updateInput.EntranceHall,
			//	"front_yard":             updateInput.FrontYard,
			//	"mud_room":               updateInput.MudRoom,
			//	"furnished_room":         updateInput.FurnishedRoom,
			//	"garden":                 updateInput.Garden,
			//	"garage":                 updateInput.Garage,
			//	"ensuite":                updateInput.Ensuite,
			//	"character_feature":      updateInput.CharacterFeature,
			//	"epc_ratings":            updateInput.EpcRatings,
			//	"pets_allowed":           updateInput.PetsAllowed,
			//	"smoking_allowed":        updateInput.SmokingAllowed,
			//	"dss_allowed":            updateInput.DssAllowed,
			//	"sharers_allowed":        updateInput.SharersAllowed,
			//	"property_type":          updateInput.PropertyType,
			//	"property_detached_type": updateInput.PropertyDetachedType,
			//}
			updateProperty := p.Db.Model(&models.Property{}).Where("id = ?", propertyId).Updates(updateInput)

			if updateProperty.Error == nil {
				response.SuccessResponse(http.StatusOK, "properties updated successfully", nil, c)
				return
			} else {
				response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}

		} else {
			if errors.Is(findProperty.Error, gorm.ErrRecordNotFound) == true {
				response.ErrorResponse(http.StatusNotFound, "property not found", c)
				return
			} else {
				response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}

		}

	} else {
		response.ErrorResponse(http.StatusBadRequest, err.Error(), c)
	}

}

func (p PropertyController) TogglePropertyVisibility(c *gin.Context) {
	var propertyVisibility request2.PropertyVisibility
	var property models.Property

	if err := c.ShouldBind(&propertyVisibility); err == nil {
		// check if property exist
		checkPropertyExist := p.Db.Model(&models.Property{}).Where("id = ?", propertyVisibility.PropertyID).First(&property)

		if checkPropertyExist.Error == nil {
			// check if property is already visible; if visible, make it invisible
			checkPropertyIsVisible := p.Db.Model(&models.Property{}).Where("id = ?", propertyVisibility.PropertyID).Where("is_visible", 1).First(&property)
			if checkPropertyIsVisible.Error == nil {
				makePropertyInvisible := p.Db.Model(&models.Property{}).Where("id = ?", propertyVisibility.PropertyID).Update("is_visible", 0)
				if makePropertyInvisible.Error == nil {
					response.SuccessResponse(http.StatusOK, "property visibility updated successfully", nil, c)
					return
				} else {
					response.ErrorResponse(http.StatusInternalServerError, "An error occurred..", c)
					return
				}

			} else {
				if errors.Is(checkPropertyIsVisible.Error, gorm.ErrRecordNotFound) {
					makePropertyVisible := p.Db.Model(&models.Property{}).Where("id = ?", propertyVisibility.PropertyID).Update("is_visible", 1)
					if makePropertyVisible.Error == nil {
						response.SuccessResponse(http.StatusOK, "property visibility updated successfully", nil, c)
						return
					} else {
						response.ErrorResponse(http.StatusInternalServerError, "An error occurred.", c)
						return
					}
				} else {
					response.ErrorResponse(http.StatusInternalServerError, "An error occurred!!", c)
					return
				}
			}

		} else {
			if errors.Is(checkPropertyExist.Error, gorm.ErrRecordNotFound) == true {

				response.ErrorResponse(http.StatusNotFound, "property not found", c)
				return

			} else {
				response.ErrorResponse(http.StatusInternalServerError, "An error occurred!", c)
				return
			}

		}

	} else {
		response.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}

}

func (p PropertyController) AddOrRemoveFavoriteProperty(c *gin.Context) {
	userId, _ := c.Get("user")
	var addFavorite request2.FavoritePropertyInput
	var favoriteProperty models.FavoriteProperty
	deletedAt := time.Now()
	if err := c.ShouldBind(&addFavorite); err == nil {
		// check if already property is already added
		findProperty := p.Db.Model(&models.FavoriteProperty{}).Where("property_id = ?", addFavorite.PropertyID).Where("user_id = ?", userId).Where("deleted_at IS NULL").First(&favoriteProperty)
		// if no error and property found in favorite, remove property from favorite
		if findProperty.Error == nil {
			// remove property
			result := p.Db.Model(&models.FavoriteProperty{}).Where("property_id = ?", addFavorite.PropertyID).Where("user_id = ?", userId).Update("deleted_at", deletedAt)

			if result.Error == nil {
				response.SuccessResponse(http.StatusOK, "property removed from favorite", nil, c)
				return
			} else {
				response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}
		} else {
			// check if property already exist update its deleted status to null
			if errors.Is(findProperty.Error, gorm.ErrRecordNotFound) {
				checkProperty := p.Db.Model(&models.FavoriteProperty{}).Where("property_id = ?", addFavorite.PropertyID).Where("user_id = ?", userId).Not("deleted_at", nil).First(&favoriteProperty)
				if checkProperty.Error == nil {
					UpdateResult := p.Db.Model(&models.FavoriteProperty{}).Where("property_id = ?", addFavorite.PropertyID).Where("user_id = ?", userId).Updates(map[string]interface{}{"deleted_at": nil})
					if UpdateResult.Error == nil {
						response.SuccessResponse(http.StatusOK, "property added to favorite", nil, c)
						return
					} else {
						response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
						return
					}
				} else {
					// else if error is equals to not found, add property to favorite
					if errors.Is(checkProperty.Error, gorm.ErrRecordNotFound) == true {
						addProperty := map[string]interface{}{
							"property_id": addFavorite.PropertyID,
							"user_id":     userId,
						}
						addPropertyResult := p.Db.Model(&models.FavoriteProperty{}).Create(addProperty)
						if addPropertyResult.Error == nil {
							response.SuccessResponse(http.StatusOK, "property added to favorite", nil, c)
							return

						} else {
							response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
							return
						}

					} else {
						response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
						return
					}
				}

			} else {
				response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
				return
			}

		}

	} else {
		response.ErrorResponse(http.StatusBadRequest, err.Error(), c)
		return
	}

}

func (p PropertyController) GetUserFavoriteProperties(c *gin.Context) {
	userId, _ := c.Get("user")
	var properties []*models.Property
	results := p.Db.Model(&models.Property{}).Joins("left join favorite_properties on favorite_properties.property_id = properties.id").Where("favorite_properties.user_id = ?", userId).Where("favorite_properties.deleted_at IS NULL").Find(&properties)

	if results.Error == nil {
		response.SuccessResponse(http.StatusOK, "properties fetches successfully", properties, c)
		return
	} else {
		if errors.Is(results.Error, gorm.ErrRecordNotFound) == true {
			response.ErrorResponse(http.StatusNotFound, "no properties found", c)
			return

		} else {
			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
			return
		}

	}
}

func (p PropertyController) GetUserUploadedProperties(c *gin.Context) {
	userId, _ := c.Get("user")

	page := c.DefaultQuery("page", "1")

	limit := c.DefaultQuery("limit", "10")
	//var user models.User
	//findUser := p.Db.Model(&models.User{}).Where("id = ?", userId).First(&user)
	preloadQuery := "PropertyImages"

	fetchProperties, err := common.Pagination("properties", page, limit, "created_at desc", map[string]interface{}{"user_id": userId}, &preloadQuery, p.Db)

	if err == nil {
		response.SuccessResponse(http.StatusOK, "properties fetched successfully", fetchProperties, c)
		return
	} else {
		if errors.Is(err, gorm.ErrRecordNotFound) == true {
			response.ErrorResponse(http.StatusNotFound, "no properties found", c)
			return
		} else {
			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
			return
		}
	}

}

func (p PropertyController) PropertyRecommendations(c *gin.Context) {
	userId, _ := c.Get("user")
	var user models.User
	var properties []models.Property

	if err := p.Db.Model(&models.Property{}).Find(&properties).Limit(300).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) == true {
			response.ErrorResponse(http.StatusNotFound, "no properties found", c)
			return
		} else {
			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
			return
		}
	}

	if err := p.Db.Model(&models.User{}).Where("id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) == true {
			response.ErrorResponse(http.StatusNotFound, "user not found", c)
			return
		} else {
			response.ErrorResponse(http.StatusInternalServerError, "An error occurred", c)
			return
		}
	}
	recommendation := Recommend(properties, user)

	response.SuccessResponse(http.StatusOK, "recommendation fetched successfully", recommendation, c)

}

func Recommend(destinations []models.Property, currDest models.User) []models.Property {
	var similarPlaces []models.Property

	for _, d := range destinations {
		sim := ComputeSimilarity(d, currDest)

		fmt.Printf("Distance to %v: %f\n", d.Address, sim)

		if sim > 0.4 {
			similarPlaces = append(similarPlaces, d)
		}

	}

	return similarPlaces

}

func ComputeSimilarity(item1 models.Property, item2 models.User) float64 {
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
