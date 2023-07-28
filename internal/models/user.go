package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator"
)

type User struct {
	Common

	Name     string `json:"name" pg:",notnull" validate:"required"`
	Email    string `json:"email" pg:",notnull" validate:"required,email"`
	CreateBy string `json:"create_by"`
}

func (u *User) Validate() error {
	if err := validator.New().Struct(u); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			fieldName := validationErrors[0].Field()
			jsonName := u.getFieldJSONName(fieldName)
			tag := validationErrors[0].Tag()
			if jsonName == "email" && tag == "email" {
				return fmt.Errorf("'%s' is not a valid email address", u.Email)

			}
			return fmt.Errorf("%s is %s", jsonName, validationErrors[0].Tag())
		}
	}
	return nil
}

func (u *User) getFieldJSONName(fieldName string) string {
	t := reflect.TypeOf(u)
	if field, found := t.Elem().FieldByName(fieldName); found {
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			return strings.Split(jsonTag, ",")[0]
		}
	}
	return fieldName
}

func (u *User) IsFieldOutputOnly(field string) bool {
	list := [...]string{
		"create_by",
	}

	for _, curr := range list {
		if curr == field {
			return true
		}
	}

	return u.Common.IsFieldOutputOnly(field)
}
