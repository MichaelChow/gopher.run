---
title: "5.10 透过Eino看AI应用开发"
date: 2025-05-19T04:49:00Z
draft: false
weight: 5010
---

# 5.10 透过Eino看AI应用开发

# 一、AI应用的本质

> 🎯 核心洞察：AI应用的本质是围绕大模型的信息流：Input（原始信息） → Process（大模型） → Output(最终信息)

![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/1c73908b-803e-4ce8-84c1-f2e3b9634e47/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB466YWFV4FRH%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005444Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJGMEQCIGPKT4xBWQOsvcusq1mdTjOKon0jKWWRXI4oXhp8wuVHAiAPtZomf3qcTl4iCboxWwadA1UY3GAtDnl641ZsfI92IiqIBAiZ%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F8BEAAaDDYzNzQyMzE4MzgwNSIMk8O2Nl0mSzR%2Bw5TfKtwDAFIOYF72QTLmSZwv1nSlC4mcnN5I2%2F%2F%2FRqkzM4oMaQipwqrm8yBaUeMvzK2FcYQbmc9kK6ONg1dgoZj3pz4jTwJtIvN9YV9x9ZaDtbfncR%2B559%2FbqAjfV10PZvELymtctl2NaeAyvbaiGPkxQnkgjydQP%2BDFOrXj4My0sTGo2peVcodFksB70G0lbHr18iE3z8AaH4gJWF03Wbp6PzpgvfaE3Jg2fl%2BNRELvK8Y60p8grzSyk9brOCnMTijKt%2FU65lpPs9dWWGwxBkKOYtPSMbDk%2FJEQUkPBPERgBkBp%2BHNa%2BjhW3i9wbmev7b98tGkVk7HwGOFW4I8eNym7WKsFbRJokqUP%2BEUvhLgBpPHXNmDTz%2BH32T7AubhwLw%2BJlho%2B9ELBdrER2yaGgs8FF8I0c%2BMgopcQLfhSbmuoE2wckh2deQyNdMjkt6EWsmKxj2ntYG0e6guVeUt9%2BUNQoVzXEBHBzCX4SMsEUHZuV%2BpEWGpEGJgNayG7KOApnNukdlkztgKBPh4B7YRW1jFcgzJC7XMMDt4TBDx5Zc0KntSLRUKCMhSAyP6S6VqUAEStHM0S9JWWABydP3NC1YvQOnX3ll83Lp4m5sYuvgph40Rn2%2BLHvCSrG63eHDJq7ZMw1rrrwwY6pgHQyYAIUtGCOcO8Uns9xcYcU4a0Orr4zf1wHfCzB1rkkl81BL2LIeVGBjoErQ%2BLG4tRJDxw1v8LOQKq0Bu0KXHXlIlQBfS8Vqhn9if%2FCl8fP8h7BfV89U1c83nUOZqq4JVd0dSqKp%2BGOcVQYYAAgRPy6ORg%2FdmO5zMulEAO6CA4ADJPzHmCV0T6kQlfoMJqOTunilc0Dz1j4gLoZoYDaXd1yz0Vighq&X-Amz-Signature=cd11e34229602da0d396239133001309aaceeac5d9bbe3a5f86419a9fc1cadf9&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)

### **关键1：组件及其实现**

**组件定义表：**

| **组件(Interface)** | **组件实现(实例)** | 分工 | 
| --- | --- | --- | 
| **ChatTemplate**Prompt模板 |   | 负责信息的格式化与增强处理 | 
| **Retriever**知识库 | Redis、ES… | 负责从知识库中召回信息，进行必要的信息检索 | 
| **ChatModel**大模型 | OpenAI、Claude… | 负责推理和信息生成 | 
| **Tool**工具 | MCP、search… | 负责进一步的信息处理 | 

