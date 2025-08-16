---
title: "5.2 components 组件"
date: 2025-08-13T01:45:00Z
draft: false
weight: 5002
---

# 5.2 components 组件

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

![](/images/24e24637-29b5-8060-8b1c-e32bf1c8d93e/image_24324637-29b5-8096-91e6-d6bca621cd22.jpg)



### **Callback**

> Eino **Callback 用户手册**：[https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/callback_manual/](https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/callback_manual/)





### **CallOption**

> **Eino CallOption 能力与规范**：[https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/call_option_capabilities/](https://www.cloudwego.io/zh/docs/eino/core_modules/chain_and_graph_orchestration/call_option_capabilities/)



