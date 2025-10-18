---
title: "2.9 error与panic"
date: 2025-04-04T01:23:00Z
draft: false
weight: 2009
---

# 2.9 error与panic

# error

- **Go的排错式编程风格**：
    - Go的错误处理有特定的规律，进行错误检查之后，检测到失败的情况往往都在成功之前。如果检测到的失败导致函数返回，成功的逻辑一般不会放在else块中而是在外层的作用域中。
    - **函数会有一种通常的形式，就是在开头有一连串的检查用来返回错误，之后跟着实际的函数体一直到最后。**
- **错误处理是包的API设计的重要部分，发生错误知识许多预料行为中的一种而已****：**
    - 有些函数总是成功返回的，尽管还有灾难性的和不可预知的场景(如运行时的内存耗尽），这类错误的表现和起因相差甚远，而且恢复的希望也很渺茫；如：**strings.Contains、strconv.FormatBool**，对所有可能的参数变量都有定义好的返回结果，不会调用失败)；
    - 其他一些函数只要符合其前置条件就能成功返回；如：time.Date始终会利用年月等构成time.Time，但**如果最后一个参数（时区）为nil会引发panic异常的宕机，标志着这是一个明显的bug**，应避免这样调用代码。
    - 对于许多其他函数，即使在高质量的代码中，也不能保证一定能够成功返回，因为有些因素不受开发者的掌控；**如：任何操作I/O的函数都一定会面对可能的错误，很多可靠的操作都可能会毫无征兆的发生错误**；
