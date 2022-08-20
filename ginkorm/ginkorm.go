package ginkorm

import (
	"database/sql"

	"github.com/zinkt/ginkweb/ginkorm/dialect"
	"github.com/zinkt/ginkweb/ginkorm/log"
	"github.com/zinkt/ginkweb/ginkorm/session"
)

// Engine负责交互前的连接测试、交互后的关闭连接
type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

// 当driver为sqlite3时，source为db路径
func NewEngine(driver, source string) (engine *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	// 保证指定的dialect存在
	d, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}

	engine = &Engine{db: db, dialect: d}
	log.Infof("Connect database %s:%s success", driver, source)
	return
}

func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error(err)
	}
	log.Info("Close database success")
}

func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect)
}
