package repository

import (
	"crud/models"
	"gorm.io/gorm"
)

type Post interface {
	Create(userID int, data models.PostInput) (int, error)
	GetAll() ([]models.Post, error)
	GetById(id int) (models.Post, error)
	Update(userID, id int, data models.PostInput) error
	Delete(userID, id int) error
}

type Comment interface {
}

type User interface {
	CreateUser(data models.UserInput) (models.User, error)
	GetUser(email, password string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
}

type Repository struct {
	Post
	Comment
	User
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
		Post: NewPostRepository(db),
	}
}
