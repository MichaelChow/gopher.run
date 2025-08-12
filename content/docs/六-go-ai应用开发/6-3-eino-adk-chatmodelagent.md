---
title: "6.3 Eino ADK: ChatModelAgent"
date: 2025-07-03T13:48:00Z
draft: false
weight: 6003
---

# 6.3 Eino ADK: ChatModelAgent



> [https://github.com/cloudwego/eino-examples](https://github.com/cloudwego/eino-examples): 提供了实用的示例来帮助上手基于Eino的AI应用开发

## 一、**ChatModelAgent**

> doc：[https://www.cloudwego.io/zh/docs/eino/core_modules/eino_adk/eino-adk-agent-实现/eino-adk-chatmodelagent/](https://www.cloudwego.io/zh/docs/eino/core_modules/eino_adk/eino-adk-agent-%E5%AE%9E%E7%8E%B0/eino-adk-chatmodelagent/)

> code：[github.com/cloudwego/eino-examples/adk/intro/chatmodel](http://github.com/cloudwego/eino-examples/adk/intro/chatmodel)



**ChatModelAgent：**

一个核心预构建的Agent，封装了ChatModel、tool。

```go
// eino/adk/chatmodel.go
type ChatModelAgent struct {
	name        string
	description string
	instruction string

	model       model.ToolCallingChatModel
	toolsConfig ToolsConfig  

	genModelInput GenModelInput

	outputKey string
	maxStep   int

	subAgents   []Agent
	parentAgent Agent

	disallowTransferToParent bool

	exit tool.BaseTool

	// runner
	once   sync.Once
	run    runFunc
	frozen uint32
}
```



**ChatModelAgentConfig：**

<!-- 列布局开始 -->

```go
// eino/adk/chatmodel.go
type ChatModelAgentConfig struct {
	Name        string
	Description string
	Instruction string

	Model model.ToolCallingChatModel

	**ToolsConfig** ToolsConfig

	// optional
	**GenModelInput** GenModelInput

	// Exit tool. Optional, defaults to nil, which will generate an Exit Action.
	// The built-in implementation is 'ExitTool'
	Exit tool.BaseTool

	// optional
	OutputKey string

	MaxStep int
}
```


---

![](/images/22524637-29b5-80e9-84e8-ff0346e83856/image_24824637-29b5-80ce-99e3-fdb27b582776.jpg)

<!-- 列布局结束 -->

**ToolsConfig：**

复用了 Eino Graph的compose.ToolsNodeConfig，详细参考：[Eino: ToolsNode&Tool 使用说明](https://www.cloudwego.io/zh/docs/eino/core_modules/components/tools_node_guide)。并额外提供了 ReturnDirectly 配置，ChatModelAgent 调用配置在 ReturnDirectly 中的 Tool 后会直接退出。

为 ChatModelAgent 配置了 ToolsConfig 后，它在内部的执行流程就遵循了 ReAct 模式：调用 ChatModel（Reason）、chatModel 返回工具调用请求（Action）、ChatModelAgent 执行工具（Act）

执行循环直到 ChatModel 判断不需要调用 Tool 结束。

**当没有配置工具时，ChatModelAgent 退化为一次 ChatModel 调用。**

```go
// github.com/cloudwego/eino/adk/chatmodel.go

type ToolsConfig struct {
    compose.ToolsNodeConfig

    // Names of the tools that will make agent return directly when the tool is called.
    // When multiple tools are called and more than one tool is in the return directly list, only the first one will be returned.
    ReturnDirectly map[string]bool
}
```



**GenModelInput:**

Agent 被调用时会使用该方法生成 ChatModel 的初始输入：

```go
type GenModelInput func(ctx context.Context, instruction string, input *AgentInput) ([]Message, error)
```

Agent 提供了默认的 GenModelInput 方法：

1. 将 Instruction 作为 system message 加到 AgentInput.Messages 前
1. 以 SessionValues 为 variables 渲染 1 中得到的 message list


**OutputKey：**

配置后 Agent 产生的最后一个 message 会被以设置的 OutputKey 为 key 添加到 SessionValues 中。



**Exit：**

效果类似 ToolReturnDirectly。当 chatModel 调用这个工具后并执行后，ChatModelAgent 将直接退出。

```go
// github.com/cloudwego/eino/adk/chatmodel.go

type ExitTool struct{}

func (et ExitTool) Info(_ context.Context) (*schema.ToolInfo, error) {
    return ToolInfoExit, nil
}

func (et ExitTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
    type exitParams struct {
       FinalResult string `json:"final_result"`
    }

    params := &exitParams{}
    err := sonic.UnmarshalString(argumentsInJSON, params)
    if err != nil {
       return "", err
    }

    err = SendToolGenAction(ctx, "exit", NewExitAction())
    if err != nil {
       return "", err
    }

    return params.FinalResult, nil
}
```



**Transfer:**

使用 SetSubAgents 为 ChatModelAgent 设置父或子 Agent 后，ChatModelAgent 会增加一个 Transfer Tool，并且在 prompt 中指示 ChatModel 在需要 transfer 时调用这个 Tool 并以 transfer 目标 AgentName 作为 Tool 输入。在此工具被调用后，Agent 会产生 TransferAction 并退出。



**AgentTool:**

方便地将 Eino ADK Agent 转化为 Tool 供 ChatModelAgent 调用:

```go
// github.com/cloudwego/eino/adk/agent_tool.go

func NewAgentTool(_ context.Context, agent Agent, options ...AgentToolOption) tool.BaseTool
```



如把之前创建的 `BookRecommendAgent` 转换为 Tool

```go
bookRecommender := NewBookRecommendAgent()
bookRecommendeTool := NewAgentTool(ctx, bookRecommender)

// other agent
a, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
    // xxx
    ToolsConfig: adk.ToolsConfig{
        ToolsNodeConfig: compose.ToolsNodeConfig{
            Tools: []tool.BaseTool{bookRecommendeTool},
        },
    },
})
```



**Interrupt&Resume：**

复用了 Eino Graph 的 Interrupt&Resume 能力。

```go
// github.com/cloudwego/eino/adk/interrupt.go

func NewInterruptAndRerunErr(extra any) error
```

定义 ToolOption 来在恢复时传递新输入：(非必须，实践时也可以根据 context、闭包等其他方式传递新输入)

```go
import (
    "github.com/cloudwego/eino/components/tool"
)

type askForClarificationOptions struct {
    NewInput *string
}

func WithNewInput(input string) tool.Option {
    return tool.WrapImplSpecificOptFn(func(t *askForClarificationOptions) {
       t.NewInput = &input
    })
}
```

工具 `ask_for_clarification` 使用了 Interrupt&Resume 能力来实现向用户“询问”。

## **二、example: 图书推荐Agent**

根据用户的输入推荐相关图书。

**🏗️ 项目架构：**

```go
chatmodel/
├── chatmodel.go          # 主程序入口：创建图书推荐代理、启用流式输出、实现检查点存储（内存存储）、支持对话恢复和继续
├── subagents/            # 代理实现
│   ├── agent.go          # 图书推荐代理：调用底层模型、配置了工具
│   ├── booksearch.go     # 图书搜索工具
│   └── ask_for_clarification.go  # 澄清问题工具
common/
├── model
│   ├── ark.go
│   └── openai.go
└── prints
    └── util.go
```



1. 创建 ChatModel: ark.go
```go
import (
	"context"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

func NewArkChatModel() model.ToolCallingChatModel {
	cm, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		APIKey:  os.Getenv("ARK_API_KEY"),
		Model:   os.Getenv("ARK_CHAT_MODEL"),
		BaseURL: os.Getenv("ARK_BASE_URL"),
	})
	if err != nil {
		log.Fatalf("openai.NewChatModel failed: %v", err)
	}
	return cm
}
```



1. utils.InferTool将本地函数转换一个tool: booksearch.go
```go
import (
    "context"
    "log"

    "github.com/cloudwego/eino/components/tool"
    "github.com/cloudwego/eino/components/tool/utils"
)

type BookSearchInput struct {
    Genre     string `json:"genre" jsonschema:"description=Preferred book genre,enum=fiction,enum=sci-fi,enum=mystery,enum=biography,enum=business"`
    MaxPages  int    `json:"max_pages" jsonschema:"description=Maximum page length (0 for no limit)"`
    MinRating int    `json:"min_rating" jsonschema:"description=Minimum user rating (0-5 scale)"`
}

type BookSearchOutput struct {
    Books []string
}

func NewBookRecommender() tool.InvokableTool {
    bookSearchTool, err := utils.InferTool("search_book", "Search books based on user preferences", func(ctx context.Context, input *BookSearchInput) (output *BookSearchOutput, err error) {
       // search code
       // ...
       return &BookSearchOutput{Books: []string{"God's blessing on this wonderful world!"}}, nil
    })
    if err != nil {
       log.Fatalf("failed to create search book tool: %v", err)
    }
    return bookSearchTool
}
```



1. 创建 ChatModelAgent: booksearch.go
```go
// eino-examples/adk/intro/chatmodel/subagents/agent.go
import (
    "context"
    "fmt"
    "log"

    "github.com/cloudwego/eino/adk"
    "github.com/cloudwego/eino/components/tool"
    "github.com/cloudwego/eino/compose"
)

func NewBookRecommendAgent() adk.Agent {
    ctx := context.Background()

    a, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
       Name:        "BookRecommender",
       Description: "An agent that can recommend books",
       Instruction: `You are an expert book recommender. Based on the user's request, use the "search_book" tool to find relevant books. Finally, present the results to the user.`,
       Model:       NewChatModel(),
       ToolsConfig: adk.ToolsConfig{
          ToolsNodeConfig: compose.ToolsNodeConfig{
             Tools: []tool.BaseTool{NewBookRecommender()},
          },
       },
    })
    if err != nil {
       log.Fatal(fmt.Errorf("failed to create chatmodel: %w", err))
    }

    return a
}
```



1. 通过 Runner 运行：chatmodel.go
```go
// eino-examples/adk/intro/chatmodel/chatmodel.go
import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/cloudwego/eino/adk"

    "github.com/cloudwego/eino-examples/adk/intro/chatmodel/subagents"
)

