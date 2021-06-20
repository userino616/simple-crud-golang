package repository

import (
	"crud/models"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(userID int, data models.PostInput) (int, error) {
	p := models.Post{UserID: userID, Title: data.Title, Body: data.Body}
	db := r.db.Create(&p)
	return p.ID, db.Error
}

func (r *PostRepository) GetAll() ([]models.Post, error) {
	var posts []models.Post
	query := r.db.Find(&posts)
	return posts, query.Error
}

func (r *PostRepository) GetById(id int) (models.Post, error) {
	var post models.Post
	query := r.db.First(&post, id)
	return post, query.Error
}

func (r *PostRepository) Delete(userID, id int) error {
	db := r.db.Where("user_id = ? AND id = ?", userID, id).Delete(&models.Post{})

	if db.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return db.Error
}

func (r *PostRepository) Update(userID, id int, data models.PostInput) error {
	p := models.Post{UserID: userID, ID: id, Title: data.Title, Body: data.Body}
	db := r.db.Model(&p).Updates(&p)

	if db.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return db.Error
}
