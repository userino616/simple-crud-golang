package repository

import (
	"crud/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(data models.UserInput) (models.User, error) {
	u := models.User{Email: data.Email, Password: data.Password}
	db := r.db.Create(&u)
	return u, db.Error
}

func (r *UserRepository) GetUser(email, password string) (models.User, error) {
	var u models.User
	db := r.db.Where("email = ? AND password = ?", email, password).First(&u)

	return u, db.Error
}

func (r *UserRepository) GetUserByEmail(email string) (models.User, error) {
	var u models.User
	db := r.db.Where("email = ?", email).First(&u)
	return u, db.Error
}
