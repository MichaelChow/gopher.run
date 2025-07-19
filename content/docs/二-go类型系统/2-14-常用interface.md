---
title: "2.14 常用interface"
date: 2025-04-18T12:51:00Z
draft: false
weight: 2014
---

# 2.14 常用interface

# 标准库常用接口

### **flag.Value{}**

```go
package flag

// Value接口代表了存储在标志内的值
type Value interface {
    String() string   // String方法用于格式化标志对应的值，可用于输出命令行帮助消息。由于有了该方法，因此每个flag.Value其实也是fmt.Stringer。
    Set(string) error  // Set方法解析了传入的字符串参数并更新标志值。可以认为Set方法是String方法的逆操作，两个方法使用同样的记法规格是一个很好的实践。
}
```

```go
// flag.Duration函数创建一个time.Duration类型的标记变量，并且允许用户通过多种用户友好的方式来设置这个变量的大小，这种方式还包括和String方法相同的符号排版形式。这种对称设计使得用户交互良好。
var period = flag.Duration("period", 1*time.Second, "sleep period")
flag.Parse()
fmt.Printf("Sleeping for %v...", *period)  // fmt包调用time.Duration的String方法打印这个时间周期是以用户友好的注解方式，而不是一个纳秒数字
time.Sleep(*period)
fmt.Println()
//$ go run sleep.go -period 1m3s
//$ Sleeping for 1m3s...
```



```go
type Celsius float64
type Fahrenheit float64

func CToF(c Celsius) Fahrenheit { return Fahrenheit(c*9.0/5.0 + 32.0) }
func FToC(f Fahrenheit) Celsius { return Celsius((f - 32.0) * 5.0 / 9.0) }

func (c Celsius) String() string { return fmt.Sprintf("%g°C", c) }

// *celsiusFlag satisfies the flag.Value interface.
type celsiusFlag struct{ Celsius }

func (f *celsiusFlag) Set(s string) error {
	var unit string
	var value float64
	fmt.Sscanf(s, "%f%s", &value, &unit) // no error check needed
	switch unit {
	case "C", "°C":
		f.Celsius = Celsius(value)
		return nil
	case "F", "°F":
		f.Celsius = FToC(Fahrenheit(value))
		return nil
	}
	return fmt.Errorf("invalid temperature %q", s)
}

// CelsiusFlag函数将所有逻辑都封装在一起
func CelsiusFlag(name string, value Celsius, usage string) *Celsius {
	f := celsiusFlag{value} // Celsius字段是一个会通过Set方法在标记处理的过程中更新的变量
	// 调用Var方法将标记加入应用的命令行标记集合中，有异常复杂命令行接口的全局变量flag.CommandLine.Programs可能有几个这个类型的变量
	// 调用Var方法将一个*celsiusFlag参数赋值给一个flag.Value参数，导致编译器去检查*celsiusFlag是否有必须的方法
	flag.CommandLine.Var(&f, name, usage)
	// 它返回一个内嵌在celsiusFlag变量f中的Celsius指针给调用者
	return &f.Celsius
}
```

### **error{}**

```go
// 实际中很少直接调用errors.New函数，因为有一个更方便的封装函数fmt.Errorf，它还会处理字符串格式化。
package fmt

import "errors"

func Errorf(format string, args ...interface{}) error {
    return errors.New(Sprintf(format, args...))
}
```

```go
package errors
// 整个errors包仅只有4行

func New(text string) error { return &errorString{text} }

type errorString struct { text string }

func (e *errorString) Error() string { return e.text }

type error interface {
    Error() string
}
```



```go
// syscall包：Go语言底层系统调用API ：定义一个实现error接口的数字类型Errno，并且在Unix平台上，Errno的Error方法会从一个字符串表中查找错误消息
package syscall

type Errno uintptr // operating system error code

var errors = [...]string{
    1:   "operation not permitted",   // EPERM
    2:   "no such file or directory", // ENOENT
    3:   "no such process",           // ESRCH
    // ...
}

func (e Errno) Error() string {
    if 0 <= int(e) && int(e) < len(errors) {
        return errors[e]
    }
    return fmt.Sprintf("errno %d", e)
}
```



```go
var err error = syscall.Errno(2)  // 类型为syscall.Errno，值为2
// Errno是一个系统调用错误的高效表示方式，它通过一个有限的集合进行描述，并且它满足标准的错误接口。我们会在第7.11节了解到其它满足这个接口的类型。
fmt.Println(err.Error()) // "no such file or directory"
fmt.Println(err)         // "no such file or directory"
```



### **sort.Interface{}**

排序操作和字符串格式化一样是很多程序经常使用的操作，我们不希望每次需要的时候都重写或者拷贝这些代码。

在很多语言中，排序算法都是和序列数据类型关联，同时排序函数和具体类型元素关联。但Go语言的sort.Sort函数不会对具体的序列和它的元素做任何假设。使用了一个接口类型sort.Interface来指定通用的排序算法和可能被排序到的序列类型之间的约定。这个接口的实现由序列的具体表示和它希望排序的元素决定，序列的表示经常是一个切片。

