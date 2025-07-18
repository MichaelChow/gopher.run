---
title: "透过Eino看AI应用开发：基本结构"
date: 2025-05-19T04:49:00Z
draft: false
---

# 透过Eino看AI应用开发：基本结构

# 一、AI应用的本质

> 🎯 核心洞察：AI应用的本质是围绕大模型的信息流：Input（原始信息） → Process（大模型） → Output(最终信息)

![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/1c73908b-803e-4ce8-84c1-f2e3b9634e47/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4667NUBS54H%2F20250718%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250718T170132Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEHgaCXVzLXdlc3QtMiJIMEYCIQCBlQeQItKlNtf0x1hUEheREOUYmL6ZmY5BlIB6nQql3QIhAMdI%2Fomb1bxAr4z8DvGPKk1X2pa3XGNqLqC7mXrtmFAjKogECJH%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1IgxXSevChp2Qw%2F6yvRwq3ANzyibS6UqpzK5lpZI3SHtzrel9ORXwuT4mSLOSfLXNS2FUO%2F1iAKsZvflhSb7ywQJ4J2s4MSFkUkizK8apjuVYXddobd40yeKvoRjsjssF5CxTUhaVQJw7vUZCjHaAsLgqkawpYzBQMNexh6BE85cT%2FAJdePP4PpnISNbdzlLKVTHBBX6lH4ukVDRkr1gHrYfBQ4daknNu7EqHs7%2Fu7Dq0urgTUfmB1GYnOfWzFcoNFKxikb4jhO0S5CdinzqDeg%2FZxV1yZZqM5lGIoLVsuuFIQ%2FSbJDRCTzu4izKoYU%2FJEjIXJZ1AuXfsKsVULE1NuKX%2BKb3SRP7dssW348S75rQgjQjHB5C2ykBnxxckpCnbeF4AN3Z2DzyQUZhpxpH2nBL0Fy%2FMTc9XqAeN%2FpZr8kROSBbzmkqvm1V%2FFEffwJMzC8DITZGzYhwoEeS%2FAY%2BbG7V5xQC8oZIVu%2FXJLe8GiDI6Q3azxjlox1%2FZJMlz7XsNgSnchUxtSs5fq5e4e9KSuxdjzXgXnDoxoLB9SwpL7iqzoY3ssIYXePCYgNjMYc%2FR%2FLgz5LUP7LSCVl77OubCkrZavhihZuPBR6Q5BMVgEzp8XiMwx5sKlD79y73Fi%2F8mjW5ZPnXnDGl5KzDEwDCZzunDBjqkAQh6Nged3kO0ISjSidiJowP7tC%2BLWrs%2BYtiboIAXABkbSdoxQAG4lJCW4FOVN7EZst6FCY4eW3jm8O3rs3mW7UE8rBPlGDPh1XIT0tfW2UpwiiK76gbL57%2F%2FJKNZVheW5ClCHOwX25QCrEh0nzXar35ter1XQqdnTVAdvJigsUrAsM4OUySVHrUWVhc%2FqBFqlXAbb5S3qboTfDCyGjfuYe0aUrbi&X-Amz-Signature=f00710323622da6c559d461ce04c7aac2f12ca43fb0bc9ccfe04d143162ac954&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)

### **关键1：组件及其实现**

**组件定义表：**

| **组件(Interface)** | **组件实现(实例)** | 分工 | 
| --- | --- | --- | 
| **ChatTemplate**Prompt模板 |   | 负责信息的格式化与增强处理 | 
| **Retriever**知识库 | Redis、ES… | 负责从知识库中召回信息，进行必要的信息检索 | 
| **ChatModel**大模型 | OpenAI、Claude… | 负责推理和信息生成 | 
| **Tool**工具 | MCP、search… | 负责进一步的信息处理 | 

