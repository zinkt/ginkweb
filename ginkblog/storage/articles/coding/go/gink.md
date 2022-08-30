### ServeHTTP接口

要使用net/http的`ListenAndServe(addr string, h Handler)`方法，需要实现ServeHTTP接口



### 理解顺序

1. 上下文context
2. 路由
3. 分组控制
4. 中间件
5. 模板