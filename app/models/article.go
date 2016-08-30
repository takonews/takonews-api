package models

import "time"

// Article is a model of news articles
type Article struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	// Article belongs to NewsSite
	NewsSite    NewsSite `gorm:"ForeignKey:NewsSiteID;AssociationForeignKey:ID"`
	NewsSiteID  uint     `gorm:"not null"`
	URL         string   `gorm:"not null"`
	Title       string   `gorm:"not null"`
	Description string
	FullText    string    `gorm:"type:text"`
	PublishedAt time.Time `gorm:"not null"`
}
