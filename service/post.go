package service

import (
	"fmt"
	"strings"

	"test/model"
	"test/repository"
)

const postServicePath = `postService: %w`

type postService struct {
	postRepo repository.IPostRepository
	voteRepo repository.IVoteRepository

	validator *validator
}

func NewPostService(postRepo repository.IPostRepository, voteRepo repository.IVoteRepository) *postService {
	return &postService{
		postRepo:  postRepo,
		voteRepo:  voteRepo,
		validator: NewValidator(),
	}
}

func (s *postService) Create(p model.PostCreateDTO) error {
	if !s.validator.StringValidate(p.Title) ||
		!s.validator.StringValidate(p.Content) {
		fmt.Println("error1")
		return fmt.Errorf(postServicePath, model.ErrIncorectData)
	}
	if err := s.postRepo.Create(p); err != nil {
		fmt.Println("error2")
		return fmt.Errorf(postServicePath, err)
	}
	return nil
}

func (s *postService) Get(id int) (model.Post, error) {
	post, err := s.postRepo.Get(id)
	if err != nil {
		return post, fmt.Errorf(postServicePath, err)
	}
	return post, nil
}

func (s *postService) List(category string) (listWithCAtegory []model.Post, err error) {
	list, err := s.postRepo.List()
	if err != nil {
		return list, fmt.Errorf(postServicePath, err)
	}

	for _, post := range list {
		if strings.Contains(post.Category, category) {
			listWithCAtegory = append(listWithCAtegory, post)
		}
	}

	return listWithCAtegory, nil
}

func (s *postService) ListCreatedByUser(uid int) ([]model.Post, error) {
	list, err := s.postRepo.GetByUser(uid)
	if err != nil {
		return list, fmt.Errorf(postServicePath, err)
	}

	return list, nil
}

func (s *postService) ListLikedByUser(uid int) (likedPosts []model.Post, err error) {
	list, err := s.List("")
	if err != nil {
		return list, err
	}

	listVote, err := s.voteRepo.GetByUser(model.Vote{
		UserID: uid,
		Vote:   true,
	})
	if err != nil {
		return likedPosts, fmt.Errorf(postServicePath, err)
	}

	for _, post := range list {
		for _, vote := range listVote {
			if post.ID == vote.PostID {
				likedPosts = append(likedPosts, post)
			}
		}
	}

	return likedPosts, nil
}
