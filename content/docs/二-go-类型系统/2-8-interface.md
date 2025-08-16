---
title: "2.8 interface"
date: 2024-12-28T15:08:00Z
draft: false
weight: 2008
---

# 2.8 interface



interface type、实现、interface变量的值、interface的类型断言

## interface即约定

- **interface类型是对具体类型行为的概括与抽象。通过使用接口，我们可以写出更加灵活和通用的函数，这些函数不用绑定在一个特定的类型实现上。**
    - 具体类型指定了它所含数据的精确布局，还暴露了基于这个精确布局的内部操作。如对于数值有算术操作，对于slice类型我们有索引、append、range等操作。
    - 具体类型还会通过其方法来提供额外的能力。
    - **接口是一种抽象类型，它并没有暴露所含数据的布局或者内部结构，当然也没有那些数据的基本操作，它所提供的仅仅是一些方法而已。**如果你拿到一个接口类型的值，你无从知道它是什么 is a，你能知道的仅仅是它能做什么 has a，或者更精确地讲，仅仅是它提供了哪些方法。
- **Go的interface类型是****隐式实现的**：**一个具体的类型无须声明它实现了哪些接口，只要提供接口所必需的方法即可。**
    - 无须改动已有类型的实现，就可以为这些类型创建新的接口，对于那些不能修改包的类型，这一点特别有用。
- **接口即约定：io.Writer接口定义了Fprintf和调用者之间的约定。**
    - **约定要求调用者提供具体类型（如*os.File、*bytes.Buffer），包含一个与其签名和行为一致的Write方法。**
    - **约定保证了Fprintf能使用任何满足io.Writer接口的参数。Fprintf只需要能调用参数的Write函数，无须假设它写入的是一个文件还是一段内存。**
        - **可取代性(substitutability)：可以把一种类型替换为满足同一接口的另一种类型的特性**，也是面向对象语言的典型特征。
        ```go
        package io
        type Writer interface {
            // Implementations must not retain p.
            Write(p []byte) (n int, err error)
        }
        ```
    - 通过使用io.Writer的interface，不必把Fprintf函数的格式化这个最困难的代码笨拙的复制一份
        ```go
        package fmt
        func Fprintf(w io.Writer, format string, args ...interface{}) (int, error)
        func Printf(format string, args ...interface{}) (int, error) {
            return Fprintf(os.Stdout, format, args...)
        }
        func Sprintf(format string, args ...interface{}) string {
            var buf bytes.Buffer
            Fprintf(&buf, format, args...)
            return buf.String()
        }
        ```
    - ByteCounter实现了io.Writer interface
        ```go
        type ByteCounter int
        func (c *ByteCounter) Write(p []byte) (int, error) {
        	*c += ByteCounter(len(p)) // convert int to ByteCounter
        	return len(p), nil
        }
        ```
- 当设计一个新包时，一个新手Go程序员会首先创建一系列接口，然后再定义满足这些接口的具体类型：这种方式会产生很多接口，但这些接口只有一个单独的实现。不要这样做。**这种接口是不必要的抽象，还有运行时的成本。**可以用导出机制（参考6.6节）来限制一个类型的哪些方法或结构体的哪些字段是对包外可见的。仅在有两个或者多个具体类型需要按统一的方式处理时才需要接口。
    - 这个规则也有特例，如果接口和类型实现出于依赖的原因不能放在同一个包里边，那么一个接口只有一个具体类型实现也是可以的。在这种情况下，**接口是一种解耦两个包的好方式**。
- **因为Go中仅在有两个或者多个类型满足的情况下才使用接口**，所以它就必然会抽象掉那些特有的实现细节。这种设计的结果就是出现了具有更简单和更少方法的接口，如io.Writer和fmt.Stringer都只有一个方法。**设计新类型时越小的接口越容易满足。一个不错的接口设计标准是： 仅要求你需要的方法。（****ask only for what you need）**
- Go语言能很好地支持面向对象编程风格，但这并不意味着你只能使用它。Go中**不是所有东西都必须是一个对象，全局函数应该有它们的位置，不完全封装的数据类型也应该有位置。**
## **interface类型**

