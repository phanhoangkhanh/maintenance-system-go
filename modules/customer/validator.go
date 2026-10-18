package customer

import (
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type ValidatorCustomer struct {
	Repo *Repository
}

// RegisterValidatorCustomer must run once at startup, before binding requests.
func RegisterValidatorCustomer(repo *Repository) error {
	if repo == nil || repo.db == nil {
		return fmt.Errorf("customer validator requires a repository with a database")
	}
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("unsupported customer validation engine")
	}
	
	customerValidator := &ValidatorCustomer{Repo: repo}
	v.RegisterValidation("uniqueName", customerValidator.UniqueName)
	v.RegisterValidation("uniqueAddress", customerValidator.UniqueAddress)
	return nil
}

func (v *ValidatorCustomer) UniqueName(fl validator.FieldLevel) bool {
	name, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	exists, err := v.Repo.ExistsByName(name)
	return err == nil && !exists
}

func (v *ValidatorCustomer) UniqueAddress(fl validator.FieldLevel) bool {
	address, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	exists, err := v.Repo.ExistsByAddress(address)
	return err == nil && !exists
}
