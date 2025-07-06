---
title: "透过Eino看AI应用开发：基本结构"
description: "关于透过Eino看AI应用开发：基本结构的详细说明"
date: 2025-05-19T04:49:00Z
draft: false

---

# 透过Eino看AI应用开发：基本结构

## 一、AI应用的本质

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



# 二、**Eino的整体结构**

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
| **核心库**[https://github.com/cloudwego/eino](https://github.com/cloudwego/eino) | **1. 核心Schema**：定义和抽象AI应用领域的领域模型的结构体。如流的读取和写入<br/>**2. Components**：8种组件类型的抽象。如各自的输入输出规范、流式处理范式和回调事件信息<br/>**3. Compose**：核心的编排能力部分。包含编排的基本组成元素和能力。<br/>&nbsp;&nbsp;a). **引擎层**：pregel和dag。区别在于节点间触发下一个节点的方式<br/>&nbsp;&nbsp;b). **API层**：<br/>&nbsp;&nbsp;&nbsp;&nbsp;◦ **Chain**：简单的链式有向图，只能向前推进<br/>&nbsp;&nbsp;&nbsp;&nbsp;◦ **Graph**：循环或非循环有向图。功能强大且灵活<br/>&nbsp;&nbsp;&nbsp;&nbsp;◦ **Workflow**：这是一种新型结构，能够**解耦数据流和控制流**，允许灵活配置数据流节点间的数据映射关系<br/>**4. Flow**：最顶层包含一些预置的编排产物，如 ReAct Agent 和 Multi Agent 等，这些都是基于compose builder能力构建的 | 
| **组件的具体实现库** [https://github.com/cloudwego/eino-ext](https://github.com/cloudwego/eino-ext) | 包括回调事件处理器的实现及可视化开发/调试工具，如IDE插件、[Eino Devops](https://github.com/cloudwego/eino-ext/tree/main/devops) | 
| **示例库** [https://github.com/cloudwego/eino-examples](https://github.com/cloudwego/eino-examples) |   | 
| **官方文档** | [https://www.cloudwego.io/zh/docs/eino/overview/bytedance_eino_practice/](https://www.cloudwego.io/zh/docs/eino/overview/bytedance_eino_practice/) | 





