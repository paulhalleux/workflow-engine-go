package expr

import "gorm.io/gorm"

// CompareExpression represents a comparison between two value expressions using a comparison operator.
type CompareExpression struct {
	Left     *string            `json:"left,omitempty"`
	Operator ComparisonOperator `json:"operator,omitempty"`
	Right    *string            `json:"right,omitempty"`
}

type ComparisonOperator string

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
	left := *e.Left
	right := *e.Right
	switch e.Operator {
	case OperatorEquals:
		if negate {
			db = db.Where(left+" != ?", right)
		} else {
			db = db.Where(left+" = ?", right)
		}
	case OperatorNotEquals:
		if negate {
			db = db.Where(left+" = ?", right)
		} else {
			db = db.Where(left+" != ?", right)
		}
	case OperatorGreaterThan:
		if negate {
			db = db.Where(left+" <= ?", right)
		} else {
			db = db.Where(left+" > ?", right)
		}
	case OperatorLessThan:
		if negate {
			db = db.Where(left+" >= ?", right)
		} else {
			db = db.Where(left+" < ?", right)
		}
	case OperatorGreaterEqual:
		if negate {
			db = db.Where(left+" < ?", right)
		} else {
			db = db.Where(left+" >= ?", right)
		}
	case OperatorLessEqual:
		if negate {
			db = db.Where(left+" > ?", right)
		} else {
			db = db.Where(left+" <= ?", right)
		}
	case OperatorIn:
		if negate {
			db = db.Where(left+" NOT IN (?)", right)
		} else {
			db = db.Where(left+" IN (?)", right)
		}
	}
	return db
}
