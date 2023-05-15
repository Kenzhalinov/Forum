package service

import (
	"database/sql"
	"errors"
	"fmt"

	"test/model"
	"test/repository"
)

const voteServicePath = `voteService: %w`

type voteService struct {
	voteRepo repository.IVoteRepository
}

func NewVoteService(voteRepo repository.IVoteRepository) *voteService {
	return &voteService{
		voteRepo: voteRepo,
	}
}

func (s *voteService) Vote(vote model.Vote) error {
	islike, err := s.voteRepo.Get(vote)
	if errors.Is(err, sql.ErrNoRows) {
		if err := s.voteRepo.Create(vote); err != nil {
			return fmt.Errorf(voteServicePath, err)
		}
		return nil

	} else if err != nil {
		fmt.Println(2)
		return fmt.Errorf(voteServicePath, err)
	}

	if err := s.voteRepo.Delete(vote); err != nil {
		return fmt.Errorf(voteServicePath, err)
	}

	if vote.Vote != islike {
		if err := s.voteRepo.Create(vote); err != nil {
			return fmt.Errorf(voteServicePath, err)
		}
		return nil
	}

	return nil
}

const voteCommentServicePath = `voteCommentService: %w`

type voteCommentService struct {
	voteCommentRepo repository.IVoteCommentRepository
}

func NewVoteCommentService(voteCommentRepo repository.IVoteCommentRepository) *voteCommentService {
	return &voteCommentService{
		voteCommentRepo: voteCommentRepo,
	}
}

func (s *voteCommentService) Vote(vote model.Vote) error {
	islike, err := s.voteCommentRepo.Get(vote)
	if errors.Is(err, sql.ErrNoRows) {
		if err := s.voteCommentRepo.Create(vote); err != nil {
			return fmt.Errorf(voteServicePath, err)
		}
		return nil

	} else if err != nil {
		fmt.Println(2)
		return fmt.Errorf(voteServicePath, err)
	}

	if err := s.voteCommentRepo.Delete(vote); err != nil {
		return fmt.Errorf(voteServicePath, err)
	}

	if vote.Vote != islike {
		if err := s.voteCommentRepo.Create(vote); err != nil {
			return fmt.Errorf(voteServicePath, err)
		}
		return nil
	}

	return nil
}
