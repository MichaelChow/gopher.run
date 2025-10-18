---
title: "5.9 Muti-Agent"
date: 2025-08-01T01:50:00Z
draft: false
weight: 5009
---

# 5.9 Muti-Agent

## **一、host j**ournal

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





## 二、eino assistant





## 三、deer-go

> [https://mp.weixin.qq.com/s/wT-UqAGxxJ0-h-zDqVXSSQ](https://mp.weixin.qq.com/s/wT-UqAGxxJ0-h-zDqVXSSQ)





## 四、manus



