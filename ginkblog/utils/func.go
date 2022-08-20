package utils

import "os"

func GetGoRunPath() string {
	path, _ := os.Getwd()
	return path
}
