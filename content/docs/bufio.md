---
title: "bufio"
date: 2025-08-09T00:38:00Z
draft: false
---

# bufio

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


### **工厂函数**

**工厂函数**是一种设计模式，命名为Newxxx()，封装**创建对象的逻辑复杂性和默认参数，返回对象实例。**替代直接使用 new 关键字或结构体字面量来创建对象。

```go
// ✅ 工厂函数方式：简洁清晰
scanner := bufio.NewScanner(os.Stdin)


// ❌ 传统方式：直接构造
scanner := &bufio.Scanner{
    r: os.Stdin,
    // 需要手动设置很多字段...
    // 容易出错，代码冗长
}
```





