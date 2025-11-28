package pagination

import "gorm.io/gorm"

type Pagination struct {
	Offset int `form:"offset,default=0"`
	Limit  int `form:"limit,default=10"`
} // @name Pagination

type PaginatedResult[T any] struct {
	TotalCount int64 `json:"totalCount"`
	Items      []T   `json:"items"`
} // @name PaginatedResult

func (p Pagination) ToGorm(db *gorm.DB) *gorm.DB {
	return db.Offset(p.Offset).Limit(p.Limit)
}
