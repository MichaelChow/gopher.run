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

是Eino 中最基础的组件类型，它允许用户在工作流中嵌入自定义的函数逻辑。Lambda 组件底层是由输入输出是否流所形成的 4 种运行函数组成，对应 4 种交互模式: Invoke、Stream、Collect、Transform。

用户构建 Lambda 时可实现其中的一种或多种，框架会根据一定的规则进行转换。

