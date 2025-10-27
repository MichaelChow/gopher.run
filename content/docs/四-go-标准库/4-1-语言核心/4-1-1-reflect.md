---
title: "4.1.1 reflect"
date: 2025-05-16T23:19:00Z
draft: false
weight: 4001
---

# 4.1.1 reflect

### **reflect.Type、reflect.Value**

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




### **example：Display，一个递归的值打印器**

- 尽量避免在一个包的API中暴露涉及反射的接口
- 示例：display函数
    - 使用前面定义的打印基础类型——基本类型、函数和chan等——元素值的formatAtom函数
    - 使用reflect.Value的方法来递归显示复杂类型的每一个成员。在递归下降过程中，path字符串，从最开始传入的起始值（这里是“e”），将逐步增长来表示是如何达到当前值（例如“e.args[0].value”）的。
    - **Slice和**Array**：** 两种的处理逻辑是一样的。Len方法返回slice或数组值中的元素个数，Index(i)获得索引i对应的元素，返回的也是一个reflect.Value；如果索引i超出范围的话将导致panic异常，这与数组或slice类型内建的len(a)和a[i]操作类似。display针对序列中的每个元素递归调用自身处理，我们通过在递归处理时向path附加“[i]”来表示访问路径。虽然reflect.Value类型带有很多方法，但是只有少数的方法能对任意值都安全调用。例如，Index方法只能对Slice、数组或字符串类型的值调用，如果对其它类型调用则会导致panic异常。
    - Struct**：** NumField方法报告结构体中成员的数量，Field(i)以reflect.Value类型返回第i个成员的值。成员列表也包括通过匿名字段提升上来的成员。为了在path添加“.f”来表示成员路径，我们必须获得结构体对应的reflect.Type类型信息，然后访问结构体第i个成员的名字。
    - Map**:** MapKeys方法返回一个reflect.Value类型的slice，每一个元素对应map的一个key。和往常一样，遍历map时顺序是随机的。MapIndex(key)返回map中key对应的value。我们向path添加“[key]”来表示访问路径。（我们这里有一个未完成的工作。其实map的key的类型并不局限于formatAtom能完美处理的类型；数组、结构体和接口都可以作为map的key。针对这种类型，完善key的显示信息是练习12.1的任务。）
    - Ptr**指针：** Elem方法返回指针指向的变量，依然是reflect.Value类型。即使指针是nil，这个操作也是安全的，在这种情况下指针是Invalid类型，但是我们可以用IsNil方法来显式地测试一个空指针，这样我们可以打印更合适的信息。我们在path前面添加“*”，并用括弧包含以避免歧义。
    - Interface**：** 再一次，我们使用IsNil方法来测试接口是否是nil，如果不是，我们可以调用v.Elem()来获取接口对应的动态值，并且打印对应的类型和值。
        ```go
        func display(path string, v reflect.Value) {
            switch v.Kind() {
            case reflect.Invalid:
                fmt.Printf("%s = invalid\n", path)
            case reflect.Slice, reflect.Array:
                for i := 0; i < v.Len(); i++ {
                    display(fmt.Sprintf("%s[%d]", path, i), v.Index(i))
                }
            case reflect.Struct:
                for i := 0; i < v.NumField(); i++ {
                    fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
                    display(fieldPath, v.Field(i))
                }
            case reflect.Map:
                for _, key := range v.MapKeys() {
                    display(fmt.Sprintf("%s[%s]", path,
                        formatAtom(key)), v.MapIndex(key))
                }
            case reflect.Ptr:
                if v.IsNil() {
                    fmt.Printf("%s = nil\n", path)
                } else {
                    display(fmt.Sprintf("(*%s)", path), v.Elem())
                }
            case reflect.Interface:
                if v.IsNil() {
                    fmt.Printf("%s = nil\n", path)
                } else {
                    fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
                    display(path+".value", v.Elem())
                }
            default: // basic types, channels, funcs
                fmt.Printf("%s = %s\n", path, formatAtom(v))
            }
        }
        ```
    - 表现：
        ```go
        type Movie struct {
            Title, Subtitle string
            Year            int
            Color           bool
            Actor           map[string]string
            Oscars          []string
            Sequel          *string
        }
        strangelove := Movie{
            Title:    "Dr. Strangelove",
            Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
            Year:     1964,
            Color:    false,
            Actor: map[string]string{
                "Dr. Strangelove":            "Peter Sellers",
                "Grp. Capt. Lionel Mandrake": "Peter Sellers",
                "Pres. Merkin Muffley":       "Peter Sellers",
                "Gen. Buck Turgidson":        "George C. Scott",
                "Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
                `Maj. T.J. "King" Kong`:      "Slim Pickens",
            },
            Oscars: []string{
                "Best Actor (Nomin.)",
                "Best Adapted Screenplay (Nomin.)",
                "Best Director (Nomin.)",
                "Best Picture (Nomin.)",
            },
        }
        Display strangelove (display.Movie):
        strangelove.Title = "Dr. Strangelove"
        strangelove.Subtitle = "How I Learned to Stop Worrying and Love the Bomb"
        strangelove.Year = 1964
        strangelove.Color = false
        strangelove.Actor["Gen. Buck Turgidson"] = "George C. Scott"
        strangelove.Actor["Brig. Gen. Jack D. Ripper"] = "Sterling Hayden"
        strangelove.Actor["Maj. T.J. \"King\" Kong"] = "Slim Pickens"
        strangelove.Actor["Dr. Strangelove"] = "Peter Sellers"
        strangelove.Actor["Grp. Capt. Lionel Mandrake"] = "Peter Sellers"
        strangelove.Actor["Pres. Merkin Muffley"] = "Peter Sellers"
        strangelove.Oscars[0] = "Best Actor (Nomin.)"
        strangelove.Oscars[1] = "Best Adapted Screenplay (Nomin.)"
        strangelove.Oscars[2] = "Best Director (Nomin.)"
        strangelove.Oscars[3] = "Best Picture (Nomin.)"
        strangelove.Sequel = nil
        Display("os.Stderr", os.Stderr)
        // Output:
        // Display os.Stderr (*os.File):
        // (*(*os.Stderr).file).fd = 2
        // (*(*os.Stderr).file).name = "/dev/stderr"
        // (*(*os.Stderr).file).nepipe = 0
        ```
