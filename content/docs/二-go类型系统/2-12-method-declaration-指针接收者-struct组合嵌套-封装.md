---
title: "2.12 method declaration、指针接收者、struct组合嵌套、封装"
date: 2024-12-28T15:07:00Z
draft: false
weight: 2012
---

# 2.12 method declaration、指针接收者、struct组合嵌套、封装



- 从90年代初开始，面向对象编程（OOP）就成为了称霸工程界和教育界的编程范式，之后几乎所有大规模被应用的语言支持了OOP，Go语言也不例外。
- **对Go而言，对象就是简单的一个值或者变量，并且拥有其方法，而方法是某种特定类型的函数**。**面向对象编程就是使用方法来描述每个数据结构的属性和操**作，于是使用者不需要了解对象本身的实现。（封装和组合）
- **Go明确区分****方法和接口****：**
    - 方法是和命名类型关联的一类函数。Go比较特殊的是**方法可以被关联到任意一种命名类型**。
    - **接口是一种抽象类型**，这种类型可以让我们以同样的方式来处理不同的固有类型，不用关心它们的具体实现，而只需要关注它们提供的方法。
# **方法声明**

- 方法的声明和普通函数的声明类似，只是在函数名字前面多了一个参数，这个参数把这个方法绑定到这个参数对应的类型上。
    - 附加的参数称为方法的接收者(receiver)，它源自早先的面向对象语言，用来描述主调方法就像向对象发送消息；
    - 相当于为这种类型定义了一个独占的方法
    - Go的接收者不使用特殊名（如this或者self），而是像其他的参数变量一样命名。由于接收者会频繁地使用，最常用的方法就是取类型名称的首字母，保持简短和一致性（如Point中的p）。
    - 使用方法的第一个好处：命名可以比函数更简短，且省略包的名字；
        ```go
        package geometry
        type Point struct{ X, Y float64 }
        func (p Point) Distance(q Point) float64 {    // Distance Point类型的方法声明，名字为：Point.Distance
        	return math.Hypot(q.X-p.X, q.Y-p.Y)         // p.Distance的表达式叫做选择器（selector）
        }
        func Distance(p, q Point) float64 {   // Distance 包级别的函数声明，名字为：geometry.Distance
        	return math.Hypot(q.X-p.X, q.Y-p.Y)
        }
        perim := geometry.Path{{1, 1}, {5, 1}, {5, 4}, {1, 1}}
        perim.Distance()           // 方法限定在类型内，名字可以比包级别函数更简短
        geometry.PathDistance(perim)  // 包级别函数名字通常需要加上类型名字来避免歧
        ```
    - **除了为struct声明方法，****Go同时可以很方便地为简单类型（如number、string、slice、map、甚至function等）声明附加的行为；同一个package下的除pointer和interface以外的任何类型，都可以声明方法。**（声明需要在同一个package下，int等类型的package无法声明到）
        ```go
        type Path []Point
        func (path Path) Distance() float64 {
            sum := 0.0
            for i := range path {
                if i > 0 {
                    sum += path[i-1].Distance(path[i])
                }
            }
            return sum
        }
        ```
### 指针接收者的方法

- 在函数中，由于函数调用时会赋值每一个实参变量，当实参太大或者需要更新一个变量时，为了避免不必要的内存开销和拷贝时间开销，此时必须使用指针来传递变量的地址；在方法中同样绑定为指针类型的方法接收者（如p *Point）；
    - 该方法的名字是(*Point).ScaleBy，不加括号的表达式会被错误的解析为*(Point.ScaleBy)
    - **如果Point有任何一个方法使用指针接收者，那么所有的Point方法都应该统一使用指针接收者（**即使有些方法并不一定需要）；
    - 方法接收者声明的类型只能是2种：命名类型(Point)、指向它们的指针(*Point)。**为防止混淆，不允许本身是指针的类型进行方法声明，interface类型也不允许；**
    - **实际开发中，方法的接收者通常都为结构体的指针类型，因为通常需要setattr等给接收者直接赋值行为；**
    ```go
    func (p *Point) ScaleBy(factor float64) {  
        p.X *= factor
        p.Y *= factor
    }
    ```
