package common

import (
	"gorm.io/gorm"
	"log"
	"math"
	"rateMyRentalBackend/dtos"
	"strconv"
)

func Paginate(page int64, limit int64, query map[string]interface{}, sort string) func(db *gorm.DB) *gorm.DB {

	offset := (page - 1) * limit

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(int(offset)).Limit(int(limit)).Where(query).Order(sort)
	}
}

func Pagination(table string, page string, limit string, sort string, query map[string]interface{}, db *gorm.DB) (*dtos.PaginationDTO, error) {
	parsePage, er := strconv.ParseInt(page, 0, 8)
	if er != nil {
		log.Fatal(er)
		return nil, er
	}
	// parse query "limit" to int
	parseLimit, e := strconv.ParseInt(limit, 0, 8)
	if e != nil {
		log.Fatal(e)
		return nil, e
	}
	var totalRows int64
	db.Table(table).Where(query).Count(&totalRows)
	totalPages := int(math.Ceil(float64(totalRows)) / float64(parseLimit))

	var data []map[string]interface{}

	if result := db.Table(table).Scopes(Paginate(parsePage, parseLimit, query, sort)).Find(&data).Error; result != nil {
		return nil, result
	}

	paginatedResult := &dtos.PaginationDTO{
		TotalPages: int64(totalPages),
		TotalDocs:  totalRows,
		Page:       parsePage,
		Limit:      parseLimit,
		Docs:       data,
	}

	return paginatedResult, nil

}
