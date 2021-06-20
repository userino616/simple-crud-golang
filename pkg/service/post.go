package service

import (
	"crud/models"
	"crud/pkg/repository"
)

type PostService struct {
	repo repository.Post
}

func NewPostService(repo repository.Post) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) Create(userID int, data models.PostInput) (int, error) {
	return s.repo.Create(userID, data)
}

func (s *PostService) GetAll() ([]models.Post, error) {
	return s.repo.GetAll()
}

func (s *PostService) GetById(id int) (models.Post, error) {
	return s.repo.GetById(id)
}

func (s *PostService) Delete(userID, id int) error {
	return s.repo.Delete(userID, id)
}

func (s *PostService) Update(userID, id int, data models.PostInput) error {
	return s.repo.Update(userID, id, data)
}
