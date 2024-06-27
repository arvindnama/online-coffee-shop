package dataUtils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Field     string
	Validator validator.Func
}

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

func NewValidation(customValidators []*CustomValidator) *Validation {
	validate := validator.New()
	for _, c := range customValidators {
		validate.RegisterValidation(c.Field, c.Validator)
	}
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
