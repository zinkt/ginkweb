package controllers

import (
	"bytes"
	"strings"

	"github.com/zinkt/ginkweb/ginkblog/models"
	"github.com/zinkt/ginkweb/ginkblog/storage"
	"github.com/zinkt/ginkweb/ginkblog/utils"
	"github.com/zinkt/ginkweb/ginkorm/log"
)

// 对每篇文章的内容
// 截取靠前一段字符，并删除每行前后的字符，便于展示
func GenerateAbstract(articles []models.Article) {
	for i := 0; i < len(articles); i++ {
		tmp := strings.Split(articles[i].Content[:utils.Min(666, len(articles[i].Content))], "\n")
		tmp = tmp[:len(tmp)-1]
		var b bytes.Buffer
		for j := 0; j < len(tmp); j++ {
			b.WriteString(strings.Trim(tmp[j], "#"))
			b.WriteByte('\n')
		}
		b.WriteString("......")
		articles[i].Content = b.String()
	}
}

// ********************* db **************************

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

// func GetArticleByCateAndTitle(cate, title string) *models.Article {
// 	s := storage.DB.NewSession()
// 	s.Model(&models.Article{})
// 	a := &models.Article{}
// 	err := s.Where("Category = ? AND Title = ?", cate, title).First(a)
// 	if err != nil {
// 		log.Infof("Article not found: Category = %s, Title = %s", cate, title)
// 		return nil
// 	}
// 	a.Viewed += 1
// 	s.Where("Category = ? AND Title = ?", cate, title).Update("Viewed", a.Viewed)
// 	return a
// }

func GetArticlesByCate(cate string) []models.Article {
	// 待缓存优化
	s := storage.DB.NewSession()
	s.Model(&models.Article{})
	var articles []models.Article
	err := s.Where("Category = ?", cate).OrderBy("LastUpdateTime").Find(&articles)
	if err != nil {
		log.Infof("Articles not found: Category = %s\n", cate)
		return nil
	}
	return articles
}

func GetAllArticles() []models.Article {
	// 待缓存优化
	s := storage.DB.NewSession()
	s.Model(&models.Article{})
	var articles []models.Article
	err := s.OrderBy("LastUpdateTime").Find(&articles)
	if err != nil {
		log.Infof("Articles not found\n")
		return nil
	}
	return articles
}
