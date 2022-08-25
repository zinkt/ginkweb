package routes

import (
	"path/filepath"

	"github.com/zinkt/ginkweb/gink"
	"github.com/zinkt/ginkweb/ginkblog/controllers"
	"github.com/zinkt/ginkweb/ginkblog/utils"
)

func InitWebRoutes() *gink.Engine {
	g := gink.Default()

	// 此处filepath.Join()会Clean掉多余的separator
	g.Static("/static", filepath.Join(utils.GetGoRunPath(), "/static"))
	g.LoadHTMLGlob(filepath.Join(utils.GetGoRunPath(), "/views/*/*.html"))

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
		archives.GET("/:aid", nil)
	}

	g.GET("/about", nil)
	return g
}
