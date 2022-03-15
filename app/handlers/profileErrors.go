package handlers

import "errors"

var (
	ErrProfileNotFound = errors.New("profile not found")
	ErrProfileNotFiled = errors.New("profile not filed")
	ErrMockIsEmpty     = errors.New("profile not filed")
)
