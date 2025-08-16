---
title: "6.2 protobuf、grpc"
date: 2025-05-21T08:54:00Z
draft: false
weight: 6002
---

# 6.2 protobuf、grpc

- **官方文档(必读)**：[https://protobuf.dev/programming-guides/proto3/](https://protobuf.dev/programming-guides/proto3/)
    - 全称 Protocol Buffers，两个版本：proto2、proto3
- 示例：
```protobuf
syntax = "proto3";  // proto版本
go_package   // 定义一个go_package

message SearchRequest {   // 数据交换的基本单元 
  string query = 1;
  int32 page_number = 2;
  int32 results_per_page = 3;
}
```

- 数据类型及对应常用语言的类型映射：
| Proto Type | Go Type | Rust Type | Java/Kotlin Type[1] | Python Type[3] | C++ Type | 
| --- | --- | --- | --- | --- | --- | 
| double | float64 | f64 | double | float | double | 
| float | float32 | f32 | float | float | float | 
| int32 | int32 | i32 | int | int | int32_t | 
| int64 | int64 | i64 | long | int/long[4] | int64_t | 
| uint32 | uint32 | u32 | int[2] | int/long[4] | uint32_t | 
| uint64 | uint64 | u64 | long[2] | int/long[4] | uint64_t | 
| sint32 | int32 | i32 | int | int | int32_t | 
| sint64 | int64 | i64 | long | int/long[4] | int64_t | 
| fixed32 | uint32 | u32 | int[2] | int/long[4] | uint32_t | 
| fixed64 | uint64 | u64 | long[2] | int/long[4] | uint64_t | 
| sfixed32 | int32 | i32 | int | int | int32_t | 
| sfixed64 | int64 | i64 | long | int/long[4] | int64_t | 
| bool | bool | bool | boolean | bool | bool | 
| string | string | ProtoString | String | str/unicode[5] | std::string | 
| bytes | []byte | ProtoBytes | ByteString | str (Python 2), bytes (Python 3) | std::string | 

- **Specifying Field Cardinality**
    - repeated
    - map等
- **Defining Services**
    ```protobuf
    service SearchService {   // 定义一个服务，包括多个rpc方法
      rpc Search(SearchRequest) returns (SearchResponse);
    }
    ```






- 官网：[https://grpc.io/](https://grpc.io/)
- 开源中国组织翻译的中文文档：[https://doc.oschina.net/grpc](https://doc.oschina.net/grpc)
