package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"test/model"
)

type sessionRepository struct {
	db *sql.DB
}

const sessionRepoPath = `sessionRepository: %w`

func NewSessionRepository(db *sql.DB) *sessionRepository {
	return &sessionRepository{
		db: db,
	}
}

func (r *sessionRepository) Create(s model.Session) error {
	query := `insert into session (user_id, cookie, expire_at) values ($1, $2, $3)`

	if _, err := r.db.Exec(query, &s.ID, &s.Cookie, &s.ExpireAt); err != nil {
		return fmt.Errorf(sessionRepoPath, err)
	}

	return nil
}

func (r *sessionRepository) GetByCookie(cookie string) (model.Session, error) {
	query := `SELECT *FROM session WHERE cookie = $1`

	var session model.Session
	err := r.db.QueryRow(query, &cookie).Scan(&session.ID, &session.Cookie, &session.ExpireAt)
	if errors.Is(err, sql.ErrNoRows) {
		return session, fmt.Errorf(sessionRepoPath, model.ErrNoCookie)
	} else if err != nil {
		return session, fmt.Errorf(sessionRepoPath, err)
	}
	return session, nil
}

func (r *sessionRepository) Delete(uid int) error {
	query := `DELETE FROM session WHERE user_id = $1`

	_, err := r.db.Exec(query, uid)
	if err != nil {
		return err
	}
	return nil
}
