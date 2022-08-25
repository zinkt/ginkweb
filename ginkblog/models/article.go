package models

type Article struct {
	Id             uint
	Title          string
	Content        string
	CateId         uint
	CreateTime     int64
	LastUpdateTime int64
	// tags? author?
}

// func GetArticleById(id int) *Article {

// }
