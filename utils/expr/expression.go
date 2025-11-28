package expr

import (
	"fmt"

	"gorm.io/gorm"
)

type Expression struct {
	And     *[]Expression      `json:"and,omitempty"`
	Or      *[]Expression      `json:"or,omitempty"`
	Not     *Expression        `json:"not,omitempty"`
	Compare *CompareExpression `json:"compare,omitempty"`
} // @name Expression

func (e Expression) IsEmpty() bool {
	return e.And == nil && e.Or == nil && e.Not == nil && e.Compare == nil
}

func NewEmptyExpression() Expression {
	return Expression{}
}

func NewAndExpression(expressions ...Expression) Expression {
	return Expression{
		And: &expressions,
	}
}

func NewOrExpression(expressions ...Expression) Expression {
	return Expression{
		Or: &expressions,
	}
}

func NewNotExpression(expression Expression) Expression {
	return Expression{
		Not: &expression,
	}
}

func NewCompareExpression(field string, operator ComparisonOperator, value fmt.Stringer) Expression {
	return Expression{
		Compare: &CompareExpression{
			Field:    &field,
			Operator: operator,
			Value:    &value,
		},
	}
}

func (e Expression) ToGorm(db *gorm.DB, negate bool) *gorm.DB {
	if e.And != nil {
		for _, expr := range *e.And {
			db = expr.ToGorm(db, negate)
		}
	} else if e.Or != nil {
		db = db.Or(func(tx *gorm.DB) *gorm.DB {
			for _, expr := range *e.Or {
				tx = expr.ToGorm(tx, negate)
			}
			return tx
		})
	} else if e.Not != nil {
		db = e.Not.ToGorm(db, !negate)
	} else if e.Compare != nil {
		db = e.Compare.ToGorm(db, negate)
	}
	return db
}
