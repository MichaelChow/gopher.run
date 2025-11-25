---
title: "5.8 adk进阶"
date: 2025-11-25T13:43:00Z
draft: false
weight: 5008
---

# 5.8 adk进阶

# 一、Plan-Execute MultiAgent范式（结构化解决问题）

## 源码解读

ADK 提供的基于「规划-执行-反思」范式的 Multi-Agent 协作模式（参考论文 **Plan-and-Solve Prompting**），旨在解决复杂任务的分步拆解、执行与动态调整问题。

通过** Planner（规划器）、Executor（执行器）和 Replanner（重规划器）** 三个核心Agent的协同工作，实现任务的结构化规划、工具调用执行、进度评估与动态replanning，最终达成用户目标。

其中：

- **规划者（Planner）**：根据用户目标，生成一个包含**详细步骤且结构化的初始任务计划**
- **执行者（Executor）**：执行当前计划中的首个步骤，调用外部工具完成具体任务。基于 `ChatModelAgent` 实现，配置工具集（如搜索、计算、数据库访问等）
    - 从 Session 中获取当前 `Plan` 和已执行步骤
    - 提取计划中的第一个未执行步骤作为目标
    - 调用工具执行该步骤，将结果存储于 Session
- **反思者（Replanner）**：评估执行进度，决定是修正计划继续交由 Executor 运行，或是结束任务


**实现方式：**组合****`SequentialAgent` 和 `LoopAgent` 

- 外层 `SequentialAgent`：先执行 `Planner` 生成初始计划，再进入执行-重规划循环
- 内层 `LoopAgent`：循环执行 `Executor` 和 `Replanner`，直至任务完成或达到最大迭代次数
<!-- 列布局开始 -->

![](/images/2b624637-29b5-8095-9c5f-d7dcc209a2be/image_29924637-29b5-8010-a8cd-c85c866075e0.jpg)

![](/images/2b624637-29b5-8095-9c5f-d7dcc209a2be/image_29d24637-29b5-8001-90c3-dbcd89419223.jpg)


---

![](/images/2b624637-29b5-8095-9c5f-d7dcc209a2be/image_29d24637-29b5-8093-8e84-d565cf6413cc.jpg)

 




---

example:

```go
import github.com/cloudwego/eino/adk/prebuilt/planexecute

// Plan-Execute 模式的科研助手
researchAssistant := planexecute.New(ctx, &planexecute.Config{
    Planner: adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
        Name: "research_planner",
        Instruction: "制定详细的研究计划，包括文献调研、数据收集、分析方法等",
        Model: gpt4Model,
    }),
    Executor: adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
        Name: "research_executor",
        ToolsConfig: adk.ToolsConfig{
            Tools: []tool.BaseTool{
                scholarSearchTool,
                dataAnalysisTool,
                citationTool,
            },
        },
    }),
    Replanner: replannerAgent,
})
```

<!-- 列布局结束 -->

Plan-Execute 模式有如下特点：

- **明确的分层架构**：通过将任务拆解为规划、执行和反思重规划三个阶段，形成层次分明的认知流程，体现了** “先思考再行动，再根据反馈调整” 的闭环认知策略，在各类场景中都能达到较好的效果**。
- **动态迭代优化**：Replanner 根据执行结果和当前进度，实时判断任务是否完成或需调整计划，支持动态重规划。该机制有效**解决了传统单次规划难以应对环境变化和任务不确定性的瓶颈**，提升了系统的鲁棒性和灵活性。
- **职责分明且松耦合**：Plan-Execute 模式由多个智能体协同工作，支持独立开发、测试和替换。模块化设计方便扩展和维护，符合工程最佳实践。
- **具备良好扩展性**：不依赖特定的语言模型、工具或 Agent，方便集成多样化外部资源，满足不同应用场景需求。


非常适合**需要多步骤推理、动态调整和工具集成的复杂任务场景；**

可能的应用场景有：

- **复杂研究分析**：通过规划分解研究问题，执行多轮数据检索与计算，动态调整研究方向和假设，提升分析深度和准确性。
- **自动化工作流管理**：将复杂业务流程拆解为结构化步骤，结合多种工具（如数据库查询、API 调用、计算引擎）逐步执行，并根据执行结果动态优化流程。
- **多步骤问题解决**：适用于需要分步推理和多工具协作的场景，如法律咨询、技术诊断、策略制定等，确保每一步执行都有反馈和调整。
- **智能助理任务执行**：支持智能助理根据用户目标规划任务步骤，调用外部工具完成具体操作，并根据重规划思考结合用户反馈调整后续计划，提升任务完成的完整性和准确性。


**源码解读：**

```go
package planexecute
// Plan 表示一个包含一系列可执行步骤的执行计划。
// 它支持 JSON 序列化和反序列化，同时提供对第一个步骤的访问。

type Plan interface {
	FirstStep() string // 返回计划中要执行的第一个步骤。
	
	json.Marshaler // 将 Plan 序列化为 JSON，用于提示模板
	json.Unmarshaler // 将 JSON 内容反序列化为 Plan，用于处理chatModel或工具调用的输出
}

// NewPlan is a function type that creates a new Plan instance.
type NewPlan func(ctx context.Context) Plan
```



**Planner：**规划者，根据目标生成plan

```go
func NewPlanner(_ context.Context, cfg *PlannerConfig) (adk.Agent, error) {
	var chatModel model.BaseChatModel
	var toolCall bool
	if cfg.ChatModelWithFormattedOutput != nil {
		chatModel = cfg.ChatModelWithFormattedOutput
	} else {
		toolCall = true
		toolInfo := cfg.ToolInfo
		if toolInfo == nil {
			toolInfo = &PlanToolInfo
		}

		var err error
		chatModel, err = cfg.ToolCallingChatModel.WithTools([]*schema.ToolInfo{toolInfo})
		if err != nil {
			return nil, err
		}
	}

	inputFn := cfg.GenInputFn  
	if inputFn == nil {
		inputFn = defaultGenPlannerInputFn  // 如果没传，使用默认的 PlannerPrompt
	}

	planParser := cfg.NewPlan
	if planParser == nil {
		planParser = defaultNewPlan
	}

	return &planner{
		toolCall:   toolCall,
		chatModel:  chatModel,
		genInputFn: inputFn,
		newPlan:    planParser,
	}, nil
}
```



```go
type PlannerConfig struct {
	// 预配置为以 Plan格式 输出的模型
	ChatModelWithFormattedOutput model.BaseChatModel
	// 当提供 ToolInfo 时，它将使用工具调用来生成 Plan 结构。Model二选一。
	ToolCallingChatModel model.ToolCallingChatModel

	// ToolInfo 定义使用工具调用时 Plan 结构的模式。默认值PlanToolInfo。
	ToolInfo *schema.ToolInfo

	// GenInputFn 是**生成规划器输入消息的函数**。默认值defaultGenPlannerInputFn。
	GenInputFn GenPlannerModelInputFn

	// 创建一个新的Plan实例，用于反序列化 模型生成的JSON输出。默认值defaultNewPlan。
	NewPlan NewPlan
}
```



默认值PlanToolInfo:

```go
// PlanToolInfo 定义了可与 ToolCallingChatModel 一起使用的 Plan tool schema，用于指导模型生成一个包含有序步骤的结构化的已排序的计划。
PlanToolInfo = schema.ToolInfo{
		Name: "Plan",
		Desc: "生成一个按顺序执行的步骤列表计划。每个步骤应该清晰、可执行，并按逻辑顺序排列。输出将用于指导执行过程。",
		ParamsOneOf: schema.NewParamsOneOfByParams(
			map[string]*schema.ParameterInfo{
				"steps": {
					Type:     schema.Array,
					ElemInfo: &schema.ParameterInfo{Type: schema.String},
					Desc:     "要遵循的不同步骤，应按排序顺序排列",
					Required: true,
				},
			},
		),
	}
```



**GenInputFn（**PlannerPrompt**）**:

