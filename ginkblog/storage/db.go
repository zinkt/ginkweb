package storage

import (
	"io/fs"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shurcooL/github_flavored_markdown"
	"github.com/zinkt/ginkweb/ginkblog/models"
	"github.com/zinkt/ginkweb/ginkblog/utils"
	"github.com/zinkt/ginkweb/ginkorm"
	"github.com/zinkt/ginkweb/ginkorm/log"
)

var DB *ginkorm.Engine

func init() {
	e, err := ginkorm.NewEngine("sqlite3", filepath.Join(utils.GetGoRunPath(), "storage", "ginkblog.db"))
	if err != nil {
		log.Error(err)
		panic(err)
	}
	DB = e
}

func CheckAndSnycArticles() {
	s := DB.NewSession()
	if !s.Model(&models.Article{}).HasTable() {
		log.Error("Article table not exists")
		_ = s.CreateTable()
	}
	// 获取已发布过的文章列表
	var loaded []models.Article
	s.OrderBy("Id DESC").Find(&loaded)
	// 组织成[Path] : *models.Article的形式确定某文件是否已发布
	loadedMap := make(map[string]*models.Article, len(loaded))
	for i := 0; i < len(loaded); i++ {
		loadedMap[loaded[i].RelativePath] = &loaded[i]
	}
	// 设置自增id
	// ???多此一举而不设置AUTOINCREMENT的原因是，Insert()直接将Article中的id以空值0插入，拆解较为复杂，暂时搁置
	var count uint
	if len(loaded) != 0 {
		count = loaded[0].Id
	}
	filepath.Walk(filepath.Join(utils.GetGoRunPath(), "storage", "articles"),
		func(fullpath string, info fs.FileInfo, err error) error {
			if err != nil {
				log.Error(err)
				return err
			}
			if info.IsDir() {
				return nil
			}
			filename := info.Name()
			filesuffix := path.Ext(filename)
			fileprefix := strings.TrimSuffix(filename, filesuffix)
			gorunpath := utils.GetGoRunPath()
			relativePath := fullpath[len(gorunpath)+len("/storage/articles/"):]
			cate := strings.Split(relativePath, "/")[0]
			// 如果是这些后缀的文件才载入	未来用配置文件的方式
			suffixes := [...]string{".md", ".txt"}
			if utils.HasSuffixes(filesuffix, suffixes[:]) {
				// 便于获取category，即tmp[1]
				// 如果已经发布过
				if art, exists := loadedMap[relativePath]; exists {
					// 若修改过，则更新
					if !art.LastUpdateTime.Equal(info.ModTime()) {
						bytes, err := ioutil.ReadFile(fullpath)
						if err != nil {
							log.Errorf("failed to read %s", fullpath)
							return err
						}
						// 解析md文章为html
						bytes = github_flavored_markdown.Markdown(bytes)
						// 更新数据库
						n, err := s.Where("Id = ?", art.Id).Update("Content", string(bytes), "LastUpdateTime", info.ModTime())
						if err != nil {
							log.Error("Failed to update")
						}
						log.Infof("Update %s success, %d row(s) affected", info.Name(), n)
					}
				} else {
					bytes, err := ioutil.ReadFile(fullpath)
					if err != nil {
						log.Error("failed to read %s", fullpath)
						return err
					}
					// 解析md文章为html
					bytes = github_flavored_markdown.Markdown(bytes)
					s.Insert(&models.Article{
						// ???
						// 若不指定id，这里id在插入时会插入空值值0
						Id:      count + 1,
						Title:   fileprefix,
						Content: string(bytes),
						// 路径的倒数第二个
						Category:       cate,
						RelativePath:   relativePath,
						CreateTime:     time.Now(),
						LastUpdateTime: info.ModTime(),
						Viewed:         0,
					})
					count++
					log.Infof("Load %s success", info.Name())
				}
			}
			return nil
		})
}