![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/e43e0c79-db4a-45dd-870b-cb7af123867b/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB466YWFV4FRH%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005444Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJGMEQCIGPKT4xBWQOsvcusq1mdTjOKon0jKWWRXI4oXhp8wuVHAiAPtZomf3qcTl4iCboxWwadA1UY3GAtDnl641ZsfI92IiqIBAiZ%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F8BEAAaDDYzNzQyMzE4MzgwNSIMk8O2Nl0mSzR%2Bw5TfKtwDAFIOYF72QTLmSZwv1nSlC4mcnN5I2%2F%2F%2FRqkzM4oMaQipwqrm8yBaUeMvzK2FcYQbmc9kK6ONg1dgoZj3pz4jTwJtIvN9YV9x9ZaDtbfncR%2B559%2FbqAjfV10PZvELymtctl2NaeAyvbaiGPkxQnkgjydQP%2BDFOrXj4My0sTGo2peVcodFksB70G0lbHr18iE3z8AaH4gJWF03Wbp6PzpgvfaE3Jg2fl%2BNRELvK8Y60p8grzSyk9brOCnMTijKt%2FU65lpPs9dWWGwxBkKOYtPSMbDk%2FJEQUkPBPERgBkBp%2BHNa%2BjhW3i9wbmev7b98tGkVk7HwGOFW4I8eNym7WKsFbRJokqUP%2BEUvhLgBpPHXNmDTz%2BH32T7AubhwLw%2BJlho%2B9ELBdrER2yaGgs8FF8I0c%2BMgopcQLfhSbmuoE2wckh2deQyNdMjkt6EWsmKxj2ntYG0e6guVeUt9%2BUNQoVzXEBHBzCX4SMsEUHZuV%2BpEWGpEGJgNayG7KOApnNukdlkztgKBPh4B7YRW1jFcgzJC7XMMDt4TBDx5Zc0KntSLRUKCMhSAyP6S6VqUAEStHM0S9JWWABydP3NC1YvQOnX3ll83Lp4m5sYuvgph40Rn2%2BLHvCSrG63eHDJq7ZMw1rrrwwY6pgHQyYAIUtGCOcO8Uns9xcYcU4a0Orr4zf1wHfCzB1rkkl81BL2LIeVGBjoErQ%2BLG4tRJDxw1v8LOQKq0Bu0KXHXlIlQBfS8Vqhn9if%2FCl8fP8h7BfV89U1c83nUOZqq4JVd0dSqKp%2BGOcVQYYAAgRPy6ORg%2FdmO5zMulEAO6CA4ADJPzHmCV0T6kQlfoMJqOTunilc0Dz1j4gLoZoYDaXd1yz0Vighq&X-Amz-Signature=f1a5e3941062c4e81ee3417fcfbfba10028e2fe1a2dbbfe57b536c99c402685e&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)

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

![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/f17d4bd3-48a9-41d6-aceb-c79844361606/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB466YWFV4FRH%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005444Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJGMEQCIGPKT4xBWQOsvcusq1mdTjOKon0jKWWRXI4oXhp8wuVHAiAPtZomf3qcTl4iCboxWwadA1UY3GAtDnl641ZsfI92IiqIBAiZ%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F8BEAAaDDYzNzQyMzE4MzgwNSIMk8O2Nl0mSzR%2Bw5TfKtwDAFIOYF72QTLmSZwv1nSlC4mcnN5I2%2F%2F%2FRqkzM4oMaQipwqrm8yBaUeMvzK2FcYQbmc9kK6ONg1dgoZj3pz4jTwJtIvN9YV9x9ZaDtbfncR%2B559%2FbqAjfV10PZvELymtctl2NaeAyvbaiGPkxQnkgjydQP%2BDFOrXj4My0sTGo2peVcodFksB70G0lbHr18iE3z8AaH4gJWF03Wbp6PzpgvfaE3Jg2fl%2BNRELvK8Y60p8grzSyk9brOCnMTijKt%2FU65lpPs9dWWGwxBkKOYtPSMbDk%2FJEQUkPBPERgBkBp%2BHNa%2BjhW3i9wbmev7b98tGkVk7HwGOFW4I8eNym7WKsFbRJokqUP%2BEUvhLgBpPHXNmDTz%2BH32T7AubhwLw%2BJlho%2B9ELBdrER2yaGgs8FF8I0c%2BMgopcQLfhSbmuoE2wckh2deQyNdMjkt6EWsmKxj2ntYG0e6guVeUt9%2BUNQoVzXEBHBzCX4SMsEUHZuV%2BpEWGpEGJgNayG7KOApnNukdlkztgKBPh4B7YRW1jFcgzJC7XMMDt4TBDx5Zc0KntSLRUKCMhSAyP6S6VqUAEStHM0S9JWWABydP3NC1YvQOnX3ll83Lp4m5sYuvgph40Rn2%2BLHvCSrG63eHDJq7ZMw1rrrwwY6pgHQyYAIUtGCOcO8Uns9xcYcU4a0Orr4zf1wHfCzB1rkkl81BL2LIeVGBjoErQ%2BLG4tRJDxw1v8LOQKq0Bu0KXHXlIlQBfS8Vqhn9if%2FCl8fP8h7BfV89U1c83nUOZqq4JVd0dSqKp%2BGOcVQYYAAgRPy6ORg%2FdmO5zMulEAO6CA4ADJPzHmCV0T6kQlfoMJqOTunilc0Dz1j4gLoZoYDaXd1yz0Vighq&X-Amz-Signature=f8c80ea761cd0cc6007c903dfbe1dd18276c0da8a8b30991fd917a2168da8a6e&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)