![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/e43e0c79-db4a-45dd-870b-cb7af123867b/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4667NUBS54H%2F20250718%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250718T170132Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEHgaCXVzLXdlc3QtMiJIMEYCIQCBlQeQItKlNtf0x1hUEheREOUYmL6ZmY5BlIB6nQql3QIhAMdI%2Fomb1bxAr4z8DvGPKk1X2pa3XGNqLqC7mXrtmFAjKogECJH%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1IgxXSevChp2Qw%2F6yvRwq3ANzyibS6UqpzK5lpZI3SHtzrel9ORXwuT4mSLOSfLXNS2FUO%2F1iAKsZvflhSb7ywQJ4J2s4MSFkUkizK8apjuVYXddobd40yeKvoRjsjssF5CxTUhaVQJw7vUZCjHaAsLgqkawpYzBQMNexh6BE85cT%2FAJdePP4PpnISNbdzlLKVTHBBX6lH4ukVDRkr1gHrYfBQ4daknNu7EqHs7%2Fu7Dq0urgTUfmB1GYnOfWzFcoNFKxikb4jhO0S5CdinzqDeg%2FZxV1yZZqM5lGIoLVsuuFIQ%2FSbJDRCTzu4izKoYU%2FJEjIXJZ1AuXfsKsVULE1NuKX%2BKb3SRP7dssW348S75rQgjQjHB5C2ykBnxxckpCnbeF4AN3Z2DzyQUZhpxpH2nBL0Fy%2FMTc9XqAeN%2FpZr8kROSBbzmkqvm1V%2FFEffwJMzC8DITZGzYhwoEeS%2FAY%2BbG7V5xQC8oZIVu%2FXJLe8GiDI6Q3azxjlox1%2FZJMlz7XsNgSnchUxtSs5fq5e4e9KSuxdjzXgXnDoxoLB9SwpL7iqzoY3ssIYXePCYgNjMYc%2FR%2FLgz5LUP7LSCVl77OubCkrZavhihZuPBR6Q5BMVgEzp8XiMwx5sKlD79y73Fi%2F8mjW5ZPnXnDGl5KzDEwDCZzunDBjqkAQh6Nged3kO0ISjSidiJowP7tC%2BLWrs%2BYtiboIAXABkbSdoxQAG4lJCW4FOVN7EZst6FCY4eW3jm8O3rs3mW7UE8rBPlGDPh1XIT0tfW2UpwiiK76gbL57%2F%2FJKNZVheW5ClCHOwX25QCrEh0nzXar35ter1XQqdnTVAdvJigsUrAsM4OUySVHrUWVhc%2FqBFqlXAbb5S3qboTfDCyGjfuYe0aUrbi&X-Amz-Signature=4452a75a8ff1f66a3631a6b823652f3659082d5ce7dafcebd7f5c88b47807dea&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)

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

![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/f17d4bd3-48a9-41d6-aceb-c79844361606/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4667NUBS54H%2F20250718%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250718T170132Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEHgaCXVzLXdlc3QtMiJIMEYCIQCBlQeQItKlNtf0x1hUEheREOUYmL6ZmY5BlIB6nQql3QIhAMdI%2Fomb1bxAr4z8DvGPKk1X2pa3XGNqLqC7mXrtmFAjKogECJH%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1IgxXSevChp2Qw%2F6yvRwq3ANzyibS6UqpzK5lpZI3SHtzrel9ORXwuT4mSLOSfLXNS2FUO%2F1iAKsZvflhSb7ywQJ4J2s4MSFkUkizK8apjuVYXddobd40yeKvoRjsjssF5CxTUhaVQJw7vUZCjHaAsLgqkawpYzBQMNexh6BE85cT%2FAJdePP4PpnISNbdzlLKVTHBBX6lH4ukVDRkr1gHrYfBQ4daknNu7EqHs7%2Fu7Dq0urgTUfmB1GYnOfWzFcoNFKxikb4jhO0S5CdinzqDeg%2FZxV1yZZqM5lGIoLVsuuFIQ%2FSbJDRCTzu4izKoYU%2FJEjIXJZ1AuXfsKsVULE1NuKX%2BKb3SRP7dssW348S75rQgjQjHB5C2ykBnxxckpCnbeF4AN3Z2DzyQUZhpxpH2nBL0Fy%2FMTc9XqAeN%2FpZr8kROSBbzmkqvm1V%2FFEffwJMzC8DITZGzYhwoEeS%2FAY%2BbG7V5xQC8oZIVu%2FXJLe8GiDI6Q3azxjlox1%2FZJMlz7XsNgSnchUxtSs5fq5e4e9KSuxdjzXgXnDoxoLB9SwpL7iqzoY3ssIYXePCYgNjMYc%2FR%2FLgz5LUP7LSCVl77OubCkrZavhihZuPBR6Q5BMVgEzp8XiMwx5sKlD79y73Fi%2F8mjW5ZPnXnDGl5KzDEwDCZzunDBjqkAQh6Nged3kO0ISjSidiJowP7tC%2BLWrs%2BYtiboIAXABkbSdoxQAG4lJCW4FOVN7EZst6FCY4eW3jm8O3rs3mW7UE8rBPlGDPh1XIT0tfW2UpwiiK76gbL57%2F%2FJKNZVheW5ClCHOwX25QCrEh0nzXar35ter1XQqdnTVAdvJigsUrAsM4OUySVHrUWVhc%2FqBFqlXAbb5S3qboTfDCyGjfuYe0aUrbi&X-Amz-Signature=02155d20a263c0349fdf007c363d935db3c00a5f541b21e8d247607783478d96&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)

