---
title: "5.1 thrift"
date: 2025-03-11T14:39:00Z
draft: false
weight: 5001
---

# 5.1 thrift

# **一、RPC与**IDL的基本概念

- **RPC** (Remote Procedure Call，/prəˈsiːdʒər/) ，即远程过程调用。
    - 费曼学习法版本：**像调用本地函数一样调用远端函数**；所以实现RPC框架需要： (通俗来讲，就是**调用远端服务的某个方法(psm.method)，并获取到对应的响应**。)
        1. 通信协议（协议的编解码、序列化）
        1. 传输通信
        1. 服务治理（如服务发现、负载均衡、熔断等）
        1. 代码生成
    - 补充：RPC 本质上**定义了一种通信的流程**，而具体的实现技术没有约束，核心需要解决的问题为**序列化与网络通信**。（如可以通过 `gob/json/pb/thrift` 来序列化和反序列化消息内容，通过`socket/http`来进行网络通信。）只要客户端与服务端在这两方面达成共识，能够做到消息正确的解析接口即可。
- 服务架构的演进：单体应用架构 → 垂直应用架构（按应用拆分） → 分布式服务架构（独立自治、可维护性高、交付速度快）
- rpc开发流程：如基于 **Thrift （/θrɪft/，n*****.***** 节俭, 节约）**的RPC服务开发，通常包括如下过程：
    1. **编写****Thrift****IDL 服务接口定义；**
    1. kitex自动生成**客户端、服务端代码；**
    1. **Server端****修改handler的业务逻辑；**
    1. **Server端编译运行服务监听端口(发布)、接口测试；**
    1. Client端编写客户端程序(overpass)，经过服务发现连接上Server端程序，发起请求并接收响应；
- 一次 rpc 调用的基本流程：
---

<!-- 列布局开始 -->

**Client：**

1. 构造请求参数，发起调用
1. 通过服务发现、负载均衡等得到服务端实例地址，并建立连接； （此步骤中包含的流程称为「**服务治理**」，通常包括并不限于**服务发现（consul） 、负载均衡、ACL（服务鉴权）、熔断、限流**等等功能。这些功能是由其他组件提供的，并不是 Thrift 框架所具有的功能。）
1. 请求参数序列化成二进制数据
1. 通过网络将数据发送给服务端
---

**Client：**

1. 接收数据
1. 反序列化出结果
1. 得到调用的结果

---

**Server：**

1. 服务端接收数据
1. 反序列化出请求参数
1. handler 处理请求并返回响应结果
1. 将响应结果序列化成二进制数据
1. 通过网络将数据返回给客户端
---



<!-- 列布局结束 -->

---

- **IDL：**Interface Definition Language ( /ˌdefɪˈnɪʃn/)，接口定义/描述语言。
    - 费曼版本：**确保双方在说的是同一个语言、同一件事**；如果我们要使用 RPC 进行调用，就需要知道对方的接口是什么，需要传什么参数，同时也需要知道返回值是什么样的。需要IDL 就是为了解决这样的问题，通过 IDL来**约定双方的协议**，就像在写代码的时候需要调用某个函数，我们需要知道 `签名`一样。
    - 优势：
        - 跨语言：中立的方式描述接口和数据结构
        - 规范：明确的定义接口和数据结构和行为
        - 代码生成：IDL编译器，根据IDL自动生成不同编程语言的代码，极大简化了开发工作
# 二、thrift语法

