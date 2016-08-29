package models

import "time"

// Article is a model of news articles
type Article struct {
	ID          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	NewsSite    NewsSite `gorm:"ForeignKey:NewsSiteID;AssociationForeignKey:ID"`
	NewsSiteID  uint     // Article belongs to NewsSite
	URL         string
	Title       string
	Description string
	FullText    string `gorm:"type:text"`
	PublishedAt time.Time
}
