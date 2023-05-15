package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"test/model"
)

type voteRepository struct {
	db *sql.DB
}

const voteRepoPath = `voteRepository: %w`

func NewVoteRepository(db *sql.DB) *voteRepository {
	return &voteRepository{
		db: db,
	}
}

func (r *voteRepository) Create(v model.Vote) error {
	query := `insert into votes (user_id, post_id, vote) values ($1, $2, $3)`
	if _, err := r.db.Exec(query, &v.UserID, &v.PostID, &v.Vote); err != nil {
		return fmt.Errorf(voteRepoPath, err)
	}
	return nil
}

func (r *voteRepository) Get(v model.Vote) (vote bool, err error) {
	query := `select vote from votes where post_id = $1 and user_id = $2`

	if err = r.db.QueryRow(query, &v.PostID, &v.UserID).Scan(&vote); err != nil {
		return false, fmt.Errorf(voteRepoPath, sql.ErrNoRows)
	}

	return vote, nil
}

func (r *voteRepository) Delete(v model.Vote) error {
	query := `delete from votes where post_id = $1 and user_id = $2`

	_, err := r.db.Exec(query, &v.PostID, &v.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (r *voteRepository) GetByUser(v model.Vote) (votes []model.Vote, err error) {
	query := `select * from votes where user_id = ? and vote = ?`

	rows, err := r.db.Query(query, &v.UserID, &v.Vote)
	if errors.Is(err, sql.ErrNoRows) {
		return []model.Vote{}, nil
	} else if err != nil {
		return nil, fmt.Errorf(voteRepoPath, err)
	}
	defer rows.Close()
	for rows.Next() {
		var vote model.Vote
		if err := rows.Scan(
			&vote.PostID,
			&vote.UserID,
			&vote.Vote,
		); err != nil {
			return nil, fmt.Errorf(voteRepoPath, err)
		}
		votes = append(votes, vote)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(voteRepoPath, err)
	}
	return votes, nil
}

type voteCommentRepository struct {
	db *sql.DB
}

const voteCommentRepoPath = `voteCommentRepository: %w`

func NewVoteCommentRepository(db *sql.DB) *voteCommentRepository {
	return &voteCommentRepository{
		db: db,
	}
}

func (r *voteCommentRepository) Create(v model.Vote) error {
	query := `insert into votes_comm (user_id, comm_id, vote) values ($1, $2, $3)`
	if _, err := r.db.Exec(query, &v.UserID, &v.PostID, &v.Vote); err != nil {
		return fmt.Errorf(voteRepoPath, err)
	}
	return nil
}

func (r *voteCommentRepository) Get(v model.Vote) (vote bool, err error) {
	query := `select vote from votes_comm where comm_id = $1 and user_id = $2`

	if err = r.db.QueryRow(query, &v.PostID, &v.UserID).Scan(&vote); err != nil {
		return false, fmt.Errorf(voteRepoPath, sql.ErrNoRows)
	}

	return vote, nil
}

func (r *voteCommentRepository) Delete(v model.Vote) error {
	query := `delete from votes_comm where comm_id = $1 and user_id = $2`

	_, err := r.db.Exec(query, &v.PostID, &v.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (r *voteCommentRepository) GetByUser(v model.Vote) (votes []model.Vote, err error) {
	query := `select * from votes_comm where comm_id = ? and vote = ?`

	rows, err := r.db.Query(query, &v.UserID, &v.Vote)
	if errors.Is(err, sql.ErrNoRows) {
		return []model.Vote{}, nil
	} else if err != nil {
		return nil, fmt.Errorf(voteRepoPath, err)
	}
	defer rows.Close()
	for rows.Next() {
		var vote model.Vote
		if err := rows.Scan(
			&vote.PostID,
			&vote.UserID,
			&vote.Vote,
		); err != nil {
			return nil, fmt.Errorf(voteRepoPath, err)
		}
		votes = append(votes, vote)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(voteRepoPath, err)
	}
	return votes, nil
}
