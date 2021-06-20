package service

import (
	"context"
	"crud/models"
	"crud/pkg/repository"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
)

type User interface {
	SignUp(user models.UserInput) error
	SignIn(user models.UserInput) (models.User, error)

	GenerateToken(user models.User) (string, error)
	ParseToken(token *jwt.Token) (int, error)

	GetOrCreateUserByEmail(email string) (models.User, error)

	GetGoogleAuthUrl() string
	GoogleAuthExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error)
	GetEmailFromGoogleAccessToken(ctx context.Context, accessToken *oauth2.Token) (string, error)
}

type Post interface {
	Create(userID int, data models.PostInput) (int, error)
	GetAll() ([]models.Post, error)
	GetById(id int) (models.Post, error)
	Update(userID, id int, data models.PostInput) error
	Delete(userID, id int) error
}

type Comment interface {
}

type Service struct {
	User
	Post
	Comment
}

func NewService(r *repository.Repository) *Service {
	return &Service{
		User: NewUserService(r.User),
		Post: NewPostService(r.Post),
	}
}
