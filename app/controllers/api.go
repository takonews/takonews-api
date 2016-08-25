package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/takonews/takonews-api/app/models"
	"github.com/takonews/takonews-api/db"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func ArticleIndex(c *gin.Context) {
	// parameters
	articles := []models.Article{}
	sort := c.Query("sort")
	fields := c.Query("fields")

	sql := db.DB
	/*
		DB processing
	*/
	// Sort
	// .Order("column [asc/desc]")
	for i, v := range strings.Split(sort, ",") {
		if v == "" {
			if i == 0 { // /articles /articles?sort=
				sql = sql.Order("published_at desc")
			} else { // /articles?sort=hoge,
				c.Status(http.StatusBadRequest)
				errorResp := ErrorResponse{Message: "wrong sort param"}
				encoder := json.NewEncoder(c.Writer)
				encoder.Encode(errorResp)
				return
			}
		} else if string(v[0]) == "-" { // sort=-hoge
			sql = sql.Order(v[1:] + " desc")
		} else { // sort=hoge
			sql = sql.Order(v + " asc")
		}
	}
	sql.Find(&articles)

	/*
		select output field
	*/
	var data [](map[string]interface{})
	fs := strings.Split(fields, ",")
	for _, v := range articles {
		if len(fs) > 1 || (len(fs) == 1 && fs[0] != "") {
			data = append(data, (&v).SelectFields(fs...))
		}
	}

	// set header
	c.Writer.Header().Set("Link", "<page=3>; rel=\"next\", <page=1>; rel=\"prev\", <page=5>; rel=\"last\"")
	c.Status(http.StatusOK)

	// write response
	var results interface{}
	if len(data) > 0 {
		results = data
	} else {
		results = articles
	}
	b, err := json.MarshalIndent(results, "", " ")
	if err != nil {
		panic(err)
	}
	_, err = c.Writer.WriteString(string(b))
	if err != nil {
		panic(err)
	}
}

func ArticleShow(c *gin.Context) {
	// params
	// articlesID := c.Param("articles_id")
	articles := []models.Article{}
	db.DB.Find(&articles).Order("created_at desc")

	// set header
	c.Writer.Header().Set("Link", "<page=3>; rel=\"next\", <page=1>; rel=\"prev\", <page=5>; rel=\"last\"")
	c.Status(http.StatusOK)

	// write response
	encoder := json.NewEncoder(c.Writer)
	encoder.Encode(articles)
}
