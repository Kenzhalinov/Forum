package model

type CommentCreateDTO struct {
	PostID  int
	UserID  int
	Content string
}
type Comment struct {
	ID       int
	PostID   int
	User     string
	Content  string
	Likes    int
	Dislikes int
}
