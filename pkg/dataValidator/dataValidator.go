package dataValidator

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/app/models"
	"gopkg.in/validator.v2"
	"reflect"
	"regexp"
)

func SetValidators() {
	validator.SetValidationFunc("interests", func(val interface{}, _ string) error {
		v := reflect.ValueOf(val)
		if v.Kind() != reflect.Slice {
			return validator.ErrUnsupported
		}
		if v.IsNil() {
			return nil
		}
		nVal := val.([]string)
		for _, value := range nVal {
			if len(value) > models.InterestSize {
				return validator.ErrLen
			}
		}
		return nil
	})
	validator.SetValidationFunc("birthday", func(val interface{}, _ string) error {
		v := reflect.ValueOf(val)
		if v.Kind() != reflect.String {
			return validator.ErrUnsupported
		}
		if v.IsZero() {
			return nil
		}
		nVal := val.(string)
		if len(nVal) > models.BirthdaySize {
			return validator.ErrLen
		}
		check, err := regexp.MatchString(models.BirthdayRexexp, nVal)
		if err != nil {
			return handlers.ErrBaseApp
		}
		if !check {
			return validator.ErrRegexp
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
