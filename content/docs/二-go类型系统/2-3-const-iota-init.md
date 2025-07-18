---
title: "2.3 const、iota、init"
date: 2025-04-03T02:02:00Z
draft: false
weight: 2003
---

# 2.3 const、iota、init

- Go中的常量就是**不变量（值不可修改）**，这样**可以防止在运行期被意外或恶意的修改**，类型只能是**Number、boolean、string/rune**；
    ```go
    const (
    	E   = 2.71828182845904523536028747135266249775724709369995957496696763 // https://oeis.org/A001113
    	Pi  = 3.14159265358979323846264338327950288419716939937510582097494459 // https://oeis.org/A000796
    	Phi = 1.61803398874989484820458683436563811772030917980576286213544862 // https://oeis.org/A001622
    )
    ```
- 如果是批量声明的常量，除了第一个外，**可省略常量右边的初始化表达式**，**表示使用前面常量的初始化表达式写法和常量类型；可省略类型，由初始化表达式推断；**
    ```go
    const (
        a = 1
        b      // 1
        c = 2
        d     // 2
    )
    ```
- 所有常量及其运算**在编译时创建**，可减少运行时的工作，也方便其他编译优化。**定义它们的表达式必须也是可被编译器求值的常量表达式**。如 `1<<3` 就是一个常量表达式，而 `math.Sin(math.Pi/4)` 则不是，**因为对 **`**math.Sin**`** 的****函数调用在运行时才会发生；常量运算**如整数除零、字符串索引越界、任何导致无效浮点数的操作等运行时错误可以在编译时发现；
    - 常量间的所有算术运算、逻辑运算和比较运算的结果也是常量，：len、cap、real、imag、complex和unsafe.Sizeof返回的结果为常量
    - 因为它们的值是在编译期就确定的，因此常量可以是构成类型的一部分。（如用于指定数组类型的长度）
        ```go
        const IPv4Len = 4
        func parseIPv4(s string) IP {
            var p [IPv4Len]byte
        }
        ```
- **Go中枚举常量： 使用枚举器/常量生成器**`**iota**`** 初始化，**由于`iota`可为表达式的一部分，而**表达式可以被隐式地重复**，所以可用于生成一组以相似规则初始化的常量，但是不用每行都写一遍初始化表达式；
    - iota 是go语言的常量计数器，只能在常量的表达式中使用，表示**从0开始自动加1**，所以Sunday=0，Monday=1，以此类推。
        ```go
        // 首先定义一个Weekday命名类型，然后为一周的每天定义了一个常量，从周日0开始。
        type Weekday int
        const (
        	Sunday Weekday = iota
        	Monday
        	Tuesday
        	Wednesday
        	Thursday
        	Friday
        	Saturday
        )
        ```
        - x<<n左移运算 等价于 乘以 2^n，**用****0填充****右边空缺的bit位**
            ```go
            // 每个常量都是1024的幂
            const (
            	_   = 1 << (10 * iota)  // 通过赋予空白标识符来忽略第一个值
            	
            	KiB      // 1024   , 1 << (10 * 1)
            	MiB     // 1048576 , 1 << (10 * 2)
            	GiB     // 1073741824, 1 << (10 * 3)
            	TiB     // 1099511627776, 1 << (10 * 4)   (exceeds 1 << 32)
            	PiB     // 1125899906842624, 1 << (10 * 5)
            	EiB     // 1152921504606846976, 1 << (10 * 6)
            	ZiB     // ZiB和YiB的值已经超出int64，但是它们依然是合法的常量  1180591620717411303424, 1 << (10 * 7)    (exceeds 1 << 64) 溢出
            	YiB     // 1208925819614629174706176, 1 << (10 * 8) 溢出
            	fmt.Println(YiB/ZiB) // "1024"    // 而且像下面的常量表达式依然有效（译注：YiB/ZiB是在编译期计算出来的，并且结果常量是1024，是Go语言int变量能有效表示的）：
            )
            ```
            ```go
            // 给一个无符号整数的最低5bit的每个bit指定一个名字
            // 使用这些常量可以用于测试、设置或清除对应的bit位的值
            type Flags uint
            const (
            	FlagUp           Flags = 1 << iota // 1 << 0，1 * 2^0
            	FlagBroadcast                      // 1 << 1，1 * 2^1
            	FlagLoopback                       // 1 << 2，1 * 2^2
            	FlagPointToPoint                   // 1 << 3，1 * 2^3
            	FlagMulticast                      // 1 << 4，1 * 2^4
            )
            ```
- 由于可将 `String` 之类的方法附加在用户定义的类型上（如结构体）， 因此它就为打印时自动格式化任意值提供了可能性。 尽管你常常会看到这种技术应用于结构体，但它对于像 `ByteSize` 之类的浮点数标量等类型也是有用的。
    表达式 `YB` 会打印出 `1.00YB`，而 `ByteSize(1e13)` 则会打印出 `9.09`。
    在这里用 `Sprintf` 实现 `ByteSize` 的 `String` 方法很安全（不会无限递归），这倒不是因为类型转换，而是它以 `%f` 调用了 `Sprintf`，它并不是一种字符串格式：`Sprintf` 只会在它需要字符串时才调用 `String` 方法，而 `%f` 需要一个浮点数值。
    ```go
    func (b ByteSize) String() string {
        switch {
        case b >= YB:
            return fmt.Sprintf("%.2fYB", b/YB)
        case b >= ZB:
            return fmt.Sprintf("%.2fZB", b/ZB)
        case b >= EB:
            return fmt.Sprintf("%.2fEB", b/EB)
        case b >= PB:
            return fmt.Sprintf("%.2fPB", b/PB)
        case b >= TB:
            return fmt.Sprintf("%.2fTB", b/TB)
        case b >= GB:
            return fmt.Sprintf("%.2fGB", b/GB)
        case b >= MB:
            return fmt.Sprintf("%.2fMB", b/MB)
        case b >= KB:
            return fmt.Sprintf("%.2fKB", b/KB)
        }
        return fmt.Sprintf("%.2fB", b)
    }
    ```