```go
// GenPlannerModelInputFn is a function type that generates input messages for the planner.
type GenPlannerModelInputFn func(ctx context.Context, userInput []adk.Message) ([]adk.Message, error)

func defaultGenPlannerInputFn(ctx context.Context, userInput []adk.Message) ([]adk.Message, error) {
	msgs, err := PlannerPrompt.Format(ctx, map[string]any{
		"input": userInput,
	})
	return msgs, nil
}


// PlannerPrompt 是规划器的提示模板。它为规划器提供上下文和指导，说明如何生成 Plan。
PlannerPrompt = prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是一个专业的规划智能体。给定一个目标，创建一个全面的分步计划来实现该目标。

## 你的任务
分析目标并生成一个战略计划，将目标分解为可管理、可执行的步骤。

## 规划要求
你计划中的每个步骤必须：
- **具体且可执行**：清晰的指令，可以无歧义地执行
- **自包含**：包含所有必要的上下文、参数和要求
- **独立可执行**：可以由智能体执行，不依赖其他步骤
- **逻辑有序**：按最优顺序排列以实现高效执行
- **目标导向**：直接有助于实现主要目标

## 规划指南
- 消除冗余或不必要的步骤
- 为每个步骤包含相关的约束、参数和成功标准
- 确保最后一步产生完整的答案或可交付成果
- 预见潜在挑战并包含缓解策略
- 构建步骤使其在逻辑上相互建立
- 提供足够的细节以确保成功执行

## 质量标准
- 计划完整性：是否涵盖了目标的所有方面？
- 步骤清晰度：每个步骤是否可以被独立理解和执行？
- 逻辑流程：步骤是否遵循合理的进展？
- 效率：这是实现目标最直接的路径吗？
- 适应性：计划能否处理意外结果或变化？`),
		schema.MessagesPlaceholder("input", false),
	)
```



**defaultPlan (实现Plan)：**

```go
type defaultPlan struct {
	// 步骤包含有序的行动列表。每一步都应清晰、可执行，并按逻辑顺序排列。
	Steps []string `json:"steps"`
}

func (p *defaultPlan) FirstStep() string { // 返回first step或""
	if len(p.Steps) == 0 {
		return ""
	}
	return p.Steps[0]
}

func (p *defaultPlan) MarshalJSON() ([]byte, error) {
	type planTyp defaultPlan
	return sonic.Marshal((*planTyp)(p))
}

func (p *defaultPlan) UnmarshalJSON(bytes []byte) error {
	type planTyp defaultPlan
	return sonic.Unmarshal(bytes, (*planTyp)(p))
}

// JSON Schema:
//
//	{
//	  "type": "object",
//	  "properties": {
//	    "steps": {
//	      "type": "array",
//	      "items": {
//	        "type": "string"
//	      },
//	      "description": "Ordered list of actions to be taken. Each step should be clear, actionable, and arranged in a logical sequence."
//	    }
//	  },
//	  "required": ["steps"]
//	}
```

自定义实现Plan：略。一般不需要？



**2. Executor（执行者）**: 执行plan中的first step

从 Session 获取计划、用户输入、已执行步骤，执行结果存入 ExecutedStepSessionKey

```go
// NewExecutor creates a new executor agent.
func NewExecutor(ctx context.Context, cfg *ExecutorConfig) (adk.Agent, error) {

	genInputFn := cfg.GenInputFn
	if genInputFn == nil {
		genInputFn = defaultGenExecutorInputFn
	}
	genInput := func(ctx context.Context, instruction string, _ *adk.AgentInput) ([]adk.Message, error) {

		plan, ok := adk.GetSessionValue(ctx, PlanSessionKey)
		if !ok {
			panic("impossible: plan not found")
		}
		plan_ := plan.(Plan)

		userInput, ok := adk.GetSessionValue(ctx, UserInputSessionKey)
		if !ok {
			panic("impossible: user input not found")
		}
		userInput_ := userInput.([]adk.Message)

		var executedSteps_ []ExecutedStep
		executedStep, ok := adk.GetSessionValue(ctx, ExecutedStepsSessionKey)
		if ok {
			executedSteps_ = executedStep.([]ExecutedStep)
		}

		in := &ExecutionContext{
			UserInput:     userInput_,
			Plan:          plan_,
			ExecutedSteps: executedSteps_,
		}

		msgs, err := genInputFn(ctx, in)
		if err != nil {
			return nil, err
		}

		return msgs, nil
	}

	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:          "Executor",
		Description:   "an executor agent",
		Model:         cfg.Model,
		ToolsConfig:   cfg.ToolsConfig,
		GenModelInput: genInput,
		MaxIterations: cfg.MaxIterations,
		OutputKey:     ExecutedStepSessionKey,
	})
	if err != nil {
		return nil, err
	}

	return agent, nil
}
```



```go
// ExecutorConfig provides configuration options for creating an executor agent.
type ExecutorConfig struct {
	// Model is the chat model used by the executor.
	Model model.ToolCallingChatModel

	// ToolsConfig specifies the tools available to the executor.
	ToolsConfig adk.ToolsConfig

	// MaxIterations defines the upper limit of ChatModel generation cycles.
	// The agent will terminate with an error if this limit is exceeded.
	// Optional. Defaults to 20.
	MaxIterations int

	// GenInputFn generates the input messages for the Executor.
	// Optional. If not provided, defaultGenExecutorInputFn will be used.
	GenInputFn GenModelInputFn
}
```



```go
// ExecutorPrompt 是执行器的提示模板。
// 它为执行器提供上下文和指导，说明如何执行任务。
ExecutorPrompt = prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是一个勤奋且细致的执行器智能体。遵循给定的计划，仔细且彻底地执行你的任务。`),
		schema.UserMessage(`## 目标
{input}
## 给定以下计划：
{plan}
## 已完成的步骤和结果
{executed_steps}
## 你的任务是执行第一步，即：
{step}`))

func defaultGenExecutorInputFn(ctx context.Context, in *ExecutionContext) ([]adk.Message, error) {

	planContent, err := in.Plan.MarshalJSON()

	userMsgs, err := ExecutorPrompt.Format(ctx, map[string]any{
		"input":          formatInput(in.UserInput),
		"plan":           string(planContent),
		"executed_steps": formatExecutedSteps(in.ExecutedSteps),
		"step":           in.Plan.FirstStep(),
	})

	return userMsgs, nil
}
```



**3. Replanner（重新规划者）：**评估进度，决定 继续执行 或 结束任务

- Plan：更新Session中的计划（继续执行）
- Respond：生成最终响应（结束任务），发送 BreakLoopAction 退出循环
```go
func NewReplanner(_ context.Context, cfg *ReplannerConfig) (adk.Agent, error) {
	planTool := cfg.PlanTool
	if planTool == nil {
		planTool = &PlanToolInfo
	}

	respondTool := cfg.RespondTool
	if respondTool == nil {
		respondTool = &RespondToolInfo
	}

	chatModel, err := cfg.ChatModel.WithTools([]*schema.ToolInfo{planTool, respondTool})
	if err != nil {
		return nil, err
	}

	planParser := cfg.NewPlan
	if planParser == nil {
		planParser = defaultNewPlan
	}

	return &replanner{
		chatModel:   chatModel,
		planTool:    planTool,
		respondTool: respondTool,
		genInputFn:  cfg.GenInputFn,
		newPlan:     planParser,
	}, nil
}
```



```go
type ReplannerConfig struct {
	// ChatModel，配置 PlanTool 和 RespondTool 来生成更新的计划或响应
	ChatModel model.ToolCallingChatModel

	// 定义Plan tool的schema。默认值PlanToolInfo（同planer）
	PlanTool *schema.ToolInfo

	// 定义RespondTool的schema。默认值RespondToolInfo
	RespondTool *schema.ToolInfo

	// 生成 Replanner的输入消息。默认值buildGenReplannerInputFn
	GenInputFn GenModelInputFn

	// 创建一个新的 Plan 实例，返回的 Plan 将用于反序列化来自 PlanTool 的模型生成的 JSON 输出。默认值defaultNewPlan
	NewPlan NewPlan
}
```



```go
var (
	// 定义RespondTool的默认schema，指示模型生成对用户的直接响应。
	RespondToolInfo = schema.ToolInfo{
		Name: "Respond",
		Desc: "生成对用户的直接响应。当你拥有提供最终答案所需的所有信息时，使用此工具。",
		ParamsOneOf: schema.NewParamsOneOfByParams(
			map[string]*schema.ParameterInfo{
				"response": {
					Type:     schema.String,
					Desc:     "提供给用户的完整响应",
					Required: true,
				},
			},
		),
	}
```



```go
	// ReplannerPrompt 是重新规划器的提示模板。
	// 它为重新规划器提供上下文和指导，说明如何重新生成 Plan。
	ReplannerPrompt = prompt.FromMessages(schema.FString,
		schema.SystemMessage(
			`你将审查实现目标的进展。分析当前状态并确定最优的下一步行动。

## 你的任务
基于上述进展，你必须选择恰好一个行动：

### 选项 1：完成（如果目标已完全实现）
调用 '{respond_tool}'，包含：
- 全面的最终答案
- 清晰总结目标如何达成的结论
- 执行过程中的关键见解

### 选项 2：继续（如果需要更多工作）
调用 '{plan_tool}'，包含一个修订的计划，该计划：
- 仅包含剩余步骤（排除已完成的步骤）
- 融入从已执行步骤中吸取的经验教训
- 解决发现的任何差距或问题
- 保持逻辑步骤顺序

