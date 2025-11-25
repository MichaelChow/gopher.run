---
title: "5.7 adk基础"
date: 2025-08-06T11:59:00Z
draft: false
weight: 5007
---

# 5.7 adk基础



# **一、什么是ADK**

Agent：一个独立的、可执行的AI任务单元，通过调用ChatModel的理解能力和预定义Tool的工具执行能力，能够自主学习完成复杂的任务。主要功能：

- 推理：分析数据、识别模式，根据逻辑和可用信息推导出结论
- 行动：执行任务
- 观察：自助收集上下文信息
- 规划：确定必要的步骤，选择最佳行动方法
- 协作：与其他AI Agent/人 进行协作


**Eino ADK**（Agent Development Kit）****是一个专为 Go 语言设计的 Agent 和 Multi-Agent 开发框架，设计上参考了 [Google-ADK](https://google.github.io/adk-docs/agents/) 中对 Agent 与协作机制的定义。该工具库提供了统一的抽象接口、灵活的组合模式和强大的协作机制，设计哲学是"简单的事情简单做，复杂的事情也能做"，**让开发者能够专注于业务逻辑的实现，而不必担心底层的技术复杂性**（如跨Agent的context传播、事件流分发和转换、任务控制权转移、中断与恢复、callback通用能力），能像搭建乐高积木一样构建复杂的AI Agent系统：

- **少写胶水**：统一接口与事件流，复杂任务拆解更自然。
- **快速编排**：预设范式 + 工作流，分分钟搭好管线。
- **更可控**：可中断、可恢复、可审计，Agent 协作过程“看得见”。
<!-- 列布局开始 -->

ADK整体模块构成：

![](/images/24724637-29b5-80af-a6a1-e96b72202555/image_29924637-29b5-8072-9f82-fc688b7a0da8.jpg)




---

封装关系：

![](/images/24724637-29b5-80af-a6a1-e96b72202555/image_29624637-29b5-8056-975d-e9cb44249dfb.jpg)



<!-- 列布局结束 -->

# 二、预定义组件

## 2.1 `ChatModelAgent 带React能力`

**最重要的预构建组件**，封装了与大语言模型的交互逻辑，实现了经典的 **ReAct**（**Reason-Act-Observe**）模式。**ChatModlAgent的行为是 非确定性的，通过LLM来动态决定 call tool/transfer another agent**。运行过程为:

<!-- 列布局开始 -->

1. 调用 LLM（Reason）
1. LLM 返回工具调用请求（Action）
1. ChatModelAgent 执行工具（Act）
1. 将工具结果返回给 LLM（Observation），结合之前的上下文**循环（loop）**生成，直到模型判断不需要调用 Tool 后结束。



---

![](/images/24724637-29b5-80af-a6a1-e96b72202555/image_29924637-29b5-80f1-bdc9-efa10e7bea1a.jpg)

<!-- 列布局结束 -->

ReAct 模式的核心是**“思考 → 行动 → 观察 → 再思考”**的闭环，解决传统 Agent **“盲目行动”**(如一次性搜集全部信息导致的信息过载)**或“推理与行动脱节”**（如凭空靠直觉决策）的痛点。

```go
// eino/adk/chatmodel.go
type ChatModelAgent struct {
	name        string
	description string
	instruction string

	model       model.ToolCallingChatModel
	toolsConfig ToolsConfig

	genModelInput GenModelInput

	outputKey     string
	maxIterations int

	subAgents   []Agent // subAgents
	parentAgent Agent

	disallowTransferToParent bool

	exit tool.BaseTool

	// runner
	once   sync.Once  // sync.Once
	run    runFunc
	frozen uint32
}
```



example：使用 ADK 快速构建具有 `ReAct` 能力的 `ChatModelAgent`

```go
import github.com/cloudwego/eino/adk

// 创建一个包含多个工具的 ReAct ChatModelAgent
chatAgent := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
    Name:        "intelligent_assistant",
    Description: "An intelligent assistant capable of using multiple tools to solve complex problems",
    Instruction: "You are a professional assistant who can use the provided tools to help users solve problems",
    Model:       openaiModel,
    ToolsConfig: adk.ToolsConfig{
        Tools: []tool.BaseTool{
            searchTool,
            calculatorTool,
            weatherTool,
        },
    }
})
```

备注：在adk上线前，`Flow/`集成工具目录下提供了基于compose.Graph的 `ReAct Agent` 和 `Host Multi Agent`。（推荐使用新的基于adk的统一定义的版本）

```go
// eino-framework/eino/adk
adk/
├── 1. 核心接口定义
│   ├── interface.go          # 定义Agent、Message等核心接口和数据结构
│   └── instruction.go        # 指令相关的接口定义
│
├── 2. 基础工具和基础设施
│   ├── utils.go              #** 异步迭代器、生成器**等核心工具函数
│   ├── call_option.go        # 调用选项和配置管理
│   └── runctx.go             # 运行时上下文管理
│
├── 3. 核心Agent实现
│   └── react.go              # ReAct Agent，实现推理和行动循环
│   ├── chatmodel.go          # ChatModel Agent，**基于上述react.go的ReactAgent，处理AI对话和工具调用**
│   ├── agent_tool.go         # 代理工具集成，支持工具调用功能
│   ├── flow.go               # 流程Agent，管理代理间的消息流转
│   ├── workflow.go           # Workflow（工作流-精确流水线）Agent，支持 顺序、并发、循环 控制 子Agent 可预测的 确定性执行流程
│   └── Custom Agent          # 通过接口实现自己的 Agent，允许定义高度定制的复杂 Agent  
│
├── 4. 执行和运行管理
│   ├── runner.go             # 代理运行器，管理代理的生命周期和执行
│   └── interrupt.go          # 中断处理，支持代理执行的中断和恢复
│
├── 5. 预构建组件
│   └── prebuilt/             # 预构建的代理和工具组件
│   ├──--- supervisor.go      # 监督者模式实现：监督者Agent控制所有通信流程和任务委托，并根据当前上下文和任务需求决定调用哪个Agent。
│   ├──--- plan_execute.go		 # 计划-执行-反思 模式：Plan Agent 生成含多个步骤的计划，Execute Agent 根据用户 query 和计划来完成任务。Execute 后会再次调用 Plan，决定完成任务 / 重新进行规划。
│
└── 6. 测试文件
    ├── *_test.go             # 各模块的单元测试
    └── ...
```

[https://www.cloudwego.io/zh/docs/eino/core_modules/eino_adk/agent_implementation/chat_model/](https://www.cloudwego.io/zh/docs/eino/core_modules/eino_adk/agent_implementation/chat_model/)

## 2.2 `WorkflowAgent: 精密的agent流水线`

区别于基于 **LLM自主决策** 的Transfer（不确定的执行流），Workflow Agents模式 采用**预设决策（代码定义的执行流，可预测、可控制）**的方式来运行子Agent。

可基于 **Sequential Agent（顺序）、Parallel Agent（并发）、Loop Agent（循环）三种基础的 Workflow Agent执行模式 进行组合嵌套，构建各种复杂的执行流程。**

### **1. Sequential Agent（顺序）**

- **线性执行**：**最基础的Workflow Agent**，严格按照SubAgents数组的顺序，依次执行一次注册的Agents后结束。
- **History 传递**：每个 Agent 的执行结果都会被添加到 History 中，后续 Agent 可以访问前面 Agent 的执行历史，形成一个线性的执行链。
- **支持 提前退出**：如果任何一个子 Agent 产生退出 / 中断动作，整个 Sequential 流程会立即终止。


可能的应用场景：

- **数据 ETL**：`ExtractAgent`（从 MySQL 抽取订单数据）→ `TransformAgent`（清洗空值、格式化日期）→ `LoadAgent`（加载到数据仓库）
- **CI / CD 流水线**：`CodeCloneAgent`（从代码仓库拉取代码）→`UnitTestAgent`（运行单元测试，用例失败时返回错误与分析报告）→`CompileAgent`（编译代码）→`DeployAgent`（部署到目标环境）
<!-- 列布局开始 -->

![](/images/24724637-29b5-80af-a6a1-e96b72202555/image_29d24637-29b5-802c-9b0e-cb0132fc463b.jpg)


---

example：

```go
import github.com/cloudwego/eino/adk

// 依次执行 制定研究计划 -> 搜索资料 -> 撰写报告
sequential := adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{
    Name: "research_pipeline",
    Description: "",
    SubAgents: []adk.Agent{
        planAgent,    // 制定研究计划
        searchAgent,  // 搜索资料
        writeAgent,   // 撰写报告
    },
})

func NewPlanAgent() adk.Agent {
	a, err := adk.NewChatModelAgent(...) // 底层使用的ChatModelAgent
	return a
}
```

<!-- 列布局结束 -->



### **2. Loop Agent（循环）**

配置中注册的 Agents基于 SequentialAgent 实现，循环的顺序依次执行配置中注册的 Agents，直到达到最大迭代次数 或 某个子 Agent 产生 ExitAction。

- **循环执行**：重复执行 SubAgents 序列，每次循环都是一个完整的 Sequential Agent 执行过程。
- **History 累积**：每次迭代的结果都会累积到 History 中，后续迭代可以访问所有历史信息。
- **条件退出**：支持通过输出包含 `ExitAction` 的事件或达到最大迭代次数来终止循环，配置 `MaxIterations=0` 时表示无限循环。

特别适用于需要 **迭代优化、反复处理或持续监控******的场景；可能的应用场景有：

- **数据同步**：`CheckUpdateAgent`（检查源库增量）→ `IncrementalSyncAgent`（同步增量数据）→ `VerifySyncAgent`（验证一致性）
- **压力测试**：`StartClientAgent`（启动测试客户端）→ `SendRequestsAgent`（发送请求）→ `CollectMetricsAgent`（收集性能指标）


<!-- 列布局开始 -->



![](/images/24724637-29b5-80af-a6a1-e96b72202555/image_29d24637-29b5-8025-a317-ec4c9f81bf64.jpg)




---

example:

```go
import github.com/cloudwego/eino/adk

// 循环执行 5 次，每次顺序为：分析当前状态 -> 提出改进方案 -> 验证改进效果
loop := adk.NewLoopAgent(ctx, &adk.LoopAgentConfig{
    Name: "iterative_optimization",
    SubAgents: []adk.Agent{
        analyzeAgent,  // 分析当前状态
        improveAgent,  // 提出改进方案
        validateAgent, // 验证改进效果
    },
    MaxIterations: 5,
})

```

<!-- 列布局结束 -->



### **3. Parallel Agent（并发）**

基于相同的输入上下文，并发执行配置中注册的 所有Agents，并等待全部完成后结束。

- **并发执行**：所有子 Agent 同时启动，在独立的 goroutine 中并行执行。Parallel 内部默认包含异常处理机制：
    - **Panic 恢复**：每个 goroutine 都有独立的 panic 恢复机制
    - **错误隔离**：单个子 Agent 的错误不会影响其他子 Agent 的执行
    - **中断处理**：支持子 Agent 的中断和恢复机制
- **共享输入**：所有子 Agent 接收调用 Pararllel Agent 相同的初始输入和上下文。
- **等待与结果聚合**：内部使用 sync.WaitGroup 等待所有子 Agent 执行完成，收集所有子 Agent 的执行结果并按接收顺序输出到 `AsyncIterator` 中。


特别适用于可以独立并行处理的任务，能够显著提高执行效率；

可能的应用场景：

- **多源数据采集**：`MySQLCollector`（采集用户表）+ `PostgreSQLCollector`（采集订单表）+ `MongoDBCollector`（采集商品评论）
- **多渠道推送**：`WeChatPushAgent`（推送到微信公众号）+ `SMSPushAgent`（发送短信）+ `AppPushAgent`（推送到 APP）


<!-- 列布局开始 -->



![](/images/24724637-29b5-80af-a6a1-e96b72202555/image_24724637-29b5-80aa-a519-cb80f0981b9f.jpg)




---

example:

```go
import github.com/cloudwego/eino/adk

// 并发执行 情感分析 + 关键词提取 + 内容摘要
parallel := adk.NewParallelAgent(ctx, &adk.ParallelAgentConfig{
    Name: "multi_analysis",
    Description: "",
    SubAgents: []adk.Agent{
        sentimentAgent,  // 情感分析
        keywordAgent,    // 关键词提取
        summaryAgent,    // 内容摘要
    },
})
```

<!-- 列布局结束 -->





# 三、基础设计

### **统一的 Agent 抽象**

ADK 的核心是一个简洁而强大的 `Agent` 接口，无论是简单的问答机器人，还是复杂的多步骤任务处理系统，都可以通过这个统一的接口加以实现。

```go
// github.com/cloudwego/eino/adk/interface.go
type Agent interface {
    Name(ctx context.Context) string  // **明确的身份**
    Description(ctx context.Context) string  // 	**清晰的职责**
    Run(ctx context.Context, input *AgentInput) ***AsyncIterator[*AgentEvent]**  // **标准化的执行方式**。返回一个**异步迭代器****（生产与消费之间没有同步控制）**。调用者可以通过这个 **AgenEvent Iterator****（**迭代器） 持续接收 Agent 产生的事件
    // 执行任务时可通过 **Context 中的 Session 暂存数据**
}
```



```go
// **AgentInput:**
type AgentInput struct {
	 Messages        []Message  // 同ChatModel，用户指令、对话历史、背景知识、样例数据等任何你希望传递给 Agent 的数据
	 EnableStreaming boo.  // 向 Agent的**建议其输出模式（并非一个强制性约束）：控制那些同时支持流式和非流式输出的组件的行为（如ChatModel），不会影响仅支持一种输出方式的组件**
}

input := &adk.AgentInput{
    Messages: []adk.Message{
       schema.UserMessage("What's the capital of France?"),
       schema.AssistantMessage("The capital of France is Paris.", nil),
       schema.UserMessage("How far is it from London? "),
    },
}

// **AgentRunOption**有一个 DesignateAgent 方法，调用该方法可以在调用多 Agent 系统时指定 Option 生效的 Agent。****func (m *MyAgent) Run(ctx context.Context, input *adk.AgentInput, opts ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
    o := &options{}
    o = adk.GetImplSpecificOptions(o, opts...)   // adk.GetImplSpecificOptions
    // run code...
}
```



```go
// AsyncIterator声明：
// github.com/cloudwego/eino/adk/utils.go
type AsyncIterator[T any] struct {   // 泛型结构体，迭代任何类型的数据。 
    ...
}
func (ai *AsyncIterator[T]) Next() (T, bool) {  // 阻塞式，每次调用 Next() ，程序会暂停执行，直到Agent 产生了一个新的 AgentEvent 或 Agent 主动关闭了迭代器（通常是Agent运行结束， 第二个返回值返回false）
    ...
}
// AsyncIterator 可以由 NewAsyncIteratorPair 创建：
func NewAsyncIteratorPair[T any]() (*AsyncIterator[T], *AsyncGenerator[T]) // 返回的 AsyncGenerator 用来生产数据


// AsyncIterator使用：常在 for 循环中处理
iter := agent.Run(xxx) // get AsyncIterator from Agent.Run

for {
    event, ok := iter.Next()
    if !ok {
        break
    }
    // handle event
}
```



**Agent.Run 通常会**** 在新的Goroutine中运行Agent，从而立刻返回AsyncIterator供调用者监听（异步任务****）**：产生新的**AgentEvent**时写入到 **Generator** 中，供 Agent 调用方在 **Iterator** 中消费:

```go
import "github.com/cloudwego/eino/adk"

func (m *Agent) Run(ctx context.Context, input *adk.AgentInput, opts ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
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

// AgentWithOptions: 支持Eino ADK Agent 做一些通用配置
// github.com/cloudwego/eino/adk/flow.go
func AgentWithOptions(ctx context.Context, agent Agent, opts ...AgentOption) Agent
```



### **异步事件驱动架构**`AsyncIterator[*AgentEvent]`

ADK 采用了异步事件流设计，通过 `AsyncIterator[*AgentEvent]` 实现**非阻塞的事件处理（unbuffed chan）**，并通过 `Runner` 框架运行 Agent：

- **实时响应**：`AgentEvent` 包含 Agent 执行过程中特定节点输出（Agent 回复、工具处理结果等等），用户可以立即看到 Agent 的思考过程和中间结果。
- **追踪执行过程**：`AgentEvent` 额外携带状态修改动作与运行轨迹，便于开发调试和理解 Agent 行为。
- **自动流程控制**：框架通过 `Runner` 自动处理中断、跳转、退出行为，无需用户额外干预。


Agent在其运行过程中产生的核心事件数据结构：

```go
// github.com/cloudwego/eino/adk/interface.go

type AgentEvent struct {
    AgentName string   
    RunPath []string // 当前 Agent 的完整调用链路（入口Agent到当前产生事件的所有AgentName）
    Output *AgentOutput
    Action *AgentAction
    Err error
}

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

// 控制多 Agent 协作，比如立刻退出、中断、跳转等
type AgentAction struct {
	Exit bool // true -> Multi-Agent 会立刻退出
	Interrupted *InterruptInfo
	TransferToAgent *TransferToAgentAction // 跳转到目标 Agent 运行
	CustomizedAction any
}
```

### **灵活的协作机制: 共享Session、移交运行Transfer、显式调用ToolCall**

Eino ADK 支持处于同一个系统内的 Agent 之间以多种方式进行协作（交换数据或触发运行）：

1. **共享 Session**：单次运行过程中持续存在的 KV 存储，用于支持跨 Agent 的状态管理和数据共享。
    ```go
    // 获取全部 SessionValues
    func GetSessionValues(ctx context.Context) map[string]any
    // 指定 key 获取 SessionValues 中的一个值，key 不存在时第二个返回值为 false，否则为 true
    func GetSessionValue(ctx context.Context, key string) (any, bool)
    // 添加 SessionValues。原SetSessionValue更名
    func AddSessionValue(ctx context.Context, key string, value any)
    // 批量添加 SessionValues
    func AddSessionValues(ctx context.Context, kvs map[string]any)
    ```
1. **移交运行（Transfer）**：携带本 Agent 输出结果上下文，将任务移交至子 Agent 继续处理。适用于智能体功能可以清晰的划分边界与层级的场景，常结合 ChatModelAgent 使用，通过 LLM 的生成结果进行动态路由。结构上，以此方式进行协作的两个 Agent 称为父子 Agent：
<!-- 列布局开始 -->

![](/images/24724637-29b5-80af-a6a1-e96b72202555/image_29a24637-29b5-8030-ac7a-f6483beec6f7.jpg)


---

example:

```go
// 设置父子 Agent 关系
func SetSubAgents(ctx context.Context, agent Agent, subAgents []Agent) (Agent, error)

// 指定目标 Agent 名称，构造 Transfer Event
func NewTransferToAgentAction(destAgentName string) *AgentAction
```



<!-- 列布局结束 -->

    

<!-- 列布局开始 -->

![](/images/24724637-29b5-80af-a6a1-e96b72202555/image_24b24637-29b5-80d2-b367-eee8840c6ab9.jpg)




---

- 每一个 Agent 产生的 AgentEvent 都会被保存到 History 中，在调用一个新 Agent 时(Workflow/ Transfer)，History 中的 AgentEvent 会被转换并拼接到 AgentInput 中。
    - 默认情况下，其他 Agent 的 Assistant Message 或 Tool Message，被转换为 User Message, 避免了当前Agent的chatModel的的上下文混乱。
    - 只有当 Event 的 RunPath “属于”当前 Agent 的 RunPath 时（ RunPathA 与 RunPathB 相同或者 RunPathA 是 RunPathB 的前缀），该 Event 才会参与构建 Agent 的 AgentInput。（过滤掉无关的AgentInput）
    - 通过 WithHistoryRewriter 可以自定义 Agent 从 History 中生成 AgentInput 的方式：
        ```go
        // github.com/cloudwego/eino/adk/flow.go
        type HistoryRewriter func(ctx context.Context, entries []*HistoryEntry) ([]Message, error)
        func WithHistoryRewriter(h HistoryRewriter) AgentOption
        ```
<!-- 列布局结束 -->



**Transfer的含义是将任务****移交运行****给子Agent（****SubAgents****），而不是委托或者分配（区别于ToolCall）。**Agent 运行时产生带有**包含 TransferAction 的 AgentEvent 后**，Eino ADK 会调用 Action 指定的 Agent。

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

TransferAction 可以使用 `NewTransferToAgentAction` 快速创建：

```go
import "github.com/cloudwego/eino/adk"
event := adk.NewTransferToAgentAction("dest agent name")
```



1. **显式调用（ToolCall）**：将 Agent 视为工具进行调用。适用于 Agent 运行仅需要明确清晰的参数而非完整运行上下文的场景，常结合 ChatModelAgent，作为工具运行后将结果返回给 ChatModel 继续处理。除此之外，ToolCall 同样支持调用符合工具接口构造的、不含 Agent 的普通工具。
<!-- 列布局开始 -->

![](/images/24724637-29b5-80af-a6a1-e96b72202555/image_29a24637-29b5-80d9-aedb-ec3e458a4e51.jpg)


---

example:

```go
// 将 Agent 转换为 Tool
func NewAgentTool(_ context.Context, agent Agent, options ...AgentToolOption) tool.BaseTool
```



<!-- 列布局结束 -->



### **Runner抽象与 Interrupted Action、Checkpoint、Resume**

**Runner: Eino ADK 中负责执行 Agent 的核心引擎**。

主要作用是**管理和控制 Agent 的整个生命周期****:**如处理多 Agent 协作，保存传递上下文等，interrupt、callback 等切面能力也均依赖 Runner 实现。

**任何 Agent 都应通过 Runner 来运行**。

```go
// github.com/cloudwego/eino/adk/runners.go
// 声明
type Runner struct {
	a               Agent
	enableStreaming bool
	store           compose.CheckPointStore
}

// 调用
runner := adk.NewRunner(ctx, adk.RunnerConfig{
		EnableStreaming: true, // you can disable streaming here
		Agent:           a,
		CheckPointStore: newInMemoryStore(),
})
```



```go
func (r *Runner) Run(ctx context.Context, messages []Message, opts ...AgentRunOption) *AsyncIterator[*AgentEvent]
    
// Query 是为了方便单次查询而提供的Run的语法糖
func (r *Runner) Query(ctx context.Context,query string, opts ...AgentRunOption) *AsyncIterator[*AgentEvent] {
	return r.Run(ctx, []Message{schema.UserMessage(query)}, opts...)
}
```



**Runner 提供运行时中断与恢复的功能:**

允许正在运行中的 Agent 主动中断并保存其当前状态，并在未来从中断点恢复执行。

使用场景：长时间等待、可暂停或需要外部输入（Human in the loop）等。多轮对话（ 多次的`runner.Query()` ）？

1. **Interrupted Action**：由 Agent 抛出`Interrupt Action` 的 `Event` 中断事件，主动通知Agent `Runner` 中断运行（拦截）。并允许携带额外信息供调用方阅读与使用。
1. **Checkpoint**：Agent `Runner` 拦截事件后，通过初始化时注册的 `CheckPointStore` 保存当前运行状态。Runner 在终止运行后会将当前运行状态（原始输入、对话历史等）以及 Agent 抛出的 InterruptInfo 以 CheckPointID 为 key 持久化到 CheckPointStore 中。
1. **Resume**：运行条件重新 ready 后，由 Agent `Runner` 从断点通过 `Resume` 方法携带恢复运行所需要的新信息，从断点处恢复运行。


**example**:

```go
// 1. 创建支持断点恢复的 Runner
runner := adk.NewRunner(ctx, adk.RunnerConfig{
    Agent:           complexAgent,
    CheckPointStore: memoryStore, // 内存状态存储
})

// 2. 开始执行
iter := runner.Query(ctx, "recommend a book to me", adk.WithCheckPointID("1"))
for {
    event, ok := iter.Next()
    if !ok {
       break
    }
    if event.Err != nil {
       log.Fatal(event.Err)
    }
    if event.Action != nil {
        // 3. 由 Agent 内部抛出 Interrupt 事件
        if event.Action.Interrupted != nil {
           ii, _ := json.MarshalIndent(event.Action.Interrupted.Data, "", "\t")
           fmt.Printf("action: interrupted\n")
           fmt.Printf("interrupt snapshot: %v", string(ii))
        }
    }
}

// 4. 从 stdin 接收用户输入
scanner := bufio.NewScanner(os.Stdin)
fmt.Print("\nyour input here: ")
scanner.Scan()
fmt.Println()
nInput := scanner.Text()

// 5. 携带用户输入信息，从断点恢复执行
iter, err := runner.Resume(ctx, "1", adk.WithToolOptions([]tool.Option{subagents.WithNewInput(nInput)}))
```



**序列化：**

**为了保存 interface 中数据的原本类型，Eino ADK 使用 gob（**[**https://pkg.go.dev/encoding/gob**](https://pkg.go.dev/encoding/gob)**）序列化运行状态**。

Eino 会**自动注册框架内置的类型，**在使用自定义类型时需要提前使用 gob.Register 或 **gob.RegisterName 注册类型**（更推荐后者，前者使用路径加类型名作为默认名字，因此类型的位置和名字均不能发生变更）。



**inMemoryStore：**

**compose.CheckPointStore interface的一个实现。**

```go
// **compose.CheckPointStore**
type CheckPointStore interface {
	Get(ctx context.Context, checkPointID string) ([]byte, bool, error)
	Set(ctx context.Context, checkPointID string, checkPoint []byte) error
}
```

```go
type inMemoryStore struct {
	mem map[string][]byte
}

func (i *inMemoryStore) Set(ctx context.Context, key string, value []byte) error {
	i.mem[key] = value
	return nil
}

func (i *inMemoryStore) Get(ctx context.Context, key string) ([]byte, bool, error) {
	v, ok := i.mem[key]
	return v, ok, nil
}

func newInMemoryStore() compose.CheckPointStore {
	return &inMemoryStore{
		mem: map[string][]byte{},
	}
}
```



**Resume:**

调用 Runner 的 Resume 接口，传入中断时的 CheckPointID 可以恢复运行：

```go
iter, err := runner.Resume(ctx, "1", adk.WithToolOptions([]tool.Option{subagents.WithNewInput(nInput)}))

// github.com/cloudwego/eino/adk/runner.go
func (r *Runner) Resume(ctx context.Context, checkPointID string, **opts ...AgentRunOption**) (*AsyncIterator[*AgentEvent], error)
```

恢复 Agent 运行需要发生中断的 Agent 实现了 ResumableAgent 接口， Runner 从 CheckPointerStore 读取运行状态并恢复运行，

其中** InterruptInfo 和上次运行配置的 EnableStreaming 会作为输入提供给 Agent**：

Resume如果向 Agent 传入新信息，**可以定义 AgentRunOption，并在调用 Runner.Resume 时传入**。

```go
// github.com/cloudwego/eino/adk/interface.go
type ResumableAgent interface {
    Agent

    Resume(ctx context.Context, info *ResumeInfo, **opts ...AgentRunOption**) *AsyncIterator[*AgentEvent]
}

// github.com/cloudwego/eino/adk/interrupt.go
type ResumeInfo struct {
		EnableStreaming bool
		*InterruptInfo
}
```



# 四、adk example

## 4.1 helloworld

### `helloworld` ChatModelAgent

7行代码：实现简单对话式ChatModelAgent

```go
model, err := ark.NewChatModel(...)
agent, err := adk.NewChatModelAgent(...）
runner := adk.NewRunner(...)

iter := runner.Query(ctx, "Hello, please introduce yourself. use chinese to answer")
for {
		event, ok := iter.Next()
		msg, err := event.Output.MessageOutput.GetMessage()
}
```



### `ChatModelAgent` 

```go
// 核心一行代码
runner := adk.NewRunner{
		...,
		CheckPointStore: newInMemoryStore(),  // map[string][]byte
	})
iter := runner.Query(ctx, "recommend a book to me", adk.WithCheckPointID("1"))
iter, err := runner.Resume(ctx, "1", adk.WithToolOptions([]tool.Option{subagents.WithNewInput(nInput)}))
```

12行代码：使用 `ChatModelAgent` 带interrupt中断和恢复、本地function tool。

```go
model, err := ark.NewChatModel(...)
bookSearchTool, err := utils.InferTool(..., func(ctx, input) (output, err) {...})
newAskForClarificationTool, err := utils.InferOptionableTool(...,func(..., opts) (output, err) {...}
agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{..., ToolsConfig: adk.ToolsConfig{...}}
runner := adk.NewRunner{
		...,
		CheckPointStore: newInMemoryStore(),  // map[string][]byte
	})
	
iter := runner.Query(ctx, "recommend a book to me", adk.WithCheckPointID("1"))
// 交互循环
for {
		event, ok := iter.Next()
		prints.Event(event)
	}
	

scanner := bufio.NewScanner(os.Stdin)
scanner.Scan()
nInput := scanner.Text()
iter, err := runner.Resume(ctx, "1", adk.WithToolOptions([]tool.Option{subagents.WithNewInput(nInput)}))
for {
		event, ok := iter.Next()
		prints.Event(event)
}
```



### `custom Agent`

7行代码：实现符合ADK定义的自定义Agent

```go
type MyAgent struct {}

func (m *MyAgent) Name() {...}
func (m *MyAgent) Description() {...}
func (m *MyAgent) Run(...) *adk.AsyncIterator[*adk.AgentEvent] {
	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()
	go func() {
		defer func() {
			e := recover()
			gen.Close()
		}()
		// agent run code
		gen.Send(&adk.AgentEvent{
			Output: &adk.AgentOutput{
				MessageOutput: &adk.MessageVariant{
					IsStreaming: false,
					Message: &schema.Message{
						Role:    schema.Assistant,
						Content: "hello world",
					},
					Role: schema.Assistant,
				},
			},
		})
	}()
}
```

## 4.2 workflow：Loop agent + Parallel agent + Sequential agent

```go
// 核心一行代码
a, err := adk.NewLoopAgent(ctx, &adk.LoopAgentConfig{... SubAgents:{[]adk.Agent{a1,a2}...}
a, err := adk.NewParallelAgent(ctx, &adk.LoopAgentConfig{... SubAgents:{[]adk.Agent{a1,a2，a3}...}
a, err := adk.NewSequentialAgent(ctx, &adk.LoopAgentConfig{... SubAgents:{[]adk.Agent{a1,a2}...}

ctx, endSpanFn := startSpanFn(ctx, "layered-supervisor", query)
endSpanFn(ctx, lastMessage)
```

`Loop` agent（循环agent）：14行代码，loop agent：1个main agent + 1个critique****agent， + cozeloop trace

```go
// cozeloop trace: eino-ext/callbacks/cozeloop   coze-dev/cozeloop-go
traceCloseFn, startSpanFn := trace.AppendCozeLoopCallbackIfConfigured(ctx)
defer traceCloseFn(ctx)
AppendCozeLoopCallbackIfConfigured() 

cm, err := ark.NewChatModel()
a1, err := adk.NewChatModelAgent()
a2, err := adk.NewChatModelAgent()
a, err := adk.NewLoopAgent(ctx, &adk.LoopAgentConfig{... SubAgents:{[]adk.Agent{a1,a2}...}

query := "briefly introduce what a multimodal embedding model is."
ctx, endSpanFn := startSpanFn(ctx, "layered-supervisor", query)
runner := adk.NewRunner()

iter := runner.Query(ctx, query)
for {
		event, ok := iter.Next()
		prints.Event(event)
}
endSpanFn(ctx, lastMessage)
```



`Parallel` agent（并行agent）：15行代码，Parallel agent：1个Stock数据收集 agent + 1个News数据收集****agent + 1个社交媒体数据收集****agent， + cozeloop trace

```go
// cozeloop trace: eino-ext/callbacks/cozeloop   coze-dev/cozeloop-go
traceCloseFn, startSpanFn := trace.AppendCozeLoopCallbackIfConfigured(ctx)
defer traceCloseFn(ctx)
AppendCozeLoopCallbackIfConfigured() 

cm, err := ark.NewChatModel()
a1, err := adk.NewChatModelAgent() // NewStockDataCollectionAgent
a2, err := adk.NewChatModelAgent() // NewNewsDataCollectionAgent
a3, err := adk.NewChatModelAgent() // NewSocialMediaInfoCollectionAgent
a, err := adk.NewParallelAgent(ctx, &adk.LoopAgentConfig{... SubAgents:{[]adk.Agent{a1,a2，a3}...}

query := "give me today's market research"
ctx, endSpanFn := startSpanFn(ctx, "layered-supervisor", query)
runner := adk.NewRunner()

iter := runner.Query(ctx, query)
for {
		event, ok := iter.Next()
		prints.Event(event)
}
endSpanFn(ctx, lastMessage)
```



`Sequential` (连续的)agent:

```go
// cozeloop trace: eino-ext/callbacks/cozeloop   coze-dev/cozeloop-go
traceCloseFn, startSpanFn := trace.AppendCozeLoopCallbackIfConfigured(ctx)
defer traceCloseFn(ctx)
AppendCozeLoopCallbackIfConfigured() 

cm, err := ark.NewChatModel()
a1, err := adk.NewChatModelAgent() // NewPlanAgent
a2, err := adk.NewChatModelAgent() // NewWriterAgent
a, err := adk.NewSequentialAgent(ctx, &adk.LoopAgentConfig{... SubAgents:{[]adk.Agent{a1,a2}...}

query := "give me today's market research"
ctx, endSpanFn := startSpanFn(ctx, "layered-supervisor", query)
runner := adk.NewRunner()

iter := runner.Query(ctx, query)
for {
		event, ok := iter.Next()
		prints.Event(event)
}
endSpanFn(ctx, lastMessage)
```



## 4.3 跨agent传递

###  `session`：跨agent传递 data and state（状态）

```go
// 核心一行代码
adk.AddSessionValue(ctx, "user-name", in.Name)  // a1
userName, _ := adk.GetSessionValue(ctx, "user-name") // a2
```

9行代码：AddSessionValue、GetSessionValue

```go
adk.AddSessionValue(ctx, "user-name", in.Name)  // a1
userName, _ := adk.GetSessionValue(ctx, "user-name") // a2

toolA, err := utils.InferTool("tool_a", "set user name", toolAFn)
toolB, err := utils.InferTool("tool_b", "set user age", toolBFn)

a, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
	...
	ToolsConfig{toolA,toolB},
	Model: model.NewChatModel()
}

r := adk.NewRunner(Agent: a)

iter := r.Query(ctx, "my name is Alice, my age is 18")
for {
		event, ok := iter.Next()
		prints.Event(event)
}
```



### `transfer移交运行`

```go
// 核心一行代码
a, err := adk.SetSubAgents(ctx, routerAgent, []adk.Agent{chatAgent, weatherAgent})
```

12行代码：通过SetSubAgents的的transfer_to_agent 实现控制权的动态选择与转移。

Agent 职责单一 模块化，可独立开发测试，子 Agent 专注各自能力；



```go
weatherTool, err := utils.InferTool(...)
a1, err := adk.NewChatModelAgent(... Tools:weatherTool) // weatherAgent
a2, err := adk.NewChatModelAgent() // chatAgent
a3, err := adk.NewChatModelAgent() // routerAgent
a, err := adk.SetSubAgents(routerAgent, []adk.Agent{chatAgent, weatherAgent}) // SetSubAgents会在 RouterAgent 中注入 transfer_to_agent

runner := adk.NewRunner(a) 
iter := runner.Query(ctx, "What's the weather in Beijing?") // transfer(转移)到 WeatherAgent
for {
		event, ok := iter.Next()
		prints.Event(event)
}

iter = runner.Query(ctx, "Book me a flight from New York to London tomorrow.") // 无匹配 Agent，RouterAgent 直接回复无法处理
for {
		event, ok := iter.Next()
		prints.Event(event)
}
```



## 4.4 加载.env的方法

### 方案1：使用.env文件配置环境变量

1. vscode 安装 `mikestead.dotenv` 扩展：支持.env .env.local .env.example等常见文件名的语法高亮。但不支持自动加载环境变量。
1. 项目根目录创建.env文件，务必将.env添加到.gitignore(否则ak/sk泄露到gitlab/github)。在.env中配置：
    1. 注意：不带双引号，不带export开头。
    1. 终端及其子进程要生效.env： `export $(grep -v '^#' .env | xargs)` 或 .env每行都要加export开头再source
        1. 直接`source .env` ；测试当前终端能生效：`echo $ARK_API_KEY` ，**但终端运行子进程时仍然读取不到**。
    ```shell
    # ark model: https://console.volcengine.com/ark
    # 必填
    # 火山云方舟 ChatModel 的 Endpoint ID
    ARK_CHAT_MODEL=""
    # 火山云方舟 向量化模型的 Endpoint ID
    ARK_EMBEDDING_MODEL=""
    # 火山云方舟的 API Key
    ARK_API_KEY=""
    ARK_BASE_URL="https://ark.cn-beijing.volces.com/api/v3/"
    # apmplus: https://console.volcengine.com/apmplus-server
    # 下面必填环境变量如果为空，则不开启 apmplus callback
    # APMPlus 的 App Key，必填
    APMPLUS_APP_KEY=""         
    # APMPlus 的 Region，选填，不填写时，默认是 cn-beijing
    APMPLUS_REGION=""
    # langfuse: https://cloud.langfuse.com/
    # 下面两个环境变量如果为空，则不开启 langfuse callback
    # Langfuse Project 的 Public Key 
    LANGFUSE_PUBLIC_KEY=""
    # Langfuse Project 的 Secret Key。 注意，Secret Key 仅可在被创建时查看一次
    LANGFUSE_SECRET_KEY=""
    # Redis Server 的地址，不填写时，默认是 localhost:6379
    REDIS_ADDR=
    OPENAI_API_KEY=""
    OPENAI_MODEL="gpt-4o-mini"
    OPENAI_BASE_URL="https://api.openai.com/v1"
    OPENAI_BY_AZURE=false
    ```
1. 从调试设置里创建 .vscode/launch.json，设置加载.env。
    ```json
    {
      "version": "0.2.0",
      "configurations": [
        {
    	    // 配置名称，显示在下拉菜单中
          "name": "Debug helloworld",
          // 调试器类型
          "type": "go",
          // 请求类型：launch（启动程序）或 attach（附加到运行中进程）
          "request": "launch",
          // Go 特定的模式：debug, test, exec 等
          "mode": "auto",
          // 要调试的程序路径
          "program": "${workspaceFolder}/adk/helloworld",
          // 从 .env 文件加载环境变量
          "envFile": "${workspaceFolder}/.env",
          // 控制台类型
          "console": "integratedTerminal",
          // 是否显示详细日志
          "showLog": false
        }
      ]
    }
    ```
1. command + shift + D启动调试，默认会加载.env。F5继续调试、F11单步进入；


go test加载.env：

.vscode/settings.json

```go
{
    "go.testEnvFile": "${workspaceFolder}/.env"
}
```

然后直接vscode中的 run test



**上述为调试/运行场景，如果需要sudo运行：**

**Terminal-Run Task**：使用VSCode的**Tasks + 环境配置**

.vscode/tasks.json：在 VSCode 内用 **Terminal-Run Task**即可执行。**等价于在终端里执行其中的command。**

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Run hello.go with sudo",
      "type": "shell",
      "command": "export $(grep -v '^#' .env | xargs) && sudo -E go run hello.go"
      "options": {
        "env": {
          "GO_ENV": "development"
        }
      },
      "problemMatcher": []
    }
  ]
}
```

### 方法2：direnv工具 + sudo -E

1. 安装 `direnv`
    ```shell
    sudo apt install direnv    # Ubuntu/Debian
    # 或
    brew install direnv        # macOS
    ```
1. 在项目根目录创建 `.envrc`
    ```shell
    export $(grep -v '^#' .env | xargs)
    ```
1. 启用 `direnv`:` direnv allow`。每次进入项目目录时，`.env` 中的变量都会被自动加载到当前 shell。
1. 在 VSCode 终端中运行：`**-E**`** 表示保留当前环境变量（包括**`**.env**`** 中的变量）。(否则，切换到sudo后前面source .env的环境变量丢失了)**
    ```shell
    sudo -E go run h.go
    ```
### 方法3：godotenv package

```shell
import "github.com/joho/godotenv"
err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or failed to load")
	}

	// 测试读取变量
	apiKey := os.Getenv("API_KEY")
	fmt.Println("API_KEY:", apiKey)
