---
title: "5.8 eino-ext"
date: 2025-09-30T04:49:00Z
draft: false
weight: 5008
---

# 5.8 eino-ext

# **一、Eino Dev**

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
![](/images/27e24637-29b5-8078-a249-ef1975df4f55/image_29324637-29b5-80b6-9bf8-dbacb2817b45.jpg)





更多文档：

- [Eino Dev 可视化编排插件功能指南](https://www.cloudwego.io/zh/docs/eino/core_modules/devops/visual_orchestration_plugin_guide/#%E7%BC%96%E6%8E%92%E7%BB%84%E4%BB%B6%E4%BB%8B%E7%BB%8D)
- [可视化开发](https://www.cloudwego.io/zh/docs/eino/overview/bytedance_eino_practice/#%E5%8F%AF%E8%A7%86%E5%8C%96%E5%BC%80%E5%8F%91-1)
## 功能2：可视化Debug

- 运行 [源码地址](https://github.com/cloudwego/eino-examples/blob/3a94b9ab0db133907636c07ef1e3cf267551725c/devops/debug/main.go) 
![](/images/27e24637-29b5-8078-a249-ef1975df4f55/image_29424637-29b5-80ea-92b5-c3f3f12db3c7.jpg)

- Eino Dev 配置调试地址，选择需要调试的Graph
- 点击 Test Run。默认从star节点开始执行，可以点击可视化graph从任意节点开始执行，看到每个节点的input、output。（类似 AI Agent版本的trace）
![](/images/27e24637-29b5-8078-a249-ef1975df4f55/image_29424637-29b5-80ea-820f-f63269084836.jpg)



![](/images/27e24637-29b5-8078-a249-ef1975df4f55/image_29424637-29b5-80c4-9238-c4e8884c9413.jpg)

![](/images/27e24637-29b5-8078-a249-ef1975df4f55/image_29424637-29b5-8002-941f-c68efa4289b9.jpg)



- 高级功能：**指定 interface 字段的实现类型**
    ![](/images/27e24637-29b5-8078-a249-ef1975df4f55/image_29524637-29b5-80e1-b502-e3bf13cc7fba.jpg)


更多文档：[Eino Dev 可视化调试插件功能指南](https://www.cloudwego.io/zh/docs/eino/core_modules/devops/visual_debug_plugin_guide/)



# 二、a2a

A2A (Agent-to-Agent) 是一个标准化的 Agent 间通信协议实现，允许不同的 AI Agent 通过统一的接口进行交互和协作。

该模块使eino-ext拓展组件实现。

**核心特性：**

- 🔄 支持同步和异步消息交互
- 📡 支持流式响应（Streaming）
- 🔔 支持 Push Notification（异步通知）
- 🔐 支持多种认证方式
- 🎯 任务状态管理和生命周期控制
- 🔌 可插拔的传输层（目前支持 JSON-RPC）
- 🧩 与 Eino ADK (Agent Development Kit) 无缝集成
**使用场景**：

1. **多 Agent 协作系统**：不同 Agent 之间需要标准化通信
1. **Agent 服务化**：将 Agent 能力封装为可远程调用的服务
1. **Agent 编排**：构建复杂的 Agent 工作流
1. **跨组织 Agent 调用**：通过标准协议实现不同组织开发的 Agent 互通


模块层级结构:

```shell
a2a/
├── models/          ** # 数据模型定义**
│   ├── task.go      # Task 相关数据结构
│   ├── message.go   # Message 相关数据结构
│   ├── artifact.go  # Artifact 数据结构
│   ├── card.go      # Agent Card 定义
│   ├── part.go      # Message Part 定义
│   ├── handler.go   # Handler 接口定义
│   └── ...
├── client/          **# A2A 客户端实现**
│   └── client.go
├── server/          **# A2A 服务端实现**
│   ├── server.go
│   ├── eventqueue.go    # 事件队列
│   ├── taskstore.go     # 任务存储
│   ├── tasklocker.go    # 任务锁
│   └── notifier.go      # 推送通知
├── transport/       **# 传输层抽象和实现**
│   ├── transport.go      # 传输层接口
│   └── jsonrpc/         # JSON-RPC 实现
│       ├── client/
│       ├── server/
│       └── core/
├── extension/      ** # 扩展集成**
│   └── eino/       # Eino ADK 集成
│       ├── server.go    # Eino Server 适配器
│       ├── client.go    # Eino Client 适配器
│       └── utils.go
├── utils/          **# 工具函数**
└── examples/       **# 示例代码**
    ├── client/
    └── server/
```

分层架构图:

```shell
┌─────────────────────────────────────────────────────────┐
│                     Application Layer                    │
│  ┌──────────────────┐        ┌──────────────────┐      │
│  │   Eino Agent     │        │   Custom App     │      │
│  └────────┬─────────┘        └────────┬─────────┘      │
└───────────┼──────────────────────────┼─────────────────┘
            │                           │
┌───────────┼──────────────────────────┼─────────────────┐
│           │      A2A Core Layer      │                  │
│  ┌────────▼─────────┐       ┌────────▼─────────┐      │
│  │   A2A Server     │       │   A2A Client     │      │
│  │  (server.go)     │       │  (client.go)     │      │
│  └────────┬─────────┘       └────────┬─────────┘      │
│           │                           │                  │
│  ┌────────▼───────────────────────────▼─────────┐      │
│  │          Models & Data Structures             │      │
│  │  Task, Message, Artifact, AgentCard, etc.    │      │
│  └────────┬──────────────────────────────────────┘      │
└───────────┼─────────────────────────────────────────────┘
            │
┌───────────┼─────────────────────────────────────────────┐
│           │      Transport Layer                         │
│  ┌────────▼─────────────────────────────────────┐      │
│  │        JSON-RPC over HTTP/HTTPS               │      │
│  │  (transport/jsonrpc/)                         │      │
│  └───────────────────────────────────────────────┘      │
└─────────────────────────────────────────────────────────┘
```





### Task（任务）

是 A2A 协议中的核心概念，代表一次完整的 Agent 交互过程。

**Task 的生命周期状态：**

```go
type TaskState string

const (
    TaskStateSubmitted     TaskState = "submitted"      // 已提交，等待处理
    TaskStateWorking       TaskState = "working"        // 正在处理
    TaskStateInputRequired TaskState = "input-required" // 需要用户输入（暂停）
    TaskStateCompleted     TaskState = "completed"      // 已完成（终态）
    TaskStateCanceled      TaskState = "canceled"       // 已取消（终态）
    TaskStateFailed        TaskState = "failed"         // 失败（终态）
    TaskStateRejected      TaskState = "rejected"       // 被拒绝（终态）
    TaskStateAuthRequired  TaskState = "auth-required"  // 需要认证（暂停）
    TaskStateUnknown       TaskState = "unknown"        // 未知状态
)

```

**Task 状态转换图：**

```plain text
                    ┌──────────────┐
                    │  submitted   │
                    └──────┬───────┘
                           │
                    ┌──────▼───────┐
            ┌───────┤   working    ├───────┐
            │       └──────┬───────┘       │
            │              │               │
    ┌───────▼─────┐ ┌──────▼──────┐ ┌─────▼────────┐
    │input-required│ │   completed  │ │auth-required │
    │  (paused)    │ │  (terminal)  │ │  (paused)    │
    └──────────────┘ └──────────────┘ └──────────────┘
            │
    ┌───────▼─────┐ ┌───────────┐  ┌──────────┐
    │  canceled   │ │  failed   │  │ rejected │
    │ (terminal)  │ │(terminal) │  │(terminal)│
    └─────────────┘ └───────────┘  └──────────┘

```

**Task 数据结构：**

```go
type Task struct {
    ID        string      // 唯一任务 ID（UUID）
    ContextID string      // 上下文 ID，用于关联多个任务
    Status    TaskStatus  // 当前状态
    Artifacts []*Artifact // 生成的工件（输出）
    History   []*Message  // 历史消息记录
    Metadata  map[string]any // 元数据
}

type TaskStatus struct {
    State     TaskState // 状态
    Message   *Message  // 关联消息
    Timestamp string    // 时间戳（ISO 8601）
}

```

### Message（消息）

Message 表示用户或 Agent 之间交换的信息。

```go
type Message struct {
    Role             Role            // "user" 或 "agent"
    Parts            []Part          // 消息内容（可多模态）
    Metadata         map[string]any  // 元数据
    ReferenceTaskIDs []string        // 引用的任务 ID
    MessageID        string          // 消息 ID
    TaskID           *string         // 所属任务 ID
    ContextID        *string         // 上下文 ID
}

```

**Part（消息片段）支持的类型：**

```go
type PartKind string

const (
    PartKindText PartKind = "text"  // 文本
    PartKindFile PartKind = "file"  // 文件
    PartKindData PartKind = "data"  // 结构化数据
)

type Part struct {
    Kind     PartKind
    Text     *string         // 文本内容
    File     *FileContent    // 文件内容（Base64 或 URI）
    Data     map[string]any  // 结构化数据
    Metadata map[string]any
}

```

### Artifact（工件）

Artifact 表示 Agent 生成的输出或中间结果。

```go
type Artifact struct {
    ArtifactID  string          // 唯一标识
    Name        string          // 名称
    Description string          // 描述
    Parts       []Part          // 内容（可多模态）
    Metadata    map[string]any  // 元数据
}

```

**使用场景：**

- 代码生成结果
- 图像/文档生成
- 分析报告
- 中间处理结果
### Agent Card（Agent 名片）

Agent Card 描述了一个 Agent 的基本信息和能力。

```go
type AgentCard struct {
    ProtocolVersion    string              // A2A 协议版本（"0.2.5"）
    Name               string              // Agent 名称
    Description        string              // 描述
    URL                string              // 服务地址
    Version            string              // Agent 版本
    Capabilities       AgentCapabilities   // 能力声明
    Skills             []AgentSkill        // 技能列表
    SecuritySchemes    map[string]*SecurityScheme
    DefaultInputModes  []string            // 支持的输入模式
    DefaultOutputModes []string            // 支持的输出模式
}

type AgentCapabilities struct {
    Streaming              bool  // 是否支持流式
    PushNotifications      bool  // 是否支持推送通知
    StateTransitionHistory bool  // 是否记录状态转换历史
}
```

### Server 架构端实现

Server 架构

```go
type A2AServer struct {
    agentCard               *models.AgentCard
    messageHandler          MessageHandler          // 普通消息处理器
    messageStreamingHandler MessageStreamingHandler // 流式消息处理器
    cancelTaskHandler       CancelTaskHandler       // 取消任务处理器
    taskEventsConsolidator  TaskEventsConsolidator  // 事件合并器
    logger                  Logger
    taskIDGenerator         func(ctx context.Context) (string, error)
    contextIDGenerator      func(ctx context.Context) (string, error)
    taskStore               TaskStore   // 任务存储
    taskLocker              TaskLocker  // 任务锁
    queue                   EventQueue  // 事件队列
    pushNotifier            PushNotifier // 推送通知器
}

```

核心 Handler 接口

```go
// 普通消息处理器（同步）
type MessageHandler func(
    ctx context.Context,
    params *InputParams,
) (*models.TaskContent, error)

// 流式消息处理器（异步）
type MessageStreamingHandler func(
    ctx context.Context,
    params *InputParams,
    writer ResponseEventWriter,
) error

// 取消任务处理器
type CancelTaskHandler func(
    ctx context.Context,
    params *InputParams,
) (*models.TaskContent, error)

// 事件合并器：将流式事件合并为最终任务状态
type TaskEventsConsolidator func(
    ctx context.Context,
    t *models.Task,
    events []models.ResponseEvent,
    handleErr error,
) *models.TaskContent

```

消息处理流程：

同步消息处理

```plain text
Client Request
     │
     ▼
┌─────────────────┐
│ SendMessage     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Lock Task       │──────┐ (new task or existing)
└────────┬────────┘      │
         │               │
         ▼               │
┌─────────────────┐      │
│ MessageHandler  │      │
└────────┬────────┘      │
         │               │
         ▼               │
┌─────────────────┐      │
│ Update TaskStore│      │
└────────┬────────┘      │
         │               │
         ▼               │
┌─────────────────┐      │
│ Unlock Task     │◄─────┘
└────────┬────────┘
         │
         ▼
   Return Result

```

流式消息处理：

```plain text
Client Request
     │
     ▼
┌──────────────────┐
│SendMessageStream │
└─────────┬────────┘
          │
          ▼
┌──────────────────┐
│ Lock Task        │
│ Reset Queue      │
└─────────┬────────┘
          │
          ├──────────────────────────┐
          │                          │
          ▼                          ▼
┌──────────────────┐      ┌──────────────────┐
│ Async Execution  │      │ Stream Reader    │
│                  │      │  (Pop Queue)     │
│ ┌──────────────┐ │      │                  │
│ │   Handler    │ │      │  ┌────────────┐ │
│ │   Execute    │─┼──┬───┼─►│ Send Event │─┼──► Client
│ └──────────────┘ │  │   │  └────────────┘ │
│                  │  │   │                  │
│ ┌──────────────┐ │  │   └──────────────────┘
│ │ Push to Queue│◄┼──┘
│ └──────────────┘ │
│                  │
│ ┌──────────────┐ │
│ │ Consolidate  │ │
│ └──────────────┘ │
│                  │
│ ┌──────────────┐ │
│ │Update & Save │ │
│ └──────────────┘ │
│                  │
│ ┌──────────────┐ │
│ │Unlock & Close│ │
│ └──────────────┘ │
└──────────────────┘

```

**关键点：**

1. **异步执行**：Handler 在独立的 goroutine 中执行
1. **事件队列**：通过队列实现生产者-消费者模式
1. **流式传输**：客户端通过 SSE (Server-Sent Events) 实时接收事件
1. **任务锁**：保证任务处理的并发安全
1. **错误恢复**：支持 panic 捕获和错误传播


TaskStore（任务存储）

```go
type TaskStore interface {
    Get(ctx context.Context, taskID string) (*models.Task, bool, error)
    Save(ctx context.Context, task *models.Task) error
}

```

**实现方式：**

- 默认：内存存储（`inMemoryTaskStore`）
- 可扩展：Redis、数据库等持久化存储


TaskLocker（任务锁）

```go
type TaskLocker interface {
    Lock(ctx context.Context, taskID string) error
    Unlock(ctx context.Context, taskID string) error
}

```

**作用：**

- 防止同一任务的并发修改
- 保证任务状态的一致性


EventQueue（事件队列）

```go
type EventQueue interface {
    Reset(ctx context.Context, taskID string) error
    Push(ctx context.Context, taskID string,
         event *models.SendMessageStreamingResponseUnion,
         err error) error
    Pop(ctx context.Context, taskID string) (
        event *models.SendMessageStreamingResponseUnion,
        err error,
        closed bool,
        popErr error)
    Close(ctx context.Context, taskID string) error
}

```

**实现：**

- 基于 channel 的内存队列
- 支持多个任务的并发队列管理
- 每个任务有独立的事件队列


PushNotifier（推送通知）

```go
type PushNotifier interface {
    Set(ctx context.Context, config *models.TaskPushNotificationConfig) error
    Get(ctx context.Context, configID string) (
        models.PushNotificationConfig, bool, error)
    SendNotification(ctx context.Context,
        event *models.SendMessageStreamingResponseUnion) error
}

```

**使用场景：**

- 长时间运行的任务
- 异步通知客户端
- Webhook 集成


服务端使用示例：

```go
import (
    "github.com/cloudwego/eino-ext/a2a/server"
    "github.com/cloudwego/eino-ext/a2a/transport/jsonrpc"
)

func main() {
    ctx := context.Background()

    // 1. 创建 Hertz HTTP 服务器
    hz := hertz_server.Default()

    // 2. 创建 JSON-RPC 注册器
    registrar, _ := jsonrpc.NewRegistrar(ctx, &jsonrpc.ServerConfig{
        Router:      hz,
        HandlerPath: "/a2a",
    })

    // 3. 注册 A2A 处理器
    server.RegisterHandlers(ctx, registrar, &server.Config{
        AgentCardConfig: server.AgentCardConfig{
            Name:        "My Agent",
            Description: "A helpful AI agent",
            URL:         "<https://example.com/a2a>",
            Version:     "1.0.0",
        },

        // 流式消息处理器
        MessageStreamingHandler: func(ctx context.Context,
                                     params *server.InputParams,
                                     writer server.ResponseEventWriter) error {
            // 处理用户输入
            userInput := params.Input

            // 发送状态更新
            writer.Write(models.ResponseEvent{
                TaskStatusUpdateEventContent: &models.TaskStatusUpdateEventContent{
                    Status: models.TaskStatus{
                        State: models.TaskStateWorking,
                    },
                },
            })

            // 生成输出
            result := processInput(userInput)

            // 发送结果
            writer.Write(models.ResponseEvent{
                TaskArtifactUpdateEventContent: &models.TaskArtifactUpdateEventContent{
                    Artifact: models.Artifact{
                        Parts: []models.Part{
                            {Kind: models.PartKindText, Text: &result},
                        },
                    },
                    LastChunk: true,
                },
            })

            return nil
        },

        // 任务取消处理器
        CancelTaskHandler: func(ctx context.Context,
                               params *server.InputParams) (*models.TaskContent, error) {
            return &models.TaskContent{
                Status: models.TaskStatus{State: models.TaskStateCanceled},
            }, nil
        },

        // 事件合并器
        TaskEventsConsolidator: consolidateEvents,
    })

    hz.Run()
}

```

### Client端实现详解

Client 架构

```go
type A2AClient struct {
    cli transport.ClientTransport
}

```

主要方法：

```go
// 获取 Agent 信息
func (c *A2AClient) AgentCard(ctx context.Context) (*models.AgentCard, error)

// 发送消息（同步）
func (c *A2AClient) SendMessage(ctx context.Context,
    params *models.MessageSendParams) (*models.SendMessageResponseUnion, error)

// 发送消息（流式）
func (c *A2AClient) SendMessageStreaming(ctx context.Context,
    params *models.MessageSendParams) (*ServerStreamingWrapper, error)

// 获取任务状态
func (c *A2AClient) GetTask(ctx context.Context,
    params *models.TaskQueryParams) (*models.Task, error)

// 取消任务
func (c *A2AClient) CancelTask(ctx context.Context,
    params *models.TaskIDParams) (*models.Task, error)

// 重新订阅任务（断线重连）
func (c *A2AClient) ResubscribeTask(ctx context.Context,
    params *models.TaskIDParams) (*ServerStreamingWrapper, error)

```

客户端使用示例：

```go
import (
    "github.com/cloudwego/eino-ext/a2a/client"
    "github.com/cloudwego/eino-ext/a2a/transport/jsonrpc"
)

func main() {
    ctx := context.Background()

    // 1. 创建传输层
    transport, _ := jsonrpc.NewTransport(ctx, &jsonrpc.ClientConfig{
        BaseURL:     "<http://localhost:8080>",
        HandlerPath: "/a2a",
    })

    // 2. 创建客户端
    cli, _ := client.NewA2AClient(ctx, &client.Config{
        Transport: transport,
    })

    // 3. 获取 Agent 信息
    card, _ := cli.AgentCard(ctx)
    fmt.Printf("Agent: %s\\n", card.Name)

    // 4. 发送流式消息
    stream, _ := cli.SendMessageStreaming(ctx, &models.MessageSendParams{
        Message: models.Message{
            Role: models.RoleUser,
            Parts: []models.Part{
                {Kind: models.PartKindText, Text: ptr("Hello, agent!")},
            },
        },
    })

    // 5. 接收流式响应
    for {
        event, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatal(err)
        }

        // 处理事件
        if event.TaskStatusUpdateEvent != nil {
            fmt.Printf("Status: %s\\n", event.TaskStatusUpdateEvent.Status.State)
        }
        if event.TaskArtifactUpdateEvent != nil {
            fmt.Printf("Artifact: %v\\n", event.TaskArtifactUpdateEvent.Artifact)
        }
    }
}

```

### Eino 集成

Eino Server 集成：**将 Eino ADK Agent 包装为 A2A 服务**。

```go
import (
    "github.com/cloudwego/eino/adk"
    einoa2a "github.com/cloudwego/eino-ext/a2a/extension/eino"
)

func main() {
    ctx := context.Background()

    // 1. 创建 Eino Agent
    agent := createMyEinoAgent()

    // 2. 创建 JSON-RPC 注册器
    registrar, _ := jsonrpc.NewRegistrar(ctx, &jsonrpc.ServerConfig{
        Router:      hertz_server.Default(),
        HandlerPath: "/agent",
    })

    // 3. 注册为 A2A 服务
    einoa2a.RegisterServerHandlers(ctx, agent, &einoa2a.ServerConfig{
        Registrar: registrar,

        // Agent 运行选项转换器
        AgentRunOptionConvertor: func(ctx context.Context,
            t *models.Task,
            input *models.Message,
            metadata map[string]any) ([]adk.AgentRunOption, error) {
            // 从 A2A Message 转换为 ADK 运行选项
            return []adk.AgentRunOption{}, nil
        },

        // Checkpoint 存储（支持中断恢复）
        CheckPointStore: myCheckpointStore,

        // 历史消息转换器
        HistoryMessageConvertor: func(ctx context.Context,
            messages []*models.Message) ([]adk.Message, error) {
            // 从 A2A Messages 转换为 ADK Messages
            return convertMessages(messages), nil
        },

        // 恢复选项转换器（用于中断后恢复）
        ResumeConvertor: func(ctx context.Context,
            t *models.Task,
            input *models.Message,
            metadata map[string]any) ([]adk.AgentRunOption, error) {
            return []adk.AgentRunOption{}, nil
        },
    })
}

```

**事件转换流程：**

```plain text
ADK AgentEvent → A2A ResponseEvent

┌──────────────────────┐
│  AgentEvent          │
├──────────────────────┤
│ - Action             │─┐
│   - Interrupted      │ │    ┌──────────────────────┐
│   - TransferToAgent  │─┼───►│ TaskStatusUpdate     │
│ - Output             │ │    │ - State              │
│   - MessageOutput    │─┘    │ - Message            │
└──────────────────────┘      └──────────────────────┘
          │
          └──────────────────► ┌──────────────────────┐
                               │ TaskArtifactUpdate   │
                               │ - Artifact           │
                               │ - LastChunk          │
                               └──────────────────────┘

```

Eino Client 集成：**将远程 A2A 服务包装为 Eino Agent。**

```go
import (
    einoa2a "github.com/cloudwego/eino-ext/a2a/extension/eino"
)

func main() {
    ctx := context.Background()

    // 1. 创建 A2A 传输层
    transport, _ := jsonrpc.NewTransport(ctx, &jsonrpc.ClientConfig{
        BaseURL:     "<http://remote-agent:8080>",
        HandlerPath: "/agent",
    })

    // 2. 创建 Eino Agent（包装 A2A Client）
    agent, _ := einoa2a.NewAgent(ctx, einoa2a.AgentConfig{
        Transport: transport,

        // 可选：自定义输入转换
        InputMessageConvertor: func(ctx context.Context,
            messages []*schema.Message) (models.Message, error) {
            return convertToA2AMessage(messages), nil
        },

        // 可选：自定义输出转换
        OutputConvertor: func(ctx context.Context,
            receiver *einoa2a.ResponseUnionReceiver,
            sender *einoa2a.AgentEventSender) {
            // 自定义从 A2A 响应到 ADK 事件的转换逻辑
        },
    })

    // 3. 像使用普通 Eino Agent 一样使用
    runner := adk.NewRunner(ctx, adk.RunnerConfig{
        Agent: agent,
    })

    iter := runner.Run(ctx, []adk.Message{
        schema.UserMessage("Hello!"),
    })

    // 处理结果
    for {
        event, ok := iter.Next()
        if !ok {
            break
        }
        handleEvent(event)
    }
}

```

**中断与恢复支持：**

```go
// Agent 执行过程中发生中断
iter := runner.Run(ctx, input)
for {
    event, ok := iter.Next()
    if !ok {
        break
    }

    // 检测到中断
    if event.Action != nil && event.Action.Interrupted != nil {
        interruptInfo := event.Action.Interrupted

        // 保存中断信息（自动保存在 CheckPointStore）
        fmt.Printf("Agent interrupted: %v\\n", interruptInfo.Data)

        // ... 等待用户输入 ...

        // 恢复执行
        resumeIter, _ := runner.Resume(ctx, interruptInfo.CheckPointID,
            einoa2a.WithResumeMessages(userResponse))
        // 继续处理
    }
}

```

### 传输层实现

Transport 接口

```go
// 客户端传输接口
type ClientTransport interface {
    AgentCard(ctx context.Context) (*models.AgentCard, error)
    SendMessage(ctx context.Context, params *models.MessageSendParams)
        (*models.SendMessageResponseUnion, error)
    SendMessageStreaming(ctx context.Context, params *models.MessageSendParams)
        (models.ResponseReader, error)
    GetTask(ctx context.Context, params *models.TaskQueryParams)
        (*models.Task, error)
    CancelTask(ctx context.Context, params *models.TaskIDParams)
        (*models.Task, error)
    ResubscribeTask(ctx context.Context, params *models.TaskIDParams)
        (models.ResponseReader, error)
    Close() error
}

// 服务端注册接口
type HandlerRegistrar interface {
    Register(context.Context, *models.ServerHandlers) error
}

```

JSON-RPC 实现：

目前支持的传输协议是 JSON-RPC over HTTP/HTTPS。

**特点：**

- 基于 CloudWeGo Hertz HTTP 框架
- 支持 SSE (Server-Sent Events) 流式传输
- 自定义 JSON-RPC 2.0 协议实现
- 支持元数据传递和中间件
**核心组件：**

```plain text
transport/jsonrpc/
├── core/
│   ├── jsonrpc.go       # JSON-RPC 协议核心
│   ├── connection.go    # 连接管理
│   ├── message.go       # 消息编解码
│   └── middleware.go    # 中间件支持
├── client/
│   ├── client.go        # HTTP 客户端
│   └── option.go        # 配置选项
└── server/
    ├── server.go        # HTTP 服务端
    └── option.go        # 配置选项

```

**消息格式：**

```json
// Request
{
  "jsonrpc": "2.0",
  "id": "req-123",
  "method": "message/send",
  "params": {
    "message": {
      "role": "user",
      "parts": [{"kind": "text", "text": "Hello"}]
    }
  }
}

// Response
{
  "jsonrpc": "2.0",
  "id": "req-123",
  "result": {
    "task": {
      "id": "task-456",
      "status": {"state": "completed"},
      ...
    }
  }
}

// Stream Event (SSE format)
data: {"message": {...}}

data: {"taskStatusUpdateEvent": {"status": {"state": "working"}}}

data: {"taskArtifactUpdateEvent": {"artifact": {...}, "lastChunk": true}}

```

