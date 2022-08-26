package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/zinkt/ginkweb/gink"
	"github.com/zinkt/ginkweb/ginkorm/log"
)

func Index(c *gink.Context) {

	c.HTML(http.StatusOK, "index/index", nil)

}

func ArticleDetailById(c *gink.Context) {
	id, err := strconv.Atoi(c.Param("aid"))
	if err != nil {
		log.Error("Failed to get article id")
	}
	article := GetArticleById(id)
	if article == nil {
		c.HTML(http.StatusNotFound, "index/404", nil)
		return
	}
	// article.Viewed += 1
	data := gink.H{
		"article": article,
	}
	c.HTML(http.StatusOK, "article/detail", data)

}

func ArticleDetailByTitle(c *gink.Context) {
	tmp := strings.Split(c.Path, "/")
	article := GetArticleByCateAndTitle(tmp[len(tmp)-2], c.Param("atitle"))
	if article == nil {
		c.HTML(http.StatusNotFound, "index/404", nil)
		return
	}
	data := gink.H{
		"article": article,
	}
	c.HTML(http.StatusOK, "article/detail", data)
}
