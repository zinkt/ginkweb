package gink

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(ctx *Context) {
		// Start timer
		t := time.Now()

		// 处理request
		ctx.Next()

		// 计算解析时间
		log.Printf("[%d] %s in %v", ctx.StatusCode, ctx.Req.RequestURI, time.Since(t))
	}
}
