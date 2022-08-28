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

func CategoryArticlesIndex(c *gink.Context) {
	tmp := strings.Split(c.Path, "/")
	articles := GetArticlesByCate(tmp[len(tmp)-1])
	if len(articles) > 5 {
		articles = articles[:5]
	}
	for i := 0; i < len(articles); i++ {
		articles[i].Content = articles[i].Content[:46] + "..."
	}
	data := gink.H{
		"articleList": articles,
		"page":        1,
	}

	c.HTML(http.StatusOK, "article/list", data)
}
func CategoryArticlesPaging(c *gink.Context) {
	tmp := strings.Split(c.Path, "/")
	curpage, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		log.Error("Failed to get page")
	}
	articles := GetArticlesByCate(tmp[len(tmp)-2])
	if len(articles) > 5 {
		articles = articles[(curpage-1)*5 : curpage*5]
	}
	for i := 0; i < len(articles); i++ {
		articles[i].Content = articles[i].Content[:46] + "..."
	}
	data := gink.H{
		"articleList": articles,
		"curpage":     curpage,
	}

	c.HTML(http.StatusOK, "article/list", data)
}