- 可以看出，反射能够访问到结构体中未导出的成员。
    - 需要当心的是这个例子的输出在不同操作系统上可能是不同的，并且随着标准库的发展也可能导致结果不同。（这也是将这些成员定义为私有成员的原因之一！）我们甚至可以用Display函数来显示reflect.Value 的内部构造（在这里设置为`*os.File`的类型描述体）。`Display("rV", reflect.ValueOf(os.Stderr))`调用的输出如下，当然不同环境得到的结果可能有差异：
        ```go
        Display rV (reflect.Value):
        (*rV.typ).size = 8
        (*rV.typ).hash = 871609668
        (*rV.typ).align = 8
        (*rV.typ).fieldAlign = 8
        (*rV.typ).kind = 22
        (*(*rV.typ).string) = "*os.File"
        (*(*(*rV.typ).uncommonType).methods[0].name) = "Chdir"
        (*(*(*(*rV.typ).uncommonType).methods[0].mtyp).string) = "func() error"
        (*(*(*(*rV.typ).uncommonType).methods[0].typ).string) = "func(*os.File) error"
        ...
        ```
- 观察下面两个例子的区别：
    - reflect.ValueOf(i):返回一个具体类型的 Value
    - reflect.ValueOf(&i): 返回一个指向i的指针，对应Ptr类型。在switch的Ptr分支中，对这个值调用 Elem 方法，返回一个Value来表示变量 i 本身，对应Interface类型。像这样一个间接获得的Value，可能代表任意类型的值，包括接口类型。display函数递归调用自身，这次它分别打印了这个接口的动态类型和值。
    ```go
    var i interface{} = 3
    Display("i", i)
    // Output:
    // Display i (int):
    // i = 3
    Display("&i", &i)
    // Output:
    // Display &i (*interface {}):
    // (*&i).type = int
    // (*&i).value = 3
    ```
- 对于目前的实现，如果遇到对象图中含有回环，Display将会陷入死循环，例如下面这个首尾相连的链表：
    ```go
    // a struct that points to itself
    type Cycle struct{ Value int; Tail *Cycle }
    var c Cycle
    c = Cycle{42, &c}
    Display("c", c)
    ```
    - Display会永远不停地进行深度递归打印：
    ```go
    Display c (display.Cycle):
    c.Value = 42
    (*c.Tail).Value = 42
    (*(*c.Tail).Tail).Value = 42
    (*(*(*c.Tail).Tail).Tail).Value = 42
    ...ad infinitum...
    ```
    - 让Display支持这类带环的数据结构需要些技巧，需要额外记录迄今访问的路径；相应会带来成本。通用的解决方案是采用 unsafe 的语言特性
    - 带环的数据结构很少会对fmt.Sprint函数造成问题，因为它很少尝试打印完整的数据结构。例如，当它遇到一个指针的时候，它只是简单地打印指针的数字值**。在打印包含自身的slice或map时可能卡住**，但是这种情况很罕见，不值得付出为了处理回环所需的开销。