- **编译器对方法接收者的实参和形参的隐式转换：****实参的类型会隐式转换为形参的类型**
    - 实参接收者和形参接收者是同一个类型（如都是T类型 或 都是*T类型）：无需隐式转换；
        ```go
        Point{1, 2}.Distance(q) //  Point
        pptr.ScaleBy(2)         // *Point
        ```
    - 实参接收者是T类型的变量，而形参接收者是*T类型：编译器会隐式地获取变量的地址；
        ```go
        p.ScaleBy(2) // implicit (&p)  **引用符号&**
        ```
    - 实参接收者是*T类型而形参接收者是T类型：编译器会隐式地解引用接收者，获得实际的取值；
        ```go
        pptr.Distance(q) // implicit (*pptr) **解引用符号***
        ```
- **Nil是一个合法的接收者：**
    - 像一些函数允许nil指针作为实参，方法的接收者也一样；尤其是当nil是类型中有意义的零值（如map和slice类型）时，如下面的nil代表空链表
    - **当定义一个类型允许nil作为接收者时，应当在文档注释中显示地标明**；
    ```go
    // *IntList的类型nil代表空列表
    type IntList struct {
        Value int
        Tail  *IntList
    }
    // Sum返回列表元素的总和
    func (list *IntList) Sum() int {
        if list == nil {
            return 0
        }
        return list.Value + list.Tail.Sum()
    }
    ```
## **通过struct组合嵌套成新类型**

- **struct**内嵌可以使我们定义字段特别多的复杂类型：可以将字段先按小类型分组，然后定义小类型的方法，最后再把它们**组合**
    - 熟悉基于类的面向对象编程语言的读者会误认为：~~ColoredPoint~~~~**is a**~~~~Point；Point类型就是ColoredPoint类型的基类，而ColoredPoint则作为子类或派生类~~；
    - **但ColoredPoint****has a Point**，Distance有一个形参Point, q不是Point，因此虽然q有一个内嵌的Point字段，但是必须显式地使用它。尝试直接传递q作为参数会报错：
        ```go
        type Point struct{ X, Y float64 }
        type ColoredPoint struct {
            Point    // 可以直接认为通过嵌入的字段就是ColoredPoint自身的字段
            Color color.RGBA   
        }
        var q = ColoredPoint{Point{5, 4}, blue}
        ```
    - 匿名字段类型可以是个指向命名类型的指针:字段和方法间接地来自于所指向的对象，可以让我们共享通用的结构以及使对象之间的关系更加动态、多样化
        ```go
        type ColoredPoint struct {
            *Point
            Color color.RGBA
        }
        p := ColoredPoint{&Point{1, 1}, red}
        q := ColoredPoint{&Point{5, 4}, blue}
        ```
    - **编译器解析一个选择器到方法的顺序(从外层到内层，类似先调用字类的方法，再调用父类的同名方法)**，如果选择器有二义性的话编译器会报错（如同一级里有两个同名的方法）:
        1. 先查找直接在这个类型里声明的方法；（第1层）
        1. 再查找从来自ColoredPoint的内嵌字段的方法；（第2层）
        1. 最后查找Point和RGBA中内嵌字段的方法。（第3层..）
    - 示例：
        - 使用两个包级别变量
            ```go
            var (
                mu sync.Mutex // 互斥锁，保护mapping
                mapping = make(map[string]string)
            )
            func Lookup(key string) string {
                mu.Lock()
                v := mapping[key]
                mu.Unlock()
                return v
            }
            ```
        - 优化版本：**组合嵌套成一个struct，sync.Mutex的Lock和Unlock方法也都被引入到了这个匿名结构中**
            ```go
            var cache = struct {
                sync.Mutex      // 组合嵌套
                mapping map[string]string
            }{
                mapping: make(map[string]string),
            }
            func Lookup(key string) string {
                cache.Lock()
                v := cache.mapping[key]
                cache.Unlock()
                return v
            }
            ```
## **方法变量和方法表达式**

