package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitModels() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Link{})
	db.AutoMigrate(&LinkStatus{})
	db.AutoMigrate(&Warning{})

	return db, nil
}
