package validator

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	slugRegex     = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
)

func Setup() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		_ = v.RegisterValidation("username", func(fl validator.FieldLevel) bool {
			return usernameRegex.MatchString(fl.Field().String())
		})

		_ = v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
			return slugRegex.MatchString(fl.Field().String())
		})

		_ = v.RegisterValidation("post_type", func(fl validator.FieldLevel) bool {
			val := fl.Field().String()
			return val == "question" || val == "concept" || val == "build-log"
		})

		_ = v.RegisterValidation("vote_type", func(fl validator.FieldLevel) bool {
			val := fl.Field().Int()
			return val == 1 || val == -1
		})
	}
}
