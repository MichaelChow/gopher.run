---
title: "5.3 components对话处理类"
date: 2025-08-13T01:45:00Z
draft: false
weight: 5003
---

# 5.3 components对话处理类

**components**/kəm'ponənt/ n. 零组件

大模型应用开发和传统应用开发最显著的区别在于大模型所具备的两大核心能力：

1. **基于语义的文本处理能力**：能够理解和生成人类语言，处理非结构化的内容语义关系；
1. **智能决策能力**：能够基于上下文进行推理和判断，做出相应的行为决策。


这两项核心能力进而催生了如下三种主要的AI应用类型：

1. **对话处理类（Chat）**：处理用户输入并生成相应回答。ChatTemplate、ChatModel
1. **文本语义处理类（RAG）**：对文本文档进行语义化处理、存储和检索。Document.Loader 、Document.Transformer Embedding Indexer Retriever
1. **决策执行类（Tool call）**：基于上下文做出决策并调用相应工具。ToolsNode Lambda。


Eino基于这三种模式将这些常用能力抽象为可复用的「组件」（Components）。

## 一、prompt

Prompt 组件是一个用于处理和格式化提示模板的组件。它的主要作用是将用户提供的变量值填充到预定义的消息模板中，生成用于与语言模型交互的标准消息格式。这个组件可用于以下场景：

