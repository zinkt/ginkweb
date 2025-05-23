package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/shurcooL/github_flavored_markdown"
	"github.com/shurcooL/github_flavored_markdown/gfmstyle"
	"github.com/zinkt/ginkweb/gink"
	"github.com/zinkt/ginkweb/ginkblog/routes"
	"github.com/zinkt/ginkweb/ginkblog/storage"
)

func main() {
	// ginkDemo()
	// mytest()
	ginkblogDemo()
}

func mytest() {
	// Serve the "/assets/gfm.css" file.
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(gfmstyle.Assets)))

	var w io.Writer = os.Stdout // It can be an http.ResponseWriter.
	markdown := []byte("# GitHub Flavored Markdown\n\nHello.")

	io.WriteString(w, `<html><head><meta charset="utf-8"><link href="/assets/gfm.css" media="all" rel="stylesheet" type="text/css" /><link href="//cdnjs.cloudflare.com/ajax/libs/octicons/2.1.2/octicons.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`)
	w.Write(github_flavored_markdown.Markdown(markdown))
	io.WriteString(w, `</article></body></html>`)
}

func ginkblogDemo() {
	storage.CheckAndSnycArticles()
	g := routes.InitWeb()
	// go func() {
	// 	if err := g.Run_https(":443"); err != nil {
	// 		log.Fatal("HTTPS failed:", err)
	// 	}
	// }()
	g.Run(":8080")
}

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func ginkDemo() {
	// 创建一个gink引擎
	engine := gink.New()
	// 设置中间件
	engine.Use(gink.Logger(), gink.Recovery())

	// 恢复
	engine.GET("/panic", func(ctx *gink.Context) {
		names := []string{"zinkt"}
		ctx.String(http.StatusOK, names[10])
	})

	// 模板
	engine.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	engine.LoadHTMLGlob("./templates/*")
	engine.Static("/assets", "./static")

	// 设置路由， 和handler函数
	engine.GET("/", func(ctx *gink.Context) {
		ctx.HTML(http.StatusOK, "cssTest.tmpl", nil)
	})
	stu1 := &student{Name: "zinkt", Age: 20}
	stu2 := &student{Name: "jason", Age: 16}
	engine.GET("/students", func(ctx *gink.Context) {
		ctx.HTML(http.StatusOK, "cssTest.tmpl", gink.H{
			"title":  "gink",
			"stuArr": [2]*student{stu1, stu2},
		})
	})
	engine.GET("/date", func(ctx *gink.Context) {
		ctx.HTML(http.StatusOK, "cssTest.tmpl", gink.H{
			"title": "gink",
			"now":   time.Date(2022, 3, 21, 9, 40, 0, 0, time.UTC),
		})
	})
	v1 := engine.Group("/v1")
	// 括号为了好看
	{
		// v1.GET("/", func(ctx *gink.Context) {
		// 	ctx.HTML(http.StatusOK, "<h1>hello this is /</h1>")
		// })

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

	// engine.GET("/assets/*filepath", func(ctx *gink.Context) {
	// 	ctx.JSON(http.StatusOK, gink.H{"filepath": ctx.Param("filepath")})
	// })

	engine.Run(":9999")
}
