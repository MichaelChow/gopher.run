---
title: "4.7.1 context"
date: 2025-07-31T01:27:00Z
draft: false
weight: 4007
---

# 4.7.1 context



> godoc：[https://pkg.go.dev/context](https://pkg.go.dev/context)
goblog：[https://go.dev/blog/context](https://go.dev/blog/context)



用于广播消息和传递数据，是 Go 语言中非常重要的并发控制机制。



在 Go 服务器中，每个传入的请求都在自己的 goroutine 中处理。请求处理器通常会启动额外的 goroutine 来访问数据库和 RPC 服务等后端。处理一个请求的一组 goroutine 通常需要访问特定于请求的值，例如最终用户的身份、授权令牌和请求的截止时间。当请求被取消或超时时，所有处理该请求的 goroutine 应该快速退出，以便系统可以回收它们正在使用的任何资源。

在谷歌，我们开发了一个 `context` 包，它使得跨 API 边界传递请求范围值、取消信号和截止日期变得容易，以便所有参与处理请求的 goroutine。该包作为 context 公开提供。本文介绍了如何使用该包，并提供了一个完整的可工作示例。



**Context 包的核心特性**：

1. **取消控制**：提供统一的取消机制
1. **超时控制**：支持截止时间和超时设置
1. **值传递**：在函数间安全传递请求范围的数据
1. **链式传播**：取消信号自动传播到子 Context
1. **线程安全**：支持并发访问
1. **资源管理**：自动清理相关资源
**设计哲学**

- **显式传递**：Context 作为第一个参数显式传递
- **不可变性**：Context 创建后不可修改
- **链式设计**：支持 Context 的嵌套和传播
- **资源安全**：自动管理相关资源


## **1. Context 接口设计**

**核心接口设计哲学：**

- **不可变性**：Context 一旦创建就不能修改
- **链式传递**：Context 可以在函数间传递
- **取消传播**：子 Context 会继承父 Context 的取消信号
```go
type Context interface {
    // 1. 截止时间
    Deadline() (deadline time.Time, ok bool)
    
    // 2. 取消信号通道
    Done() <-chan struct{}
    
    // 3. 错误信息
    Err() error
    
    // 4. 值传递
    Value(key any) any
}
```



## **2. 基础 Context 实现**

### **emptyCtx - 空 Context**

```go
type emptyCtx struct{}

func (emptyCtx) Deadline() (deadline time.Time, ok bool) {
    return  // 返回零值，表示无截止时间
}

func (emptyCtx) Done() <-chan struct{} {
    return nil  // 返回 nil，表示永不取消
}

func (emptyCtx) Err() error {
    return nil  // 返回 nil，表示无错误
}

func (emptyCtx) Value(key any) any {
    return nil  // 返回 nil，表示无值
}
```



### **Background 和 TODO**

Background returns a non-nil, empty [Context](https://pkg.go.dev/context#Context). It is **never canceled, has no values, and has no deadline**.  Background 返回一个非空的空 Context。它永远不会被取消，没有值，也没有截止时间。

It is typically used by the main function, initialization, and tests, and as the top-level Context for incoming requests.通常由 main 函数、初始化和测试使用，以及作为传入请求的顶层 Context。

```go
func Background() Context {
    return backgroundCtx{}  // 根 Context，returns a non-nil, empty Context
}
```



```go
type backgroundCtx struct{ emptyCtx }
type todoCtx struct{ emptyCtx }


func TODO() Context {
    return todoCtx{}  // 临时 Context
}
```

## **3. 可取消 Context 实现**

**cancelCtx 结构**

```go
type cancelCtx struct {
    Context  // 嵌入父 Context
    
    mu       sync.Mutex            // 保护以下字段
    done     atomic.Value          // chan struct{}，懒加载
    children map[canceler]struct{} // 子 Context 集合
    err      error                 // 取消原因
    cause    error                 // 取消原因（详细）
}
```

**取消机制**

```go
func (c *cancelCtx) cancel(removeFromParent bool, err, cause error) {
    c.mu.Lock()
    if c.err != nil {
        c.mu.Unlock()
        return // 已经取消
    }
    
    c.err = err
    c.cause = cause
    
    // 关闭 done 通道
    d, _ := c.done.Load().(chan struct{})
    if d == nil {
        c.done.Store(closedchan)
    } else {
        close(d)
    }
    
    // 取消所有子 Context
    for child := range c.children {
        child.cancel(false, err, cause)
    }
    c.children = nil
    c.mu.Unlock()
    
    // 从父 Context 中移除
    if removeFromParent {
        removeChild(c.Context, c)
    }
}
```

## **4. WithCancel 实现**

**创建可取消 Context**

```go
func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
    c := withCancel(parent)
    return c, func() { c.cancel(true, Canceled, nil) }
}

func withCancel(parent Context) *cancelCtx {
    if parent == nil {
        panic("cannot create context from nil parent")
    }
    
    c := &cancelCtx{}
    c.propagateCancel(parent, c)  // 设置取消传播
    return c
}
```

**取消传播机制**

```go
func (c *cancelCtx) propagateCancel(parent Context, child canceler) {
    c.Context = parent
    
    done := parent.Done()
    if done == nil {
        return // 父 Context 永不取消
    }
    
    // 检查父 Context 是否已经取消
    select {
    case <-done:
        child.cancel(false, parent.Err(), Cause(parent))
        return
    default:
    }
    
    // 尝试将子 Context 添加到父 Context 的子集合中
    if p, ok := parentCancelCtx(parent); ok {
        p.mu.Lock()
        if p.err != nil {
            // 父 Context 已经取消
            child.cancel(false, p.err, p.cause)
        } else {
            if p.children == nil {
                p.children = make(map[canceler]struct{})
            }
            p.children[child] = struct{}{}
        }
        p.mu.Unlock()
        return
    }
    
    // 启动 goroutine 监听父 Context 的取消信号
    goroutines.Add(1)
    go func() {
        select {
        case <-parent.Done():
            child.cancel(false, parent.Err(), Cause(parent))
        case <-child.Done():
        }
    }()
}
```



## **5. 定时 Context 实现**

**timerCtx 结构**

```go
type timerCtx struct {
    cancelCtx
    timer *time.Timer  // 定时器
    
    deadline time.Time  // 截止时间
}

func (c *timerCtx) Deadline() (deadline time.Time, ok bool) {
    return c.deadline, true
}
```

**WithDeadline 实现**

```go
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
    if parent == nil {
        panic("cannot create context from nil parent")
    }
    
    // 检查父 Context 的截止时间
    if cur, ok := parent.Deadline(); ok && cur.Before(d) {
        return WithCancel(parent)  // 父 Context 截止时间更早
    }
    
    c := &timerCtx{
        deadline: d,
    }
    c.cancelCtx.propagateCancel(parent, c)
    
    dur := time.Until(d)
    if dur <= 0 {
        // 截止时间已过
        c.cancel(true, DeadlineExceeded, nil)
        return c, func() { c.cancel(false, Canceled, nil) }
    }
    
    c.mu.Lock()
    defer c.mu.Unlock()
    if c.err == nil {
        c.timer = time.AfterFunc(dur, func() {
            c.cancel(true, DeadlineExceeded, nil)
        })
    }
    
    return c, func() { c.cancel(true, Canceled, nil) }
}
```

## **6. 值传递 Context 实现**

**valueCtx 结构**

```go
type valueCtx struct {
    Context
    key, val any
}

func (c *valueCtx) Value(key any) any {
    if c.key == key {
        return c.val
    }
    return value(c.Context, key)  // 向上查找
}
```

**WithValue 实现**

```go
func WithValue(parent Context, key, val any) Context {
    if parent == nil {
        panic("cannot create context from nil parent")
    }
    if key == nil {
        panic("nil key")
    }
    if !reflectlite.TypeOf(key).Comparable() {
        panic("key is not comparable")
    }
    
    return &valueCtx{parent, key, val}
}
```



## **7. 错误类型**

**预定义错误**

```go
var Canceled = errors.New("context canceled")
var DeadlineExceeded error = deadlineExceededError{}

type deadlineExceededError struct{}

func (deadlineExceededError) Error() string { 
    return "context deadline exceeded" 
}
func (deadlineExceededError) Timeout() bool { 
    return true 
}
func (deadlineExceededError) Temporary() bool { 
    return true 
}
```



## **8. 实际使用示例**

**基本使用**

```go
func processWithTimeout(ctx context.Context) error {
    // 创建带超时的 Context
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    // 执行操作
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-time.After(10 * time.Second):
        return errors.New("operation completed")
    }
}
```

**值传递**

```go
type userKey struct{}

func processUser(ctx context.Context, userID string) {
    // 将用户信息存储到 Context
    ctx = context.WithValue(ctx, userKey{}, userID)
    
    // 在函数链中传递
    processUserData(ctx)
}

func processUserData(ctx context.Context) {
    // 从 Context 中获取用户信息
    if userID, ok := ctx.Value(userKey{}).(string); ok {
        fmt.Printf("Processing user: %s\n", userID)
    }
}
```

**取消传播**

```go
func worker(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d: context canceled\n", id)
            return
        case <-time.After(time.Second):
            fmt.Printf("Worker %d: working...\n", id)
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // 启动多个 worker
    for i := 0; i < 3; i++ {
        go worker(ctx, i)
    }
    
    // 5 秒后取消所有 worker
    time.Sleep(5 * time.Second)
    cancel()
    
    time.Sleep(time.Second) // 等待 worker 退出
}
```



## **9. 性能优化**

**原子操作**

```go
// 使用 atomic.Value 存储 done 通道
done atomic.Value  // of chan struct{}

func (c *cancelCtx) Done() <-chan struct{} {
    d := c.done.Load()
    if d != nil {
        return d.(chan struct{})
    }
    
    c.mu.Lock()
    defer c.mu.Unlock()
    d = c.done.Load()
    if d == nil {
        d = make(chan struct{})
        c.done.Store(d)
    }
    return d.(chan struct{})
}
```

**懒加载**

```go
// done 通道懒加载，只在需要时创建
var closedchan = make(chan struct{})

func init() {
    close(closedchan)
}
```



## **10. 设计亮点**

**1. 不可变性**

```go
// Context 一旦创建就不能修改，确保线程安全
ctx := context.Background()
ctx = context.WithValue(ctx, key, value)  // 创建新的 Context
```

**2. 链式传递**

```go
// Context 可以在函数间安全传递
func f1(ctx context.Context) {
    ctx = context.WithTimeout(ctx, time.Second)
    f2(ctx)
}

func f2(ctx context.Context) {
    // 继承 f1 的超时设置
    select {
    case <-ctx.Done():
        return
    }
}
```

**3. 取消传播**

```go
// 父 Context 取消时，所有子 Context 都会被取消
parent, cancel := context.WithCancel(context.Background())
child1, _ := context.WithCancel(parent)
child2, _ := context.WithCancel(parent)

cancel()  // 取消父 Context
// child1 和 child2 都会被自动取消
```



**4. 资源管理**

```go


// 使用 defer 确保资源清理
func process(ctx context.Context) {
    ctx, cancel := context.WithTimeout(ctx, time.Second)
    defer cancel()  // 确保资源被清理
    
    // 处理逻辑
}
```



### 函数式选项的设计模式 Functional Options Pattern

优雅、灵活的方式来配置对象和函数参数：

```go
// 自动清理资源
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
    return WithDeadline(parent, time.Now().Add(timeout))
}

// 创建带超时的上下文
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```



如果使用**构造函数参数：**

```go
// 如果context包这样设计：
type Context struct {
    deadline time.Time
    cancel   chan struct{}
}

// 构造函数需要很多参数
func NewContext(parent Context, deadline time.Time, cancel chan struct{}) Context {
    // ...
}

// 使用时就变得很复杂：
deadline := time.Now().Add(5 * time.Second)
cancelChan := make(chan struct{})
ctx := context.NewContext(context.Background(), deadline, cancelChan)
```

如果使用**Setter方法：**

```go
// 如果使用Setter方法：
ctx := context.NewContext()
ctx.SetDeadline(time.Now().Add(5 * time.Second))
ctx.SetCancelChan(make(chan struct{}))

// 问题：
// 1. 需要多行代码
// 2. 状态可能不一致
// 3. 没有返回值（无法获取cancel函数）
```



**WithTimeout的优势:**

1. **函数式编程风格**
```go
// 一行代码完成所有配置
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

// 对比传统方式：
// 1. 创建Context
// 2. 设置deadline
// 3. 设置cancel机制
// 4. 处理错误
```

1. **不可变性（Immutability）**
```go
// WithTimeout返回新的Context，不修改原Context
parent := context.Background()
ctx1, cancel1 := context.WithTimeout(parent, 5*time.Second)
ctx2, cancel2 := context.WithTimeout(parent, 10*time.Second)

// parent保持不变
// ctx1和ctx2是独立的，互不影响
```

1. **链式调用**
```go
// 可以链式组合多个WithXxx
ctx := context.Background()
ctx = context.WithValue(ctx, "user_id", "123")
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

// 或者更简洁的写法：
ctx, cancel := context.WithTimeout(
    context.WithValue(context.Background(), "user_id", "123"),
    5*time.Second,
)
defer cancel()
```