func main() {
    ctx := context.Background()
    a := subagents.NewBookRecommendAgent()
    runner := adk.NewRunner(ctx, adk.RunnerConfig{
       Agent: a,
    })
    iter := runner.Query(ctx, "recommend a fiction book to me")
    for {
       event, ok := iter.Next()
       if !ok {
          break
       }
       if event.Err != nil {
          log.Fatal(event.Err)
       }
       msg, err := event.Output.MessageOutput.GetMessage()
       if err != nil {
          log.Fatal(err)
       }
       fmt.Printf("\nmessage:\n%v\n======", msg)
    }
}
```



1. 工具 `ask_for_clarification` 使用了 Interrupt&Resume 能力来实现向用户“询问”。
```go
import (
    "context"
    "log"

    "github.com/cloudwego/eino/components/tool"
    "github.com/cloudwego/eino/components/tool/utils"
    "github.com/cloudwego/eino/compose"
)

type askForClarificationOptions struct {
    NewInput *string
}

func WithNewInput(input string) tool.Option {
    return tool.WrapImplSpecificOptFn(func(t *askForClarificationOptions) {
       t.NewInput = &input
    })
}

type AskForClarificationInput struct {
    Question string `json:"question" jsonschema:"description=The specific question you want to ask the user to get the missing information"`
}

