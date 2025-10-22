package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
}

func CreateValidator() (*validator.Validate, error) {
	var validate = validator.New()

	if err := validate.RegisterValidation("wf_version", validateWfVersion); err != nil {
		return nil, err
	}

	return validate, nil
}

func Validate(v interface{}) ([]*ValidationError, error) {
	validate, err := CreateValidator()
	if err != nil {
		return nil, errors.New("failed to create validator")
	}

	enl := en.New()
	var uni = ut.New(enl, enl)
	trans, _ := uni.GetTranslator("en")
	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		return nil, err
	}

	err = validate.Struct(v)
	if err != nil {
		var validationErrors []*ValidationError
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, &ValidationError{
				Field:   err.Field(),
				Tag:     err.Tag(),
				Value:   fmt.Sprintf("%v", err.Value()),
				Message: err.Translate(trans),
			})
		}
		return validationErrors, nil
	}

	return nil, nil
}

func validateWfVersion(fl validator.FieldLevel) bool {
	version := fl.Field().String()
	// Simple semantic versioning check: MAJOR.MINOR.PATCH(-draft.X)?
	var major, minor, patch, draft int
	n, err := fmt.Sscanf(version, "%d.%d.%d-draft.%d", &major, &minor, &patch, &draft)
	if err == nil && n == 4 {
		return true
	}
	n, err = fmt.Sscanf(version, "%d.%d.%d", &major, &minor, &patch)
	if err == nil && n == 3 {
		return true
	}
	return false
}