- **一个接口类型定义了一套方法**；如果一个具体类型要实现该接口，那么**必须实现接口类型定义中的所有方法**。
    - **io.Writer负责所有可以写入字节的类型的抽象：文件、内存缓冲区、网络连接、HTTP客户端、打包器(archiver)、散列器(hasher)等**。
    - Reader抽象了所有可以读取字节的类型，Closer抽象了所有可以关闭的类型; 如文件或者网络连接; **Go中单方法接口的命名习惯: xxer**
        ```go
        package io
        type Reader interface {
            Read(p []byte) (n int, err error)
        }
        type Closer interface {
            Close() error
        }
        ```
    - 已有接口组合嵌套得到的新接口: 
        ```go
        type ReadWriter interface {
            Reader
            Writer
        }
        type ReadWriteCloser interface {
            Reader
            Writer
            Closer
        }
        ```
## **实现interface**

- 如果一个类型实现了一个**interface**要求的所有方法集合，那么这个类型实现了这个接口。
    - `*os.File`类型实现了io.Reader、Writer、Closer、ReadWriter interface
    - `*bytes.Buffer`实现了Reader、Writer、ReadWriter interface，但是它没有实现Closer interface因为它不具有Close方法；
- 为了简化表述，**通常说一个具体类型“是一个”(is–a)特定的接口类型**，这其实代表着该具体类型实现了该接口。
    - 如：*bytes.Buffer是一个io.Writer;
    - *os.File是一个io.ReaderWriter。
- **interface的赋值规则：仅当一个表达式实现了一个接口时，这个表达式才可以赋给该接口**。当右侧表达式也是一个interface时，该规则同样有效：
    ```go
    var w io.Writer
    w = os.Stdout           // OK: *os.File 有Write方法
    w = new(bytes.Buffer)   // OK: *bytes.Buffer 有Write方法
    w = time.Second         // 编译错误: time.Duration 缺少Write 方法
    var rwc io.ReadWriteCloser
    rwc = os.Stdout         // OK: *os.File 有Read, Write, Close方法
    rwc = new(bytes.Buffer) // 编译错误:  *bytes.Buffer 没有Close方法
    ```
    ```go
    w = rwc                 // OK: io.ReadWriteCloser 有Write方法
    rwc = w
    ```
- **一个类型有某一个方法**：
    ```go
    type IntSet struct { /* ... */ }
    // 接收者是一个指针类型
    func (*IntSet) String() string
    // 所以不能在一个不能寻址的IntSet值上调用这个方法
    var _ = IntSet{}.String() // compile error: String requires *IntSet receiver
    // 但可以在一个IntSet变量上(可寻址)调用这个方法：
    var s IntSet
    var _ = s.String() // OK: s is a variable and &s has a String method
    ```
- **interface封装了具体类型和数据：只有通过interface暴露的方法才可以调用，具体类型的其他方法则无法通过接口来调用：**
    - 一个拥有更多方法的**interface**，给了我们它所指向数据的更多信息，**当然也给实现这个interface提出更高的门槛**。如io.ReadWriter，与io.Reader相比；
- **空接口类型 interface{}**：它完全不包含任何方法，对具体类型没有任何要求，所以我们可以把任何值赋给空接口类型，用来接受任意类型的参数；
    ```go
    // **builtin.go**: any is an alias for interface{} and is equivalent to interface{} in all ways.
    type any = interface{}
    ```
    ```go
    func Println(a ...any) (n int, err error) {    // 用来接受任何类型的参数
    	return Fprintln(os.Stdout, a...)
    }
    ```
    - 即使我们创建了一个指向布尔值、浮点数、字符串、map、指针或者其他类型的interface{}接口，也无法直接使用其中的值，毕竟这个接口不包含任何方法。我们需要一个方法从空接口中使用类型断言来还原出实际值；
    ```go
    any = true 
    any = 12.34
    any = "hello"
    any = map[string]int{"one": 1}
    any = new(bytes.Buffer)
    ```
