package model

import (
	"errors"
	"time"
)

type Session struct {
	ID       int
	Cookie   string
	ExpireAt time.Time
}

var (
	ErrSessionIsExpired = errors.New("cookie time is expired")
	ErrNoCookie         = errors.New("cookie is not found")
)
