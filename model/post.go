package model

import "errors"

type Post struct {
	ID       int
	User     string
	Title    string
	Content  string
	Category string
	Likes    int
	Dislikes int
}
type PostCreateDTO struct {
	UserID   int
	Title    string
	Content  string
	Category string
}

var ErrPostIsNotFound = errors.New("post is not found")
