---
title: "1.11 go package、mod"
date: 2024-12-28T15:09:00Z
draft: false
weight: 1011
---

# 1.11 go package、mod

- 今天一个中等规模的程序可能包含10000个函数，但是作者可能只须思考它们其中的10%，甚至不需要设计函数，因为绝大部分都是其他人来写的，然后通过包来复用。
- Go标准库目前已有304个包，为大多数的程序提供了必要的基础构件。Go社区有很多成熟的包被设计、共享、重用和改进，可以通过 [http://godoc.org](http://godoc.org/) 检索。
    ```shell
    // 查看标准包的具体数目
    $ go list std | wc -l
    304
    ```
- 任何包管理系统的目的都是通过对关联的特性进行分类，组织成便于理解和修改的单元，使其与程序的其他包保持独立，从而有助于设计和维护大型的程序。模块化允许包在不同的项目中共享、复用，在组织中发布，或者在全世界范围内使用。**支持模块化、封装、单独编译和代码重用**。
- 每个包定义了一个不同的命名空间作为它的标识符。每个名字关联一个具体的包，它让我们在为类型、函数等选取短小而且清晰的名字的同时，不与程序的其他部分冲突。
- 包通过控制名字是否导出使其对包外可见来提供封装能力。限制包成员的可见性，从而隐藏API后面的辅助函数和类型，允许包的维护者修改包的实现而不影响包外部的代码。限制变量的可见性也可以隐藏变量，这样使用者仅可以通过导出函数来对其访问和更新，他们可以保留自己的不变量以及在并发程序中实现互斥的访问。
    - 名字的外部可见性控制，一个简单的规则是：**如果一个名字是大写字母开头的，那么该名字是导出的**（译注：因为汉字不区分大小写，因此汉字开头的名字是没有导出的）。
- Go程序的编译比其他语言要快，即便从零开始编译也如此。主要原因有三：
    - 原因1：所有的导入都必须在每一个源文件的**开头**进行显式声明，这样编译器在确定依赖性的时候就不需要读取和处理整个文件；
    - 原因2：包的依赖性形成有向无环图，因为没有环，所以包可以独立甚至并行编译；**每个导入声明从当前包向导入的包建立一个依赖。如果这些依赖形成一个循环，go build工具会报错**。
    - 原因3：Go包编译输出的目标文件不仅记录它自己的导出信息，还记录它所依赖包的导出信息。当编译一个包时，编译器必须从每一个导入中读取一个目标文件，但是不会超出这些文件（译注：很多都是重复的间接依赖）。
- 导入路径：每一个包都通过一个唯一的字符串进行标识（称为导入路径、包名），用在import声明中；**为了避免冲突，除了标准库中的包之外，其他包的导入路径应该以互联网域名作为路径开始**，这样也方便查找包。
    - **导入的包可以通过空行进行分组，通常表示不同领域和方面的包**。
    - 导入顺序不重要，但按照惯例每一组都按照字母进行排序。（gofmt和goimports工具都会自动进行分组并排序。）
        - Go在坚持其强硬的导入的包必须有使用和代码格式化规则的设计哲学下，通过提供**goimports****工具****和****gofmt工具**，使得程序员在编辑器保存时，自动添加或删除导入的包、自动格式化Go源文件，实现无缝体验。（类似透明加解密）
    - 导入重命名：
        ```sql
        import (
            "crypto/rand"
        		 mrand "math/rand"
            "golang.org/x/net/html"
            "github.com/go-sql-driver/mysql"
        )
        ```
        - 如果需要把两个名字一样的包（如math/rand和crypto/rand）导入到第三个包中，导入声明就必须至少为其中的一个指定一个替代名字来避免冲突；
        - 如果有时用到自动生成的代码，导入的包名字非常冗长，使用一个替代名字可能更方便。同样的缩写名字要一直用下去，以避免产生混淆；
        - 空导入：当需要对包级别的变量执行初始化表达式求值，并执行它的init函数，但又未使用包时，会有“unused import”编译错误。Go的`_`空白标识符，并不能被访问(区别于python)；
            ```go
            package png /
            import  _ "image/png" // register PNG decoder 
            // 最终的效果是，主程序只需要匿名导入特定图像驱动包就可以用image.Decode解码对应格式的图像了。
            func Decode(r io.Reader) (image.Image, error)
            func DecodeConfig(r io.Reader) (image.Config, error)
            func init() {
                const pngHeader = "\x89PNG\r\n\x1a\n"
                image.RegisterFormat("png", pngHeader, Decode, DecodeConfig)
            }
            ```
- **包的import声明**：在每一个Go源文件的开头都需要进行包声明。通常包名是导入路径的最后一段，但有3个例外：
    - 例外1：如果包定义一个可执行程序，总是使用包名字main，这是告诉go build的信号，它必须调用连接器生成可执行文件；`package main`
    - 例外2：包所在的目录中可能有一些文件名字以_test.go结尾，包名中会出现以_test结尾。这样一个目录中有两个包：一个普通的，加上一个外部测试包。**_test后缀的外部拓展包由go test独立编译，并且指明文件属于哪个包**。**外部测试包一般用来避免测试代码中的循环导入依赖；**
    - 例外3：一些依赖版本号的管理工具会在导入路径后追加版本号信息，例如“gopkg.in/yaml.v2”。这种情况下**包的名字并不包含版本号后缀**，而是yaml。
        ```go
        // 数据库包database/sql也是采用了类似的技术，让用户可以根据自己需要选择导入必要的数据库驱动。例如：
        import (
            "database/sql"
            _ "github.com/lib/pq"              // enable support for Postgres
            _ "github.com/go-sql-driver/mysql" // enable support for MySQL
        )
        db, err = sql.Open("postgres", dbname) // OK
        db, err = sql.Open("mysql", dbname)    // OK
        db, err = sql.Open("sqlite3", dbname)  // returns error: unknown driver "sqlite3"
        ```
- **包的命名约定：**
    - **当创建一个包，一般要用短小的包名，但也不能短到像加了密一样**。标准库中最常用的包有bufio、bytes、flag、fmt、http、io、json、os、sort、sync和time等包；
    - 尽可能保持命名的**可读性和无歧义**。如不要把辅助工具包命名为util，imageutil、ioutilis等名称更具体和清晰。
    - 避免选择经常用于相关的局部变量的包名，或者迫使使用者使用重命名导入，如使用以path命名的包。
    - 包名通常使用统一的形式。标准包bytes、errors和strings使用复数来避免覆盖响应的预声明类型，使用go/types这个形式，来避免和关键字的冲突。
    - 避免使用有其他含义的包名，如temperature用tempconv，能和strconv等类似。
    - 包成员命名，需要同时考虑包名和成员名两个部分如何很好地组合命名。下面有一些例子：bytes.Equal、flag.Int、http.Get、json.Marshal
    - 还有一个以New命名的函数用于创建实例。
        ```go
        package rand // "math/rand"
        // 这可能导致一些名字重复，如template.Template或rand.Rand，这就是为什么这些种类的包名往往特别短的原因之一。
        type Rand struct{ /* ... */ }
        func New(source Source) *Rand
        ```
    - 在另一个极端，像net/http包有将近二十种类型和更多的函数，但包中最重要的成员名字却依然是保持简单明了的：Get、Post、Handle、Error、Client、Server等。
- Go的命令行接口使用“瑞士军刀”风格，带有十几个子命令，如get、run、build和fmt。可以运行go help来查看内置文档的索引
    ```go
    $ go
    ...
        build            compile packages and dependencies
        clean            remove object files
        doc              show documentation for package or symbol
        env              print Go environment information
        fmt              run gofmt on package sources
        get              download and install packages and dependencies
        install          compile and install packages and dependencies
        list             list packages
        run              compile and run Go program
        test             test packages
        version          print Go version
        vet              run go tool vet on packages
    Use "go help [command]" for more information about a command.
    ...
    ```
- 为了让配置操作最小化，go工具非常依赖惯例。如: 
    - 给定一个Go源文件的名称，Go语言的工具可以找到源文件对应的包，**因为每个目录只包含了单一的包**，并且包的导入路径对应于工作区的目录结构。
    - 给定一个包的导入路径，Go语言的工具可以找到存放目标文件的对应目录。也可以根据导入路径找到存储代码的仓库的远程服务器URL。
- **工作区结构：****大多数的Go语言用户只需要进行唯一的配置是GOPATH**，它指定工作空间的根。当需要切换到不同的工作空间时，更新GOPATH变量的值即可。如在编写本书时将GOPATH设置为`$HOME/gobook`： GOPATH有三个子目录。
    - src子目录包含源文件。每一个包放在一个目录中，该目录相对于$GOPATH/src的名字是包的导入路径，如 [gopl.io/ch1/helloworld](http://gopl.io/ch1/helloworld)；
    - pkg子目录是构建工具存储编译后的包的位置；
    - bin子目录放置像helloworld这样的可执行程序；
    ```go
    $ export GOPATH=$HOME/gobook
    $ go get gopl.io/...
    ```
- `go env`命令用于查看Go语言工具涉及的所有环境变量的值，包括未设置环境变量的默认值。GOOS环境变量用于指定目标操作系统（例如android、linux、darwin或windows），GOARCH环境变量用于指定处理器的类型，例如amd64、386或arm等。虽然GOPATH环境变量是唯一必须要设置的，但是其它环境变量也会偶尔用到。
    ```go
    $ go env
    GOPATH="/home/gopher/gobook"
    GOROOT="/usr/local/go"
    GOARCH="amd64"
    GOOS="darwin"
    ...
    ```
- **包的下载**：go get命令可以下载单一的包，也可以使用...符号来下载子树或仓库。在go get完成包的下载之后，它会构建它们，然后安装库和相应的命令。
- **包的构建：**`go build`命令编译命令行参数指定的每个包。如果包是一个库，结果会被舍弃；这可以用于检测包是可以正确编译的。如果包的名字是main，`go build`将调用链接器在当前目录创建一个可执行程序，可执行程序的名字取导入路径的最后一段。
    - **由于每个目录只包含一个包，因此每个可执行程序或者叫Unix命令都需要放到一个独立的目录中**。如 [golang.org/x/tools/cmd/godoc](http://golang.org/x/tools/cmd/godoc) 命令
    - 如果包名是main，可执行程序的名字来自第一个.go文件名的主体部分；
    - 对于即用即抛型的程序，我们需要在构建之后尽快运行。go run命令将这两步合并起来；
    - 默认情况下，go build命令构建所有需要的包以及它们所有的依赖性，然后丢弃除了最终可执行程序之外的所有编译后的代码。依赖性分析和编译本身都非常快，**但当项目增长到数十个包和数十万行代码的时候，重新编译依赖性的时间明显变慢，也许数秒钟的时间，即使依赖的部分根本没有改变过**；
    - `go install`命令和`go build`命令很相似，但是它会保存每个包的编译成果，而不是将它们都丢弃。被编译的包会被保存到$GOPATH/pkg目录下，目录路径和 src目录路径对应，可执行程序被保存到$GOPATH/bin目录。（很多用户会将$GOPATH/bin添加到可执行程序的搜索列表中。）还有，`go install`命令和`go build`命令都不会重新编译没有发生变化的包，这可以使后续构建更快捷
    - 为了方便编译依赖的包，`go build -i`命令将安装每个目标所依赖的包。
    - 针对不同操作系统或CPU的交叉构建也是很简单的。只需要设置好目标对应的GOOS和GOARCH，然后运行构建命令即可。下面交叉编译的程序将输出它在编译时的操作系统和CPU类型：
        ```go
        func main() {
            fmt.Println(runtime.GOOS, runtime.GOARCH)
        }
        // 下面以64位和32位环境分别编译和执行：
        $ go build gopl.io/ch10/cross
        $ ./cross
        darwin amd64
        $ GOARCH=386 go build gopl.io/ch10/cross
        $ ./cross
        darwin 386
        ```
    - 有些包可能需要针对不同平台和处理器类型使用不同版本的代码文件，以便于处理底层的可移植性问题或为一些特定代码提供优化。如果一个文件名包含了一个操作系统或处理器类型名字，例如net_linux.go或asm_amd64.s，Go语言的构建工具将只在对应的平台编译这些文件。还有一个特别的构建注释参数可以提供更多的构建过程控制。例如，文件中可能包含下面的注释：
        ```shell
        // +build linux darwin
        ```
    - 在包声明和包注释的前面，该构建注释参数告诉`go build`只在编译程序对应的目标操作系统是Linux或Mac OS X时才编译这个文件。下面的构建注释则表示不编译这个文件：
        ```shell
        // +build ignore
        ```
- **包的文档化：Go风格强烈鼓励有良好的包API文档。每一个导出的包成员的声明以及包声明自身应该立刻使用注释来描述它的目的和用途。**
    - **Go文档注释应当保持简洁，文档需要像代码一样维护**。使用声明的包名作为开头的第一句注释通常是总结。函数参数和其他的标识符并不需要用引号或括号特别标注。
        ```go
        // **Fprintf** formats according to a format specifier and writes to w.
        // It returns the number of bytes written and any write error encountered.
        // 第一行通常是摘要说明，以被注释者的名字开头。注释中函数的参数或其它的标识符并不需要额外的引号或其它标记注明。
        func Fprintf(w io.Writer, format string, a ...interface{}) (int, error)
        // Fprintf函数格式化的细节在fmt包文档中描述。如果注释后紧跟着包声明语句，那注释对应整个包的文档。包文档对应的注释只能有一个（译注：其实可以有多个，它们会组合成一个包文档注释），包注释可以出现在任何一个源文件中。
        // 如果包的注释内容比较长，一般会放到一个独立的源文件中；fmt包注释就有300行之多。这个专门用于保存包文档的源文件通常叫doc.go。
        ```
    - 比较长的包注释可以使用一个单独的注释文件，fmt的注释超过300行，文件名通常叫doc.go；
        ```go
        $ go doc time  // 打印其后所指定的实体的声明与文档注释，如一个包
        $ go doc time.Since  // 某个具体的包成员
        $ go doc time.Duration.Seconds  // 一个方法
        $ go doc json.decode  // 该命令并不需要输入完整的包导入路径或正确的大小写
        ```
    - godoc的在线服务 [https://godoc.org](https://godoc.org/) ，包含了成千上万的开源包的检索工具。
        ```go
        $ godoc -http :8000 // 其中-analysis=type和-analysis=pointer命令行标志参数用于打开文档和代码中关于静态分析的结果。
        // 在浏览器查看 http://localhost:8000/pkg 页面
        ```
- **内部包：**Go中的包是最重要的封装机制。**没有导出的标识符只在同一个包内部可以访问，而导出的标识符则是面向全世界任何地方都是可见的。**
    - 但一个中间地带是很有帮助的，这种方式定义标识符可以被一个小的可信任的包集合访问，但不是所有人可以访问；如：当我们计划将一个大的包拆分为很多小的更容易维护的子包，但是我们并不想将内部的子包结构也完全暴露出去； 
    - **internal包：Go语言的构建工具对包含internal名字的路径段的包导入路径做了特殊处理，一个internal包只能被和internal目录有同一个父目录的包所导入**。如，net/http/internal/chunked内部包只能被net/http/httputil或net/http包导入，但是不能被net/url包导入。不过net/url包却可以导入net/http/httputil包。
        ```go
        net/http
        net/http/internal/chunked
        net/http/httputil
        net/url
        ```
- **包的查询：go list工具上报可用包的信息**
    ```go
    $ go list github.com/go-sql-driver/mysql
    github.com/go-sql-driver/mysql
    $ go list ...  // 列出工作区中的所有包
    $ go list gopl.io/ch3/...  // 或者是特定子目录下的所有包
    $ go list ...xml...  // 或者是和某个主题相关的所有包
    $ go list -json hash  // 获取每个包完整的元信息
    ```
- **包的初始化：****首先是解决包级变量的依赖顺序**，**然后****按照包级变量声明出现的顺序依次初始化**：
    ```go
    var a = b + c // a 第三个初始化, 为 3
    var b = f()   // b 第二个初始化, 为 2, 通过调用 f (依赖c)
    var c = 1     // c 第一个初始化, 为 1
    func f() int { return c + 1 }
    ```
    如果包中含有多个.go源文件，它们将按照发给编译器的顺序进行初始化，Go语言的构建工具**首先会将.go文件根据文件名排序****，然后****依次调用编译器编译**。
    对于在包级别声明的变量，如果**有初始化表达式则用表达式初始化**，还有一些没有初始化表达式的（如某些表格数据初始化并不是一个简单的赋值过程）**可以用一个特殊的init初始化函数来简化初始化工作**。每个文件都可以包含多个init初始化函数。这样的init初始化函数**除了不能被调用或引用外**，其他行为和普通函数类似。在每个**文件中的init**初始化函数，**在程序开始执行时****按照它们声明的顺序****被自动调用**。
    ```go
    func init() { /* ... */ }
    ```
    - 每个包在解决依赖的前提下，以导入声明的顺序初始化，每个包只会被初始化一次。因此，如果一个p包导入了q包，那么在p包初始化的时候可以认为q包必然已经初始化过了。
    - 初始化工作是**自下而上进行的**，**main包最后被初始化**。可以确保在main函数执行之前，所有依赖的包都已经完成初始化工作了。


