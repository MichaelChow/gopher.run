---
title: "5.5 ToolsNode & Tool & Lambda & eino-ext"
date: 2025-08-28T13:33:00Z
draft: false
weight: 5005
---

# 5.5 ToolsNode & Tool & Lambda & eino-ext



# 一、tool

> 💡 引用：
> - [Eino: ToolsNode&Tool 使用说明](https://www.cloudwego.io/zh/docs/eino/core_modules/components/tools_node_guide/)
> - [如何创建一个 tool ?](https://www.cloudwego.io/zh/docs/eino/core_modules/components/tools_node_guide/how_to_create_a_tool/)
> - [tool ext：googlesearch、mcp](https://www.cloudwego.io/zh/docs/eino/ecosystem_integration/tool/)

## 1. tool定义

一个用于扩展模型能力的组件，它允许模型调用外部工具来完成特定的任务。

应用场景：

- 让模型能够获取实时信息（如搜索引擎、天气查询等）
- 使模型能够执行特定的操作（如数据库操作、API 调用等）
- 扩展模型的能力范围（如数学计算、代码执行等）
- 与外部系统集成（如知识库查询、插件系统等）
### **interface定义 与**ToolInfo struct

**interface：**

```go
// eino/compose/tool/interface.go

// 基础工具interface，提供工具信息
type BaseTool interface {
		// 获取工具的描述信息*schema.ToolInfo，用于提供给大模型
    Info(ctx context.Context) (*schema.ToolInfo, error)
}

// 支持同步调用的工具interface
type InvokableTool interface {
    BaseTool
    // 同步执行工具
    // 参数：上下文对象，用于传递请求级别的信息和Callback Manager，JSON 格式的参数字符串、工具执行的选项
    // 返回值：string工具调用结果
    InvokableRun(ctx context.Context, argumentsInJSON string, opts ...Option) (string, error)
}

// 支持流式输出的工具interface
type StreamableTool interface {
    BaseTool
    // 流式执行工具
    // 参数：同上
    // 返回值：*schema.StreamReader[string]流式工具调用结果
    StreamableRun(ctx context.Context, argumentsInJSON string, opts ...Option) (*schema.StreamReader[string], error)
}
```

**设计特点：**分层设计、JSON接口、流式支持、选项模式



**ToolInfo struct：**

BaseTool interface的Info() 返回**ToolInfo struct，告诉大模型如何具体构造符合约束的function call参数：**

```go
// eino/schema/tool.go

type ToolInfo struct {
    // 工具的唯一名称，用于清晰地表达其用途
    Name string
    // 用于告诉模型如何/何时/为什么使用这个工具；可以在描述中包含少量示例
    Desc string
    // 工具接受的参数定义，可以通过两种方式描述：
    // 1. 使用 ParameterInfo：schema.NewParamsOneOfByParams(params)
    // 2. 使用 OpenAPIV3：schema.NewParamsOneOfByOpenAPIV3(openAPIV3)
    *ParamsOneOf
}
```



方式 1 - map[string]*ParameterInfo：用map，key 即为参数名，value 则是这个参数的详细约束。

简单直观，当参数由开发者通过编码的方式手动维护时常用。

```go
// 结构定义详见: https://github.com/cloudwego/eino/blob/main/schema/tool.go
type ParameterInfo struct {
    Type DataType    // The type of the parameter.
    ElemInfo *ParameterInfo    // The element type of the parameter, only for array.
    SubParams map[string]*ParameterInfo    // The sub parameters of the parameter, only for object.
    Desc string    // The description of the parameter.
    Enum []string    // The enum values of the parameter, only for string.
    Required bool    // Whether the parameter is required.
}

// example
map[string]*schema.ParameterInfo{
    "name": &schema.ParameterInfo{
        Type: schema.String,
        Required: true,
    },
    "age": &schema.ParameterInfo{
        Type: schema.Integer,
    },
    "gender": &schema.ParameterInfo{
        Type: schema.String,    
        Enum: []string{"male", "female"},
    },
}
```

方式2：**openapi3.Schema。****一般是直接使用**`**utils.GoStruct2ParamsOneOf**`** 来构建ToolInfo 或 直接用 **`**utils.InferTool()**`** 直接构建 tool**，而不由开发者自行直接调用构造此结构体。

Eino 提供了在结构体中通过 go tag 描述参数约束的方式，并提供了 GoStruct2ParamsOneOf 方法来生成一个 struct 的参数约束:

```go
func GoStruct2ParamsOneOf[T any](opts ...Option) (*schema.ParamsOneOf, error)
```

```go
import (
    "context"
    "github.com/cloudwego/eino/components/tool/utils"
)

type User struct {
    Name   string `json:"name" jsonschema:"required,description=the name of the user"`
    Age    int    `json:"age" jsonschema:"description=the age of the user"`
    Gender string `json:"gender" jsonschema:"enum=male,enum=female"`
}

params, err := utils.GoStruct2ParamsOneOf[User]()
```

如果 tool 是对一些 openapi 的封装，则可以通过 导出openapi.json文件来生成。



### **公共 Option**

Tool 组件使用 ToolOption 来定义可选参数， ToolsNode 没有抽象公共的 option。每个具体的实现可以定义自己的特定 Option，通过 WrapToolImplSpecificOptFn 函数包装成统一的 ToolOption 类型。

```go
package tool

// Option defines call option for InvokableTool or StreamableTool component, which is part of component interface signature.
// Each tool implementation could define its own options struct and option funcs within its own package,
// then wrap the impl specific option funcs into this type, before passing to InvokableRun or StreamableRun.
type Option struct {
	implSpecificOptFn any
}

// WrapImplSpecificOptFn wraps the impl specific option functions into Option type.
// T: the type of the impl specific options struct.
// Tool implementations are required to use this function to convert its own option functions into the unified Option type.
// For example, if the tool defines its own options struct:
//
//	type customOptions struct {
//	    conf string
//	}
//
// Then the tool needs to provide an option function as such:
//
//	func WithConf(conf string) Option {
//	    return WrapImplSpecificOptFn(func(o *customOptions) {
//			o.conf = conf
//		}
//	}
func WrapImplSpecificOptFn[T any](optFn func(*T)) Option {
	return Option{
		implSpecificOptFn: optFn,
	}
}
```

**设计亮点：**类型安全、延迟应用、默认值支持

## 2. tool使用-原子项

### **ToolsNode 使用example**

通常不会被单独使用，一般用于编排之中接在 ChatModel 之后。

```go
// 创建工具节点，一个工具节点可包含多种工具
toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
    Tools: []tool.BaseTool{
        searchTool,     // 搜索工具
        weatherTool,    // 天气查询工具
        calculatorTool, // 计算器工具
    },
})
if err != nil {
    return err
}

// 在 Chain 中使用
chain := compose.NewChain[*schema.Message, []*schema.Message]()
chain.AppendToolsNode(toolsNode)

// graph 中
graph := compose.NewGraph[*schema.Message, []*schema.Message]()
graph.AddToolsNode(toolsNode)
```

### **Option 使用example**

```go
import "github.com/cloudwego/eino/components/tool"
// 定义 Option 结构体
type MyToolOptions struct {
    Timeout time.Duration
    MaxRetries int
    RetryInterval time.Duration
}

// 定义 Option 函数
func WithTimeout(timeout time.Duration) tool.Option {
    return tool.WrapImplSpecificOptFn(func(o *MyToolOptions) {
        o.Timeout = timeout
    })
}

type MyTool struct {
    options MyToolOptions
}

func (t *MyTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
    // 省略具体实现
    return nil, err
}
func (t *MyTool) StreamableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (*schema.StreamReader[string], error) {
    // 省略具体实现
    return nil, err
}

func (t *MyTool) InvokableRun(ctx context.Context, argument string, opts ...tool.Option) (string, error) {
	// 将执行编排时传入的自定义配置设置到MyToolOptions中
    tmpOptions := tool.GetImplSpecificOptions(&t.options, opts...)

    // 根据tmpOptions中Timeout的值处理Timeout逻辑
	return "", nil
}
```

执行编排时，可以使用 compose.WithToolsNodeOption() 传入 ToolsNode 相关的Option设置，ToolsNode下的所有 Tool 都能接收到

```go
streamReader, err := graph.Stream(ctx, []*schema.Message{
    {
        Role:    schema.User,
        Content: "hello",
    },
}, compose.WithToolsNodeOption(compose.WithToolOption(WithTimeout(10 * time.Second))))
```

### **Callback 使用example**

```go
import (
    "context"

    callbackHelper "github.com/cloudwego/eino/utils/callbacks"
    "github.com/cloudwego/eino/callbacks"
    "github.com/cloudwego/eino/compose"
    "github.com/cloudwego/eino/components/tool"
)

// 创建 callback handler
handler := &callbackHelper.ToolCallbackHandler{
    OnStart: func(ctx context.Context, info *callbacks.RunInfo, input *tool.CallbackInput) context.Context {
       fmt.Printf("开始执行工具，参数: %s\n", input.ArgumentsInJSON)
       return ctx
    },
    OnEnd: func(ctx context.Context, info *callbacks.RunInfo, output *tool.CallbackOutput) context.Context {
       fmt.Printf("工具执行完成，结果: %s\n", output.Response)
       return ctx
    },
    OnEndWithStreamOutput: func(ctx context.Context, info *callbacks.RunInfo, output *schema.StreamReader[*tool.CallbackOutput]) context.Context {
       fmt.Println("工具开始流式输出")
       go func() {
          defer output.Close()

          for {
             chunk, err := output.Recv()
             if errors.Is(err, io.EOF) {
                return
             }
             if err != nil {
                return
             }
             fmt.Printf("收到流式输出: %s\n", chunk.Response)
          }
       }()
       return ctx
    },
}

// 使用 callback handler
helper := callbackHelper.NewHandlerHelper().
    Tool(handler).
    Handler()
 
/*** compose a chain
* chain := NewChain
* chain.appendxxx().
*       appendxxx().
*       ...
*/

// 在运行时使用
runnable, err := chain.Compile()
if err != nil {
    return err
}
result, err := runnable.Invoke(ctx, input, compose.WithCallbacks(helper))
```

## 3. tool使用-已有实现

### ext-tool search

1. `bingsearch`** Tool example：实现了**`tool.InvokableTool` 接口，可配置的搜索参数
**使用example：**

```go
import (
	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino-ext/components/tool/bingsearch"
)

// 设置 Bing Search API 密钥
bingSearchAPIKey := os.Getenv("BING_SEARCH_API_KEY")
	
// 创建上下文
ctx := context.Background()
	
// 创建 Bing Search 工具
bingSearchTool, err := bingsearch.NewTool(ctx, &bingsearch.Config{
		APIKey: bingSearchAPIKey,
		Cache:  5 * time.Minute,
})
	
if err != nil {
	log.Fatalf("Failed to create tool: %v", err)
}
// ... 配置并使用 ToolsNode
```

使用Config struct 进行tool配置：

```go
type Config struct {
// Config represents the Bing search tool configuration.
type Config struct {
    ToolName string `json:"tool_name"` // optional, default is "bing_search"
    ToolDesc string `json:"tool_desc"` // optional, default is "search web for information by bing"

    APIKey     string     `json:"api_key"`     // required
    Region     Region     `json:"region"`      // optional, default: ""
    MaxResults int        `json:"max_results"` // optional, default: 10
    SafeSearch SafeSearch `json:"safe_search"` // optional, default: SafeSearchModerate
    TimeRange  TimeRange  `json:"time_range"`  // optional, default: nil

    Headers    map[string]string `json:"headers"`     // optional, default: map[string]string{}
    Timeout    time.Duration     `json:"timeout"`     // optional, default: 30 * time.Second
    ProxyURL   string            `json:"proxy_url"`   // optional, default: ""
    Cache      time.Duration     `json:"cache"`       // optional, default: 0 (disabled)
    MaxRetries int               `json:"max_retries"` // optional, default: 3
}
```

**请求 Schema & 响应 Schema**

```go
type SearchResponse struct {
    Results []*searchResult `json:"results" jsonschema_description:"The results of the search"`
}

type searchResult struct {
    Title       string `json:"title" jsonschema_description:"The title of the search result"`
    URL         string `json:"url" jsonschema_description:"The link of the search result"`
    Description string `json:"description" jsonschema_description:"The description of the search result"`
}
```

1. **Googlesearch tool：**实现了`tool.InvokableTool` 接口，通过 Google Custom Search API 进行网络搜索
**使用example：**

```go
import "github.com/cloudwego/eino-ext/components/tool/googlesearch"

ctx := context.Background()
googleAPIKey := os.Getenv("GOOGLE_API_KEY")
googleSearchEngineID := os.Getenv("GOOGLE_SEARCH_ENGINE_ID")
if googleAPIKey == "" || googleSearchEngineID == "" {
   log.Fatal("[GOOGLE_API_KEY] and [GOOGLE_SEARCH_ENGINE_ID] must set")
}

// create tool
searchTool, err := googlesearch.NewTool(ctx, &googlesearch.Config{
   APIKey:         googleAPIKey,  // Google API 密钥
   SearchEngineID: googleSearchEngineID,  // 搜索引擎 ID
   Lang:           "zh-CN",  // 可选：搜索界面语言
   Num:            5,  // 可选：每页结果数量
   
   BaseURL:        "custom-base-url",     // 可选：自定义 API 基础 URL, default: https://customsearch.googleapis.com
   ToolName:       "google_search",       // 可选：工具名称
   ToolDesc:       "google search tool",  // 可选：工具描述
})
if err != nil {
   log.Fatal(err)
}

// prepare 搜索的请求params
req := googlesearch.SearchRequest{
   Query: "Golang concurrent programming",  // 搜索关键词
   Num:   3,  // 返回结果数量
   Lang:  "en",  // 搜索语言
}

args, err := json.Marshal(req)
if err != nil {
   log.Fatal(err)
}

// do search
resp, err := searchTool.InvokableRun(ctx, string(args))
if err != nil {
   log.Fatal(err)
}

var searchResp googlesearch.SearchResult
  if err := json.Unmarshal([]byte(resp), &searchResp); err != nil {
    log.Fatal(err)
}

// Print results
fmt.Println("Search Results:")
fmt.Println("==============")
for i, result := range searchResp.Items {
   fmt.Printf("\n%d. Title: %s\n", i+1, result.Title)
   fmt.Printf("   Link: %s\n", result.Link)
   fmt.Printf("   Desc: %s\n", result.Desc)
}
fmt.Println("")
fmt.Println("==============")

// seems like:
// Search Results:
// ==============
// 1. Title: My Concurrent Programming book is finally PUBLISHED!!! : r/golang
//    Link: https://www.reddit.com/r/golang/comments/18b86aa/my_concurrent_programming_book_is_finally/
//    Desc: Posted by u/channelselectcase - 398 votes and 46 comments
// 2. Title: Concurrency — An Introduction to Programming in Go | Go Resources
//    Link: https://www.golang-book.com/books/intro/10
//    Desc:
// 3. Title: The Comprehensive Guide to Concurrency in Golang | by Brandon ...
//    Link: https://bwoff.medium.com/the-comprehensive-guide-to-concurrency-in-golang-aaa99f8bccf6
//    Desc: Update (November 20, 2023) — This article has undergone a comprehensive revision for enhanced clarity and conciseness. I’ve streamlined the…
// ==============
```



1. `duckduckgo` search tool：实现了`tool.InvokableTool` 接口，通过 DuckDuckGo 搜索引擎进行网络搜索。DuckDuckGo 是一个注重隐私的搜索引擎，不会追踪用户的搜索行为，**无需 api key 鉴权即可直接使用**。
```go
 import (
		 "github.com/cloudwego/eino-ext/components/tool/duckduckgo"
    "github.com/cloudwego/eino-ext/components/tool/duckduckgo/ddgsearch"
)

ctx := context.Background()

// **init search client**
tool, err := duckduckgo.NewTool(ctx, &duckduckgo.Config{
	ToolName:    "duckduckgo_search",     // 工具名称
   ToolDesc:    "search web for information by duckduckgo", // 工具描述
   Region:     ddgsearch.RegionCN, // 搜索地区
   MaxResults: 10, // 每页结果数量
   SafeSearch:  ddgsearch.SafeSearchOff, // 安全搜索级别
   TimeRange:   ddgsearch.TimeRangeAll,  // 时间范围
   
   DDGConfig: &ddgsearch.Config{  // DuckDuckGo 配置
	   Timeout:    10 * time.Second,
	   Cache:      true,
	   MaxRetries: 5,
   },
})
if err != nil {
   log.Fatalf("NewTool of duckduckgo failed, err=%v", err)
}

// 搜索参数
searchReq := &duckduckgo.SearchRequest{
    Query: "Golang programming development", // 搜索关键词
    Page:  1, // 页码
}

jsonReq, err := json.Marshal(searchReq)
if err != nil {
        log.Fatalf("Marshal of search request failed, err=%v", err)
}

// Execute search
resp, err := tool.InvokableRun(ctx, string(jsonReq))
if err != nil {
   log.Fatalf("Search of duckduckgo failed, err=%v", err)
}

var searchResp duckduckgo.SearchResponse
  if err := json.Unmarshal([]byte(resp), &searchResp); err != nil {
    log.Fatalf("Unmarshal of search response failed, err=%v", err)
}

// Print results
fmt.Println("Search Results:")
fmt.Println("==============")
for i, result := range searchResp.Results {
  fmt.Printf("\n%d. Title: %s\n", i+1, result.Title)
  fmt.Printf("   Link: %s\n", result.Link)
  fmt.Printf("   Description: %s\n", result.Description)
}
```



### ext-tool **httprequest**

实现了`tool.InvokableTool` 接口，支持bind ChatModel来发起 GET、POST、PUT 和 DELETE 请求，可配置请求头和 HttpClient

**发起GET 请求**使用example：

```go
import (
	"github.com/bytedance/sonic"
	req "github.com/cloudwego/eino-ext/components/tool/httprequest/get"
)


// Configure the GET tool
config := &req.Config{
// Headers is optional
Headers: map[string]string{
	"User-Agent": "MyCustomAgent",
},
// HttpClient is optional
HttpClient: &http.Client{
		Timeout:   30 * time.Second,
	  Transport: &http.Transport{},
	},
}

ctx := context.Background()

// Create the GET tool
tool, err := req.NewTool(ctx, config)
if err != nil {
	log.Fatalf("Failed to create tool: %v", err)
}

// Prepare the GET request payload
request := &req.GetRequest{
	URL: "https://jsonplaceholder.typicode.com/posts",
}

jsonReq, err := sonic.Marshal(request)
if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
}

// Execute the GET request using the InvokableTool interface
resp, err := tool.InvokableRun(ctx, string(jsonReq))
if err != nil {
	log.Fatalf("GET request failed: %v", err)
}

fmt.Println(resp)
```



### ext-tool **Commandline**

[https://www.cloudwego.io/zh/docs/eino/ecosystem_integration/tool/tool_commandline/](https://www.cloudwego.io/zh/docs/eino/ecosystem_integration/tool/tool_commandline/)

### ext-tool **Browseruse**

[Tool - Browseruse](https://www.cloudwego.io/zh/docs/eino/ecosystem_integration/tool/tool_browseruse/)

### ext-tool **sequentialthinking**

[Tool - sequentialthinking](https://www.cloudwego.io/zh/docs/eino/ecosystem_integration/tool/tool_sequentialthinking/)

### ext-tool **wikipedia**

实现了`tool.InvokableTool` 接口，维基百科搜索工具

```go
import (
	"github.com/cloudwego/eino-ext/components/tool/wikipedia"
	"github.com/cloudwego/eino/components/tool"
)

ctx := context.Background()

// 创建工具配置
// 下面所有这些参数都是默认值，仅作用法展示
config := &wikipedia.Config{
		UserAgent:   "eino (https://github.com/cloudwego/eino)",
		DocMaxChars: 2000,
		Timeout:     15 * time.Second,
		TopK:        3,
		MaxRedirect: 3,
		Language:    "en",
}

// 创建搜索工具
t, err := wikipedia.NewTool(ctx, config)
if err != nil {
	log.Fatal("Failed to create tool:", err)
}

// 与 Eino 的 ToolsNode 一起使用
tools := []tool.BaseTool{t}
// ... 配置并使用 ToolsNode
```



### mcp tool

[Model Context Protocol(MCP)](https://modelcontextprotocol.io/introduction)是 Anthropic 推出的供模型访问的资源的标准化开放协议。

Eino封装 实现了 Eino InvokableTool 接口，可以直接访问已有 MCP Server 上的资源。

![](/images/25d24637-29b5-80e4-b275-e8ec8dbc3df3/image_29f24637-29b5-804e-83f5-de05e05aa1be.jpg)

使用example：

```go
import (
    "github.com/mark3labs/mcp-go/client"  // 使用开源 sdk mark3labs/mcp-go
    mcpp "github.com/cloudwego/eino-ext/components/tool/mcp"
)

func getMCPTool(ctx context.Context) []tool.BaseTool {
   //  创建mcp client
		// sse client  needs to manually start asynchronous communication
		// while stdio does not require it.
   cli, err := client.NewSSEMCPClient("http://localhost:12345/sse")
   // stdio client
	 cli, err := client.NewStdioMCPClient(myCommand, myEnvs, myArgs...)
	 // 其他创建 Client 的方法（比如 InProcess），更多信息可以参考：https://mcp-go.dev/transports
   if err != nil {
      log.Fatal(err)
   }
   err = cli.Start(ctx)
   if err != nil {
	   log.Fatal(err)
    }
		
		// 用户需要自行完成 client 初始化。考虑到 client 的复用，封装假设 client 已经完成和 Server 的 [Initialize](https://spec.modelcontextprotocol.io/specification/2024-11-05/basic/lifecycle/)。
    initRequest := mcp.InitializeRequest{}
    initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
    initRequest.Params.ClientInfo = mcp.Implementation{
       Name:    "example-client",
       Version: "1.0.0",
    }
		 _, err = cli.Initialize(ctx, initRequest)
     if err != nil {
        log.Fatal(err)
     }
        
     // 使用 Client 创建 Eino Tool：
     tools, err := mcpp.GetTools(ctx, &mcpp.Config{
		     Cli: cli
					ToolNameList: []string{"name"}, // 可选。支持使用 Name 筛选 Server 提供的 Tools，避免调用预期外的 Tools。
		})
    if err != nil {
        log.Fatal(err)
     }
     return tools
}

// 在react agent中直接使用
agent, err := react.NewAgent(ctx, &react.AgentConfig{
    Model:                 llm,
    ToolsConfig:           compose.ToolsNodeConfig{Tools: mcpTools},
})
```

代码参考： [https://github.com/cloudwego/eino-ext/blob/main/components/tool/mcp/examples/mcp.go](https://github.com/cloudwego/eino-ext/blob/main/components/tool/mcp/examples/mcp.go)

## 4. tool使用-自行实现

### **1. utils.InferTool把本地函数转为**tool.InvokableTool

开发场景中经常需要把一个现有的本地函数（如`AddUser()`）封装成 Eino tool.InvokableTool，用来 bind 到ChatModel上。

1. **可使用 utils.InferTool 来更简洁的构建****（**~~**NewTool()的语法糖**~~**）**：
- 参数约束直接维护在 input struct的tag参数定义中（参考上方 `GoStruct2ParamsOneOf`）时（参数结构体和描述信息同源，无需维护两份信息）。
```go
// 通过**utils.InferTool**把AddUser()封装成tool
func createTool() (tool.InvokableTool, error) {aqaswq
    return utils.InferTool("add_user", "add user", AddUser)
}
```

- 使用example：
```go
import (
    "github.com/cloudwego/eino/components/tool"
    "github.com/cloudwego/eino/components/tool/utils"
    "github.com/cloudwego/eino/schema"
)

type User struct {
    Name   string `json:"name" jsonschema:"required,description=the name of the user"`
    Age    int    `json:"age" jsonschema:"description=the age of the user"`
    Gender string `json:"gender" jsonschema:"enum=male,enum=female"`
}

type Result struct {
    Msg string `json:"msg"`
}

func AddUser(ctx context.Context, user *User) (*Result, error) {
    // some logic
}
```

- 内部实现：
```go
// InferTool creates an InvokableTool from a given function by inferring the ToolInfo from the function's request parameters.
// End-user can pass a SchemaCustomizerFn in opts to customize the go struct tag parsing process, overriding default behavior.
func InferTool[T, D any](toolName, toolDesc string, i InvokeFunc[T, D], opts ...Option) (tool.InvokableTool, error) {
	ti, err := goStruct2ToolInfo[T](toolName, toolDesc, opts...) // goStruct2ToolInfo
	if err != nil {
		return nil, err
	}

	return NewTool(ti, i, opts...), nil // 调用NewTool
}
```



1. **使用 InferOptionableTool 方法**
Option 机制是 Eino 提供的一种在运行时传递动态参数的机制。

当开发者要实现一个**需要自定义 option 参数时**，则可使用 InferOptionableTool 这个方法，增加了一个 option 参数。

```go
func useInInvoke() {
    ctx := context.Background()
    tl, _ := utils.InferOptionableTool("invoke_infer_optionable_tool", "full update user info", updateUserInfoWithOption)

    content, _ := tl.InvokableRun(ctx, `{"name": "bruce lee"}`, WithUserInfoOption("hello world"))

    fmt.Println(content) // Msg is "hello world", because WithUserInfoOption change the UserInfoOption.Field1
}
```

示例如下（改编自 [cloudwego/eino/components/tool/utils/invokable_func_test.go](https://github.com/cloudwego/eino/blob/main/components/tool/utils/invokable_func_test.go)）：

```go
import (
    "fmt"
    "context"
    "github.com/cloudwego/eino/components/tool"
    "github.com/cloudwego/eino/components/tool/utils"
    "github.com/cloudwego/eino/schema"
)

type UserInfoOption struct {
    Field1 string
}

func WithUserInfoOption(s string) tool.Option {
    return tool.WrapImplSpecificOptFn(func(t *UserInfoOption) {
        t.Field1 = s
    })
}

func updateUserInfoWithOption(_ context.Context, input *User, opts ...tool.Option) (output *UserResult, err error) {
    baseOption := &UserInfoOption{
        Field1: "test_origin",
    }
    // handle option
    option := tool.GetImplSpecificOptions(baseOption, opts...)
    return &Result{
        Msg:  option.Field1,
    }, nil
}
```



~~**附录：使用 NewTool() 构建。**~~

~~适合简单的工具实现，通过定义工具信息和处理函数来创建 Tool。但需要在 ToolInfo 中手动定义参数信息（ParamsOneOf），和实际的参数结构（TodoAddParams）是分开定义的。这样不仅~~~~**造成了代码的冗余，而且在参数发生变化时需要同时修改两处地方，容易导致不一致，维护起来也比较麻烦**~~~~。~~

~~当一个函数满足下面这种函数签名时，就可以用 NewTool 把其变成一个 InvokableTool ：~~

```go
type InvokeFunc[T, D any] func(ctx context.Context, input T) (output D, err error)
```

~~NewTool 的方法如下：~~

```go
// 代码见: github.com/cloudwego/eino/components/tool/utils/invokable_func.go
func NewTool[T, D any](desc *schema.ToolInfo, i InvokeFunc[T, D], opts ...Option) tool.InvokableTool
```

~~同理 NewStreamTool 可创建 StreamableTool~~

~~`AddUser()`~~~~ example:~~

```go
import (
    "context"

    "github.com/cloudwego/eino/components/tool"
    "github.com/cloudwego/eino/components/tool/utils"
    "github.com/cloudwego/eino/schema"
)

// 处理函数
func AddTodoFunc(_ context.Context, params *TodoAddParams) (string, error) {
    // Mock处理逻辑
    return `{"msg": "add todo success"}`, nil
}

func getAddTodoTool() tool.InvokableTool {
    // 工具信息
    info := &schema.ToolInfo{
        Name: "add_todo",
        Desc: "Add a todo item",
        ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
            "content": {
                Desc:     "The content of the todo item",
                Type:     schema.String,
                Required: true,
            },
            "started_at": {
                Desc: "The started time of the todo item, in unix timestamp",
                Type: schema.Integer,
            },
            "deadline": {
                Desc: "The deadline of the todo item, in unix timestamp",
                Type: schema.Integer,
            },
        }),
    }

    // 使用NewTool创建工具
    return utils.NewTool(info, AddTodoFunc)
}
```



### **2. 直接实现 InvokableTool interface**

对于需要更多自定义逻辑的场景，可以通过手动实现 InvokableTool interface 来创建：

```go
import (
    "context"

    "github.com/cloudwego/eino/components/tool"
    "github.com/cloudwego/eino/schema"
)

type ListTodoTool struct {}

func (lt *ListTodoTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
    return &schema.ToolInfo{
        Name: "list_todo",
        Desc: "List all todo items",
        ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
            "finished": {
                Desc:     "filter todo items if finished",
                Type:     schema.Boolean,
                Required: false,
            },
        }),
    }, nil
}

func (lt *ListTodoTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
    // Mock调用逻辑
    return `{"todos": [{"id": "1", "content": "在2024年12月10日之前完成Eino项目演示文稿的准备工作", "started_at": 1717401600, "deadline": 1717488000, "done": false}]}`, nil
}
```

备注：由于大模型给出的 function call 参数始终是一个 string，对应到 Eino 框架中，tool 的调用参数入参也就是一个序列化成 string 的 json。因此，这种方式需要开发者自行处理参数的反序列化，并且调用的结果也用 string 的方式返回。



# 二、 Lambda组件

> [https://www.cloudwego.io/zh/docs/eino/core_modules/components/lambda_guide/](https://www.cloudwego.io/zh/docs/eino/core_modules/components/lambda_guide/)

自定义组件-Lambda组件：支持**自定义的函数逻辑。**

是Eino 中最基础的组件类型，它允许用户在工作流中嵌入。Lambda 组件底层是由输入输出是否流对应 4 种交互模式（4种函数）: Invoke、Stream、Collect、Transform。



**Lambda组件定义**：核心是 `Lambda` 结构体，封装了用户提供的 Lambda 函数，用户可通过构建方法创建一个 Lambda 组件：

```go
// eino/compose/types_lambda.go

// Lambda is the node that wraps the user provided lambda function.
// It can be used as a node in Graph or Chain (include Parallel and Branch).
// Create a Lambda by using AnyLambda/InvokableLambda/StreamableLambda/CollectableLambda/TransformableLambda.
// eg.
//
//	lambda := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
//		return input, nil
//	})
type Lambda struct {
	executor *composableRunnable
}
```



**Lambda组件的构建方法：**

- Eino的组件接口的统一规范:一个组件的可调用方法需要有3个入参和2个出参：`func (ctx, input, …option) (output, error)`
- 但在使用 Lambda 的场景中通常使用匿名函数
**不使用自定义 Option:**

```go
// input 和 output 类型为自定义的任何类型
lambda := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
    // some logic
})

// input 可以是任意类型；output 必须是 *schema.StreamReader[O]，其中 O 可以是任意类型
lambda := compose.StreamableLambda(func(ctx context.Context, input string) (output *schema.StreamReader[string], err error) {
    // some logic
})

// input 必须是 *schema.StreamReader[I]，其中 I 可以是任意类型；output 可以是任意类型
lambda := compose.CollectableLambda(func(ctx context.Context, input *schema.StreamReader[string]) (output string, err error) {
    // some logic
})

// input 必须是 *schema.StreamReader[I]，其中 I 可以是任意类型；output 必须是 *schema.StreamReader[O]，其中 O 可以是任意类型
lambda := compose.TransformableLambda(func(ctx context.Context, input *schema.StreamReader[string]) (output *schema.StreamReader[string], err error) {
    // some logic
})
```

**使用自定义 Option:**

```go
type Options struct {
    Field1 string
}
type MyOption func(*Options)

lambda := compose.InvokableLambdaWithOption(
    func(ctx context.Context, input string, opts ...MyOption) (output string, err error) {
        // 处理 opts
        // some logic
    }
)

// AnyLambda 允许同时实现多种交互模式的 Lambda 函数类型：
type Options struct {
    Field1 string
}

type MyOption func(*Options)

// input 和 output 类型为自定义的任何类型
lambda, err := compose.AnyLambda(
    // Invoke 函数
    func(ctx context.Context, input string, opts ...MyOption) (output string, err error) {
        // some logic
    },
    // Stream 函数
    func(ctx context.Context, input string, opts ...MyOption) (output *schema.StreamReader[string], err error) {
        // some logic
    },
    // Collect 函数
    func(ctx context.Context, input *schema.StreamReader[string], opts ...MyOption) (output string, err error) {
        // some logic
    },
    // Transform 函数
    func(ctx context.Context, input *schema.StreamReader[string], opts ...MyOption) (output *schema.StreamReader[string], err error) {
        // some logic
    },
)
```



**Graph 中使用：**

```go
graph := compose.NewGraph[string, *MyStruct]()
graph.AddLambdaNode(
    "node1",
    compose.InvokableLambda(func(ctx context.Context, input string) (*MyStruct, error) {
        // some logic
    }),
)
```

**Chain 中使用：**

```go
chain := compose.NewChain[string, string]()
chain.AppendLambda(compose.InvokableLambda(func(ctx context.Context, input string) (string, error) {
    // some logic
}))
```



**两个内置的 Lambda：**

```go
// ToList 是一个内置的 Lambda，用于将单个输入元素转换为包含该元素的切片（数组）：
// 创建一个 ToList Lambda
lambda := compose.ToList[*schema.Message]()

// 在 Chain 中使用
chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
chain.AppendChatModel(chatModel)  // chatModel 返回 *schema.Message
chain.AppendLambda(lambda)        // 将 *schema.Message 转换为 []*schema.Message
```



```go

// MessageParser 是一个内置的 Lambda，用于将 JSON 消息（通常由 LLM 生成）解析为指定的结构体：
// 定义解析目标结构体
type MyStruct struct {
    ID int `json:"id"`
}

// 创建解析器
parser := schema.NewMessageJSONParser[*MyStruct](&schema.MessageJSONParseConfig{
    ParseFrom: schema.MessageParseFromContent,
    ParseKeyPath: "", // 如果仅需要 parse 子字段，可用 "key.sub.grandsub"
})

// 创建解析 Lambda
parserLambda := compose.MessageParser(parser)

// 在 Chain 中使用
chain := compose.NewChain[*schema.Message, *MyStruct]()
chain.AppendLambda(parserLambda)

// 使用示例
runner, err := chain.Compile(context.Background())
parsed, err := runner.Invoke(context.Background(), &schema.Message{
    Content: `{"id": 1}`,
})
// parsed.ID == 1

// MessageParser 支持从消息内容（Content）或工具调用结果（ToolCall）中解析数据，这在意图识别等场景中常用：

// 从工具调用结果解析
parser := schema.NewMessageJSONParser[*MyStruct](&schema.MessageJSONParseConfig{
    ParseFrom: schema.MessageParseFromToolCall,
})
```



# **三、eino-ext：Eino Dev**

安装（当cursor/trae搜不到时）：

1. 从vscode 插件市场搜索**Eino Dev，**下载.vsix 文件；
1. 打开Cursor，按 Cmd+Shift+P (macOS) 打开命令面板，输入 Install from VSIX 安装；
1. 底部找到 Eino Dev
备注：golang版本已修复请求1分钟超时问题，vscode版本不存在该问题；



## 功能1：可视化编排&代码生成

拖拽组件 实现Graph的编排并生成代码。支持导入导出。

orchestration  /ˌɔːrkɪ'streɪʃn/ n. 管弦乐编曲；和谐的结合



这里直接导入看效果：[https://github.com/cloudwego/eino-examples/blob/764d04fbf360878c5109d024239b2432caa30b47/quickstart/eino_assistant/eino/knowledge_indexing.json](https://github.com/cloudwego/eino-examples/blob/764d04fbf360878c5109d024239b2432caa30b47/quickstart/eino_assistant/eino/knowledge_indexing.json)。

编排组件包括：

- **Graph（图）**
    - **Node（节点）**
        - **Component（组件）**
            - **Slot（插槽）**
![](/images/25d24637-29b5-80e4-b275-e8ec8dbc3df3/image_29324637-29b5-80b6-9bf8-dbacb2817b45.jpg)





更多文档：

- [Eino Dev 可视化编排插件功能指南](https://www.cloudwego.io/zh/docs/eino/core_modules/devops/visual_orchestration_plugin_guide/#%E7%BC%96%E6%8E%92%E7%BB%84%E4%BB%B6%E4%BB%8B%E7%BB%8D)
- [可视化开发](https://www.cloudwego.io/zh/docs/eino/overview/bytedance_eino_practice/#%E5%8F%AF%E8%A7%86%E5%8C%96%E5%BC%80%E5%8F%91-1)
## 功能2：可视化Debug

- 运行 [源码地址](https://github.com/cloudwego/eino-examples/blob/3a94b9ab0db133907636c07ef1e3cf267551725c/devops/debug/main.go) 
![](/images/25d24637-29b5-80e4-b275-e8ec8dbc3df3/image_29424637-29b5-80ea-92b5-c3f3f12db3c7.jpg)

- Eino Dev 配置调试地址，选择需要调试的Graph
- 点击 Test Run。默认从star节点开始执行，可以点击可视化graph从任意节点开始执行，看到每个节点的input、output。（类似 AI Agent版本的trace）
![](/images/25d24637-29b5-80e4-b275-e8ec8dbc3df3/image_29424637-29b5-80ea-820f-f63269084836.jpg)



![](/images/25d24637-29b5-80e4-b275-e8ec8dbc3df3/image_29424637-29b5-80c4-9238-c4e8884c9413.jpg)

![](/images/25d24637-29b5-80e4-b275-e8ec8dbc3df3/image_29424637-29b5-8002-941f-c68efa4289b9.jpg)



- 高级功能：**指定 interface 字段的实现类型**
    ![](/images/25d24637-29b5-80e4-b275-e8ec8dbc3df3/image_29524637-29b5-80e1-b502-e3bf13cc7fba.jpg)


更多文档：[Eino Dev 可视化调试插件功能指南](https://www.cloudwego.io/zh/docs/eino/core_modules/devops/visual_debug_plugin_guide/)



