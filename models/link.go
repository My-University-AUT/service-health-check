package models

import (
	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	URL    string `json:"url" validate:"required"`
	UserID uint   `json:"user_id"`
}

func (link *Link) CreateLink(db *gorm.DB) error {
	result := db.Create(link)

	return result.Error
}

func GetLink(db *gorm.DB, userId uint) ([]Link, error) {
	var links []Link
	result := db.Where(&Link{UserID: userId}).Find(&links)

	return links, result.Error
}