```



## 附：老版本 Eino React Agent（基于compose.Graph）

> 王德政: 
eino flow/下的react agent 和adk/ 下的 react agent基本没区别，都是 react 模式；
adk 下的 ChatModelAgent 是符合 adk.Agent 接口的，接口更易用一些；
如果要使用 adk 相关能力建议用 adk 目录下的， 如果单独使用也建议用 adk 目录下的；



备注：不推荐。推荐使用后面新上线的adk/chatModelAgent，有更规范的agent定义的interface，其封装了adk/react.go的ReAct能力。

```go
// 通过 flow/agent/react 包提供完整实现
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


# 五、a2a

Agent2Agent协议是Google提出的一种开放标准，旨在实现AI Agent之间的无缝通信与协作。

<!-- 列布局开始 -->

Eino ADK提供了双向封装能力：

- 将 **Local Agent** 发布为 **A2A Server**
- 将A2A Client连接的 **Remote A2A Server** 转换为 **Local Agent**
example：还在 alpha 版本中。[https://github.com/cloudwego/eino-ext/tree/a2a/v0.0.1-alpha.4/a2a/extension/eino/examples](https://github.com/cloudwego/eino-ext/tree/a2a/v0.0.1-alpha.4/a2a/extension/eino/examples)




---

![](/images/24724637-29b5-80af-a6a1-e96b72202555/image_29d24637-29b5-80ff-9aed-d70fa7ff3e66.jpg)

![](/images/24724637-29b5-80af-a6a1-e96b72202555/image_29d24637-29b5-808c-94e1-c30bc5db3799.jpg)



<!-- 列布局结束 -->





<!-- 列布局开始 -->

![](/images/24724637-29b5-80af-a6a1-e96b72202555/image_29d24637-29b5-804b-8760-eb29c6a38f00.jpg)


---

![](/images/24724637-29b5-80af-a6a1-e96b72202555/image_29d24637-29b5-8032-8db7-ed4aa3d18787.jpg)



<!-- 列布局结束 -->



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



引用：

- [Eino ADK：一文搞定 AI Agent 核心设计模式，从 0 到 1 搭建智能体系统](https://mp.weixin.qq.com/s/p_QqDN6m2anHAE97P2Q2bw?forceh5=1)
- [如何构建 MultiAgent——Eino adk 与 a2a 实践 - 王德政](https://www.bilibili.com/video/BV1qixrzFEWo/?spm_id_from=333.1391.0.0)
