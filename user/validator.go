package user

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

//MOdels has :binding:"required,inRoleGroup" -> clarify func here and register in module-provider
//Docs: https://pkg.go.dev/github.com/go-playground/validator/v10#hdr-Required
func RegisterValidatorUser() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
    v.RegisterValidation("inRoleGroup", InRoleGroup)
  }
}

var InRoleGroup validator.Func = func(fl validator.FieldLevel) bool {
  role, ok := fl.Field().Interface().(string)
  if ok {
    validRoles := []string{"admin", "manager", "driver", "technician", "operator"} 
    for _, r := range validRoles {
      if role == r {
        return true
      }
    }
    return false
  }
  return false
}

