package test

import (
	"testing"

	"github.com/0xirvan/url-shortener/validation"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

var Validate = validation.InitValidator()

func TestCreateUserValidationSuccess(t *testing.T) {
	err := Validate.Struct(&validation.CreateUser{
		Name:     "John Doe",
		Email:    "jhon@doe.com",
		Password: "password123",
		Role:     "user",
	})

	if err != nil {
		errs := err.(validator.ValidationErrors)
		t.Errorf("Validation failed for the following fields:")
		for _, fieldError := range errs {
			t.Errorf("Field: %s, Error: %s", fieldError.Field(), fieldError.Translate(validation.Translator))
		}
	}

	assert.NoError(t, err, "Validation should pass for valid input")
}

func TestCreateUserValidationFailure(t *testing.T) {
	err := Validate.Struct(&validation.CreateUser{
		Name:     "Jo",
		Email:    "invalid-email",
		Password: "short",
		Role:     "invalid-role",
	})

	if err != nil {
		errs := err.(validator.ValidationErrors)
		t.Error(validation.TranslateError(errs))
	}

	assert.Error(t, err, "Validation should fail for invalid input")
}