```go
package sort

type Interface interface {
		// 一个内置的排序算法需要知道三个方法
    Len() int  // 序列的长度
    Less(i, j int) bool // i, j are indices of sequence elements  // 表示两个元素比较的结果
    Swap(i, j int)  // 一种交换两个元素的方式
}


// 对一个字符串切片进行排序，实现上述接口
type StringSlice []string
func (p StringSlice) Len() int           { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// 由于对字符串切片的排序是很常用，所以sort包提供了StringSlice类型，也提供了Strings函数能让上面这些调用简化成sort.Strings(names)
```





```go
package sort

// sort包定义了一个非导出的struct类型reverse，内嵌了一个sort.Interface
type reverse struct{ sort.Interface } // that is, sort.Interface

// reverse的Less方法：调用了内嵌的sort.Interface值的Less方法，但是通过交换索引的方式使排序结果变成逆序
func (r reverse) Less(i, j int) bool { return r.Interface.Less(j, i) }

// reverse的另外两个方法Len和Swap，隐式地由原有内嵌的sort.Interface提供


// sort.Reverse函数值使用了组合（GO面向对象的重要特性，类似90年代语言中的继承，但更优）
func Reverse(data Interface) Interface { return reverse{data} }

```



sorting

```go

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

// 变量tracks包含了一个播放列表
// 尽管可以直接存储Tracks值，但**如果每个元素都是****指针 *Track (而不是Track类型)会更快**（指针是一个机器字码长度而Track类型可能是八个或更多）
var tracks = []*Track{
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

// 打印成一个表格
func printTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "------", "-----", "----", "------")
	for _, t := range tracks {
		fmt.Fprintf(tw, format, t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
	tw.Flush() // calculate column widths and print table  格式化整个表格并且将它写向os.Stdout（标准输出）
}

// 为了能 by Artist字段 对播放列表进行排序，定义一个新的带有必须的Len，Less和Swap方法的切片类型(像对StringSlice那样)。
type byArtist []*Track

func (x byArtist) Len() int           { return len(x) }                    // 切片的长度
func (x byArtist) Less(i, j int) bool { return x[i].Artist < x[j].Artist } // 比较两个元素的大小
func (x byArtist) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }          // 交换两个元素

// !+yearcode
type byYear []*Track

func (x byYear) Len() int           { return len(x) }
func (x byYear) Less(i, j int) bool { return x[i].Year < x[j].Year }
func (x byYear) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

//!-yearcode

func main() {
	fmt.Println("byArtist:")
	// 为了调用通用的排序程序，我们必须先将tracks转换为新的byArtist类型，它定义了具体的排序：
	sort.Sort(byArtist(tracks))
	printTracks(tracks)

	fmt.Println("\nReverse(byArtist):")
	sort.Sort(sort.Reverse(byArtist(tracks))) // 这里不需要定义一个有颠倒Less方法的新类型byReverseArtist
	printTracks(tracks)

	fmt.Println("\nbyYear:")
	sort.Sort(byYear(tracks))
	printTracks(tracks)

	fmt.Println("\nCustom sort:")
	// 定义一个多层的排序函数: 它主要的排序键是标题，第二个键是年，第三个键是运行时间Length，其中这个排序使用了匿名排序函数
	sort.Sort(customSort{tracks, func(x, y *Track) bool {
		if x.Title != y.Title {
			return x.Title < y.Title
		}
		if x.Year != y.Year {
			return x.Year < y.Year
		}
		if x.Length != y.Length {
			return x.Length < y.Length
		}
		return false
	}})
	//!-customcall
	printTracks(tracks)
}


// !+customcode
type customSort struct {
	t    []*Track
	less func(x, y *Track) bool
}

func (x customSort) Len() int           { return len(x.t) }
func (x customSort) Less(i, j int) bool { return x.less(x.t[i], x.t[j]) }
func (x customSort) Swap(i, j int)      { x.t[i], x.t[j] = x.t[j], x.t[i] }
```

### **http.Handler{}**

```go
// net/http包

package http

// Handler接口，包含一个ServeHTTP方法
type Handler interface {
    ServeHTTP(w ResponseWriter, r *Request)
}

// ListenAndServe函数
func ListenAndServe(address string, h Handler) error
```

*http1.go*

```go
func main() {
    db := database{"shoes": 50, "socks": 5}
    log.Fatal(http.ListenAndServe("localhost:8000", db))
}
// dollars类型
type dollars float32

// dollars类型的String方法
func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

// database类型
type database map[string]dollars

// database类型的ServeHTTP方法  -> 实现了Handler接口
func (db database) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    for item, price := range db {
        fmt.Fprintf(w, "%s: %s\n", item, price)
    }
}
```



*http2.go*

```go
// ServeHTTP方法：在http1.go的基础上，switch case
func (db database) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    switch req.URL.Path {
    case "/list":
        for item, price := range db {
            fmt.Fprintf(w, "%s: %s\n", item, price)
        }
    case "/price":
        item := req.URL.Query().Get("item")
        price, ok := db[item]
        if !ok {
            w.WriteHeader(http.StatusNotFound) // 404
            fmt.Fprintf(w, "no such item: %q\n", item)
            return
        }
        fmt.Fprintf(w, "%s\n", price)
    default:
        w.WriteHeader(http.StatusNotFound) // 404
        fmt.Fprintf(w, "no such page: %s\n", req.URL)
    }
}
```



*http3.go*

```go
// 再看net/http包
// NewServeMux 分配并返回一个新的 ServeMux
func NewServeMux() *ServeMux {
	return &ServeMux{}
}
```

