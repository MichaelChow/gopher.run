---
title: "5.3 kiteX概览"
date: 2025-03-23T11:54:00Z
draft: false
weight: 5003
---

# 5.3 kiteX概览

# **一、Kitex架构设计**

**Kitex [kaɪt’eks] （kite*****n.***** 风筝）**：字节跳动开源的 Go 微服务 RPC 框架，具有**高性能**、**强可扩展**的特点，在字节内部已广泛使用。如果对微服务性能有要求，又希望定制扩展融入企业内部的治理体系，Kitex 会是一个不错的选择。

[https://github.com/cloudwego/kitex/blob/develop/README_cn.md](https://github.com/cloudwego/kitex/blob/develop/README_cn.md)

- 云原生大图：[https://landscape.cncf.io/](https://landscape.cncf.io/)
![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/d0729b68-cbc3-4863-b6f3-8a8b7c24cebd/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4663EUZYGTC%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005336Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJHMEUCIAjeWKeMNLOfdOKxbnYNJ7KN7SkNX5Niw%2BNM8qsRwf2YAiEA3BU%2FWHNrek5%2FXBuWZDEqTnwfN1ReY1e5leCq9bhZqRIqiAQImf%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FARAAGgw2Mzc0MjMxODM4MDUiDHrd8%2F7aRT2isYMQHCrcAx8pTLsvm%2BM1x%2Bj1nH0lNOvQcW6iE5tzFGCoRglFbaZiy64fyO%2BRz%2Bd0CiUVt1GI5E9CbXYjNTnVY7bW3h6FDBABn6b%2BsHcaBlOpzllnnb5icahRsrBFWfnHuOC5rv%2BuTmYnq0xpVTfChtzKLUKoFZqtKKjYSpQji75ZiBEu2D%2BJe1kwZmZar490Ie01734i%2FtXaZXzefuLIOwL58lPvih9PEb%2FxbhhptRmLL0%2BQL79q1FsBcQyD8CfyE%2FX335H3VS%2F3IstLoRr%2F4AiGCyL%2FzStveFWp0AkvdgswUMzombenMkzX3PwwlhR%2Bu10zvF%2FgoH82SjBELsTS2NPVV6lRFzSdNMtt1S9JY%2Ff%2FdP7AEffE80ru7MrUq2sxEY%2BgSqLnmMtH4U2wAiBT74GXP%2FqY2fuMEMV12fCI8XajGQoS%2BbZno2JuLrAmrtYXgtDmueP%2FxKQwof76Z%2FR1hlvJWdcHRaz7SQCa8Eb7XtgrmQu24cQ9DQ5YtfZzUcyorWH7YlW9uHmG29j7a69DkqxKyGyJ021LsMPM2NQgT6KEziCWVNOU0lAvKdxFEKNob3vBTlYuIhvuupvIXBRIHaoE1uNix1ljaZ4YJKuSDO9o9HDZ%2F3j0u68%2B4xTS9K9ktpsBMK%2B668MGOqUBDM8dYgEBdLZZCCf8KpQt1KI1gkxE71pUTBaaUgWIhtuTOBmEV4zB4CHTkAt6Ep9J9Fg5CeQi947X1X5nejDJCArum83GpZ%2FjwKmPC8gQKRHezbllZWaEa30OXY2QzDuE5D%2BYhc4KOzgScsRWNxaf30E4oRvbXw3gEllNWJ5Pp79qEVf%2B%2FPctiP3SjXpDM5ynPYsunDS1fyXz3WVNwaD3BVkIPafn&X-Amz-Signature=1ef647d4ca7b2bc0c4b78efbc0fa83f07c0be72b96afa098870007dce12da6cf&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)

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

![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/b354ddfb-a234-4386-9bbf-8ce56e4643ba/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB466UL7PYHDF%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005338Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJIMEYCIQDLqZdM34NXvrD8kRgbczNCGzTbx7V1XVBUfxHt28PzZQIhANpkihq0sH6b2nHJpBnIRga2TH%2FirUs2gEKthSo344qjKogECJn%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1IgxGogXfdYVcUgoOHUsq3APNZB4UNpTUkCrRJ%2BmlfFExnc%2BVwe144jS0pq29ALc57LttQ1kYlCMLHLWvL%2B59s1jt9YIGKr5qrjc%2Bzgx7TyFRqLi0%2BtsW%2BSjvrFQ3gUyJ%2BPepQBgoYhkbaHUZHTrHM3gDMY1D9Toj0%2B8JjyktZ2F7eeflInKwM1Yy5JquAt2OGe%2B6cGwN8VBVOzGK91E9EPY1%2FNt8O51Nl%2BjnE4PWtTPsWy56skb1Us2yDRh5oUQ%2B05qbmdt3lMJHuC4gz3UbXgY4xRq9TgsSCgiAfsWRey1ExfAFdcGe8AA3chLjkD1pQU9UPKv9%2FOiKKU7yI9604ghFmcOToswoM%2F6SWwM0vVlY1E0Nk8iMpYO4uALjQwQNiNzwj7F7NchHItEpuZMsmzNJbbze1R%2BW7QncuXlyXFDhZXssmTcQgxIQIg2QdeCNsN8lV4vzltW28J2zzctBhofhCBdfAihAUnq%2F4Fn7mnzjh3zCtbNfHFa3unAOpsGluCLg71zVHxqFWqma7YAh2cPFdRrcw0EtGaqx7ptTiQVnyP%2F%2Fbl4W8e0PSbVZfO5AZsq7azhVPklyqPfohVHzTHwHaRIKD67H6r3R7IqTXGgQk1agDOONs%2BpAZH%2FSvZS1Aht3kiGfQWkJcVoaWjCKuuvDBjqkASdxXy0XmQk3YeNDIQK3Ra9JX9hDvQ3WzXsY9UKaSMgHNGcCPoKJTA%2BZtVBsJ4Onc6rXzgk6%2BnityOlX9i8UvgNMcTpIxwNRWUAFi%2FePrK82Q8iJkPjXP0iYeJDeb2cj36DkjDe7UXfUv5NazvACJ8vaMs%2B%2FuGTAQca4aDnK6AkLjGB1b7X%2BCkT844wVKDkZ5maVkSyRC1V8%2B6fYM2QqQrQj%2Bvli&X-Amz-Signature=28e3d4273b61d2f2ce565b7410f9e2103d8147b624038857d67b452877a64337&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)






---

![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/d963c311-b70b-4e7b-9939-9dbd8887d010/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4666E4O2D4W%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005338Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJHMEUCIGR3SzaugHTJoP86CnxC3qOCN7bIj2H1mxN0yeSpBIwcAiEAm1fU75Y4L9MwhA%2FkPLa8I1B%2Bz5lOBBDT%2Beyg8u%2BsiycqiAQImf%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FARAAGgw2Mzc0MjMxODM4MDUiDMswd9Mi093DxQbXPircAzsAlUmazx6rxsAdRFO1Q0Iv6mKlgn9kZmdoLMypRgVdkT5wZkwmFpVJWULcOJi1PhMAXZK62RfQbMkYPbpFjjPw%2B4wUk6RdzCFxG6xR5uItnmFlmPHXA1LvtjMf37PGgtyU5HLUffB94VO91HN3fkk%2BP0OOzGK4Sr1IJtO5L17S76qLKKOWPDjOyH3PQzMA7rFssQ5FfsOo3o0sK1uOtwGTTNuh0HmhnCHHp9HTfqvpfLytIc%2BQZGXgv8jgGeBqyy%2BAshGnfJ3eD1b93Zw8uqDsFSgNg88unDP7K9fhNBE7kmfzHtMol2n1sCv5GK96n2ZFk2lRdgd209dYe2wwQcNb8i5qDyB04pePhyaFN5j5ESwtuoUu4yEvkRxK8PA0jMkZ7JGizjcsO7rjm%2FS%2B0J7UQM1qQWepBa6tKiV1eXTjRH6JZw6smSgTGhS1zmIeZP7O5XfPKZlJ2CQjt4OheACEqFCOFo5f4RJvITUtvtV0Lst7zmL6aqKnvLqVCUTngTn5yxJBZRjqp4KbuqDPuOqDY5Zn2znMVFSGBSZnhk3wXajXM6iYKQ9t2CTvbQCL75GGB3%2BC%2FXxWe2J3%2FCbYwsVdsmx2dn1w4xNW%2FbLATv8R3qnh814McblZN7qLMOG668MGOqUBuHoBW0sPH%2BlukRuHOPpxRhrU7b84q%2BKunjEb3XBH%2FU%2Fyd1SBQNqaiqzav09SEhETQX0E16S9J2vTO3g%2FL4bQ1TR%2B9cSgwWR8fTaZgNgHSrJK3i%2BHRoMQjjV%2Fo9hYdi%2F8W5tNEGLixFxPBZuyKTnwMc74ESP1EGslpMGXGrZXiMa1dFjUhioOUr4YKlmpeL2d08eQs7Wgut%2BpiuSIpMkywLVg1Vxf&X-Amz-Signature=46cb67e2536785f8b94af84c760c4c3eda345813ccde2ffd38b03a27e4fd8167&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)

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

![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/3e5d320e-0265-431f-a9b3-c2b141fb711d/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4663EUZYGTC%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005336Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJHMEUCIAjeWKeMNLOfdOKxbnYNJ7KN7SkNX5Niw%2BNM8qsRwf2YAiEA3BU%2FWHNrek5%2FXBuWZDEqTnwfN1ReY1e5leCq9bhZqRIqiAQImf%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FARAAGgw2Mzc0MjMxODM4MDUiDHrd8%2F7aRT2isYMQHCrcAx8pTLsvm%2BM1x%2Bj1nH0lNOvQcW6iE5tzFGCoRglFbaZiy64fyO%2BRz%2Bd0CiUVt1GI5E9CbXYjNTnVY7bW3h6FDBABn6b%2BsHcaBlOpzllnnb5icahRsrBFWfnHuOC5rv%2BuTmYnq0xpVTfChtzKLUKoFZqtKKjYSpQji75ZiBEu2D%2BJe1kwZmZar490Ie01734i%2FtXaZXzefuLIOwL58lPvih9PEb%2FxbhhptRmLL0%2BQL79q1FsBcQyD8CfyE%2FX335H3VS%2F3IstLoRr%2F4AiGCyL%2FzStveFWp0AkvdgswUMzombenMkzX3PwwlhR%2Bu10zvF%2FgoH82SjBELsTS2NPVV6lRFzSdNMtt1S9JY%2Ff%2FdP7AEffE80ru7MrUq2sxEY%2BgSqLnmMtH4U2wAiBT74GXP%2FqY2fuMEMV12fCI8XajGQoS%2BbZno2JuLrAmrtYXgtDmueP%2FxKQwof76Z%2FR1hlvJWdcHRaz7SQCa8Eb7XtgrmQu24cQ9DQ5YtfZzUcyorWH7YlW9uHmG29j7a69DkqxKyGyJ021LsMPM2NQgT6KEziCWVNOU0lAvKdxFEKNob3vBTlYuIhvuupvIXBRIHaoE1uNix1ljaZ4YJKuSDO9o9HDZ%2F3j0u68%2B4xTS9K9ktpsBMK%2B668MGOqUBDM8dYgEBdLZZCCf8KpQt1KI1gkxE71pUTBaaUgWIhtuTOBmEV4zB4CHTkAt6Ep9J9Fg5CeQi947X1X5nejDJCArum83GpZ%2FjwKmPC8gQKRHezbllZWaEa30OXY2QzDuE5D%2BYhc4KOzgScsRWNxaf30E4oRvbXw3gEllNWJ5Pp79qEVf%2B%2FPctiP3SjXpDM5ynPYsunDS1fyXz3WVNwaD3BVkIPafn&X-Amz-Signature=e5586445e7d55919dc423b7fb244b75023da3db17dfe6dc67c3275ebf43195bc&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)

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









