package validator

import (
	"github.com/go-playground/locales/en"
	universal "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/en"
)

// CustomValidator holds the validator instance
type CustomValidator struct {
	validator  *validator.Validate
	translator universal.Translator
}

// NewValidator creates and returns a new CustomValidator
func NewValidator() *CustomValidator {
	en := en.New()
	un := universal.New(en, en)

	vl := validator.New()
	tr, _ := un.GetTranslator("en")

	translations.RegisterDefaultTranslations(vl, tr)

	return &CustomValidator{
		validator:  vl,
		translator: tr,
	}
}

// Validate validates the given interface and returns a map of errors
func (cv *CustomValidator) Validate(i interface{}) []string {
	if err := cv.validator.Struct(i); err != nil {
		errors := []string{}
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Translate(cv.translator))
		}

		return errors
	}

	return nil
}