### **关键3：数据流处理难题**

> 💡 **事件驱动设计，整个图的执行过程中的每一个步骤，甚至每一帧数据都被视作独立的事件**

**数据流处理表：**

| 结构 | 示例 | 
| --- | --- | 
| 流复制 | 复制成多个流，因为数据流一旦被消费可能就无法再被重复使用 | 
| 流合并 | 将多个独立的数据流合并成一个单一的数据流 | 
| 流拼接 | 将流式数据中的各帧数据组合成一个完整的数据集 | 
| 流装箱 | 转换回流式数据（由多帧数据组成一个序列） | 

![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/32ef152c-6e51-4aa7-b4dd-fe9ac5f87d4b/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4667NUBS54H%2F20250718%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250718T170132Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEHgaCXVzLXdlc3QtMiJIMEYCIQCBlQeQItKlNtf0x1hUEheREOUYmL6ZmY5BlIB6nQql3QIhAMdI%2Fomb1bxAr4z8DvGPKk1X2pa3XGNqLqC7mXrtmFAjKogECJH%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1IgxXSevChp2Qw%2F6yvRwq3ANzyibS6UqpzK5lpZI3SHtzrel9ORXwuT4mSLOSfLXNS2FUO%2F1iAKsZvflhSb7ywQJ4J2s4MSFkUkizK8apjuVYXddobd40yeKvoRjsjssF5CxTUhaVQJw7vUZCjHaAsLgqkawpYzBQMNexh6BE85cT%2FAJdePP4PpnISNbdzlLKVTHBBX6lH4ukVDRkr1gHrYfBQ4daknNu7EqHs7%2Fu7Dq0urgTUfmB1GYnOfWzFcoNFKxikb4jhO0S5CdinzqDeg%2FZxV1yZZqM5lGIoLVsuuFIQ%2FSbJDRCTzu4izKoYU%2FJEjIXJZ1AuXfsKsVULE1NuKX%2BKb3SRP7dssW348S75rQgjQjHB5C2ykBnxxckpCnbeF4AN3Z2DzyQUZhpxpH2nBL0Fy%2FMTc9XqAeN%2FpZr8kROSBbzmkqvm1V%2FFEffwJMzC8DITZGzYhwoEeS%2FAY%2BbG7V5xQC8oZIVu%2FXJLe8GiDI6Q3azxjlox1%2FZJMlz7XsNgSnchUxtSs5fq5e4e9KSuxdjzXgXnDoxoLB9SwpL7iqzoY3ssIYXePCYgNjMYc%2FR%2FLgz5LUP7LSCVl77OubCkrZavhihZuPBR6Q5BMVgEzp8XiMwx5sKlD79y73Fi%2F8mjW5ZPnXnDGl5KzDEwDCZzunDBjqkAQh6Nged3kO0ISjSidiJowP7tC%2BLWrs%2BYtiboIAXABkbSdoxQAG4lJCW4FOVN7EZst6FCY4eW3jm8O3rs3mW7UE8rBPlGDPh1XIT0tfW2UpwiiK76gbL57%2F%2FJKNZVheW5ClCHOwX25QCrEh0nzXar35ter1XQqdnTVAdvJigsUrAsM4OUySVHrUWVhc%2FqBFqlXAbb5S3qboTfDCyGjfuYe0aUrbi&X-Amz-Signature=a410246037710cbb179cebd2d8f0c3d017e3bbc70383a7dd4b779308000349a4&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)



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



