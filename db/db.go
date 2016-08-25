package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // database driver init
	"github.com/takonews/takonews-api/config"
)

// DB is like a conn variable
var DB *gorm.DB

func init() {
	dbConfig := config.Config.DB

	var err error
	DB, err = gorm.Open(dbConfig.Adapter, fmt.Sprintf("%v:%v@/%v?parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Name))
	if err != nil {
		panic(err)
	}
}