- 文档：[https://github.com/apache/thrift](https://github.com/apache/thrift)
    - 官方文档：[https://thrift.apache.org/docs/idl](https://thrift.apache.org/docs/idl)
    - 用户总结文档（必读）：[https://diwakergupta.github.io/thrift-missing-guide/](https://diwakergupta.github.io/thrift-missing-guide/)
- Thrift **2007年由Facebook开源的RPC框架**，其定义了自己的IDL（接口描述语言）和RPC消息协议，与编程语言无关，可实现跨语言通信。
    - 核心协议：TBinary
    - 常用组合：Framed TBinary
    - 缺陷：**由于开源时间早于微服务框架流行时间，对微服务的支持并不是很好，需要基于源码定制。但由于字节、美团等大厂早期已经使用了Thrift，迁移到gRPC成本大，所以大厂一直沿用了Thrift**
    ![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/df47951e-b71d-4a7c-bed3-28c9301c0325/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB466T5TKXK6X%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005325Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJHMEUCIAEM3qLlNDIXncZSxbWDhAOL2gAC1mJ7UFQq49o1wPUNAiEAqcILoIC4DqoHFzhth5PnAEanS9cURGV2ljqUGLWgec0qiAQImf%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FARAAGgw2Mzc0MjMxODM4MDUiDKJERAEMFTB5ayow5SrcA%2BHrvTZnnni9gsKIujp9A52AEQWxcUkFpXSL4vNEqYiMy9x0kPFrYD3DTUf7x3OOQr9iwRs6SexXxiU5rgOC%2BSnFQFowJfNhuNjnbq5ClGUl3%2BI7sLhoyYOrtb5sxkXyID1wghBca%2BCEfkAzMRkQi2VE48YTjJHS%2BsVKxKCN7NtL1Wmd1Tl4v4TKbFRtP8ZfS4G38JniXnrKQ7cvYcapavUMbQ2FE42VEeURAGFXdmjaLK5e5%2BXj1zz%2FyOO0%2Bs30vdj%2FuDlHFkhEwqQa82b6eJMPtsWtWpaSXE%2BCBl1%2BHNjsiKv41B8vqo6ReIkh8AXz80ZvBhdekVLf%2FXEof1aQUFmE5vpkuxIWzCaEVlNAazi0%2FGoN%2BCysiqTybALJs0JfTdoJAXj9RB3Gc6XFIM2d4C6R34yoxMUk6ItsdECCyaZ1hfjEriHoH53EsI5rqAk8ehuFnbSWVcAEH%2Bw9fRRfcTSqPg7lewacf4rSSSMhzZsi%2B0KMcntDqeRXL2dL%2Ffhv63jy%2BUzr1ERD1Z2mu%2BqMscO0sUkkLBcolrcctuY0%2BWJigXNdBfknYQ3UOllY8BFSia9k2FuVvwOXEexykBCgQq94ZZb%2FaMNRS5gNwyoerWSUuoJSaTYdS9xRiI5KMLa668MGOqUBuiqYNqzKqz88du1oRsrwUrltjquQBNQ6eMQm8O95aMoMb%2Ft7ZrJo8B3hp%2BlE7NjJVYzlBpTGzj9bNZx6NW87H624eWyFKen%2B8FPqSUjuFgLzhaOL3xNlHNKaIG%2BFlv0UDtxXwjJ1bM99zgpdnkYuzj1UAH0i1X2wr7nkwyBStJzogYRUuD%2BjCYny5x9o9W4HjXUrlPe2jag3f3F0SHaOSv151QKO&X-Amz-Signature=e22c2003106d45b51ce7bb3162bca03cfb7568c75c46b09864bee3ba7a60a730&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)




Kitex 默认支持 `thrift` 和 `proto3` 两种 IDL。本文简单介绍 Thrift IDL 语法，

> 注意：Thrift 是一款 RPC 框架，其使用的 IDL 以 .thrift 为后缀，故常常也使用 thrift 来表示 IDL，请根据上下文判断语意。

## 类型

### **1. 基本类型**

Thrift IDL 有以下几种基本类型：

- **bool**: 布尔型
- **byte**: 有符号字节
- **i8**: 8位有符号整型（低版本不支持）
- **i16**: 16位有符号整型
- **i32**: 32位有符号整型
- **i64**: 64位有符号整型 （**注意**：Thrift IDL 没有无符号整数类型。因为许多编程语言中没有原生的无符号整数类型。）
- **double**: 64位浮点型
- **string**: 字符串（编码方式未知或二进制字符串）
- **binary**: 表示无编码要求的 byte 二进制数组。因此是字节数组情况下（比如 json/pb 序列化后数据在 thrift rpc 间传输）请使用 binary 类型，不要使用 string 类型；
    - Golang 的实现 string/binary 都使用 string 类型进行存储，string 底层只是字节数组不保证是 UTF-8 编码的，可能与其他语言的行为不一致。
### **2. 容器类型**

Thrift 提供的容器是强类型容器，映射到大多数编程语言中常用的容器类型。具体包括以下三种容器：

- **list< t1 >**: 元素类型为 t1 的有序列表，允许元素重复。Translates to an STL vector, Java ArrayList, native arrays in scripting languages, etc.
- **set< t1 >**: 元素类型为 t1 的无序表，不允许元素重复。
- **map<t1,t2>**: 键类型为 t1，值类型为 t2 的 map。


### **3. 枚举类型**

Thrift 提供了枚举类型

- 编译器默认从 0 开始赋值
- 可以对某个变量进行赋值（整数）
- 不支持嵌套的 enum
```c++
enum TweetType {
    TWEET, //
    RETWEET = 2, //
    DM = 0xa,
    REPLY
}

struct TweetType {
		1: optional TweeType tweeType = TweeType.TWEET
}
```

### **4. 类型定义**

Thrift 支持类似 C/C++ 的类型定义，**注意：typedef** 定义的末尾没有分号

```c++
typedef i32 MyInteger
typedef Tweet ReTweet
```



### **5. 常量定义**

Thrift 内定义常量的方式如下：

```c++
const i32 INT_CONST = 1234;

const map<string,string> MAP_CONST = {
    "hello": "world",
    "goodnight": "moon"
}
```



### **7. Struct 及 Requiredness 说明**

Struct 由不同的 fields 构成，其中每个 **field** 有唯一的整型 **id**，类型 **type**，名字 **name** 和 一个可选择设置的默认值 **default value**。

1. **Field id**：每个 field 必须有一个正整数的标志符
1. type：包括三种类型: 
    1. required：必填字段，如果对端没有收到该字段会返回error。（从维护角度不建议用 required 修饰字段）
        - 注意：该修饰在 thrift 官方的解释如下，期望被 set，在 Golang 的实现里如果某个字段为 nil，实际还是会编码，所以对端收到的是空struct。
    1. optional：可选字段。未赋值时不会编码，对端也不够早
        - 若没有设置该字段且**没有默认值**的话，则不对该字段进行序列化
        - 对于非指针字段，**需要调用** NewXXX 方法来初始化结构体**才能填入默认值**，不能用 &XXX{} 方式
    1. default：不加修饰则是 default 类型，即使未赋值也会编码
        - 注意：**发送方发 nil，接收方会构造默认值**，如果希望接收方同样接收 nil 需要用 optional 修饰
```c++
struct Location {
    1: required double latitude;
    2: required double longitude;
}

struct Tweet {
    1: required i32 userId;
    2: required string userName;
    3: required string text;
    4: optional Location loc; // Struct的定义内可以包含其他 Struct
    16: optional string language = "english" // 可设置默认值
}
```

**注意：**

1. Thrift 不支持嵌套定义 Struct
1. **如果 struct 已经在使用了，请不要更改各个 field 的 id 和 type**
1. 如果没有特殊需求，**建议都使用 optional**。由于 Kitex 需要保留和 apache 官方 Go 实现兼容性，也保留了对于 required 和 default 修饰的字段处理逻辑的不合理之处。例如 Request(struct类型).User(struct类型).Name(string类型) 这样一个结构，如果 User 和 Name 都是 required，但 client 端没有给 User 赋值（即 request.User == nil), 在 client 端编码不会报错，但是会将 User 的 id 和 type(struct) 写入，在 server 端解码时会初始化 User（即 request.User != nil），但是在继续解码 User 时读不到 Name 字段，就会报错。
### **8. Exception**

Exception 与 struct 类似，但它被用来集成 目标编程语言 中的异常处理机制。Exception 内定义的所有 field 的名字都是唯一的。

## IDL示例

### **1. 注释**

Thrift 支持 c风格的多行注释 和 c++/Java 风格的单行注释

```c++
/*
* This is a multi-line comment.
* Just like in C.
*/

// C++/Java style single-line comments work just as well.

```

### **2. 命名空间**namespace

Thrift 的命名空间与 C++ 的 namespace 和 go 的 package 类似，提供了一种组织（隔离）代码的方式，也可避免类型定义内名字冲突的问题。

Thrift 提供了针对不同语言的 namespace 定义方式，各语言建议单独定义：

```c++
namespace cpp com.example.project
namespace java com.example.project
namespace go com.example.project
```

### **3. Include**

引用其他IDL文件中的定义，通过include可以很少的实现模块化、复用、拆分结构体定义，方便管理、维护 IDL。用户可利用文件名作为前缀对具体定义进行访问。

```c++
include "tweet.thrift" ...

struct TweetSearchResult {
    1: list<tweet.Tweet> tweets;
}
```

### **4. RPC Service定义**

Thrift 内的 service 定义在语义上和 oop 内的接口是相同的。代码生成工具会根据 service 的定义生成 client 和 service 端的接口实现。

- Service 的参数和返回值类型可以是 基础类型 或者 struct
> oneway 本身不具有可靠性，且在处理上比较特殊会带来一些隐患，不建议使用

```c++
service Twitter {
    // A method definition looks like C code. It has a return type, arguments,
    // and optionally a list of exceptions that it may throw. Note that argument
    // lists and exception list are specified using the exact same syntax as
    // field lists in structs.

    void ping(); // 1
    bool postTweet(1:Tweet tweet); // 2
    TweetSearchResult searchTweets(1:string query); // 3

    // The 'oneway' modifier indicates that the client only makes a request and does not wait for any response at all. Oneway methods MUST be void.
    oneway void zip() // 4
}

```

### **5. IDL 示例**

以下为简单的 thrift idl 示例，包含 common.thrift 和 service.thrift 两个文件。

- common.thrift：包含各种类型的使用和 struct 的定义。
```c++
namespace go example.common

// typedef
typedef i32 TestInteger

// Enum
enum TestEnum {
    Enum1 = 1,
    Enum2,
    Enum3 = 10,
}

// Constant
const i32 TestIntConstant = 1234;

// Struct
struct TestStruct {
    1: bool sBool
    2: required bool sBoolReq
    3: optional bool sBoolOpt
    4: list<string> sListString
    5: set<i16> sSetI16
    6: map<i32,string> sMapI32String
}
```



- service.thrift：引用 common.thrift，定义 service。
    - tMethod：接收一个类型为 TestRequest 的参数，返回一个类型为 TestResponse 的返回值。
```c++
namespace go example.service

include "common.thrift"

struct TestRequest {
   1: string msg
   2: common.TestStruct s
}

struct TestResponse {
   1: string msg
   2: common.TestStruct s
}

service TestService {
   TestResponse tMethod(1: TestRequest req)
}

```



### **6. Kitex Thrift IDL 规范**

为满足服务调用的规范，Kitex 对 IDL 的定义提出了一些必须遵守的要求：

- 方法只能拥有一个参数，并且这个参数类型必须是自定义的 Struct 类型，参数类型名字使用驼峰命名法，通常为：`XXXRequest`
- 方法的返回值类型必须是自定义的 Struct 类型，不可以为 void，使用驼峰命名法，通常为：`XXXResponse`


## Thrift消息协议-TBinary编码

- 补充：由于JSON数据格式传输体积大，编码低效，性能损失较多，而RPC场景对性能有更高的要求，所以很少采用TJSON协议。




