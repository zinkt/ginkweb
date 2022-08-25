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
func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}