![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/619bb96a-a3af-40a4-ba76-c8fd81ea78dd/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4667NUBS54H%2F20250718%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250718T170132Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEHgaCXVzLXdlc3QtMiJIMEYCIQCBlQeQItKlNtf0x1hUEheREOUYmL6ZmY5BlIB6nQql3QIhAMdI%2Fomb1bxAr4z8DvGPKk1X2pa3XGNqLqC7mXrtmFAjKogECJH%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1IgxXSevChp2Qw%2F6yvRwq3ANzyibS6UqpzK5lpZI3SHtzrel9ORXwuT4mSLOSfLXNS2FUO%2F1iAKsZvflhSb7ywQJ4J2s4MSFkUkizK8apjuVYXddobd40yeKvoRjsjssF5CxTUhaVQJw7vUZCjHaAsLgqkawpYzBQMNexh6BE85cT%2FAJdePP4PpnISNbdzlLKVTHBBX6lH4ukVDRkr1gHrYfBQ4daknNu7EqHs7%2Fu7Dq0urgTUfmB1GYnOfWzFcoNFKxikb4jhO0S5CdinzqDeg%2FZxV1yZZqM5lGIoLVsuuFIQ%2FSbJDRCTzu4izKoYU%2FJEjIXJZ1AuXfsKsVULE1NuKX%2BKb3SRP7dssW348S75rQgjQjHB5C2ykBnxxckpCnbeF4AN3Z2DzyQUZhpxpH2nBL0Fy%2FMTc9XqAeN%2FpZr8kROSBbzmkqvm1V%2FFEffwJMzC8DITZGzYhwoEeS%2FAY%2BbG7V5xQC8oZIVu%2FXJLe8GiDI6Q3azxjlox1%2FZJMlz7XsNgSnchUxtSs5fq5e4e9KSuxdjzXgXnDoxoLB9SwpL7iqzoY3ssIYXePCYgNjMYc%2FR%2FLgz5LUP7LSCVl77OubCkrZavhihZuPBR6Q5BMVgEzp8XiMwx5sKlD79y73Fi%2F8mjW5ZPnXnDGl5KzDEwDCZzunDBjqkAQh6Nged3kO0ISjSidiJowP7tC%2BLWrs%2BYtiboIAXABkbSdoxQAG4lJCW4FOVN7EZst6FCY4eW3jm8O3rs3mW7UE8rBPlGDPh1XIT0tfW2UpwiiK76gbL57%2F%2FJKNZVheW5ClCHOwX25QCrEh0nzXar35ter1XQqdnTVAdvJigsUrAsM4OUySVHrUWVhc%2FqBFqlXAbb5S3qboTfDCyGjfuYe0aUrbi&X-Amz-Signature=d8ce9be2bedbd1b76a768d4007d7d3186441e23c0c3d1ea513ce00a3fb7e34ce&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)

| 代码仓库 | **Feature** | 
| --- | --- | 
| **核心库**[https://github.com/cloudwego/eino](https://github.com/cloudwego/eino) | **1. 核心Schema**：定义和抽象AI应用领域的领域模型的结构体。如流的读取和写入<br/>**2. Components**：8种组件类型的抽象。如各自的输入输出规范、流式处理范式和回调事件信息<br/>**3. Compose**：核心的编排能力部分。包含编排的基本组成元素和能力。<br/>&nbsp;&nbsp;a). **引擎层**：pregel和dag。区别在于节点间触发下一个节点的方式<br/>&nbsp;&nbsp;b). **API层**：<br/>&nbsp;&nbsp;&nbsp;&nbsp;◦ **Chain**：简单的链式有向图，只能向前推进<br/>&nbsp;&nbsp;&nbsp;&nbsp;◦ **Graph**：循环或非循环有向图。功能强大且灵活<br/>&nbsp;&nbsp;&nbsp;&nbsp;◦ **Workflow**：这是一种新型结构，能够**解耦数据流和控制流**，允许灵活配置数据流节点间的数据映射关系<br/>**4. Flow**：最顶层包含一些预置的编排产物，如 ReAct Agent 和 Multi Agent 等，这些都是基于compose builder能力构建的 | 
| **组件的具体实现库** [https://github.com/cloudwego/eino-ext](https://github.com/cloudwego/eino-ext) | 包括回调事件处理器的实现及可视化开发/调试工具，如IDE插件、[Eino Devops](https://github.com/cloudwego/eino-ext/tree/main/devops) | 
| **示例库** [https://github.com/cloudwego/eino-examples](https://github.com/cloudwego/eino-examples) |   | 
| **官方文档** | [https://www.cloudwego.io/zh/docs/eino/overview/bytedance_eino_practice/](https://www.cloudwego.io/zh/docs/eino/overview/bytedance_eino_practice/) | 





