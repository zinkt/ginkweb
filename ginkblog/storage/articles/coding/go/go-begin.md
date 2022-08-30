## mod包管理

### 网络包

若没自动获取，则```go get github.com/mattn/go-sqlite3 ```

### 本地包

go mod 中

```replace gink => ../gink v0.0.0```

### 更好的方式

#### 在同一个项目下

**注意**：在一个项目（project）下我们是可以定义多个包（package）的。

##### 目录结构

现在的情况是，我们在`moduledemo/main.go`中调用了`mypackage`这个包。

```
moduledemo
├── go.mod
├── main.go
└── mypackage
    └── mypackage.go
```

##### 导入包

这个时候，我们需要在`moduledemo/go.mod`中按如下定义：

```
module moduledemo

go 1.14
```

然后在`moduledemo/main.go`中按如下方式导入`mypackage`

```
package main

import (
    "fmt"
    "moduledemo/mypackage"  // 导入同一项目下的mypackage包
)
func main() {
    mypackage.New()
    fmt.Println("main")
}
```

#### 举个例子

举一反三，假设我们现在有文件目录结构如下：

```
└── bubble
    ├── dao
    │   └── mysql.go
    ├── go.mod
    └── main.go
```

其中`bubble/go.mod`内容如下：

```
module github.com/q1mi/bubble

go 1.14
```

`bubble/dao/mysql.go`内容如下：

```
package dao

import "fmt"

func New(){
    fmt.Println("mypackage.New")
}
```

`bubble/main.go`内容如下：

```
package main

import (
    "fmt"
    "github.com/q1mi/bubble/dao"
)
func main() {
    dao.New()
    fmt.Println("main")
}
```

#### 不在同一个项目下

##### 目录结构

```
├── moduledemo
│   ├── go.mod
│   └── main.go
└── mypackage
    ├── go.mod
    └── mypackage.go
```

##### 导入包

这个时候，`mypackage`也需要进行module初始化，即拥有一个属于自己的`go.mod`文件，内容如下：

```
module mypackage

go 1.14
```

然后我们在`moduledemo/main.go`中按如下方式导入：

```
import (
    "fmt"
    "mypackage"
)
func main() {
    mypackage.New()
    fmt.Println("main")
}
```

因为这两个包不在同一个项目路径下，你想要导入本地包，并且这些包也没有发布到远程的github或其他代码仓库地址。这个时候我们就需要在`go.mod`文件中使用`replace`指令。

在调用方也就是`packagedemo/go.mod`中按如下方式指定使用**相对路径**来寻找`mypackage`这个包。

```
module moduledemo

go 1.14


require "mypackage" v0.0.0
replace "mypackage" => "../mypackage"
```

#### 举个例子

最后我们再举个例子巩固下上面的内容。

我们现在有文件目录结构如下：

```
├── p1
│   ├── go.mod
│   └── main.go
└── p2
    ├── go.mod
    └── p2.go
```

`p1/main.go`中想要导入`p2.go`中定义的函数。

`p2/go.mod`内容如下：

```
module liwenzhou.com/q1mi/p2

go 1.14
```

`p1/main.go`中按如下方式导入

```
import (
    "fmt"
    "liwenzhou.com/q1mi/p2"
)
func main() {
    p2.New()
    fmt.Println("main")
}
```

因为我并没有把`liwenzhou.com/q1mi/p2`这个包上传到`liwenzhou.com`这个网站，我们只是想导入本地的包，这个时候就需要用到`replace`这个指令了。

`p1/go.mod`内容如下：

```
module github.com/q1mi/p1

go 1.14


require "liwenzhou.com/q1mi/p2" v0.0.0
replace "liwenzhou.com/q1mi/p2" => "../p2"
```

此时，我们就可以正常编译`p1`这个项目了。





## 指令

### 更换默认代理网址

cmd: ```go env -w GOPROXY=https://goproxy.cn,direct```

### 指定所用的模块在本地的位置

cmd: ```go mod edit -replace example.com/greetings=../greetings```

设置GOBIN（install之后的位置）```go env -w GOBIN=/somewhere/else/bin```
恢复前一个设定```go env -u GOBIN```

### 模块

初始化：```go mod init example.com/greetings```

为了调用本地的模块，修改模块路径（否则会去仓库下）：  
```go mod edit -replace example.com/greetings=../greetings```

同步模块的依赖：
```go mod tidy```

### 测试

测试当前目录下以_test.go结尾的文件中的以Test开头的函数    
```go test```