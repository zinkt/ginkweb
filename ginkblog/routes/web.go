package routes

import (
	"html/template"
	"path/filepath"

	"github.com/zinkt/ginkweb/gink"
	"github.com/zinkt/ginkweb/ginkblog/controllers"
	"github.com/zinkt/ginkweb/ginkblog/utils"
)

func InitWeb() *gink.Engine {
	g := gink.Default()

	// 注册中间件
	g.Use(gink.Logger(), gink.Recovery())
	// 添加自定义模板渲染函数
	g.SetFuncMap(template.FuncMap{
		"FormatAsDate": utils.FormatAsDate,
	})

	// 此处filepath.Join()会Clean掉多余的separator，插入"OS specific Separator"
	g.Static("/static", filepath.Join(utils.GetGoRunPath(), "static"))
	g.LoadHTMLGlob(filepath.Join(utils.GetGoRunPath(), "views", "*", "*.html"))

	g.GET("/", controllers.Index)

	category := g.Group("/category")
	{
		category.GET("/coding", nil)
		category.GET("/share", nil)
		category.GET("/thinking", nil)
	}

	archives := g.Group("/archives")
	{
		archives.GET("/", nil)
		archives.GET("/:aid", controllers.ArticleDetail)
	}

	g.GET("/about", nil)

	return g
}