func NewAskForClarificationTool() tool.InvokableTool {
    t, err := utils.InferOptionableTool(
       "ask_for_clarification",
       "Call this tool when the user's request is ambiguous or lacks the necessary information to proceed. Use it to ask a follow-up question to get the details you need, such as the book's genre, before you can use other tools effectively.",
       func(ctx context.Context, input *AskForClarificationInput, opts ...tool.Option) (output string, err error) {
          o := tool.GetImplSpecificOptions[askForClarificationOptions](nil, opts...)
          if o.NewInput == nil {
             return "", compose.NewInterruptAndRerunErr(input.Question)
          }
          return *o.NewInput, nil
       })
    if err != nil {
       log.Fatal(err)
    }
    return t
}
```



在 Runner 中配置 CheckPointStore（例子中使用最简单的 InMemoryStore），并在调用 Agent 时传入 CheckPointID (用来在恢复时使用)。

eino Graph 在中断时，会把 Graph 的 InterruptInfo 放入 Interrupted.Data 中：

```go
func main() {
    ctx := context.Background()
    a := internal.NewBookRecommendAgent()
    runner := adk.NewRunner(ctx, adk.RunnerConfig{
       Agent:           a,
       CheckPointStore: newInMemoryStore(),
    })
    iter := runner.Query(ctx, "recommend a book to me", adk.WithCheckPointID("1"))
    for {
       event, ok := iter.Next()
       if !ok {
          break
       }
       if event.Err != nil {
          log.Fatal(event.Err)
       }
       if event.Action != nil && event.Action.Interrupted != nil {
          fmt.Printf("\ninterrupt happened, info: %+v\n", event.Action.Interrupted.Data.(*compose.InterruptInfo).RerunNodesExtra["ToolNode"])
          continue
       }
       msg, err := event.Output.MessageOutput.GetMessage()
       if err != nil {
          log.Fatal(err)
       }
       fmt.Printf("\nmessage:\n%v\n======\n\n", msg)
    }
    
    // xxxxxx
}
```



之后向用户询问新输入并恢复运行

```go
func main(){
    // xxx
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Print("new input is:\n")
    scanner.Scan()
    nInput := scanner.Text()

    iter, err := runner.Resume(ctx, "1", adk.WithToolOptions([]tool.Option{chatmodel.WithNewInput(nInput)}))
    if err != nil {
        log.Fatal(err)
    }
    for {
        event, ok := iter.Next()
        if !ok {
           break
        }
        if event.Err != nil {
           log.Fatal(event.Err)
        }
        msg, err := event.Output.MessageOutput.GetMessage()
        if err != nil {
           log.Fatal(err)
        }
        fmt.Printf("\nmessage:\n%v\n======\n\n", msg)
    }
}

