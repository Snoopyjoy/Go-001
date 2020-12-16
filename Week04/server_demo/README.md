## TASK
按照自己的构想，写一个项目满足基本的目录结构和工程，代码需要包含对数据层、业务层、API 注册，以及 main 函数对于服务的注册和启动，信号处理，使用 Wire 构建依赖。可以使用自己熟悉的框架。

## 准备工作
### 安装protoc
```shell
go get google.golang.org/protobuf/cmd/protoc-gen-go \
         google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

## TODO
- [x] API层，通过proto文件生成grpc服务接口  
- [x] 使用wire构建依赖  
- [x] 服务注册与启动
- [x] 信号处理
- [x] 分层架构

## 业务分层
### service层职责
1. 接口层，实现grpc server接口。
2. 应用层，同时实现biz业务逻辑、其他微服务接口的组装，DTO到领域对象的转换。
3. 消息队列消费
### biz层职责
业务逻辑，各业务逻辑应互相独立。如果有依赖放到service层去组装。
repo接口定义。
### data层职责
实现repo接口，从db获取数据并实现缓存控制。

## 实现
api实现一个profile服务，提供一个查询用户profile的接口。返回用户昵称和用户的金币数。
biz实现一个user业务提供用户基础数据，一个coin业务提供用户的金币数据。
data实现user、coin表的读取与缓存

## 目录结构
```
├── README.md
├── api
│   └── profile
│       ├── profile.pb.go
│       ├── profile.proto
│       └── profile_grpc.pb.go
├── cmd
│   └── profile
│       └── main.go
├── configs
│   └── profile
│       └── app_dev.yaml
├── go.mod
├── go.sum
├── internal
│   ├── biz
│   │   ├── coin.go
│   │   └── user.go
│   ├── data
│   │   ├── coin.go
│   │   └── user.go
│   └── service
│       ├── application.go
│       ├── service.go
│       ├── wire.go
│       └── wire_gen.go
├── pkg
│   ├── app
│   │   ├── app.go
│   │   ├── wire.go
│   │   └── wire_gen.go
│   ├── cache
│   │   └── cache.go
│   └── db
│       └── db.go
└── scripts
    └── gen_pb.sh
```

## 服务调试
进入`service_demo/cmd/profile`目录，执行代码`go run main.go`启动服务。  
ctrl+c优雅停止服务。  
### 使用grpc调试工具调试接口
安装grpcui
```shell
go get github.com/fullstorydev/grpcui/...
go install github.com/fullstorydev/grpcui/cmd/grpcui
```
调试服务
```shell
grpcui -plaintext 127.0.0.1:8080
```
在浏览器中打开grpcui输出的调试网址。  
在grpcui的UI中选择需要调试的接口进行调试。
