package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/takonews/takonews-api/app/models"
	"github.com/takonews/takonews-api/db"
)

// ArticleIndex show articles
// Available Query Parameters
// * sort
// * fields
func ArticleIndex(c *gin.Context) {
	// parameters
	articles := []models.Article{}
	sort := c.Query("sort")
	fields := c.Query("fields")

	/*
		DB processing
	*/
	// Sort
	// .Order("column [asc/desc]")
	sql := db.DB
	sorts := strings.Split(sort, ",")
	sql, err := OrderArticles(sql, sorts...)
	if err != nil {
		c.Status(http.StatusBadRequest)
		errorResp := ErrorResponse{Message: "wrong sort param"}
		encoder := json.NewEncoder(c.Writer)
		err = encoder.Encode(errorResp)
		if err != nil {
			panic(err)
		}

		return
	}
	sql.Find(&articles)

	/*
		select output field
	*/
	fs := strings.Split(fields, ",")
	results := SelectArticles(&articles, fs...)

	// set header
	c.Writer.Header().Set("Link", "<page=3>; rel=\"next\", <page=1>; rel=\"prev\", <page=5>; rel=\"last\"")
	c.Status(http.StatusOK)

	// write response
	b, err := json.MarshalIndent(results, "", " ")
	if err != nil {
		panic(err)
	}
	_, err = c.Writer.WriteString(string(b))
	if err != nil {
		panic(err)
	}
}

// ArticleShow show article details
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
	err := encoder.Encode(articles)
	if err != nil {
		c.Status(http.StatusBadRequest)
		errorResp := ErrorResponse{Message: "wrong sort param"}
		encoder := json.NewEncoder(c.Writer)
		err = encoder.Encode(errorResp)
		if err != nil {
			panic(err)
		}
	}
}

// SelectArticles extends gorm.DB.Select function
func SelectArticles(articles *[]models.Article, fields ...string) (results []map[string]interface{}) {
	for _, v := range *articles {
		if len(fields) > 1 || (len(fields) == 1 && fields[0] != "") {
			results = append(results, (&v).SelectFields(fields...))
		} else {
			results = append(results, (&v).SelectFields(
				"id",
				"title",
				"news_site_id",
				"published_at",
				"url",
			))
		}
	}

	return results
}

// OrderArticles extends gorm.DB.Order function
func OrderArticles(db *gorm.DB, sorts ...string) (*gorm.DB, error) {
	var dbRet = db
	var err error

	for i, v := range sorts {
		if v == "" {
			if i == 0 { // /articles /articles?sort= /articles?sort=,hoge
				dbRet = dbRet.Order("published_at desc")
			} else { // /articles?sort=hoge,
				return nil, err
			}
		} else if string(v[0]) == "-" { // sort=-hoge
			dbRet = dbRet.Order(v[1:] + " desc")
		} else { // sort=hoge
			dbRet = dbRet.Order(v + " asc")
		}
	}

	return dbRet, nil
}
