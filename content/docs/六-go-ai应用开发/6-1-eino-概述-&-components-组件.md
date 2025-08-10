---
title: "6.1 Eino 概述 & Components 组件"
date: 2025-05-19T04:49:00Z
draft: false
weight: 6001
---

# 6.1 Eino 概述 & Components 组件

## 一、什么是Eino

> [https://mp.weixin.qq.com/s/Hyjpic0EMmmCHnxARjjUHA](https://mp.weixin.qq.com/s/Hyjpic0EMmmCHnxARjjUHA)

> 🎯 核心洞察：AI应用的本质是围绕大模型的信息流：Input（原始信息） → Process（大模型） → Output(最终信息)

![](/images/1f824637-29b5-80b2-be2f-c09695ffa6b1/image_22424637-29b5-80ee-92b7-c196e9371af5.jpg)

### **关键1：组件及其实现**

**组件定义表：**

| **组件(Interface)** | **组件实现(实例)** | 分工 | 
| --- | --- | --- | 
| **ChatTemplate**Prompt模板 |   | 负责信息的格式化与增强处理 | 
| **Retriever**知识库 | Redis、ES… | 负责从知识库中召回信息，进行必要的信息检索 | 
| **ChatModel**大模型 | OpenAI、Claude… | 负责推理和信息生成 | 
| **Tool**工具 | MCP、search… | 负责进一步的信息处理 | 

![](/images/1f824637-29b5-80b2-be2f-c09695ffa6b1/image_22424637-29b5-80a5-ac30-e2b892da475f.jpg)

### **关键2：组件之间的连接（信息流编排）**

> 💡 根据业务逻辑需求有效地串联起来，确保信息能够有效流动
> 1. local（框架内部）：由这些核心部分组成的数据流或信息流是无状态的
> 1. external（框架外部）：如外部存储或外部工具的API等，可插拔的（可以随时被替换）

**信息流结构表：**

| 结构 | 示例 | 
| --- | --- | 
| 扇出（一对多的拆解） | 从起始节点到Retriever | 
| 扇入（多对一的合并） | Retrieval再汇聚到ChatModel | 
| 分支判断 | 信息输出 or 继续信息处理 | 
| 回环 | 从Tool返回ChatModel | 

![](/images/1f824637-29b5-80b2-be2f-c09695ffa6b1/image_22424637-29b5-8097-bfb7-fcb390fb65bf.jpg)

### **关键3：数据流处理难题**

> 💡 **事件驱动设计，整个图的执行过程中的每一个步骤，甚至每一帧数据都被视作独立的事件**

**数据流处理表：**

| 结构 | 示例 | 
| --- | --- | 
| 流复制 | 复制成多个流，因为数据流一旦被消费可能就无法再被重复使用 | 
| 流合并 | 将多个独立的数据流合并成一个单一的数据流 | 
| 流拼接 | 将流式数据中的各帧数据组合成一个完整的数据集 | 
| 流装箱 | 转换回流式数据（由多帧数据组成一个序列） | 

![](/images/1f824637-29b5-80b2-be2f-c09695ffa6b1/image_22424637-29b5-8009-b3a2-f4de24d6aa1f.jpg)



> 附：还有Go的类型安全优势、AI大模型的统一集成等

---

### Eino的整体结构

> 💡 **开源社区优秀的AI应用开发框架**
> - [https://github.com/langchain-ai/langchain](https://github.com/langchain-ai/langchain) & [https://github.com/langchain-ai/langgraph](https://github.com/langchain-ai/langgraph)（python/js）
>     - 文档：[https://python.langchain.com/docs/introduction/](https://python.langchain.com/docs/introduction/)
>     - 官网：[https://www.langchain.com/](https://www.langchain.com/) 字节自己也在用
> - [https://github.com/run-llama/llama_index](https://github.com/run-llama/llama_index)（python/js）
>     - 官网：[https://www.llamaindex.ai/](https://www.llamaindex.ai/) 
>     - 文档：[https://docs.llamaindex.ai/en/stable/](https://docs.llamaindex.ai/en/stable/)

![](/images/1f824637-29b5-80b2-be2f-c09695ffa6b1/image_22524637-29b5-808b-b11f-dd19b975fd16.jpg)

| 代码仓库 | **Feature** | 
| --- | --- | 
| **核心库**[https://github.com/cloudwego/eino](https://github.com/cloudwego/eino) | **1. 核心Schema /'skiːmə/**：定义和抽象AI应用领域的领域模型的结构体。如流的读取和写入<br/>**2. Components /kəm'ponənt/ n. 组件**：8种组件类型的抽象。如各自的输入输出规范、流式处理范式和回调事件信息<br/>**3. Compose /kəm'poʊz/ vt 组成**：核心的编排能力部分。包含编排的基本组成元素和能力。<br/>&nbsp;&nbsp;a). **引擎层**：pregel和dag。区别在于节点间触发下一个节点的方式<br/>&nbsp;&nbsp;b). **API层**：<br/>&nbsp;&nbsp;&nbsp;&nbsp;◦ **Chain**：简单的链式有向图，只能向前推进<br/>&nbsp;&nbsp;&nbsp;&nbsp;◦ **Graph**：循环或非循环有向图。功能强大且灵活<br/>&nbsp;&nbsp;&nbsp;&nbsp;◦ **Workflow**：这是一种新型结构，能够**解耦数据流和控制流**，允许灵活配置数据流节点间的数据映射关系<br/>**4. Flow**：最顶层包含一些预置的编排产物，如 ReAct Agent 和 Multi Agent 等，这些都是基于compose builder能力构建的 | 
| **组件的具体实现库** [https://github.com/cloudwego/eino-ext](https://github.com/cloudwego/eino-ext) | 包括回调事件处理器的实现及可视化开发/调试工具，如IDE插件、[Eino Devops](https://github.com/cloudwego/eino-ext/tree/main/devops) | 
| **示例库** [https://github.com/cloudwego/eino-examples](https://github.com/cloudwego/eino-examples) |   | 
| **官方文档** | [https://cloudwego.cn/zh/docs/eino/overview/bytedance_eino_practice/](https://cloudwego.cn/zh/docs/eino/overview/bytedance_eino_practice/) | 

![](/images/1f824637-29b5-80b2-be2f-c09695ffa6b1/image_24b24637-29b5-8029-ba59-e3baef4fb70a.jpg)



### **Eino Components 组件**

大模型应用开发和传统应用开发最显著的区别在于大模型所具备的两大核心能力：

- **基于语义的文本处理能力**：能够理解和生成人类语言，处理非结构化的内容语义关系
- **智能决策能力**：能够基于上下文进行推理和判断，做出相应的行为决策
这两项核心能力进而催生了如下三种主要的AI应用类型。Eino基于这三种模式将这些常用能力抽象为可复用的「组件」（Components）

1. **对话处理类（Chat）**：处理用户输入并生成相应回答。ChatTemplate、ChatModel
1. **文本语义处理类（RAG）**：对文本文档进行语义化处理、存储和检索。Document.Loader 、Document.Transformer Embedding Indexer Retriever
1. **决策执行类（Tool call）**：基于上下文做出决策并调用相应工具。ToolsNode Lambda


## 二、**对话处理类（Chat）Components组件**

### **ChatTemplate**

> [https://www.cloudwego.io/zh/docs/eino/core_modules/components/chat_template_guide/](https://www.cloudwego.io/zh/docs/eino/core_modules/components/chat_template_guide/)

一个用于处理和格式化提示模板的组件。

主要作用：将用户提供的变量值填充到预定义的消息模板中，生成用于与语言模型交互的标准消息格式。

应用场景：

- 构建结构化的系统提示
- 处理多轮对话的模板 (包括 history)
- 实现可复用的提示模式
### **ChatModel**

> [https://www.cloudwego.io/zh/docs/eino/core_modules/components/chat_model_guide/](https://www.cloudwego.io/zh/docs/eino/core_modules/components/chat_model_guide/)

一个用于与大语言模型交互的组件。

主要作用：将用户的输入消息发送给语言模型，并获取模型的响应。

应用场景：

- 自然语言对话
- 文本生成和补全
- 工具调用的参数生成
- 多模态交互（文本、图片、音频等）
## 三、**文本语义处理类（RAG）Components组件**

### Document.Loader

> [https://www.cloudwego.io/zh/docs/eino/core_modules/components/document_loader_guide/](https://www.cloudwego.io/zh/docs/eino/core_modules/components/document_loader_guide/)

一个用于加载文档的组件。

应用场景：

- 从网络 URL 加载网页内容
- 读取本地 PDF、Word 等格式的文档
**Document.Parser**

一个用于解析文档内容的工具包。它不是一个独立的组件，而是作为 Document Loader 的内部工具。

应用场景：

- 解析不同格式的文档内容（如文本、PDF、Markdown 等）
- 根据文件扩展名自动选择合适的解析器 (eg：ExtParser)
- 为解析后的文档添加元数据信息
### **Document.Transformer**

> [https://www.cloudwego.io/zh/docs/eino/core_modules/components/document_transformer_guide/](https://www.cloudwego.io/zh/docs/eino/core_modules/components/document_transformer_guide/)

一个用于文档转换和处理的组件。它的主要作用是对输入的文档进行各种转换操作，如分割、过滤、合并等，从而得到满足特定需求的文档。

应用场景：

- 将长文档分割成小段落以便于处理
- 根据特定规则过滤文档内容
- 对文档内容进行结构化转换
- 提取文档中的特定部分
### **Embedding**

> [https://www.cloudwego.io/zh/docs/eino/core_modules/components/embedding_guide/](https://www.cloudwego.io/zh/docs/eino/core_modules/components/embedding_guide/)

一个用于将文本转换为向量表示的组件。主要作用是将文本内容映射到向量空间，使得语义相似的文本在向量空间中的距离较近。

应用场景：

- 文本相似度计算
- 语义搜索
- 文本聚类分析
### **Indexer**

> [https://www.cloudwego.io/zh/docs/eino/core_modules/components/indexer_guide/](https://www.cloudwego.io/zh/docs/eino/core_modules/components/indexer_guide/)

一个用于存储和索引文档的组件。它的主要作用是将文档及其向量表示存储到后端存储系统中，并提供高效的检索能力。

应用场景：

- 构建向量数据库，以用于语义关联搜索
### **Retriever**

> [https://www.cloudwego.io/zh/docs/eino/core_modules/components/retriever_guide/](https://www.cloudwego.io/zh/docs/eino/core_modules/components/retriever_guide/)

一个用于从各种数据源检索文档的组件。

主要作用：根据用户的查询（query）从文档库中检索出最相关的文档。

应用场景：

- 基于向量相似度的文档检索
- 基于关键词的文档搜索
- 知识库问答系统 (rag)


## 四、**决策执行类（Tool call） Components组件**

### **ToolsNode&Tool**

> [https://www.cloudwego.io/zh/docs/eino/core_modules/components/tools_node_guide/](https://www.cloudwego.io/zh/docs/eino/core_modules/components/tools_node_guide/)
[https://www.cloudwego.io/zh/docs/eino/core_modules/components/tools_node_guide/how_to_create_a_tool/](https://www.cloudwego.io/zh/docs/eino/core_modules/components/tools_node_guide/how_to_create_a_tool/)



一个用于扩展模型能力的组件，它允许模型调用外部工具来完成特定的任务。

应用场景：

- 让模型能够获取实时信息（如搜索引擎、天气查询等）
- 使模型能够执行特定的操作（如数据库操作、API 调用等）
- 扩展模型的能力范围（如数学计算、代码执行等）
- 与外部系统集成（如知识库查询、插件系统等）


### Lambda

> [https://www.cloudwego.io/zh/docs/eino/core_modules/components/lambda_guide/](https://www.cloudwego.io/zh/docs/eino/core_modules/components/lambda_guide/)

是Eino 中最基础的组件类型，它允许用户在工作流中嵌入自定义的函数逻辑。Lambda 组件底层是由输入输出是否流所形成的 4 种运行函数组成，对应 4 种交互模式: Invoke、Stream、Collect、Transform。

用户构建 Lambda 时可实现其中的一种或多种，框架会根据一定的规则进行转换。



## 五、Eino Compose 组合/**编排**

**编排：**对Eino Components 组件（原子能力）进行串联、组合

大模型应用的开发的特点：自定义的业务逻辑几乎都是**仅仅对『原子能力』的组合串联**。

Eino 对「编排」有着这样的洞察：

- 编排要独立在业务逻辑之上的清晰的一层，**不能让业务逻辑融入到编排中**。
    - 业务逻辑复杂度封装到组件内部，上层的编排层拥有更全局的视角，让**逻辑层次变得非常清晰。**
- 大模型应用的核心是 “对提供原子能力的组件” 进行组合串联，**组件是编排的 “第一公民”**。
- 抽象视角看编排：编排是在构建一张网络，数据则在这个网络中流动，网络的每个节点都对流动的数据有格式/内容的要求，一个能顺畅流动的数据网络，关键就是 “**上下游节点间的数据格式是否对齐**？”。
    - 提供了 “类型对齐” 的开发方式的强化，降低开发者心智负担，把 golang 的**类型安全**特性发挥出来
- 业务场景的复杂度会反映在编排产物的复杂性上，只有**横向的治理能力**才能让复杂场景不失控。
    - 提供了切面能力，callback 机制支持了基于节点的**统一治理能力。**
- 大模型是会持续保持高速发展的，大模型应用也是，只有**具备扩展能力的应用才拥有生命力**。
    - 提供了 call option 的机制，**扩展性**是快速迭代中的系统最基本的诉
于是，Eino 提供了 “基于 Graph 模型 (graph.AddXXXNode() + graph.AddEdge()) 的，以**组件**为原子节点的，以**上下游类型对齐**为基础的编排” 的解决方案。

而在现实的大多数业务场景中，往往仅需要 “按顺序串联” 即可，因此，Eino 封装了接口更易于使用的 `Chain`。Eino中的Chain 是对 Graph 的封装，除了 “环” 之外，Chain 暴露了几乎所有 Graph 的能力。



### **Chain/Graph 编排**

> [https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/chain_graph_introduction/#chain](https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/chain_graph_introduction/#chain) 
Eino**编排的设计理念**：[https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/orchestration_design_principles/](https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/orchestration_design_principles/)



### **Eino 流式编程**

> **Eino 流式编程要点**：[https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/stream_programming_essentials/](https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/stream_programming_essentials/)

![](/images/1f824637-29b5-80b2-be2f-c09695ffa6b1/image_24324637-29b5-8096-91e6-d6bca621cd22.jpg)



### **Callback**

> Eino **Callback 用户手册**：[https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/callback_manual/](https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/callback_manual/)





### **CallOption**

> **Eino CallOption 能力与规范**：[https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/call_option_capabilities/](https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/call_option_capabilities/)



