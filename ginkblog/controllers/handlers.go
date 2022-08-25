package controllers

import (
	"net/http"

	"github.com/zinkt/ginkweb/gink"
)

func Index(c *gink.Context) {

	c.HTML(http.StatusOK, "index/index", nil)

}

func ArticleDetail(c *gink.Context) {

}
