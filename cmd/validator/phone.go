package validator

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func RegisterPhoneValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
			return regexp.MustCompile(`^1[3-9]\d{9}$`).MatchString(fl.Field().String())
		})
	}
}
