---
title: "4.8 数据序列化"
date: 2025-07-31T01:26:00Z
draft: false
weight: 4008
---

# 4.8 数据序列化

**编码和解码**

```go
/usr/local/go/src/
├── encoding/           # 编码包 - 数据转换
│   ├── json/           # JSON处理 - 最常用的序列化
│   │   ├── decode.go   # JSON解码
│   │   ├── encode.go   # JSON编码
│   │   ├── fold.go     # JSON折叠
│   │   ├── scanner.go  # JSON扫描器
│   │   ├── stream.go   # JSON流
│   │   ├── table.go    # JSON表
│   │   └── 其他JSON文件
│   ├── xml/            # XML处理 - 结构化数据
│   │   ├── read.go     # XML读取
│   │   ├── write.go    # XML写入
│   │   ├── marshal.go  # XML序列化
│   │   └── 其他XML文件
│   ├── base64/         # Base64编码 - 二进制编码
│   │   └── base64.go   # Base64实现
│   ├── hex/            # 十六进制 - 二进制表示
│   │   └── hex.go      # 十六进制实现
│   ├── binary/         # 二进制编码 - 底层编码
│   │   └── binary.go   # 二进制实现
│   ├── asn1/           # ASN.1编码 - 密码学标准
│   │   └── asn1.go     # ASN.1实现
│   └── 其他编码
└── 其他序列化
```









