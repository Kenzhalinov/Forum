package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"test/model"
)

var ErrPostIsNotFound = errors.New("post is not found")

type postRepository struct {
	db *sql.DB
}

const postRepoPath = `postRepository: %w`

func NewPostRepository(db *sql.DB) *postRepository {
	return &postRepository{
		db: db,
	}
}

func (r *postRepository) Create(p model.PostCreateDTO) error {
	query := `INSERT INTO posts(user_id, title, content, category) VALUES($1, $2, $3, $4)`

	if _, err := r.db.Exec(query, &p.UserID, &p.Title, &p.Content, &p.Category); err != nil {
		return fmt.Errorf(postRepoPath, err)
	}
	return nil
}

func (r *postRepository) Get(id int) (post model.Post, err error) {
	query := `SELECT t1.id, t2.login, t1.title, t1.content, t1.category,
	CASE WHEN t4.likes IS NULL THEN 0 ELSE t4.likes END,
	CASE WHEN t5.dislikes IS NUll THEN 0 ELSE t5.dislikes END  
	FROM posts t1
	LEFT JOIN users t2 on t1.user_id = t2.id
	LEFT JOIN (SELECT post_id, COUNT(vote) AS likes FROM votes WHERE vote IS TRUE GROUP BY post_id ) t4 on t1.id = t4.post_id
	LEFT JOIN (SELECT post_id, COUNT(vote) AS dislikes FROM votes WHERE vote IS FALSE GROUP BY post_id ) t5 on t1.id = t5.post_id
	WHERE t1.id = $1`

	err = r.db.QueryRow(query, &id).Scan(&post.ID, &post.User, &post.Title, &post.Content, &post.Category, &post.Likes, &post.Dislikes)
	if errors.Is(err, sql.ErrNoRows) {
		return post, fmt.Errorf(postRepoPath, ErrPostIsNotFound)
	} else if err != nil {
		return post, err
	}
	return post, nil
}

func (r *postRepository) List() (posts []model.Post, err error) {
	query := `SELECT t1.id, t2.login, t1.title, t1.content, t1.category,
	CASE WHEN t4.likes IS NULL THEN 0 ELSE t4.likes END,
	CASE WHEN t5.dislikes IS NULL THEN 0 ELSE t5.dislikes END
	FROM posts t1
	LEFT JOIN users t2 on t1.user_id = t2.id
	LEFT JOIN (SElECT post_id, COUNT(vote) AS likes FROM votes WHERE vote IS TRUE GROUP BY post_id) t4 on t1.id = t4.post_id 
	LEFT JOIN (SElECT post_id, COUNT(vote) AS dislikes FROM votes WHERE vote IS FALSE GROUP BY post_id) t5 on t1.id = t5.post_id`

	rows, err := r.db.Query(query)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf(postRepoPath, ErrPostIsNotFound)
	} else if err != nil {
		return nil, fmt.Errorf(postRepoPath, err)
	}
	defer rows.Close()

	for rows.Next() {
		var post model.Post
		if err = rows.Scan(
			&post.ID,
			&post.User,
			&post.Title,
			&post.Content,
			&post.Category,
			&post.Likes,
			&post.Dislikes,
		); err != nil {
			return nil, fmt.Errorf(postRepoPath, err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(postRepoPath, err)
	}
	return posts, nil
}

func (r *postRepository) GetByUser(uid int) (posts []model.Post, err error) {
	query := `select t1.id, t2.login, t1.title, t1.content, t1.category,
	case when t4.likes is null then 0 else t4.likes end,
	case when t5.dislikes is null then 0 else t5.dislikes end
	from posts t1
	left join users t2 on t1.user_id = t2.id
	left join (select post_id, count(vote) as likes from votes where vote is true group by post_id) t4 on t1.id = t4.post_id
	left join (select post_id, count(vote) as dislikes from votes where vote is false group by post_id) t5 on t1.id = t5.post_id
	where t1.user_id = ?`

	rows, err := r.db.Query(query, uid)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf(postRepoPath, model.ErrPostIsNotFound)
	} else if err != nil {
		return nil, fmt.Errorf(postRepoPath, err)
	}
	defer rows.Close()
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(
			&post.ID,
			&post.User,
			&post.Title,
			&post.Content,
			&post.Category,
			&post.Likes,
			&post.Dislikes,
		); err != nil {
			return nil, fmt.Errorf(postRepoPath, err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(postRepoPath, err)
	}
	return posts, nil
}
