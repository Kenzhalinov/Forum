package model

import "errors"

type User struct {
	ID       int
	Email    string
	Login    string
	Password string
}

var (
	ErrIncorectData         = errors.New("data is not valid")
	ErrUserNotExist         = errors.New("user is not exists")
	ErrUserIncorrectPasword = errors.New("data is not valid")
)