## 规划要求
你计划中的每个步骤必须：
- **具体且可执行**：清晰的指令，可以无歧义地执行
- **自包含**：包含所有必要的上下文、参数和要求
- **独立可执行**：可以由智能体执行，不依赖其他步骤
- **逻辑有序**：按最优顺序排列以实现高效执行
- **目标导向**：直接有助于实现主要目标

## 规划指南
- 消除冗余或不必要的步骤
- 根据新信息调整策略
- 为每个步骤包含相关的约束、参数和成功标准

## 决策标准
- 原始目标是否已完全满足？
- 是否还有剩余的要求或子目标？
- 结果是否表明需要调整策略？
- 还需要哪些具体行动？`),
		schema.UserMessage(`## 目标
{input}

## 原始计划
{plan}

## 已完成的步骤和结果
{executed_steps}`),
	)
)
```



1. **整体流程**: SequentialAgent里面嵌套一个LoopAgent
```go
func New(ctx context.Context, cfg *Config) (adk.Agent, error) {
	loop, err := adk.NewLoopAgent(ctx, &adk.LoopAgentConfig{
		Name:          "execute_replan",
		SubAgents:     []adk.Agent{cfg.Executor, cfg.Replanner},
		MaxIterations: maxIterations,
	})
	
	return adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{
		Name:      "plan_execute_replan",
		SubAgents: []adk.Agent{cfg.Planner, loop},
	})
}
```



## example：`plan-execute-replan`

`plan-execute-replan` agent：

```go
planAgent := planexecute.NewPlanner(cm) // eino/adk/prebuilt/planexecute
return planexecute.NewPlanner(ctx, &planexecute.PlannerConfig{
		ToolCallingChatModel: model.NewChatModel(),
	})

executeAgent := planexecute.NewExecutor(cm)
return planexecute.NewExecutor(ctx, &planexecute.ExecutorConfig{
		Model: model.NewChatModel(),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: travelTools,
			},
		},

		GenInputFn: xxx

planexecute.NewReplanner(cm)
return planexecute.NewReplanner(ctx, &planexecute.ReplannerConfig{
		ChatModel: model.NewChatModel(),
	})


entryAgent, err := planexecute.New(ctx, &planexecute.Config{
		Planner:       planAgent,
		Executor:      executeAgent,
		Replanner:     replanAgent,
		MaxIterations: 20,
	})
	
r := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: entryAgent,
})
	
query := `Plan a 3-day trip to Beijing in Next Month. I need flights from New York, hotel recommendations, and must-see attractions.
Today is 2025-09-09.`

iter := r.Query(ctx, query)
for {
		event, ok := iter.Next()
		prints.Event(event)
}
```

第二层-ToolKit：

```go
// GetAllTravelTools returns all travel-related tools
func GetAllTravelTools(ctx context.Context) ([]tool.BaseTool, error) {
	weatherTool, err := NewWeatherTool(ctx)
	flightTool, err := NewFlightSearchTool(ctx)
	hotelTool, err := NewHotelSearchTool(ctx)
	attractionTool, err := NewAttractionSearchTool(ctx)
	askForClarificationTool := NewAskForlarificationTool()

	return []tool.BaseTool{weatherTool, flightTool, hotelTool, attractionTool, askForClarificationTool}, nil
}

travelTools, err := tools.GetAllTravelTools(ctx)
planexecute.NewExecutor(ctx, &planexecute.ExecutorConfig{
		Model: model.NewChatModel(),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: travelTools,
			},
		},
```

第二层-GenInputFn：输入生成函数，将 ***planexecute.ExecutionContext上下文**转换为 LLM 可理解的消息

```go
GenInputFn: func(ctx context.Context, in *planexecute.ExecutionContext) ([]adk.Message, error) {
			planContent, err_ := in.Plan.MarshalJSON()

			firstStep := in.Plan.FirstStep()

			msgs, err_ := executorPrompt.Format(ctx, map[string]any{
				"input":          formatInput(in.UserInput),
				"plan":           string(planContent),
				"executed_steps": formatExecutedSteps(in.ExecutedSteps),
				"step":           firstStep,
			})
			
			return msgs, nil
		},
```



## `example：integration-excel-agent`

