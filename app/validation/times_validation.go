package validation

import (
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
)

func RegisterTimesValidation(validate *validator.Validate) {
	validate.RegisterValidation("isbirthdate", func(fl validator.FieldLevel) bool {
		field := fl.Field()

		if field.Kind() == reflect.Ptr && field.IsNil() {
			return false // nil will not pass in this validation
			// if you want to pass the nil, required to use omitnil validation at the front
		}

		var birthDate time.Time
		if field.Kind() == reflect.Ptr {
			birthDate = field.Elem().Interface().(time.Time)
		} else {
			birthDate = field.Interface().(time.Time)
		}

		now := time.Now()
		age := now.Year() - birthDate.Year()

		if now.YearDay() < birthDate.YearDay() {
			age--
		}

		return age > 12 && age <= 120
	})
	validate.RegisterValidation("notfuture", func(fl validator.FieldLevel) bool {
		field := fl.Field()

		if field.Kind() == reflect.Ptr && field.IsNil() {
			return false // nil will not pass in this validation
			// if you want to pass the nil, required to use omitnil validation at the front
		}

		var date time.Time
		if field.Kind() == reflect.Ptr {
			date = field.Elem().Interface().(time.Time)
		} else {
			date = field.Interface().(time.Time)
		}

		return date.Before(time.Now()) || date.Equal(time.Now())
	})
	validate.RegisterValidation("notpast", func(fl validator.FieldLevel) bool {
		field := fl.Field()

		if field.Kind() == reflect.Ptr && field.IsNil() {
			return false
		}

		var date time.Time
		if field.Kind() == reflect.Ptr {
			date = field.Elem().Interface().(time.Time)
		} else {
			date = field.Interface().(time.Time)
		}

		return date.After(time.Now()) || date.Equal(time.Now())
	})
}