### **示例: 编码/序列化为S表达式**

- Display是一个用于显示结构化数据的调试工具，但是它并不能将任意的Go语言对象编码为通用消息然后用于进程间通信。
- Go语言的标准库支持了包括JSON、XML和ASN.1等多种编码格式。还有另一种依然被广泛使用的格式是S表达式格式，采用Lisp语言的语法。但是和其他编码格式不同的是，Go语言自带的标准库并不支持S表达式，主要是因为它没有一个公认的标准规范。
- 示例：定义一个包用于将任意的Go语言对象编码/序列化为S表达式格式，它支持以下结构：
    ```go
    42          integer
    "hello"     string（带有Go风格的引号）
    foo         symbol（未用引号括起来的名字）
    (1 2 3)     list  （括号包起来的0个或多个元素）
    ```
    ```go
    func encode(buf *bytes.Buffer, v reflect.Value) error {
        switch v.Kind() {
        case reflect.Invalid:
            buf.WriteString("nil")
        case reflect.Int, reflect.Int8, reflect.Int16,
            reflect.Int32, reflect.Int64:
            fmt.Fprintf(buf, "%d", v.Int())
        case reflect.Uint, reflect.Uint8, reflect.Uint16,
            reflect.Uint32, reflect.Uint64, reflect.Uintptr:
            fmt.Fprintf(buf, "%d", v.Uint())
        case reflect.String:
            fmt.Fprintf(buf, "%q", v.String())
        case reflect.Ptr:
            return encode(buf, v.Elem())
        case reflect.Array, reflect.Slice: // (value ...)
            buf.WriteByte('(')
            for i := 0; i < v.Len(); i++ {
                if i > 0 {
                    buf.WriteByte(' ')
                }
                if err := encode(buf, v.Index(i)); err != nil {
                    return err
                }
            }
            buf.WriteByte(')')
        case reflect.Struct: // ((name value) ...)
            buf.WriteByte('(')
            for i := 0; i < v.NumField(); i++ {
                if i > 0 {
                    buf.WriteByte(' ')
                }
                fmt.Fprintf(buf, "(%s ", v.Type().Field(i).Name)
                if err := encode(buf, v.Field(i)); err != nil {
                    return err
                }
                buf.WriteByte(')')
            }
            buf.WriteByte(')')
        case reflect.Map: // ((key value) ...)
            buf.WriteByte('(')
            for i, key := range v.MapKeys() {
                if i > 0 {
                    buf.WriteByte(' ')
                }
                buf.WriteByte('(')
                if err := encode(buf, key); err != nil {
                    return err
                }
                buf.WriteByte(' ')
                if err := encode(buf, v.MapIndex(key)); err != nil {
                    return err
                }
                buf.WriteByte(')')
            }
            buf.WriteByte(')')
        default: // float, complex, bool, chan, func, interface
            return fmt.Errorf("unsupported type: %s", v.Type())
        }
        return nil
    }
    ```
    - Marshal函数是对encode的包装，以保持和encoding/...下其它包有着相似的API：
        ```go
        // Marshal encodes a Go value in S-expression form.
        func Marshal(v interface{}) ([]byte, error) {
            var buf bytes.Buffer
            if err := encode(&buf, reflect.ValueOf(v)); err != nil {
                return nil, err
            }
            return buf.Bytes(), nil
        }
        ```
