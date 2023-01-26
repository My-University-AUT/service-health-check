package models

import "gorm.io/gorm"

type LinkStatus struct {
	gorm.Model
	StatusCode int `json:"status_code"`
	LinkID     uint
}

func (linkStatus *LinkStatus) CreateLinkStatus(db *gorm.DB) error {
	result := db.Create(linkStatus)

	return result.Error
}
