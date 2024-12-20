# LXTIAN-BLOG 雷小天博客

## 技术栈

1. 开发语言：**[golang](https://golang.google.cn/) ^go1.22.6**
2. 编程框架：**[go-zero](https://go-zero.dev/) ^1.7.2**
3. 服务发现：**[etcd](https://etcd.io/) ^3.5.10**

## 运行

### 启动开发环境

```shell
# 配置环境变量
cp -n .env.example .env || true
# 查看 .env 文件的格式是否正确
cat -A .env
# 重新加载环境变量 确保 .env 文件被正确解析
docker-compose config
# 如果发现问题，可以手动修复或者运行 
dos2unix .env
# 新版本禁用 Docker BuildKit调试
DOCKER_BUILDKIT=0 docker-compose up -d --build
# 运行
docker compose up -d

# 配置中心
/gateway
Telemetry:
  Name: web-rpc
  Endpoint: ""
  Batcher: jaeger
  Sampler: 1.0
  
ShortLink:
  Url: ""
  Key: ""
  Domain: ""
  Protocol: ""
  
  
/web
Telemetry:
  Name: web-rpc
  Endpoint: ""
  Batcher: jaeger
  Sampler: 1.0
```

### 安装/更新项目依赖

```shell
go mod tidy
go mod vendor
```

### 停止开发环境

```shell
docker compose down
```

### 运行单元测试

```shell
编写测试单元文件hello-test.go
package testunit

import (
	"fmt"
	"testing"
)
func TestHello(t *testing.T) {
	input := "Hello, 世界" 
    fmt.Println("say:", input)
}
```

## 开发

### 工具安装
#### goctl工具安装

* 查看go版本  
```
go version
```
* 安装
1. 如果 go 版本在 1.16 以前，则使用如下命令安装：  
```GO111MODULE=on go get -u github.com/zeromicro/go-zero/tools/goctl@latest ```
2. 如果 go 版本在 1.16 及以后，则使用如下命令安装：  
```go install github.com/zeromicro/go-zero/tools/goctl@latest```
* 验证  
```goctl --version```


### 新增服务
##### api服务生成
```shell
# 创建工作空间并进入该目录
$ mkdir -p ~/lxtian-blog/api && cd ~/lxtian-blog/api
# 执行指令生成 demo 服务
$ goctl api new demo
Done.
```

##### grpc服务生成
```shell
# 创建工作空间并进入该目录
$ mkdir -p ~/lxtian-blog/rpc && cd ~/lxtian-blog/rpc
# 执行指令生成 demo 服务
$ goctl rpc new demo
Done.
```

### CLI代码生产
#### api代码生成
```shell
# 切换到api服务目录下执行
cd ~/lxtian-blog/gateway
$ goctl api go -api gateway.api -dir .
Done.
```

#### grpc代码生成
```shell
# 切换到rpc服务目录下执行
cd ~/lxtian-blog/rpc/member
$ goctl rpc protoc member.proto --go_out=. --go-grpc_out=. --zrpc_out=. -m
Done.
```

#### mysql model代码生成
```shell
# 这将为 article 表生成相应的 Go model 代码，并放在 ./model 目录中。
cd ~/lxtian-blog/rpc/web
$ goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/your_database" -table="article" -dir="./model"
# -table="article,user"：指定要生成的多个表，用逗号分隔表名。
$ goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/your_database" -table="article,user" -dir="./model"
# -table="article_*"：使用通配符匹配表名，这里会生成所有以 article_ 开头的表的 model
goctl model mysql datasource -url="root:password@tcp(127.0.0.1:3306)/your_database" -table="article_*" -dir="./model"
Done.
```

#### mongo model代码生成
```shell
# 这将为 article 表生成相应的 Go model 代码，并放在 ./model 目录中。
cd ~/lxtian-blog/rpc/web
# 生成文档名称为 Article 的 mongo 代码
$ goctl model mongo --type article --dir .
Done.
```

### windows注入环境变量

```shell
#etcd环境变量,多个逗号隔开
$env:ETCD_HOSTS="127.0.0.1:2379"

#mysql环境变量
$env:DB_HOST="127.0.0.1"
$env:DB_PORT="3306"
$env:DB_DATABASE="lxtblog"
$env:DB_USERNAME="root"
$env:DB_PASSWORD="root"

#mongodb环境变量
$env:MONGODB_HOST="127.0.0.1"
$env:MONGODB_PORT="27017"
$env:MONGODB_DATABASE="lxtblog"
$env:MONGODB_USERNAME=""
$env:MONGODB_PASSWORD=""

#redis环境变量
$env:REDIS_HOST="127.0.0.1:6379"
$env:REDIS_TYPE="node"
$env:REDIS_PASS=""
$env:REDIS_TLS=false
```


### 数据库封装操作
#### mongodb案例
```shell
# mongodb获取文章内容
conn := model.NewArticleModel(l.svcCtx.MongoUri, l.svcCtx.Config.MongoDB.DATABASE, "txy_article")
contentId, ok := article["title"].(string)
if !ok {
    // 处理类型断言失败的情况
    return nil, errors.New("content_id is not a string")
}
res, err := conn.FindOne(l.ctx, contentId)
if err != nil {
    return nil, err
}
article["content"] = res.Content
```

#### redis案例
```shell
# redis实例
err = l.svcCtx.Rds.SetCtx(l.ctx, "key", "hello world")
if err != nil {
    logc.Error(l.ctx, err)
}
v, err := l.svcCtx.Rds.GetCtx(l.ctx, "key")
if err != nil {
    logc.Error(l.ctx, err)
}
fmt.Println(v)
```

### 内存缓存
#### 案例
```shell
// 测试缓存
l.svcCtx.Cache.Set("userInfo", map[string]interface{}{
    "id":   txyUser.Id,
    "name": "雷小天",
    "age":  18,
    "sex":  1,
})
v, exist := l.svcCtx.Cache.Get("userInfo")
if !exist {
    // deal with not exist
    return nil, errors.New("deal with not exist:数据不存在")
}
value, ok := v.(map[string]interface{})
if !ok {
    // deal with type error
    return nil, errors.New("deal with type error:数据类型错误")
}
fmt.Println("value:", value)
```