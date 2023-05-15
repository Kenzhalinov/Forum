package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"test/model"
)

const commentsRepoPath = `commetsRepo: %w`

type commetsRepository struct {
	db *sql.DB
}

func NewCommetsRepository(db *sql.DB) *commetsRepository {
	return &commetsRepository{db: db}
}

func (r *commetsRepository) Create(com model.CommentCreateDTO) error {
	query := `INSERT INTO comments (post_id,user_id,content) VALUES($1,$2,$3)`
	_, err := r.db.Exec(query, com.PostID, com.UserID, com.Content)
	if err != nil {
		return fmt.Errorf(commentsRepoPath, err)
	}
	return nil
}

func (r *commetsRepository) GetByPost(postID int) ([]model.Comment, error) {
	query := `SELECT t1.id, t1.post_id, t3.login, t1.content,
	CASE WHEN t4.likes IS NULL THEN 0 ELSE t4.likes END,
	CASE WHEN t5.dislikes IS NUll THEN 0 ELSE t5.dislikes END
	FROM comments t1
	left join users t3 on t3.id = t1.user_id
	LEFT JOIN (SELECT comm_id, COUNT(vote) AS likes FROM votes_comm WHERE vote IS TRUE GROUP BY comm_id ) t4 on t4.comm_id = t1.id
	LEFT JOIN (SELECT comm_id, COUNT(vote) AS dislikes FROM votes_comm WHERE vote IS FALSE GROUP BY comm_id ) t5 on t5.comm_id = t1.id
	WHERE post_id = $1`

	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf(commentsRepoPath, err)
	}
	defer rows.Close()
	commets := []model.Comment{}
	for rows.Next() {
		var commet model.Comment
		if err := rows.Scan(&commet.ID, &commet.PostID, &commet.User, &commet.Content, &commet.Likes, &commet.Dislikes); errors.Is(err, sql.ErrNoRows) {
			return []model.Comment{}, nil
		} else if err != nil {
			return nil, fmt.Errorf(commentsRepoPath, err)
		}
		commets = append(commets, commet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(commentsRepoPath, err)
	}
	return commets, nil
}
