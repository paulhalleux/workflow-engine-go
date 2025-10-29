package utils

import (
	"encoding/json"

	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/datatypes"
)

func StructToJSON(s *structpb.Struct) (datatypes.JSON, error) {
	if s == nil {
		return datatypes.JSON([]byte("{}")), nil
	}

	b, err := json.Marshal(s.AsMap())
	if err != nil {
		return nil, err
	}

	return datatypes.JSON(b), nil
}
