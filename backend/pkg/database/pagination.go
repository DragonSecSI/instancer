package database

import (
	"gorm.io/gorm"
)

func Paginate(db *gorm.DB, page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page < 1 {
			page = 1
		}
		if pageSize < 1 {
			pageSize = 50
		}

		offset := (page - 1) * pageSize
		limit := pageSize

		return db.Offset(offset).Limit(limit)
	}
}
