---
title: "4.1 语言核心"
date: 2025-07-31T01:28:00Z
draft: false
weight: 4001
---

# 4.1 语言核心

1. **内置类型和函数 (Builtin)**
理解基本类型和内置函数

```go
/usr/local/go/src/builtin/
├── builtin.go          # 语言内置定义
    - 基本类型系统：int, uint, float, complex, bool, string, byte, rune
    - 内置函数：len, cap, append, make, new, delete, panic, recover
    - 内置接口：error, comparable
    - 内置常量：true, false, iota, nil
```

1. **类型系统 (Type System)**
理解反射机制和不安全操作

```go
/usr/local/go/src/
├── reflect/            # 反射系统 - 类型系统的核心
│   ├── type.go         # 类型信息表示
│   ├── value.go        # 值操作和类型转换
│   ├── makefunc.go     # 动态函数创建
│   ├── swapper.go      # 切片元素交换
│   ├── deepequal.go    # 深度相等比较
│   └── 其他反射核心文件
├── unsafe/             # 不安全操作 - 底层类型操作
│   └── unsafe.go       # 指针操作、类型转换
└── internal/           # 内部实现
    ├── reflectlite/    # 轻量级反射
    └── 其他内部包
```



1. **运行时系统 (Runtime)**
理解运行时系统（内存管理、GC、调度器）

```go
/usr/local/go/src/runtime/
├── 内存管理 (Memory Management)
│   ├── malloc.go       # 内存分配器入口
│   ├── mheap.go        # 堆内存管理
│   ├── mcentral.go     # 中心缓存管理
│   ├── mcache.go       # 线程本地缓存
│   ├── mspan.go        # 内存跨度管理
│   └── mstats.go       # 内存统计
├── 垃圾回收 (Garbage Collection)
│   ├── gc.go           # GC主控制器
│   ├── mgc.go          # GC标记阶段
│   ├── mgcmark.go      # GC标记工作
│   ├── mgcscavenge.go  # GC清理阶段
│   └── mheap.go        # GC堆管理
├── 协程调度 (Goroutine Scheduler)
│   ├── scheduler.go    # 调度器主逻辑
│   ├── proc.go         # 处理器管理
│   ├── stack.go        # 栈管理
│   ├── stubs.go        # 汇编存根
│   └── 其他调度文件
├── 内置类型运行时 (Built-in Type Runtime)
│   ├── string.go       # 字符串运行时
│   ├── slice.go        # 切片运行时
│   ├── map.go          # 映射运行时
│   ├── chan.go         # 通道运行时
│   ├── iface.go        # 接口运行时
│   └── type.go         # 类型系统运行时
├── 系统接口 (System Interface)
│   ├── os_linux.go     # Linux系统接口
│   ├── os_windows.go   # Windows系统接口
│   ├── os_darwin.go    # macOS系统接口
│   └── 其他系统文件
└── 其他运行时组件
```