文档：[https://mp.weixin.qq.com/s/787AJLf2czPn4L-FAnB9zA](https://mp.weixin.qq.com/s/787AJLf2czPn4L-FAnB9zA)

跑起来：

1. 根据readme.md，配置ARK_VISION_MODEL、ARK_VISION_API_KEY、EXCEL_AGENT_PYTHON_EXECUTABLE_PATH=python3 
1. 在.env同级目录安装venv虚拟环境，并安装好readme的依赖
1. 手动加载env和source venv/bin/activate，最后以sudo身份读取env并启动。codeagent写出的python代码，找不到questions.csv?


```go
export $(grep -v '^#' .env | xargs) && source venv/bin/activate &&  sudo -E go run '/Users/xxx/eino-examples/adk/multiagent/integration-excel-agent'
```



核心代码-第一层：

```go
p, err := planner.NewPlanner(ctx, operator)
e, err := executor.NewExecutor(ctx, operator)
rp, err := replanner.NewReplanner(ctx, operator)
planExecuteAgent, err := planexecute.New(ctx, &planexecute.Config{
		Planner:       p,
		Executor:      e,
		Replanner:     rp,
		MaxIterations: 20,
	})

reportAgent, err := report.NewReportAgent(ctx, operator)

agent, err := adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{
		Name:        "SequentialAgent",
		Description: "sequential agent",
		SubAgents: []adk.Agent{
			planExecuteAgent, reportAgent,
		},
	})

query := schema.UserMessage("请帮我将 questions.csv 表格中的第一列提取到一个新的 csv 中")
runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: true,
	})

ctx = params.InitContextParams(ctx)
params.AppendContextParams(ctx, map[string]interface{}{  // 传参？
		params.FilePathSessionKey:            inputFileDir,
		params.WorkDirSessionKey:             workdir,
		params.UserAllPreviewFilesSessionKey: utils.ToJSONString(previews),
		params.TaskIDKey:                     uuid,
	})
iter := runner.Run(ctx, []*schema.Message{query})
prints.Event(event)

// Operator → Tool → LLM
operator := &LocalOperator{}  // commandline.Operator 接口的实现，提供底层文件与命令执行能力。**不直接作为tool，用于被多个Tool复用。**
```



核心代码-第二层-NewPlanner：

plannerPromptTemplate：SystemPrompt + UserPrompt

```go
var (
	plannerPromptTemplate = prompt.FromMessages(schema.Jinja2,
		schema.SystemMessage(`你是一位专门从事 Excel 数据处理任务的专家规划师。你的目标是理解用户需求并将其分解为清晰的、逐步执行的计划。

**1. 理解目标：**
- 仔细分析用户的请求，确定最终目标。
- 识别输入数据（Excel 文件）和期望的输出格式。

**2. 交付物：**
- 最终输出应该是一个表示计划的 JSON 对象，包含步骤列表。
- 每个步骤必须是对执行该步骤的代理的清晰、简洁的指令。

**3. 计划分解原则：**
- **粒度：** 将任务分解为尽可能小的逻辑步骤。例如，不要使用"处理数据"，而应使用"读取 Excel 文件"、"过滤掉缺失值的行"、"计算 'Sales' 列的平均值"等。
- **顺序：** 步骤应按正确的执行顺序排列。
- **清晰度：** 每个步骤应该明确无误，易于执行该步骤的代理理解。

**4. 输出格式（少样本示例）：**
以下是一个良好计划的示例：
用户请求："请计算附件 'sales_data.xlsx' 文件中每个产品类别的平均销售额，并生成报告。"
{
  "steps": [
    {
      "instruction": "将 'sales_data.xlsx' 文件读取到 pandas DataFrame 中。"
    },
    {
      "instruction": "按 'Product Category' 对 DataFrame 进行分组，并计算每个组的 'Sales' 列的平均值。"
    },
    {
      "instruction": "总结每个产品类别的平均销售额，并以表格形式呈现结果。"
    }
  ]
}

**5. 限制条件：**
- 不要在计划中直接生成代码。
- 确保计划是逻辑合理且可实现的。
- 最后一步应该始终是生成报告或提供最终结果。
`),
		schema.UserMessage(`
用户查询：{{ user_query }}
当前时间：{{ current_time }}
文件预览（如果文件具有 xlsx 扩展名，预览将提供前 20 行的具体内容，否则仅提供文件路径）：
{{ file_preview }}
`),
	)
)
```

```go
sc, err := generic.PlanToolInfo.ToJSONSchema()

cm := utils.NewChatModel(ctx,
		utils.WithMaxTokens(4096),  // 计划很简洁，同时也节省token
		utils.WithTemperature(0),  // 几乎确定性，相同输入产生相同输出
		utils.WithTopP(0),  // 只考虑最可能的 token
		utils.WithDisableThinking(true),  // 规划任务需要直接输出，不需要中间推理。节省token，提高效率。
		utils.WithResponseFormatJsonSchema(&openai.ChatCompletionResponseFormatJSONSchema{
			Name:        generic.PlanToolInfo.Name,
			Description: generic.PlanToolInfo.Desc,
			JSONSchema:  sc,  // 上述json schema，确保输出符合JSON Schema
			Strict:      true,
		}),
	)

a, err := planexecute.NewPlanner(ctx, &planexecute.PlannerConfig{
    ChatModelWithFormattedOutput: cm,  // 确保输出符合 JSON Schema
    **GenInputFn**: func(ctx context.Context, userInput []adk.Message) ([]adk.Message, error) {  // prompt，使用上述prompt模板
				pf, _ := params.GetTypedContextParams[string](ctx, params.UserAllPreviewFilesSessionKey)
				msgs, err := plannerPromptTemplate.Format(ctx, map[string]any{ 
					"user_query":   utils.FormatInput(userInput),
					"current_time": utils.GetCurrentTime(),
					"file_preview": pf,
				})
				return msgs, nil
		},  // 将用户输入转换为模型输入
    NewPlan: func(...) { return &generic.Plan{} },  // 创建空的计划对象
})

// Wrapper的逻辑在：agents/wrap_plan.go，放在agents/planner/下更合适？
return agents.NewWrite2PlanMDWrapper(a, op)  // **custom adk.agent, 包装器（Wrapper）模式：在不修改原智能体a的情况下，增加“将计划写入 Markdown 文件”的功能。**
```



`generic/plan.go`：实现plan interface（这里不使用defaultPlan）：

```go
type Step struct {
	Index int    `json:"index"`
	Desc  string `json:"desc"`
}

type Plan struct {
	Steps []Step `json:"steps"`
}

func (p *Plan) FirstStep() string {
	if len(p.Steps) == 0 {
		return ""
	}
	stepStr, _ := sonic.MarshalString(p.Steps[0])
	return stepStr
}

func (p *Plan) MarshalJSON() ([]byte, error) {
	type Alias Plan  // **使用类型别名避免循环依赖**
	return json.Marshal((*Alias)(p))
}

func (p *Plan) UnmarshalJSON(bytes []byte) error {
	type Alias Plan
	a := (*Alias)(p)
	return json.Unmarshal(bytes, a)
}

// 定义工具信息，约束 LLM 输出格式（包括参数类型校验等）
var PlanToolInfo = &schema.ToolInfo{  
	Name: "create_plan",
	Desc: "生成一个结构化的、分步骤的执行计划来解决给定的复杂任务。计划中的每个步骤必须分配给专门的智能体，并且必须有清晰、可执行的描述。",
	ParamsOneOf: schema.NewParamsOneOfByParams(
		map[string]*schema.ParameterInfo{
			"steps": {
				Type: schema.Array,
				Desc:     "要遵循的不同步骤，应按排序顺序排列",
				ElemInfo: &schema.ParameterInfo{
					Type: schema.Object,
					SubParams: map[string]*schema.ParameterInfo{
						"index": {
							Type:     schema.Integer,
							Desc:     "该步骤在整个计划中的序号。**必须从 1 开始，并且每个后续步骤必须恰好递增 1。**",
							Required: true,
						},
						"desc": {
							Type:     schema.String,
							Desc:     "该步骤要执行的具体任务的清晰、简洁和可执行的描述。它应该是分配给智能体的直接指令。",
							Required: true,
						},
					},
				},
				Required: true,
			},
		},
	),
}
sc := PlanToolInfo.ToJSONSchema() // 将 ParamsOneOf 转换为 JSON Schema，Eino v0.6开始统一使用json schema（行业标准）
```



`generic/full_plan.go`: 对基础 Plan 的扩展，用于表示和跟踪任务的完整执行状态

```go
type FullPlan struct {
	TaskID     int           `json:"task_id,omitempty"`
	Status     PlanStatus    `json:"status,omitempty"`  // plan执行状态
	AgentName  string        `json:"agent_name,omitempty"`  
	Desc       string        `json:"desc,omitempty"` // 负责的 Agent 名称
	ExecResult *SubmitResult `json:"exec_result,omitempty"` // plan执行结果
}

type PlanStatus string

const (
	PlanStatusTodo    PlanStatus = "todo"
	PlanStatusDoing   PlanStatus = "doing"
	PlanStatusDone    PlanStatus = "done"
	PlanStatusFailed  PlanStatus = "failed"
	PlanStatusSkipped PlanStatus = "skipped"
)
// 状态映射：将英文状态映射为中文显示
var (
	PlanStatusMapping = map[PlanStatus]string{
		PlanStatusTodo:    "待执行",
		PlanStatusDoing:   "执行中",
		PlanStatusDone:    "已完成",
		PlanStatusFailed:  "执行失败",
		PlanStatusSkipped: "已跳过",
	}
)

// 格式化方法：将状态转换为中文，生成 Markdown 格式字符串
// 1. **[已完成]** 读取 Excel 文件
// ### 执行结果
// 文件读取成功，共 100 行数据
func (p *FullPlan) String() string {
	status, ok := PlanStatusMapping[p.Status]
	if !ok {
		status = string(p.Status)
	}
	res := fmt.Sprintf("%d. **[%s]** %s", p.TaskID, status, p.Desc)
	if p.ExecResult != nil {
		res += fmt.Sprintf("\n%s", p.ExecResult.String())
	}
	return res
}

// 生成 Markdown 任务列表格式
// - [x] 1. 读取 Excel 文件
// - [ ] 2. 处理数据
func (p *FullPlan) PlanString(n int) string {
	if p.Status != PlanStatusDoing && p.Status != PlanStatusTodo {
		return fmt.Sprintf("- [x] %d. %s", n, p.Desc)
	}
	return fmt.Sprintf("- [ ] %d. %s", n, p.Desc)
}

// 工具函数：将 FullPlan 列表转换为 Markdown 字符串
// ### 任务计划
// - [x] 1. 读取 Excel 文件
// - [x] 2. 处理数据
// - [ ] 3. 生成报告
func FullPlan2String(plan []*FullPlan) string {
	var planStr = "### 任务计划\n"
	for i, p := range plan {
		planStr += p.PlanString(i+1) + "\n"
	}
	return planStr
}

// 将计划写入 plan.md
func Write2PlanMD(ctx context.Context, op commandline.Operator, wd string, plan []*FullPlan) error {
	planStr := FullPlan2String(plan)
	filePath := filepath.Join(wd, "plan.md")
	return op.WriteFile(ctx, filePath, planStr)
}
```



`generic/``**submit_result**``.go`：定义了任务执行结果的数据结构和相关工具函数，用于表示和格式化任务的最终执行结果。

- SubmitResult 作为 FullPlan 的执行结果字段
- 在 agents/wrap_plan.go 中，为每个已执行步骤创建 SubmitResult
```go
type SubmitResult struct {
	IsSuccess *bool               `json:"is_success,omitempty"`
	Result    string              `json:"result,omitempty"`
	Files     []*SubmitResultFile `json:"files,omitempty"`
}

type SubmitResultFile struct {
	Path string `json:"path,omitempty"`
	Desc string `json:"desc,omitempty"`
}

func (s *SubmitResult) String() string {
	res := fmt.Sprintf("### 执行结果\n%s", s.Result)
	if len(s.Files) > 0 {
		res += "\n#### 中间产物"
	}
	for _, f := range s.Files {
		res += fmt.Sprintf("\n- 描述：%s, 路径：%s", f.Desc, f.Path)
	}
	return res
}

func ListDir(dir string) ([]*SubmitResultFile, error) {
	var resp []*SubmitResultFile

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(d.Name(), ".") {
			return nil
		}
		if path == dir {
			return nil
		}
		if d.IsDir() {
			next := filepath.Join(dir, d.Name())
			nextResp, err := ListDir(next)
			if err != nil {
				return err
			}
			resp = append(resp, nextResp...)
			return nil
		}
		resp = append(resp, &SubmitResultFile{
			Path: filepath.Join(filepath.Dir(dir), d.Name()),
		})
		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

```



`tools/wrap.go`：装饰器模式，（Decorator Pattern），在不修改原始tool的情况下增强功能。

**执行流程图：**

```go
工具调用
↓
预处理阶段

- ToolRequestRepairJSON (修复 JSON)
↓
执行原始工具
- baseTool.InvokableRun()
↓
后处理阶段
- FilePostProcess (格式化输出)
- EditFilePostProcess (简化响应)
↓
返回处理后的响应
```

| 工具 | 预处理 | 后处理 | 用途 | 
| --- | --- | --- | --- | 
| bash | RepairJSON | FilePostProcess | 修复 JSON + 格式化命令输出 | 
| python_runner | RepairJSON | FilePostProcess | 修复 JSON + 格式化 Python 输出 | 
| edit_file | RepairJSON | EditFilePostProcess | 修复 JSON + 简化成功消息 | 
| read_file | RepairJSON | nil | 仅修复 JSON | 
| tree | RepairJSON | nil | 仅修复 JSON | 
| submit_result | RepairJSON | nil |   | 

```go

type ToolRequestPreprocess func(ctx context.Context, baseTool tool.InvokableTool, toolArguments string) (string, error)

type ToolResponsePostprocess func(ctx context.Context, baseTool tool.InvokableTool, toolResponse, toolArguments string) (string, error)

func NewWrapTool(t tool.InvokableTool, preprocess []ToolRequestPreprocess, postprocess []ToolResponsePostprocess) tool.InvokableTool {
	return &wrapTool{
		baseTool:    t,
		preprocess:  preprocess,
		postprocess: postprocess,
	}
}

type wrapTool struct {
	baseTool    tool.InvokableTool
	preprocess  []ToolRequestPreprocess
	postprocess []ToolResponsePostprocess
}

func (w *wrapTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return w.baseTool.Info(ctx)
}

func (w *wrapTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	for _, pre := range w.preprocess {
		var err error
		argumentsInJSON, err = pre(ctx, w.baseTool, argumentsInJSON)
		if err != nil {
			log.Printf("[WrapTool.PreProcess] failed to process response: %v", err)
		}
	}

	resp, err := w.baseTool.InvokableRun(ctx, argumentsInJSON, opts...)
	if err != nil {
		return "", err
	}

	for _, post := range w.postprocess {
		resp, err = post(ctx, w.baseTool, resp, argumentsInJSON)
		if err != nil {
			log.Printf("[WrapTool.PostProcess] failed to process response: %v", err)
			return resp, err
		}
	}

	return resp, nil
}

func ToolRequestRepairJSON(ctx context.Context, baseTool tool.InvokableTool, toolArguments string) (string, error) {
	return utils.RepairJSON(toolArguments), nil
}

type runResult struct {
	Command    string            `json:"command"`
	StdOut     []*stdoutData     `json:"stdout"`
	StdErr     []*stderrData     `json:"stderr"`
	FileChange []*fileChangeData `json:"file_change"`
	ErrData    []*errData        `json:"err_data"`
}

type stdoutData struct {
	Stdout string `json:"stdout"`
}

type stderrData struct {
	Stderr string `json:"stderr"`
}

type errData struct {
	Err string `json:"err"`
}

type fileChangeData struct {
	FileType       string          `json:"file_type"`
	Path           string          `json:"path"`
	Type           string          `json:"type"`
	Uri            string          `json:"uri"`
	Url            string          `json:"url"`
	MultiMediaInfo *multiMediaInfo `json:"multi_media_info,omitempty"`
}

type multiMediaInfo struct {
	MediaType      string `json:"media_type"`
	AdditionalType string `json:"additional_type"`
	AdditionalInfo string `json:"additional_info"`
}

func FilePostProcess(ctx context.Context, baseTool tool.InvokableTool, toolResponse, toolArguments string) (string, error) {
	rawResult := runResult{}
	if err := json.Unmarshal([]byte(toolResponse), &rawResult); err != nil {
		return toolResponse, nil
	}

	type fileOutputFormat struct {
		FileType string `json:"Change subject (file/directory)"`
		Path     string `json:"File/directory relative path"`
		Type     string `json:"Change type (create/delete/update)"`
	}
	var (
		stdOut     []string
		fileChange []fileOutputFormat
		stdErr     []string
	)
	for _, item := range rawResult.StdOut {
		if item != nil {
			stdOut = append(stdOut, item.Stdout)
		}
	}
	for _, item := range rawResult.FileChange {
		if item != nil {
			fileChange = append(fileChange, fileOutputFormat{
				FileType: item.FileType,
				Path:     item.Path,
				Type:     item.Type,
			})
		}
	}
	for _, item := range rawResult.StdErr {
		if item != nil {
			stdErr = append(stdErr, item.Stderr)
		}
	}
	for _, item := range rawResult.ErrData {
		if item != nil {
			stdErr = append(stdErr, item.Err)
		}
	}

	var output string
	if len(fileChange) > 0 {
		fcText, _ := jsoniter.MarshalToString(fileChange)
		output += "fileChange: \n" + fcText + "\n"
	}
	if len(stdErr) > 0 {
		output += "shell command stderr and warnings:" + strings.Join(stdErr, "\n") + "\n"
	}
	if len(rawResult.StdOut) > 0 {
		output += "shell command stdout: " + strings.Join(stdOut, "\n") + "\n"
	}

	return output, nil
}

func EditFilePostProcess(ctx context.Context, baseTool tool.InvokableTool, toolResponse, toolArguments string) (string, error) {
	return fmt.Sprintf("Write file: %s success!", toolResponse), nil
}

func isImage(uri string) bool {
	ext := filepath.Ext(uri)
	for _, e := range []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".tiff", ".heic"} {
		if ext == e {
			return true
		}
	}
	return false
}
```





核心代码-第二层executor：

```go
var executorPrompt = prompt.FromMessages(schema.FString,
	schema.SystemMessage(`你是一个认真细致的执行Agent。请仔细遵循给定的计划并彻底执行任务。

可用工具：
- CodeAgent：这是一个专门处理 Excel 文件的代码代理。它接收分步计划，通过生成 Python 代码（利用 pandas 进行数据分析/操作，matplotlib 进行绘图/可视化，openpyxl 进行 Excel 读写）来处理每个任务，并按顺序执行任务。当需要为 Excel 操作进行分步 Python 编码时，React Agent应该调用它，以确保精确、高效地完成任务。

注意事项：
- 不要转移给其他代理，仅使用工具。
`),
	schema.UserMessage(`## 目标
{input}
## 给定以下计划：
{plan}
## 已完成步骤和结果
{executed_steps}
## 你的任务是执行第一步，即：
{step}`))

cm, err := utils.NewChatModel(ctx,
		utils.WithMaxTokens(4096),
		utils.WithTemperature(float32(0)),
		utils.WithTopP(float32(0)),
	)
ca, err := newCodeAgent(ctx, operator)
sa, err := newWebSearchAgent(ctx)

a, err := planexecute.NewExecutor(ctx, &planexecute.ExecutorConfig{
		Model: cm,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{
					adk.NewAgentTool(ctx, ca),  // CodeAgent as tool
					adk.NewAgentTool(ctx, sa),  // WebSearchAgent as tool
				},
			},
		},
		MaxIterations: 20,
		GenInputFn: func(ctx context.Context, in *planexecute.ExecutionContext) ([]adk.Message, error) { // prompt，使用上述executorPrompt模板
			planContent, err := in.Plan.MarshalJSON()
			return executorPrompt.Format(ctx, map[string]any{
				"input":          utils.FormatInput(in.UserInput),
				"plan":           string(planContent),
				"executed_steps": utils.FormatExecutedSteps(in.ExecutedSteps),
				"step":           in.Plan.FirstStep(),
			})
		},
	})
```



核心代码-第二层replanner：

```go
var (
	replannerPromptTemplate = prompt.FromMessages(schema.Jinja2,
		schema.SystemMessage(`你是一位专门从事 Excel 数据处理任务的专家规划师。你的目标是理解用户需求并将其分解为清晰的、分步骤的计划。

**1. 理解目标：**
- 仔细分析用户的请求，确定最终目标。
- 识别输入数据（Excel 文件）和期望的输出格式。

**2. 交付物：**
- 最终输出应该是一个表示计划的 JSON 对象，包含步骤列表。
- 每个步骤必须是对执行该步骤的智能体的清晰且简洁的指令。

**3. 计划分解原则：**
- **粒度：** 将任务分解为尽可能小的逻辑步骤。例如，不要使用"处理数据"，而应使用"读取 Excel 文件"、"过滤掉缺失值的行"、"计算 'Sales' 列的平均值"等。
- **顺序：** 步骤应按正确的执行顺序排列。
- **清晰度：** 每个步骤应该明确无误，易于执行该步骤的智能体理解。

**4. 输出格式（少样本示例）：**
以下是一个良好计划的示例：
用户请求："请计算附件 'sales_data.xlsx' 文件中每个产品类别的平均销售额，并生成报告。"
{
  "steps": [
    {
      "instruction": "将 'sales_data.xlsx' 文件读取到 pandas DataFrame 中。"
    },
    {
      "instruction": "按 'Product Category' 对 DataFrame 进行分组，并计算每个组的 'Sales' 列的平均值。"
    },
    {
      "instruction": "总结每个产品类别的平均销售额，并以表格形式呈现结果。"
    }
  ]
}

**5. 限制：**
- 不要在计划中直接生成代码。
- 确保计划是逻辑合理且可实现的。
- 最后一步应该始终是生成报告或提供最终结果。

**6. 重新规划：**  // 明确指出如何判断是否终止loop
- 如果当前计划已完成，调用 'submit_result' 工具。
- 如果计划需要修改或扩展，使用新计划调用 'create_plan' 工具。  create_plan不是一个eino tool，只用于传递结构化数据（计划），而不是执行实际逻辑
`),
		schema.UserMessage(`
用户查询：{{ user_query }}
当前时间：{{ current_time }}
文件预览：
{{ file_preview }}
已执行步骤：{{ executed_steps }}
剩余步骤：{{ remaining_steps }}
`),
	)
)

cm, err := utils.NewChatModel(ctx,
		utils.WithMaxTokens(4096),
		utils.WithTopP(0),
		utils.WithTemperature(1.0),  // 支持随机性（创造性调整）
		utils.WithDisableThinking(true),
	)

respondInfo, err := tools.NewToolSubmitResult(op).Info(context.Background())
a, err := planexecute.NewReplanner(ctx, &planexecute.ReplannerConfig{
		ChatModel:   cm,
		PlanTool:    generic.PlanToolInfo,  // **用于创建/修改计划的tool**
		RespondTool: respondInfo,  //  **用于提交结果并结束任务的tool**
		GenInputFn:  func(ctx context.Context, in *planexecute.ExecutionContext) ([]adk.Message, error) {
				pf, _ := params.GetTypedContextParams[string](ctx, params.UserAllPreviewFilesSessionKey)
				plan, ok := in.Plan.(*generic.Plan)
			
				// remove the first step
				plan.Steps = plan.Steps[1:]
				planStr, err := sonic.MarshalString(plan)
			
				userInput, err := sonic.MarshalString(in.UserInput)
			
				return replannerPromptTemplate.Format(ctx, map[string]any{
					"current_time":    utils.GetCurrentTime(),
					"file_preview":    pf,
					"user_query":      userInput,
					"remaining_steps": planStr,
					"executed_steps":  utils.FormatExecutedSteps(in.ExecutedSteps),
				})
			},
	
		NewPlan: func(ctx context.Context) planexecute.Plan {
			return &generic.Plan{}
		},
	})
return agents.NewWrite2PlanMDWrapper(a, op) // 同planner
```



核心代码-第二层NewReportAgent：

```go
cm, err := utils.NewChatModel(ctx,
		utils.WithMaxTokens(15000),
		utils.WithTemperature(0.1),
		utils.WithTopP(1),
	)
	
var imageReaderTool tool.InvokableTool
imageReaderTool = tools.NewToolImageReader(visionModel) // 视觉模型
preprocess := []tools.ToolRequestPreprocess{tools.ToolRequestRepairJSON}
	agentTools := []tool.BaseTool{
		tools.NewWrapTool(tools.NewBashTool(operator), preprocess, nil),
		tools.NewWrapTool(tools.NewTreeTool(operator), preprocess, nil),
		tools.NewWrapTool(tools.NewEditFileTool(operator), preprocess, nil),
		tools.NewWrapTool(tools.NewReadFileTool(operator), preprocess, nil),
		tools.NewWrapTool(tools.NewToolSubmitResult(operator), preprocess, nil),
	}
	agentTools = append(agentTools, tools.NewWrapTool(imageReaderTool, preprocess, nil))
	

ra, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name: "Report",
		Description: `这是一个报告代理，负责从给定的文件路径读取文件，并基于其内容生成综合报告。
其工作流程包括读取文件、分析数据和信息、总结关键发现和洞察，并生成清晰、简洁的报告以回答用户的查询。
如果文件包含图表或可视化内容，代理将在报告中适当引用它们。当需要从指定文件生成详细的数据驱动报告时，React 代理应该调用此子代理。`,
		Instruction: `你是一个报告代理。你的任务是读取给定文件路径的文件，并基于其内容生成综合报告。

**工作流程：**
1.  读取由"输入文件路径"和"工作目录"指定的文件内容。
2.  分析文件中的数据和信息。
3.  总结关键发现和洞察。
4.  生成清晰、简洁的报告以回答用户的查询。
5.  如果有任何图表或可视化内容，请在报告中引用它们。
6.  如果工作完成，必须在结束前调用 SubmitResult 工具。
`,
		Model: cm,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: agentTools,
			},
			ReturnDirectly: tools.SubmitResultReturnDirectly,
		},
		GenModelInput: func(ctx context.Context, instruction string, input *adk.AgentInput) ([]adk.Message, error) {
			planExecuteResult := input.Messages
			if len(input.Messages) > 0 && input.Messages[len(input.Messages)-1].Role == schema.Tool {
				planExecuteResult = []*schema.Message{input.Messages[len(input.Messages)-1]}
			}

			fp, ok := params.GetTypedContextParams[string](ctx, params.FilePathSessionKey)

			plan, ok := utils.GetSessionValue[*generic.Plan](ctx, planexecute.PlanSessionKey)
			if !ok {
				return nil, fmt.Errorf("plan not found")
			}

			planStr, err := json.MarshalIndent(plan, "", "\t")

			wd, ok := params.GetTypedContextParams[string](ctx, params.WorkDirSessionKey)

			files, err := generic.ListDir(wd)

			tpl := prompt.FromMessages(schema.Jinja2,
				schema.SystemMessage(instruction),
				schema.UserMessage(`
User Query: {{ user_query }}
Input File Path: {{ file_path }}
Working Directory: {{ work_dir }}
Working Directory Files: {{ work_dir_files }}
Current Time: {{ current_time }}

**Plan Details:**
{{ plan }}
`))

			msgs, err := tpl.Format(ctx, map[string]any{
				"file_path":      fp,
				"work_dir":       wd,
				"work_dir_files": utils.ToJSONString(files),
				"user_query":     utils.FormatInput(planExecuteResult),
				"plan":           string(planStr),
				"current_time":   utils.GetCurrentTime(),
			})
			if err != nil {
				return nil, err
			}

			return msgs, nil
		},
		MaxIterations: 20,
	})
	
	
func ToolRequestRepairJSON(ctx context.Context, baseTool tool.InvokableTool, toolArguments string) (string, error) { // 修复json
	return utils.RepairJSON(toolArguments), nil
}
func RepairJSON(input string) string {
	input = strings.TrimPrefix(input, "<|FunctionCallBegin|>")
	input = strings.TrimSuffix(input, "<|FunctionCallEnd|>")
	input = strings.TrimPrefix(input, "<think>")
	output, err := jsonrepair.JSONRepair(input)
	if err != nil {
		return input
	}

	return output
}
```



核心代码-第二层read_image tool：

```go
var (
	toolImageReaderInfo = &schema.ToolInfo{
		Name: "image_reader",
		Desc: "Tool for describing image content",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"query": {
				Type:     "string",
				Desc:     "Questions posed about the image",
				Required: true,
			},
			"image_path": {
				Type:     "string",
				Desc:     "The path of the image file",
				Required: true,
			},
		}),
	}
)

func NewToolImageReader(visionModel model.BaseChatModel) tool.InvokableTool {
	return &localToolImageReader{visionModel: visionModel}
}

type localToolImageReader struct {
	visionModel model.BaseChatModel
}

func (t *localToolImageReader) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return toolImageReaderInfo, nil
}
// **直接实现 InvokableTool interface**
func (t *localToolImageReader) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	var params struct {
		Query     string `json:"query"`
		ImagePath string `json:"image_path"`
	}

	f, err := os.Open(params.ImagePath)

	defer f.Close()
	fc, err := io.ReadAll(f)

  // **图片转成base64**
	mimeType := http.DetectContentType(fc)
	b64 := base64.StdEncoding.EncodeToString(fc)
	url := fmt.Sprintf("data:%s;base64,%s", mimeType, b64)
	
	// **带base64图片的prompt**
	msgs := []*schema.Message{  // 
		schema.SystemMessage(""), // TODO: fill system prompt
		schema.UserMessage(params.Query),
		{
			Role: schema.User,
			UserInputMultiContent: []schema.MessageInputPart{
				{
					Type: schema.ChatMessagePartTypeImageURL,
					Image: &schema.MessageInputImage{
						MessagePartCommon: schema.MessagePartCommon{
							URL:      &url,
							MIMEType: mimeType,
						},
						Detail: "",
					},
				},
			},
		},
	}

	resp, err := t.visionModel.Generate(ctx, msgs) // 调用prompt

	return resp.Content, nil
}

```



**Excel Agent**：是一个“看得懂 Excel 的智能助手”，它先把问题拆解成步骤，再一步步执行并校验结果。它能理解用户问题与上传的文件内容，提出可行的解决方案，并选择合适的工具（系统命令、生成并运行 Python 代码、网络查询等等）完成任务。

- **更稳定的产出质量**，通过“规划—执行—反思”闭环减少漏项与错误
- **更强的可扩展性**，各 Agent 独立构建，低耦合利于迭代更新。
- **更少的人工操作**，把复杂繁琐的 Excel 处理工作交给 Agent 自动完成。


**架构图**：

- **规划者（Planner）**：分析用户输入，拆解用户问题为可执行的计划
- **执行者（Executor）**：正确执行当前计划中的首个步骤
    - **CodeAgent**：接收来自 Executor 的指令，调用多种工具（例如读写文件，运行 python 代码等）完成任务
    - **WebSearchAgent**：接收来自 Executor 的指令，进行网络搜索
- **反思者（Replanner）**：根据 Executor 执行的结果和现有规划，决定继续执行、调整规划或完成执行
- **ReportAgent**：根据运行过程与结果，生成总结性质的报告
![](/images/2b624637-29b5-8095-9c5f-d7dcc209a2be/image_2ad24637-29b5-80ba-81b7-deda9fc1f5d3.jpg)

**运行动线图**：

![](/images/2b624637-29b5-8095-9c5f-d7dcc209a2be/image_2ad24637-29b5-8095-a15b-f6ef2743a5cf.jpg)











Eino中的Multi-Agent自定义架构要如何设计与实现？[使用Eino框架实现DeerFlow系统](https://mp.weixin.qq.com/s?__biz=Mzg2MTc0Mjg2Mw%3D%3D&mid=2247495153&idx=1&sn=e207794d53c6bc8256c5f8784aa13218&scene=21#wechat_redirect)

# 二、Supervisor MultiAgent范式（中心化协调模式）

## 源码解读

ADK 提供的一种中心化 Multi-Agent 协作模式，旨在为集中决策与分发执行的通用场景提供解决方案。由一个 Supervisor Agent（监督者） 和多个 SubAgent （子 Agent）组成，其中：

- **Supervisor Agent**：作为协作核心， 负责任务的分配（如基于规则或 LLM 决策）、子 Agent 完成后的结果汇总与下一步决策。
- **SubAgents：**专注于执行具体任务。子 Agent 完成后，`WithDeterministicTransferTo` 触发 Transfer 事件，将任务转让回 Supervisor，确保在完成后自动将任务控制权交回 Supervisor。


Supervisor 模式有如下特点：

- **中心化控制**：Supervisor 统一管理子 Agent，可根据输入与子 Agent 执行结果动态调整任务分配。
- **确定性回调**：子 Agent 执行完毕后会将运行结果返回到 Supervisor Agent，避免协作流程中断。
- **松耦合扩展**：子 Agent 可独立开发、测试和替换，方便拓展与维护。只需确保**实现 Agent 接口并绑定到 Supervisor**，即可接入协作流程。


非常适合于**动态协调多个专业 Agent 完成复杂任务**的场景；

可能的应用场景有：

- **科研项目管理**：Supervisor 分配 **调研、实验、报告撰写** 任务给不同子 Agent。
- **客户服务流程**：Supervisor 根据用户问题类型，分配给技术支持、售后、销售等子 Agent。




<!-- 列布局开始 -->

![](/images/2b624637-29b5-8095-9c5f-d7dcc209a2be/image_29d24637-29b5-8081-88a3-d7a7cb251f67.jpg)

![](/images/2b624637-29b5-8095-9c5f-d7dcc209a2be/image_29d24637-29b5-80f6-8746-dbe75cebd817.jpg)


---

example:

```go
import github.com/cloudwego/eino/adk/prebuilt/supervisor

// 科研项目管理：创建一个监督者模式的 multi-agent
// 包含 research（调研），experimentation（实验），report（报告）三个子 Agent
supervisor, err := supervisor.New(ctx, &supervisor.Config{
    SupervisorAgent: supervisorAgent,
    SubAgents: []adk.Agent{
        researchAgent,
        experimentationAgent,
        reportAgent,
    },
})
```



<!-- 列布局结束 -->

完整example：[example：integration-project-manager](https://www.notion.so/2b62463729b580959c5fd7dcc209a2be#2982463729b580228d42f9b03c4d426f) 



`WithDeterministicTransferTo`

是 Eino ADK 提供的 Agent 增强工具，用于为 Agent 注入任务转让（Transfer）能力 。它允许开发者为目标 Agent 预设固定的任务转让路径，当该 Agent 完成任务（未被中断）时，会自动生成 Transfer 事件，将任务流转到预设的目标 Agent。

是构建 Supervisor Agent 协作模式的基础，**确保子 Agent 在执行完毕后能可靠地将任务控制权交回监督者（Supervisor）**，形成“分配-执行-反馈”的闭环协作流程。

```go
// 包装方法
func AgentWithDeterministicTransferTo(_ context.Context, config *DeterministicTransferConfig) Agent

// 配置详情
type DeterministicTransferConfig struct {
    Agent        Agent          // 被增强的目标 Agent
    ToAgentNames []string       // 任务完成后转让的目标 Agent 名称列表
}
```



## `example：supervisor`

```go
// 核心一行代码
sv := supervisor.New(Supervisor: sv,SubAgents:  []adk.Agent{searchAgent, mathAgent}

ctx, endSpanFn := startSpanFn(ctx, "Supervisor", query)
endSpanFn(ctx, lastMessage)
```

supervisor agent

```go
traceCloseFn, startSpanFn := trace.AppendCozeLoopCallbackIfConfigured(ctx)
defer traceCloseFn(ctx)

sv, err := adk.NewChatModelAgent()


searchAgent, err := buildSearchAgent(ctx)
mathAgent, err := buildMathAgent(ctx)
sv := supervisor.New(Supervisor: sv,SubAgents:  []adk.Agent{searchAgent, mathAgent} // adk/prebuilt/supervisor

query := "find US and New York state GDP in 2024. what % of US GDP was New York state?"
runner := adk.NewRunner(sv)

ctx, endSpanFn := startSpanFn(ctx, "Supervisor", query)
iter := runner.Query(ctx, query)

for {
		event, hasEvent := iter.Next()
		prints.Event(event)
}
endSpanFn(ctx, lastMessage)
```



## `example：layered-supervisor`

```go
// 核心一行代码
sv, err := supervisor.New(ctx, &supervisor.Config{
		Supervisor: sv,
		SubAgents:  []adk.Agent{searchAgent, mathAgent},
	})

mathAgent := supervisor.New(ctx, &supervisor.Config{
		Supervisor: mathA,
		SubAgents:  []adk.Agent{sa, ma, da},
	})
```



1个supervisor agent下有嵌套1个supervisor subagent

```go
sv, err := supervisor.New(ctx, &supervisor.Config{
		Supervisor: sv,
		SubAgents:  []adk.Agent{searchAgent, mathAgent},
	})

mathAgent := supervisor.New(ctx, &supervisor.Config{
		Supervisor: mathA,
		SubAgents:  []adk.Agent{sa, ma, da},
	})


query := "find US and New York state GDP in 2024. what % of US GDP was New York state? " +
		"Then multiply that percentage by 1.589."
ctx, endSpanFn := startSpanFn(ctx, "layered-supervisor", query)

iter := adk.NewRunner(ctx, adk.RunnerConfig{
		EnableStreaming: true,
		Agent:           sv,
	}).Query(ctx, query)
	
var lastMessage adk.Message
for {
		event, hasEvent := iter.Next()
		if !hasEvent {
			break
		}

		prints.Event(event)

		if event.Output != nil {
			lastMessage, _, err = adk.GetMessage(event)
		}
	}

endSpanFn(ctx, lastMessage)

// wait for all span to be ended
time.Sleep(5 * time.Second)
```



## `example：integration-project-manager`

![](/images/2b624637-29b5-8095-9c5f-d7dcc209a2be/image_29a24637-29b5-80a9-b46f-d6006b62d853.jpg)

**详情**：[https://mp.weixin.qq.com/s/p_QqDN6m2anHAE97P2Q2bw?forceh5=1](https://mp.weixin.qq.com/s/p_QqDN6m2anHAE97P2Q2bw?forceh5=1)

`ProjectManagerAgent`：项目开发经理Agent（使用 Supervisor 模式）：根据动态的用户输入，路由并协调多个负责不同维度工作的子智能体开展工作。

1. `ResearchAgent`(调研Agent): 负责调研并生成可行方案。支持中断后从用户处接收额外的上下文信息来提高调研方案生成的准确性。
1. `CodeAgent`(编码 Agent)：使用知识库工具，召回相关知识作为参考，生成高质量的代码。
1. `ReviewAgent`(评论 Agent)：使用顺序工作流编排问题分析、评价生成、评价验证三个步骤，对调研结果/编码结果进行评审，给出合理的评价，供项目经理进行决策。
    1. questionAnalysisAgent
    1. generateReviewAgent
    1. reviewValidationAgent
```go
// 核心一行代码
supervisorAgent, err := supervisor.New(ctx, &supervisor.Config{
		Supervisor: s,
		SubAgents:  []adk.Agent{researchAgent, codeAgent, reviewAgent},
	})
	
researchAgent :=  带webSearchTool、newAskForClarificationTool
codeAgent := 带knowledgeBaseTool（召回RAG知识库）
reviewAgent := adk.NewSequentialAgent(
		SubAgents:   []adk.Agent{questionAnalysisAgent, generateReviewAgent, reviewValidationAgent},
	})


runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           supervisorAgent,
		EnableStreaming: true,
		CheckPointStore: newInMemoryStore(),
	})
	
// 循环中断和恢复
for !finished {
		var iter *adk.AsyncIterator[*adk.AgentEvent]

		if !interrupted {
			iter = runner.Query(ctx, query, adk.WithCheckPointID(checkpointID))
		} else {
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("\ninput additional context for web search: ")
			scanner.Scan()
			fmt.Println()
			nInput := scanner.Text()

			iter, err = runner.Resume(ctx, checkpointID, adk.WithToolOptions([]tool.Option{agents.WithNewInput(nInput)}))
			if err != nil {
				log.Fatal(err)
			}
		}

		interrupted = false

		for {
			event, ok := iter.Next()
			if !ok {
				if !interrupted {
					finished = true
				}
				break
			}
			if event.Err != nil {
				log.Fatal(event.Err)
			}
			if event.Action != nil {
				if event.Action.Interrupted != nil {
					interrupted = true
				}
				if event.Action.Exit {
					finished = true
				}
			}
			prints.Event(event)
		}
	}
```



**核心代码:**

```go
// Init chat model for agents
tcm, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{..})

// Init research agent
researchAgent, err := agents.NewResearchAgent(ctx, tcm)

// Init code agent
codeAgent, err := agents.NewCodeAgent(ctx, tcm)

// Init technical agent
reviewAgent, err := agents.NewReviewAgent(ctx, tcm)
 
// Init project manager agent
s, err := agents.NewProjectManagerAgent(ctx, tcm)


// Combine agents into ADK supervisor pattern
// Supervisor: project manager
// Sub-agents: researcher / coder / reviewer
supervisorAgent, err := supervisor.New(ctx, &supervisor.Config{
   Supervisor: s,
   SubAgents:  []adk.Agent{researchAgent, codeAgent, reviewAgent},
})

// Init Agent runner
runner := adk.NewRunner(ctx, adk.RunnerConfig{
   Agent:           supervisorAgent,
   EnableStreaming: true,// enable stream output
   CheckPointStore: newInMemoryStore(),// enable checkpoint for interrupt & resume
})

query := "please generate a simple ai chat project with python."
checkpointID := "1"

// Start runner with a new checkpoint id
iter := runner.Query(ctx, query, adk.WithCheckPointID(checkpointID))
interrupted := false
for {
       event, ok := iter.Next()
       if !ok {
          break
       }
       if event.Err != nil {
          log.Fatal(event.Err)
       }
       if event.Action != nil && event.Action.Interrupted != nil {
          interrupted = true
       }
       prints.Event(event)
    }

    if !interrupted {
       return
}

// interrupt and ask for additional user context
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Print("\ninput additional context for web search: ")
    scanner.Scan()
    fmt.Println()
    nInput := scanner.Text()

// Resume by checkpoint id, with additional user context injection
    iter, err = runner.Resume(ctx, checkpointID, adk.WithToolOptions([]tool.Option{agents.WithNewInput(nInput)}))
    if err != nil {
       log.Fatal(err)
    }
    for {
       event, ok := iter.Next()
       if !ok {
          break
       }
       if event.Err != nil {
          log.Fatal(event.Err)
       }
       prints.Event(event)
    }
}
```



如果不用Eion，从0开发：

| **设计点** | **基于 Eino ADK 开发** | **传统开发模式** | 
| --- | --- | --- | 
| **Agent 抽象** | 统一定义，职责独立，代码整洁，便于各 Agent 分头开发 | 没有统一定义，团队协作开发效率差，后期维护成本高 | 
| **输入输出** | 有统一定义，全部基于事件驱动运行过程通过 iterator 透出，所见即所得 | 没有统一定义，输入输出混乱运行过程只能手动加日志，不利于调试 | 
| **Agent 协作** | 框架自动传递上下文 | 通过代码手动传递上下文 | 
| **中断恢复能力** | 仅需在 Runner 中注册 CheckPointStore 提供断点数据存储介质 | 需要从零开始实现，解决序列化与反序列化、状态存储与恢复等问题 | 
| **Agent 模式** | 多种成熟模式开箱即用 | 需要从零开始实现 | 

# 三、deepAgent MultiAgent范式

## 源码解读

deepAgent论文地址（中国人大学生和小红书实习员工的论文）：[https://arxiv.org/html/2510.21618?_immersive_translate_auto_translate=1](https://arxiv.org/html/2510.21618?_immersive_translate_auto_translate=1)

- 对比 Plan-and-Execute
    - 优势：DeepAgents 将 Plan/RePlan 作为工具供主 Agent 自由调用，可以在任务中跳过不必要的规划，整体上减少模型调用次数、降低耗时与成本。
    - 劣势：任务规划与委派由一次模型调用完成，对模型能力要求更高，提示词调优也相对更困难
- 对比 Supervisor（ReAct）
    - 优势：DeepAgents 通过内置 WriteTodos 强化任务拆解与规划；同时隔离多 Agents 上下文，在大规模、多步骤任务中通常效果更优。
    - 劣势：制定计划与调用子 Agent 会带来额外的模型请求，增加耗时与 token 成本；若任务拆分不合理，可能对效果产生反作用。


[https://www.cloudwego.io/zh/docs/eino/core_modules/eino_adk/agent_implementation/deepagents/](https://www.cloudwego.io/zh/docs/eino/core_modules/eino_adk/agent_implementation/deepagents/)

![](/images/2b624637-29b5-8095-9c5f-d7dcc209a2be/image_2ad24637-29b5-80e6-9ec0-eba159b11624.jpg)



## `example：deepAgent`范式







# 四、Human-in-the-Loop

### 源码解读

![](/images/2b624637-29b5-8095-9c5f-d7dcc209a2be/image_2b324637-29b5-8076-986e-c9ae1ebce42b.jpg)



**三个主要参与者之间按时间顺序的交互流程：**

![](/images/2b624637-29b5-8095-9c5f-d7dcc209a2be/image_2b324637-29b5-806b-ad56-dd5f6834c548.jpg)







**理解human-in-the-loop的需求：**

![](/images/2b624637-29b5-8095-9c5f-d7dcc209a2be/image_2b324637-29b5-801c-942d-dec9ffe7049c.jpg)

因此，总结我们的目标是：

1. 帮助开发者尽可能轻松地回答上述问题。
1. 帮助最终用户尽可能轻松地回答上述问题。
1. 使框架能够自动并开箱即用地回答上述问题。
## approval **审批模式**

非常适合不可逆操作，如删除文件、数据库修改、金融交易。

InvokableApprovableTool 是 eino-examples 提供的一个 tool 装饰器，可以为任意的 InvokableTool 加上“审批中断”功能。

Eino 用 CheckPointStore 来保存 Agent 中断时的运行状态。用 CheckPointID 来唯一标识和串联“中断前”和“中断后”的两次（或多次）运行。

- 这里用的 InMemoryStore，保存在内存中。
- 实际使用中，**推荐用分布式存储比如 redis**。
用InterruptID(event.Action.Interrupted.InterruptContexts[0].ID) 来标识“哪里发生了中断”。这里直接打印在了终端上，实际使用中，可能需要作为 HTTP 响应返回给前端。

## review-and-edit **审查与编辑模式**

允许在执行前进行人工审查和原地编辑工具调用参数。非常适合纠正误解。



## feedback-loop **反馈循环模式**

agent 生成内容，人类提供定性反馈以进行改进。



## Follow-up **追问模式**

agent识别出不充分的工具输出并请求澄清或下一步行动。

