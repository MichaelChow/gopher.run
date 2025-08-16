---
title: "4.4.1 bufio"
date: 2025-08-09T00:38:00Z
draft: false
weight: 4004
---

# 4.4.1 bufio

### 源码解读-bufio

`bufio`包：实现了缓冲的 I/O，它包装了 io.Reader 或 io.Writer 对象，提供缓冲功能和一些文本 I/O 的便利方法。

```go
// Package bufio implements buffered I/O. It wraps an io.Reader or io.Writer
// object, creating another object (Reader or Writer) that also implements
// the interface but provides buffering and some help for textual I/O.
package bufio
```

**核心常量:**

```go
const (
    defaultBufSize = 4096  // 默认缓冲区大小
)

const minReadBufferSize = 16
const maxConsecutiveEmptyReads = 100
```

**Reader 结构体:**

```go
type Reader struct {
    buf          []byte    // **内部缓冲区**
    rd           io.Reader // 底层读取器
    r, w         int       // **缓冲区中数据的读取和写入位置**
    err          error     // 错误状态
    lastByte     int       // 最后读取的字节，用于 UnreadByte **支持回退操作**
    lastRuneSize int       // 最后读取的 rune 大小，用于 UnreadRune **支持回退操作**
}
```

**fill - 填充缓冲区**

```go
func (b *Reader) fill() {
    // **数据滑动：将未读取的数据移到缓冲区开头**
    if b.r > 0 {
        copy(b.buf, b.buf[b.r:b.w])
        b.w -= b.r
        b.r = 0
    }

    // **批量读取：从底层读取器读取数据到缓冲区**
    for i := maxConsecutiveEmptyReads; i > 0; i-- {
        n, err := b.rd.Read(b.buf[b.w:])
        // **错误处理：处理读取错误和空读取**
        if n < 0 {
            panic(errNegativeRead)
        }
        b.w += n
        if err != nil {
            b.err = err
            return
        }
        if n > 0 {
            return
        }
    }
    b.err = io.ErrNoProgress
}
```

**Read - 读取数据**

```go
func (b *Reader) Read(p []byte) (n int, err error) {
    n = len(p)
    if n == 0 {
        if b.Buffered() > 0 {
            return 0, nil
        }
        return 0, b.readErr()
    }
    
    if b.r == b.w {
        // 缓冲区为空
        if b.err != nil {
            return 0, b.readErr()
        }
        if len(p) >= len(b.buf) {
            // **大读取优化：如果读取大小超过缓冲区，直接读取到目标缓冲区**
            n, b.err = b.rd.Read(p)
            if n > 0 {
                b.lastByte = int(p[n-1])
                b.lastRuneSize = -1
            }
            return n, b.readErr()
        }
        // 填充缓冲区
        b.r = 0
        b.w = 0
        n, b.err = b.rd.Read(b.buf)
        if n == 0 {
            return 0, b.readErr()
        }
        b.w += n
    }

    // **缓冲区复用：小读取从缓冲区获取数据**
    n = copy(p, b.buf[b.r:b.w])
    b.r += n
    b.lastByte = int(b.buf[b.r-1])
    b.lastRuneSize = -1
    return n, nil
}
```



**Peek - 预览数据**

```go
// 预览数据而不移动读取位置
func (b *Reader) Peek(n int) ([]byte, error) {
    if n < 0 {
        return nil, ErrNegativeCount
    }

    b.lastByte = -1
    b.lastRuneSize = -1

    // 确保缓冲区有足够数据
    for b.w-b.r < n && b.w-b.r < len(b.buf) && b.err == nil {
        b.fill()
    }

    if n > len(b.buf) {
        return b.buf[b.r:b.w], ErrBufferFull
    }

    var err error
    if avail := b.w - b.r; avail < n {
        n = avail
        err = b.readErr()
        if err == nil {
            err = ErrBufferFull
        }
    }
    return b.buf[b.r : b.r+n], err
}
```



**Writer 结构体**

```go
type Writer struct {
    err error   // 错误状态
    buf []byte  // 缓冲区
    n   int     // 缓冲区中已写入的字节数
    wr  io.Writer // 底层写入器
}
```

**Write - 写入数据**

