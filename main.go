package main

import (
	"net/http"

	"github.com/zinkt/ginkweb/gink"
)

func main() {
	// 创建一个gink引擎
	engine := gink.New()

	engine.Use(gink.Logger())

	// 设置路由， 和handler函数
	engine.GET("/index", func(ctx *gink.Context) {
		ctx.HTML(http.StatusOK, "<h1>this is index</h1>")
	})

	v1 := engine.Group("/v1")
	{
		v1.GET("/", func(ctx *gink.Context) {
			ctx.HTML(http.StatusOK, "<h1>hello this is /</h1>")
		})

		v1.GET("/hello", func(ctx *gink.Context) {
			// eg. expect /hello?name=zinkt
			ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
		})
	}

	v2 := engine.Group("/v2")
	{
		v2.GET("/hello/:name", func(ctx *gink.Context) {
			// eg. expect /hello/zinkt
			ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Param("name"), ctx.Path)
		})
		v2.POST("/login", func(ctx *gink.Context) {
			ctx.JSON(http.StatusOK, gink.H{
				"username": ctx.PostForm("username"),
				"password": ctx.PostForm("password"),
			})
		})
	}

	engine.GET("/assets/*filepath", func(ctx *gink.Context) {
		ctx.JSON(http.StatusOK, gink.H{"filepath": ctx.Param("filepath")})
	})

	engine.Run(":9999")
}