### **关键3：数据流处理难题**

> 💡 **事件驱动设计，整个图的执行过程中的每一个步骤，甚至每一帧数据都被视作独立的事件**

**数据流处理表：**

| 结构 | 示例 | 
| --- | --- | 
| 流复制 | 复制成多个流，因为数据流一旦被消费可能就无法再被重复使用 | 
| 流合并 | 将多个独立的数据流合并成一个单一的数据流 | 
| 流拼接 | 将流式数据中的各帧数据组合成一个完整的数据集 | 
| 流装箱 | 转换回流式数据（由多帧数据组成一个序列） | 

![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/32ef152c-6e51-4aa7-b4dd-fe9ac5f87d4b/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB466YWFV4FRH%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005444Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJGMEQCIGPKT4xBWQOsvcusq1mdTjOKon0jKWWRXI4oXhp8wuVHAiAPtZomf3qcTl4iCboxWwadA1UY3GAtDnl641ZsfI92IiqIBAiZ%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F8BEAAaDDYzNzQyMzE4MzgwNSIMk8O2Nl0mSzR%2Bw5TfKtwDAFIOYF72QTLmSZwv1nSlC4mcnN5I2%2F%2F%2FRqkzM4oMaQipwqrm8yBaUeMvzK2FcYQbmc9kK6ONg1dgoZj3pz4jTwJtIvN9YV9x9ZaDtbfncR%2B559%2FbqAjfV10PZvELymtctl2NaeAyvbaiGPkxQnkgjydQP%2BDFOrXj4My0sTGo2peVcodFksB70G0lbHr18iE3z8AaH4gJWF03Wbp6PzpgvfaE3Jg2fl%2BNRELvK8Y60p8grzSyk9brOCnMTijKt%2FU65lpPs9dWWGwxBkKOYtPSMbDk%2FJEQUkPBPERgBkBp%2BHNa%2BjhW3i9wbmev7b98tGkVk7HwGOFW4I8eNym7WKsFbRJokqUP%2BEUvhLgBpPHXNmDTz%2BH32T7AubhwLw%2BJlho%2B9ELBdrER2yaGgs8FF8I0c%2BMgopcQLfhSbmuoE2wckh2deQyNdMjkt6EWsmKxj2ntYG0e6guVeUt9%2BUNQoVzXEBHBzCX4SMsEUHZuV%2BpEWGpEGJgNayG7KOApnNukdlkztgKBPh4B7YRW1jFcgzJC7XMMDt4TBDx5Zc0KntSLRUKCMhSAyP6S6VqUAEStHM0S9JWWABydP3NC1YvQOnX3ll83Lp4m5sYuvgph40Rn2%2BLHvCSrG63eHDJq7ZMw1rrrwwY6pgHQyYAIUtGCOcO8Uns9xcYcU4a0Orr4zf1wHfCzB1rkkl81BL2LIeVGBjoErQ%2BLG4tRJDxw1v8LOQKq0Bu0KXHXlIlQBfS8Vqhn9if%2FCl8fP8h7BfV89U1c83nUOZqq4JVd0dSqKp%2BGOcVQYYAAgRPy6ORg%2FdmO5zMulEAO6CA4ADJPzHmCV0T6kQlfoMJqOTunilc0Dz1j4gLoZoYDaXd1yz0Vighq&X-Amz-Signature=28d3f12d4f8657cac5ac5bb2fc2e102b0de16801dd1cd51639d329e0743940db&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)



> 附：还有Go的类型安全优势、AI大模型的统一集成等

---



# 二、Eino的整体结构

