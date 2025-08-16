---
title: "6.3 kiteX概览"
date: 2025-03-23T11:54:00Z
draft: false
weight: 6003
---

# 6.3 kiteX概览

# **一、Kitex架构设计**

**Kitex [kaɪt’eks] （kite*****n.***** 风筝）**：字节跳动开源的 Go 微服务 RPC 框架，具有**高性能**、**强可扩展**的特点，在字节内部已广泛使用。如果对微服务性能有要求，又希望定制扩展融入企业内部的治理体系，Kitex 会是一个不错的选择。

[https://github.com/cloudwego/kitex/blob/develop/README_cn.md](https://github.com/cloudwego/kitex/blob/develop/README_cn.md)

- 云原生大图：[https://landscape.cncf.io/](https://landscape.cncf.io/)
![](/images/1bf24637-29b5-8036-8428-c988f31a86eb/image_1fb24637-29b5-8045-a8ed-fb8ad4ed8657.jpg)

**特点：**

1. **高性能：**使用自研的高性能网络库 [Netpoll](https://github.com/cloudwego/netpoll)，性能相较 go net 具有显著优势。
1. **扩展性：**提供了较多的扩展接口以及默认扩展实现，使用者也可以根据需要自行定制扩展，具体见下面的框架扩展。
1. **多消息协议：**RPC 消息协议默认支持 **Thrift**、**Kitex Protobuf**、**gRPC**。Thrift 支持 Buffered 和 Framed 二进制协议，与支持原生 Thrift 协议的多语言框架都能互通； Kitex Protobuf 是 Kitex 自定义的 Protobuf 消息协议，协议格式类似 Thrift；gRPC 是对 gRPC 消息协议的支持，可以与 gRPC 互通。除此之外，使用者也可以扩展自己的消息协议，目前社区也提供了 Dubbo 协议的支持，可以与 Dubbo 互通。
1. **多传输协议：**传输协议封装消息协议进行 RPC 互通，传输协议可以额外透传元信息，用于服务治理，Kitex 支持的传输协议有 **TTHeader**、**HTTP2**。TTHeader 可以和 Thrift、Kitex Protobuf 结合使用；HTTP2 目前主要是结合 gRPC 协议使用，后续也会支持 Thrift。
1. **多种消息类型：**支持 **PingPong**、**Oneway**、**双向 Streaming**。其中 Oneway 目前只对 Thrift 协议支持，双向 Streaming 只对 gRPC 支持，后续会考虑支持 Thrift 的双向 Streaming。
1. **服务治理：**支持服务注册/发现、负载均衡、熔断、限流、重试、监控、链路跟踪、日志、诊断等服务治理模块，大部分均已提供默认扩展，使用者可选择集成。
1. **代码生成：**Kitex 内置代码生成工具，可支持生成 **Thrift**、**Protobuf** 以及脚手架代码。




# 二、学习路径

以Gomall系列教程入手：[Go 微服务电商项目实战](https://www.bilibili.com/video/BV1Nx4y1Y75p?spm_id_from=333.788.videopod.sections&vd_source=b79a821abfa98b8b24b71c819080ae85)

<!-- 列布局开始 -->

![](/images/1bf24637-29b5-8036-8428-c988f31a86eb/image_1fd24637-29b5-80ba-9797-f9592b2bbfa1.jpg)






---

![](/images/1bf24637-29b5-8036-8428-c988f31a86eb/image_1fd24637-29b5-802d-8729-edd4373bfb5f.jpg)

<!-- 列布局结束 -->



开发环境：

```shell
mkdir hello_world
cd hello_world 
go mod init github.com/cloudwego/biz-demo/gomall/hello_word
go get -u github.com/cloudwego/hertz
touch main.go
go mod tidy
go run main.go
```



[http://localhost:8888/hello](http://localhost:8888/hello)

![](/images/1bf24637-29b5-8036-8428-c988f31a86eb/image_1fd24637-29b5-8058-b297-caa532ebc095.jpg)

# **三、kiteX-开发流程: 单个服务**

## Kitex 开发环境

**环境：**

1. go最新版本
1. 设置国内代理：`go env -w GOPROXY=https://goproxy.cn`
1. 在安装代码生成工具前，确保 `GOPATH` 环境变量已经被正确地定义（例如 `export GOPATH=~/go`）并且将`$GOPATH/bin`添加到 `PATH` 环境变量之中（例如 `export PATH=$GOPATH/bin:$PATH`）；请勿将 `GOPATH` 设置为当前用户没有读写权限的目录。
### **代码生成工具**

Kitex 中使用到的代码生成工具包括 IDL 编译器与 kitex tool。了解更多有关代码生成工具的内容，参见[代码生成](https://www.cloudwego.io/zh/docs/kitex/tutorials/code-gen/)

### **IDL 编译器 (Thrift 可跳过)**

IDL 编译器能够解析 IDL 并生成对应的序列化和反序列化代码，Kitex 支持 Thrift 和 protobuf 这两种 IDL，这两种 IDL 的解析分别依赖于 thriftgo 与 protoc。

- thrift 依赖的 thriftgo 在 安装 kitex 工具的时候会安装（最新版本会去掉对 thriftgo 的依赖，[Kitex 去 Apache Thrift](https://www.cloudwego.io/zh/docs/kitex/best-practice/remove_apache_codec/)），**无需手动安装**
- protobuf 编译器安装可见 [protoc](https://github.com/protocolbuffers/protobuf/releases)
### **kitex tool**

1. 安装**kitex**：kitex 是 Kitex 框架提供的用于生成代码的一个命令行工具。目前，kitex 支持 thrift 和 protobuf 的 IDL，并支持生成一个服务端项目的骨架。kitex 的使用需要依赖于 IDL 编译器确保你已经完成 IDL 编译器的安装。
```shell
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
```

1. 查看版本： `kitex --version`  v0.12.3


## 运行 hello，world

1. git clone [https://github.com/cloudwego/kitex-examples.git](https://github.com/cloudwego/kitex-examples.git)
1. `cd kitex-examples/hello`
1. 运行服务端代码：`go run .`
1. 另启一个终端运行客户端代码：`go run ./client`
## Step1: thrift IDL中添加一个接口定义

**新增方法：**你已经完成现有示例代码的运行，接下来添加你自己实现的方法并运行起来吧。

打开 `hello` 目录下的 `hello.thrift` 文件，你会看到以下内容：

```c++
namespace go api

struct Request {
        1: string message
}

struct Response {
        1: string message
}

service Hello {
    Response echo(1: Request req)
}
```

现在让我们为新方法分别定义一个新的请求和响应，`AddRequest` 和 `AddResponse`，并在 `service Hello` 中增加 `add` 方法：

```c++
namespace go api

struct Request {
    1: string message
}

struct Response {
    1: string message
}

struct AddRequest {
  	1: i64 first
  	2: i64 second
}

struct AddResponse {
  	1: i64 sum
}

service Hello {
    Response echo(1: Request req)
    AddResponse add(1: AddRequest req)
}
```

## **Step2：**`kitex`**生成代码**

运行如下命令后，`kitex` 工具根据 `hello.thrift` 内容自动更新代码文件。  /ˈmɑːdʒuːl/  *n.* 单元, 单位，飞船独立舱

```plain text
kitex -module "github.com/cloudwego/kitex-examples" -service a.b.c hello.thrift
```

执行完上述命令后，`kitex` 工具将更新下述文件

1. 更新 `./handler.go`，在里面增加一个 `Add` 方法的基本实现
1. 更新 `./kitex_gen`，里面有框架运行所必须的代码文件
## **Step3：更新**`handler`的业**务逻辑**

上述步骤完成后，`./handler.go` 中会自动补全一个 `Add` 方法的基本实现，类似如下代码：

```go
// Add implements the HelloImpl interface.
func (s *HelloImpl) Add(ctx context.Context, req *api.AddRequest) (resp *api.AddResponse, err error) {
        // TODO: Your code here...
        return
}
```

这个方法对应我们在 `hello.thrift` 中新增的 `Add` 方法，我们要做的就是增加我们想要的业务逻辑代码。例如返回请求参数相加后的结果：

```go
// Add implements the HelloImpl interface.
func (s *HelloImpl) Add(ctx context.Context, req *api.AddRequest) (resp *api.AddResponse, err error) {
        // TODO: Your code here...
        resp = &api.AddResponse{Sum: req.First + req.Second}
        return
}
```

## **Step4：新增Client调用**

服务端已经有了 `Add` 方法的处理，现在让我们在客户端增加对 `Add` 方法的调用。

在 `./client/main.go` 中你会看到类似如下的 `for` 循环：

```go
for {
        req := &api.Request{Message: "my request"}
        resp, err := client.Echo(context.Background(), req)
        if err != nil {
                log.Fatal(err)
        }
        log.Println(resp)
        time.Sleep(time.Second)
}
```

现在让我们在里面增加 `Add` 方法的调用：

```go
for {
        req := &api.Request{Message: "my request"}
        resp, err := client.Echo(context.Background(), req)
        if err != nil {
                log.Fatal(err)
        }
        log.Println(resp)
        time.Sleep(time.Second)
        // 添加代码
        addReq := &api.AddRequest{First: 512, Second: 512}
        addResp, err := client.Add(context.Background(), addReq)
        if err != nil {
                log.Fatal(err)
        }
        log.Println(addResp)
        time.Sleep(time.Second)
}
```

## **Step5：编译运行、接口测试**

按照第一次运行示例代码的方法，分别重新运行服务端与客户端代码，看到类似如下输出，代表运行成功：

2025/03/15 21:33:46 Response({Message:my request})
2025/03/15 21:33:47 AddResponse({Sum:1024})
2025/03/15 21:33:48 Response({Message:my request})
2025/03/15 21:33:49 AddResponse({Sum:1024})