```go
func main() {
    db := database{"shoes": 50, "socks": 5}
    mux := http.NewServeMux()
    // db.list：func(w http.ResponseWriter, req *http.Request) 方法类型的值（方法值）
    // http.HandlerFunc(db.list)：http.HandlerFunc是一个类型，为类型转换
    mux.Handle("/list", http.HandlerFunc(db.list))
    mux.Handle("/price", http.HandlerFunc(db.price))
    log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
    for item, price := range db {
        fmt.Fprintf(w, "%s: %s\n", item, price)
    }
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
    item := req.URL.Query().Get("item")
    price, ok := db[item]
    if !ok {
        w.WriteHeader(http.StatusNotFound) // 404
        fmt.Fprintf(w, "no such item: %q\n", item)
        return
    }
    fmt.Fprintf(w, "%s\n", price)
}
```



*http4.go*

```go
// 再看net/http包
// http.HandlerFunc：是一个实现了接口http.Handler的方法的函数类型
// 是一个让函数值满足一个接口的适配器，这里函数和这个接口仅有的方法有相同的函数签名
type HandlerFunc func(w ResponseWriter, r *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}

// 为了方便，net/http包提供了一个全局的ServeMux实例DefaultServerMux、包级别的http.Handle、http.HandleFunc函数。
// 现在，为了使用DefaultServeMux作为服务器的主handler，我们不需要将它传给ListenAndServe函数；nil值就可以工作。

// DefaultServeMux is the default [ServeMux] used by [Serve].
var defaultServeMux ServeMux

var DefaultServeMux = &defaultServeMux

// If there is no registered handler that applies to the request,
// Handler returns a “page not found” handler and an empty pattern.
func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string) {
	if use121 {
		return mux.mux121.findHandler(r)
	}
	h, p, _, _ := mux.findHandler(r)
	return h, p
}

// HandleFunc registers the handler function for the given pattern.
// If the given pattern conflicts, with one that is already registered, HandleFunc
// panics.
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	if use121 {
		mux.mux121.handleFunc(pattern, handler)
	} else {
		mux.register(pattern, HandlerFunc(handler))
	}
}
```

```go
func main() {
    db := database{"shoes": 50, "socks": 5}
    http.HandleFunc("/list", db.list)
    http.HandleFunc("/price", db.price)
    log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
```

### **示例：创建一个表达式求值器接口**

- 假定这里的表达式语言由浮点数符号（小数点）；二元操作符+，-，*， 和/；一元操作符-x和+x；调用pow(x,y)，sin(x)，和sqrt(x)的函数；当然也有括号和标准的优先级运算符。所有的值都是float64类型。
```go
// An Expr is an arithmetic expression.
type Expr interface{}
```



*eval.go*

```go
// A Var identifies a variable, e.g., x.
type Var string

// A literal is a numeric constant, e.g., 3.141.
type literal float64

// A unary represents a unary operator expression, e.g., -x.
type unary struct {
    op rune // one of '+', '-'
    x  Expr
}

// A binary represents a binary operator expression, e.g., x+y.
type binary struct {
    op   rune // one of '+', '-', '*', '/'
    x, y Expr
}

// A call represents a function call expression, e.g., sin(x).
type call struct {
    fn   string // one of "pow", "sin", "sqrt"
    args []Expr
}

// 将变量的名字映射成对应的值
type Env map[Var]float64

type Expr interface {
    // Eval returns the value of this Expr in the environment env.
    // 根据给定的environment变量，返回表达式的值
    Eval(env Env) float64
}
```



```go
func (v Var) Eval(env Env) float64 {
    return env[v]
}

func (l literal) Eval(_ Env) float64 {
    return float64(l)
}
```



```go
// 具体实现
func (u unary) Eval(env Env) float64 {
    switch u.op {
    case '+':
        return +u.x.Eval(env)
    case '-':
        return -u.x.Eval(env)
    }
    panic(fmt.Sprintf("unsupported unary operator: %q", u.op))
}

func (b binary) Eval(env Env) float64 {
    switch b.op {
    case '+':
        return b.x.Eval(env) + b.y.Eval(env)
    case '-':
        return b.x.Eval(env) - b.y.Eval(env)
    case '*':
        return b.x.Eval(env) * b.y.Eval(env)
    case '/':
        return b.x.Eval(env) / b.y.Eval(env)
    }
    panic(fmt.Sprintf("unsupported binary operator: %q", b.op))
}

func (c call) Eval(env Env) float64 {
    switch c.fn {
    case "pow":
        return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
    case "sin":
        return math.Sin(c.args[0].Eval(env))
    case "sqrt":
        return math.Sqrt(c.args[0].Eval(env))
    }
    panic(fmt.Sprintf("unsupported function call: %s", c.fn))
}
```



