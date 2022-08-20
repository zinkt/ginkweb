package gink

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node       //eg. roots['GET'] roots['POST']
	handlers map[string]HandlerFunc //eg. handlers['GET-/p/:lang/doc']
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 将pattern按'/'分割为parts并返回
// 当某part以*开头时，不再转化后续parts（暂时），表示匹配了所有后续内容
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 添加method方法的pattern到路由中，并绑定handler
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method] // 根据方法获取Trie树根节点
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// 根据path（此处为用户传入的路径）匹配到一个对应的节点，并解析出匹配符（':''*'）的参数
// 返回对应节点和参数
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		log.Println("kidding? no method?")
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		// 解析出':'和'*'通配符的参数
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			// 如pattern=*filepath，path=css/zinkt.css
			// 解析结果为"filepath" : "css/zinkt.css"
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				// 解析到*即停止
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// 根据router处理ctx
func (r *router) handle(ctx *Context) {
	n, params := r.getRoute(ctx.Method, ctx.Path)
	if n != nil {
		ctx.Params = params // 在ctx中保存根据pattern解析的params
		key := ctx.Method + "-" + n.pattern
		ctx.handlers = append(ctx.handlers, r.handlers[key]) //添加对应的handler，并在下面的Next()中一起调用
	} else {
		ctx.handlers = append(ctx.handlers, func(ctx *Context) {
			ctx.String(http.StatusNotFound, "404 NOT FOUND: %s\n", ctx.Path)
		})
	}
	ctx.Next()
}