```





# 附：c**hain agent examples**

## **example：todoagent**

在构建 Agent 时，ToolsNode 是一个核心组件，它负责管理和执行工具调用。ToolsNode 可以集成多个工具，并提供统一的调用接口。它支持同步调用（Invoke）和流式调用（Stream）两种方式，能够灵活地处理不同类型的工具执行需求。

```go
import (
    "context"

    "github.com/cloudwego/eino/components/tool"
    "github.com/cloudwego/eino/compose"
)

conf := &compose.ToolsNodeConfig{
    Tools: []tool.BaseTool{tool1, tool2},  // 工具可以是 InvokableTool 或 StreamableTool
}
toolsNode, err := compose.NewToolNode(context.Background(), conf)
```

完整示例：

```go
import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/cloudwego/eino-ext/components/model/openai"
    "github.com/cloudwego/eino/components/tool"
    "github.com/cloudwego/eino/compose"
    "github.com/cloudwego/eino/schema"
)

func main() {
    // 初始化 tools
    todoTools := []tool.BaseTool{
        getAddTodoTool(),                               // NewTool 构建
        updateTool,                                     // InferTool 构建
        &ListTodoTool{},                                // 实现Tool接口
        searchTool,                                     // 官方封装的工具
    }

    // 创建并配置 ChatModel
    chatModel, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
        Model:       "gpt-4",
        APIKey:      os.Getenv("OPENAI_API_KEY"),
    })
    if err != nil {
        log.Fatal(err)
    }
    // 获取工具信息并绑定到 ChatModel
    toolInfos := make([]*schema.ToolInfo, 0, len(todoTools))
    for _, tool := range todoTools {
        info, err := tool.Info(ctx)
        if err != nil {
            log.Fatal(err)
        }
        toolInfos = append(toolInfos, info)
    }
    err = chatModel.BindTools(toolInfos)
    if err != nil {
        log.Fatal(err)
    }


    // 创建 tools 节点
    todoToolsNode, err := compose.NewToolNode(context.Background(), &compose.ToolsNodeConfig{
        Tools: todoTools,
    })
    if err != nil {
        log.Fatal(err)
    }

    // 构建完整的处理链
    chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
    chain.
        AppendChatModel(chatModel, compose.WithNodeName("chat_model")).
        AppendToolsNode(todoToolsNode, compose.WithNodeName("tools"))

    // 编译并运行 chain
    agent, err := chain.Compile(ctx)
    if err != nil {
        log.Fatal(err)
    }

    // 运行示例
    resp, err := agent.Invoke(ctx, []*schema.Message{
        {
           Role:    schema.User,
           Content: "添加一个学习 Eino 的 TODO，同时搜索一下 cloudwego/eino 的仓库地址",
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    // 输出结果
    for _, msg := range resp {
        fmt.Println(msg.Content)
    }
}
```

## example：程序员鼓励师chat

使用ChatModel构建一个简单的"程序员鼓励师" LLM 应用。包括：创建ChatTemplate、创建 ChatModel、运行ChatModel

> 代码库：[https://github.com/cloudwego/eino-examples/tree/main/quickstart/chat](https://github.com/cloudwego/eino-examples/tree/main/quickstart/chat)

1. **创建ChatTemplate (template.go)**
对话是通过 `schema.Message` 来表示，含以下重要字段：

- `Role`: 消息的角色，可以是：
    - `system`: 系统指令，用于设定模型的行为和角色
    - `user`: 用户的输入
    - `assistant`: 模型的回复 /ə'sɪstənt/ *n.* 助手
    - `tool`: 工具调用的结果
- `Content`: 消息的具体内容
**关键特性**：

- **参数化**：使用 {role}, {style}, {question} 等占位符
- **对话历史**：**通过 MessagesPlaceholder 支持多轮对话**
- **格式化**：使用 FString 格式进行参数替换
```go
// eino-examples/quickstart/chat/template.go

import (
    "context"

    "github.com/cloudwego/eino/components/prompt"
    "github.com/cloudwego/eino/schema"
)

// 创建模板，使用 FString 格式
template := prompt.FromMessages(schema.FString,
   // 系统消息模板
   schema.SystemMessage("你是一个{role}。你需要用{style}的语气回答问题。你的目标是帮助程序员保持积极乐观的心态，提供技术建议的同时也要关注他们的心理健康。"),

   // 插入需要的对话历史（新对话的话这里不填）
   schema.MessagesPlaceholder("chat_history", true),

   // 用户消息模板
   schema.UserMessage("问题: {question}"),
)

// 使用模板生成消息
messages, err := template.Format(context.Background(), map[string]any{
   "role":     "程序员鼓励师",
   "style":    "积极、温暖且专业",
   "question": "我的代码一直报错，感觉好沮丧，该怎么办？",
   // 对话历史（这个例子里模拟两轮对话历史）
   "chat_history": []*schema.Message{
      schema.UserMessage("你好"),
      schema.AssistantMessage("嘿！我是你的程序员鼓励师！记住，每个优秀的程序员都是从 Debug 中成长起来的。有什么我可以帮你的吗？", nil),
      schema.UserMessage("我觉得自己写的代码太烂了"),
      schema.AssistantMessage("每个程序员都经历过这个阶段！重要的是你在不断学习和进步。让我们一起看看代码，我相信通过重构和优化，它会变得更好。记住，Rome wasn't built in a day，代码质量是通过持续改进来提升的。", nil),
   },
})
```

1. **创建 ChatModel (模型抽象 ollama.go)**
Ollama 是一个开源的本地大语言模型运行框架，支持多种开源模型。

llama****/'lɑːmə/ n. 美洲驼；无峰驼

```go
// eino-examples/quickstart/chat/ollama.go

import (
    "github.com/cloudwego/eino-ext/components/model/ollama"
)


chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
    BaseURL: "http://localhost:11434", // Ollama 服务地址
    Model:   "llama2",                 // 模型名称
})
```

**统一接口**：model.ToolCallingChatModel

**设计优势:**

- **可插拔**：可以轻松切换不同的模型提供商
- **统一接口**：所有模型都实现相同的接口
- **配置化**：通过配置对象管理模型参数
```go
func createOllamaChatModel(ctx context.Context) model.ToolCallingChatModel {
    chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
        BaseURL: "http://localhost:11434",
        Model:   "llama2:7b",
    })
    return chatModel
}
```



1. **运行ChatModel**
Eino ChatModel 提供了两种运行模式：

- 输出完整消息(generate)
- 输出消息流(stream): 让 ChatModel 提供类似打字机的输出效果，使用户更早得到模型响应，提升用户体验。
**生成模式 vs 流式模式 (generate.go)**

```go
// 生成模式：一次性返回完整结果
func generate(ctx context.Context, llm model.ToolCallingChatModel, in []*schema.Message) *schema.Message {
    result, err := llm.Generate(ctx, in)
    return result
}

