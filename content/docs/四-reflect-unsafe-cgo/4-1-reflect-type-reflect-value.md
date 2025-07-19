---
title: "4.1 reflect.Type、reflect.Value"
date: 2025-05-16T23:19:00Z
draft: false
weight: 4001
---

# 4.1 reflect.Type、reflect.Value

### 为什么使用反射

- **反射(reflection)的定义**：**在编译时不知道变量的具体类型的情况下，能够在运行时更新变量和检查它们的值、调用它们的方法和它们支持的内在操作。**（费曼学习法： Java反射在切面等的应用、性能损耗等）反射也让我们可以把类型当做头等值。
- **为什么使用反射**
    - 需要编写能统一处理不同值类型的函数，而这些类型可能无法共享同一个接口
    - 类型可能没有确定的表示方式
    - 某些类型在设计函数时可能还不存在
- 示例：
    - fmt.Fprintf的例子：一个典型的反射应用，它能够对任意类型的值进行格式化和打印输出，包括用户自定义的类型
    - Sprint函数使用switch类型分支实现的局限性：
        - 组合类型的数目基本是无穷的
        - 即使能识别底层类型（如map[string][]string），也无法匹配具体的用户定义类型（如url.Values）
        - 无法包含所有可能的类型，会造成对外部库的循环依赖
        ```go
        func Sprint(x interface{}) string {
            type stringer interface {
                String() string
            }
            switch x := x.(type) {
            case stringer:
                return x.String()
            case string:
                return x
            case int:
                return strconv.Itoa(x)
            // ...similar cases for int16, uint32, and so on...
            case bool:
                if x {
                    return "true"
                }
                return "false"
            default:
                // array, chan, func, map, pointer, slice, struct
                return "???"
            }
        }
        ```


## **reflect.Type、reflect.Value**

- 反射由reflect包提供的，定义了如下两个重要的类型
- **reflect.Type类型****：**
    - 表示Go类型的interface，其有许多方法来区分类型以及检查它们的组成部分；与接口的类型描述信息相对应，用于标识接口值的动态类型。
    - reflect.TypeOf()：**入参**接受任意的 interface{} 类型，并以 reflect.Type 形式返回其动态类型：
        - reflect.TypeOf(3)：将一个具体的值（3）转为接口类型会有一个隐式的接口转换操作，它会创建一个包含两个信息的接口值：操作数的动态类型（这里是 int）和它的动态的值
            ```go
            t := reflect.TypeOf(3)  // a reflect.Type
            fmt.Println(t.String()) // "int"
            fmt.Println(t)          // "int"
            var w io.Writer = os.Stdout
            fmt.Println(reflect.TypeOf(w)) // "*os.File"
            ```
        - 因为打印一个接口的动态类型对于调试和日志是有帮助的，fmt.Printf 提供了一个缩写 %T 参数，内部使用 reflect.TypeOf 来输出：
            ```go
            fmt.Printf("%T\n", 3) // "int"
            ```
- **reflect.Value类型****：**
    - 入参接受任意的 interface{} 类型，并返回一个装载着其动态值的 reflect.Value； 
    - 和reflect.TypeOf 类似，返回的结果也是具体的类型，但reflect.Value 也可以持有一个接口值；
    - 和 reflect.Type类似，reflect.Value 也满足 fmt.Stringer 接口，但除非 Value 持有的是字符串，否则 String 方法只返回其类型。而使用 fmt 包的 %v 标志参数会对 reflect.Values 特殊处理。(贝叶斯学习法：生产环境代码如果大量使用%v打印大型结构体，由于使用反射会导致大量性能损耗，甚至引发故障)
        ```go
        v := reflect.ValueOf(3) // a reflect.Value
        fmt.Println(v)          // "3"
        fmt.Printf("%v\n", v)   // "3"
        fmt.Println(v.String()) // NOTE: "<int Value>"
        ```
    - 对 Value 调用 Type 方法将返回具体类型所对应的 reflect.Type：
        ```go
        t := v.Type()           // a reflect.Type
        fmt.Println(t.String()) // "int"
        ```
    - reflect.ValueOf 的逆操作是 reflect.Value.Interface 方法。它返回一个 interface{} 类型，装载着与 reflect.Value 相同的具体值：
        ```go
        v := reflect.ValueOf(3) // a reflect.Value
        x := v.Interface()      // an interface{}
        i := x.(int)            // an int
        fmt.Printf("%d\n", i)   // "3"
        ```
    - reflect.Value 和 interface{} 都能装载任意的值。所不同的是：
        - 一个空的接口隐藏了值内部的表示方式和所有方法，因此只有我们知道具体的动态类型才能使用类型断言来访问内部的值（就像上面那样），内部值我们没法访问。
        - 而相比之下，一个 Value 则有很多方法来检查其内容，无论它的具体类型是什么。
- 示例：使用反射替代Sprint的类型分支
    - 使用 reflect.Value 的 Kind 方法来替代之前的类型 switch。虽然还是有无穷多的类型，但是它们的 kinds 类型却是有限的：
        - Bool、String 和 所有数字类型的基础类型；
        - Array 和 Struct 对应的聚合类型；
        - Chan、Func、Ptr、Slice 和 Map 对应的引用类型；
        - interface 类型；
        - 表示空值的 Invalid 类型。（空的 reflect.Value 的 kind 即为 Invalid。）
        ```go
        package format
        import (
            "reflect"
            "strconv"
        )
        // Any formats any value as a string.
        func Any(value interface{}) string {
            return formatAtom(reflect.ValueOf(value))
        }
        // formatAtom formats a value without inspecting its internal structure.
        func formatAtom(v reflect.Value) string {
            switch v.Kind() {
            case reflect.Invalid:
                return "invalid"
            case reflect.Int, reflect.Int8, reflect.Int16,
        	        reflect.Int32, reflect.Int64:
                return strconv.FormatInt(v.Int(), 10)
            case reflect.Uint, reflect.Uint8, reflect.Uint16,
                reflect.Uint32, reflect.Uint64, reflect.Uintptr:
                return strconv.FormatUint(v.Uint(), 10)
            // ...floating-point and complex cases omitted for brevity...
            case reflect.Bool:
                return strconv.FormatBool(v.Bool())
            case reflect.String:
                return strconv.Quote(v.String())
            case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
                return v.Type().String() + " 0x" +
                    strconv.FormatUint(uint64(v.Pointer()), 16)
            default: // reflect.Array, reflect.Struct, reflect.Interface
                return v.Type().String() + " value"
            }
        }
        ```




