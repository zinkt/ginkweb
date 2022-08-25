package controllers

import (
	"github.com/zinkt/ginkweb/ginkblog/models"
	"github.com/zinkt/ginkweb/ginkblog/storage"
)

func GetArticleById(id int) *models.Article {
	s := storage.DB.NewSession()
	s.Model(&models.Article{})
	a := &models.Article{}
	s.Where("Id = ?", id).First(a)
	return a
}