// 流式模式：实时返回每个 token
func stream(ctx context.Context, llm model.ToolCallingChatModel, in []*schema.Message) *schema.StreamReader[*schema.Message] {
    result, err := llm.Stream(ctx, in)
    return result
}
```



**流式处理 (stream.go)：逐 token 处理**

```go
// eino-examples/quickstart/chat/stream.go

import (
    "io"
    "log"

    "github.com/cloudwego/eino/schema"
)

func reportStream(sr *schema.StreamReader[*schema.Message]) {
    defer sr.Close()

    i := 0
    for {
       message, err := sr.Recv()
       if err == io.EOF { // 流式输出结束
          return
       }
       if err != nil {
          log.Fatalf("recv failed: %v", err)
       }
       // 处理每个 token
       log.Printf("message[%d]: %+v\n", i, message)
       i++
    }
}
```





E*ino Assistant*



[https://www.cloudwego.io/zh/docs/eino/overview/bytedance_eino_practice/](https://www.cloudwego.io/zh/docs/eino/overview/bytedance_eino_practice/)



[「火山引擎豆包模型」](https://console.volcengine.com/ark)：需要实名认证后购买使用，每人有 50万免费Tokens额度

![](/images/22524637-29b5-80e9-84e8-ff0346e83856/image_22624637-29b5-809d-9038-f6e6c7298821.jpg)



![](/images/22524637-29b5-80e9-84e8-ff0346e83856/image_22924637-29b5-80e6-a91a-fd24357da730.jpg)

![](/images/22524637-29b5-80e9-84e8-ff0346e83856/image_22924637-29b5-803c-bdfb-f6c0f749a704.jpg)







