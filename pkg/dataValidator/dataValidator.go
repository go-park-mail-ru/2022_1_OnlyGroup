package dataValidator

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/app/models"
	"gopkg.in/validator.v2"
	"reflect"
	"regexp"
	"time"
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

		if len(nVal) == 0 {
			return nil
		}

		timeValidate, err := time.Parse("2006-01-02", nVal)
		if err != nil {
			return validator.ErrInvalid
		}

		topLimit, err := time.Parse("2006-01-02", models.BirthdayTopLimit)
		if err != nil {
			return validator.ErrInvalid
		}

		topLimit, err := time.Parse("2006-01-02", models.BirthdayBottomLimit)
		if err != nil {
			return validator.ErrInvalid
		}

		if time.timeValidate

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
