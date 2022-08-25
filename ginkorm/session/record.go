package session

import (
	"errors"
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

// 接受 2 种入参，平铺开来的键值对和 map 类型的键值对
// kv 为 map[string]interface{}
// 或者 kv list : "Name", "Tom", "Age", 18
func (s *Session) Update(kv ...interface{}) (int64, error) {
	// 判断传入的类型，如果不是map类型，则会转换为map
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Delete() (int64, error) {
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

// ******** 链式 ********
// s.Where("Age > 18").Limit(3).Find(&users)

// 向clause新增limit
func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

// 返回一条记录
// 如	u := &User{}
// 		_ = s.OrderBy("Age DESC").First(u)
// 根据传入的类型，利用反射构造切片，调用 Limit(1) 限制返回的行数，
// 调用 Find 方法获取到查询结果，若成功查到，则将信息写入切片
func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("First() RECORD NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}
