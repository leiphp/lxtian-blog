# LXTIAN-BLOG 雷小天博客

## 技术栈

1. 开发语言：**[golang](https://golang.google.cn/) ^go1.22.6**
2. 编程框架：**[go-zero](https://go-zero.dev/) ^1.7.2**
3. 服务发现：**[etcd](https://etcd.io/) ^3.5.10**

## 运行

### 启动开发环境

```shell
cp -n .env.example .env || true
docker compose up -d
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