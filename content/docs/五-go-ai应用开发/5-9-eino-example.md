---
title: "5.9 eino-example"
date: 2025-07-03T13:48:00Z
draft: false
weight: 5009
---

# 5.9 eino-example



# 一、仓库结构

> [https://github.com/cloudwego/eino-examples](https://github.com/cloudwego/eino-examples): 包含了 Eino 框架的示例和演示代码，提供了实用的示例来帮助开发者更好地理解和使用 Eino 的功能。

```markdown
- **adk/:** xxx
- **components/**: cloudwego/eino-ext 中各种组件的使用示例
    - 包含不同类型组件的实现和使用方式
    - 展示如何使用和自定义 Eino 的扩展组件
- **compose/**: Eino 编排能力的使用示例
    - 展示如何使用 Graph 和 Chain 进行编排
    - 提供不同组件组合的模式
    - 展示各种编排场景和最佳实践
- **devops/：**xxx
- **flow/**: Eino flow 模块的使用示例
    - 包含基于流的编程模式演示
    - 展示如何实现和管理数据流
    - 包含流处理的示例
- **quickstart/**: 用户文档中的快速入门示例
    - 帮助新用户快速上手的基础示例
    - 包含与官方文档中相同的演示代码
```

# 二、**quickstart/**

### quickstart/chat

由于默认的OpenAI模型需要翻墙、账号和绑定信用卡（实测招行的万事达信用卡被拒绝）等，较不方便。

<!-- 列布局开始 -->

![](/images/22524637-29b5-80e9-84e8-ff0346e83856/image_27e24637-29b5-800b-8a3c-ecf33419e9d6.jpg)


---

![](/images/22524637-29b5-80e9-84e8-ff0346e83856/image_27e24637-29b5-8031-bfeb-e47f3a070c31.jpg)



<!-- 列布局结束 -->

这里使用**本地的Ollama**（llama****/'lɑːmə/ n. 美洲驼；无峰驼），一个开源的本地大语言模型运行框架，支持多种开源模型。



**安装ollama:**

```shell
# 安装
brew install ollama
# 启动
ollama serve
# 确认正常运行
curl -s http://localhost:11434/api/tags
```



**跑通：**

```shell
cd eino-examples/quickstart/chat
go run .
# 默认流式输出
```

![](/images/22524637-29b5-80e9-84e8-ff0346e83856/image_27e24637-29b5-8005-9074-fba7685116d5.jpg)



**源码解析：**

- 简单的eino函数调用，熟悉最基本的调用流程
- 单步调试看执行流：eino-ext中ollama chatModel的实现
- **Pull Request：**llm → cm
### quickstart/todoagent



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





# 二、components/



# 三、compose/



# 四、flow/



# 五、devops/



# 六、adk/

详见：[四、adk example](https://www.notion.so/2472463729b580afa6a1e96b72202555#2472463729b580caa58efb2a949867de) 





# 七、agent app

## **host j**ournal

> 💡 Multi Agent 系统由多个协同工作的 Agent 组成，每个 Agent 都有其特定的职责和专长。通过 Agent 间的交互与协作，可以处理更复杂的任务，实现分工协作。这种方式特别适合需要多个专业领域知识结合的场景。

**实际案例：日记助手 多Agent**

**架构优势**

- **职责分离**：
    - **Host** - 负责意图识别，决定调用哪个专家
    - **Specialists** - 专家智能体，负责具体任务执行
    - **Summarizer** - 汇总多个专家的输出（可选）
- **模块化**：每个专家可以独立开发和部署
- **可扩展**：易于添加新的专家
- **专业化**：每个专家可以针对特定任务优化


**工作流程：**

用户输入 → Host意图识别 → 路由到专家 → 专家执行 → 汇总结果 → 返回用户



**1. 创建 Host**

- Host 使用强大的模型进行意图识别
- SystemPrompt 定义了 Host 的职责范围
- Host 会分析用户输入，决定调用哪个专家
```go
func newHost(ctx context.Context, baseURL, apiKey, modelName string) (*host.Host, error) {
    chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
        BaseURL: baseURL,
        Model:   modelName,
        ByAzure: true,
        APIKey:  apiKey,
    })
    if err != nil {
        return nil, err
    }

    return &host.Host{
        ChatModel:    chatModel,
        SystemPrompt: "You can read and write journal on behalf of the user. When user asks a question, always answer with journal content.",
    }, nil
}
```



1. **创建写日记专家**
- 专家使用专门的模型和参数
- 通过 Chain 编排处理流程
- 定义明确的 AgentMeta 信息
```go
func newWriteJournalSpecialist(ctx context.Context) (*host.Specialist, error) {
    chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
        BaseURL: "http://localhost:11434",
        Model:   "llama3-groq-tool-use",
        Options: &api.Options{
            Temperature: 0.000001, // 低温度确保输出稳定
        },
    })
    if err != nil {
        return nil, err
    }

    // 创建处理链：重写用户查询 → 写入文件
    chain := compose.NewChain[[]*schema.Message, *schema.Message]()
    
    // 第一步：重写用户查询，提取日记内容
    chain.AppendLambda(compose.InvokableLambda(func(ctx context.Context, input []*schema.Message) ([]*schema.Message, error) {
        systemMsg := &schema.Message{
            Role:    schema.System,
            Content: "You are responsible for preparing the user query for insertion into journal. The user's query is expected to contain the actual text the user want to write to journal, as well as convey the intention that this query should be written to journal. You job is to remove that intention from the user query, while preserving as much as possible the user's original query, and output ONLY the text to be written into journal",
        }
        return append([]*schema.Message{systemMsg}, input...), nil
    })).
        AppendChatModel(chatModel).
        AppendLambda(compose.InvokableLambda(func(ctx context.Context, input *schema.Message) (string, error) {
            // 写入文件
            now := time.Now()
            dateStr := now.Format("2006-01-02")
            filename := fmt.Sprintf("journal_%s.txt", dateStr)
            
            content := fmt.Sprintf("%s\n", input.Content)
            err := os.WriteFile(filename, []byte(content), 0644)
            if err != nil {
                return "", err
            }
            
            return fmt.Sprintf("Journal written successfully: %s", input.Content), nil
        }))

    r, err := chain.Compile(ctx)
    if err != nil {
        return nil, err
    }

    return &host.Specialist{
        AgentMeta: host.AgentMeta{
            Name:        "write_journal",
            IntendedUse: "write user's content to journal file",
        },
        Invokable: func(ctx context.Context, input []*schema.Message, opts ...agent.AgentOption) (*schema.Message, error) {
            return r.Invoke(ctx, input, agent.GetComposeOptions(opts...)...)
        },
    }, nil
}
```

**3. 创建读日记专家**

```go
func newReadJournalSpecialist(ctx context.Context) (*host.Specialist, error) {
    chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
        BaseURL: "http://localhost:11434",
        Model:   "llama3-groq-tool-use",
        Options: &api.Options{
            Temperature: 0.000001,
        },
    })
    if err != nil {
        return nil, err
    }

    // 创建处理链：读取文件 → 格式化输出
    chain := compose.NewChain[[]*schema.Message, *schema.Message]()
    chain.AppendLambda(compose.InvokableLambda(func(ctx context.Context, input []*schema.Message) (string, error) {
        now := time.Now()
        dateStr := now.Format("2006-01-02")
        filename := fmt.Sprintf("journal_%s.txt", dateStr)
        
        content, err := os.ReadFile(filename)
        if err != nil {
            if os.IsNotExist(err) {
                return "No journal entries found for today.", nil
            }
            return "", err
        }
        
        return string(content), nil
    }))

    r, err := chain.Compile(ctx)
    if err != nil {
        return nil, err
    }

    return &host.Specialist{
        AgentMeta: host.AgentMeta{
            Name:        "view_journal_content",
            IntendedUse: "read and display journal content",
        },
        Invokable: func(ctx context.Context, input []*schema.Message, opts ...agent.AgentOption) (*schema.Message, error) {
            return r.Invoke(ctx, input, agent.GetComposeOptions(opts...)...)
        },
    }, nil
}
```



**4. 创建问答专家**

```go
func newAnswerWithJournalSpecialist(ctx context.Context) (*host.Specialist, error) {
    chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
        BaseURL: "http://localhost:11434",
        Model:   "llama3-groq-tool-use",
        Options: &api.Options{
            Temperature: 0.000001,
        },
    })
    if err != nil {
        return nil, err
    }

    // 创建图：加载日记 → 提取查询 → 模板 → 模型 → 回答
    graph := compose.NewGraph[[]*schema.Message, *schema.Message]()

    // 加载日记节点
    if err = graph.AddLambdaNode("journal_loader", compose.InvokableLambda(func(ctx context.Context, input []*schema.Message) (string, error) {
        now := time.Now()
        dateStr := now.Format("2006-01-02")
        return loadJournal(dateStr)
    }), compose.WithOutputKey("journal")); err != nil {
        return nil, err
    }

    // 提取查询节点
    if err = graph.AddLambdaNode("query_extractor", compose.InvokableLambda(func(ctx context.Context, input []*schema.Message) (string, error) {
        return input[len(input)-1].Content, nil
    }), compose.WithOutputKey("query")); err != nil {
        return nil, err
    }

    // 创建模板
    systemTpl := `Answer user's query based on journal content: {journal}`
    chatTpl := prompt.FromMessages(schema.FString,
        schema.SystemMessage(systemTpl),
        schema.UserMessage("{query}"),
    )
    if err = graph.AddChatTemplateNode("template", chatTpl); err != nil {
        return nil, err
    }

    if err = graph.AddChatModelNode("model", chatModel); err != nil {
        return nil, err
    }

    // 连接节点
    if err = graph.AddEdge("journal_loader", "template"); err != nil {
        return nil, err
    }
    if err = graph.AddEdge("query_extractor", "template"); err != nil {
        return nil, err
    }
    if err = graph.AddEdge("template", "model"); err != nil {
        return nil, err
    }
    if err = graph.AddEdge(compose.START, "journal_loader"); err != nil {
        return nil, err
    }
    if err = graph.AddEdge(compose.START, "query_extractor"); err != nil {
        return nil, err
    }
    if err = graph.AddEdge("model", compose.END); err != nil {
        return nil, err
    }

    r, err := graph.Compile(ctx)
    if err != nil {
        return nil, err
    }

    return &host.Specialist{
        AgentMeta: host.AgentMeta{
            Name:        "answer_with_journal",
            IntendedUse: "load journal content and answer user's question with it",
        },
        Invokable: func(ctx context.Context, input []*schema.Message, opts ...agent.AgentOption) (*schema.Message, error) {
            return r.Invoke(ctx, input, agent.GetComposeOptions(opts...)...)
        },
    }, nil
}
```

**5. 组装 Multi-Agent**

```go
func main() {
    ctx := context.Background()
    
    // 创建 Host
    h, err := newHost(ctx, "your_base_url", "your_api_key", "gpt-4")
    if err != nil {
        panic(err)
    }

    // 创建专家们
    writer, err := newWriteJournalSpecialist(ctx)
    if err != nil {
        panic(err)
    }

    reader, err := newReadJournalSpecialist(ctx)
    if err != nil {
        panic(err)
    }

    answerer, err := newAnswerWithJournalSpecialist(ctx)
    if err != nil {
        panic(err)
    }

    // 组装 Multi-Agent
    hostMA, err := host.NewMultiAgent(ctx, &host.MultiAgentConfig{
        Host: *h,
        Specialists: []*host.Specialist{
            writer,
            reader,
            answerer,
        },
    })
    if err != nil {
        panic(err)
    }

    // 创建回调处理器
    cb := &logCallback{}

    // 交互循环
    for {
        println("\n\nYou: ")

        var message string
        scanner := bufio.NewScanner(os.Stdin)
        for scanner.Scan() {
            message += scanner.Text()
            break
        }

        if err := scanner.Err(); err != nil {
            panic(err)
        }

        if message == "exit" {
            return
        }

        msg := &schema.Message{
            Role:    schema.User,
            Content: message,
        }

        // 流式调用
        out, err := hostMA.Stream(ctx, []*schema.Message{msg}, host.WithAgentCallbacks(cb))
        if err != nil {
            panic(err)
        }

        defer out.Close()

        println("\nAnswer:")

        for {
            msg, err := out.Recv()
            if err != nil {
                if err == io.EOF {
                    break
                }
            }

            print(msg.Content)
        }
    }
}
```



**高级配置**

**1. 自定义 StreamToolCallChecker**

```go
func customStreamToolCallChecker(ctx context.Context, sr *schema.StreamReader[*schema.Message]) (bool, error) {
    defer sr.Close()
    for {
        msg, err := sr.Recv()
        if err != nil {
            if errors.Is(err, io.EOF) {
                break
            }
            return false, err
        }

        if len(msg.ToolCalls) > 0 {
            return true, nil
        }
    }
    return false, nil
}

// 在创建 Host 时使用
host := &host.Host{
    ChatModel:              chatModel,
    SystemPrompt:           "Your system prompt",
    StreamToolCallChecker:  customStreamToolCallChecker,
}
```

1. **配置 Summarizer**
当 Host 同时选择多个专家时，需要 Summarizer 来汇总结果：

```go
hostMA, err := host.NewMultiAgent(ctx, &host.MultiAgentConfig{
    Host: *h,
    Specialists: []*host.Specialist{
        writer,
        reader,
        answerer,
    },
    Summarizer: &host.Summarizer{
        ChatModel:    summarizerModel,
        SystemPrompt: "Summarize the outputs from multiple specialists into a coherent response.",
    },
})
```



## eino assistant



## deer-go

> [https://mp.weixin.qq.com/s/wT-UqAGxxJ0-h-zDqVXSSQ](https://mp.weixin.qq.com/s/wT-UqAGxxJ0-h-zDqVXSSQ)



### manus



## **todoagent**

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

## 程序员鼓励师chat

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







