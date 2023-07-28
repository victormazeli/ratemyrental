package common

import (
	"gorm.io/gorm"
	"math"
	"rateMyRentalBackend/database/models"
	"rateMyRentalBackend/http/response"
	"strconv"
)

func Paginate(page int64, limit int64, query map[string]interface{}, sort string) func(db *gorm.DB) *gorm.DB {

	offset := (page - 1) * limit

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(int(offset)).Limit(int(limit)).Where(query).Order(sort)
	}
}

//func Pagination(model interface{}, page string, limit string, sort string, query map[string]interface{}, preloadQuery *string, db *gorm.DB) (*response.PaginationDTO, error) {
//	parsePage, er := strconv.ParseInt(page, 0, 8)
//	if er != nil {
//		log.Fatal(er)
//		return nil, er
//	}
//	// parse query "limit" to int
//	parseLimit, e := strconv.ParseInt(limit, 0, 8)
//	if e != nil {
//		log.Fatal(e)
//		return nil, e
//	}
//	var totalRows int64
//	db.Model(model).Where(query).Count(&totalRows)
//	totalPages := int(math.Ceil(float64(totalRows)) / float64(parseLimit))
//
//	var data []map[string]interface{}
//
//	if preloadQuery != nil {
//		preloadAssociation := *preloadQuery
//		fmt.Println(preloadAssociation)
//		if result := db.Model(model).Preload("PropertyImages").Scopes(Paginate(parsePage, parseLimit, query, sort)).Find(&data).Error; result != nil {
//			return nil, result
//		}
//	} else {
//		if result := db.Model(model).Scopes(Paginate(parsePage, parseLimit, query, sort)).Find(&data).Error; result != nil {
//			return nil, result
//		}
//	}
//
//	paginatedResult := &response.PaginationDTO{
//		Docs:       data,
//		TotalPages: int64(totalPages),
//		TotalDocs:  totalRows,
//		Page:       parsePage,
//		Limit:      parseLimit,
//	}
//
//	return paginatedResult, nil
//
//}

func Pagination(model interface{}, page string, limit string, sort string, query map[string]interface{}, preloadQuery *string, db *gorm.DB) (*response.PaginationDTO, error) {
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
	db.Model(&models.Property{}).Where(query).Count(&totalRows)

	totalPages := int(math.Ceil(float64(totalRows)) / float64(parseLimit))

	var data []models.Property
	offset := (parsePage - 1) * parseLimit
	queryDB := db.Model(&models.Property{}).Offset(int(offset)).Limit(int(parseLimit)).Order(sort)

	if preloadQuery != nil {
		preloadAssociation := *preloadQuery
		queryDB = db.Preload(preloadAssociation).Offset(int(offset)).Limit(int(parseLimit)).Order(sort)
	}

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
