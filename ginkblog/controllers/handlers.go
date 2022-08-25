package controllers

import (
	"net/http"
	"strconv"

	"github.com/zinkt/ginkweb/gink"
	"github.com/zinkt/ginkweb/ginkorm/log"
)

func Index(c *gink.Context) {

	c.HTML(http.StatusOK, "index/index", nil)

}

func ArticleDetail(c *gink.Context) {
	id, err := strconv.Atoi(c.Param("aid"))
	if err != nil {
		log.Error("Failed to get article id")
	}
	article := GetArticleById(id)

	data := gink.H{
		"article": article,
	}
	c.HTML(http.StatusOK, "article/detail", data)
}
