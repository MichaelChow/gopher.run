---
title: "text/template"
date: 2025-08-05T00:51:00Z
draft: false
---

# text/template



> [https://pkg.go.dev/text/template](https://pkg.go.dev/text/template)



package template implements data-driven templates for generating textual output.

text/template 是 Go 标准库提供的数据驱动的文本模板引擎，用于生成动态文本内容。



模板是通过将其应用于数据结构来执行的。模板中的注释引用数据结构中的元素（通常是结构体的字段或映射的键），以控制执行并导出要显示的值。模板的执行会遍历结构，并将光标（用点'.'表示，称为"dot"）设置为执行过程中当前结构位置的值。



模板的输入文本是任何格式的 UTF-8 编码文本。"Actions"（数据评估或控制结构）由"{{"and"}}"分隔；所有动作外的文本都会原样复制到输出中。

- 默认情况下，当模板执行时，所有在操作符之间的文本都会原封不动地复制。
- 修剪空白字符：如果操作符的左定界符（默认为"{{"）紧跟着一个减号和空格，那么会从紧邻的前一个文本中删除所有尾随的空格。如果右定界符（"}}"）前面有空格和减号，会从紧邻的下一个文本中删除所有开头的空格。
- 空白字符的定义与 Go 中的定义相同：空格、水平制表符、回车和换行。


### **基本结构**

**1. Template 对象**

```go
type Template struct {
    name string
    *parse.Tree
    *common
    leftDelim  string
    rightDelim string
}
```

**2. 模板语法**

- **分隔符**: {{ 和 }}
- **变量**: {{.VariableName}}
- **管道**: {{.Value | function}}
- **控制结构**: {{if}}, {{range}}, {{with}}


为什么设计成{{.v}}，而不是更简洁的{v}? Go 语言"明确优于简洁"的设计哲学

```go
// 问题示例
template := `
用户信息:
姓名: {name}
邮箱: {email}
地址: {address}
价格: {price} 元
`

// 问题：
// 1. 普通文本中的大括号会被误解析
// 2. 难以区分变量和文本
// 3. 解析器复杂度增加
```

```go
// 优势示例
template := `
用户信息:
姓名: {{.name}}
邮箱: {{.email}}
地址: {{.address}}
价格: {{.price}} 元
`

// 优势：
// 1. 明确的语法边界
// 2. 易于解析和识别
// 3. 与主流模板引擎一致
// 4. 支持复杂语法扩展
```

**主流模板引擎对比：**

| 模板引擎 | 语言 | 变量语法 | 条件语法 | 循环语法 | 函数调用 | 
| --- | --- | --- | --- | --- | --- | 
| **Go text/template** | go | {{.variable}} | {{if .condition}} | {{range .items}} | {{.value \| function}} | 
| **Mustache** |   | {{variable}} | {{#condition}} | {{#items}} | 不支持 | 
| **Handlebars** | JavaScript | {{variable}} | {{#if condition}} | {{#each items}} | {{helper value}} | 
| **Django/Jinja2** | python | {{ variable }} | {% if condition %} | {% for item in items %} | {{ function(value) }} | 
| **EJS** | JavaScript | <%= variable %> | <% if (condition) { %> | <% for (item of items) { %> | <%= function(value) %> | 
| **Thymeleaf** |   | ${variable} | th:if="${condition}" | th:each="item : items" | ${#functions.function(value)} | 

