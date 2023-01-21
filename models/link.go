package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	URL            string       `json:"url" validate:"required"`
	ErrorThreshold int64        `json:"error_threshold" validate:"required"`
	UserID         uint         `json:"user_id"`
	LinkStatuses   []LinkStatus `json:"link_statuses"`
	Warnings       []Warning    `json:"warnings"`
}

type LinkStat struct {
	LinkID     int    `json:"link_id"`
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
	Count      int    `json:"count"`
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

func GetLinksStat(db *gorm.DB, userId uint) ([]LinkStat, error) {
	lastDay := time.Now().Add(-24 * time.Hour)
	var linkStat []LinkStat
	query := db.Select("status_code, link_id, url, count(*) as count").Table("links").Where("links.user_id = ?", userId).Joins("left join link_statuses on link_statuses.link_id = links.id").Where("link_statuses.created_at > ?", lastDay).Group("status_code, link_id").Find(&linkStat)

	log.Println("result i want to show", query.Error, linkStat)
	return linkStat, query.Error
}

func GetLinksStatByLinkID(db *gorm.DB, userId uint, linkID string) ([]LinkStat, error) {
	lastDay := time.Now().Add(-24 * time.Hour)
	var linkStat []LinkStat
	query := db.Select("status_code, link_id, url, count(*) as count").Table("links").Where("links.user_id = ? and links.id = ?", userId, linkID).Joins("left join link_statuses on link_statuses.link_id = links.id").Where("link_statuses.created_at > ?", lastDay).Group("status_code, link_id").Find(&linkStat)

	log.Println("result i want to show", query.Error, linkStat)
	return linkStat, query.Error
}