- 下面是Marshal对12.3节的strangelove变量编码后的结果：
    - 整个输出编码为一行中以减少输出的大小，但是也很难阅读。下面是对S表达式手动格式化的结果
    ```go
    ((Title "Dr. Strangelove") (Subtitle "How I Learned to Stop Worrying and Lo
    ve the Bomb") (Year 1964) (Actor (("Grp. Capt. Lionel Mandrake" "Peter Sell
    ers") ("Pres. Merkin Muffley" "Peter Sellers") ("Gen. Buck Turgidson" "Geor
    ge C. Scott") ("Brig. Gen. Jack D. Ripper" "Sterling Hayden") ("Maj. T.J. \
    "King\" Kong" "Slim Pickens") ("Dr. Strangelove" "Peter Sellers"))) (Oscars
    ("Best Actor (Nomin.)" "Best Adapted Screenplay (Nomin.)" "Best Director (N
    omin.)" "Best Picture (Nomin.)")) (Sequel nil))
    ```
    - 编写一个S表达式的美化格式化函数:
    ```go
    ((Title "Dr. Strangelove")
     (Subtitle "How I Learned to Stop Worrying and Love the Bomb")
     (Year 1964)
     (Actor (("Grp. Capt. Lionel Mandrake" "Peter Sellers")
             ("Pres. Merkin Muffley" "Peter Sellers")
             ("Gen. Buck Turgidson" "George C. Scott")
             ("Brig. Gen. Jack D. Ripper" "Sterling Hayden")
             ("Maj. T.J. \"King\" Kong" "Slim Pickens")
             ("Dr. Strangelove" "Peter Sellers")))
     (Oscars ("Best Actor (Nomin.)"
              "Best Adapted Screenplay (Nomin.)"
              "Best Director (Nomin.)"
              "Best Picture (Nomin.)"))
     (Sequel nil))
    ```
    - 和fmt.Print、json.Marshal、Display函数类似，sexpr.Marshal函数处理带环的数据结构也会陷入死循环。


### **通过reflect.Value修改值**

- 前面介绍通过反射来读取变量，下面接受如何通过反射机制来修改变量。
- **一个变量就是一个可寻址的内存空间，里面存储了一个值，并且存储的值可以通过内存地址来更新：**
    - Go语言中类似x、x.f[1]和*p形式的表达式都可以表示变量
    - 但x + 1和f(2)则不是变量
- 对于reflect.Values也有类似的区别。有一些reflect.Values是可取地址的；其它一些则不可以。考虑以下的声明语句：
    - 实际上，所有通过reflect.ValueOf(x)返回的reflect.Value都是不可取地址的。
        - a对应的变量不可取地址。因为a中的值仅仅是整数2的拷贝副本。
        - b中的值也同样不可取地址。
        - c中的值还是不可取地址，它只是一个指针`&x`的拷贝。
            ```go
            x := 2                   // value   type    variable?
            a := reflect.ValueOf(2)  // 2       int     no
            b := reflect.ValueOf(x)  // 2       int     no
            c := reflect.ValueOf(&x) // &x      *int    no
            d := c.Elem()            // 2       int     yes (x)
            ```
    - **通过调用reflect.ValueOf(&x).Elem()，来获取任意变量x对应的可取地址的Value。**
        - **对于d，它是c的解引用方式生成的，指向另一个变量，因此是可取地址的。**
        - 调用reflect.Value的CanAddr方法来判断其是否可以被取地址：
            ```go
            fmt.Println(a.CanAddr()) // "false"
            fmt.Println(b.CanAddr()) // "false"
            fmt.Println(c.CanAddr()) // "false"
            fmt.Println(d.CanAddr()) // "true"
            ```
        - 每当我们通过指针间接地获取的reflect.Value都是可取地址的，即使开始的是一个不可取地址的Value。在反射机制中，所有关于是否支持取地址的规则都是类似的。
            - 如，slice的索引表达式e[i]将隐式地包含一个指针，它就是可取地址的，即使开始的e表达式不支持也没有关系。
            - 以此类推，reflect.ValueOf(e).Index(i)对应的值也是可取地址的，即使原始的reflect.ValueOf(e)不支持也没有关系。
- 要从变量对应的可取地址的reflect.Value来访问变量需要三个步骤。
    - 第一步是调用Addr()方法，它返回一个Value，里面保存了指向变量的指针。
    - 然后是在Value上调用Interface()方法，也就是返回一个interface{}，里面包含指向变量的指针。
    - 最后，如果我们知道变量的类型，我们可以使用类型的断言机制将得到的interface{}类型的接口强制转为普通的类型指针。这样我们就可以通过这个普通指针来更新变量了：
        ```go
        x := 2
        d := reflect.ValueOf(&x).Elem()   // d refers to the variable x
        px := d.Addr().Interface().(*int) // px := &x
        *px = 3                           // x = 3
        fmt.Println(x)                    // "3"
        ```
    - 或者，不使用指针，而是通过调用可取地址的reflect.Value的reflect.Value.Set方法来更新对应的值：
    - 当我们用Display显示os.Stdout结构时，我们发现反射可以越过Go语言的导出规则的限制读取结构体中未导出的成员，比如在类Unix系统上os.File结构体中的fd int成员。然而，利用反射机制并不能修改这些未导出的成员：
        ```go
        stdout := reflect.ValueOf(os.Stdout).Elem() // *os.Stdout, an os.File var
        fmt.Println(stdout.Type())                  // "os.File"
        fd := stdout.FieldByName("fd")
        fmt.Println(fd.Int()) // "1"
        fd.SetInt(2)          // panic: unexported field
        ```
    - 一个可取地址的reflect.Value会记录一个结构体成员是否是未导出成员，如果是的话则拒绝修改操作。因此，CanAddr方法并不能正确反映一个变量是否是可以被修改的。另一个相关的方法CanSet是用于检查对应的reflect.Value是否是可取地址并可被修改的：
        ```go
        fmt.Println(fd.CanAddr(), fd.CanSet()) // "true false"
        ```


