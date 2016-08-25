package models

import (
	"reflect"
	"time"
)

// Article is a model of news articles
type Article struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	NewsSite    NewsSite   `gorm:"ForeignKey:NewsSiteID;AssociationForeignKey:ID" json:"-"`
	NewsSiteID  uint       `json:"news_site_id"` //`gorm:"unique_index"` // Article belongs to NewsSite
	URL         string     `json:"url"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Query       string     `json:"query"`
	PublishedAt time.Time  `json:"published_at"`
}

func fieldSet(fields ...string) map[string]bool {
	set := make(map[string]bool, len(fields))
	for _, s := range fields {
		set[s] = true
	}

	return set
}

// SelectFields select column
// return value is not Article but map!!!
func (s *Article) SelectFields(fields ...string) map[string]interface{} {
	fs := fieldSet(fields...)
	rt, rv := reflect.TypeOf(*s), reflect.ValueOf(*s)
	out := make(map[string]interface{}, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		jsonKey := field.Tag.Get("json")
		if fs[jsonKey] {
			out[jsonKey] = rv.Field(i).Interface()
		}
	}

	return out
}
