package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"city-pulse/app/daos"
	"city-pulse/app/models"
	"city-pulse/internal/config"
)

var (
	ErrEmailExists  = errors.New("bu email zaten kayıtlı")
	ErrInvalidCreds = errors.New("email veya şifre hatalı")
)

type UserService struct {
	dao *daos.UserDAO
}

func NewUserService(dao *daos.UserDAO) *UserService {
	return &UserService{dao: dao}
}

func (s *UserService) Register(ctx context.Context, email, username, password string) (*models.User, error) {
	existing, err := s.dao.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		Username:     username,
		PasswordHash: string(hash),
	}

	if err := s.dao.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, *models.User, error) {
	user, err := s.dao.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, ErrInvalidCreds
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, ErrInvalidCreds
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      user.ID.Hex(),
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(72 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	})

	tokenStr, err := token.SignedString([]byte(config.Config.JWTSecret))
	if err != nil {
		return "", nil, err
	}

	return tokenStr, user, nil
}
