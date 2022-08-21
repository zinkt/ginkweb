package session

import (
	"database/sql"
	"strings"

	"github.com/zinkt/ginkweb/ginkorm/clause"
	"github.com/zinkt/ginkweb/ginkorm/dialect"
	"github.com/zinkt/ginkweb/ginkorm/log"
	"github.com/zinkt/ginkweb/ginkorm/schema"
)

type Session struct {
	db *sql.DB
	// 带占位符的sql语句
	sql strings.Builder
	// SQL 语句中占位符的对应值
	sqlVars []interface{}

	dialect dialect.Dialect
	// 对象映射成的表
	refTable *schema.Schema

	clause clause.Clause
}

// 新建session
func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

// 清除已有的sql
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

// 返回session对应的db
func (s *Session) DB() *sql.DB {
	return s.db
}

// 设置准备执行的sql及其占位符values
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// **** 封装Exec()、Query() 和 QueryRow()用于清空sql和打印日志 ****

// 执行sql
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// 查询返回一条记录
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// 查询返回多行记录
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
