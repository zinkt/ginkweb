package gink

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 构建JSON数据时，显得更简洁
type H map[string]interface{}

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// 从Req中导出的信息
	Path   string
	Method string
	// 解析后的路由参数
	Params map[string]string
	// response info
	StatusCode int
	// 中间件
	handlers []HandlerFunc
	// 记录当前执行到第几个中间件
	index int
	// 为了通过Context访问engine中的HTML模板
	engine *Engine
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

// index是记录当前执行到第几个中间件，当在中间件中调用Next方法时，
// 控制权交给了下一个中间件，直到调用到最后一个中间件，
// 然后再从后往前，调用每个中间件在Next方法之后定义的部分。
func (ctx *Context) Next() {
	ctx.index++
	s := len(ctx.handlers)
	for ; ctx.index < s; ctx.index++ {
		ctx.handlers[ctx.index](ctx)
	}
}

// 获取的参数是预定义的
// 如：v2.GET("/hello/:name"...  ...ctx.Param("name")
func (ctx *Context) Param(key string) string {
	value := ctx.Params[key]
	return value
}

// 获取url中?之后的请求参数
func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

// 获取post方式的form表单参数
func (ctx *Context) PostForm(key string) string {
	return ctx.Req.FormValue(key)
}

// 设置状态码
func (ctx *Context) Status(code int) {
	ctx.StatusCode = code
	ctx.Writer.WriteHeader(code)
}

// 设置http请求头部信息，如内容格式 "Content-Type", "text/html"
func (ctx *Context) SetHeader(key string, value string) {
	ctx.Writer.Header().Set(key, value)
}

// 向ResponseWriter写回纯byte的数据
func (ctx *Context) Data(code int, data []byte) {
	ctx.Status(code)
	ctx.Writer.Write(data)
}

// 向ResponseWriter写回text/html格式的数据
// 描述???
func (ctx *Context) HTML(code int, name string, data interface{}) {
	ctx.SetHeader("Content-Type", "text/html")
	ctx.Status(code)
	// ExecuteTemplate applies the template associated with t
	// that has the given name to the specified data object
	// and writes the output to wr
	// 执行模板：将data填入ctx.engine.htmlTemplates中名为name的模板
	if err := ctx.engine.htmlTemplates.ExecuteTemplate(ctx.Writer, name, data); err != nil {
		ctx.Fail(500, err.Error())
	}
}

// 向ResponseWriter写回text/plain格式的数据
// 数据源：以format格式的values
func (ctx *Context) String(code int, format string, values ...interface{}) {
	ctx.SetHeader("Content-Type", "text/plain")
	ctx.Status(code)
	ctx.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 向ResponseWriter写回JSON格式的数据
// 数据源：obj
func (ctx *Context) JSON(code int, obj interface{}) {
	ctx.SetHeader("Content-Type", "application/json")
	ctx.Status(code)
	encoder := json.NewEncoder(ctx.Writer)
	if err := encoder.Encode((obj)); err != nil {
		http.Error(ctx.Writer, err.Error(), 500)
	}
}

// 向ResponseWriter写回失败信息
func (ctx *Context) Fail(code int, errormsg string) {
	ctx.Status(code)
	ctx.Writer.Write([]byte(errormsg))
}
