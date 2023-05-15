package repository

import (
	"test/model"
	"test/repository/sqlite"
)

type IUserRepository interface {
	Create(model.User) error
	GetByLogin(login string) (model.User, error)
}

type ISessionRepository interface {
	Create(s model.Session) error
	GetByCookie(cookie string) (model.Session, error)
	Delete(uid int) error
}

type IPostRepository interface {
	Create(model.PostCreateDTO) error
	Get(id int) (model.Post, error)
	List() ([]model.Post, error)
	GetByUser(uid int) ([]model.Post, error)
}

type ICommentRepository interface {
	Create(model.CommentCreateDTO) error
	GetByPost(postId int) ([]model.Comment, error)
}

type IVoteRepository interface {
	Create(model.Vote) error
	Get(model.Vote) (bool, error)
	Delete(model.Vote) error
	GetByUser(model.Vote) ([]model.Vote, error)
}

type IVoteCommentRepository interface {
	Create(model.Vote) error
	Get(model.Vote) (bool, error)
	Delete(model.Vote) error
	GetByUser(model.Vote) ([]model.Vote, error)
}

type Manager struct {
	User IUserRepository
	Sess ISessionRepository
	Comm ICommentRepository
	Post IPostRepository
	Vote IVoteRepository
	Voco IVoteCommentRepository
}

func NewManagerRepository() *Manager {
	db := sqlite.NewDBSqlite()

	userRepo := sqlite.NewUserRepository(db)
	sessRepo := sqlite.NewSessionRepository(db)
	voteRepo := sqlite.NewVoteRepository(db)
	commRepo := sqlite.NewCommetsRepository(db)
	postRepo := sqlite.NewPostRepository(db)
	vocoRepo := sqlite.NewVoteCommentRepository(db)

	return &Manager{
		User: userRepo,
		Sess: sessRepo,
		Post: postRepo,
		Comm: commRepo,
		Vote: voteRepo,
		Voco: vocoRepo,
	}
}
