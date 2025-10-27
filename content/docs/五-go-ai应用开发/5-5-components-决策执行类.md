---
title: "5.5 components-决策执行类"
date: 2025-08-28T13:33:00Z
draft: false
weight: 5005
---

# 5.5 components-决策执行类



## 一、tool

> [https://www.cloudwego.io/zh/docs/eino/core_modules/components/tools_node_guide/](https://www.cloudwego.io/zh/docs/eino/core_modules/components/tools_node_guide/)
[https://www.cloudwego.io/zh/docs/eino/core_modules/components/tools_node_guide/how_to_create_a_tool/](https://www.cloudwego.io/zh/docs/eino/core_modules/components/tools_node_guide/how_to_create_a_tool/)

一个用于扩展模型能力的组件，它允许模型调用外部工具来完成特定的任务。

应用场景：

- 让模型能够获取实时信息（如搜索引擎、天气查询等）
- 使模型能够执行特定的操作（如数据库操作、API 调用等）
- 扩展模型的能力范围（如数学计算、代码执行等）
- 与外部系统集成（如知识库查询、插件系统等）


### **核心接口设计 (interface.go)**

**基础工具接口:**

```go
type BaseTool interface {
	Info(ctx context.Context) (*schema.ToolInfo, error)
}
```

**可调用工具接口:**

```go
type InvokableTool interface {
	BaseTool

	// InvokableRun call function with arguments in JSON format
	InvokableRun(ctx context.Context, argumentsInJSON string, opts ...Option) (string, error)
}
```

**流式工具接口:**

```go
type StreamableTool interface {
	BaseTool

	StreamableRun(ctx context.Context, argumentsInJSON string, opts ...Option) (*schema.StreamReader[string], error)
}
```

**设计特点：**

- **分层设计**：从BaseTool到具体实现，支持渐进式功能扩展
- **JSON接口**：统一使用JSON格式进行参数传递和结果返回
- **流式支持**：支持流式响应，适合大数据量处理
- **选项模式**：通过可变参数opts支持灵活配置
### **选项模式设计 (option.go)**

**泛型选项支持**

```go
func WrapImplSpecificOptFn[T any](optFn func(*T)) Option {
	return Option{
		implSpecificOptFn: optFn,
	}
}

func GetImplSpecificOptions[T any](base *T, opts ...Option) *T {
	if base == nil {
		base = new(T)
	}

	for i := range opts {
		opt := opts[i]
		if opt.implSpecificOptFn != nil {
			optFn, ok := opt.implSpecificOptFn.(func(*T))
			if ok {
				optFn(base)
			}
		}
	}

	return base
}
```

**设计亮点：**

- **类型安全**：通过泛型确保选项的类型安全
- **延迟应用**：选项函数在需要时才被调用
- **默认值支持**：支持基础选项的默认值设置
### **工具函数系统 (utils/)**

**通用工具函数 (common.go)**

```go
func marshalString(resp any) (string, error) {
	if rs, ok := resp.(string); ok {
		return rs, nil
	}
	return sonic.MarshalString(resp)
}
```

**可调用工具函数 (invokable_func.go)**

```go
// InvokeFunc is the function type for the tool.
type InvokeFunc[T, D any] func(ctx context.Context, input T) (output D, err error)

// OptionableInvokeFunc is the function type for the tool with tool option.
type OptionableInvokeFunc[T, D any] func(ctx context.Context, input T, opts ...tool.Option) (output D, err error)

// InferTool creates an InvokableTool from a given function by inferring the ToolInfo from the function's request parameters.
func InferTool[T, D any](toolName, toolDesc string, i InvokeFunc[T, D], opts ...Option) (tool.InvokableTool, error) {
	ti, err := goStruct2ToolInfo[T](toolName, toolDesc, opts...)
	if err != nil {
		return nil, err
	}

	return NewTool(ti, i, opts...), nil
}
```

**核心功能：**

- **类型推断**：自动从Go结构体推断工具信息
- **函数包装**：将普通函数包装成工具接口
- **Schema生成**：自动生成OpenAPI Schema


**流式工具函数 (streamable_func.go)**

```go
// StreamFunc is the function type for the streamable tool.
type StreamFunc[T, D any] func(ctx context.Context, input T) (output *schema.StreamReader[D], err error)

// InferStreamTool creates an StreamableTool from a given function by inferring the ToolInfo from the function's request parameters
func InferStreamTool[T, D any](toolName, toolDesc string, s StreamFunc[T, D], opts ...Option) (tool.StreamableTool, error) {
	ti, err := goStruct2ToolInfo[T](toolName, toolDesc, opts...)
	if err != nil {
		return nil, err
	}

	return NewStreamTool(ti, s, opts...), nil
}
```



**创建选项 (create_options.go)**

```go
type toolOptions struct {
	um UnmarshalArguments
	m  MarshalOutput
	sc SchemaCustomizerFn
}

// WithUnmarshalArguments wraps the unmarshal arguments option.
func WithUnmarshalArguments(um UnmarshalArguments) Option {
	return func(o *toolOptions) {
		o.um = um
	}
}

// WithMarshalOutput wraps the marshal output option.
func WithMarshalOutput(m MarshalOutput) Option {
	return func(o *toolOptions) {
		o.m = m
	}
}

// WithSchemaCustomizer sets a user-defined schema customizer for inferring tool parameter from tagged go struct.
func WithSchemaCustomizer(sc SchemaCustomizerFn) Option {
	return func(o *toolOptions) {
		o.sc = sc
	}
}
```

**配置选项：**

- **参数解析**：自定义参数的反序列化逻辑
- **结果序列化**：自定义结果的序列化逻辑
- **Schema定制**：自定义OpenAPI Schema的生成逻辑


**错误处理 (error_handler.go)**

```go
type ErrorHandler func(context.Context, error) string

// WrapToolWithErrorHandler wraps any BaseTool with custom error handling.
func WrapToolWithErrorHandler(t tool.BaseTool, h ErrorHandler) tool.BaseTool {
	ih := &infoHelper{info: t.Info}
	var s tool.StreamableTool
	if st, ok := t.(tool.StreamableTool); ok {
		s = st
	}
	if it, ok := t.(tool.InvokableTool); ok {
		if s == nil {
			return WrapInvokableToolWithErrorHandler(it, h)
		} else {
			return &combinedErrorWrapper{
				infoHelper: ih,
				errorHelper: &errorHelper{
					i: it.InvokableRun,
					h: h,
				},
				streamErrorHelper: &streamErrorHelper{
					s: s.StreamableRun,
					h: h,
				},
			}
		}
	}
	if s != nil {
		return WrapStreamableToolWithErrorHandler(s, h)
	}
	return t
}
```

**错误处理特性：**

- **统一处理**：为所有工具类型提供统一的错误处理
- **类型感知**：自动检测工具类型并应用相应的错误处理
- **组合支持**：支持同时实现多个接口的工具
### **回调机制设计 (callback_extra.go)**

**回调数据结构**

```go
type CallbackInput struct {
	// ArgumentsInJSON is the arguments in json format for the tool.
	ArgumentsInJSON string
	// Extra is the extra information for the tool.
	Extra map[string]any
}

type CallbackOutput struct {
	// Response is the response for the tool.
	Response string
	// Extra is the extra information for the tool.
	Extra map[string]any
}
```

**类型转换工具**

```go
func ConvCallbackInput(src callbacks.CallbackInput) *CallbackInput {
	switch t := src.(type) {
	case *CallbackInput:
		return t
	case string:
		return &CallbackInput{ArgumentsInJSON: t}
	default:
		return nil
	}
}
```



## 二、 Lambda组件

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



