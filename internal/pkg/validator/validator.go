package validator

import (
	"github.com/go-playground/validator/v10"
	"sync"
)

var (
	once     sync.Once
	validate *validator.Validate
)

func Validate() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
	})
	return validate
}
