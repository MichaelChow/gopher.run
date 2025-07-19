---
title: "5.4 kiteX-开发流程: 多个服务"
date: 2025-05-22T13:28:00Z
draft: false
weight: 5004
---

# 5.4 kiteX-开发流程: 多个服务

商品服用调用库存服务查询库存，API 服务调用商品服务查询商品信息，对前端或用户暴露 HTTP 接口供查询商品信息。

本节中，我们使用 Kitex 完成了 RPC 服务端与客户端的开发，实现了 RPC 调用，开发流程总结如下：

1. 服务端:编写 IDL，使用 kitex 生成代码后填充服务业务逻辑即可运行
1. 客户端：使用与服务端相同的 IDL，使用 kitex 生成代码后，创建客户端示例，构造请求参数后即可发起调用
本例中，我们仅演示了 Kitex 的基本使用方法，Kitex 还为我们提供了各种微服务治理特性，你可以在[指南](https://www.cloudwego.io/zh/docs/kitex/tutorials/)中获取更多信息，或查看示例代码小节解锁更多高级用法。

## **Step1：编写thrift IDL接口定义**

1. **创建项目目录:**mkdir example_shop、 cd example_shop、`mkdir idl`、`cd idl`
1. 商品服务使用item.thrift定义接口，库存服务使用stock.thrift定义接口，同时创建 `base.thrift` 定义公共数据结构。
```c++
namespace go example.shop.base

struct BaseResp {
    1: string code
    2: string msg
}
```

item.thrift：

```c++
namespace go example.shop.item

include "base.thrift"

struct Item {
    1: i64 id
    2: string title
    3: string description
    4: i64 stock
}

struct GetItemReq {
    1: required i64 id
}

struct GetItemResp {
    1: Item item

    255: base.BaseResp baseResp
}

service ItemService{
    GetItemResp GetItem(1: GetItemReq req)
}
```

stock.thrift

```c++
namespace go example.shop.stock

include "base.thrift"

struct GetItemStockReq {
    1: required i64 item_id
}

struct GetItemStockResp {
    1: i64 stock

    255: base.BaseResp BaseResp
}

service StockService {
    GetItemStockResp GetItemStock(1:GetItemStockReq req)
}
```

## Step2: kitex**生成代码、拉取依赖**

两个IDL 定义了两个服务，执行两次 kitex 命令：

```shell
kitex -module example_shop idl/item.thrift
kitex -module example_shop idl/stock.thrift
```

生成的代码分两部分：

1. 结构体的编解码序列化代码，由 IDL 编译器生成；
1. 由 kitex 工具在前者产物上叠加，生成用于创建和发起 RPC 调用的桩代码。它们默认都在 `kitex_gen` 目录下。


上面生成的代码并不能直接运行，需要自己完成 `NewClient` （overpass）和 `NewServer` 的构建。kitex 命令行工具提供了 `-service` 参数能直接生成带有脚手架的代码，接下来让我们为商品服务和库存服务分别生成脚手架。



1. 首先为两个 RPC 服务分别单独创建目录：`mkdir -p rpc/item rpc/stock`
1. 再分别进入各自的目录中，执行如下命令生成代码：
```shell
// item 目录下执行
kitex -module example_shop -service example.shop.item -use example_shop/kitex_gen ../../idl/item.thrift

// stock 目录下执行
kitex -module example_shop -service example.shop.stock -use example_shop/kitex_gen ../../idl/stock.thrift
```

kitex 默认会将代码生成到执行命令的目录下，kitex 的命令中：

- `-module module name`：参数表明生成代码的 `go mod` 中的 `module name`，在本例中为 `example_shop`
- `-service service name`： 参数表明我们要生成脚手架代码，后面紧跟的 `example.shop.item` 或 `example.shop.stock` 为该服务的名字。
- `-use example_shop/kitex_gen`: 参数表示让 kitex 不生成 `kitex_gen` 目录，而使用该选项给出的 `import path`。在本例中因为第一次已经生成 `kitex_gen` 目录了，后面都可以复用。
- 最后一个参数则为该服务的 IDL 文件


生成代码后的项目结构如下：

```plain text
.
├── go.mod // **go module 文件**
├── go.sum
├── idl   // **示例 idl 存放的目录**
│   ├── base.thrift
│   ├── item.thrift
│   └── stock.thrift
├── kitex_gen
│   └── example
│       └── shop
│           ├── base
│           │   ├── base.go // **根据 IDL 生成的编解码文件，由 IDL 编译器生成**
│           │   ├── k-base.go // kitex 专用的一些拓展内容
│           │   └── k-consts.go
│           ├── item
│           │   ├── item.go // **根据 IDL 生成的编解码文件，由 IDL 编译器生成**
│           │   ├── itemservice //** kitex 封装代码主要在这里**
│           │   │   ├── client.go
│           │   │   ├── invoker.go
│           │   │   ├── itemservice.go
│           │   │   └── server.go
│           │   ├── k-consts.go
│           │   └── k-item.go // kitex 专用的一些拓展内容
│           └── stock
│               ├── k-consts.go
│               ├── k-stock.go // kitex 专用的一些拓展内容
│               ├── stock.go // **根据 IDL 生成的编解码文件，由 IDL 编译器生成**
│               └── stockservice // **kitex 封装代码主要在这里**
│                   ├── client.go
│                   ├── invoker.go
│                   ├── server.go
│                   └── stockservice.go
└── rpc
    ├── item
    │   ├── build.sh   // **用来编译的脚本，一般情况下不需要更改**
    │   ├── handler.go // **服务端的业务逻辑都放在这里，这也是我们需要更改和编写的文件**
    │   ├── kitex_info.yaml
    │   ├── main.go
    │   └── script
    │       └── bootstrap.sh
    └── stock
        ├── build.sh   // 用来编译项目的脚本，一般情况下不需要更改
        ├── handler.go // **服务端的业务逻辑都放在这里，这也是我们需要更改和编写的文件**
        ├── kitex_info.yaml
        ├── main.go    // **服务启动函数，一般在这里做一些资源初始化的工作，可以更改**
        └── script
            └── bootstrap.sh

```





1. 回到项目根目录cd `example_shop`。 使用 `go mod tidy` 命令拉取项目依赖
如果遇到类似如下两种报错：

```plain text
github.com/apache/thrift/lib/go/thrift: ambiguous import: found package github.com/apache/thrift/lib/go/thrift in multiple modules

github.com/cloudwego/kitex@v0.X.X/pkg/utils/thrift.go: not enough arguments in call to t.tProt.WriteMessageBegin
```

先执行一遍下述命令，再继续操作：

```plain text
go mod edit -droprequire=github.com/apache/thrift/lib/go/thrift
go mod edit -replace=github.com/apache/thrift=github.com/apache/thrift@v0.13.0
```

这是因为 thrift 官方在 0.14 版本对 thrift 接口做了 breaking change，导致生成代码不兼容。

若想要升级 kitex 版本，执行 `go get -v github.com/cloudwego/kitex@latest` 即可

## **Step3：更新**`handler`**商品服务的业务逻辑**

需要编写的服务端逻辑都在 `handler.go` 这个文件中。两个服务，对应了两个 `handler.go：`

商品服务的服务端逻辑 `rpc/item/handler.go`

注意：这里的 `GetItem` 函数就对应了我们之前在 `item.thrift` IDL 中定义的 `GetItem` 方法。对应自动生成的：basic/example_shop/kitex_gen/example/shop/item/item.go

注意：经过查询，阿里的HSF RPC框架，没有IDL的概念

```go
package main

import (
	"context"
	item "example_shop/kitex_gen/example/shop/item"
)

// ItemServiceImpl implements the last service interface defined in the IDL.
type ItemServiceImpl struct{}

// GetItem implements the ItemServiceImpl interface.
func (s *ItemServiceImpl) GetItem(ctx context.Context, req *item.GetItemReq) (resp *item.GetItemResp, err error) {
	// TODO: Your code here...
	return
}
```



注意到KiteX根据IDL接口定义自动生成的item.go，GetStock()方法，并不符合Go语言规范的Stock()

![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/c4a05b38-aa9f-4cf8-8ab2-143033db5a04/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4663MIDRUW3%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005341Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJIMEYCIQDtJMNb4o8mzx%2BtIYHkgiesScztn0CNJe6ISGerZ%2BRVPgIhAPhVCCuMq%2BvThJYPpgYRHTOUkJD4lORpvZTHZWm3mABOKogECJn%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1IgzzZ%2FNfadZa7nwf0fIq3ANer2Oh%2Fer3APiCuu9Ibc%2Fe%2BOcIyjf9uT0l%2FrRK%2F8lvkmVO%2BkW1Bv5IRmYxsQzjXNfPvq6IzCi2BYPFWt6MDm3ix%2BNhs8XoYl1a1Wa4YBvLD6HLj%2FRx0EUAsxKOwrFvWf2h9HEWWoV2ccUOdcIr%2Fa6%2Bw71x6%2Fyt0Hd5eVByx%2BtjRWZq2JDkhQDwDRRaSbJbTpR3ccDMhKVQim3skzSqI7fVKbRuQpKFYuEs1EqygszhAhVDdGG9Fmr36AHIt8OPeNYWzUwd1EMgcfssNudcefFuUaQVUW6MMoV9WcqnGZVOaZ1QKctfSj8cm%2B%2BSa6vwXAMmGsigf6C7EigaWhHsaJGFzavOMnnhTgTuBFEhN8hqqUWK66gfuaTzKEOabdgIxlmnV381IGnW0uEwu4okLh%2B2TJRQJdigwNVlGB7lB06MSn7lGz3xedDWeFj3S4GNqe00YwhtLFVDtMnY0DPFiOJcAqWOXLgcMlxWyoQckK3mqTzKnLLAoJ%2BUMXUXTQ8zghNhQ1atj9bLKrFKJy3Jmwe49abVTHw1tDXj2pG%2BmHPyrH2q1eOTRhJGSE1gbD3K4yRb32KdSfO5tRBfhBSG85BggoAiDcE0FkjS3vPAhaOetYpM3L6pDrAYQ3sfNTDXuuvDBjqkAUzb8oIQsKzbRTAO47zN6xu%2BoN6ASfDqUdLNIeGYv2xHf7pZj5TKZnvLS18PNyK%2FFECuEn4AeacPuztG4b3xXb8MVAj%2FKutIAcfQ90fZ90aqCnoRj%2BGcFLrMQFQfEpp44OAcqhYRY6bcdD4Vt93aUyKbZtFQx57nDTL9HFFjAPBnNyQIAhAmUB08U3IZ97RHfNVNgM299XjpfxIZCLZkP06b8zDH&X-Amz-Signature=9fbea095a5cdf711bc5869cde7ba4ed45a136a22071f5d9d7adf445602c9320d&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)



main.go：代码很简单，即使用 kitex 生成的代码创建一个 `server` 服务端，并调用其 `Run` 方法开始运行。通常使用 `main.go` 进行一些项目初始化，如加载配置等。

```go
package main

import (
    item "example_shop/kitex_gen/example/shop/item/itemservice"
    "log"
)

func main() {
    svr := item.NewServer(new(ItemServiceImpl))

    err := svr.Run()

    if err != nil {
       log.Println(err.Error())
    }
}
```

## **Step4：运行商品服务、接口测试**

kitex 也为我们生成了编译脚本，即 `build.sh` :

1. 定义了一个变量 `RUN_NAME`，用于指定生成的可执行文件的名称，值为我们在 IDL 中指定的 `namespace`。本例中为 `example.shop.item`
1. 创建 `output` 目录，此后的编译出的二进制文件放在 `output/bin` 下。同时将 `script` 目录下的项目启动脚本复制进去
1. 根据环境变量 `IS_SYSTEM_TEST_ENV` 的值判断生成普通可执行文件或测试可执行文件。值为 1 则代表使用 `go test -c` 生成测试文件，否则正常使用 `go build` 命令编译。
```go
#!/usr/bin/env bash
RUN_NAME="example.shop.item"

mkdir -p output/bin
cp script/* output/
chmod +x output/bootstrap.sh

if [ "$IS_SYSTEM_TEST_ENV" != "1" ]; then
    go build -o output/bin/${RUN_NAME}
else
    go test -c -covermode=set -o output/bin/${RUN_NAME} -coverpkg=./...
fi
```



直接执行 `sh build.sh` 即可编译项目。

编译成功后，生成 `output` 目录：

```plain text
output
├── bin // 存放二进制可执行文件
│   └── example.shop.item
└── bootstrap.sh // 运行文件的脚本

```

执行 `sh output/bootstrap.sh` 即可启动编译后的二进制文件。

输出类似以下命令，代表运行成功：

```plain text
2024/01/19 22:12:18.758245 server.go:83: [Info] KITEX: server listen at addr=[::]:8888
```

在上面的日志输出中，`addr=[::]:8888` 代表我们的服务运行在本地的 8888 端口，此参数可以在创建 `server` 时传入 `option` 配置来修改，更多服务端配置见 [Server Option](https://www.cloudwego.io/zh/docs/kitex/tutorials/options/server_options/)。

## **Step5：编写API服务，**对外暴露 HTTP 接口

编写 API 服务用于调用刚刚运行起来的商品服务，并对外暴露 HTTP 接口。

### **1. 创建 client**

1. `mkdir api`  `cd api`
1. `kitex_gen` 目录下，Kitex 已经为我们封装了创建客户端的代码。
1. `item.NewClient` 用于创建 `client`，其第一个参数为调用的 *服务名*，第二个参数为 *options*，用于传入可选参数， 此处的 `client.WithHostPorts` 用于指定服务端的地址，我们可以在运行商品服务时发现其监听在本地的 8888 端口，所以我们指定 8888 端口。
```go
import (
	"example_shop/kitex_gen/example/shop/item/itemservice"
	"github.com/cloudwego/kitex/client"
	...
)
...
c, err := itemservice.NewClient("example.shop.item", client.WithHostPorts("0.0.0.0:8888"))
if err != nil {
	log.Fatal(err)
}
```

### **2. 调用服务**

首先创建了一个请求 `req` , 然后通过 `cli.GetItem` 发起了调用。

1. 其第一个参数为 `context.Context`，通过通常用其传递信息或者控制本次调用的一些行为，你可以在后续章节中找到如何使用它。
1. 其第二个参数为本次调用的请求参数。
1. 其第三个参数为本次调用的 `options` ，Kitex 提供了一种 `callopt` 机制，顾名思义——调用参数 ，有别于创建 client 时传入的参数，这里传入的参数仅对此次生效。 此处的 `callopt.WithRPCTimeout` 用于指定此次调用的超时（通常不需要指定，此处仅作演示之用）同样的，你可以在[基本特性](https://www.cloudwego.io/zh/docs/kitex/tutorials/basic-feature)一节中找到更多的参数。
```go
import "example_shop/kitex_gen/example/shop/item"
...
req := item.NewGetItemReq()
req.Id = 1024
resp, err := cli.GetItem(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
if err != nil {
  log.Fatal(err)
}
```

### **3. 暴露 HTTP 接口**

可以使用 `net/http` 或`Hertz`等其他框架来对外提供 HTTP 接口。

### **4. 测试接口**

打开游览器访问 `localhost:8889/api/item`，看到如下信息，代表请求成功。

```go
GetItemResp({Item:Item({Id:1024 Title:Kitex Description:Kitex is an excellent framework! Stock:0}) BaseResp:BaseResp({Code: Msg:})})
```

## Step6：运行库存服务

在更常见的场景下，一次 RPC 调用远不能实现业务需求，故我们再添加库存服务，模拟更常见的场景。

库存服务的代码我们已经生成，只需要补充业务逻辑即可，与商品服务类似，业务代码位于 `rpc/stock/handler.go`，我们补充以下逻辑：

```go
package main

import (
    "context"

    stock "example_shop/kitex_gen/example/shop/stock"
)

// StockServiceImpl implements the last service interface defined in the IDL.
type StockServiceImpl struct{}

// GetItemStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) GetItemStock(ctx context.Context, req *stock.GetItemStockReq) (resp *stock.GetItemStockResp, err error) {
    resp = stock.NewGetItemStockResp()
    resp.Stock = req.GetItemId()
    return
}
```

由于之前的商品服务和 API 服务分别占用了 8888 和 8889 端口，故我们在库存服务的 `main.go` 中修改监听的端口：

```go
package main

import (
    "log"
    "net"

    stock "example_shop/kitex_gen/example/shop/stock/stockservice"

    "github.com/cloudwego/kitex/server"
)

func main() {
    addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8890")
    svr := stock.NewServer(new(StockServiceImpl), server.WithServiceAddr(addr))

    err := svr.Run()

    if err != nil {
       log.Println(err.Error())
    }
}
```

接下来另启终端使用 `go run .` 运行库存服务，看到以下输出代表运行成功：

```plain text
2024/01/21 00:09:47.076192 server.go:83: [Info] KITEX: server listen at addr=127.0.0.1:8890
```

## Step7：更新handler**商品服务的业务逻辑**

接下来我们补充商品服务，实现对库存服务的调用，与 API 服务类似，我们只需要创建客户端后构造参数发起调用即可。

为了实现 client 的复用，我们在 `ItemServiceImpl` 中补充字段 `stockservice.Client` ，在 `rpc/item/handler.go` 中我们补充以下方法：

```go
package main

import (
    "context"
    "log"

    item "example_shop/kitex_gen/example/shop/item"
	"example_shop/kitex_gen/example/shop/stock"
    "example_shop/kitex_gen/example/shop/stock/stockservice"

    "github.com/cloudwego/kitex/client"
)

// ItemServiceImpl implements the last service interface defined in the IDL.
type ItemServiceImpl struct{
  	stockCli stockservice.Client
}

func NewStockClient(addr string) (stockservice.Client, error) {
    return stockservice.NewClient("example.shop.stock", client.WithHostPorts(addr))
}

// GetItem implements the ItemServiceImpl interface.
func (s *ItemServiceImpl) GetItem(ctx context.Context, req *item.GetItemReq) (resp *item.GetItemResp, err error) {
    resp = item.NewGetItemResp()
    resp.Item = item.NewItem()
    resp.Item.Id = req.GetId()
    resp.Item.Title = "Kitex"
    resp.Item.Description = "Kitex is an excellent framework!"

    stockReq := stock.NewGetItemStockReq()
    stockReq.ItemId = req.GetId()
    stockResp, err := s.stockCli.GetItemStock(context.Background(), stockReq)
    if err != nil {
       log.Println(err)
       stockResp.Stock = 0
    }
    resp.Item.Stock = stockResp.GetStock()
    return
}
```

为 `ItemServiceImpl` 补充了库存服务的客户端，我们需要初始化后才能使用，我们在 `rpc/item/main.go` 中完成初始化操作：

```go
package main

import (
    "log"

    item "example_shop/kitex_gen/example/shop/item/itemservice"
)

func main() {
    itemServiceImpl := new(ItemServiceImpl)
    stockCli, err := NewStockClient("0.0.0.0:8890")
    if err != nil {
       log.Fatal(err)
    }
    itemServiceImpl.stockCli = stockCli

    svr := item.NewServer(itemServiceImpl)

    err = svr.Run()

    if err != nil {
       log.Println(err.Error())
    }
}
```

由于库存服务跑在 8890 端口，所以我们指定 8890 端口创建客户端。

至此，商品服务代码编写完整，参照上文重新编译启动商品服务，看到如下输出代表运行成功：

```shell
2024/01/21 00:18:29.522546 server.go:83: [Info] KITEX: server listen at addr=[::]:8888
```

## Step8: **测试接口**

打开游览器访问 `localhost:8889/api/item`，看到如下信息，代表请求成功。

```shell
GetItemResp({Item:Item({Id:1024 Title:Kitex Description:Kitex is an excellent framework! Stock:1024}) BaseResp:BaseResp({Code: Msg:})})
```

可以看到 `Stock: 1024`，代表我们的商品服务成功请求了库存服务并响应 API 服务。

## Step9：**接入服务注册中心**

为了更贴近真实环境，接下来为我们的服务接入注册中心，在本例中选择了 etcd 作为注册中心，etcd 的安装与使用可参考 [etcd.io](https://etcd.io/) 或使用下述 `docker compose` 文件，接下来默认你已经安装并启动 etcd 服务实例。

```shell
version: "3"

services:
  etcd:
    image: bitnami/etcd:3.5
    container_name: etcd
    ports:
      - 2379:2379
      - 2380:2380
    volumes:
      - ./etcd/data:/bitnami/etcd-data
    environment:
      - TZ=Asia/Shanghai
      - ALLOW_NONE_AUTHENTICATION=yes
      - ETCD_ADVERTISE_CLIENT_URLS=http://etcd:2379
```

Kitex 作为一款微服务框架，也为我们提供了服务治理的能力。当然，在服务注册与发现的场景中也为 etcd 做了适配，可见 [etcd 注册中心使用文档](https://www.cloudwego.io/zh/docs/kitex/tutorials/third-party/service_discovery/etcd/)。同时 Kitex 还提供了其他常见注册中心适配，文档可见[服务发现](https://www.cloudwego.io/zh/docs/kitex/tutorials/third-party/service_discovery/)。

首先，我们需要拉取依赖，在项目根目录下执行如下命令：

```shell
go get github.com/kitex-contrib/registry-etcd
```

### **服务注册**

在 Kitex 中，为服务注册中心抽象出了如下接口：

```go
type Registry interface {
    Register(info *Info) error
    Deregister(info *Info) error
}

```



该接口包括注册服务和注销服务两个方法，两个方法都需要传入服务的信息。只要实现了该接口都可用作服务注册中心使用，所以你也可以自定义服务注册中心，详情可见[服务注册扩展](https://www.cloudwego.io/zh/docs/kitex/tutorials/framework-exten/registry/)。

在 Kitex 中使用注册中心的流程比较简单，大致分为两步：

1. 创建 Registry
1. 在创建服务实例时通过 Option 参数指定 Registry 与服务基本信息
使用 etcd 时，可以使用以下函数创建 Registry：

```go
func NewEtcdRegistry(endpoints []string, opts ...Option) (registry.Registry, error)

func NewEtcdRegistryWithAuth(endpoints []string, username, password string) (registry.Registry, error)

func NewEtcdRegistryWithRetry(endpoints []string, retryConfig *retry.Config, opts ...Option) (registry.Registry, error)
```



指定 Registry 与服务基本信息分别使用 `server.WithRegistry` 与 `server.WithServerBasicInfo` 作为 option 参数在创建服务实例时传入。

本例中，库存服务和商品服务需要对外提供服务，故我们将这两个服务注册进注册中心。



1. **注册库存服务**
在 `rpc/stock/main.go` 中补充以下逻辑：

```go
package main

import (
    "log"
    "net"

    stock "example_shop/kitex_gen/example/shop/stock/stockservice"

    "github.com/cloudwego/kitex/pkg/rpcinfo"
    "github.com/cloudwego/kitex/server"
    etcd "github.com/kitex-contrib/registry-etcd"
)

func main() {
  	// 使用时请传入真实 etcd 的服务地址，本例中为 127.0.0.1:2379
    r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
    if err != nil {
       log.Fatal(err)
    }

    addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8890")
    svr := stock.NewServer(new(StockServiceImpl),
       server.WithServiceAddr(addr),
       // 指定 Registry 与服务基本信息
       server.WithRegistry(r),
       server.WithServerBasicInfo(
          &rpcinfo.EndpointBasicInfo{
             ServiceName: "example.shop.stock",
          },
       ),
    )

    err = svr.Run()

    if err != nil {
       log.Println(err.Error())
    }
}

```



1. **注册商品服务**
在 `rpc/item/main.go` 中补充以下逻辑：

```go
package main

import (
    "log"

    item "example_shop/kitex_gen/example/shop/item/itemservice"

    "github.com/cloudwego/kitex/pkg/rpcinfo"
    "github.com/cloudwego/kitex/server"
    etcd "github.com/kitex-contrib/registry-etcd"
)

func main() {
  	// 使用时请传入真实 etcd 的服务地址，本例中为 127.0.0.1:2379
    r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
    if err != nil {
       log.Fatal(err)
    }

    itemServiceImpl := new(ItemServiceImpl)
    stockCli, err := NewStockClient("0.0.0.0:8890")
    if err != nil {
        log.Fatal(err)
    }
    itemServiceImpl.stockCli = stockCli

    svr := item.NewServer(itemServiceImpl,
       // 指定 Registry 与服务基本信息
       server.WithRegistry(r),
       server.WithServerBasicInfo(
          &rpcinfo.EndpointBasicInfo{
             ServiceName: "example.shop.item",
          }),
    )

    err = svr.Run()

    if err != nil {
       log.Println(err.Error())
    }
}

```

1. **测试**
补充完成代码后我们分别启动这两个服务，两个服务的终端都会出现类似输出：

```plain text
2024/02/04 18:25:22.958727 etcd_registry.go:274: [Info] start keepalive lease 694d8d736e961295 for etcd registry

```

我们再使用 etcdctl 确认是否注册成功，执行 `etcdctl get --prefix "kitex"` 会有如下输出：

```plain text
kitex/registry-etcd/example.shop.item/192.168.196.240:8888
{"network":"tcp","address":"192.168.196.240:8888","weight":10,"tags":null}
kitex/registry-etcd/example.shop.stock/127.0.0.1:8890
{"network":"tcp","address":"127.0.0.1:8890","weight":10,"tags":null}

```

如果都能正确输出，代表我们的服务注册成功。



### **服务发现**

在 Kitex 中，为服务发现中心抽象出了如下接口：

```go
type Resolver interface {
    Target(ctx context.Context, target rpcinfo.EndpointInfo) string
    Resolve(ctx context.Context, key string) (Result, error)
    Diff(key string, prev, next Result) (Change, bool)
    Name() string
}

```



只要实现了该接口都可用作服务发现中心使用，所以你也可以自定义服务发现中心，详情可见[服务发现扩展](https://www.cloudwego.io/zh/docs/kitex/tutorials/framework-exten/service_discovery/)。

在 Kitex 中使用发现中心的流程比较简单，大致分为两步：

1. 创建 Resolver
1. 在创建服务调用客户端时通过 Option 参数指定 Resolver。
使用 etcd 时，可以使用以下函数创建 Resolver：

```go
func NewEtcdResolver(endpoints []string, opts ...Option) (discovery.Resolver, error)

func NewEtcdResolverWithAuth(endpoints []string, username, password string) (discovery.Resolver, error)

```



指定 Resolver 使用 `client.WithResolver()` 作为 option 参数在创建 client 时传入。

本例中，商品服务和 API 服务需要调用其他服务，故我们为这两个服务使用服务发现中心。

1. **商品服务接入**
在前面的代码中，我们将商品服务调用库存服务的逻辑放在了 `rpc/item/handler.go` 中，故我们在此文件的 `NewStockClient` 中添加逻辑：

```go
func NewStockClient() (stockservice.Client, error) {
    // 使用时请传入真实 etcd 的服务地址，本例中为 127.0.0.1:2379
    r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
    if err != nil {
       log.Fatal(err)
    }
    return stockservice.NewClient("example.shop.stock", client.WithResolver(r)) // 指定 Resolver
}

```



1. **API 服务接入**
API 服务只有一个文件，我们在 `api/main.go` 的 `main` 函数中直接添加相关逻辑即可：

```go
func main() {
  	// 使用时请传入真实 etcd 的服务地址，本例中为 127.0.0.1:2379
	resolver, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}

	c, err := itemservice.NewClient("example.shop.item", client.WithResolver(resolver)) // 指定 Resolver
	if err != nil {
		log.Fatal(err)
	}
  	cli = c

    hz := server.New(server.WithHostPorts("localhost:8889"))

	hz.GET("/api/item", Handler)

	if err := hz.Run(); err != nil {
		log.Fatal(err)
	}
}

```



### **测试**

为商品服务添加代码后需要重新启动，然后再启动 API 服务，打开游览器访问 `localhost:8889/api/item`，看到如下信息，代表请求成功。

```plain text
GetItemResp({Item:Item({Id:1024 Title:Kitex Description:Kitex is an excellent framework! Stock:1024}) BaseResp:BaseResp({Code: Msg:})})

```



