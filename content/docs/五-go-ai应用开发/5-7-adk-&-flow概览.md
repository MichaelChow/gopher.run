---
title: "5.7 adk & flow概览"
date: 2025-08-06T11:59:00Z
draft: false
weight: 5007
---

# 5.7 adk & flow概览



## **一、什么是 Agent、ADK？**

> Agent：通过结合大语言模型的理解能力（**ChatModel**）和预定义工具的执行能力（**Tool**），自主地完成复杂的任务。

Eino 把最常用的**大模型应用模式**封装成**简单、易用的工具：**

- `Flow/`：在该集成工具目录下提供了基于compose.Graph的 `ReAct Agent` 和 `Host Multi Agent`。
- `ADK/`：Eino Agent Development Kit，**基于Eino已有组件生态的面向Agent开发的框架**。参考 [Google-ADK](https://google.github.io/adk-docs/agents/) 的设计，**相较于Eino Graph 大幅简化了Agent、Multi-Agent开发**。它通过内置能力高效协调多智能体交互：跨智能体上下文传播、流式数据兼容与动态转换、任务控制权转移、中断与恢复机制、通用切面编程特性。适用场景广泛、模型无关、部署无关，并提供完善的生产级应用的治理能力，助力开发者搭建 **对话智能体、非对话智能体、复杂任务、工作流**等多种多样的 Agent 应用。
![](/images/24724637-29b5-80af-a6a1-e96b72202555/image_24724637-29b5-8089-87f7-e37a36c1713d.jpg)

```go
// eino-framework/eino/adk
adk/
├── 1. 核心接口定义
│   ├── interface.go          # 定义Agent、Message等核心接口和数据结构
│   └── instruction.go        # 指令相关的接口定义
│
├── 2. 基础工具和基础设施
│   ├── utils.go              # 异步迭代器、生成器等核心工具函数
│   ├── call_option.go        # 调用选项和配置管理
│   └── runctx.go             # 运行时上下文管理
│
├── 3. 核心代理实现
│   ├── chatmodel.go          # 聊天模型代理，处理AI对话和工具调用
│   ├── agent_tool.go         # 代理工具集成，支持工具调用功能
│   ├── flow.go               # 流程代理，管理代理间的消息流转
│   ├── workflow.go           # 工作流代理，支持顺序、并行、循环执行
│   └── react.go              # ReAct代理，实现推理和行动循环
│
├── 4. 执行和运行管理
│   ├── runner.go             # 代理运行器，管理代理的生命周期和执行
│   └── interrupt.go          # 中断处理，支持代理执行的中断和恢复
│
├── 5. 预构建组件
│   └── prebuilt/             # 预构建的代理和工具组件
│   ├──--- supervisor.go      # 监督者模式实现
│
└── 6. 测试文件
    ├── *_test.go             # 各模块的单元测试
    └── ...
```



### **Agent Interface**

**Agent Interface：**

实现此接口的 Struct 可被视为一个 Agent

1. 从入参中 AgentInput、AgentRunOption、Context Session（可选），**获取任务详情及相关数据**。
1. **执行任务**，并将**执行过程、执行结果**输出到 **AgenEvent Iterator**
1. 执行任务时，可通过 **Context 中的 Session 暂存数据**
```go
// github.com/cloudwego/eino/adk/interface.go
type Agent interface {
    Name(ctx context.Context) string  // Agent 的名称，作为 Agent 的标识
    Description(ctx context.Context) string  // 	Agent 的职能描述信息，主要用于让其他的 Agent 了解和判断该 Agent 的职责或功能
    Run(ctx context.Context, input *AgentInput) *AsyncIterator[*AgentEvent]  // Agent 的核心执行方法，返回一个迭代器，调用者可以通过这个迭代器持续接收 Agent 产生的事件
}
```



**AgentInput:**

Run 方法接收 AgentInput 作为 Agent 的输入

```go
type AgentInput struct {
	 Messages        []Message  // 同ChatModel，用户指令、对话历史、背景知识、样例数据等任何你希望传递给 Agent 的数据
	 EnableStreaming bool.  // 向 Agent的**建议其输出模式（并非一个强制性约束）：控制那些同时支持流式和非流式输出的组件的行为（如ChatModel），不会影响仅支持一种输出方式的组件**
}
```

```go
import (
    "github.com/cloudwego/eino/adk"
    "github.com/cloudwego/eino/schema"
)

input := &adk.AgentInput{
    Messages: []adk.Message{
       schema.UserMessage("What's the capital of France?"),
       schema.AssistantMessage("The capital of France is Paris.", nil),
       schema.UserMessage("How far is it from London? "),
    },
}
```



**AgentRunOption:**

由 Agent 实现定义，可以在请求维度修改 Agent 配置或者控制 Agent 行为。

```go
// github.com/cloudwego/eino/adk/call_option.go
// func WrapImplSpecificOptFn[T any](optFn func(*T)) AgentRunOption
// func GetImplSpecificOptions[T any](base *T, opts ...AgentRunOption) *T

import "github.com/cloudwego/eino/adk"

type options struct {
    modelName string
}
 
func WithModelName(name string) adk.AgentRunOption {   // 在请求维度修改调用的模型
    return adk.WrapImplSpecificOptFn(func(t *options) {   // adk.WrapImplSpecificOptFn
       t.modelName = name
    })
}

func (m *MyAgent) Run(ctx context.Context, input *adk.AgentInput, opts ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
    o := &options{}
    o = adk.GetImplSpecificOptions(o, opts...)   // adk.GetImplSpecificOptions
    // run code...
}

// adk.AgentRunOption 有一个 DesignateAgent 方法，调用该方法可以在调用多 Agent 系统时指定 Option 生效的 Agent。
```



**Agent.Run()**: 

```go
// github.com/cloudwego/eino/adk/utils.go

type AsyncIterator[T any] struct {   // 泛型结构体，迭代任何类型的数据。 
    ...
}

func (ai *AsyncIterator[T]) Next() (T, bool) {  // 阻塞式，每次调用 Next() ，程序会暂停执行，直到Agent 产生了一个新的 AgentEvent 或 Agent 主动关闭了迭代器（通常是Agent运行结束， 第二个返回值返回false）
    ...
}
```

AsyncIterator 常在 for 循环中处理：

```go
iter := myAgent.Run(xxx) // get AsyncIterator from Agent.Run

for {
    event, ok := iter.Next()
    if !ok {
        break
    }
    // handle event
}
```

AsyncIterator 可以由 `NewAsyncIteratorPair` 创建：

```go
// github.com/cloudwego/eino/adk/utils.go

func NewAsyncIteratorPair[T any]() (*AsyncIterator[T], *AsyncGenerator[T]) // 返回的 AsyncGenerator 用来生产数据
```

Agent.Run **返回*****AsyncIterator**[*AgentEvent]：**一个异步迭代器（异步指生产与消费之间没有同步控制）**，迭代器类型被固定为AsyncIterator[***AgentEvent**]，目的是旨在**让调用者实时地接收到 被调用Agent 产生的一系列 AgentEvent**。

**因此 Agent.Run 通常会在 Goroutine 中运行 Agent 从而立刻返回 AsyncIterator 供调用者监听（异步任务）**：产生新的**AgentEvent**时写入到 **Generator** 中，供 Agent 调用方在 **Iterator** 中消费。

```go
import "github.com/cloudwego/eino/adk"

func (m *MyAgent) Run(ctx context.Context, input *adk.AgentInput, opts ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
    // handle input
    iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()
    go func() {        // goroutine
       defer func() {
          // recover code
          gen.Close()
       }()
       // agent run code
       // gen.Send(event)
    }()
    return iter
}
```



**AgentWithOptions:**

支持Eino ADK Agent 做一些通用配置

```go
// github.com/cloudwego/eino/adk/flow.go
func AgentWithOptions(ctx context.Context, agent Agent, opts ...AgentOption) Agent
```



### **AgentEvent**

Agent 在其运行过程中产生的核心事件数据结构，包含了 Agent 的元信息、输出、行为和报错。

```go
// github.com/cloudwego/eino/adk/interface.go

type AgentEvent struct {
    AgentName string    //  产生了当前的 AgentEvent 的 Agent 实例
    RunPath []string // 当前 Agent 的完整调用链路（入口Agent到当前产生事件的所有AgentName）
    Output *AgentOutput
    Action *AgentAction
    Err error
}
```

**AgentOutput:**

```go
// github.com/cloudwego/eino/adk/interface.go

type AgentOutput struct {
    MessageOutput *MessageVariant // Message 输出

    CustomizedOutput any // 其他类型的输出
}

type MessageVariant struct {  // **核心结构**
    IsStreaming bool // 标志位,true -> MessageStream, false -> Message

    Message       Message
    MessageStream MessageStream
    // message role: Assistant or Tool
    Role schema.RoleType  // 消息的角色（常用的元数据）
    // only used when Role is Tool
    ToolName string  // 如果消息角色是 Tool ，这个字段会直接提供工具的名称（常用的元数据）
}
```

**AgentAction：**

控制多 Agent 协作，比如立刻退出、中断、跳转等

```go
type AgentAction struct {
	Exit bool // true -> Multi-Agent 会立刻退出

	Interrupted *InterruptInfo

	TransferToAgent *TransferToAgentAction // 跳转到目标 Agent 运行

	CustomizedAction any
}
```



### **多 Agent 协作: 上下文传递**

**History:**

每一个 Agent 产生的 AgentEvent 都会被保存到 History 中，在调用一个新 Agent 时(Workflow/ Transfer)，History 中的 AgentEvent 会被转换并拼接到 AgentInput 中。

默认情况下，其他 Agent 的 Assistant Message 或 Tool Message，被转换为 User Message, 避免了当前Agent的chatModel的的上下文混乱。

只有当 Event 的 RunPath “属于”当前 Agent 的 RunPath 时（ RunPathA 与 RunPathB 相同或者 RunPathA 是 RunPathB 的前缀），该 Event 才会参与构建 Agent 的 AgentInput。（过滤掉无关的AgentInput）

![](/images/24724637-29b5-80af-a6a1-e96b72202555/image_24b24637-29b5-80d2-b367-eee8840c6ab9.jpg)

**WithHistoryRewrite：**

通过 AgentWithOptions 可以自定义 Agent 从 History 中生成 AgentInput 的方式：

```go
// github.com/cloudwego/eino/adk/flow.go

type HistoryRewriter func(ctx context.Context, entries []*HistoryEntry) ([]Message, error)

func WithHistoryRewriter(h HistoryRewriter) AgentOption
```



**SessionValues:**

在一次运行中持续存在的全局临时 KV 存储，一次运行中的任何 Agent 可以在任何时间读写 SessionValues。

Eino ADK 提供了三种方法访问 SessionValues：

```go
// github.com/cloudwego/eino/adk/runctx.go
// 获取全部 SessionValues
func GetSessionValues(ctx context.Context) map[string]any
// 设置 SessionValues
func SetSessionValue(ctx context.Context, key string, value any)
// 指定 key 获取 SessionValues 中的一个值，key 不存在时第二个返回值为 false，否则为 true
func GetSessionValue(ctx context.Context, key string) (any, bool)
```



### **多Agent协作：Transfer SubAgents**

Agent 运行时产生带有**包含 TransferAction 的 AgentEvent 后**，Eino ADK 会调用 Action 指定的 Agent

**Transfer的含义是将****任务移交****给子 Agent，而不是委托或者分配（区别与ToolCall）**

TransferAction 可以使用 `NewTransferToAgentAction` 快速创建：

```go
import "github.com/cloudwego/eino/adk"

event := adk.NewTransferToAgentAction("dest agent name")
```



运行前需要先调用 `SetSubAgents` 将可能的子 Agent 注册到 Eino ADK 中：

```go
// github.com/cloudwego/eino/adk/flow.go
func SetSubAgents(ctx context.Context, agent Agent, subAgents []Agent) (Agent, error)
```



运行时动态地注册父子 Agent: 

如果 Agent 实现了 `OnSubAgents` 接口，`SetSubAgents` 中会调用相应的方法向 Agent 注册。

```go
// github.com/cloudwego/eino/adk/interface.go
type OnSubAgents interface {
    OnSetSubAgents(ctx context.Context, subAgents []Agent) error
    OnSetAsSubAgent(ctx context.Context, parent Agent) error

    OnDisallowTransferToParent(ctx context.Context) error
}
```

![](/images/24724637-29b5-80af-a6a1-e96b72202555/image_24724637-29b5-800d-bb70-d019f6664baf.jpg)

### **多Agent协作：Workflow**

提供了 顺序、并行和循环三种工作流模式，供用户灵活组合出不同的工作流图。

在 Workflow Agent 中，每个 Agent 拿到相同的 AgentInput 输入，按照预先设定好的拓扑结构所表达的顺序依次运行。

1. **Sequential**
- Agent 间协作方式：Transfer
- AgentInput 的上下文策略：上游 Agent 全对话
- 决策自主性：预设决策
将用户提供的 SubAgents 列表，组合成按照顺序依次执行的 Sequential Agent，其中的 Name 和 Description 作为 Sequential Agent 的名称标识和描述。

Sequential Agent 执行时，将 SubAgents 列表，按照顺序依次执行，直至将所有 Agent 执行一遍后结束。

注： 由于 Agent 只能获取到上游 Agent 的全对话，后执行的 Agent 看不到先执行的 Agent 的 AgentEvent 输出。

```go
type SequentialAgentConfig struct {
    Name        string
    Description string
    SubAgents   []Agent
}

func NewSequentialAgent(ctx context.Context, config *SequentialAgentConfig) (Agent, error) {
    // omit code
}
```

![](/images/24724637-29b5-80af-a6a1-e96b72202555/image_24724637-29b5-8036-95e3-f9775b548ba1.jpg)

1. **Parallel**
> Agent 间协作方式：Transfer
AgentInput 的上下文策略：上游 Agent 全对话
决策自主性：预设决策

将用户提供的 SubAgents 列表，组合成基于相同上下文，并发执行的 Parallel Agent，其中的 Name 和 Description 作为 Parallel Agent 的名称标识和描述。

Parallel Agent 执行时，将 SubAgents 列表，并发执行，待所有 Agent 执行完成后结束。

```go
// eino/adk/workflow.go
type ParallelAgentConfig struct {
    Name        string
    Description string
    SubAgents   []Agent
}

func NewParallelAgent(ctx context.Context, config *ParallelAgentConfig) (Agent, error) {
    // omit code
}
```

![](/images/24724637-29b5-80af-a6a1-e96b72202555/image_24724637-29b5-80aa-a519-cb80f0981b9f.jpg)

**3. Loop**

- Agent 间协作方式：Transfer
- AgentInput 的上下文策略：上游 Agent 全对话
- 决策自主性：预设决策
将用户提供的 SubAgents 列表，按照数组顺序依次执行，循环往复，组合成 Loop Agent，其中的 Name 和 Description 作为 Loop Agent 的名称标识和描述。

Sequential Agent 执行时，将 SubAgents 列表，并发执行，待所有 Agent 执行完成后结束。

```go
// eino/adk/workflow.go
type LoopAgentConfig struct {
    Name        string
    Description string
    SubAgents   []Agent

    MaxIterations int
}

func NewLoopAgent(ctx context.Context, config *LoopAgentConfig) (Agent, error) {
    // omit code
}
```

![](/images/24724637-29b5-80af-a6a1-e96b72202555/image_24724637-29b5-8084-9ff7-dc03c60bcf1c.jpg)

### **多Agent协作：AgentAsTool**

- Agent 间协作方式：ToolCall
- AgentInput 的上下文策略：全新任务描述
- 决策自主性：自主决策
将一个 Agent 转换成 Tool，被其他的 Agent 当成普通的 Tool 使用。

注：一个 Agent 能否将其他 Agent 当成 Tool 进行调用，取决于自身的实现。adk 中提供的 ChatModelAgent 支持 AgentAsTool 的功能

```go
// eino/adk/agent_tool.go
func NewAgentTool(_ context.Context, agent Agent, options ...AgentToolOption) tool.BaseTool {
    // omit code
}
```

下图展示了 Agent1 把 Agent2、Agent3 当成 Tool 进行调用的过程，类似 Function Stack Call，即在 Agent1 运行过程中，将 Agent2、Agent3 当成工具函数来进行调用。

- AgentAsTool 可作为 Supervisor Multi-Agent 的一种实现方式
    ![](/images/24724637-29b5-80af-a6a1-e96b72202555/image_24724637-29b5-8004-bc2d-fd76266fadd5.jpg)


### **Agent 扩展**

**Runner:** 

是 Eino ADK 中负责**执行 Agent 的核心引擎**。

主要作用是管理和控制 Agent 的整个生命周期，如处理多 Agent 协作，保存传递上下文等，interrupt、callback 等切面能力也均依赖 Runner 实现。任何 Agent 都应通过 Runner 来运行。

```go
// github.com/cloudwego/eino/adk/runners.go
type Runner struct {
	a               Agent
	enableStreaming bool
	store           compose.CheckPointStore
}
```

```go
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		EnableStreaming: true, // you can disable streaming here
		Agent:           a,
		CheckPointStore: newInMemoryStore(),
	})
```



**Interrupt Action & Resume Action:**

允许一个正在运行的 Agent 主动中断其执行，保存当前状态，并在稍后从中断点恢复执行。

这对于处理需要外部输入、长时间等待或可暂停的任务流非常有用。

产生包含 Interrupted Action 的 AgentEvent 来主动中断 Runner 的运行：

```go
// github.com/cloudwego/eino/adk/interface.go
type AgentAction struct {
    // other actions
    Interrupted *InterruptInfo
    // other actions
}

// github.com/cloudwego/eino/adk/interrupt.go
type InterruptInfo struct { // 附带自定义的中断信息（如向调用者说明中断原因等），传递给调用者、恢复Agent 运行时重新传递给中断的 Agent
    Data any
}
```



**状态持久化 (Checkpoint)：**

Runner 在终止运行后会将当前运行状态（原始输入、对话历史等）以及 Agent 抛出的 InterruptInfo 以 CheckPointID 为 key 持久化到 CheckPointStore 中。

为了保存 interface 中数据的原本类型，Eino ADK 使用 gob（[https://pkg.go.dev/encoding/gob](https://pkg.go.dev/encoding/gob)）序列化运行状态。因此在使用自定义类型时需要提前使用 gob.Register 或 **gob.RegisterName 注册类型**（更推荐后者，前者使用路径加类型名作为默认名字，因此类型的位置和名字均不能发生变更）。Eino 会**自动注册框架内置的类型**。

```go
// github.com/cloudwego/eino/adk/runner.go
type RunnerConfig struct {
    // other fields
    CheckPointStore CheckPointStore
}

// github.com/cloudwego/eino/adk/interrupt.go
type CheckPointStore interface {
    Set(ctx context.Context, key string, value []byte) error
    Get(ctx context.Context, key string) ([]byte, bool, error)
}

// github.com/cloudwego/eino/adk/interrupt.go
func WithCheckPointID(id string) AgentRunOption {
	return WrapImplSpecificOptFn(func(t *options) {
		t.checkPointID = &id
	})
}
```



**Resume:**

调用 Runner 的 Resume 接口传入中断时的 CheckPointID 可以恢复运行：

```go
// github.com/cloudwego/eino/adk/runner.go
func (r *Runner) Resume(ctx context.Context, checkPointID string, opts ...AgentRunOption) (*AsyncIterator[*AgentEvent], error)
```

恢复 Agent 运行需要发生中断的 Agent 实现了 ResumableAgent 接口， Runner 从 CheckPointerStore 读取运行状态并恢复运行，其中 InterruptInfo 和上次运行配置的 EnableStreaming 会作为输入提供给 Agent：

Resume 如果向 Agent 传入新信息，可以定义 AgentRunOption，在调用 Runner.Resume 时传入。

```go
// github.com/cloudwego/eino/adk/interface.go
type ResumableAgent interface {
    Agent

    Resume(ctx context.Context, info *ResumeInfo, opts ...AgentRunOption) *AsyncIterator[*AgentEvent]
}

// github.com/cloudwego/eino/adk/interrupt.go
type ResumeInfo struct {
		EnableStreaming bool
		*InterruptInfo
}
```



### Eino React Agent（基于compose.Graph）

> 官网：[https://react-lm.github.io/](https://react-lm.github.io/)

**ReAct**（Reasoning + Acting）是一种 AI Agent模式：用户输入 → 模型推理 → 工具调用 → 结果反馈 → 模型推理 → ... → 最终答案。

- **推理阶段（Reasoning）：**AI 模型分析用户问题、决定是否需要调用工具、选择合适的工具和参数
- **行动阶段（Acting）：**执行选定的工具、获取工具执行结果、将结果作为上下文传递给下一轮推理当ChatModelAgent没有配置工具时，退化为一次 ChatModel 调用。
![](/images/24724637-29b5-80af-a6a1-e96b72202555/image_24b24637-29b5-80b4-b588-caf259473f8a.jpg)



**Eino React Agent** 是实现了 [React 模式](https://react-lm.github.io/) 的Agent框架，用户可以用来快速灵活地构建并调用 React Agent。

React Agent 底层使用 compose.Graph 作为编排方案，主要包含两个节点：**ChatModel、Tools**

**节点拓扑&数据流图：**

![](/images/24724637-29b5-80af-a6a1-e96b72202555/image_24224637-29b5-80fe-8dc4-f1252fad003b.jpg)

用户输入 → ChatModel → Tools → ChatModel → ... → 最终结果

- 所有历史消息都会放入 state 中
- 在传递给 ChatModel 前，会通过 MessageModifier 处理消息
- 直到 ChatModel 返回的消息中不再有 tool call，则返回最终消息


通过 flow/agent/react 包提供完整实现：

```go
// eino/flow/agent/react/react.go
type Agent struct {
    runnable         compose.Runnable[[]*schema.Message, *schema.Message]
    graph            *compose.Graph[[]*schema.Message, *schema.Message]
    graphAddNodeOpts []compose.GraphAddNodeOpt
}

type AgentConfig struct {
    ToolCallingModel model.ToolCallingChatModel  // 支持工具调用的聊天模型
    ToolsConfig      compose.ToolsNodeConfig     // 工具配置
    MessageModifier  MessageModifier             // 消息修改器
    MaxStep          int                         // 最大步数限制
    ToolReturnDirectly map[string]struct{}       // 直接返回的工具
    StreamToolCallChecker func(...)              // 流式工具调用检查器
}
```



使用示例：

```go
// 创建 ReAct 代理
agent, err := react.NewAgent(ctx, &react.AgentConfig{
    ToolCallingModel: chatModel,
    ToolsConfig: compose.ToolsNodeConfig{
        Tools: []tool.BaseTool{restaurantTool, dishTool},
    },
    MaxStep: 12,
})

// 使用代理
msg, err := agent.Generate(ctx, []*schema.Message{
    {Role: schema.User, Content: "推荐北京最好的川菜餐厅"},
})
```

**带中断的 ReAct:**

```go
// 支持用户干预的 ReAct 实现
for {
    result, err := runner.Invoke(ctx, input)
    if err == nil {
        fmt.Printf("最终结果: %s", result.Content)
        break
    }
    
    // 处理中断，允许用户修改工具调用参数
    info, ok := compose.ExtractInterruptInfo(err)
    if ok {
        // 用户确认或修改工具调用参数
        // 继续执行
    }
}
```



**初始化配置**

**基础配置：**

```go
import (
    "github.com/cloudwego/eino-ext/components/model/openai"
    "github.com/cloudwego/eino/flow/agent/react"
    "github.com/cloudwego/eino/components/tool"
    "github.com/cloudwego/eino/schema"
)

// 创建 OpenAI 模型
openaiModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
    APIKey: os.Getenv("OPENAI_API_KEY"),
    Model:  "gpt-3.5-turbo",
})

// 创建工具
weatherTool := tool.NewTool("get_weather", "获取天气信息", func(ctx context.Context, input string) (string, error) {
    return "晴天，温度25°C", nil
})

// 创建 React Agent
agent, err := react.NewAgent(ctx, &react.AgentConfig{
    ToolCallingModel: openaiModel,
    ToolsConfig: compose.ToolsNodeConfig{
        Tools: []tool.BaseTool{weatherTool},
    },
    MaxStep: 12, // 最大步数限制
})
```



**高级配置选项:**

```go
agent, err := react.NewAgent(ctx, &react.AgentConfig{
    ToolCallingModel: openaiModel,
    ToolsConfig: compose.ToolsNodeConfig{
        Tools: []tool.BaseTool{weatherTool, searchTool},
    },
    MessageModifier: func(ctx context.Context, messages []*schema.Message) []*schema.Message {
        // 自定义消息处理逻辑
        return messages
    },
    MaxStep: 12,
    ToolReturnDirectly: map[string]struct{}{
        "final_answer": {}, // 某些工具调用后直接返回
    },
    StreamToolCallChecker: customToolCallChecker, // 自定义流式工具调用检查器
})
```

**调用方式:**

**1. Generate（同步调用）**

```go
outMessage, err := agent.Generate(ctx, []*schema.Message{
    schema.UserMessage("北京今天天气怎么样？"),
})

if err != nil {
    log.Fatal(err)
}

fmt.Println("回答:", outMessage.Content)
```

1. **Stream（流式调用）**
```go
msgReader, err := agent.Stream(ctx, []*schema.Message{
    schema.UserMessage("写一个 golang 的 hello world 程序"),
})

if err != nil {
    log.Fatal(err)
}

for {
    msg, err := msgReader.Recv()
    if err != nil {
        if errors.Is(err, io.EOF) {
            break
        }
        log.Printf("接收错误: %v\n", err)
        return
    }
    
    fmt.Print(msg.Content)
}
```

**3. 使用回调（Callbacks）**

```go
// 构建回调处理器
callback := react.BuildAgentCallback(
    &template.ModelCallbackHandler{
        OnStart: func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
            fmt.Printf("模型开始处理: %s\n", info.Name)
            return ctx
        },
        OnEnd: func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
            fmt.Printf("模型处理完成: %s\n", info.Name)
            return ctx
        },
    },
    &template.ToolCallbackHandler{
        OnStart: func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
            fmt.Printf("工具开始执行: %s\n", info.Name)
            return ctx
        },
        OnEnd: func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
            fmt.Printf("工具执行完成: %s\n", info.Name)
            return ctx
        },
    },
)

// 使用回调调用
outMessage, err := agent.Generate(ctx, []*schema.Message{
    schema.UserMessage("查询天气"),
}, react.WithComposeOptions(compose.WithCallbacks(callback)))
```



**在 Graph/Chain 中使用**

Agent 可以作为 Lambda 嵌入到其他 Graph 中：

```go
agent, _ := react.NewAgent(ctx, &react.AgentConfig{
    ToolCallingModel: chatModel,
    ToolsConfig: compose.ToolsNodeConfig{
        Tools: []tool.BaseTool{weatherTool, searchTool},
    },
    MaxStep: 40,
})

// 创建 Chain
chain := compose.NewChain[[]*schema.Message, string]()
agentLambda, _ := compose.AnyLambda(agent.Generate, agent.Stream, nil, nil)

chain.
    AppendLambda(agentLambda).
    AppendLambda(compose.InvokableLambda(func(ctx context.Context, input *schema.Message) (string, error) {
        fmt.Println("获得 Agent 响应:", input.Content)
        return input.Content, nil
    }))

r, _ := chain.Compile(ctx)
res, _ := r.Invoke(ctx, []*schema.Message{{Role: schema.User, Content: "hello"}})
```



**运行过程分析**

当用户输入："我在海淀区，给我推荐一些菜，需要有口味辣一点的菜，至少推荐有 2 家餐厅"

**第一步：ChatModel 推理**

- 模型分析用户需求
- 决定调用 query_restaurants 工具
- 参数：{"location":"海淀区","topn":2}
**第二步：Tools 执行**

- 执行餐厅查询工具
- 返回 2 家海淀区餐厅信息
**第三步：ChatModel 再次推理**

- 基于餐厅信息，决定查询菜品
- 为每个餐厅调用 query_dishes 工具
- 并发执行多个工具调用
**第四步：Tools 并发执行**

- 同时查询两家餐厅的菜品
- 返回详细的菜品信息
**第五步：ChatModel 最终整合**

- 整合所有信息
- 生成最终推荐结果


## 二、Tool

### 基本结构

一个 agent 要调用 tool，需要有两步：

1. 大模型根据 tool 的功能和参数需求构建调用参数 → 需要tool 的功能介绍和调用这个 tool 所需要的参数信息，要求大模型能理解生成的function call 参数是否符合约束
1. 实际调用 tool
源码解读：

```go
type BaseTool interface {
    Info(ctx context.Context) (*schema.ToolInfo, error) // 要求有 Info() 接口返回 schema.ToolInfo
}

type InvokableTool interface {. // 一次性返回
    BaseTool

    // InvokableRun call function with arguments in JSON format
    InvokableRun(ctx context.Context, argumentsInJSON string, opts ...Option) (string, error)
}

type StreamableTool interface { // 流式返回
    BaseTool

    StreamableRun(ctx context.Context, argumentsInJSON string, opts ...Option) (*schema.StreamReader[string], error)
}
```



两种参数约束的表达：

- 方式 1 - map[string]*ParameterInfo：用map，key 即为参数名，value 则是这个参数的详细约束。简单直观，当参数由开发者通过编码的方式手动维护时常用。
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

- 方式2：**openapi3.Schema。**一般不由开发者自行构建此结构体，而是使用一些方法来生成。
    - 使用 GoStruct2ParamsOneOf 生成
    - 通过 openapi.json 文件生成


> **如何创建一个 tool：**[https://cloudwego.cn/zh/docs/eino/core_modules/components/tools_node_guide/how_to_create_a_tool/](https://cloudwego.cn/zh/docs/eino/core_modules/components/tools_node_guide/how_to_create_a_tool/)

### **方式1：直接实现 Tool 接口**

由于 tool 的定义都是接口，因此最直接实现一个 tool 的方式即实现接口。

对于需要更多自定义逻辑的场景，可以通过实现 Tool 接口来创建。

以 InvokableTool 为例：

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



### **方式2：utils.InferTool把本地函数转为 tool**

如我们代码中本身已经有了一个 `AddUser` 的方法，但为了让大模型可以自主决策如何调用这个方法，我们经常需要把这个方法变成一个 tool 并 bind 到大模型上。

Eino 中提供了 `NewTool` 的方法来把一个函数转成 tool，同时，针对为参数约束通过结构体的 tag 来表示的场景提供了 InferTool 的方法，让构建的过程更加简单。



**方式一：使用 NewTool 构建。**适合简单的工具实现，通过定义工具信息和处理函数来创建 Tool

**不足**：需要在 ToolInfo 中手动定义参数信息（ParamsOneOf），和实际的参数结构（TodoAddParams）是分开定义的。这样不仅造成了代码的冗余，而且在参数发生变化时需要同时修改两处地方，容易导致不一致，维护起来也比较麻烦。



当一个函数满足下面这种函数签名时，就可以用 NewTool 把其变成一个 InvokableTool ：

> 同理 NewStreamTool 可创建 StreamableTool

```go
type InvokeFunc[T, D any] func(ctx context.Context, input T) (output D, err error)
```

NewTool 的方法如下：

```go
// 代码见: github.com/cloudwego/eino/components/tool/utils/invokable_func.go
func NewTool[T, D any](desc *schema.ToolInfo, i InvokeFunc[T, D], opts ...Option) tool.InvokableTool
```



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



**方式二：使用 InferTool 构建:**更加简洁，通过结构体的 tag 来定义参数信息，就能实现参数结构体和描述信息同源，无需维护两份信息：

从 NewTool 中可以看出，构建一个 tool 的过程需要分别传入 ToolInfo 和 InvokeFunc， ParamsOneOf 的部分和 InvokeFunc 的 input 参数需要保持一致。

- ToolInfo 中包含 ParamsOneOf 的部分，这代表着函数的入参约束
- InvokeFunc 的函数签名中也有 input 的参数


更优雅的解决方法是 “参数约束直接维护在 input 参数类型定义中”，可参考上方 `GoStruct2ParamsOneOf` 的介绍。

当参数约束信息包含在 input 参数类型定义中时，就可以使用 InferTool 来实现，函数签名如下：

```go
func InferTool[T, D any](toolName, toolDesc string, i InvokeFunc[T, D], opts ...Option) (tool.InvokableTool, error)
```

```go
import (
    "context"

    "github.com/cloudwego/eino/components/tool/utils"
)

// 参数结构体
type TodoUpdateParams struct {
    ID        string  `json:"id" jsonschema:"description=id of the todo"`
    Content   *string `json:"content,omitempty" jsonschema:"description=content of the todo"`
    StartedAt *int64  `json:"started_at,omitempty" jsonschema:"description=start time in unix timestamp"`
    Deadline  *int64  `json:"deadline,omitempty" jsonschema:"description=deadline of the todo in unix timestamp"`
    Done      *bool   `json:"done,omitempty" jsonschema:"description=done status"`
}

// 处理函数
func UpdateTodoFunc(_ context.Context, params *TodoUpdateParams) (string, error) {
    // Mock处理逻辑
    return `{"msg": "update todo success"}`, nil
}

// 使用 InferTool 创建工具
updateTool, err := utils.InferTool(
    "update_todo", // tool name 
    "Update a todo item, eg: content,deadline...", // tool description
    UpdateTodoFunc)
```



**方式三**：使用 **InferOptionableTool** 方法

Option 机制是 Eino 提供的一种在运行时传递动态参数的机制。当开发者要实现一个需要自定义 option 参数时则可使用 InferOptionableTool 这个方法，相比于 InferTool 对函数签名的要求，这个方法的签名增加了一个 option 参数，签名如下：

```go
func InferOptionableTool[T, D any](toolName, toolDesc string, i OptionableInvokeFunc[T, D], opts ...Option) (tool.InvokableTool, error)
```



### 方式3：直接**使用 eino-ext 中的提供的通用tool**

> [https://github.com/cloudwego/eino-ext/tree/main/components/tool](https://github.com/cloudwego/eino-ext/tree/main/components/tool)

```go
import (
    "github.com/cloudwego/eino-ext/components/tool/duckduckgo"
)


// 创建 duckduckgo Search 工具
searchTool, err := duckduckgo.NewTool(ctx, &duckduckgo.Config{})
```



### 方式4：**使用 MCP 协议**

MCP（Model Context Protocol）是一个开放的模型上下文协议，现在越来越多的工具和平台都在基于这套协议把自身的能力暴露给大模型调用，eino 可以把基于 MCP 提供的可调用工具作为 tool，这将极大扩充 tool 的种类。

在 Eino 中使用 MCP 提供的 tool 非常方便：

```go
import (
    "fmt"
    "log"
    "context"
    "github.com/mark3labs/mcp-go/client"
    mcpp "github.com/cloudwego/eino-ext/components/tool/mcp"
)

func getMCPTool(ctx context.Context) []tool.BaseTool {
        cli, err := client.NewSSEMCPClient("http://localhost:12345/sse")
        if err != nil {
                log.Fatal(err)
        }
        err = cli.Start(ctx)
        if err != nil {
                log.Fatal(err)
        }

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

        tools, err := mcpp.GetTools(ctx, &mcpp.Config{Cli: cli})
        if err != nil {
                log.Fatal(err)
        }

        return tools
}
```



