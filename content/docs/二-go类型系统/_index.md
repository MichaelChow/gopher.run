---
title: "二、Go类型系统"
date: 2024-11-25T01:47:00Z
draft: false
weight: 2000
---

# 二、Go类型系统

Go语言将数据类型分为四类：

| **基础类型（basic type）** | boolean、numeric、string，是Go语言世界数据的原子。 | 
| --- | --- | 
| **组合类型（aggregate type）** | array**、**struct，由基础类型组合，值由内存中的一组变量构成，是Go语言世界数据的分子。 | 
| **引用类型（reference type）** | pointer、slice、map、function、**channel**，都**是间接指向****程序变量或状态，**操作所引用数据的全部效果会遍及该数据的全部引用； | 
| **接口类型（interface type）** |   | 

> *计算机底层全是bit，而****实际操作则是基于大小固定的单元中的数值，称为字（word），如整数、浮点数、比特数组、内存地址等****；进而构成更大的聚合类型；Go的数据类型宽泛，向下匹配硬件特性，向上满足程序员所需。*



## Go的内存对齐

```go
// 计算机最底层是 bit（二进制位）
// 每个 bit 只能是 0 或 1

// 在 Go 中，我们可以操作 bit
var flags uint8 = 0b10101010  // 8 个 bit
var mask uint8 = 0b00000001   // 1 个 bit

// bit 操作
result := flags & mask         // 按位与
result = flags | mask          // 按位或
result = flags ^ mask          // 按位异或
result = flags << 1            // 左移
result = flags >> 1            // 右移

// bit 的物理实现: 在硬件层面，bit 通过以下方式表示：
// - 电压高低（高电压 = 1，低电压 = 0）
// - 磁化方向（北 = 1，南 = 0）
// - 光信号（有光 = 1，无光 = 0）
```



```go
// 字是计算机处理的基本单位
// 现代计算机通常是 64 位（8 字节）

// Go 中的基本数据类型对应不同的字大小
var (
    // 8 位（1 字节）
    b byte = 255        // uint8
    
    // 16 位（2 字节）
    s int16 = 32767
    
    // 32 位（4 字节）
    i int32 = 2147483647
    f float32 = 3.14
    
    // 64 位（8 字节）
    l int64 = 9223372036854775807
    d float64 = 3.14159265359
    
    // 平台相关（32 位或 64 位）
    n int = 42  // 在 64 位系统上是 64 位
)

// 内存地址
// 内存地址也是基于字的
var x int = 42
var ptr *int = &x

// 指针的大小取决于平台
fmt.Printf("Pointer size: %d bytes\n", unsafe.Sizeof(ptr))
// 64 位系统：8 字节
// 32 位系统：4 字节
```



```go
// 数组是相同类型元素的固定大小聚合
// 内存布局：连续的字
// [int][int][int][int][int]
// 每个 int 占用一个字（8 字节）
var arr [5]int = [5]int{1, 2, 3, 4, 5}


// 内存布局：
// [ID: 4字节][Age: 1字节][填充: 3字节][Name: 16字节]
// 总大小：24 字节（考虑内存对齐）
type Person struct {
    ID   int32   // 4 字节
    Age  uint8   // 1 字节
    Name string  // 16 字节（指针 + 长度）
}


// 切片本身：24 字节
// 底层数组：40 字节（5 * 8）
type slice struct {
    ptr *int    // 8 字节：指向底层数组的指针
    len int     // 8 字节：长度
    cap int     // 8 字节：容量
}
var s []int = []int{1, 2, 3, 4, 5}
```



```go
// 整数类型
var (
    // 8 位
    i8 int8 = 127
    
    // 16 位
    i16 int16 = 32767
    
    // 32 位
    i32 int32 = 2147483647
    
    // 64 位
    i64 int64 = 9223372036854775807
    
    // 平台相关
    i int = 42  // 32 位系统：32 位，64 位系统：64 位
)

// 浮点数类型
var (
    f32 float32 = 3.14  // 32 位 IEEE 754
    f64 float64 = 3.14159265359  // 64 位 IEEE 754
)

// 复数类型
var (
    c64 complex64 = 1 + 2i   // 64 位（32 + 32）
    c128 complex128 = 1 + 2i  // 128 位（64 + 64）
)
```



```go
// 高级抽象
// 字符串：底层是字节数组，但提供高级接口
str := "Hello, 世界"
fmt.Println(len(str))        // 13（字节数）
fmt.Println(utf8.RuneCountInString(str))  // 8（字符数）

// 切片：动态数组的抽象
slice := []int{1, 2, 3, 4, 5}
slice = append(slice, 6)    // 动态增长
slice = slice[:3]           // 动态收缩

// 映射：哈希表的抽象
m := map[string]int{"a": 1, "b": 2}
m["c"] = 3                 // 动态插入
delete(m, "a")             // 动态删除

// 类型安全
// 编译时类型检查
var x int = 42
var y string = "hello"
// x = y  // 编译错误：类型不匹配

// 运行时类型检查
func processValue(v interface{}) {
    switch v.(type) {
    case int:
        fmt.Println("Integer:", v)
    case string:
        fmt.Println("String:", v)
    default:
        fmt.Println("Unknown type")
    }
}

```



**内存布局优化：利用硬件特性**

```go
// 结构体字段顺序影响缓存性能
type CacheFriendly struct {
    ID    int64   // 8 字节
    Value float64 // 8 字节
    Flag  bool    // 1 字节
    // 填充 7 字节以保持对齐
}

type CacheUnfriendly struct {
    Flag  bool    // 1 字节
    // 填充 7 字节
    ID    int64   // 8 字节
    Value float64 // 8 字节
}
```

**内存管理：从 bit 到对象:**

**内存层次：**

**内存层次**

```go
// 1. bit 级别：物理内存
// 2. 字级别：CPU 寄存器
// 3. 页级别：操作系统内存管理
// 4. 对象级别：Go 运行时内存管理

// Go 的内存分配器
type mspan struct {
    startAddr uintptr  // 起始地址
    npages    uintptr  // 页数
    allocBits *gcBits  // 分配位图
    // ... 更多字段
}
```

**垃圾回收：**

```go
// Go 的 GC 工作在对象级别，但底层是 bit 操作
type gcWork struct {
    // 工作队列
    wbuf1, wbuf2 *workbuf
    // 标记位图
    markrootNext uint32
    markrootDone uint32
}
```







## 📚 目录

- [2.1 numeric、boolean](2-1-numeric-boolean/)
- [2.2 string](2-2-string/)
- [2.4 array、slice、append](2-4-array-slice-append/)
- [2.5 pointer](2-5-pointer/)
- [2.6 map、make、new](2-6-map-make-new/)
- [2.7 泛型 generics](2-7-泛型-generics/)
- [2.8 struct、组合嵌套、json](2-8-struct-组合嵌套-json/)
- [2.9 func declaration、递归、多值返回、可变参数](2-9-func-declaration-递归-多值返回-可变参数/)
- [2.10 函数变量、匿名函数、闭包](2-10-函数变量-匿名函数-闭包/)
- [2.11 error、panic、recover、defer ](2-11-error-panic-recover-defer-/)
- [2.12 method declaration、指针接收者、struct组合嵌套、封装](2-12-method-declaration-指针接收者-struct组合嵌套-封装/)
- [2.13 interface type、实现、interface变量的值、interface的类型断言](2-13-interface-type-实现-interface变量的值-interface的类型断言/)
- [2.14 常用interface](2-14-常用interface/)



























