package service

import (
	"regexp"
	"strings"
)

type validator struct{}

func NewValidator() *validator {
	return &validator{}
}

func (v *validator) StringValidate(s string) bool {
	if len(strings.TrimSpace(s)) != 0 {
		return true
	} else {
		return false
	}
}

func (v *validator) EmailValidate(s string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(s)
}
