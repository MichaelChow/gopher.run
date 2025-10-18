---
title: "4.5 系统编程"
date: 2025-07-31T01:27:00Z
draft: false
weight: 4005
---

# 4.5 系统编程

1. **操作系统接口**
```go
/usr/local/go/src/
├── os/                 # 操作系统接口
│   ├── dir.go          # 目录操作
│   ├── env.go          # 环境变量
│   ├── error.go        # 错误定义
│   ├── file.go         # 文件操作
│   ├── path.go         # 路径操作
│   ├── proc.go         # 进程操作
│   ├── stat.go         # 文件状态
│   ├── types.go        # 类型定义
│   └── 其他系统文件
├── path/               # 路径操作
│   ├── match.go        # 路径匹配
│   └── path.go         # 路径操作
├── filepath/           # 文件路径
│   ├── match.go        # 文件路径匹配
│   ├── path.go         # 文件路径操作
│   └── 其他路径文件
└── 其他系统接口
```

1. **系统调用**
```go
/usr/local/go/src/
├── syscall/            # 系统调用
│   ├── syscall.go      # 系统调用接口
│   ├── syscall_linux.go # Linux系统调用
│   ├── syscall_windows.go # Windows系统调用
│   └── 其他系统调用文件
└── 其他底层接口
```





## 📚 目录

- [4.5.1 os](4-5-1-os/)