- 构建结构化的系统提示
- 处理多轮对话的模板 (包括 history)
- 实现可复用的提示模式
> [https://www.cloudwego.io/zh/docs/eino/core_modules/components/chat_template_guide/](https://www.cloudwego.io/zh/docs/eino/core_modules/components/chat_template_guide/)

### 核心接口、结构设计

- **接口抽象清晰**：ChatTemplate定义了明确的行为契约
```go
type ChatTemplate interface {
    Format(ctx context.Context, vs map[string]any, opts ...Option) ([]*schema.Message, error)
}

// 单一职责：专注于模板到消息的转换
// 依赖注入：通过context传递运行时信息；ctx：上下文对象，用于传递请求级别的信息，同时也用于传递 Callback Manager
// 灵活配置：支持变量映射和选项配置；vs：变量值映射，用于填充模板中的占位符；opts：可选参数，用于配置格式化行为；
// 错误处理：明确的错误返回机制

// - []*schema.Message：格式化后的消息列表
```

- **实现简洁高效**：DefaultChatTemplate通过组合模式实现复杂功能
```go
type DefaultChatTemplate struct {
    templates  []schema.MessagesTemplate  // 模板集合
    formatType schema.FormatType          // 格式化类型
}

// 组合模式：通过组合多个MessagesTemplate实现复杂模板
// 策略模式：支持多种格式化类型（FString、Jinja2、GoTemplate等）
// 不可变性：构造后模板和格式类型不可修改
```

- DefaultChatTemplate的**Format方法**：
```go
func (t *DefaultChatTemplate) Format(ctx context.Context, vs map[string]any, _ ...Option) (result []*schema.Message, err error) {
    // 1. 回调机制初始化
    ctx = callbacks.EnsureRunInfo(ctx, t.GetType(), components.ComponentOfPrompt)
    ctx = callbacks.OnStart(ctx, &CallbackInput{...})  // OnStart
    
    // 2. 错误处理defer
    defer func() {
        if err != nil {
            _ = callbacks.OnError(ctx, err)  // OnError
        }
    }()
    
    // 3. 模板遍历处理
    result = make([]*schema.Message, 0, len(t.templates))
    for _, template := range t.templates {
        msgs, err := template.Format(ctx, vs, t.formatType)
        if err != nil {
            return nil, err
        }
        result = append(result, msgs...)
    }
    
    // 4. 成功回调
    _ = callbacks.OnEnd(ctx, &CallbackOutput{...})  // OnEnd
    return result, nil
}
```

- DefaultChatTemplate的**构造函数FromMessages**：把多个 message 变成一个 chat template
```go
func FromMessages(formatType schema.FormatType, templates ...schema.MessagesTemplate) *DefaultChatTemplate {
    return &DefaultChatTemplate{
        templates:  templates,
        formatType: formatType,
    }
}
```

- **集成性强**：无缝集成到Chain和Graph编排系统
```go
// **Chain集成**
chain := compose.NewChain[map[string]any, []*schema.Message]()
chain.AppendChatTemplate(template)

// Graph集成
graph.AddChatTemplateNode("template", template)
```

- **可观测性好**：完整的回调机制支持监控和调试
    - **切面编程**：通过callbacks实现横切关注点
    - **生命周期管理**：OnStart、OnEnd、OnError完整覆盖
    - **组件标识**：通过ComponentOfPrompt标识组件类型
- **扩展性优秀**：支持多种格式化策略和自定义模板


### option设计模式

**Option接口:**

```go
type Option struct {
    implSpecificOptFn any
}
```

**WrapImplSpecificOptFn:**

```go
func WrapImplSpecificOptFn[T any](optFn func(*T)) Option {
    return Option{
        implSpecificOptFn: optFn,
    }
}
```

**GetImplSpecificOptions**：

```go
func GetImplSpecificOptions[T any](base *T, opts ...Option) *T {
    if base == nil {
        base = new(T)
    }
    
    for i := range opts {
        opt := opts[i]
        if opt.implSpecificOptFn != nil {
            s, ok := opt.implSpecificOptFn.(func(*T))
            if ok {
                s(base)
            }
        }
    }
    
    return base
}
```

**函数式选项模式（Functional Options Pattern）:**

```go
// 选项构造函数
func WithUserID(uid int64) Option {
    return WrapImplSpecificOptFn[implOption](func(i *implOption) {
        i.userID = uid
    })
}

func WithName(n string) Option {
    return WrapImplSpecificOptFn[implOption](func(i *implOption) {
        i.name = n
    })
}

// 链式调用：支持多个选项的组合使用
// 默认值处理：自动处理nil值，提供默认初始化
// 类型安全：通过泛型确保选项类型匹配
```

### callback机制

**核心接口设计**

```go
// 回调输入结构
type CallbackInput struct {
    Variables  map[string]any                    // 模板变量
    Templates  []schema.MessagesTemplate         // 模板集合
    Extra      map[string]any                    // 扩展信息
}

// 回调输出结构
type CallbackOutput struct {
    Result     []*schema.Message                 // 格式化结果
    Templates  []schema.MessagesTemplate         // 模板集合
    Extra      map[string]any                    // 扩展信息
}
```

**回调处理器设计:**

**PromptCallbackHandler接口:**

```go
type PromptCallbackHandler struct {
    OnStart func(ctx context.Context, runInfo *callbacks.RunInfo, input *prompt.CallbackInput) context.Context
    OnEnd   func(ctx context.Context, runInfo *callbacks.RunInfo, output *prompt.CallbackOutput) context.Context
    OnError func(ctx context.Context, runInfo *callbacks.RunInfo, err error) context.Context
}

// 函数式接口：每个回调都是独立的函数
// 类型安全：使用具体的类型而非interface{}
// 可选实现：支持部分回调方法的实现
```

Handler**构建器模式:**

```go
// 使用HandlerBuilder构建回调处理器
handler := callbacks.NewHandlerBuilder().
    OnStartFn(func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
        // 处理开始事件
        return ctx
    }).
    OnEndFn(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
        // 处理结束事件
        return ctx
    }).
    OnErrorFn(func(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
        // 处理错误事件
        return ctx
    }).
    Build()
```



### example

**基础用法**：

```go
// 创建模板
template := prompt.FromMessages(
    schema.FString,
    &schema.Message{Content: "Hello, {name}!"},
    &schema.Message{Content: "How are you?"}
)

// 格式化
msgs, err := template.Format(context.Background(), map[string]any{"name": "World"})
```

**高级组合:**

```go
// 系统消息模板
systemTemplate := &schema.Message{Content: "You are a helpful assistant."}

// 用户消息模板
userTemplate := &schema.Message{Content: "User: {user_input}"}

// 组合模板
chatTemplate := prompt.FromMessages(schema.Jinja2, systemTemplate, userTemplate)
```

## 二、model

### **接口层次结构**

```go
// 基础接口 - 提供核心能力
type BaseChatModel interface {
    Generate(ctx context.Context, input []*schema.Message, opts ...Option) (*schema.Message, error)
    Stream(ctx context.Context, input []*schema.Message, opts ...Option) (*schema.StreamReader[*schema.Message], error)
}

// 扩展接口 - 增加工具调用能力
type ToolCallingChatModel interface {
    BaseChatModel
    WithTools(tools []*schema.ToolInfo) (ToolCallingChatModel, error)
}
```



**策略模式 + 装饰器模式**：

- BaseChatModel 定义了基础策略接口
- ToolCallingChatModel 通过装饰器模式扩展功能
- 避免了继承的复杂性，保持了接口的纯净性
**函数式选项模式**：

```go
// 灵活的配置选项系统
type Option struct {
    apply func(opts *Options)
    implSpecificOptFn any
}

// 使用示例
model.Generate(ctx, messages, 
    WithTemperature(0.7),
    WithMaxTokens(1000),
    WithTools(tools))
```



### **回调机制设计**

1. **回调数据结构**
```go
// 输入数据结构
type CallbackInput struct {
    Messages []*schema.Message    // 输入消息
    Tools []*schema.ToolInfo      // 可用工具
    Config *Config                // 模型配置
    Extra map[string]any          // 扩展信息
}

// 输出数据结构  
type CallbackOutput struct {
    Message *schema.Message       // 生成的消息
    Config *Config                // 使用的配置
    TokenUsage *TokenUsage        // Token使用统计
    Extra map[string]any          // 扩展信息
}
```

**2. 类型转换机制**

框架提供了智能的类型转换函数，体现了**类型安全**的设计思想：

```go
// 自动类型转换，支持多种输入类型
func ConvCallbackInput(src callbacks.CallbackInput) *CallbackInput {
    switch t := src.(type) {
    case *CallbackInput:           // 直接类型匹配
        return t
    case []*schema.Message:        // 从消息数组转换
        return &CallbackInput{Messages: t}
    default:
        return nil
    }
}
```



### **核心设计哲学**

**1. 不可变性设计**

```go
// ToolCallingChatModel 返回新实例，避免状态污染
WithTools(tools []*schema.ToolInfo) (ToolCallingChatModel, error)
```

1. **并发安全:**
废弃了 BindTools 方法，避免并发问题：

- **状态污染**：直接修改当前实例的内部状态
- **非原子操作**：BindTools 和 Generate 之间不是原子操作
- **竞态条件**：多个goroutine同时调用时会产生数据竞争
```go
// Deprecated: Please use ToolCallingChatModel interface instead, which provides a safer way to bind tools
// without the concurrency issues and tool overwriting problems that may arise from the BindTools method.
type ChatModel interface {
	BaseChatModel

	// BindTools bind tools to the model.
	// BindTools before requesting ChatModel generally.
	// notice the non-atomic problem of BindTools and Generate.
	BindTools(tools []*schema.ToolInfo) error
}

// ToolCallingChatModel extends BaseChatModel with tool calling capabilities.
// It provides a WithTools method that returns a new instance with
// the specified tools bound, avoiding state mutation and concurrency issues.
type ToolCallingChatModel interface {
	BaseChatModel

	// WithTools returns a new ToolCallingChatModel instance with the specified tools bound.
	// This method does not modify the current instance, making it safer for concurrent use.
	WithTools(tools []*schema.ToolInfo) (ToolCallingChatModel, error)
}
```

**具体场景示例：**

```go
// 危险的使用方式
model := &MyChatModel{}
go func() {
    model.BindTools(tools1)  // goroutine 1 绑定工具1
}()

go func() {
    model.BindTools(tools2)  // goroutine 2 绑定工具2
}()

// 此时调用 Generate，无法确定使用的是哪个工具集
result, err := model.Generate(ctx, messages)
```



**WithTools 方法的安全设计：**

```go
// 安全方法：返回新实例，不修改原状态
WithTools(tools []*schema.ToolInfo) (ToolCallingChatModel, error)
```

**安全特性：**

1. **不可变性**：原实例状态保持不变
1. **实例隔离**：每个调用返回独立的实例
1. **无状态共享**：避免了goroutine间的状态竞争
```go
// 安全的使用方式
baseModel := &MyChatModel{}

// 每个goroutine获得独立的模型实例
go func() {
    model1, _ := baseModel.WithTools(tools1)
    result1, _ := model1.Generate(ctx, messages1)
}()

go func() {
    model2, _ := baseModel.WithTools(tools2)
    result2, _ := model2.Generate(ctx, messages2)
}()

// 原实例保持不变，可以继续使用
result3, _ := baseModel.Generate(ctx, messages3)
```



**3. 扩展性设计**

```go
// 支持实现特定的选项扩展
func WrapImplSpecificOptFn[T any](optFn func(*T)) Option {
    return Option{implSpecificOptFn: optFn}
}
```



### **最佳实践建议**

**1. 模型使用**

```go
// 推荐：使用 ToolCallingChatModel
toolModel, err := baseModel.WithTools(tools)
if err != nil {
    return err
}

// 避免：直接使用已废弃的 ChatModel
// model.BindTools(tools) // 不推荐
```

1. **回调处理**
```go
// 创建类型安全的回调处理器
handler := &ModelCallbackHandler{
    OnStart: func(input *CallbackInput) {
        // 处理开始事件
    },
    OnEnd: func(output *CallbackOutput) {
        // 处理结束事件
    },
}
```

**3. 选项配置**

```go
// 使用函数式选项模式
opts := GetCommonOptions(nil,
    WithTemperature(0.7),
    WithMaxTokens(1000),
    WithTools(tools),
)
```



