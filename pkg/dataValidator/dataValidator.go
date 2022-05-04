package dataValidator

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/app/models"
	"gopkg.in/validator.v2"
	"math"
	"reflect"
	"regexp"
	"time"
)

func SetValidators() {
	validator.SetValidationFunc("birthday", func(val interface{}, _ string) error {
		v := reflect.ValueOf(val)
		if v.Kind() != reflect.Struct {
			return validator.ErrUnsupported
		}
		if v.IsZero() {
			return nil
		}
		nVal := val.(time.Time)

		age := int(math.Floor(time.Now().Sub(nVal).Hours() / 24 / 365))
		if age < models.BirthdayBottomLimit || age > models.BirthdayTopLimit {
			return handlers.ErrValidateProfile
		}

		return nil
	})
	validator.SetValidationFunc("password", func(val interface{}, _ string) error {
		v := reflect.ValueOf(val)
		if v.Kind() != reflect.String {
			return validator.ErrUnsupported
		}
		if v.IsZero() {
			return nil
		}
		nVal := val.(string)

		if len(nVal) > models.PasswordMaxLength || len(nVal) < models.PasswordMinLength {
			return validator.ErrLen
		}
		match, err := regexp.MatchString(models.PasswordPatternLowerCase, nVal)
		if err != nil {
			return handlers.ErrBaseApp
		}
		if !match {
			return validator.ErrRegexp
		}
		match, err = regexp.MatchString(models.PasswordPatternUpperCase, nVal)
		if err != nil {
			return handlers.ErrBaseApp
		}
		if !match {
			return validator.ErrRegexp
		}
		match, err = regexp.MatchString(models.PasswordPatternNumber, nVal)
		if err != nil {
			return handlers.ErrBaseApp
		}
		if !match {
			return validator.ErrRegexp
		}

		return nil
	})
}
