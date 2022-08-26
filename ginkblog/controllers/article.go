package controllers

import (
	"github.com/zinkt/ginkweb/ginkblog/models"
	"github.com/zinkt/ginkweb/ginkblog/storage"
	"github.com/zinkt/ginkweb/ginkorm/log"
)

func GetArticleById(id int) *models.Article {
	s := storage.DB.NewSession()
	s.Model(&models.Article{})
	a := &models.Article{}
	err := s.Where("Id = ?", id).First(a)
	if err != nil {
		log.Infof("Article not found: Id = %d", id)
		return nil
	}
	a.Viewed += 1
	s.Where("Id = ?", id).Update("Viewed", a.Viewed)
	return a
}

func GetArticleByCateAndTitle(cate, title string) *models.Article {
	s := storage.DB.NewSession()
	s.Model(&models.Article{})
	a := &models.Article{}
	err := s.Where("Category = ? AND Title = ?", cate, title).First(a)
	if err != nil {
		log.Infof("Article not found: Category = %s, Title = %s", cate, title)
		return nil
	}
	a.Viewed += 1
	s.Where("Category = ? AND Title = ?", cate, title).Update("Viewed", a.Viewed)
	return a
}
