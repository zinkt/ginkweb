package gink

import (
	"html/template"
	"net/http"
	"path"
	"strings"
)

// HandlerFunc defines the request handler used by gink
type HandlerFunc func(*Context)

type (
	// 只是包装了一层，addRoute()时加上前缀
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc
		parent      *RouterGroup
		engine      *Engine // 所有group共享一个engine，便于访问router
	}

	// Engine实现了ServeHTTP的接口
	Engine struct {
		*RouterGroup
		router *router
		groups []*RouterGroup // 保存所有group
		// for template : html render
		htmlTemplates *template.Template
		//自定义渲染函数
		funcMap template.FuncMap
	}
)

// 构造器
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

// 用于创建新RouterGroup
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup) // 创建后需要加入全局groups
	return newGroup
}

// 静态文件处理
// http.FileServer()方法返回的是fileHandler实例，fileHandler结构体实现了Handler接口中的ServerHTTP()方法。
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	// absPath 如 /v1/assets/css/zinkt.css
	wholePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(wholePath, http.FileServer(fs))
	return func(ctx *Context) {
		// 此处获取真实用户请求的参数
		file := ctx.Param("filepath")
		// 检查文件是否可用（存在/有权限读取）
		if _, err := fs.Open(file); err != nil {
			ctx.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(ctx.Writer, ctx.Req)
	}
}

// xxx.com/assets
// 如/assets/js/zinkt.js 为 relativePath
// http.Dir()方法会返回http.Dir类型用于将字符串路径转换为文件系统
// root若为相对路径，则是相对于入口函数
func (group *RouterGroup) Static(relativePath string, root string) {
	// 处理路径：系统路径root ---> web相对路径relativePath
	s := http.Dir(root)
	handler := group.createStaticHandler(relativePath, s)
	urlPattern := path.Join(relativePath, "/*filepath")
	// 为静态文件，注册GET处理函数
	group.GET(urlPattern, handler)
}

// **************模板渲染*****************

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// ParseGlob() 解析pattern匹配到的所有HTML文件，将结果模板们与调用这个函数的*template相关联
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// 中间件函数
// 向group中添加中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// group.prefix是本group的前缀
func (group *RouterGroup) addRoute(method string, component string, handler HandlerFunc) {
	pattern := group.prefix + component
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// 实现了ServeHTTP接口，用于接管所有http请求
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	// 将所有前缀符合的中间件加入到ctx中
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	ctx := newContext(w, req)
	ctx.handlers = middlewares
	ctx.engine = engine
	engine.router.handle(ctx)
}
