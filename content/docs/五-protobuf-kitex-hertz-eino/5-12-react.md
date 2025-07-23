---
title: "5.12 ReAct"
date: 2025-07-17T12:52:00Z
draft: false
weight: 5012
---

# 5.12 ReAct

[https://react-lm.github.io/](https://react-lm.github.io/)

# **ReAct: Synergizing Reasoning and Acting in Language ModelsReAct: 结合语言模型的推理和行动**

[Shunyu Yao](https://ysymyth.github.io/), [Jeffrey Zhao](http://descrip.github.io/), [Dian Yu](https://diandyu.github.io/), [Nan Du](https://research.google/people/104844/), [Izhak Shafran](https://research.google/people/IzhakShafran/), [Karthik Narasimhan](https://www.cs.princeton.edu/~karthikn/), [Yuan Cao](https://research.google/people/YuanCao/)姚舜宇, 乔希· Zhao, 于典, 但_nan, 艾扎克·沙弗兰, 卡尔希克·纳拉辛汉, 曹元

[**[Paper]**](https://arxiv.org/abs/2210.03629)

[**[Code]**](https://github.com/ysymyth/ReAct)

[**[Blogpost]**](https://ai.googleblog.com/2022/11/react-synergizing-reasoning-and-acting.html)

[**[BibTex]**](https://react-lm.github.io/files/bib.txt)

[论文][代码][博客][BibTeX]

![](/images/23324637-29b5-8050-a76f-d9e5c73ca2f7/image_23324637-29b5-80b0-95cc-ecb84e3bee5f.jpg)

Language models are getting better at reasoning (e.g. chain-of-thought prompting) and acting (e.g. WebGPT, SayCan, ACT-1), but these two directions have remained separate.语言模型在推理（例如：链式思考提示）和行动（例如：WebGPT, SayCan, ACT-1）方面变得越来越好，但这两个方向一直保持独立。**ReAct asks, what if these two fundamental capabilities are combined?ReAct 提问，如果将这两种基本能力结合起来会怎样？**

## **Abstract 摘要**

While large language models (LLMs) have demonstrated impressive capabilities across tasks in language understanding and interactive decision making, their abilities for reasoning (e.g. chain-of-thought prompting) and acting (e.g. action plan generation) have primarily been studied as separate topics. In this paper, we explore the use of LLMs to generate both reasoning traces and task-specific actions in an interleaved manner, allowing for greater synergy between the two: reasoning traces help the model induce, track, and update action plans as well as handle exceptions, while actions allow it to interface with external sources, such as knowledge bases or environments, to gather additional information. We apply our approach, named ReAct, to a diverse set of language and decision making tasks and demonstrate its effectiveness over state-of-the-art baselines, as well as improved human interpretability and trustworthiness over methods without reasoning or acting components. Concretely, on question answering (HotpotQA) and fact verification (Fever), ReAct overcomes issues of hallucination and error propagation prevalent in chain-of-thought reasoning by interacting with a simple Wikipedia API, and generates human-like task-solving trajectories that are more interpretable than baselines without reasoning traces. On two interactive decision making benchmarks (ALFWorld and WebShop), ReAct outperforms imitation and reinforcement learning methods by an absolute success rate of 34% and 10% respectively, while being prompted with only one or two in-context examples.尽管大规模语言模型（LLMs）在语言理解和交互决策任务中展现了令人印象深刻的性能，但它们的推理能力（例如：链式思考提示）和执行能力（例如：行动方案生成）主要被分别研究。在本文中，我们探索了使用 LLMs 以交错的方式生成推理痕迹和任务特定行动，从而在两者之间实现更大的协同作用：推理痕迹帮助模型生成、跟踪和更新行动方案，以及处理异常情况，而行动则允许模型与外部来源（如知识库或环境）进行交互，以获取额外信息。我们将我们的方法命名为 ReAct，并将其应用于多种语言和决策任务，展示了其在与最先进的基线相比的有效性，以及在没有推理或执行组件的方法中，提高了人类的可解释性和可信度。 具体而言，在问答（HotpotQA）和事实验证（Fever）任务中，ReAct 通过与简单的维基百科 API 交互，克服了链式推理中普遍存在的幻觉和错误传播问题，生成了比没有推理痕迹的基线更具可解释性的类人类任务解决轨迹。在两个交互式决策基准测试（ALFWorld 和 WebShop）中，ReAct 的绝对成功率分别比模仿学习和强化学习方法高出 34%和 10%，仅需一个或两个上下文示例即可。

## **ReAct Prompting ReAct 提示**

A ReAct prompt consists of few-shot task-solving trajectories, with human-written text reasoning traces and actions, as well as environment observations in response to actions (see examples in paper appendix!)ReAct 提示由少量示例任务解决轨迹组成，包含人类撰写的文本推理痕迹和动作，以及对动作的环境观察（见论文附录中的示例！）ReAct prompting is intuitive and flexible to design, and achieves state-of-the-art few-shot performances across a variety of tasks, from question answering to online shopping!ReAct 提示设计直观且灵活，能够在各种任务中实现最先进的少量示例性能，从问答到在线购物！

![](/images/23324637-29b5-8050-a76f-d9e5c73ca2f7/image_23324637-29b5-8052-a652-d6870adf7a1a.jpg)

### **HotpotQA Example HotpotQA 示例**

The reason-only baseline (i.e. chain-of-thought) suffers from misinformation (in red) as it is not grounded to external environments to obtain and update knowledge, and has to rely on limited internal knowledge.仅推理基线（即思维链）会受到误导信息（用红色标注）的影响，因为它不依赖外部环境来获取和更新知识，只能依靠有限的内部知识。The act-only baseline suffers from the lack of reasoning, unable to synthesize the final answer despite having the same actions and observation as ReAct in this case.仅执行基线缺乏推理能力，即使在这种情况下拥有与 ReAct 相同的操作和观察，也无法综合得出最终答案。In contrast, ReAct solves the task with a interpretable and factual trajectory.相比之下，ReAct 通过可解释且符合事实的轨迹解决了任务。

![](/images/23324637-29b5-8050-a76f-d9e5c73ca2f7/image_23324637-29b5-802a-9123-e3b95f621a2a.jpg)

### **ALFWorld Example ALFWorld 示例**

For decision making tasks, we design human trajectories with sparse reasoning traces, letting the LM decide when to think vs. act.对于决策任务，我们设计了包含稀疏推理痕迹的人类轨迹，让语言模型决定何时思考而非行动。ReAct isn't perfect --- below is a failure example on ALFWorld. However, ReAct format allows easy human inspection and behavior correction by changing a couple of model thoughts, an exciting novel approach to human alignment!ReAct 并不完美——下面是一个在 ALFWorld 中的失败示例。然而，ReAct 格式允许通过更改几个模型的想法进行简单的手动检查和行为修正，这是一种令人兴奋的新方法，用于实现与人类的对齐！

![](/images/23324637-29b5-8050-a76f-d9e5c73ca2f7/image_23324637-29b5-80c7-bfc4-f74d7520dc16.jpg)

## **ReAct Finetuning: Initial ResultsReAct 微调：初步结果**

Prompting has limited context window and learning support. Initial finetuning results on HotpotQA using ReAct prompting trajectories suggest: (1) ReAct is the best fintuning format across model sizes; (2) ReAct finetuned smaller models outperform prompted larger models!提示的上下文窗口和学习支持有限。使用 ReAct 提示轨迹对 HotpotQA 进行初始微调结果表明：（1）ReAct 是各种模型大小的最佳微调格式；（2）微调的小模型优于提示的大模型！

![](/images/23324637-29b5-8050-a76f-d9e5c73ca2f7/image_23324637-29b5-80fb-be2c-d2fc6bc3b14b.jpg)

