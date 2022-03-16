package main

import (
	"net/http"

	"github.com/zinkt/ginkweb/gink"
)

func main() {
	engine := gink.New()

	engine.GET("/", func(c *gink.Context) {
		c.HTML(http.StatusOK, "<h1>hello my gee</h1>")
	})

	engine.POST("/login", func(c *gink.Context) {
		c.JSON(http.StatusOK, gink.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
			"fuck":     "you",
		})
	})

	engine.Run(":9999")
}