### **示例: 解码/反序列化S表达式**

标准库中encoding/...下每个包中提供的Marshal编码函数都有一个对应的Unmarshal函数用于解码。例如，我们在4.5节中看到的，要将包含JSON编码格式的字节slice数据解码为我们自己的Movie类型（§12.3），我们可以这样做：

```go
data := []byte{/* ... */}
var movie Movie
err := json.Unmarshal(data, &movie)

```

Unmarshal函数使用了反射机制类修改movie变量的每个成员，根据输入的内容为Movie成员创建对应的map、结构体和slice。

现在让我们为S表达式编码实现一个简易的Unmarshal，类似于前面的json.Unmarshal标准库函数，对应我们之前实现的sexpr.Marshal函数的逆操作。我们必须提醒一下，一个健壮的和通用的实现通常需要比例子更多的代码，为了便于演示我们采用了精简的实现。我们只支持S表达式有限的子集，同时处理错误的方式也比较粗暴，代码的目的是为了演示反射的用法，而不是构造一个实用的S表达式的解码器。

词法分析器lexer使用了标准库中的text/scanner包将输入流的字节数据解析为一个个类似注释、标识符、字符串面值和数字面值之类的标记。输入扫描器scanner的Scan方法将提前扫描和返回下一个记号，对于rune类型。大多数记号，比如“(”，对应一个单一rune可表示的Unicode字符，但是text/scanner也可以用小的负数表示记号标识符、字符串等由多个字符组成的记号。调用Scan方法将返回这些记号的类型，接着调用TokenText方法将返回记号对应的文本内容。

因为每个解析器可能需要多次使用当前的记号，但是Scan会一直向前扫描，所以我们包装了一个lexer扫描器辅助类型，用于跟踪最近由Scan方法返回的记号。

*gopl.io/ch12/sexpr*

```go
type lexer struct {
    scan  scanner.Scanner
    token rune // the current token
}

func (lex *lexer) next()        { lex.token = lex.scan.Scan() }
func (lex *lexer) text() string { return lex.scan.TokenText() }

func (lex *lexer) consume(want rune) {
    if lex.token != want { // NOTE: Not an example of good error handling.
        panic(fmt.Sprintf("got %q, want %q", lex.text(), want))
    }
    lex.next()
}

```

现在让我们转到语法解析器。它主要包含两个功能。第一个是read函数，用于读取S表达式的当前标记，然后根据S表达式的当前标记更新可取地址的reflect.Value对应的变量v。

```go
func read(lex *lexer, v reflect.Value) {
    switch lex.token {
    case scanner.Ident:
        // The only valid identifiers are
        // "nil" and struct field names.
        if lex.text() == "nil" {
            v.Set(reflect.Zero(v.Type()))
            lex.next()
            return
        }
    case scanner.String:
        s, _ := strconv.Unquote(lex.text()) // NOTE: ignoring errors
        v.SetString(s)
        lex.next()
        return
    case scanner.Int:
        i, _ := strconv.Atoi(lex.text()) // NOTE: ignoring errors
        v.SetInt(int64(i))
        lex.next()
        return
    case '(':
        lex.next()
        readList(lex, v)
        lex.next() // consume ')'
        return
    }
    panic(fmt.Sprintf("unexpected token %q", lex.text()))
}

```

我们的S表达式使用标识符区分两个不同类型，结构体成员名和nil值的指针。read函数值处理nil类型的标识符。当遇到scanner.Ident为“nil”是，使用reflect.Zero函数将变量v设置为零值。而其它任何类型的标识符，我们都作为错误处理。后面的readList函数将处理结构体的成员名。

一个“(”标记对应一个列表的开始。第二个函数readList，将一个列表解码到一个聚合类型中（map、结构体、slice或数组），具体类型依然于传入待填充变量的类型。每次遇到这种情况，循环继续解析每个元素直到遇到于开始标记匹配的结束标记“)”，endList函数用于检测结束标记。

最有趣的部分是递归。最简单的是对数组类型的处理。直到遇到“)”结束标记，我们使用Index函数来获取数组每个元素的地址，然后递归调用read函数处理。和其它错误类似，如果输入数据导致解码器的引用超出了数组的范围，解码器将抛出panic异常。slice也采用类似方法解析，不同的是我们将为每个元素创建新的变量，然后将元素添加到slice的末尾。