- **方法变量**(method value):****通常在一个表达式里同时使用和调用方法，但也可以分开。p.Distance叫作“选择器”，选择器会返回一个**方法变量（类似 函数变量）**->一个将方法（Point.Distance）绑定到特定接收器变量的函数
    ```go
    p := Point{1, 2}
    q := Point{4, 6}
    distanceFromP := p.Distance        // method value
    distanceFromP(q)     // "5"
    ```
- **方法表达式**(method expression): 写成T.f、（*T).f；其中T是类型，是一种函数变量，把原来方法的接收者替换成函数的第一个形参，**因此它可以像平常的函数变量一样调用**。
    ```go
    distance := Point.Distance   // **方法表达式**，类型为：func(Point, Point) float64
    distance(p, q))  // "5" 和直接使用选择器相比，**使用方法表达式需要用第一个额外参数来指定接收器 （方法降级为类似的函数）**
    ```
- 方法表达式的使用场景之一：根据选择来调用接收器各不相同的方法
    ```go
    type Point struct{ X, Y float64 }
    func (p Point) Add(q Point) Point { return Point{p.X + q.X, p.Y + q.Y} }
    func (p Point) Sub(q Point) Point { return Point{p.X - q.X, p.Y - q.Y} }
    type Path []Point
    func (path Path) TranslateBy(offset Point, add bool) {
        var op func(p, q Point) Point  // **变量op，函数值类型，代表加法、减法等操作**
        if add {
            op = Point.Add   // 赋值Add函数
        } else {
            op = Point.Sub  // 赋值Add函数
        }
        for i := range path {
            // Call either path[i].Add(offset) or path[i].Sub(offset).
            path[i] = op(path[i], offset)
        }
    }
    ```
# **封装**

- **封装（数据隐藏）的变量或方法**：对象的变量或方法，不能通过对象访问到的，**对调用方是不可见；**
    - Go只有一种方式控制命名的**包外可见性**：名字首字母大写包外可导出，否则包外不可导出。**所以要封装一个对象，必须使用结构体**。
        - 与封装相反的是，但有时候需要暴露一些内部内容。Go语言也允许导出的字段。**当然，一旦导出就必须要面对API的兼容问题，因此最初的决定需要慎重，要考虑到之后维护的复杂程度，将来发生变化的可能性，以及变化对原本代码质量的影响等。**
    - 无论是在函数内的代码还是方法内的代码，**结构体类型内的字段对于同一个包中的所有代码默认都是包内可见的**。（**Go中封装的单元是包而不是类型**）
- 封装提供了三个优点：
    - 不再需要更多的语句用来检查变量的值（因为封装后调用方不能直接修改对象的变量值）；
    - 能防止使用方依赖的属性发生改变，使得设计者可以更加灵活地改变API的实现而不破坏兼容性；
        ```go
        type Buffer struct {
            buf     []byte
            initial [64]byte
            /* ... */
        }
        // Grow 方法按需扩展缓冲区的大小，保证n个字节的空间
        func (b *Buffer) Grow(n int) {
            if b.buf == nil {
                b.buf = b.initial[:0] // 最初使用预分配的空间
            }
            if len(b.buf)+n > cap(b.buf) {
                buf := make([]byte, b.Len(), 2*cap(b.buf) + n)
                copy(buf, b.buf)
                b.buf = buf
            }
        }
        ```
    - 能防止使用者肆意地改变对象内的变量（因为对象的变量只能被同一个包内的函数修改）；
        ```go
        // 允许调用方来增加counter变量的值c.n，并且允许将这个值reset为0，但是不允许随便设置这个值（译注：因为压根就访问不到）
        type Counter struct { n int }
        func (c *Counter) N() int     { return c.n }
        func (c *Counter) Increment() { c.n++ }
        func (c *Counter) Reset()     { c.n = 0 }
        ```
    - **仅仅用来获得或者修改内部变量的函数称为getter**/'gɛtɚ/ **和setter**ˈsetər/ ；Go在命名getter方法的时候，通常将Get前缀省略。这个简洁的命名习惯也同样适用在其他冗余的前缀上，比如Fetch、Find和Lookup。
        ```go
        package log
        type Logger struct {
            flags  int
            // ...
        }
        // func (l *Logger) GetFlags() int
        func (l *Logger) Flags() int
        func (l *Logger) SetFlags(flag int)
        ```
