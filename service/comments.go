package service

import (
	"fmt"

	"test/model"
	"test/repository"
)

const commentsServicePath = `commentsService: %w`

type commentsService struct {
	comRepo   repository.ICommentRepository
	validator *validator
}

func NewCommentsService(comRepo repository.ICommentRepository) *commentsService {
	return &commentsService{
		comRepo:   comRepo,
		validator: NewValidator(),
	}
}

func (s *commentsService) Create(comm model.CommentCreateDTO) error {
	if !s.validator.StringValidate(comm.Content) {
		return model.ErrIncorectData
	}

	err := s.comRepo.Create(comm)
	if err != nil {
		return err
	}
	return nil
}

func (s *commentsService) GetByPost(postID int) ([]model.Comment, error) {
	comments, err := s.comRepo.GetByPost(postID)
	if err != nil {
		return nil, fmt.Errorf(commentsServicePath, err)
	}
	return comments, nil
}
