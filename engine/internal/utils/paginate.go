package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GormScopeFactory func(db *gorm.DB) *gorm.DB

func Paginate(ctx *gin.Context) GormScopeFactory {
	return func(db *gorm.DB) *gorm.DB {
		limitParam := ctx.DefaultQuery("limit", "10")
		offsetParam := ctx.DefaultQuery("offset", "0")

		limit, err := strconv.Atoi(limitParam)
		if err != nil || limit <= 0 {
			limit = 10
		}

		offset, err := strconv.Atoi(offsetParam)
		if err != nil || offset < 0 {
			offset = 0
		}

		return db.Limit(limit).Offset(offset)
	}
}

func WithScope(db *gorm.DB, scopeFactory *GormScopeFactory) *gorm.DB {
	if scopeFactory != nil {
		return db.Scopes(*scopeFactory)
	}
	return db
}
