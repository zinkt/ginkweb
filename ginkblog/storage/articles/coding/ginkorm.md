### sqlite3使用

```bash
sqlite3 gink.db #进入数据库
.head on #打开列名开关
.table #显示表
.schema User #查看建表sql
```

### log

封装了log标准库，主要部分：

```go
var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex
)
```

### Session

封装了database/sql 标准库的`Exec()`、`Query()` 和 `QueryRow()` 三个方法，主要目的是：

1. 统一打印日志 （包括 执行的SQL 语句和错误日志） 
2. 执行完成后，清空 `(s *Session).sql` 和 `(s *Session).sqlVars` 两个变量。这样 Session 可以复用，开启一次会话，可以执行多次 SQL。 

### Engine

主要完成数据库连接检查和关闭

```go
func NewEngine(driver, source string) (engine *Engine, err error)
func (engine *Engine) Close()
func (engine *Engine) NewSession() *session.Session 
```

### Dialect

将Go语言的类型映射为数据库中的类型

#### 抽象出各个数据库中的差异部分	dialect.go

```
var dialectsMap = map[string]Dialect{}

type Dialect interface {
	DataTypeOf(typ reflect.Value) string
	TableExistSQL(tableName string) (string, []interface{})
}
func RegisterDialect(name string, dialect Dialect) 
func GetDialect(name string) (dialect Dialect, ok bool)
```

#### 增加对sqlite3的支持(实现dialect接口)	sqlite3.go

```

func init() {
	RegisterDialect("sqlite3", &sqlite3{})
}
func (s *sqlite3) DataTypeOf(typ reflect.Value) string
func (s *sqlite3) TableExistSQL(tableName string) (string, []interface{})
```

### Schema

对象和表的转换

```
type Schema struct {
	Model interface{}
	// Model对应的类名
	Name       string
	Fields     []*Field
	FieldNames []string
	// 记录字段名和field的映射关系，无需变量field
	fieldMap map[string]*Field
}
```

