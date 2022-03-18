package main

import (
	"net/http"

	"github.com/zinkt/ginkweb/gink"
)

func main() {
	// 创建一个gink引擎
	engine := gink.New()

	// 设置路由， 和handler函数
	engine.GET("/", func(ctx *gink.Context) {
		ctx.HTML(http.StatusOK, "<h1>hello my gink</h1>")
	})
	engine.GET("/hello", func(ctx *gink.Context) {
		// eg. expect /hello?name=zinkt
		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
	})
	engine.GET("/hello/:name", func(ctx *gink.Context) {
		// eg. expect /hello/zinkt
		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Param("name"), ctx.Path)
	})

	engine.GET("/assets/*filepath", func(ctx *gink.Context) {
		ctx.JSON(http.StatusOK, gink.H{"filepath": ctx.Param("filepath")})
	})

	engine.Run(":9999")
}
