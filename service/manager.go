package service

import (
	"test/model"
	"test/repository"
)

type IUserService interface {
	Authenticate(login, password string) (model.User, error)
	Authorizate(cookie string) (userID int, err error)
	CreateSession(user model.User) (string, error)
	DeleteSession(uid int) error
	Create(model.User) error
}

type IPostService interface {
	Create(model.PostCreateDTO) error
	Get(id int) (model.Post, error)
	List(string) ([]model.Post, error)
	ListCreatedByUser(uid int) ([]model.Post, error)
	ListLikedByUser(uid int) ([]model.Post, error)
}

type ICommentService interface {
	Create(comm model.CommentCreateDTO) error
	GetByPost(postID int) ([]model.Comment, error)
}

type IVoteService interface {
	Vote(vote model.Vote) error
}

type IVoteCommentService interface {
	Vote(vote model.Vote) error
}

type Manager struct {
	User IUserService
	Post IPostService
	Vote IVoteService
	Comm ICommentService
	Voco IVoteCommentService
}

func NewManagerService(repo *repository.Manager) *Manager {
	return &Manager{
		User: NewUserService(repo.User, repo.Sess),
		Post: NewPostService(repo.Post, repo.Vote),
		Comm: NewCommentsService(repo.Comm),
		Vote: NewVoteService(repo.Vote),
		Voco: NewVoteCommentService(repo.Voco),
	}
}
