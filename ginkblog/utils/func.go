package utils

import (
	"fmt"
	"os"
	"time"
)

func GetGoRunPath() string {
	path, _ := os.Getwd()
	return path
}

// str是否有在suffixes中的后缀
func HasSuffixes(str string, suffixes []string) bool {
	for _, v := range suffixes {
		if len(str) >= len(v) && str[len(str)-len(v):] == v {
			return true
		}
	}
	return false
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Plus(a, b int) int {
	return a + b
}

func Minus(a, b int) int {
	return a - b
}
