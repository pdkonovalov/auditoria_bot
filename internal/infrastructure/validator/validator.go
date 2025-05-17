package validator

import (
	playground "github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *playground.Validate
}

func New() (*Validator, error) {
	return &Validator{
		validate: playground.New(),
	}, nil
}

func (v *Validator) ContactInfo(value string) bool {
	return v.Phone(value) || v.Email(value) || v.Url(value)
}

func (v *Validator) Phone(value string) bool {
	err := v.validate.Var(value, "e164")
	if err != nil {
		return false
	}
	return true
}

func (v *Validator) Email(value string) bool {
	err := v.validate.Var(value, "email")
	if err != nil {
		return false
	}
	return true
}

func (v *Validator) Url(value string) bool {
	err := v.validate.Var(value, "url")
	if err != nil {
		return false
	}
	return true
}