- **判定是否实现接口只需要比较具体类型和接口类型的方法，所以没必要在具体类型的定义中声明这种关系。断言 *byte.Buffer必须实现 io.Writer：**
    ```go
    var _ io.Writer = (*bytes.Buffer)(nil)
    ```
- 实际开发中，方法的接收者通常都为结构体的指针类型，因为通常需要setattr等给接收者直接赋值行为；
- **从具体类型出发，提取其共性而得出的每一种分组方式都可以表示为一种接口类型。**
    - 与基于类的语言（它们显式地声明了一个类型实现的所有接口）不同的是，**Go可以在需要时才定义新的抽象和分组，并且不用修改原有类型的定义**。**当需要使用另一个作者写的包里的具体类型时，这一点特别有用。**当然，还需要这些具体类型在底层是真正有共性的。
    ```go
    type Text interface {
        Pages() int
        Words() int
        PageSize() int
    }
    type Audio interface {
        Stream() (io.ReadCloser, error)
        RunningTime() time.Duration
        Format() string // e.g., "MP3", "WAV"
    }
    type Video interface {
        Stream() (io.ReadCloser, error)
        RunningTime() time.Duration
        Format() string // e.g., "MP4", "WMV"
        Resolution() (x, y int)
    }
    // **可以后期再定义一个Streamer接口，来代表Audio、Video之间相同的部分，而不必对已经开发上线的这两个类型做改变**
    type Streamer interface {
        Stream() (io.ReadCloser, error)
        RunningTime() time.Duration
        Format() string
    }
    ```
# interface value

- 对于像Go这样的静态类型语言，类型仅仅是一个编译时的概念，所以类型不是一个值。在我们的概念模型中，用**类型描述符**来提供每个类型的具体信息（如类型的名字和方法）。
- 从概念上讲，**一个interface变量的值**有两部分：
    - 一个具体的类型（**接口的动态类型**）：用对应的类型描述符来表述；
    - 该类型的一个值（**接口的动态值**）；
- 初始化为interface类型的零值，即将动态类型和值都设置为nil。
    ```go
    var w io.Writer
    w.Write([]byte("hello")) // panic: nil pointer dereference
    ```
    - 一个接口值是否是nil取决于它的动态类型，所以现在这是一个nil接口值。w==nil或者w! =nil来检测一个接口值是否是nil。
    - 调用一个nil接口的任何方法都会导致崩溃：
    ![](/images/16a24637-29b5-80ae-98f6-d5d27bf5867c/image_1d824637-29b5-8011-bb7b-e3b4c31fdd6e.jpg)
- 如下赋值**把一个具体类型隐式转换为一个interface类型**，等价的显式转换代码：io.Writer(os.Stdout)
    ```go
    w = os.Stdout
    w.Write([]byte("hello")) // 实际调用：(*os.File).Write方法 
    ```
    - 接口值的动态类型会设置为: 指针类型*os.File的类型描述符，动态值会设置为os.Stdout的副本，即一个指向代表进程的标准输出的os.File类型的指针 fd int = 1(stdout)
    ![](/images/16a24637-29b5-80ae-98f6-d5d27bf5867c/image_1d824637-29b5-8029-a330-e340036436e0.jpg)
    - 一般来讲，在编译时我们无法知道一个接口值的动态类型会是什么，所以**通过接口来做调用必然需要使用动态分发。**编译器必须生成一段代码来从类型描述符拿到名为Write的方法地址，再间接调用该方法地址。调用的接收者就是接口值的动态值，即os.Stdout，所以实际效果与直接调用等价：
        ```go
        os.Stdout.Write([]byte("hello")) 
        ```
    - 备注：在RUST中通过所有权系统管理内存，非GC方式，可以在编译时检测拦截这类错误）