```go
func (b *Writer) Write(p []byte) (nn int, err error) {
    for len(p) > b.Available() && b.err == nil {
        var n int
        if b.Buffered() == 0 {
            // **大写入优化：如果缓冲区为空且写入量大，直接写入底层**
            n, b.err = b.wr.Write(p)
        } else {
            // **缓冲写入：小写入先写入缓冲区**
            n = copy(b.buf[b.n:], p)
            b.n += n
            b.Flush()
        }
        nn += n
        p = p[n:]
    }
    if b.err != nil {
        return nn, b.err
    }
    n := copy(b.buf[b.n:], p)
    b.n += n
    nn += n
    return nn, nil
}
```



**Flush - 刷新缓冲区**

```go
func (b *Writer) Flush() error {
    if b.err != nil {
        return b.err
    }
    if b.n == 0 {
        return nil
    }
    n, err := b.wr.Write(b.buf[0:b.n])
    if n < b.n && err == nil {
        err = io.ErrShortWrite
    }
    if err != nil {
        if n > 0 && n < b.n {
            copy(b.buf[0:b.n-n], b.buf[n:b.n])
        }
        b.n -= n
        b.err = err
        return err
    }
    b.n = 0
    return nil
}
```





**Scanner 结构体，**实现了非常常用的**文本扫描器 Scanner**。

Scanner 通过缓冲区和分词函数，将输入流分割成一系列“token”（如行、单词、字节等）。

```go
type Scanner struct {
    r            io.Reader // 输入源
    split        SplitFunc // **分词函数（如何分割token）**
    maxTokenSize int       // token最大长度
    token        []byte    // 当前token
    buf          []byte    // 缓冲区
    start, end   int       // buf中未处理数据的起止
    err          error     // 错误状态
    empties      int       // 连续空token计数
    scanCalled   bool      // 是否已调用Scan
    done         bool      // 是否扫描结束
}
```



**SplitFunc 分词函数：**定义如何把输入切分成一个个 token

**内置分词器**：

- **ScanLines：**按行分割。去除行尾的 \n 或 \r\n，最后一行即使没有换行符也会返回。
- **ScanWords：**按空白字符分割。返回每个单词。
- **ScanBytes：**按字节作为一个token分割。
- **ScanRunes：按**每个UTF-8编码的rune作为一个 token分割。支持错误处理。
```go
type SplitFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)
```

**NewScanner：**创建一个新的 Scanner，默认按行分割（ScanLines）

```go
func NewScanner(r io.Reader) *Scanner
```

**NewScanner:**创建一个新的 Scanner，默认按行分割（ScanLines）

```go
func NewScanner(r io.Reader) *Scanner
```

**Scan:**推进到下一个 token，成功返回 true，失败/结束返回 false。每次调用后，可用 s.Bytes() 或 s.Text() 获取当前 token

```go
func (s *Scanner) Scan() bool
```

**Bytes / Text :**获取当前 token 的字节切片或字符串

```go
func (s *Scanner) Bytes() []byte
func (s *Scanner) Text() string
```

**Split:**设置分词函数（如按行、按单词、按字节等）

```go
func (s *Scanner) Split(split SplitFunc)
```

**Buffer:**设置自定义缓冲区和最大 token 长度。

```go
func (s *Scanner) Buffer(buf []byte, max int)
```

**Err:**返回扫描过程中遇到的第一个非 EOF 错误

```go
func (s *Scanner) Err() error
```

**工作流程简述**

1. **初始化**：创建 Scanner，设置分词函数和缓冲区。
1. **循环扫描**：每次调用 Scan()，Scanner 会：
- 检查缓冲区是否有足够数据
- 调用分词函数尝试分割 token
- 如果不够，自动扩容缓冲区并从底层读取更多数据
- 处理各种边界情况（如 EOF、token 过长、错误等）
1. **获取结果**：用 Bytes() 或 Text() 获取当前 token。


### **工厂函数/构造函数**

**工厂函数**是一种设计模式，命名为Newxxx()，封装**创建对象的逻辑复杂性和默认参数，返回对象实例。**替代直接使用 new 关键字或结构体字面量来创建对象。



