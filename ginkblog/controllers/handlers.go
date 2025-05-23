package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/zinkt/ginkweb/gink"
	"github.com/zinkt/ginkweb/ginkblog/utils"
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
		"article":        article,
		"requireCatalog": true,
	}
	c.HTML(http.StatusOK, "article/detail", data)

}

func CategoryArticleListIndex(c *gink.Context) {
	tmp := strings.Split(c.Path, "/")
	cate := tmp[len(tmp)-1]
	CategoryArticleListPageN(c, cate, 1)
}
func CategoryArticleListPaging(c *gink.Context) {
	tmp := strings.Split(strings.TrimSuffix(c.Path, "/"), "/")
	cate := tmp[len(tmp)-2]
	curpage, err := strconv.Atoi(c.Param("page"))
	if err != nil || curpage < 1 {
		log.Error("Failed to get page")
		c.HTML(http.StatusNotFound, "index/404", nil)
		return
	}
	CategoryArticleListPageN(c, cate, curpage)
}

func CategoryArticleListPageN(c *gink.Context, cate string, n int) {
	articles := GetArticlesByCate(cate)
	totalPage := (len(articles) + 4) / 5
	if articles == nil || n > totalPage {
		c.HTML(http.StatusNotFound, "index/404", nil)
		return
	}
	articles = articles[(n-1)*5 : utils.Min(n*5, len(articles))]
	GenerateAbstract(articles)
	data := gink.H{
		"articleList":    articles,
		"curpage":        n,
		"category":       cate,
		"totalPage":      totalPage,
		"requireCatalog": false,
	}
	c.HTML(http.StatusOK, "article/list", data)
}

func ArticleListPage(c *gink.Context) {
	curpage, err := strconv.Atoi(c.Param("page"))
	if err != nil || curpage < 1 {
		log.Error("Failed to get page")
		c.HTML(http.StatusNotFound, "index/404", nil)
		return
	}
	articles := GetAllArticles()
	totalPage := (len(articles) + 4) / 5
	if articles == nil || curpage > totalPage {
		c.HTML(http.StatusNotFound, "index/404", nil)
		return
	}
	articles = articles[(curpage-1)*5 : utils.Min(curpage*5, len(articles))]
	GenerateAbstract(articles)
	data := gink.H{
		"articleList":    articles,
		"curpage":        curpage,
		"totalPage":      totalPage,
		"requireCatalog": false,
	}
	c.HTML(http.StatusOK, "article/list", data)
}

// ChatRoomPage renders the chat room HTML page.
func ChatRoomPage(c *gink.Context) {
	c.HTML(http.StatusOK, "room.html", nil)
}