- 把一个*bytes.Buffer类型的值赋给了接口值
    ```go
    w = new(bytes.Buffer)
    w.Write([]byte("hello"))
    ```
    - 动态类型现在是*bytes.Buffer，动态值现在则是一个指向新分配缓冲区的指针
    ![](/images/16a24637-29b5-80ae-98f6-d5d27bf5867c/image_1d824637-29b5-8055-bd73-fd73e63ca67c.jpg)
    - 所以调用的是(*bytes.Buffer).Write方法，方法的接收者是缓冲区的地址
- 把nil赋值给w，把动态类型和动态值都设置为nil
    ```go
    w = nil
    ```
- 一个接口值可以指向多个任意大的动态值；持有time.Time类型的接口值：
    ![](/images/16a24637-29b5-80ae-98f6-d5d27bf5867c/image_1d824637-29b5-80e0-9f7b-e734aa5f616b.jpg)
- 调试过程中，可使用fmt包的%T动作（使用反射来获取接口动态类型的名称）得到接口值的动态类型；
    ```go
    var w io.Writer
    fmt.Printf("%T\n", w) // "<nil>"
    w = os.Stdout
    fmt.Printf("%T\n", w) // "*os.File"
    w = new(bytes.Buffer)
    fmt.Printf("%T\n", w) // "*bytes.Buffer"
    ```
- **接口值可比较的，用**==和!＝比较操作符，所以可以作为map的键或者作为switch语句的操作数
    - 如果两个接口值都是nil值或动态类型和值都相同，两个接口值相等。
    - **如果两个接口值的动态类型相同，但是这个动态值是不可比较的（如切片），比较时会panic；**
        ```go
        var x insterface = []int{1,2,3}
        x==x    // panic
        ```
    - 从这点来看，接口类型是非平凡的。其他类型要么是可以安全比较的（比如基础类型和指针），要么是完全不可比较的（比如slice、map和函数）；但当比较接口值或者其中包含接口值的聚合类型时，我们必须小心崩溃的可能性。当把接口作为map的键或者switch语句的操作数时，也存在类似的风险。仅在能确认接口值包含的动态值可以比较时，才比较接口值。
- 常见Bug易混淆：**nil 接口（不包含任何信息）不等于 仅仅动态值为nil的接口值**
    - 以下代码当debug设置为true时，主函数收集函数f的输出到一个缓冲区中
    - 当设置debug为false时，我们会觉得仅仅是不再收集输出，但实际上会导致程序在调用out.Write时panic
    - 尽管一个空的*bytes.Buffer指针拥有的方法满足了该接口，但它并不满足该接口所需的一些行为。这个调用违背了(*bytes.Buffer).Write的一个隐式的前置条件，即接收者不能为空，所以把空指针赋给这个接口就是一个错误。
    - 解决方案是把main函数中的buf类型修改为io.Writer，从而避免在最开始就把一个功能不完整的值赋给一个接口。
    ![](/images/16a24637-29b5-80ae-98f6-d5d27bf5867c/image_1d824637-29b5-809f-91bb-d1b0310f2a36.jpg)
    ```go
    const debug = true
    var buf *bytes.Buffer // **bug代码**
    var buf io.Writer // **修复代码**
    if debug {
       buf = new(bytes.Buffer) // enable collection of output
    }
    f(buf) // NOTE: subtly incorrect!
    if debug {
           // ...use buf...
    }
    func f(out io.Writer) {   // bug代码：动态类型：*bytes.Buffer，动态值: nil，**非nil接口**。修复代码：**io.Writer、 nil**
        if out != nil {      // 防御性检查，非nil接口，仅动态值为nil，进入if
            out.Write([]byte("done!\n"))  // 动态分发机制实际调用：(*bytes.Buffer).Write，但此时接受者值为nil； **panic**: runtime error:  nil pointer dereference
    	  }
    }
    ```
# **interface变量的类型断言**

