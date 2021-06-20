package models

import (
	"encoding/xml"
	"errors"
)

type Post struct {
	XMLName xml.Name `xml:"post" gorm:"-" json:"-"`
	ID      int      `json:"id" xml:"id,attr" gorm:"primaryKey"`
	UserID  int      `json:"userId" xml:"user_id" gorm:"not null"`
	Title   string   `json:"title" xml:"title" gorm:"not null"`
	Body    string   `json:"body" xml:"body" gorm:"not null"`
}

type PostInput struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

var PostCreateValidationError error = errors.New("wrong data")

func (p *PostInput) Validate() error {
	if p.Title == "" || p.Body == "" {
		return PostCreateValidationError
	}
	return nil
}
