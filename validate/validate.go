package validate

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	SwitchFunc validator.Func = func(fl validator.FieldLevel) bool {
		if fl.Top().Elem().FieldByName("Switch").Interface().(bool) {
			return true
		}
		return validator.New().Var(fl.Field().Interface(), "required") == nil
	}
)

// Register Custom Validate
func RegisterCustomValidate() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("switch", SwitchFunc)
	}
}

// determines if the given string is a valid json
//	@param `val` interface{}
//	@return bool
func IsJson(val interface{}) bool {
	return validator.New().Var(val, "json") == nil
}

// determines if the given string is a valid UUID
//	@param `val` interface{}
//	@return bool
func IsUuid(val interface{}) bool {
	return validator.New().Var(val, "uuid") == nil
}