- 类型断言是一个作用在接口值上的操作，写出来类似于`x.(T)`；其中x是一个接口类型的表达式，而T是一个类型（称为断言类型）。**类型断言会检查作为操作数的动态类型是否满足指定的断言类型**。
    - 如果**断言为**一个**具体类型T**：类型断言会检查x的动态类型**是否就是T**。
        - 如果检查成功，类型断言的结果就是x的动态值，类型当然就是T。换句话说，类型断言就是用来从它的操作数中把具体类型的值提取出来的操作。
        - 如果检查失败，那么操作panic崩溃；
        ```go
        var w io.Writer
        w = os.Stdout
        f := w.(*os.File)      // success: f == os.Stdout  // 断言成功：值是x的动态值，类型是T
        c := w.(*bytes.Buffer) // panic: interface holds *os.File, not *bytes.Buffer
        ```
    - 如果断言为一个**接口类型T**：类型断言检查x的动态类型**是否满足T**。
        - 如果检查成功，动态值并没有提取出来，结果仍然是一个接口值，接口值的类型和值部分也没有变更，只是结果的类型为接口类型T。换句话说，类型断言是一个接口值表达式，从一个接口类型变为拥有另外一套方法的接口类型（通常方法数量是增多），但保留了接口值中的动态类型和动态值部分。
        - 如下类型断言代码中，w和rw都持有os.Stdout，于是所有对应的动态类型都是
            ```go
            var w io.Writer
            w = os.Stdout
            rw := w.(io.ReadWriter) // success: *os.File has both Read and Write  // 断言成功：一个有相同动态类型和值部分的接口值，但是结果为类型T
            w = new(ByteCounter)
            rw = w.(io.ReadWriter) // panic: *ByteCounter has no Read method
            ```
    - **如果操作数是一个空接口值，类型断言都失败**。很少需要从一个接口类型向一个要求更宽松的类型做类型断言，该宽松类型的接口方法比原类型的少，而且是其子集。**因为**除了在操作nil之外的情况下，在其他情况下**这种操作与赋值一致**。
        ```go
        w = rw             // io.ReadWriter is assignable to io.Writer
        w = rw.(io.Writer) // fails only if rw == nil
        ```
- 我们经常无法确定一个接口值的动态类型，这时就需要检测它是否是某一个特定类型。
    - **如果类型断言出现在需要两个结果的赋值表达式（比如下的代码）中，那么断言不会在失败时崩溃**，而是会多返回一个布尔型的ok变量返回值来指示断言是否成功。
    - 如果操作失败，ok为false，而第一个返回值为断言类型的零值，在这个例子中就是*bytes.Buffer的空指针。
        ```go
        var w io.Writer = os.Stdout
        f, ok := w.(*os.File)      // **返回第二个结果：断言成功标志ok** success:  ok, f == os.Stdout
        b, ok := w.(*bytes.Buffer) // failure，此时不会panic: !ok, b == nil
        // **if语句的扩展格式让这个变的很简洁**
        if f, ok := w.(*os.File); ok {
            // ...use f...
        }
        ```
    - 当类型断言的操作数是一个变量时，有时你会看到返回值的名字与操作数变量名一致，原有的值就被新的返回值掩盖了，如：
        ```go
        if **w**, ok := **w**.(*os.File); ok {
            // ...use w...
        }
        ```
### **使用类型断言来识别错误**