在循环处理结构体和map每个元素时必须解码一个(key value)格式的对应子列表。对于结构体，key部分对于成员的名字。和数组类似，我们使用FieldByName找到结构体对应成员的变量，然后递归调用read函数处理。对于map，key可能是任意类型，对元素的处理方式和slice类似，我们创建一个新的变量，然后递归填充它，最后将新解析到的key/value对添加到map。

```go
func readList(lex *lexer, v reflect.Value) {
    switch v.Kind() {
    case reflect.Array: // (item ...)
        for i := 0; !endList(lex); i++ {
            read(lex, v.Index(i))
        }

    case reflect.Slice: // (item ...)
        for !endList(lex) {
            item := reflect.New(v.Type().Elem()).Elem()
            read(lex, item)
            v.Set(reflect.Append(v, item))
        }

    case reflect.Struct: // ((name value) ...)
        for !endList(lex) {
            lex.consume('(')
            if lex.token != scanner.Ident {
                panic(fmt.Sprintf("got token %q, want field name", lex.text()))
            }
            name := lex.text()
            lex.next()
            read(lex, v.FieldByName(name))
            lex.consume(')')
        }

    case reflect.Map: // ((key value) ...)
        v.Set(reflect.MakeMap(v.Type()))
        for !endList(lex) {
            lex.consume('(')
            key := reflect.New(v.Type().Key()).Elem()
            read(lex, key)
            value := reflect.New(v.Type().Elem()).Elem()
            read(lex, value)
            v.SetMapIndex(key, value)
            lex.consume(')')
        }

    default:
        panic(fmt.Sprintf("cannot decode list into %v", v.Type()))
    }
}

func endList(lex *lexer) bool {
    switch lex.token {
    case scanner.EOF:
        panic("end of file")
    case ')':
        return true
    }
    return false
}

```

最后，我们将解析器包装为导出的Unmarshal解码函数，隐藏了一些初始化和清理等边缘处理。内部解析器以panic的方式抛出错误，但是Unmarshal函数通过在defer语句调用recover函数来捕获内部panic（§5.10），然后返回一个对panic对应的错误信息。

```go
// Unmarshal parses S-expression data and populates the variable
// whose address is in the non-nil pointer out.
func Unmarshal(data []byte, out interface{}) (err error) {
    lex := &lexer{scan: scanner.Scanner{Mode: scanner.GoTokens}}
    lex.scan.Init(bytes.NewReader(data))
    lex.next() // get the first token
    defer func() {
        // NOTE: this is not an example of ideal error handling.
        if x := recover(); x != nil {
            err = fmt.Errorf("error at %s: %v", lex.scan.Position, x)
        }
    }()
    read(lex, reflect.ValueOf(out).Elem())
    return nil
}

```

生产实现不应该对任何输入问题都用panic形式报告，而且应该报告一些错误相关的信息，例如出现错误输入的行号和位置等。尽管如此，我们希望通过这个例子来展示类似encoding/json等包底层代码的实现思路，以及如何使用反射机制来填充数据结构。



### **示例：获取结构体字段标签**

在4.5节我们使用构体成员标签用于设置对应JSON对应的名字。其中json成员标签让我们可以选择成员的名字和抑制零值成员的输出。在本节，我们将看到如何通过反射机制类获取成员标签。

对于一个web服务，大部分HTTP处理函数要做的第一件事情就是展开请求中的参数到本地变量中。我们定义了一个工具函数，叫params.Unpack，通过使用结构体成员标签机制来让HTTP处理函数解析请求参数更方便。

