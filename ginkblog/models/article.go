package models

import (
	"time"
)

type Article struct {
	Id       uint `ginkorm:"PRIMARY KEY"`
	Title    string
	Content  string
	Category string
	// dialect 中有对time.Time的转换
	CreateTime     time.Time
	LastUpdateTime time.Time
	Viewed         uint64
	// author?
}
