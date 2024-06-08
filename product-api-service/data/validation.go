package data

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// ValidationError is a Wrapper on top of FieldError
type ValidationError struct {
	// [learning]: This is how we extend a struct with another struct
	// unlike typescript there is no extend key word ,
	// instead first line inside struct is the base struct.
	validator.FieldError
}

type ValidationErrors []ValidationError

func (v ValidationError) Error() string {
	return fmt.Sprintf(
		"key: '%s' Error: Field validation for '%s' failed on the '%s' tag",
		v.Namespace(),
		v.Field(),
		v.Tag(),
	)
}

// convert validation errors to slice of strings
func (v ValidationErrors) Errors() []string {

	errs := []string{}

	for _, err := range v {
		errs = append(errs, err.Error())
	}
	return errs
}

type Validation struct {
	validate *validator.Validate
}

func NewValidation() *Validation {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return &Validation{validate}
}

func (v *Validation) Validate(i interface{}) ValidationErrors {

	err := v.validate.Struct(i)

	if err == nil {
		return nil
	}

	var returnErrs ValidationErrors
	// [learning]: <>.(<>) is the syntax for typecasting
	// here v.validate.Struct(i) is type casted to validator.ValidationErrors
	for _, err := range err.(validator.ValidationErrors) {
		ve := ValidationError{err}
		returnErrs = append(returnErrs, ve)
	}
	return returnErrs
}

func validateSKU(fl validator.FieldLevel) bool {
	// sku format: xxxx-xxxx-xxxx
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	return len(matches) == 1
}