- I/O会因为很多原因失败，但有三类原因通常必须单独处理：文件已存储（创建操作）、文件没找到（读取操作）、权限不足，对应三个帮助函数：
    ```go
    package os
    func IsExist(err error) bool  // 文件已经存在（对于创建操作）
    func IsNotExist(err error) bool  // 找不到文件（对于读取操作）
    func IsPermission(err error) bool // 权限拒
    ```
    - 一个幼稚的实现会通过检查错误消息关键字：在单元测试中还算够用，但对于生产级的代码则远远不够，同样类型的错误可能会用完全不同的错误消息来报告。
        ```go
        func IsNotExist(err error) bool {   
            return strings.Contains(err.Error(), "file does not exist") // 缺乏经验的实现: 检查错误消息是否包含了特定的子字符串
        }
        ```
    - 更可靠的方法：用专门的类型(如PathError)来表示结构化的错误值，使用类型断言来检查错误类型（远远多于一个简单的字符串的细节）：
        - 如果错误消息已被fmt.Errorf这类的方法合并到一个大字符串中，那么PathError的结构信息就丢失了。错误识别通常必须在失败操作发生时马上处理，而不是等到错误消息返回给调用者之后。
        ```go
        package os
        type PathError struct {
            Op   string
            Path string
            Err  error
        }
        func (e *PathError) Error() string {
            return e.Op + " " + e.Path + ": " + e.Err.Error()
        }
        func IsNotExist(err error) bool {
            if pe, ok := err.(*PathError); ok {   // 类型断言
                err = pe.Err
            }
            return err == syscall.ENOENT || err == ErrNotExist
        }
        _, err := os.Open("/no/such/file") // err: "open /no/such/file: No such file or directory"
        fmt.Println(os.IsNotExist(err)) // "true"
        ```
### **通过接口类型断言来查询特性**

- 必须将str转换为[]byte(字节slice)，需要进行内存分配和内存复制，又会被马上抛弃，性能低下；
    ```go
    func writeHeader(w io.Writer, contentType string) error {
        if _, err := w.Write([]byte("Content-Type: ")); err != nil {  
            return err
        }
        if _, err := w.Write([]byte(contentType)); err != nil {
            return err
        }
    }
    ```
- 深入net/http包查看，可以看到w对应的动态类型还支持一个能高效写入字符串的WriteString方法，这个方法避免了临时内存的分配和复制；
- 需要前置判断w的动态类型是否有这个WriteString方法：可以定义一个新的接口，只包含WriteString方法；然后使用类型断言来判断w的动态类型是否满足这个新接口；
    - 为了避免代码重复，我们把检查挪到了工具函数writeString()中。实际上，标准库提供了io.WriteString，而且这也是向io.Writer写入字符串的推荐方法；
        ```go
        func writeString(w io.Writer, s string) (n int, err error) {
            type stringWriter interface {               // 声明一个stringWriter{}
                WriteString(string) (n int, err error)
            }
            if sw, ok := w.(stringWriter); ok {     // 使用接口的类型断言判断是否满足该接口
                return sw.WriteString(s) // avoid a copy
            }
            return w.Write([]byte(s)) // allocate temporary copy
        }
        func writeHeader(w io.Writer, contentType string) error {
            if _, err := writeString(w, "Content-Type: "); err != nil {
                return err
            }
            if _, err := writeString(w, contentType); err != nil {
                return err
            }
        }
        ```
    - **前面提到，Go中接口的实现是隐式的：一个具体的类型是否满足stringWriter接口仅仅由它拥有的方法来决定，Go中没有一个具体类型与一个接口类型之间的显示关系声明**: 这个例子中比较古怪的地方是并没有一个标准的接口定义了WriteString方法并且指定它应满足的规范；即如果一个类型满足下面的接口，那么WriteString(s)必须与Write([]byte(s))等效。
        - 尽管io.WriteString文档中提到了这个假定，但在调用它的函数的文档中就很少提到这个假定了。给一个特定类型多定义一个方法，就隐式地接受了一个特性约定。**Go语言的初学者，特别是那些具有强类型语言背景的人，会对这种缺乏显式约定的方式感到不安，但在实践中很少产生问题。除了空接口interface{}，接口类型很少意外巧合地被实现。**
        ```go
        interface {
            io.Writer
            WriteString(s string) (n int, err error)
        }
        ```
