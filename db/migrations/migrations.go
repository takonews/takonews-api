package migrations

import (
	"fmt"

	"github.com/takonews/takonews-api/app/models"
	"github.com/takonews/takonews-api/db"
)

func init() {
	fmt.Println("============Migration Begin=============")
	db.DB.LogMode(true)
	db.DB.AutoMigrate(&models.Article{}, &models.NewsSite{})
	db.DB.LogMode(false)
	fmt.Println("==========Migration Complete============")
}
