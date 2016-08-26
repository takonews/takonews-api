package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takonews/takonews-api/app/controllers"
	"github.com/takonews/takonews-api/config"
)

func Router() *http.ServeMux {
	router := gin.Default()

	users := config.Config.Secret.Users
	accounts := make(map[string]string)
	for _, v := range users {
		accounts[v["name"]] = v["password"]
	}
	authorized := router.Group("/")
	authorized.Use(gin.BasicAuth(accounts))
	authorized.GET("/api/v2/articles", controllers.ArticleIndex)
	authorized.GET("/api/v2/articles/:articles_id", controllers.ArticleShow)

	var mux = http.NewServeMux()
	mux.Handle("/", router)

	return mux
}