- 前面的writeString函数使用类型断言来判定一个更普适接口类型的值(io.Writer)是否满足一个更专用的接口类型(stringWriter)，如果满足，则可以使用后者所定义的方法。这种技术不仅适用于io.ReadWriter这种标准接口，还适用于stringWriter这种自定义类型。这个方法也用在了fmt.Printf中，用于从通用类型中识别出error或者fmt.Stringer：
    ```go
    package fmt
    func formatOneValue(x interface{}) string {  // 把单个操作数转换为一个字符串
        if err, ok := x.(error); ok {
            return err.Error()
        }
        if str, ok := x.(Stringer); ok {
            return str.String()
        }
        // ...all other types...
    }
    ```
### **类型分支type switch 简化一长串的类型断言if-else语句**

接口有两种不同的风格：

- **强调方法，而不是具体类型****：接口上的各种方法突出了满足这个接口的具体类型之间的相似性，但隐藏了各个具体类型的布局和各自特有的功能。 （子类型多态(subtype polymorphism)）**
    - **如io.Reader, io.Writer, fmt.Stringer, sort.Interface, http.Handler, error**
- **强调满足这个接口的具体类型，而不是这个接口的方法（何况经常没有），也不注重信息隐藏**：充分利用了接口值能够容纳各种具体类型的能力，它把接口作为这些类型的联合(union)来使用。类型断言用来在运行时区分这些类型并分别处理。把这种风格的接口使用方式称为**可识别联合(discriminated union)**。**（特设多态(ad hoc polymorphism)）**
    ```go
    import "database/sql"
    func listTracks(db sql.DB, artist string, minYear, maxYear int) {
        result, err := db.Exec(
            "SELECT * FROM tracks WHERE artist = ? AND ? <= year AND year <= ?",  // 使用SQL字面量替换在查询字符串中的每个'?'，可避免SQL注入攻击
            artist, minYear, maxYear)
        // ...
    }
    func sqlQuote(x interface{}) string {  		// 将每个参数值转为对应的SQL字面量
        if x == nil {
            return "NULL"
        } else if _, ok := x.(int); ok {
            return fmt.Sprintf("%d", x)
        } else if _, ok := x.(uint); ok {
            return fmt.Sprintf("%d", x)
        } else if b, ok := x.(bool); ok {
            if b {
                return "TRUE"
            }
            return "FALSE"
        } else if s, ok := x.(string); ok {
            return sqlQuoteString(s)       // 转义
        } else {
            panic(fmt.Sprintf("unexpected type %T: %v", x, x))
        }
    }
    ```
    - 类型分支的最简单形式与普通分支语句类似，两个的差别是操作数改为**x.(type)**（注意：这里直接写关键词type，而不是一个特定类型），每个分支是一个或者多个类型。
        - 类型分支的分支判定基于接口值的动态类型，其中nil分支需要x==nil，而default分支则在其他分支都没有满足时才运行。
        - 注意，在原来的代码中，bool和string分支的逻辑需要访问由类型断言提取出来的原始值。这个需求比较典型，所以类型分支语句也有一种扩展形式，它用来把每个分支中提取出来的原始值绑定到一个新的变量x（与类型断言类似，重用变量名也很普遍）：
        - 与switch语句类似，类型分支也隐式创建了一个词法块，所以声明一个新变量叫x并不与外部块中的变量x冲突。每个分支也会隐式创建各自的词法块。
        - 用类型分支的扩展形式重写后的sqlQuote就更加清晰易读了：
        - 尽管sqlQuote支持任意类型的实参，但仅当实参类型能够符合类型分支中的一个时才能正常运行到结束，对于其他情况就会崩溃并抛出一条“unexpected type”（非期望类型）消息。**表面上x的类型是interface{}，实际上我们把它当作int、uint、bool、string和nil的一个可识别联合。**
        ```go
        func sqlQuote(x interface{}) string {
            switch x := x.(type) {
            case nil:
                return "NULL"
            case int, uint:
                return fmt.Sprintf("%d", x) // x has type interface{} here.
            case bool:
                if x {
                    return "TRUE"
                }
                return "FALSE"
            case string:
                return sqlQuoteString(x) // (not shown)
            default:
                panic(fmt.Sprintf("unexpected type %T: %v", x, x))   // panic
            }
        }
        ```
