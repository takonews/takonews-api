package controllers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
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
	//
	// query parameters
	//
	var sort []string
	var fields string
	var startDate time.Time
	var endDate time.Time
	var title string
	var offset int
	var limit int

	sort = strings.Split(c.Query("sort"), ",")
	fields = c.DefaultQuery("fields", "id,title,news_site_id,published_at,url")
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
	title = c.Query("title")
	offset, _ = strconv.Atoi(c.DefaultQuery("page[offset]", "0"))
	limit, _ = strconv.Atoi(c.DefaultQuery("page[limit]", "20"))

	//
	// building SQL query
	//
	sql := db.DB.Model(models.Article{})

	// filter
	sql = sql.Where("published_at BETWEEN ? AND ?", startDate, endDate)
	sql = sql.Where("title LIKE ?", "%"+title+"%")

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

	// limit
	sql = sql.Limit(limit)

	// offset
	sql = sql.Offset(offset)

	// select
	sql = sql.Select(fields)

	// find
	rows, err := sql.Rows()

	// select output fields
	var results []map[string]interface{}
	cols, err := rows.Columns()

	for rows.Next() {
		var row = make([]interface{}, len(cols))
		var rowp = make([]interface{}, len(cols))
		for i := 0; i < len(cols); i++ {
			rowp[i] = &row[i]
		}

		rows.Scan(rowp...)

		rowMap := make(map[string]interface{})
		for i, col := range cols {
			switch row[i].(type) {
			case []byte:
				row[i] = string(row[i].([]byte))
				num, err := strconv.Atoi(row[i].(string))
				if err == nil {
					row[i] = num
				}
			}
			rowMap[col] = row[i]
		}

		results = append(results, rowMap)
	}

	//
	// set header
	//
	// X-Total-Count
	var count int
	sql.Limit(-1).Offset(-1).Count(&count)
	c.Writer.Header().Set("X-Total-Count", strconv.Itoa(count))

	// Link
	var noLimitArticles []models.Article
	sql.Limit(-1).Offset(-1).Find(&noLimitArticles)
	noLimitCount := len(noLimitArticles)
	var links []string
	path := c.Request.URL.Path
	m, _ := url.ParseQuery(c.Request.URL.RawQuery)
	delete(m, "page[\"offset\"]")
	params := url.Values{}
	for k, v := range m {
		params.Set(k, v[0])
	}
	firstParams := params
	firstParams.Set("page[\"offset\"]", "0")
	linkFirstURL := path + "?" + firstParams.Encode()
	lastParams := params
	if noLimitCount-limit > 0 {
		lastParams.Set("page[\"offset\"]", strconv.Itoa(noLimitCount-limit))
	} else {
		lastParams = firstParams
	}
	linkLastParams := path + "?" + lastParams.Encode()
	prevParams := params
	if offset-limit >= 0 {
		prevParams.Set("page[\"offset\"]", strconv.Itoa(offset-limit))
	} else {
		prevParams.Set("page[\"offset\"]", strconv.Itoa(offset))
	}
	linkPrevParams := path + "?" + prevParams.Encode()
	nextParams := params
	if noLimitCount >= (offset + limit) {
		nextParams.Set("page[\"offset\"]", strconv.Itoa(offset+limit))
	} else {
		nextParams.Set("page[\"offset\"]", strconv.Itoa(offset))
	}
	linkNextParams := path + "?" + nextParams.Encode()
	links = append(links, "<"+linkFirstURL+">; rel=\"first\"")
	links = append(links, "<"+linkLastParams+">; rel=\"last\"")
	links = append(links, "<"+linkPrevParams+">; rel=\"prev\"")
	links = append(links, "<"+linkNextParams+">; rel=\"next\"")
	c.Writer.Header().Set("Link", strings.Join(links, ", "))

	// Content-Type
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
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

// OrderArticles extends gorm.DB.Order function
func OrderArticles(db *gorm.DB, sort ...string) (*gorm.DB, error) {
	var dbRet = db
	var err error

	for i, v := range sort {
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
