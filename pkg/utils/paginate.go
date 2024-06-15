package utils

import (
	"math"

	"github.com/reeversedev2/zalanda-warehouse-service/pkg/pagination"
	"gorm.io/gorm"
)

func Paginate(value interface{}, pagination *pagination.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = int64(totalPages)

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}