### **无类型常量**

- Go中6种**未明确类型的常量类型，**编译器为这些没有明确基础类型的数字常量提供比基础类型更高精度的算术运算（你可以认为至少有**256bit的运算精度**）
    - 无类型的布尔型
    - 无类型的整数
    - 无类型的字符、
    - 无类型的浮点数
    - 无类型的复数
    - 无类型的字符串
- 通过延迟明确常量的具体类型，无类型的常量不仅可以提供更高的运算精度，而且可以直接用于更多的表达式而不需要显式的类型转换。
    ```go
    // math.Pi无类型的浮点数常量，可以直接用于任意需要浮点数或复数的地方
    var x float32 = math.Pi
    var y float64 = math.Pi
    var z complex128 = math.Pi
    // 对于常量面值，不同的写法可能会对应不同的类型。
    // 如0、0.0、0i和\u0000虽然有着相同的常量值，但是它们分别对应无类型的整数、无类型的浮点数、无类型的复数和无类型的字符等不同的常量类型。
    // true和false也是无类型的布尔类型，字符串面值常量是无类型的字符串类型。
    var f float64 = 212
    fmt.Println((f - 32) * 5 / 9)     // "100"; (f - 32) * 5 is a float64
    fmt.Println(5 / 9 * (f - 32))     // "0";   5/9 is an untyped integer, 0
    fmt.Println(5.0 / 9.0 * (f - 32)) // "100"; 5.0/9.0 is an untyped float
    // 只有常量可以是无类型的。
    // 当一个无类型的常量被赋值给一个变量的时候，就像下面的第一行语句，或者出现在有明确类型的变量声明的右边，如下面的其余三行语句，无类型的常量将会被隐式转换为对应的类型，如果转换合法的话。
    var f float64 = 3 + 0i // untyped complex -> float64
    f = 2                  // untyped integer -> float64
    f = 1e123              // untyped floating-point -> float64
    f = 'a'                // untyped rune -> float64
    // 上面的语句相当于:
    var f float64 = float64(3 + 0i)
    f = float64(2)
    f = float64(1e123)
    f = float64('a'
    // 无论是隐式或显式转换，将一种类型转换为另一种类型都要求目标可以表示原始值。
    // 对于浮点数和复数，可能会有舍入处理：
    const (
        deadbeef = 0xdeadbeef // untyped int with value 3735928559
        a = uint32(deadbeef)  // uint32 with value 3735928559
        b = float32(deadbeef) // float32 with value 3735928576 (rounded up)
        c = float64(deadbeef) // float64 with value 3735928559 (exact)
        d = int32(deadbeef)   // compile error: constant overflows int32
        e = float64(1e309)    // compile error: constant overflows float64
        f = uint(-1)          // compile error: constant underflows uint
    )
    // 对于一个没有显式类型的变量声明（包括简短变量声明），常量的形式将隐式决定变量的默认类型
    i := 0      // untyped integer;        implicit int(0)
    r := '\000' // untyped rune;           implicit rune('\000')
    f := 0.0    // untyped floating-point; implicit float64(0.0)
    c := 0i     // untyped complex;        implicit complex128(0i)
    // 注意有一点不同：无类型整数常量转换为int，它的内存大小是不确定的，但是无类型浮点数和复数常量则转换为内存大小明确的float64和complex128。 
    // 如果不知道浮点数类型的内存大小是很难写出正确的数值算法的，因此Go语言不存在整型类似的不确定内存大小的浮点数和复数类型。
    // 如果要给变量一个不同的类型，我们必须显式地将无类型的常量转化为所需的类型，或给声明的变量指定明确的类型，像下面例子这样：
    var i = int8(0)
    var i int8 = 0
    // 当尝试将这些无类型的常量转为一个接口值时，这些默认类型将显得尤为重要，因为要靠它们明确接口对应的动态类型。
    fmt.Printf("%T\n", 0)      // "int"
    fmt.Printf("%T\n", 0.0)    // "float64"
    fmt.Printf("%T\n", 0i)     // "complex128"
    fmt.Printf("%T\n", '\000') // "int32" (rune)
    ```


### init

- 变量的初始化与常量类似，**但其初始值也可以是在运行时才被计算的一般表达式**。
    ```go
    var (
    	home   = os.Getenv("HOME")
    	user   = os.Getenv("USER")
    	gopath = os.Getenv("GOPATH")
    )
    ```
### `**init**`** 函数**

- 每个源文件都可以通过定义自己的无参数 `init` 函数来设置一些必要的状态。 （其实每个文件都可以拥有多个 `init` 函数。）
- 而它的结束就意味着初始化结束： 只有该包中的所有变量声明都通过它们的初始化器求值后 `init` 才会被调用， 而那些 `init` 只有在所有已导入的包都被初始化后才会被求值。
- 除了那些不能被表示成声明的初始化外，`init` 函数还常被用在程序真正开始执行前，检验或校正程序的状态。
```go
func init() {
	if user == "" {
		log.Fatal("$USER not set")
	}
	if home == "" {
		home = "/home/" + user
	}
	if gopath == "" {
		gopath = home + "/go"
	}
	// gopath 可通过命令行中的 --gopath 标记覆盖掉。
	flag.StringVar(&gopath, "gopath", gopath, "override default GOPATH")
}
```







