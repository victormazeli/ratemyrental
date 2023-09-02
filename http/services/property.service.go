package services

import (
	"gorm.io/gorm"
	"math"
	"rateMyRentalBackend/config"
	"rateMyRentalBackend/database/models"
	"rateMyRentalBackend/http/response"
	"strconv"
)

type PropertyService struct {
	Db  *gorm.DB
	Env *config.Env
}

func (ps PropertyService) Pagination(page string, limit string, query map[string]interface{}, searchTerm string) (*response.PaginationDTO, error) {
	parsePage, er := strconv.ParseInt(page, 0, 64)
	if er != nil {
		return nil, er
	}

	parseLimit, e := strconv.ParseInt(limit, 0, 64)
	if e != nil {
		return nil, e
	}

	var totalRows int64
	var data []models.Property

	queryDB := ps.Db.Preload("PropertyImages").Where(query)

	// Apply search term to multiple columns
	if searchTerm != "" {
		searchTerm = "%" + searchTerm + "%"
		queryDB = queryDB.Where("property_title LIKE ? OR longitude LIKE ? OR latitude LIKE ? OR state LIKE ? OR city LIKE ? OR postal_code LIKE ? OR country LIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm, searchTerm)
	}

	// Count the total rows after applying the search term
	queryDB.Model(&models.Property{}).Count(&totalRows)

	totalPages := int(math.Ceil(float64(totalRows)) / float64(parseLimit))

	offset := (parsePage - 1) * parseLimit
	sort := "created_at desc"
	queryDB = queryDB.Offset(int(offset)).Limit(int(parseLimit)).Order(sort)

	if result := queryDB.Find(&data).Error; result != nil {
		return nil, result
	}

	paginatedResult := &response.PaginationDTO{
		Docs:       data,
		TotalPages: int64(totalPages),
		TotalDocs:  totalRows,
		Page:       parsePage,
		Limit:      parseLimit,
	}

	return paginatedResult, nil
}
