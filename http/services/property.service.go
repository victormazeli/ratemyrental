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

func (ps PropertyService) Pagination(page string, limit string, query map[string]interface{}) (*response.PaginationDTO, error) {
	parsePage, er := strconv.ParseInt(page, 0, 64)
	if er != nil {
		return nil, er
	}

	// parse query "limit" to int
	parseLimit, e := strconv.ParseInt(limit, 0, 64)
	if e != nil {
		return nil, e
	}

	var totalRows int64

	// Explicitly set the table name for the model before using it in the query
	ps.Db.Model(&models.Property{}).Where(query).Count(&totalRows)

	totalPages := int(math.Ceil(float64(totalRows)) / float64(parseLimit))

	var data []models.Property
	offset := (parsePage - 1) * parseLimit
	preloadAssociation := "PropertyImages"
	sort := "created_at desc"
	queryDB := ps.Db.Preload(preloadAssociation).Offset(int(offset)).Limit(int(parseLimit)).Order(sort)

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
