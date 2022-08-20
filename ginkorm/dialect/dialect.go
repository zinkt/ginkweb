package dialect

import "reflect"

var dialectMap = map[string]Dialect{}

type Dialect interface {
	// 将Go数据类型转化为该数据库的数据类型
	DataTypeOf(typ reflect.Value) string
	// 判断表是否存在
	TableExistSQL(tableName string) (string, []interface{})
}

func RegisterDialect(name string, dialect Dialect) {
	dialectMap[name] = dialect
}

func GetDialect(name string) (Dialect, bool) {
	d, ok := dialectMap[name]
	return d, ok
}
