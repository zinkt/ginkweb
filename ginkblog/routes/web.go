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
		"Plus":         utils.Plus,
		"Minus":        utils.Minus,
		"Unescaped":    utils.Unescaped,
	})

	// 此处filepath.Join()会Clean掉多余的separator，插入"OS specific Separator"
	g.Static("/static", filepath.Join(utils.GetGoRunPath(), "static"))
	g.Static("/.well-known", filepath.Join(utils.GetGoRunPath(), ".well-known"))
	g.StaticFile("/favicon.ico", filepath.Join(utils.GetGoRunPath(), "static", "img", "favicon.ico"))
	g.LoadHTMLGlob(filepath.Join(utils.GetGoRunPath(), "views", "*", "*.html"))

	g.GET("/", controllers.Index)

	// 不相信用户输入doge
	category := g.Group("/category")
	{
		category.GET("/coding", controllers.CategoryArticleListIndex)
		category.GET("/share", controllers.CategoryArticleListIndex)
		category.GET("/thinking", controllers.CategoryArticleListIndex)

		category.GET("/coding/:page", controllers.CategoryArticleListPaging)
		category.GET("/share/:page", controllers.CategoryArticleListPaging)
		category.GET("/thinking/:page", controllers.CategoryArticleListPaging)

	}

	archives := g.Group("/archives")
	{
		// archives.GET("/", nil)
		archives.GET("/:aid", controllers.ArticleDetailById)
	}

	// g.GET("/about", nil)

	// WebSocket chat route
	g.GET("/ws/chat", controllers.Chat)

	// Page to test chat
	g.GET("/chat", controllers.ChatRoomPage)

	return g
}