> 💡 **开源社区优秀的AI应用开发框架**
> - [https://github.com/langchain-ai/langchain](https://github.com/langchain-ai/langchain) & [https://github.com/langchain-ai/langgraph](https://github.com/langchain-ai/langgraph)（python/js）
>     - 文档：[https://python.langchain.com/docs/introduction/](https://python.langchain.com/docs/introduction/)
>     - 官网：[https://www.langchain.com/](https://www.langchain.com/) 字节自己也在用
> - [https://github.com/run-llama/llama_index](https://github.com/run-llama/llama_index)（python/js）
>     - 官网：[https://www.llamaindex.ai/](https://www.llamaindex.ai/) 
>     - 文档：[https://docs.llamaindex.ai/en/stable/](https://docs.llamaindex.ai/en/stable/)



![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/619bb96a-a3af-40a4-ba76-c8fd81ea78dd/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB466YWFV4FRH%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005444Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJGMEQCIGPKT4xBWQOsvcusq1mdTjOKon0jKWWRXI4oXhp8wuVHAiAPtZomf3qcTl4iCboxWwadA1UY3GAtDnl641ZsfI92IiqIBAiZ%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F8BEAAaDDYzNzQyMzE4MzgwNSIMk8O2Nl0mSzR%2Bw5TfKtwDAFIOYF72QTLmSZwv1nSlC4mcnN5I2%2F%2F%2FRqkzM4oMaQipwqrm8yBaUeMvzK2FcYQbmc9kK6ONg1dgoZj3pz4jTwJtIvN9YV9x9ZaDtbfncR%2B559%2FbqAjfV10PZvELymtctl2NaeAyvbaiGPkxQnkgjydQP%2BDFOrXj4My0sTGo2peVcodFksB70G0lbHr18iE3z8AaH4gJWF03Wbp6PzpgvfaE3Jg2fl%2BNRELvK8Y60p8grzSyk9brOCnMTijKt%2FU65lpPs9dWWGwxBkKOYtPSMbDk%2FJEQUkPBPERgBkBp%2BHNa%2BjhW3i9wbmev7b98tGkVk7HwGOFW4I8eNym7WKsFbRJokqUP%2BEUvhLgBpPHXNmDTz%2BH32T7AubhwLw%2BJlho%2B9ELBdrER2yaGgs8FF8I0c%2BMgopcQLfhSbmuoE2wckh2deQyNdMjkt6EWsmKxj2ntYG0e6guVeUt9%2BUNQoVzXEBHBzCX4SMsEUHZuV%2BpEWGpEGJgNayG7KOApnNukdlkztgKBPh4B7YRW1jFcgzJC7XMMDt4TBDx5Zc0KntSLRUKCMhSAyP6S6VqUAEStHM0S9JWWABydP3NC1YvQOnX3ll83Lp4m5sYuvgph40Rn2%2BLHvCSrG63eHDJq7ZMw1rrrwwY6pgHQyYAIUtGCOcO8Uns9xcYcU4a0Orr4zf1wHfCzB1rkkl81BL2LIeVGBjoErQ%2BLG4tRJDxw1v8LOQKq0Bu0KXHXlIlQBfS8Vqhn9if%2FCl8fP8h7BfV89U1c83nUOZqq4JVd0dSqKp%2BGOcVQYYAAgRPy6ORg%2FdmO5zMulEAO6CA4ADJPzHmCV0T6kQlfoMJqOTunilc0Dz1j4gLoZoYDaXd1yz0Vighq&X-Amz-Signature=c974afc42b1e0991d8656eb70371ea49d229ebda4dbd674d038e8fe98113c35f&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)

| 代码仓库 | **Feature** | 
| --- | --- | 
| **核心库**[https://github.com/cloudwego/eino](https://github.com/cloudwego/eino) | **1. 核心Schema**：定义和抽象AI应用领域的领域模型的结构体。如流的读取和写入<br/>**2. Components**：8种组件类型的抽象。如各自的输入输出规范、流式处理范式和回调事件信息<br/>**3. Compose**：核心的编排能力部分。包含编排的基本组成元素和能力。<br/>&nbsp;&nbsp;a). **引擎层**：pregel和dag。区别在于节点间触发下一个节点的方式<br/>&nbsp;&nbsp;b). **API层**：<br/>&nbsp;&nbsp;&nbsp;&nbsp;◦ **Chain**：简单的链式有向图，只能向前推进<br/>&nbsp;&nbsp;&nbsp;&nbsp;◦ **Graph**：循环或非循环有向图。功能强大且灵活<br/>&nbsp;&nbsp;&nbsp;&nbsp;◦ **Workflow**：这是一种新型结构，能够**解耦数据流和控制流**，允许灵活配置数据流节点间的数据映射关系<br/>**4. Flow**：最顶层包含一些预置的编排产物，如 ReAct Agent 和 Multi Agent 等，这些都是基于compose builder能力构建的 | 
| **组件的具体实现库** [https://github.com/cloudwego/eino-ext](https://github.com/cloudwego/eino-ext) | 包括回调事件处理器的实现及可视化开发/调试工具，如IDE插件、[Eino Devops](https://github.com/cloudwego/eino-ext/tree/main/devops) | 
| **示例库** [https://github.com/cloudwego/eino-examples](https://github.com/cloudwego/eino-examples) |   | 
| **官方文档** | [https://www.cloudwego.io/zh/docs/eino/overview/bytedance_eino_practice/](https://www.cloudwego.io/zh/docs/eino/overview/bytedance_eino_practice/) | 





