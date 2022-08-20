package gink

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func trace() string {
	var pcs [32]uintptr
	// number of entries written to pc
	n := runtime.Callers(3, pcs[:]) // 跳过了前三个调用

	var str strings.Builder
	str.WriteString("\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
func Recovery() HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", msg+trace())
				ctx.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		ctx.Next()
	}
}