首先，我们看看如何使用它。下面的search函数是一个HTTP请求处理函数。它定义了一个匿名结构体类型的变量，用结构体的每个成员表示HTTP请求的参数。其中结构体成员标签指明了对于请求参数的名字，为了减少URL的长度这些参数名通常都是神秘的缩略词。Unpack将请求参数填充到合适的结构体成员中，这样我们可以方便地通过合适的类型类来访问这些参数。

*gopl.io/ch12/search*

```go
import "gopl.io/ch12/params"

// search implements the /search URL endpoint.
func search(resp http.ResponseWriter, req *http.Request) {
    var data struct {
        Labels     []string `http:"l"`
        MaxResults int      `http:"max"`
        Exact      bool     `http:"x"`
    }
    data.MaxResults = 10 // set default
    if err := params.Unpack(req, &data); err != nil {
        http.Error(resp, err.Error(), http.StatusBadRequest) // 400
        return
    }

    // ...rest of handler...
    fmt.Fprintf(resp, "Search: %+v\n", data)
}

```

下面的Unpack函数主要完成三件事情。第一，它调用req.ParseForm()来解析HTTP请求。然后，req.Form将包含所有的请求参数，不管HTTP客户端使用的是GET还是POST请求方法。

下一步，Unpack函数将构建每个结构体成员有效参数名字到成员变量的映射。如果结构体成员有成员标签的话，有效参数名字可能和实际的成员名字不相同。reflect.Type的Field方法将返回一个reflect.StructField，里面含有每个成员的名字、类型和可选的成员标签等信息。其中成员标签信息对应reflect.StructTag类型的字符串，并且提供了Get方法用于解析和根据特定key提取的子串，例如这里的http:"..."形式的子串。

*gopl.io/ch12/params*

```go
// Unpack populates the fields of the struct pointed to by ptr
// from the HTTP request parameters in req.
func Unpack(req *http.Request, ptr interface{}) error {
    if err := req.ParseForm(); err != nil {
        return err
    }

    // Build map of fields keyed by effective name.
    fields := make(map[string]reflect.Value)
    v := reflect.ValueOf(ptr).Elem() // the struct variable
    for i := 0; i < v.NumField(); i++ {
        fieldInfo := v.Type().Field(i) // a reflect.StructField
        tag := fieldInfo.Tag           // a reflect.StructTag
        name := tag.Get("http")
        if name == "" {
            name = strings.ToLower(fieldInfo.Name)
        }
        fields[name] = v.Field(i)
    }

    // Update struct field for each parameter in the request.
    for name, values := range req.Form {
        f := fields[name]
        if !f.IsValid() {
            continue // ignore unrecognized HTTP parameters
        }
        for _, value := range values {
            if f.Kind() == reflect.Slice {
                elem := reflect.New(f.Type().Elem()).Elem()
                if err := populate(elem, value); err != nil {
                    return fmt.Errorf("%s: %v", name, err)
                }
                f.Set(reflect.Append(f, elem))
            } else {
                if err := populate(f, value); err != nil {
                    return fmt.Errorf("%s: %v", name, err)
                }
            }
        }
    }
    return nil
}

```

最后，Unpack遍历HTTP请求的name/valu参数键值对，并且根据更新相应的结构体成员。回想一下，同一个名字的参数可能出现多次。如果发生这种情况，并且对应的结构体成员是一个slice，那么就将所有的参数添加到slice中。其它情况，对应的成员值将被覆盖，只有最后一次出现的参数值才是起作用的。

populate函数小心用请求的字符串类型参数值来填充单一的成员v（或者是slice类型成员中的单一的元素）。目前，它仅支持字符串、有符号整数和布尔型。其中其它的类型将留做练习任务。

```go
func populate(v reflect.Value, value string) error {
    switch v.Kind() {
    case reflect.String:
        v.SetString(value)

    case reflect.Int:
        i, err := strconv.ParseInt(value, 10, 64)
        if err != nil {
            return err
        }
        v.SetInt(i)

    case reflect.Bool:
        b, err := strconv.ParseBool(value)
        if err != nil {
            return err
        }
        v.SetBool(b)

    default:
        return fmt.Errorf("unsupported kind %s", v.Type())
    }
    return nil
}

```

如果我们上上面的处理程序添加到一个web服务器，则可以产生以下的会话：

