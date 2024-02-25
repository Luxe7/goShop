package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	ok, _ := regexp.MatchString(`^1(3[0-9]|4[57]|5[0-35-9]|7[0678]|8[0-9])\d{8}$`, mobile)
	return ok
}
