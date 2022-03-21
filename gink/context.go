package gink

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int
	// 中间件
	handlers []HandlerFunc
	index    int //记录当前执行到第几个中间件
	// 模板部分：为了通过Context访问engine中的HTML模板
	engine *Engine
}

func (ctx *Context) Param(key string) string {
	value, _ := ctx.Params[key]
	return value
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

func (ctx *Context) PostForm(key string) string {
	return ctx.Req.FormValue(key)
}

func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get((key))
}
func (ctx *Context) Status(code int) {
	ctx.StatusCode = code
	ctx.Writer.WriteHeader(code)
}
func (ctx *Context) SetHeader(key string, value string) {
	ctx.Writer.Header().Set(key, value)
}
func (ctx *Context) String(code int, format string, values ...interface{}) {
	ctx.SetHeader("Content-Type", "text/plain")
	ctx.Status(code)
	ctx.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}
func (ctx *Context) JSON(code int, obj interface{}) {
	ctx.SetHeader("Content-Type", "application/json")
	ctx.Status(code)
	encoder := json.NewEncoder(ctx.Writer)
	if err := encoder.Encode((obj)); err != nil {
		http.Error(ctx.Writer, err.Error(), 500)
	}
}
func (ctx *Context) Data(code int, data []byte) {
	ctx.Status(code)
	ctx.Writer.Write(data)
}
func (ctx *Context) HTML(code int, name string, data interface{}) {
	ctx.SetHeader("Content-Type", "text/html")
	ctx.Status(code)
	if err := ctx.engine.htmlTemplates.ExecuteTemplate(ctx.Writer, name, data); err != nil {
		ctx.Fail(500, err.Error())
	}
}

func (ctx *Context) Fail(code int, errormsg string) {
	ctx.Status(code)
	ctx.Writer.Write([]byte(errormsg))
}
