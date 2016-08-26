package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/takonews/takonews-api/app/models"
	"github.com/takonews/takonews-api/db"
)

// ArticleIndex show articles
// Available Query Parameters
// * sort
// * fields
// * filter
func ArticleIndex(c *gin.Context) {
	// parameters
	var sort []string
	var fields []string
	var startDate time.Time
	var endDate time.Time

	sort = strings.Split(c.Query("sort"), ",")
	fields = strings.Split(c.Query("fields"), ",")
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	now := time.Now().In(loc)
	if c.Query("start-date") == "" { // default: today:00:00:00
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	} else {
		startDate, err = time.Parse("2006-01-02", c.Query("start-date"))
		startDate = startDate.Add(-9 * time.Hour).In(loc)
	}
	if c.Query("end-date") == "" { // default: tomorrow:00:00:00
		endDate = time.Date(now.Year(), now.Month(), now.Day(), 24, 0, 0, 0, loc)
	} else {
		endDate, err = time.Parse("2006-01-02", c.Query("end-date"))
		endDate = endDate.Add((-9 + 24) * time.Hour).In(loc)
	}

	/*
		DB processing
	*/
	articles := []models.Article{}
	sql := db.DB

	// filter
	sql = sql.Where("published_at >= ?", startDate)
	sql = sql.Where("published_at <= ?", endDate)

	// sort
	sql, err = OrderArticles(sql, sort...)

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

	// find
	sql.Find(&articles)

	/*
		select output field
	*/
	results := SelectArticles(&articles, fields...)

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