```powershell
$ go build gopl.io/ch12/search
$ ./search &
$ ./fetch 'http://localhost:12345/search'
Search: {Labels:[] MaxResults:10 Exact:false}
$ ./fetch 'http://localhost:12345/search?l=golang&l=programming'
Search: {Labels:[golang programming] MaxResults:10 Exact:false}
$ ./fetch 'http://localhost:12345/search?l=golang&l=programming&max=100'
Search: {Labels:[golang programming] MaxResults:100 Exact:false}
$ ./fetch 'http://localhost:12345/search?x=true&l=golang&l=programming'
Search: {Labels:[golang programming] MaxResults:10 Exact:true}
$ ./fetch 'http://localhost:12345/search?q=hello&x=123'
x: strconv.ParseBool: parsing "123": invalid syntax
$ ./fetch 'http://localhost:12345/search?q=hello&max=lots'
max: strconv.ParseInt: parsing "lots": invalid syntax

```



### **示例：显示一个类型的方法集**

最后一个例子是使用reflect.Type来打印任意值的类型和枚举它的方法：

*gopl.io/ch12/methods*

```go
// Print prints the method set of the value x.
func Print(x interface{}) {
    v := reflect.ValueOf(x)
    t := v.Type()
    fmt.Printf("type %s\n", t)

    for i := 0; i < v.NumMethod(); i++ {
        methType := v.Method(i).Type()
        fmt.Printf("func (%s) %s%s\n", t, t.Method(i).Name,
            strings.TrimPrefix(methType.String(), "func"))
    }
}

```

reflect.Type和reflect.Value都提供了一个Method方法。每次t.Method(i)调用将一个reflect.Method的实例，对应一个用于描述一个方法的名称和类型的结构体。每次v.Method(i)方法调用都返回一个reflect.Value以表示对应的值（§6.4），也就是一个方法是帮到它的接收者的。使用reflect.Value.Call方法（我们这里没有演示），将可以调用一个Func类型的Value，但是这个例子中只用到了它的类型。

这是属于time.Duration和`*strings.Replacer`两个类型的方法：

```go
methods.Print(time.Hour)
// Output:
// type time.Duration
// func (time.Duration) Hours() float64
// func (time.Duration) Minutes() float64
// func (time.Duration) Nanoseconds() int64
// func (time.Duration) Seconds() float64
// func (time.Duration) String() string

methods.Print(new(strings.Replacer))
// Output:
// type *strings.Replacer
// func (*strings.Replacer) Replace(string) string
// func (*strings.Replacer) WriteString(io.Writer, string) (int, error)
```



### 慎用反射reflect包

反射是一个功能和表达能力都很强大的工具，但应该谨慎使用它，具体有三个原因:

- 原因1：是基于反射的代码是很脆弱的。
    - **能导致编译器报告类型错误的每种写法，在反射中都有一个对应的误用方法。**编译器在编译时就能向你报告这个错误，而反射错误则要等到执行时才以panic崩溃的方式来报告，而这可能发生在写完代码很久之后了。
    - 如readList函数尝试从输入读取一个字符串然后填充一个int类型的变量，那么调用reflect.Value.SetString就会崩溃。很多使用反射的程序都有类似的风险，所以对每一个reflect.Value都需要仔细注意它的类型、是否可寻址、是否可设置。
    - 避免这种因反射而导致的脆弱性的问题的最好方法**：将所有的反射相关的使用控制在包的内部，如果可能的话避免在包的API中直接暴露reflect.Value类型，这样可以限制一些非法输入**。
        - 如果无法做到这一点，在每个有风险的操作前指向额外的类型检查。
        - 以标准库中的代码为例，当fmt.Printf收到一个非法的操作数时，它并不会抛出panic异常，而是打印相关的错误信息。程序虽然还有BUG，但是会更加容易诊断。
            ```go
            fmt.Printf("%d %s\n", "hello", 42) // "%!d(string=hello) %!s(int=42)"
            ```
    - 反射同样降低了程序的安全性，还影响了自动化重构和分析工具的准确性，因为它们无法识别运行时才能确认的类型信息。
- 原因2：类型其实也算是某种形式的文档，而反射的相关操作则无法做静态类型检查，所以大量使用反射的代码是很难理解的。对于接受interface{}或者reflect.Value的函数，一定要写清楚期望的参数类型和其他限制条件（即不变量）。
- 原因3：**基于反射的函数会比为特定类型优化的函数慢一两个数量级**。在一个典型的程序中，大部分函数与整体性能无关，所以为了让程序更清晰可以使用反射。测试就很合适使用反射，因为大部分测试都使用小数据集。但对于关键路径上的函数，则最好避免使用反射。
