package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/steven230500/cartagena-api/internal/domain"
	"github.com/steven230500/cartagena-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
	ErrEmailTaken         = errors.New("email already registered")
)

const tokenTTL = 7 * 24 * time.Hour

type AuthService struct {
	users     repository.UserRepository
	jwtSecret []byte
}

func NewAuthService(users repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{users: users, jwtSecret: []byte(jwtSecret)}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (domain.User, error) {
	if len(password) < 8 {
		return domain.User{}, ErrWeakPassword
	}

	if _, _, err := s.users.GetByEmail(ctx, email); err == nil {
		return domain.User{}, ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}

	return s.users.Create(ctx, email, string(hash))
}

// Login devuelve el JWT y el usuario si las credenciales son correctas.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, domain.User, error) {
	user, hash, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return "", domain.User{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return "", domain.User{}, ErrInvalidCredentials
	}

	token, err := s.issueToken(user.ID)
	if err != nil {
		return "", domain.User{}, err
	}

	return token, user, nil
}

func (s *AuthService) issueToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ParseToken valida el JWT y devuelve el user id (subject) si es válido.
func (s *AuthService) ParseToken(tokenString string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}
	return claims.Subject, nil
}

func (s *AuthService) Me(ctx context.Context, userID string) (domain.User, error) {
	return s.users.GetByID(ctx, userID)
}
