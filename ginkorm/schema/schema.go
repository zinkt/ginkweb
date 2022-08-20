package schema

import (
	"go/ast"
	"reflect"

	"github.com/zinkt/ginkweb/ginkorm/dialect"
)

// 代表某一列（字段）
type Field struct {
	Name string
	Type string
	Tag  string
}

// 代表一张表
type Schema struct {
	Model interface{}
	// Model对应的类名
	Name       string
	Fields     []*Field
	FieldNames []string
	// 记录字段名和field的映射关系，无需变量field
	fieldMap map[string]*Field
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// 因为入参是一个对象的指针，因此需要 reflect.Indirect() 获取指针指向的实例
// 将dest指针指向的对象解析为d形式下的schema返回
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model: dest,
		// 以结构体名作为表名
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		// 判断	是否是内置类型 && 是否是可导出的（开头是否大写）
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("ginkorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}

	}
	return schema
}
