---
title: "4.4 IO"
date: 2025-08-16T10:18:00Z
draft: false
weight: 4004
---

# 4.4 IO

1. **格式化I/O**
```go
/usr/local/go/src/
├── fmt/                # 格式化I/O
│   ├── doc.go          # 文档
│   ├── errors.go       # 错误处理
│   ├── format.go       # 格式化
│   ├── print.go        # 打印
│   ├── scan.go         # 扫描
│   └── 其他格式化文件
└── 其他I/O包
```

1. **底层I/O**
```go
/usr/local/go/src/
├── io/                 # I/O接口
│   ├── io.go           # 核心I/O接口
│   ├── multi.go        # 多读取器/写入器
│   └── 其他I/O文件
├── bufio/              # 缓冲I/O
│   ├── bufio.go        # 缓冲I/O实现
│   ├── scan.go         # 扫描器
│   └── 其他缓冲文件
└── 其他I/O实现
```



## 📚 目录

- [4.4.1 bufio](4-4-1-bufio/)

