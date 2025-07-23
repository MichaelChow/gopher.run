---
title: "5.5 hertZ概览"
date: 2025-03-11T14:41:00Z
draft: false
weight: 5005
---

# 5.5 hertZ概览



## 一、hert架构设计

官方文档：[https://www.cloudwego.io/zh/docs/hertz/](https://www.cloudwego.io/zh/docs/hertz/)

开源地址：[https://github.com/cloudwego/hertz/blob/develop/README_cn.md](https://github.com/cloudwego/hertz/blob/develop/README_cn.md)



Hertz [həːts，即物理学家赫兹] 是字节跳动开源的 Go 微服务 HTTP 框架，在设计之初参考了其他开源框架[gin](https://github.com/gin-gonic/gin)（公司内部二次开发的ginx)、 [fasthttp](https://github.com/valyala/fasthttp)、[echo](https://github.com/labstack/echo) 的优势， 并结合字节跳动内部的需求，使其具有高易用性、高性能、高扩展性等特点，目前在字节跳动内部已广泛使用。 如今越来越多的微服务选择使用 Golang，如果对微服务性能有要求，又希望框架能够充分满足内部的可定制化需求，Hertz 会是一个不错的选择。

![](/images/1b324637-29b5-802f-aee7-f3ba0f0e2313/image_1b424637-29b5-8073-ac2e-e1a617fcba27.jpg)



### **框架特点**

- **高易用性**：在开发过程中，快速写出来正确的代码往往是更重要的。因此，在 Hertz 在迭代过程中，积极听取用户意见，持续打磨框架，希望为用户提供一个更好的使用体验，帮助用户更快的写出正确的代码。
- **高性能**：Hertz 默认使用自研的高性能网络库 Netpoll，在一些特殊场景相较于 go net，Hertz 在 QPS、时延上均具有一定优势。关于性能数据，可参考下图 Echo 数据。
    四个框架的对比:
    ![](/images/1b324637-29b5-802f-aee7-f3ba0f0e2313/image_1b424637-29b5-802d-8e32-c8348cae4880.jpg)
    三个框架的对比:
    ![](/images/1b324637-29b5-802f-aee7-f3ba0f0e2313/image_1b424637-29b5-8045-ad48-f4aecf0887fa.jpg)
    关于详细的性能数据，可参考 [https://github.com/cloudwego/hertz-benchmark](https://github.com/cloudwego/hertz-benchmark)。
- **高扩展性：**Hertz 采用了分层设计，提供了较多的接口以及默认的扩展实现，用户也可以自行扩展。同时得益于框架的分层设计，框架的扩展性也会大很多。目前仅将稳定的能力开源给社区，更多的规划参考 [RoadMap](https://github.com/cloudwego/hertz/blob/main/ROADMAP.md)。
- **多协议支持：**Hertz 框架原生提供 HTTP1.1、ALPN 协议支持。除此之外，由于分层设计，Hertz 甚至支持自定义构建协议解析逻辑，以满足协议层扩展的任意需求。
- **网络层切换能力：**Hertz 实现了 Netpoll 和 Golang 原生网络库 间按需切换能力，用户可以针对不同的场景选择合适的网络库，同时也支持以插件的方式为 Hertz 扩展网络库实现。


### hz工具

代码结构

```go
.
├── biz                                // **business 层，存放业务逻辑相关流程**
│   ├── handler                        // **存放 handler 文件**
│   │   ├── hello                      // hello/example **对应 thrift idl 中定义的 namespace**；而对于 protobuf idl，则是对应 go_package 的最后一级
│   │   │   └── example
│   │   │       └── hello_service.go   // handler 文件，用户在该文件里实现 IDL service 定义的方法，update 时会查找当前文件已有的 handler 并在尾部追加新的 handler
│   │   └── ping.go                    // 默认携带的 ping handler，用于生成代码快速调试，无其他特殊含义
│   ├── model                          // idl 内容相关的生成代码
│   │   └── hello                      // hello/example 对应 thrift idl 中定义的 namespace；而对于 protobuf idl，则是对应 go_package
│   │       └── example
│   │           └── hello.go           // thriftgo 的产物，包含 hello.thrift 定义的内容的 go 代码，update 时会重新生成
│   └── router                         // idl 中定义的路由相关生成代码
│       ├── hello                      // hello/example 对应 thrift idl 中定义的 namespace；而对于 protobuf idl，则是对应 go_package 的最后一级
│       │   └── example
│       │       ├── hello.go           // hz 为 hello.thrift 中定义的路由生成的路由注册代码；每次 update 相关 idl 会重新生成该文件
│       │       └── middleware.go      // 默认中间件函数，hz 为每一个生成的路由组都默认加了一个中间件；update 时会查找当前文件已有的 middleware 在尾部追加新的 middleware
│       └── register.go                // 调用注册每一个 idl 文件中的路由定义；当有新的 idl 加入，在更新的时候会自动插入其路由注册的调用；勿动
├── go.mod                             // go.mod 文件，如不在命令行指定，则默认使用相对于 GOPATH 的相对路径作为 module 名
├── idl                                // 用户定义的 idl，位置可任意
│   └── hello.thrift
├── main.go                            // 程序入口
├── router.go                          // 用户自定义除 idl 外的路由方法
├── router_gen.go                      // hz 生成的路由注册代码，用于调用用户自定义的路由以及 hz 生成的路由
├── .hz                                // hz 创建代码标志，无需改动
├── build.sh                           // 程序编译脚本，Windows 下默认不生成，可直接使用 go build 命令编译程序
├── script
│   └── bootstrap.sh                   // 程序运行脚本，Windows 下默认不生成，可直接运行 main.go
└── .gitignore
```



**hz client:**

```go
.
├── biz                                  // business 层，存放业务逻辑相关流程
│   └── model                            // idl 内容相关的生成代码
│       └── hello                        // hello/example 对应 thrift idl 中定义的 namespace；而对于 protobuf idl，则是对应 go_package
│           └── example
│               ├── hello.go             // thriftgo 的产物，包含 hello.thrift 定义的内容的 go 代码，update 时会重新生成
│               └── hello_service        // 基于 idl 生成的类似 RPC 形式的 http 请求一键调用代码，可与 hz 生成的 server 代码直接互通
│                   ├── hello_service.go
│                   └── hertz_client.go
│
├── go.mod                               // go.mod 文件，如不在命令行指定，则默认使用相对于 GOPATH 的相对路径作为 module 名
└── idl                                  // 用户定义的 idl，位置可任意
    └── hello.thrift
```











安装环境：

1. GO安装最新版本；
1. 确保 `GOPATH` 环境变量已经被正确地定义（例如 `export GOPATH=~/go`）并且将 `$GOPATH/bin` 添加到 `PATH` 环境变量之中（例如 `export PATH=$GOPATH/bin:$PATH`）；请勿将 `GOPATH` 设置为当前用户没有读写权限的目录。
1. 设置国内代理：go env -w GOPROXY=https://goproxy.cn
1. goland 安装 Thrift Support 插件； 
1. 安装命令行工具 hz（**代码自动生成工具**）：
    1. hz 是 Hertz 框架提供的一个用于生成代码的命令行工具，可以用于生成 Hertz 项目的脚手架。
    1. 安装 hz：`go install github.com/cloudwego/hertz/cmd/hz@latest`
    1. 更多 hz 使用方法可参考: [hz](https://www.cloudwego.io/zh/docs/hertz/tutorials/toolkit/)
## 二、hertZ-开发流程

### Step1：编写thrift IDL 接口定义

1. 创建项目文件夹： `mkdir hertz_demo`  `cd hertz_demo`
1. 创建`hello.thrift`
    1. `mkdir idl`； 
    1. touch `idl/hello.thrift`
    ```go
    // idl/hello.thrift
    namespace go hello.example
    struct HelloReq {
        1: string Name (api.query="name"); // 添加 api 注解为方便进行参数绑定
    }
    struct HelloResp {
        1: string RespBody;
    }
    service HelloService {
        HelloResp HelloMethod(1: HelloReq request) (api.get="/hello");
    }
    ```
### Step2：hz生成代码

1. 生成代码: 详细参考[这里](https://www.cloudwego.io/zh/docs/hertz/tutorials/toolkit/usage/)。 执行完毕后, 会在当前目录下生成 Hertz 项目的脚手架, 自带一个 `ping` 接口用于测试。
    ```go
    # 在 GOPATH 外执行，需要指定 go mod 名
    hz new -module hertz/demo -idl idl/hello.thrift
    # 整理 & 拉取依赖
    go mod tidy
    ```
### Step3：修改handler的业务逻辑

```go
// handler path: biz/handler/hello/example/hello_service.go
// 其中 "hello/example" 是 thrift idl 的 namespace
// "hello_service.go" 是 thrift idl 中 service 的名字，所有 service 定义的方法都会生成在这个文件中

// HelloMethod .
// @router /hello [GET]
func HelloMethod(ctx context.Context, c *app.RequestContext) {
        var err error
        var req example.HelloReq
        err = c.BindAndValidate(&req)
        if err != nil {
                c.String(400, err.Error())
                return
        }

        resp := new(example.HelloResp)

        // 你可以修改整个函数的逻辑，而不仅仅局限于当前模板
        resp.RespBody = "hello," + req.Name // 添加的逻辑

        c.JSON(200, resp)
}
```

### Step4：编译运行、接口测试

1. 编译项目： `go build`
1. 运行项目并测试： `./demo`
    ```shell
    2022/05/17 21:47:09.626332 engine.go:567: [Debug] HERTZ: Method=GET    absolutePath=/ping   --> handlerName=main.main.func1 (num=2 handlers)
    2022/05/17 21:47:09.629874 transport.go:84: [Info] HERTZ: HTTP server listening on address=[::]:8888
    ```
1. 接下来，我们可以对接口进行测试：
    ```shell
    curl http://127.0.0.1:8888/ping
    {"message":"pong"}
    curl "http://127.0.0.1:8888/hello?name=chow"
    {"RespBody":"Hello chow!"}
    ```
### Step5：迭代需求 （相同）1. 更新 thrift IDL 接口定义

```go
// idl/hello.thrift
namespace go hello.example

struct HelloReq {
    1: string Name (api.query="name");
}

struct HelloResp {
    1: string RespBody;
}

struct OtherReq {
    1: string Other (api.body="other");
}

struct OtherResp {
    1: string Resp;
}


service HelloService {
    HelloResp HelloMethod(1: HelloReq request) (api.get="/hello");
    OtherResp OtherMethod(1: OtherReq request) (api.post="/other");
}

service NewService {
    HelloResp NewMethod(1: HelloReq request) (api.get="/new");
}
```

2. hz生成代码

1. 切换到执行 new 命令的目录，更新修改后的 thrift idl:  `hz update -idl idl/hello.thrift`
    1. 在 biz/handler/hello/example/hello_service.go 新增了新的方法
    1. 在 biz/handler/hello/example 下新增了文件 new_service.go 以及对应的 “NewMethod” 方法
3. 修改handler的业务逻辑

1. 修改 “OtherMethod” 接口的业务逻辑：
    ```go
    // OtherMethod .
    // @router /other [POST]
    func OtherMethod(ctx context.Context, c *app.RequestContext) {
         var err error
         // example.OtherReq 对应的 model 文件也会重新生成
         var req example.OtherReq
         err = c.BindAndValidate(&req)
         if err != nil {
             c.String(400, err.Error())
             return
         }
         resp := new(example.OtherResp)
         // 增加的逻辑
         resp.Resp = "Other method: " + req.Other
         c.JSON(200, resp)
    }
    ```


4. 编译运行、接口测试



1. 编译项目：`go build`
1. 运行：./demo
1. 接口测试：
    ```go
    curl --location --request POST 'http://127.0.0.1:8888/other' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "Other": "other method"
    }'
    返回: {"Resp":"Other method: other method"}
    ```


path不存在时，返回：`404 page not found`

![](/images/1b324637-29b5-802f-aee7-f3ba0f0e2313/image_1b724637-29b5-8037-b94e-dde40dba5e48.jpg)

