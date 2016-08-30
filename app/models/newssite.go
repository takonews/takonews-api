package models

import "github.com/jinzhu/gorm"

// NewsSite is a model of news sites
type NewsSite struct {
	gorm.Model // include standard field
	Name       string
	URL        string
}
