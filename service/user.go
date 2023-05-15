package service

import (
	"fmt"
	"time"

	"test/model"
	"test/repository"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo  repository.IUserRepository
	sessRepo  repository.ISessionRepository
	validator *validator
}

func NewUserService(userRepo repository.IUserRepository, sessRepo repository.ISessionRepository) *userService {
	return &userService{
		userRepo:  userRepo,
		sessRepo:  sessRepo,
		validator: NewValidator(),
	}
}

const (
	cookieExpireTime = 15 * time.Minute
	userServicePath  = `userService: %w`
)

func (s *userService) Authenticate(login, password string) (model.User, error) {
	if !s.validator.StringValidate(login) ||
		!s.validator.StringValidate(password) {
		return model.User{}, fmt.Errorf(userServicePath, model.ErrIncorectData)
	}
	dbuser, err := s.userRepo.GetByLogin(login)
	if err != nil {
		return dbuser, fmt.Errorf(userServicePath, err)
	}
	if err = s.checkPasswordHash(password, dbuser.Password); err != nil {
		return dbuser, fmt.Errorf(userServicePath, err)
	}
	return dbuser, nil
}

func (s *userService) Authorizate(cookie string) (userID int, err error) {
	sess, err := s.sessRepo.GetByCookie(cookie)
	if err != nil {
		return 0, fmt.Errorf(userServicePath, err)
	}
	if sess.ExpireAt.Before((time.Now())) {
		return 0, fmt.Errorf(userServicePath, model.ErrSessionIsExpired)
	}
	return sess.ID, nil
}

func (s *userService) CreateSession(user model.User) (string, error) {
	cookie, err := s.generateCookie()
	if err != nil {
		return "", fmt.Errorf(userServicePath, err)
	}

	sess := model.Session{
		ID:       user.ID,
		Cookie:   cookie,
		ExpireAt: time.Now().Add(cookieExpireTime),
	}

	err = s.sessRepo.Delete(sess.ID)
	if err != nil {
		return "", fmt.Errorf(userServicePath, err)
	}

	if err := s.sessRepo.Create(sess); err != nil {
		return "", fmt.Errorf(userServicePath, err)
	}
	return cookie, nil
}

func (s *userService) DeleteSession(uid int) error {
	if err := s.sessRepo.Delete(uid); err != nil {
		return fmt.Errorf(userServicePath, err)
	}

	return nil
}

func (s *userService) Create(user model.User) error {
	if !s.validator.StringValidate(user.Login) ||
		!s.validator.StringValidate(user.Email) ||
		!s.validator.EmailValidate(user.Email) ||
		!s.validator.StringValidate(user.Password) {
		return fmt.Errorf(userServicePath, model.ErrIncorectData)
	}
	hashPass, err := s.hashPassword(user.Password)
	if err != nil {
		return fmt.Errorf(userServicePath, err)
	}

	user.Password = hashPass

	if err := s.userRepo.Create(user); err != nil {
		return fmt.Errorf(userServicePath, err)
	}
	return nil
}

func (s *userService) generateCookie() (string, error) {
	u2, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return u2.String(), nil
}

func (s *userService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *userService) checkPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
