package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

// 使用 log.Lshortfile 支持显示文件名和代码行号
var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex
)

// log 方法，默认打印日志到os.std
// [info ] 颜色为蓝色，[error] 为红色
var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

// log 层级
const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

// 设置层级
func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}
	// 大于1时 取消输出error
	if level > ErrorLevel {
		errorLog.SetOutput(ioutil.Discard)
	}
	// 大于0时 取消输出info
	if level > InfoLevel {
		infoLog.SetOutput(ioutil.Discard)
	}
}
