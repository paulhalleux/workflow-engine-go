package expr

import (
	"fmt"

	"gorm.io/gorm"
)

// CompareExpression represents a comparison between two value expressions using a comparison operator.
type CompareExpression struct {
	Field    *string            `json:"field,omitempty"`
	Operator ComparisonOperator `json:"operator,omitempty"`
	Value    *fmt.Stringer      `json:"value,omitempty"`
} // @name CompareExpression

type ComparisonOperator string // @name ComparisonOperator

const (
	OperatorEquals       ComparisonOperator = "="
	OperatorNotEquals    ComparisonOperator = "!="
	OperatorGreaterThan  ComparisonOperator = ">"
	OperatorLessThan     ComparisonOperator = "<"
	OperatorGreaterEqual ComparisonOperator = ">="
	OperatorLessEqual    ComparisonOperator = "<="
	OperatorIn           ComparisonOperator = "in"
)

func (e CompareExpression) ToGorm(db *gorm.DB, negate bool) *gorm.DB {
	field := *e.Field
	value := (*e.Value).String()
	switch e.Operator {
	case OperatorEquals:
		if negate {
			db = db.Where(field+" != ?", value)
		} else {
			db = db.Where(field+" = ?", value)
		}
	case OperatorNotEquals:
		if negate {
			db = db.Where(field+" = ?", value)
		} else {
			db = db.Where(field+" != ?", value)
		}
	case OperatorGreaterThan:
		if negate {
			db = db.Where(field+" <= ?", value)
		} else {
			db = db.Where(field+" > ?", value)
		}
	case OperatorLessThan:
		if negate {
			db = db.Where(field+" >= ?", value)
		} else {
			db = db.Where(field+" < ?", value)
		}
	case OperatorGreaterEqual:
		if negate {
			db = db.Where(field+" < ?", value)
		} else {
			db = db.Where(field+" >= ?", value)
		}
	case OperatorLessEqual:
		if negate {
			db = db.Where(field+" > ?", value)
		} else {
			db = db.Where(field+" <= ?", value)
		}
	case OperatorIn:
		if negate {
			db = db.Where(field+" NOT IN (?)", value)
		} else {
			db = db.Where(field+" IN (?)", value)
		}
	}
	return db
}
