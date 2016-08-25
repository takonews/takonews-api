package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takonews/takonews-api/app/controllers"
)

func Router() *http.ServeMux {
	router := gin.Default()
	router.GET("/api/v2/articles", controllers.ArticleIndex)
	router.GET("/api/v2/articles/:articles_id", controllers.ArticleShow)

	var mux = http.NewServeMux()
	mux.Handle("/", router)

	return mux
}