```go
// 对于表格中的每一条记录，这个测试会解析它的表达式然后在环境变量中计算它，输出结果
func TestEval(t *testing.T) {
    tests := []struct {
        expr string
        env  Env
        want string
    }{
        {"sqrt(A / pi)", Env{"A": 87616, "pi": math.Pi}, "167"},
        {"pow(x, 3) + pow(y, 3)", Env{"x": 12, "y": 1}, "1729"},
        {"pow(x, 3) + pow(y, 3)", Env{"x": 9, "y": 10}, "1729"},
        {"5 / 9 * (F - 32)", Env{"F": -40}, "-40"},
        {"5 / 9 * (F - 32)", Env{"F": 32}, "0"},
        {"5 / 9 * (F - 32)", Env{"F": 212}, "100"},
    }
    var prevExpr string
    for _, test := range tests {
        // Print expr only when it changes.
        if test.expr != prevExpr {
            fmt.Printf("\n%s\n", test.expr)
            prevExpr = test.expr
        }
        expr, err := Parse(test.expr)
        if err != nil {
            t.Error(err) // parse error
            continue
        }
        got := fmt.Sprintf("%.6g", expr.Eval(test.env))
        fmt.Printf("\t%v => %s\n", test.env, got)
        if got != test.want {
            t.Errorf("%s.Eval() in %v = %q, want %q\n",
            test.expr, test.env, got, test.want)
        }
    }
}
```

静态错误就是不用运行程序就可以检测出来的错误。

通过将静态检查和动态的部分分开，我们可以快速的检查错误并且对于多次检查只执行一次而不是每次表达式计算的时候都进行检查。



