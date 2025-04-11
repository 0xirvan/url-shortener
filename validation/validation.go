package validation

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var Translator ut.Translator

func InitValidator() *validator.Validate {
	validate := validator.New()

	// Initialize the translator
	english := en.New()
	uni := ut.New(english, english)
	Translator, _ = uni.GetTranslator("en")

	// Register English translations
	en_translations.RegisterDefaultTranslations(validate, Translator)

	// Register custom validation functions
	validate.RegisterValidation("password", Password)
	validate.RegisterTranslation("password", Translator, func(ut ut.Translator) error {
		return ut.Add("password", "Password must be at least 8 characters long, contain letter amda number", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("password")
		return t
	})

	return validate
}

func TranslateError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			errors[fieldError.Field()] = fieldError.Translate(Translator)
		}
	}

	return errors
}