- 如果当函数调用发生错误时，返回一个附加的结果作为错误值，作为最后一个结果返回。
    - **错误的结果类型通常为error：nil代表成功；非nil代表错误，且非nil的错误类型有一个错误消息字符串；输出错误消息可通过调用它的Error方法、fmt.Println(err)、fmt.Printf(”%v”,err)**
        ```go
        type error interface {
        	Error() string
        }
        ```
        - 当函数返回非nil的error时，其他的返回值都将是undefined的而且应该忽略。但有些函数在调用出错的情况下会返回部分有用的结果。如：在读取一个文件的时候发生错误，调用Read函数后会返回成功读取的字节数和错误值，先处理不完整的返回结果再处理错误；（因此文档中需要清晰的说明返回值的含义）
        - 如，`os.Open` 实现了error接口，返回一个 `os.PathError`。包含了出错的文件名、操作和触发的操作系统错误，即便在产生该错误的调用 和输出的错误信息相距甚远时，它也会非常有用，这比苍白的“不存在该文件或目录”更具说明性。
            ```go
            // PathError 记录一个错误以及产生该错误的路径和操作。
            type PathError struct {
            	Op string    // "open"、"unlink" 等等。
            	Path string  // 相关联的文件。
            	Err error    // 由系统调用返回。
            }
            func (e *PathError) Error() string {
            	return e.Op + " " + e.Path + ": " + e.Err.Error()
            }
            ```
            `PathError`的 `Error` 会生成如下错误信息：
            ```plain text
            open /etc/passwx: no such file or directory
            ```
            错误字符串应尽可能地指明它们的来源，例如产生该错误的包名前缀。例如在 `image` 包中，由于未知格式导致解码错误的字符串为“image: unknown format”。若调用者关心错误的完整细节，可使用类型选择或者类型断言来查看特定错误，并抽取其细节。 对于 `PathErrors`，它应该还包含检查内部的 `Err` 字段以进行可能的错误恢复。
            ```go
            for try := 0; try < 2; try++ {
            	file, err = os.Create(filename)
            	if err == nil {
            		return
            	}
            	if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOSPC {
            		deleteTempFiles()  // 恢复一些空间。
            		continue
            	}
            	return
            }
            ```
            这里的第二条 `if` 是另一种[类型断言](https://go-zh.org/doc/effective_go.html#%E6%8E%A5%E5%8F%A3%E8%BD%AC%E6%8D%A2)。若它失败， `ok` 将为 `false`，而 `e` 则为`nil`. 若它成功，`ok` 将为 `true`，这意味着该错误属于 `*os.PathError` 类型，而 `e` 能够检测关于该错误的更多信息。
    - 如果错误只有一种情况，错误的结果降级为bool类型的ok；
        ```go
        // cache.Lookup失败的唯一原因是key不存在
        value, ok := cache.Lookup(key)
        if !ok {
            // ...cache[key] does not exist…
        }
        ```
    **错误是值：与许多其他语言不同，Go的常规错误处理使用 普通的值 而非 异常 来报告错误；而Go的panic异常只用于针对程序bug导致的预期之外的错误；Go使用常规的控制流机制（if、return）处理错误逻辑，这要求更加小心谨慎，而这恰恰是设计的要点；**
    - **原因：****exception****异常会嵌入带有错误消息的控制流去处理它，通常会导致预期外的结果；错误会以复杂、难以理解的、无用的、无法帮助定位错误的堆栈跟踪信息报告给最终用户（大多是关于程序结构方面的，而不是简单明了的错误信息）；**
    - **Go与Java的异常处理代码对比：**Java默认返回了难懂的堆栈，Go返回了简洁清晰的err描述。如：[AES256/CBC/PKCS5Padsing解密](https://developer.open-douyin.com/m/docs/resource/zh-CN/local-life/develop/OpenAPI/preparation/decrypt)
### **错误处理策略**

1. 传递错误：
- 将错误原文传递下去:
    ```go
    resp, err := http.Get(url)
    if err != nil{
        return nil, err
    }
    ```
- 为原始的错误消息添加额外的上下文，提供一个从最根本问题到总体故障的清晰因果链：
    ```go
    doc, err := html.Parse(resp.Body)
    resp.Body.Close()
    if err != nil {
        return nil, fmt.Errorf("parsing %s as HTML: %v", url ,err)  // fmt.Errorf函数使用fmt.Sprintf格式化错误信息并返回
    }
    ```
    - NASA的事故调查案例：genesis: crashed: no parachute: G-switch failed: bad relay orientation；
        - genesis: crashed - 主错误，表示系统/程序崩溃
        - no parachute - 次级错误，可能表示缺少安全恢复机制
        - G-switch failed - 具体故障点，G可能指goroutine或重力(G-force)
        - bad relay orientation - 根本原因，继电器方向错误
    - 因为Go的错误消息频繁地串联起来（error chain，错误链），消息字符串首字母避免使用大写和换行；错误结果可能会很长使用grep查找；
    - 设计一个错误消息的时候应当慎重，确保每一条消息的描述都是有意义的且包含充足的相关信息；并且保持一致性的形式和错误处理方式；
    - 编写错误信息时，我们要确保错误信息对问题细节的描述是详尽的。尤其是要注意错误信息表达的一致性，即相同的函数或同包内的同一组函数返回的错误在构成和处理方式上是相似的。
        - 如：os包保证每一个文件操作（比如os.Open或针对打开的文件的Read、Write或Close方法）返回的错误不仅包括错误的信息（没有权限、路径不存在等）还包含文件的名字，因此调用者在构造错误消息的时候不需要再包含这些信息
1. 重试：对于偶然的、不可预测的错误，在短暂的间隔后对操作进行重试，超出一定的重试次数和限定的时间后再报错退出。
    ```go
    func WaitForServer(url string) error {
    	const timeout = 1 * time.Minute
    	deadline := time.Now().Add(timeout)
    	for tries := 0; time.Now().Before(deadline); tries++ {
    		_, err := http.Head(url)
    		if err == nil {
    			return nil // success
    		}
    		log.Printf("server not responding (%s); retrying...", err)
    		time.Sleep(time.Second << uint(tries)) // exponential back-off
    	}
    	return fmt.Errorf("server %s failed to respond after %s", url, timeout)   // 官方这段重试代码似乎不符合Go的排错式编程风格？
    }
    ```
1. 如果错误发生后，程序无法继续运行，调用者能够向上传播错误，然后在main中输出错误后优雅地停止程序；除非遇到了bug，才能直接在库函数中结束程序。
    ```go
    if err := WaitForServer(url); err != nil {
    		fmt.Fprintf(os.Stderr, "Site is down: %v\n", err)
    		os.Exit(1)
    }
    ```
    - 一个更加方便简洁的方法是通过调用log.Fatalf实现相同的效果。和所有的日志函数一样，它默认会将时间和日期作为前缀添加到错误消息前。
        ```go
        if err := WaitForServer(url); err != nil { 
        		log.Fatalf("Site is down: %v\n", err)      // **2006/01/02 15:04:05****Site is down: no such domain: bad.gopl.io**
        		log.SetPrefix("wait: ")                   // 我们可以设置log的前缀信息屏蔽时间信息，一般而言，前缀信息会被设置成命令名
        		log.SetFlags(0)
        	}
        ```
1. 在一些错误情况下，只记录下错误信息，然后程序继续运行。
    ```go
    if err := WaitForServer(url); err != nil {
    	log.Printf("ping failed: %v; networking disabled", err)  // log包中的所有函数会为没有换行符的字符串增加换行符
    	// fmt.Fprintf(os.Stderr, "ping failed: %v; networking disabled\n", err)
    }
    ```
1. 在某些罕见的情况下，我们可以直接安全地忽略掉整个日志：
    ```go
    dir, err := os.MkdirTemp("", "scratch")
    if err != nil {
    	return fmt.Errorf("failed to create temp dir: %v", err)
    }
    // 使用临时目录
    os.RemoveAll(dir) // 忽略错误，$TMPDIR临时目录 会被操作系统周期性的删除
    ```


### **特殊的错误：文件结束标识（EOF）**

如果要从一个文件中读取n个字节的数据，调用者必须把读取到文件尾的情况，区别于遇到其他错误的操作。

- 如果n是文件本身的长度，任何错误都代表操作失败；
- 如果n小于文件的长度，调用者会重复的读取固定大小的块直到文件耗尽；
为此，io包保证任何由文件结束引起的读取错误，始终都将会得到一个与众不同的错误 `io.EOF`，定义如下。

```go
package io

import "errors"

// 当没有更多输入时，将会返回EOF
var EOF = errors.New("EOF")
```

调用者可以使用一个简单的比较操作来检测这种情况，在下面的循环中，不断从标准输入中读取字符。

对于其他错误，我们可能需要同时得到错误相关的本质原因和数量信息，因此一个固定的错误值并不能满足我们的需求，可使用类型断言来更加系统的区分某个错误值。

```go
// 因为文件结束这种错误不需要更多的描述，所以io.EOF有固定的错误信息——“EOF”

in := bufio.NewReader(os.Stdin)
for {
    r, _, err := in.ReadRune()
    if err == io.EOF {
        break // finished reading
    }
    if err != nil {
        return fmt.Errorf("read failed:%v", err)
    }
    // ...use r…
}
```



### 使用类型断言区分错误值

见7.11节

# **panic宕机**

- Go的类型系统会在编译时捕获很多错误，但运行时错误只能在运行时检查（如**数组越界访问、解引用空指针**等）。当Go运行时检测到这些错误，就会发生painc宕机；
- 一个典型的panic宕机发生时：
    - 正常的程序执行会终止
    - goroutine中的所有defer延迟函数会执行
    - 然后程序会异常退出并留下一条日志消息。日志消息包括宕机的值，这往往代表某种错误消息，每一个goroutine都会在宕机的时候显示一个函数调用的栈跟踪消息。通常可以借助这条日志消息来诊断问题的原因而不需要再一次运行该程序，因此报告一个发生宕机的程序bug时，总是会加上这条消息。
- **如果碰到逻辑上“不可能发生”的状况（如语句执行到逻辑上不可能到达的地方时），手动调用内置的panic宕机函数是最好的处理方式；**
    ```go
    switch s := suit(drawCard()); s {
    case "Spades":                                // ...
    case "Hearts":                                // ...
    case "Diamonds":                              // ...
    case "Clubs":                                 // ...
    // 如，当程序到达了某条逻辑上不可能到达的路径
    default:
        panic(fmt.Sprintf("invalid suit %q", s)) // Joker?
    }
    ```
- 设置函数的断言是一个良好的习惯，但除非你能够提供更多的有效的错误消息或者能够很快地检测出错误，否则在运行时检测断言就毫无意义、多此一举。
    ```go
    // 除非你能提供更多的错误信息，或者能更快速的发现错误
    // 否则不需要使用断言，编译器在运行时会帮你检查代码。
    func Reset(x *Buffer) {
        if x == nil {
            panic("x is nil") // unnecessary! 没必要
        }
        x.elements = nil
    }
    ```
- 尽管Go语言的panic宕机机制和其他语言的异常很相似，但宕机的使用场景不尽相同。
    - 由于宕机会引起程序异常退出，因此**只有在发生严重的错误时才会使用宕机**（如遇到与预想的逻辑不一致的代码）；
    - 用心的程序员会将所有可能会发生异常退出的情况考虑在内，以证实bug的存在。强健的代码会优雅地处理“预期的”错误（如错误的输入、配置或者I/O失败等），这时最好能够使用错误值来加以区分。
- 函数regexp.Compile编译了一个高效的正则表达式。如果调用时给的模式参数不合法则会报错，但是检查这个错误本身没有必要且相当烦琐，因为调用者知道这个特定的调用是不会失败的。在此情况下，使用宕机来处理这种不可能发生的错误是比较合理的。
    - 在程序源码中大多数正则表达式是字符串字面值（string literals），因此regexp包提供了包装函数regexp.MustCompile检查字面值输入的合法性。(函数名中的Must前缀是一种针对此类函数的命名约定，比如template.Must.)
        ```go
        func Compile(expr string) (*Regexp, error) { /* ... */ }
        func MustCompile(expr string) *Regexp {
            re, err := Compile(expr)
            if err != nil {
                panic(err)
            }
            return re
        }
        ```
    - 包装函数MustCompile 使得初始化一个包级别的正则表达式变量（带有一个编译的正则表达式）变得更加方便
        ```go
        var httpSchemeRE = regexp.MustCompile(`^https?:`) //"http:" or "https:" 正则表达式字面量由于含有大量特殊字符，通常用原生字符串字面量反引号表示，防止转义等逻辑，详见字符串章节
        ```


- 当宕机发生时，所有的延迟函数以倒序执行，从栈最上面的函数开始一直返回至main函数
    ```go
    func main() {
    	f(3)
    }
    func f(x int) {
    	fmt.Printf("f(%d)\n", x+0/x)      // f(3)、f(2)、f(1)、f(0) panics
    	defer fmt.Printf("defer %d\n", x) // defer 3、defer 2、defer 1
    	// 
    	f(x - 1)
    }
    ```
    ![](/images/1cb24637-29b5-8053-94a7-f39703c289c5/image_1d324637-29b5-8024-ac09-ec3ffdf9d440.jpg)
- runtime包提供了转储栈的方法使程序员可以诊断错误。在main函数中延迟printStack的执行：
    - 熟悉其他语言的异常机制的读者可能会对runtime.Stack能够输出函数栈信息感到吃惊，**因为栈应该已经不存在了。但事实上，Go语言的宕机机制让延迟执行的函数在栈清理之前调用。**
    ```go
    func main() {
    	// 通过在main函数中延迟调用printStack输出堆栈信息
    	defer printStack()
    	f(3)
    }
    func printStack() {
    	var buf [4096]byte
    	// 为了方便诊断问题，runtime包允许程序员输出堆栈信息
    	n := runtime.Stack(buf[:], false)
    	os.Stdout.Write(buf[:n])
    }
    func f(x int) {
    	fmt.Printf("f(%d)\n", x+0/x)      // f(3)、f(2)、f(1)、f(0) panics
    	defer fmt.Printf("defer %d\n", x) // defer 3、defer 2、defer 1，多条defer语句的执行顺序与声明顺序相反
    	f(x - 1)
    }
    ```
    ![](/images/1cb24637-29b5-8053-94a7-f39703c289c5/image_1d324637-29b5-80e7-89ab-e23d4bea92cf.jpg)
# **recover 恢复**

- **退出程序通常是正确处理宕机的方式，但在一定情况下是进行recover恢复至少可以在程序退出前理清当前混乱的情况。**
    - 如：当Web服务器遇到一个未知错误时，可以先关闭所有连接，这总比让客户端阻塞在那里要好，而在开发阶段可以向客户端返回当前遇到的错误。
    - 如果内置的recover函数在延迟函数的内部调用，而且这个包含defer语句的函数发生宕机，recover会终止当前的宕机状态并且返回宕机的值。函数不会从之前宕机的地方继续运行而是正常返回。如果recover在其他任何情况下运行，则它没有任何效果且返回nil。
    - **不加区分的恢复所有的panic异常，不是可取的做法**。因为在panic之后，无法保证包级变量的状态仍然和我们预期一致；如：对数据结构的一次重要更新没有被完整完成、文件或者网络连接没有被关闭、获得的锁没有被释放。此外，如果写日志时产生的panic被不加区分的恢复，可能会导致漏洞被忽略；
        ```go
        func Parse(input string) (s *Syntax, err error) {
        		// deferred函数帮助Parse从panic中恢复
            defer func() {
        		    // recover() 返回 panic value，赋值给err变量接收，返回给调用者
        		    // 也可以通过调用runtime.Stack往错误信息中添加完整的堆栈调用信息
                if p := recover(); p != nil {
                    err = fmt.Errorf("internal error: %v", p)
                }
            }()
            // ...parser...
        }
        ```
    - 虽然把对panic的处理都集中在一个包下，有助于简化对复杂和不可以预料问题的处理，但作为被广泛遵守的规范，你不应该试图去恢复其他包引起的panic。
    - 公有的API应该将函数的运行失败**作为error返回，而不是panic**。同样的，你也不应该恢复一个由他人开发的函数引起的panic，比如说调用者传入的回调函数，因为你无法确保这样做是安全的。但有时我们很难完全遵循规范：
        - 如：net/http包中提供了一个web服务器，将收到的请求分发给用户提供的处理函数。很显然，我们不能因为某个处理函数引发的panic异常，杀掉整个web服务器进程；
        - web服务器遇到处理函数导致的panic时会调用recover，输出堆栈信息，继续运行。这样的做法在实践中很便捷，但也会引起资源泄漏，或是因为recover操作，导致其他问题
    - 基于以上原因，安全的做法是: **只恢复应该被恢复的panic异常，****这些被****recover****恢复异常所占的比例应该****尽可能的低**
        - 为了标识某个panic是否应该被恢复，我们可以将**panic value设置成特殊类型**
        - 然后在recover时对panic value进行检查。如果发现panic value是特殊类型，就将这个panic作为error处理，如果不是，则按照正常的panic进行处理
        - 有些情况下，我们无法恢复。某些致命错误会导致Go在运行时终止程序，如内存不足。
            ```go
            // Copied from gopl.io/ch5/outline2.
            func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
            	if pre != nil {
            		pre(n)
            	}
            	for c := n.FirstChild; c != nil; c = c.NextSibling {
            		forEachNode(c, pre, post)
            	}
            	if post != nil {
            		post(n)
            	}
            }
            // !+
            // soleTitle returns the text of the first non-empty title element
            // in doc, and an error if there was not exactly one.
            func soleTitle(doc *html.Node) (title string, err error) {
              // 将panic value设置成特殊类型（作为标记）
            	type bailout struct{}
            	// deferred函数调用recover，并检查panic value
            	defer func() {
            		switch p := recover(); p {
            		case nil:
            			// no panic
            		// 当panic value是bailout{}类型时（预期发生的错误），deferred函数生成一个error返回给调用者
            		// 请注意：对可预期的错误采用了panic（应该使用error），不符合Go语言风格，这里仅是为了向读者演示这种机制
            		case bailout{}:
            			// "expected" panic
            			err = fmt.Errorf("multiple title elements")
            		// 当panic value是其他non-nil值时，表示发生了未知的panic异常，deferred函数将调用panic函数并将当前的panic value作为参数传入。等同于recover没有做任何操作
            		default:
            			panic(p) // unexpected panic; carry on panicking
            		}
            	}()
            	// Bail out of recursion if we find more than one non-empty title.
            	forEachNode(doc, func(n *html.Node) {
            		if n.Type == html.ElementNode && n.Data == "title" &&
            			n.FirstChild != nil {
            			// 如果检测到有多个<title>：调用panic，阻止函数继续递归，并将特殊类型bailout作为panic的参数
            			if title != "" {
            				panic(bailout{}) // multiple title elements
            			}
            			title = n.FirstChild.Data
            		}
            	}, nil)
            	if title == "" {
            		return "", fmt.Errorf("no title element")
            	}
            	return title, nil
            }
            //!-
            func title(url string) error {
            	resp, err := http.Get(url)
            	if err != nil {
            		return err
            	}
            	// Check Content-Type is HTML (e.g., "text/html; charset=utf-8").
            	ct := resp.Header.Get("Content-Type")
            	if ct != "text/html" && !strings.HasPrefix(ct, "text/html;") {
            		resp.Body.Close()
            		return fmt.Errorf("%s has type %s, not text/html", url, ct)
            	}
            	doc, err := html.Parse(resp.Body)
            	resp.Body.Close()
            	if err != nil {
            		return fmt.Errorf("parsing %s as HTML: %v", url, err)
            	}
            	title, err := soleTitle(doc)
            	if err != nil {
            		return err
            	}
            	fmt.Println(title)
            	return nil
            }
            func main() {
            	for _, arg := range os.Args[1:] {
            		if err := title(arg); err != nil {
            			fmt.Fprintf(os.Stderr, "title: %v\n", err)
            		}
            	}
            }
            ```


- 当 `panic` 被调用后（包括不明确的运行时错误，例如切片检索越界或类型断言失败）， 程序将立刻终止当前函数的执行，并开始回溯Go程的栈，运行任何被推迟的函数。 若回溯到达Go程栈的顶端，程序就会终止。不过我们可以用内建的 `recover` 函数来重新或来取回Go程的控制权限并使其恢复正常执行。
- 调用 `recover` 将停止回溯过程，并返回传入 `panic` 的实参。 由于在回溯时只有被推迟函数中的代码在运行，因此 `recover` 只能在被推迟的函数中才有效。
- `recover` 的一个应用就是在服务器中终止失败的Go程而无需杀死其它正在执行的Go程。
    ```go
    func server(workChan <-chan *Work) {
    	for work := range workChan {
    		go safelyDo(work)
    	}
    }
    func safelyDo(work *Work) {
    	defer func() {
    		if err := recover(); err != nil {
    			log.Println("work failed:", err)
    		}
    	}()
    	do(work)
    }
    ```
    - 在此例中，若 `do(work)` 触发了Panic，其结果就会被记录， 而该Go程会被干净利落地结束，不会干扰到其它Go程。我们无需在推迟的闭包中做任何事情， `recover` 会处理好这一切。
    - 由于直接从被推迟函数中调用 `recover` 时不会返回 `nil`， 因此被推迟的代码能够调用本身使用了 `panic` 和 `recover` 的库函数而不会失败。例如在 `safelyDo` 中，被推迟的函数可能在调用 `recover` 前先调用记录函数，而该记录函数应当不受Panic状态的代码的影响。
    - 通过恰当地使用恢复模式，`do` 函数（及其调用的任何代码）可通过调用 `panic` 来避免更坏的结果。我们可以利用这种思想来简化复杂软件中的错误处理。 让我们看看 `regexp` 包的理想化版本，它会以局部的错误类型调用 `panic` 来报告解析错误。以下是一个 `error` 类型的 `Error` 方法和一个 `Compile` 函数的定义：
        ```go
        // Error 是解析错误的类型，它满足 error 接口。
        type Error string
        func (e Error) Error() string {
        	return string(e)
        }
        // error 是 *Regexp 的方法，它通过用一个 Error 触发Panic来报告解析错误。
        func (regexp *Regexp) error(err string) {
        	panic(Error(err))
        }
        // Compile 返回该正则表达式解析后的表示。
        func Compile(str string) (regexp *Regexp, err error) {
        	regexp = new(Regexp)
        	// doParse will panic if there is a parse error.
        	defer func() {
        		if e := recover(); e != nil {
        			regexp = nil    // 清理返回值。
        			err = e.(Error) // 若它不是解析错误，将重新触发Panic。
        		}
        	}()
        	return regexp.doParse(str), nil
        }
        ```
- 若 `doParse` 触发了Panic，恢复块会将返回值设为 `nil` —被推迟的函数能够修改已命名的返回值。在 `err` 的赋值过程中， 我们将通过断言它是否拥有局部类型 `Error` 来检查它。若它没有， 类型断言将会失败，此时会产生运行时错误，并继续栈的回溯，仿佛一切从未中断过一样。 该检查意味着若发生了一些像索引越界之类的意外，那么即便我们使用了 `panic` 和 `recover` 来处理解析错误，代码仍然会失败。
- 通过适当的错误处理，`error` 方法（由于它是个绑定到具体类型的方法， 因此即便它与内建的 `error` 类型名字相同也没有关系） 能让报告解析错误变得更容易，而无需手动处理回溯的解析栈：
    ```go
    if pos == 0 {
    	re.error("'*' illegal at start of expression")
    }
    ```
- 尽管这种模式很有用，但它应当仅在包内使用。`Parse` 会将其内部的 `panic` 调用转为 `error` 值，它并不会向调用者暴露出 `panic`。这是个值得遵守的良好规则。
- 顺便一提，这种重新触发Panic的惯用法会在产生实际错误时改变Panic的值。 然而，不管是原始的还是新的错误都会在崩溃报告中显示，因此问题的根源仍然是可见的。 这种简单的重新触发Panic的模型已经够用了，毕竟他只是一次崩溃。 但若你只想显示原始的值，也可以多写一点代码来过滤掉不需要的问题，然后用原始值再次触发Panic。 这里就将这个练习留给读者了。


# defer

- **defer延迟函数调用**/dɪ'fɝd/ ：观察代码中大量重复的resp.Body.Close()调用（Go的GC回收的是不被使用的内存，不包括操作系统层面的资源，如关闭文件、关闭网络连接），以保证title函数在任何执行路径下都会关闭网络连接。**Go独有的defer机制可以让代码变得简单****, 只需要在调用普通函数或方法前加上关键字defer；**
    - **defer能确保你在不断添加新的return时，不会忘记关闭文件、关闭网络连接；**
        ```go
        // 关闭网络连接
        func title(url string) error {
        	resp, err := http.Get(url)
        	if err != nil {
        		return err
        	}
          defer resp.Body.Close()
          
        	ct := resp.Header.Get("Content-Type")
        	if ct != "text/html" && !strings.HasPrefix(ct, "text/html;") {
        		return fmt.Errorf("%s has type %s, not text/html", url, ct)  		// resp.Body.Close()
        	}
        	doc, err := html.Parse(resp.Body)  	// resp.Body.Close()
        	if err != nil {
        		return fmt.Errorf("parsing %s as HTML: %v", url, err)
        	}
        	...
        	return nil
        }
        ```
        ```go
        // 关闭文件
        func Contents(filename string) (string, error) {
        	f, err := os.Open(filename)
        	if err != nil {
        		return "", err
        	}
        	defer f.Close()  // f.Close 会在我们结束后运行。
        	var result []byte
        	buf := make([]byte, 100)
        	for {
        		n, err := f.Read(buf[0:])
        		result = append(result, buf[0:n]...) // append 将在后面讨论。
        		if err != nil {
        			if err == io.EOF {
        				break
        			}
        			return "", err  // 我们在这里返回后，f 就会被关闭。
        		}
        	}
        	return string(result), nil // 我们在这里返回后，f 就会被关闭。
        }
        ```
    - **defer能让你的打开和关闭代码能成对的写在一块而使得代码清晰明了；**如打开/关闭、连接/断开连接、加锁/释放锁
        ```go
        // 处理互斥锁
        var mu sync.Mutex
        var m = make(map[string]int)
        func lookup(key string) int {
            mu.Lock()
            defer mu.Unlock()
            return m[key]
        }
        ```
    - **defer函数的实参表达式（如果该函数为方法则还包括接收者），在defer语句执行时就会求值**；（无需担心变量等到defer延迟执行时被改变）
    - **defer语句所在的函数的****所有函数出口执行完毕时****（包括return正常结束、由于panic导致的异常结束），所有defer后的函数才会被倒序执行（**按照后进先出（LIFO）的顺序执行**）；**
- 调试复杂程序时，defer机制也常被用于记录何时进入和退出函数：
    ```go
    func bigSlowOperation() {
    	// 函数值会在bigSlowOperation退出时被调用
    	// 只通过一条语句就能控制函数的入口和所有的出口，甚至可以记录函数的运行时间
    	// 注意（很微妙）：defer语句后有圆括号，否则不会执行函数值（本该在退出时执行的，永远不会被执行），而本该在进入时执行的操作会在退出时执行 // 2025/01/02 19:34:25 enter bigSlowOperation
    	defer trace("bigSlowOperation")() // don't forget the extra parentheses
    	// ...lots of work...
    	time.Sleep(3 * time.Second) // simulate slow operation by sleeping
    }
    func trace(msg string) func() {
    	start := time.Now()
    	log.Printf("enter %s", msg) // 2025/01/02 19:33:45 enter bigSlowOperation
    	// 返回一个函数值
    	return func() { log.Printf("exit %s (%s)", msg, time.Since(start)) } // 2025/01/02 19:33:55 exit bigSlowOperation (10.00141025s)
    }
    func triple(x int) (result int) {
    	// 1. 调用匿名函数/函数值（有圆括号），先计算表达式 result = 4
    	defer func() { result += x }()
    	// 2. 先计算表达式（调用double(x)）
    	// 5. 执行triple函数的defer：result = 8 + 4
    	return double(x)
    }
    func double(x int) (result int) {
    	// 3. 调用匿名函数/函数值（有圆括号），先计算表达式 x = 4
    	defer func() { fmt.Printf("double(%d) = %d\n", x, result) }()
    	// 4. 先计算表达式 result = 8
    	// 5. 执行double函数的defer：double(4) = 8
    	return x + x
    }
    //!-main
    func main() {
    	bigSlowOperation()
    	_ = double(4)          // double(4) = 8
    	fmt.Println(triple(4)) // double(4) = 8  12
    }
    /*
    !+output
    $ go build gopl.io/ch5/trace
    $ ./trace
    2015/11/18 09:53:26 enter bigSlowOperation
    2015/11/18 09:53:36 exit bigSlowOperation (10.000589217s)
    !-output
    *
    ```
- 代码bug: 会导致系统的文件描述符耗尽，因为在所有文件都被处理之前，没有文件会被关闭
    ```go
    // defer延迟的载体为所在的函数，只有在函数执行完毕后，这些被defer延迟的函数才会执行
    for _, filename := range filenames {
        f, err := os.Open(filename)
        if err != nil {
            return err
        }
        defer f.Close() // NOTE: risky; could run out of file descriptors
        // ...process f…
    }
    // 修复：将循环体中的defer语句单独包一个函数。在每次循环时，调用这个函数。
    for _, filename := range filenames {
        if err := doFile(filename); err != nil {
            return err
        }
    }
    func doFile(filename string) error {
        f, err := os.Open(filename)
        if err != nil {
            return err
        }
        defer f.Close()
        // ...process f…
    }
    ```
- fetch：
```go

func fetch(url string) (filename string, n int64, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	local := path.Base(resp.Request.URL.Path)
	if local == "/" {
		local = "index.html"
	}
	// 通过os.Create打开文件进行写入
	f, err := os.Create(local)
	if err != nil {
		return "", 0, err
	}
	n, err = io.Copy(f, resp.Body)
	// Close file, but prefer error from Copy, if any.
	// 在关闭文件时，我们没有对f.close采用defer机制，因为这会产生一些微妙的错误。原因：
	// 许多文件系统，尤其是NFS，写入文件时发生的错误会被延迟到文件关闭时反馈。所以如果没有检查文件关闭时的反馈信息（是否写入时发生错误），可能会导致数据丢失，而我们还误以为写入操作成功
	// 优先返回io.Copy的error，其次返回f.close的closeErr给调用者。因为它先于f.close发生，更有可能接近问题的本质
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	return local, n, err
}

func main() {
	for _, url := range os.Args[1:] {
		local, n, err := fetch(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch %s: %v\n", url, err)
			continue
		}
		fmt.Fprintf(os.Stderr, "%s => %s (%d bytes).\n", url, local, n)
	}
}
```



