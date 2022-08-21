package session

import (
	"reflect"

	"github.com/zinkt/ginkweb/ginkorm/clause"
)

// 调用例子：
// s := ginkorm.NewEngine("sqlite3", "gink.db").NewSession()
// u1 := &User{Name: "Tom", Age: 18}
// u2 := &User{Name: "Sam", Age: 25}
// s.Insert(u1, u2, ...)
func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		// recordValues：多组RecordValues，与sql而言等价于 (A1, A2, A3, ...), (B1, B2, B3, ...), (C1, C2, C3, ...)
		recordValues = append(recordValues, table.RecordValues(value))
	}

	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// var users []User
// s.Find(&users);
// 传入一个切片指针，查询的结果保存在切片中
func (s *Session) Find(values interface{}) error {
	// 反射获取对象
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	// 构造sql，并查询到符合条件的记录rows
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	// 遍历每一行记录，利用反射创建destType的实例dest
	for rows.Next() {
		dest := reflect.New(destType).Elem()
		// 并将dest的所有字段平铺开，构造切片values
		var values []interface{}
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		// 将该行记录每一列的值依次赋值给values中的每一个字段
		if err := rows.Scan(values...); err != nil {
			return err
		}
		// 将 dest 添加到切片destSlice中。循环直到所有的记录都添加到切片destSlice中
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}
