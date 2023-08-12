package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/pawpaw2022/simplebank/util"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		// Check if the currency is valid
		return util.IsSupportedCurrency(currency)
	}
	return false
}
