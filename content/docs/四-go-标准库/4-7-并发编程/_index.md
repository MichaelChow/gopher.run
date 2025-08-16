---
title: "4.7 并发编程"
date: 2025-07-31T01:37:00Z
draft: false
weight: 4007
---

# 4.7 并发编程

1. **同步原语**
```go
/usr/local/go/src/
├── sync/               # 同步原语 - 并发编程核心
│   ├── cond.go         # 条件变量
│   ├── map.go          # 并发映射
│   ├── mutex.go        # 互斥锁
│   ├── once.go         # 一次性执行
│   ├── pool.go         # 对象池
│   ├── rwmutex.go      # 读写锁
│   ├── waitgroup.go    # 等待组
│   ├── atomic.go       # 原子操作
│   └── 其他同步文件
├── context/            # 上下文 - 并发控制
│   └── context.go      # 上下文实现
└── atomic/             # 原子操作 - 无锁编程
    ├── doc.go          # 文档
    ├── value.go        # 原子值
    ├── 64bit.go        # 64位原子操作
    └── 其他原子文件
```