- **直接类型名**：用于简单的构造函数，通常只需要少量参数
- **New + 类型名**：用于复杂的构造函数，通常需要多个参数或复杂的初始化逻辑
```go
// ✅ 工厂函数方式：简洁清晰
scanner := bufio.NewScanner(os.Stdin)


// ❌ 传统方式：直接构造
scanner := &bufio.Scanner{
    r: os.Stdin,
    // 需要手动设置很多字段...
    // 容易出错，代码冗长
}


// 简单构造函数命名（用于简单对象）
func User(name string) *User {
    return &User{Name: name}
}

func SystemMessage(content string) *Message {
    return &Message{Role: System, Content: content}
}
```



### 迭代器模式

是一种行为型设计模式，它提供了一种方法顺序访问一个聚合对象中的各个元素，而又不暴露其内部的表示。

**模式优势：**

1. **解耦**：将遍历逻辑与聚合对象分离
1. **支持多种遍历方式**：可以轻松实现不同的遍历策略
1. **简化聚合接口**：聚合对象不需要暴露内部结构
1. **支持并发**：异步迭代器支持并发访问
1. **类型安全**：使用Go泛型确保类型安全


**for-range循环（Go语言内置迭代器）：**

```go
// 切片迭代
numbers := []int{1, 2, 3, 4, 5}
for index, value := range numbers {
    fmt.Printf("索引: %d, 值: %d\n", index, value)
}
```

**这个例子如何体现迭代器模式：**

1. **抽象迭代器接口**：for-range语法隐藏了迭代器的复杂性
1. **具体迭代器**：Go运行时自动为切片创建迭代器
1. **聚合对象**：numbers切片是被迭代的数据集合
1. **客户端代码**：for循环体是使用迭代器的客户端
底层实现：

```go
// Go编译器会将for-range转换为类似这样的代码
for i := 0; i < len(numbers); i++ {
    index := i
    value := numbers[i]
    fmt.Printf("索引: %d, 值: %d\n", index, value)
}
```





**bufio.Scanner（显式迭代器）：**

```go
scanner := bufio.NewScanner(reader)
for scanner.Scan() {
    line := scanner.Text()
    fmt.Printf("读取行: %s\n", line)
}
```

**这个例子如何体现迭代器模式：**

1. **迭代器接口**：bufio.Scanner实现了迭代器接口
1. **迭代方法**：Scan()方法移动到下一个元素
1. **访问方法**：Text()方法获取当前元素
1. **状态管理**：Scanner内部维护迭代状态
**Scanner的迭代器接口设计：**

```go
// Scanner的简化接口（实际实现更复杂）
type Scanner interface {
    Scan() bool      // 移动到下一个元素，返回是否成功
    Text() string    // 获取当前元素
    Err() error      // 获取迭代过程中的错误
}
```





**迭代器模式的**包含以下核心组件：

1. **Iterator（迭代器接口）:**定义访问和遍历元素的接口
```go
type Iterator[T any] interface {
    Next() (T, bool)    // 移动到下一个元素
    Current() T         // 获取当前元素
    HasNext() bool      // 是否还有下一个元素
}
```

1. **ConcreteIterator（具体迭代器）**：实现迭代器接口，跟踪当前遍历位置
```go
type SliceIterator[T any] struct {
    data  []T
    index int
}

func (it *SliceIterator[T]) Next() (T, bool) {
    if it.index >= len(it.data) {
        var zero T
        return zero, false
    }
    value := it.data[it.index]
    it.index++
    return value, true
}
```

1. **Aggregate（聚合对象）**：定义创建迭代器对象的接口
```go
type Aggregate[T any] interface {
    CreateIterator() Iterator[T]
}

type SliceAggregate[T any] struct {
    data []T
}

func (sa *SliceAggregate[T]) CreateIterator() Iterator[T] {
    return &SliceIterator[T]{data: sa.data, index: 0}
}
```

1. **Client（客户端）**：使用迭代器访问聚合对象：
```go
func main() {
    numbers := &SliceAggregate[int]{data: []int{1, 2, 3, 4, 5}}
    iterator := numbers.CreateIterator()
    
    for iterator.HasNext() {
        if value, ok := iterator.Next(); ok {
            fmt.Printf("值: %d\n", value)
        }
    }
}
```



