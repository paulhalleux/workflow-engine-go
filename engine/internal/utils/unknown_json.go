package utils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type UnknownJson map[string]interface{} // @name UnknownJson

func (a *UnknownJson) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *UnknownJson) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

func UnknownJsonFromMap(m *map[string]interface{}) *UnknownJson {
	if m == nil {
		return nil
	}
	uj := UnknownJson(*m)
	return &uj
}

func (a *UnknownJson) ToMap() *map[string]interface{} {
	if a == nil {
		return nil
	}
	m := map[string]interface{}(*a)
	return &m
}
