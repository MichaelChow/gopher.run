---
title: "4.2 示例：Display，一个递归的值打印器"
date: 2025-05-17T14:07:00Z
draft: false
weight: 4002
---

# 4.2 示例：Display，一个递归的值打印器

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


