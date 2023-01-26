package models

import "gorm.io/gorm"

type Warning struct {
	gorm.Model
	LinkID uint
}

type WarningResponse struct {
	URL              string
	LinkId           string
	WarningCreatedAt string
	Threshold        string
}

func (warning *Warning) CreateWarning(db *gorm.DB) error {
	result := db.Create(warning)

	return result.Error
}

func GetWarning(db *gorm.DB, userId uint) ([]WarningResponse, error) {
	var warnings []WarningResponse
	result := db.Select("links.url as url, links.id as link_id, warnings.created_at as warning_created_at, links.error_threshold as threshold").Table("warnings").Joins("left join links on warnings.link_id = links.id").Where("links.user_id = ?", userId).Find(&warnings)
	if result.Error != nil {
		return nil, result.Error
	}

	return warnings, nil
}