[http://localhost:8000/plot?expr=pow(2,sin(y)](http://localhost:8000/plot?expr=pow(2%2Csin(y)))*pow(2,sin(x))/12

![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/efa6386e-1adb-4004-80a8-e08568a74ebd/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB466QI7XIU3K%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T010045Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJGMEQCIDNO88LyTX1ZvOoYAtTPJ%2FFxbgMq8fxMTnR1DMPNA3q5AiA6YQymu3VgknhdqWctJwnArzaZ4I7Yd2nZbGIcVw1ELSqIBAiZ%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F8BEAAaDDYzNzQyMzE4MzgwNSIMXfJcrcvPNJkPH768KtwD4dVlJzoYFt6jvPZCxwKQG8lzEr9PksvvjkcfFLprW9UJNHK1w9aHpf8gqPi%2BXlubL7gFFVzfjfgAuXimXksWR23gW8dP0yna1FFOAeURNreKPNS4dI9Px1qmHEIJsAnc8k4onq5%2BnYwawnZVR3EVEaSUJtE13PXrwyQhAXm4kjmmg9zt9ZygwOuBymtmZNrQzTGpBzJD0BZ6DXluUKrvEnCifoDsbOMFoEKIkhbpI%2BXsPmZNscvxnDwieMOue%2BBzi5zy4BoPb5nlrOjcNixk3JoTVHMKdz9ohFAno8ZUM29hmp4xlwA0005LvmBdBdzNndcs13xtgYVdgYPX4l5dDj5%2FAR3q%2Fhe0zu2jc0J2yD70sX4rnIYDaj0DkNc6tHyFbB4scG6vrBRlfm7ASapXBF%2FUSHa3K2CtIOXCeyCybHo0MmrWBvCVahWnwZUluZ9iix%2BwK%2BezG5bU6Yhws1fzotBEOwJss3KzQwSFcF8yVkwF4I6igAfAxCVksCG7GnX1ibKkaP7MrgZwXb6DWSc8aRqADPU9Dk18wrVzdz1aTpavtB%2F%2F5KbzvUwjiq%2F6y6o9QFv3MBweeyIkWIYc%2FKFJYNTUGOYUgnq3AhgoPBGG%2B7pkJQk4mbVMj9ZcfsUwxbrrwwY6pgFvsEkuXyKHFy7rlInW5SR5Jm6NbQYdkNjpOYI8gFDeoE1St5r6fGwmm9pl7soGDpoGDyzP9BoAABBltG2n9WSLaMFA3o02rTmIxnUiVDzlFluB0QF0ftgd%2FGzfuaZH6yhM0iK9ZnSmtTTRVfEaZ%2BI6XbW4JwO4rj6J9LgRqONoB3BmEtqRM6r%2FSbbF2nUDBrEH5iJIMOz9QDed9Y1fCLmu95h0adBh&X-Amz-Signature=c44e3b1617b16a00a22150901f3ff04b797ca120505a823519382b32f38dc43b&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)

![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/d0bc3ca4-aa99-42ab-9e0b-451aa45171dc/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB466QI7XIU3K%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T010045Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJGMEQCIDNO88LyTX1ZvOoYAtTPJ%2FFxbgMq8fxMTnR1DMPNA3q5AiA6YQymu3VgknhdqWctJwnArzaZ4I7Yd2nZbGIcVw1ELSqIBAiZ%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F8BEAAaDDYzNzQyMzE4MzgwNSIMXfJcrcvPNJkPH768KtwD4dVlJzoYFt6jvPZCxwKQG8lzEr9PksvvjkcfFLprW9UJNHK1w9aHpf8gqPi%2BXlubL7gFFVzfjfgAuXimXksWR23gW8dP0yna1FFOAeURNreKPNS4dI9Px1qmHEIJsAnc8k4onq5%2BnYwawnZVR3EVEaSUJtE13PXrwyQhAXm4kjmmg9zt9ZygwOuBymtmZNrQzTGpBzJD0BZ6DXluUKrvEnCifoDsbOMFoEKIkhbpI%2BXsPmZNscvxnDwieMOue%2BBzi5zy4BoPb5nlrOjcNixk3JoTVHMKdz9ohFAno8ZUM29hmp4xlwA0005LvmBdBdzNndcs13xtgYVdgYPX4l5dDj5%2FAR3q%2Fhe0zu2jc0J2yD70sX4rnIYDaj0DkNc6tHyFbB4scG6vrBRlfm7ASapXBF%2FUSHa3K2CtIOXCeyCybHo0MmrWBvCVahWnwZUluZ9iix%2BwK%2BezG5bU6Yhws1fzotBEOwJss3KzQwSFcF8yVkwF4I6igAfAxCVksCG7GnX1ibKkaP7MrgZwXb6DWSc8aRqADPU9Dk18wrVzdz1aTpavtB%2F%2F5KbzvUwjiq%2F6y6o9QFv3MBweeyIkWIYc%2FKFJYNTUGOYUgnq3AhgoPBGG%2B7pkJQk4mbVMj9ZcfsUwxbrrwwY6pgFvsEkuXyKHFy7rlInW5SR5Jm6NbQYdkNjpOYI8gFDeoE1St5r6fGwmm9pl7soGDpoGDyzP9BoAABBltG2n9WSLaMFA3o02rTmIxnUiVDzlFluB0QF0ftgd%2FGzfuaZH6yhM0iK9ZnSmtTTRVfEaZ%2BI6XbW4JwO4rj6J9LgRqONoB3BmEtqRM6r%2FSbbF2nUDBrEH5iJIMOz9QDed9Y1fCLmu95h0adBh&X-Amz-Signature=2f59fe5a1c5fdcc01c79cf7b75a493d989307d1b86274ec4e7286eef748b71ed&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)







## **示例: 基于标记的XML解码**

encoding/xml包也提供了一个更低层的基于标记的API用于XML解码。在基于标记的样式中，解析器消费输入并产生一个标记流；四个主要的标记类型－StartElement，EndElement，CharData，和Comment－每一个都是encoding/xml包中的具体类型。每一个对(*xml.Decoder).Token的调用都返回一个标记。



```go
package xml

type Name struct {
    Local string // e.g., "Title" or "id"
}

type Attr struct { // e.g., name="value"
    Name  Name
    Value string
}

// A Token includes StartElement, EndElement, CharData,
// and Comment, plus a few esoteric types (not shown).
// 这个没有方法的Token接口也是一个**可识别联合**的例子。
type Token interface{}
type StartElement struct { // e.g., <name>
    Name Name
    Attr []Attr
}
type EndElement struct { Name Name } // e.g., </name>
type CharData []byte                 // e.g., <p>CharData</p>
type Comment []byte                  // e.g., <!-- Comment -->

type Decoder struct{ /* ... */ }
func NewDecoder(io.Reader) *Decoder
func (*Decoder) Token() (Token, error) // returns next Token in sequence
```



传统的接口如io.Reader的目的是隐藏满足它的具体类型的细节，这样就可以创造出新的实现：在这个实现中每个具体类型都被统一地对待。

相反，满足可识别联合的具体类型的集合**被设计为确定和暴露****（**而不是隐藏）。可识别联合的类型几乎没有方法，操作它们的函数使用一个类型分支的switch case集合来进行表述，这个case集合中每一个switch case都有不同的逻辑。



### xmlselect.go

```go
// 获取和打印在一个XML文档树中确定的元素下找到的文本
// Xmlselect prints the text of selected elements of an XML document.
package main

import (
    "encoding/xml"
    "fmt"
    "io"
    "os"
    "strings"
)

func main() {
    dec := xml.NewDecoder(os.Stdin)
    var stack []string // stack of element names
    // 循环每遇到一个StartElement时，它把这个元素的名称压到一个栈里，并且每次遇到EndElement时，它将名称从这个栈中推出：保证了StartElement和EndElement的序列可以被完全的匹配
    for {
        tok, err := dec.Token()
        if err == io.EOF {
            break
        } else if err != nil {
            fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
            os.Exit(1)
        }
        switch tok := tok.(type) {
        case xml.StartElement:
            stack = append(stack, tok.Name.Local) // push
        case xml.EndElement:
            stack = stack[:len(stack)-1] // pop
        case xml.CharData:
            if containsAll(stack, os.Args[1:]) {
                fmt.Printf("%s: %s\n", strings.Join(stack, " "), tok)
            }
        }
    }
}

// containsAll reports whether x contains the elements of y, in order.
func containsAll(x, y []string) bool {
    for len(y) <= len(x) {
        if len(y) == 0 {
            return true
        }
        if x[0] == y[0] {
            y = y[1:]
        }
        x = x[1:]
    }
    return false
}
```





# **接口与其它类型**

### **接口**

Go中的接口为指定对象的行为提供了一种方法：如果某样东西可以完成**这个**， 那么它就可以用在**这里**。我们已经见过许多简单的示例了；通过实现 `String` 方法，我们可以自定义打印函数，而通过 `Write` 方法，`Fprintf` 则能对任何对象产生输出。在Go代码中， 仅包含一两种方法的接口很常见，且其名称通常来自于实现它的方法， 如 `io.Writer` 就是实现了 `Write` 的一类对象。

每种类型都能实现多个接口。例如一个实现了 `sort.Interface` 接口的集合就可通过 `sort` 包中的例程进行排序。该接口包括 `Len()`、`Less(i, j int) bool` 以及 `Swap(i, j int)`，另外，该集合仍然可以有一个自定义的格式化器。 以下特意构建的例子 `Sequence` 就同时满足这两种情况。

```go
type Sequence []int

// Methods required by sort.Interface.
// sort.Interface 所需的方法。
func (s Sequence) Len() int {
    return len(s)
}
func (s Sequence) Less(i, j int) bool {
    return s[i] < s[j]
}
func (s Sequence) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

// Method for printing - sorts the elements before printing.
// 用于打印的方法 - 在打印前对元素进行排序。
func (s Sequence) String() string {
    sort.Sort(s)
    str := "["
    for i, elem := range s {
        if i > 0 {
            str += " "
        }
        str += fmt.Sprint(elem)
    }
    return str + "]"
}
```

### **类型转换**

`Sequence` 的 `String` 方法重新实现了 `Sprint` 为切片实现的功能。若我们在调用 `Sprint` 之前将 `Sequence` 转换为纯粹的 `[]int`，就能共享已实现的功能。

```plain text
func (s Sequence) String() string {
	sort.Sort(s)
	return fmt.Sprint([]int(s))
}

```

该方法是通过类型转换技术，在 `String` 方法中安全调用 `Sprintf` 的另个一例子。若我们忽略类型名的话，这两种类型（`Sequence`和 `[]int`）其实是相同的，因此在二者之间进行转换是合法的。 转换过程并不会创建新值，它只是值暂让现有的时看起来有个新类型而已。 （还有些合法转换则会创建新值，如从整数转换为浮点数等。）

在Go程序中，为访问不同的方法集而进行类型转换的情况非常常见。 例如，我们可使用现有的 `sort.IntSlice` 类型来简化整个示例：

```plain text
type Sequence []int

// // 用于打印的方法 - 在打印前对元素进行排序。
func (s Sequence) String() string {
	sort.IntSlice(s).Sort()
	return fmt.Sprint([]int(s))
}

```

现在，不必让 `Sequence` 实现多个接口（排序和打印）， 我们可通过将数据条目转换为多种类型（`Sequence`、`sort.IntSlice` 和 `[]int`）来使用相应的功能，每次转换都完成一部分工作。 这在实践中虽然有些不同寻常，但往往却很有效。

### **接口转换与类型断言**

[类型选择](https://go-zh.org/doc/effective_go.html#%E7%B1%BB%E5%9E%8B%E9%80%89%E6%8B%A9)是类型转换的一种形式：它接受一个接口，在选择 （`switch`）中根据其判断选择对应的情况（`case`）， 并在某种意义上将其转换为该种类型。以下代码为 `fmt.Printf` 通过类型选择将值转换为字符串的简化版。若它已经为字符串，我们需要该接口中实际的字符串值； 若它有 `String` 方法，我们则需要调用该方法所得的结果。

```go
type Stringer interface {
	String() string
}

var value interface{} // 调用者提供的值。
switch str := value.(type) {
case string:
	return str
case Stringer:
	return str.String()
}

```

第一种情况获取具体的值，第二种将该接口转换为另一个接口。这种方式对于混合类型来说非常完美。

若我们只关心一种类型呢？若我们知道该值拥有一个 `string` 而想要提取它呢？ 只需一种情况的类型选择就行，但它需要**类型断言**。类型断言接受一个接口值， 并从中提取指定的明确类型的值。其语法借鉴自类型选择开头的子句，但它需要一个明确的类型， 而非 `type` 关键字：

```go
value.(typeName)
```

而其结果则是拥有静态类型 `typeName` 的新值。该类型必须为该接口所拥有的具体类型， 或者该值可转换成的第二种接口类型。要提取我们知道在该值中的字符串，可以这样：

```go
str := value.(string)
```

但若它所转换的值中不包含字符串，该程序就会以运行时错误崩溃。为避免这种情况， 需使用“逗号, ok”惯用测试它能安全地判断该值是否为字符串：

```go
str, ok := value.(string)
if ok {
	fmt.Printf("字符串值为 %q\n", str)
} else {
	fmt.Printf("该值非字符串\n")
}

```

若类型断言失败，`str` 将继续存在且为字符串类型，但它将拥有零值，即空字符串。

作为对能量的说明，这里有个 `if-else` 语句，它等价于本节开头的类型选择。

```go
if str, ok := value.(string); ok {
	return str
} else if str, ok := value.(Stringer); ok {
	return str.String()
}

```

### **通用性**

若某种现有的类型仅实现了一个接口，且除此之外并无可导出的方法，则该类型本身就无需导出。 仅导出该接口能让我们更专注于其行为而非实现，其它属性不同的实现则能镜像该原始类型的行为。 这也能够避免为每个通用接口的实例重复编写文档。

在这种情况下，构造函数应当返回一个接口值而非实现的类型。例如在 `hash` 库中，`crc32.NewIEEE` 和 `adler32.New` 都返回接口类型 `hash.Hash32`。要在Go程序中用Adler-32算法替代CRC-32， 只需修改构造函数调用即可，其余代码则不受算法改变的影响。

同样的方式能将 `crypto` 包中多种联系在一起的流密码算法与块密码算法分开。 `crypto/cipher` 包中的 `Block` 接口指定了块密码算法的行为， 它为单独的数据块提供加密。接着，和 `bufio` 包类似，任何实现了该接口的密码包都能被用于构造以 `Stream` 为接口表示的流密码，而无需知道块加密的细节。

`crypto/cipher` 接口看其来就像这样：

```go
type Block interface {
	BlockSize() int
	Encrypt(src, dst []byte)
	Decrypt(src, dst []byte)
}

type Stream interface {
	XORKeyStream(dst, src []byte)
}

```

这是计数器模式CTR流的定义，它将块加密改为流加密，注意块加密的细节已被抽象化了。

```go
// NewCTR 返回一个 Stream，其加密/解密使用计数器模式中给定的 Block 进行。
// iv 的长度必须与 Block 的块大小相同。
func NewCTR(block Block, iv []byte) Stream

```

`NewCTR` 的应用并不仅限于特定的加密算法和数据源，它适用于任何对 `Block` 接口和 `Stream` 的实现。因为它们返回接口值， 所以用其它加密模式来代替CTR只需做局部的更改。构造函数的调用过程必须被修改， 但由于其周围的代码只能将它看做 `Stream`，因此它们不会注意到其中的区别。

### **接口和方法**

由于几乎任何类型都能添加方法，因此几乎任何类型都能满足一个接口。一个很直观的例子就是 `http` 包中定义的 `Handler` 接口。任何实现了 `Handler` 的对象都能够处理HTTP请求。

```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

```

`ResponseWriter` 接口提供了对方法的访问，这些方法需要响应客户端的请求。 由于这些方法包含了标准的 `Write` 方法，因此 `http.ResponseWriter` 可用于任何 `io.Writer` 适用的场景。`Request` 结构体包含已解析的客户端请求。

为简单起见，我们假设所有的HTTP请求都是GET方法，而忽略POST方法， 这种简化不会影响处理程序的建立方式。这里有个短小却完整的处理程序实现， 它用于记录某个页面被访问的次数。

```go
// 简单的计数器服务。
type Counter struct {
	n int
}

func (ctr *Counter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctr.n++
	fmt.Fprintf(w, "counter = %d\n", ctr.n)
}

```

（紧跟我们的主题，注意 `Fprintf` 如何能输出到 `http.ResponseWriter`。） 作为参考，这里演示了如何将这样一个服务器添加到URL树的一个节点上。

```plain text
import "net/http"
...
ctr := new(Counter)
http.Handle("/counter", ctr)

```

但为什么 `Counter` 要是结构体呢？一个整数就够了。 An integer is all that's needed. （接收者必须为指针，增量操作对于调用者才可见。）

```go
// 简单的计数器服务。
type Counter int

func (ctr *Counter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	*ctr++
	fmt.Fprintf(w, "counter = %d\n", *ctr)
}

```

当页面被访问时，怎样通知你的程序去更新一些内部状态呢？为Web页面绑定个信道吧。

```go
// 每次浏览该信道都会发送一个提醒。
// （可能需要带缓冲的信道。）
type Chan chan *http.Request

func (ch Chan) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ch <- req
	fmt.Fprint(w, "notification sent")
}

```

最后，假设我们需要输出调用服务器二进制程序时使用的实参 `/args`。 很简单，写个打印实参的函数就行了。

```plain text
func ArgServer() {
	fmt.Println(os.Args)
}

```

我们如何将它转换为HTTP服务器呢？我们可以将 `ArgServer` 实现为某种可忽略值的方法，不过还有种更简单的方法。 既然我们可以为除指针和接口以外的任何类型定义方法，同样也能为一个函数写一个方法。 `http` 包包含以下代码：

```go
// HandlerFunc 类型是一个适配器，它允许将普通函数用做HTTP处理程序。
// 若 f 是个具有适当签名的函数，HandlerFunc(f) 就是个调用 f 的处理程序对象。
type HandlerFunc func(ResponseWriter, *Request)

// ServeHTTP calls f(c, req).
func (f HandlerFunc) ServeHTTP(w ResponseWriter, req *Request) {
	f(w, req)
}

```

`HandlerFunc` 是个具有 `ServeHTTP` 方法的类型， 因此该类型的值就能处理HTTP请求。我们来看看该方法的实现：接收者是一个函数 `f`，而该方法调用 `f`。这看起来很奇怪，但不必大惊小怪， 区别在于接收者变成了一个信道，而方法通过该信道发送消息。

为了将 `ArgServer` 实现成HTTP服务器，首先我们得让它拥有合适的签名。

```go
// 实参服务器。
func ArgServer(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, os.Args)
}

```

`ArgServer` 和 `HandlerFunc` 现在拥有了相同的签名， 因此我们可将其转换为这种类型以访问它的方法，就像我们将 `Sequence` 转换为 `IntSlice` 以访问 `IntSlice.Sort` 那样。 建立代码非常简单：

```plain text
http.Handle("/args", http.HandlerFunc(ArgServer))

```

当有人访问 `/args` 页面时，安装到该页面的处理程序就有了值 `ArgServer` 和类型 `HandlerFunc`。 HTTP服务器会以 `ArgServer` 为接收者，调用该类型的 `ServeHTTP` 方法，它会反过来调用 `ArgServer`（通过 `f(c, req)`），接着实参就会被显示出来。

在本节中，我们通过一个结构体，一个整数，一个信道和一个函数，建立了一个HTTP服务器， 这一切都是因为接口只是方法的集和，而几乎任何类型都能定义方法。

## **接口内嵌**

Go并不提供典型的，类型驱动的子类化概念，但通过将类型<**内嵌**到结构体或接口中， 它就能“借鉴”部分实现。

接口内嵌非常简单。我们之前提到过 `io.Reader` 和 `io.Writer` 接口，这里是它们的定义。

```plain text
type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

```

`io` 包也导出了一些其它接口，以此来阐明对象所需实现的方法。 例如 `io.ReadWriter` 就是个包含 `Read` 和 `Write` 的接口。我们可以通过显示地列出这两个方法来指明 `io.ReadWriter`， 但通过将这两个接口内嵌到新的接口中显然更容易且更具启发性，就像这样：

```plain text
// ReadWriter 接口结合了 Reader 和 Writer 接口。
type ReadWriter interface {
	Reader
	Writer
}

```

正如它看起来那样：`ReadWriter` 能够做任何 `Reader` **和** `Writer` 可以做到的事情，它是内嵌接口的联合体 （它们必须是不相交的方法集）。只有接口能被嵌入到接口中。

同样的基本想法可以应用在结构体中，但其意义更加深远。`bufio` 包中有 `bufio.Reader` 和 `bufio.Writer` 这两个结构体类型， 它们每一个都实现了与 `io` 包中相同意义的接口。此外，`bufio` 还通过结合 `reader/writer` 并将其内嵌到结构体中，实现了带缓冲的 `reader/writer`：它列出了结构体中的类型，但并未给予它们字段名。

```plain text
// ReadWriter 存储了指向 Reader 和 Writer 的指针。
// 它实现了 io.ReadWriter。
type ReadWriter struct {
	*Reader  // *bufio.Reader
	*Writer  // *bufio.Writer
}

```

内嵌的元素为指向结构体的指针，当然它们在使用前必须被初始化为指向有效结构体的指针。 `ReadWriter` 结构体和通过如下方式定义：

```plain text
type ReadWriter struct {
	reader *Reader
	writer *Writer
}

```

但为了提升该字段的方法并满足 `io` 接口，我们同样需要提供转发的方法， 就像这样：

```plain text
func (rw *ReadWriter) Read(p []byte) (n int, err error) {
	return rw.reader.Read(p)
}

```

而通过直接内嵌结构体，我们就能避免如此繁琐。 内嵌类型的方法可以直接引用，这意味着 `bufio.ReadWriter` 不仅包括 `bufio.Reader` 和 `bufio.Writer` 的方法，它还同时满足下列三个接口： `io.Reader`、`io.Writer` 以及 `io.ReadWriter`。

还有种区分内嵌与子类的重要手段。当内嵌一个类型时，该类型的方法会成为外部类型的方法， 但当它们被调用时，该方法的接收者是内部类型，而非外部的。在我们的例子中，当 `bufio.ReadWriter` 的 `Read` 方法被调用时， 它与之前写的转发方法具有同样的效果；接收者是 `ReadWriter` 的 `reader` 字段，而非 `ReadWriter` 本身。

内嵌同样可以提供便利。这个例子展示了一个内嵌字段和一个常规的命名字段。

```plain text
type Job struct {
	Command string
	*log.Logger
}

```

`Job` 类型现在有了 `Log`、`Logf` 和 `*log.Logger` 的其它方法。我们当然可以为 `Logger` 提供一个字段名，但完全不必这么做。现在，一旦初始化后，我们就能记录 `Job` 了：

```plain text
job.Log("starting now...")

```

`Logger` 是 `Job` 结构体的常规字段， 因此我们可在 `Job` 的构造函数中，通过一般的方式来初始化它，就像这样：

```plain text
func NewJob(command string, logger *log.Logger) *Job {
	return &Job{command, logger}
}

```

或通过复合字面：

```plain text
job := &Job{command, log.New(os.Stderr, "Job: ", log.Ldate)}

```

若我们需要直接引用内嵌字段，可以忽略包限定名，直接将该字段的类型名作为字段名， 就像我们在 `ReaderWriter` 结构体的 `Read` 方法中做的那样。 若我们需要访问 `Job` 类型的变量 `job` 的 `*log.Logger`， 可以直接写作 `job.Logger`。若我们想精炼 `Logger` 的方法时， 这会非常有用。

```plain text
func (job *Job) Logf(format string, args ...interface{}) {
	job.Logger.Logf("%q: %s", job.Command, fmt.Sprintf(format, args...))
}

```

内嵌类型会引入命名冲突的问题，但解决规则却很简单。首先，字段或方法 `X` 会隐藏该类型中更深层嵌套的其它项 `X`。若 `log.Logger` 包含一个名为 `Command` 的字段或方法，`Job` 的 `Command` 字段会覆盖它。

其次，若相同的嵌套层级上出现同名冲突，通常会产生一个错误。若 `Job` 结构体中包含名为 `Logger` 的字段或方法，再将 `log.Logger` 内嵌到其中的话就会产生错误。然而，若重名永远不会在该类型定义之外的程序中使用，那就不会出错。 这种限定能够在外部嵌套类型发生修改时提供某种保护。 因此，就算添加的字段与另一个子类型中的字段相冲突，只要这两个相同的字段永远不会被使用就没问题。



