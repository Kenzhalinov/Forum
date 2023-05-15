package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"test/model"
)

type userRepository struct {
	db *sql.DB
}

const userRepoPath = `userRepository: %w`

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(u model.User) error {
	query := `INSERT INTO users(email, login, password) VALUES($1, $2, $3)`
	if _, err := r.db.Exec(query, &u.Email, &u.Login, &u.Password); err != nil {
		return fmt.Errorf(userRepoPath, err)
	}
	return nil
}

func (r *userRepository) GetByLogin(login string) (user model.User, err error) {
	query := `SELECT * FROM users WHERE login = $1`

	err = r.db.QueryRow(query, login).Scan(&user.ID, &user.Email, &user.Login, &user.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return user, fmt.Errorf(userRepoPath, model.ErrUserNotExist)
	} else if err != nil {
		return user, err
	}
	return user, nil
}
