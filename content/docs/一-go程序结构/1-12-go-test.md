---
title: "1.12 go test"
date: 2024-12-28T15:09:00Z
draft: false
weight: 1012
---

# 1.12 go test



莫里斯·威尔克斯(Maurice Wilkes)设计和制造了世界上第一台存储程序式计算机EDSAC，在1949年有一次实验室爬楼梯时有一个顿悟：我强烈地意识到我余生的很大一部分时间都将用来寻找我程序中的错误。”今天的软件项目比威尔克斯年代的要庞大、复杂得多，并且在使软件复杂度可以控制的技术上面，人们投入了大量的精力。其中有两种技术尤其有效：

- 第一种是代码在被正式部署前需要进行**代码评审**。
- 另一个就是：**测试。**（费曼学习法：最小可行性验证、字节的产品ABTest）。这里一般是指自动化测试，**即编写简单的程序来确保程序（产品代码）在该测试中针对特定输入产生预期的输出。**这些测试通常要么是经过精心设计之后用来检测某种功能，要么是随机性的，用来扩大测试的覆盖面。
- Go的测试方法看上去相对比较低级，它依赖于命令go test和一些能用go test运行的测试函数的编写约定。这个相对轻量级的机制对单纯的测试很有效，并且这种方式也很自然地扩展到基准测试和文档系统的示例。这和编写常规的Go代码没有区别。
### go test

- **在一个包目录中，以_test.go结尾的文件不是go build命令编译的目标，而是go test编译的目标。**
- 在`*_test.go`文件中，有三种类型的函数需要特殊对待：测试函数、基准测试（benchmark）函数、示例函数：
    - **功能测试函数：以Test前缀命名的函数，用来检测一些程序逻辑的正确性，go test运行测试函数，并且报告结果是PASS还是FAIL；**
    - **基准测试函数：以Benchmark前缀命名的函数，用于测试一些函数的性能；go test命令会多次运行基准测试函数以计算出一个平均的执行时间。**
    - **示例函数：以Example前缀命名的函数，用来提供提供一个由编译器保证正确性的示例文档。**
- go test工具扫描*_test.go文件来寻找这三种特殊函数，**并生成一个临时的main包来调用它们**，然后编译/构建（go build）和运行(go run)，并汇报结果，最后清空临时文件。
### **功能测试函数  Test***

- 必须以Test开头，且可选的后缀名称也必须以大写字母开头:
    ```go
    func TestSin(t *testing.T) { /* ... */ }
    func TestCos(t *testing.T) { /* ... */ }
    func TestLog(t *testing.T) { /* ... */ }
    ```
- 每个测试函数必须导入testing包。测试函数有如下的签名，**参数t提供了汇报测试失败和日志记录的功能**
    ```go
    import testing
    func TestName(t *testing.T) {
        // ...
    }
    ```


**示例：**

1. 编写函数
    ```go
    // Package word provides utilities for word games.
    package word
    // IsPalindrome reports whether s reads the same forward and backward.
    // (Our first attempt.)
    func IsPalindrome(s string) bool {
        for i := range s {
            if s[i] != s[len(s)-1-i] {
                return false
            }
        }
        return true
    }
    ```
1. 编写测试函数
    ```go
    package word
    import "testing"
    func TestPalindrome(t *testing.T) {
        if !IsPalindrome("detartrated") {
            t.Error(`IsPalindrome("detartrated") = false`)
        }
        if !IsPalindrome("kayak") {
            t.Error(`IsPalindrome("kayak") = false`)
        }
    }
    func TestNonPalindrome(t *testing.T) {
        if IsPalindrome("palindrome") {
            t.Error(`IsPalindrome("palindrome") = true`)
        }
    }
    ```


1. 发现线上Bug：一个法国名为“Noelle Eve Elleon”的用户会抱怨IsPalindrome函数不能识别“été”。另外一个来自美国中部用户的抱怨则是不能识别“A man, a plan, a canal: Panama.”；
1. 定位原因：
    - **先写go test测试用例（运行go test比手动测试bug报告中的内容要快得多）**，然后确保它触发的错误和用户bug报告里面的一致，以定位到bug原因；
        ```go
        func TestFrenchPalindrome(t *testing.T) {
            if !IsPalindrome("été") {
                t.Error(`IsPalindrome("été") = false`)  // 原因：非ASCII字符byte类型无法正确处理，用rune类型
            }
        }
        func TestCanalPalindrome(t *testing.T) {
            input := "A man, a plan, a canal: Panama"   
            if !IsPalindrome(input) {      // 原因：没有忽略空格和字母的大小写
                t.Errorf(`IsPalindrome(%q) = false`, input)
            }
        }
        ```
    - 参数`-v`可用于打印每个测试函数的名字和运行时间：
        ```go
        $ go test -v
        === RUN TestPalindrome
        --- PASS: TestPalindrome (0.00s)
        === RUN TestNonPalindrome
        --- PASS: TestNonPalindrome (0.00s)
        === RUN TestFrenchPalindrome
        --- FAIL: TestFrenchPalindrome (0.00s)
            word_test.go:28: IsPalindrome("été") = false
        === RUN TestCanalPalindrome
        --- FAIL: TestCanalPalindrome (0.00s)
            word_test.go:35: IsPalindrome("A man, a plan, a canal: Panama") = false
        FAIL
        exit status 1
        FAIL    gopl.io/ch11/word1  0.017s
        ```
    - 参数`-run`对应一个正则表达式，只有测试函数名被它正确匹配的测试函数才会被`go test`测试命令运行：
        ```go
        $ go test -v -run="French|Canal"
        === RUN TestFrenchPalindrome
        --- FAIL: TestFrenchPalindrome (0.00s)
            word_test.go:28: IsPalindrome("été") = false
        === RUN TestCanalPalindrome
        --- FAIL: TestCanalPalindrome (0.00s)
            word_test.go:35: IsPalindrome("A man, a plan, a canal: Panama") = false
        FAIL
        exit status 1
        FAIL    gopl.io/ch11/word1  0.014s
        ```
    - `go test`命令如果没有参数指定包那么将默认采用当前目录对应的包（和`go build`命令一样）。我们可以用下面的命令构建和运行测试。
1. 修复bug：
    ```go
    // Package word provides utilities for word games.
    package word
    import "unicode"
    // IsPalindrome reports whether s reads the same forward and backward.
    // Letter case is ignored, as are non-letters.
    func IsPalindrome(s string) bool {
        var letters []rune
        for _, r := range s {
            if unicode.IsLetter(r) {
                letters = append(letters, unicode.ToLower(r))
            }
        }
        for i := range letters {
            if letters[i] != letters[len(letters)-1-i] {
                return false
            }
        }
        return true
    }
    ```
1. 回归测试：在提交代码更新之前，**使用不带参数的go test命令以运行全部的测试用例(回归测试)**，**以确保修复失败测试的同时没有引入新的bug**；


- 和其他编程语言或测试框架的assert断言不同，t.Errorf调用也没有引起panic异常或停止测试的执行。即使表格中前面的数据导致了测试的失败，表格后面的测试数据依然会运行测试，因此在一个测试中我们可能了解多个失败的信息。如果我们真的需要停止测试，或许是因为初始化失败或可能是早先的错误导致了后续错误等原因，我们可以使用t.Fatal或t.Fatalf停止当前测试函数。它们必须在和测试函数同一个goroutine内调用。
- **测试失败的信息一般的形式：“f(x) = y, want z”**。其中f(x)解释了失败的操作和对应的输入，y是实际的运行结果，z是期望的正确的结果。要避免无用和冗余的信息。测试的作者应该要努力帮助程序员诊断测试失败的原因。
- 测试样例：
    - **基于测试用例表的测试方式**：将之前的所有测试用例合并到了一个测试表格中很直观
        ```go
        func TestIsPalindrome(t *testing.T) {
            var tests = []struct {
                input string
                want  bool
            }{
                {"", true},
                {"a", true},
                {"aa", true},
                {"ab", false},
                {"kayak", true},
                {"detartrated", true},
                {"A man, a plan, a canal: Panama", true},
                {"Evil I did dwell; lewd did I live.", true},
                {"Able was I ere I saw Elba", true},
                {"été", true},
                {"Et se resservir, ivresse reste.", true},
                {"palindrome", false}, // non-palindrome
                {"desserts", false},   // semi-palindrome
            }
            for _, test := range tests {
                if got := IsPalindrome(test.input); got != test.want {
                    t.Errorf("IsPalindrome(%q) = %v", test.input, got)
                }
            }
        }
        ```
    - 随机测试：通过构建随机输入来扩展测试的覆盖范围。
        - 由于随机测试的不确定性，在遇到测试用例失败的情况下，一定要记录足够的信息以便于复现问题，如记录伪随机数生成器的种子
- 在测试的代码里面不要调用log.Fatal或者os.Exit，因为这两个调用会阻止跟踪的过程，这两个函数的调用可以认为是main函数的特权
- 黑盒测试假设测试者对包的了解仅通过公开的API和文档，而包的内部逻辑则是不透明的。相反，白盒测试可以访问包的内部函数和数据结构，并且可以做一些常规用户无法做到的观察和改动。黑盒测试通常更加健壮，每次程序更新后基本不需要修改。白盒测试可以对实现的特定之处提供更详细的覆盖测试。
    - TestIsPalindrome函数仅调用导出的函数IsPalindrome，所以它是一个黑盒测试。
    - TestEcho函数调用echo函数并且更新全局变量out，无论函数echo还是变量out都是未导出的，所以它是一个白盒测试。
    - 可以使用易于测试的伪实现来替换部分产品代码。这种Mock模拟的伪实现的优点是更易于配置、预测和观察，并且更可靠。它们还能够避免带来副作用，比如更新产品数据库或者刷信用卡。
- **外部测试包**：低级别包的测试导入了高级别包会导致包循环引用而产生编译错误，将这个测试函数定义在外部测试包中来解决这个问题。在net/url目录中，外部测试包的声明是url_test 独立的一个包；
    - 有时候，外部测试包需要对待测试包拥有特殊的访问权限，例如为了避免循环引用，一个白盒测试必须存在于一个单独的包中。在这种情况下，我们使用一种小技巧：在包内测试文件_test.go中添加一些函数声明，将包内部的功能暴露给外部测试。这些文件也因此为测试提供了包的一个“后门”。如果一个源文件存在的唯一目的就在于此，并且自己不包含任何测试，它们一般称作export_test.go。
- **编写有效测试**：
    - 其他语言的框架提供了识别测试函数的机制（一般通过反射或者元数据），在测试前后执行测试“启动”和“销毁”的钩子，以及为常规的断言、值比较、错误消息格式化和终止失败的测试（一般通过抛出异常的方式）提供工具方法的库。但导致的结果是这些测试看上去像是用一门其他的语言编写的。
    - Go对测试的看法是完全不同的。它期望测试的编写者自己来做大部分的工作，通过定义函数来避免重复，就像他们为普通程序所做的那样。测试的过程不是死记硬背地填表格；测试也是有用户界面的，虽然它的用户也是它的维护者。一个好的测试不会在发生错误时崩溃而是输出该问题一个简洁、清晰的现象描述，以及其他与上下文相关的信息。**理想情况下，维护者通过测试输出结果，而不需要再通过阅读源代码来探究测试失败的原因（同打Log）**。一个好的测试不应该在发现一次测试失败后就终止，而是要在一次运行中尝试报告多个错误，因为错误发生的方式本身会揭露错误的原因。
        ```go
        import (
            "fmt"
            "strings"
            "testing"
        )
        // A poor assertion function.
        func assertEqual(x, y int) {
            if x != y {
                panic(fmt.Sprintf("%d != %d", x, y))
            }
        }
        func TestSplit(t *testing.T) {
            words := strings.Split("a:b:c", ":")
            assertEqual(len(words), 3)     // 断言函数犯了过早抽象的错误：仅仅测试两个整数是否相同，而没能根据上下文提供更有意义的错误信息。
            // ...
        }
        func TestSplit(t *testing.T) {
            s, sep := "a:b:c", ":"
            words := strings.Split(s, sep)
            if got, want := len(words), 3; got != want {    
                t.Errorf("Split(%q, %q) returned %d words, want %d",         // **不仅报告了调用的具体函数、它的输入和结果的意义；并且打印的真实返回的值和期望返回的值；并且即使断言失败依然会继续尝试运行更多的测试**
                    s, sep, got, want)
            }
            // ...
        }
        ```
- 避免脆弱的测试：如果一个应用在遇到新的合法输入的情况下经常崩溃，那么这个程序是有缺陷的；如果在程序发生可靠的改动的时候测试用例奇怪地失败了，那么这个测试用例也是脆弱的。**最脆弱的测试在产品代码发生任何改动的时候都会失败，无论这些改动是好是坏**，这些测试通常称为**变化探测器(changedetector)或现状探测器(status quo test)**，并且处理它们花费的时间将会使得它们曾经带来的好处消失殆尽。
- 从本质上看，测试不可能是完整的。著名计算机科学家EdsgerDijkstra说，“测试能证明bug存在，而无法证明bug不存在。” 无论有多少测试都无法证明一个包是没有bug的。在最好的情况下，测试可以增强了我们的信心，这些包是可以在很多重要的场景下正常工作的。（类似渗透测试相对的证明系统的安全性）；
- **测试覆盖率**：对待测程序执行的测试的比例称为测试的覆盖率；
    - 语句的覆盖率是指在测试中至少被运行一次的代码占总代码数的比例。**在运行每个测试前，它将待测代码拷贝一份并做修改，在每个词法块都会设置一个布尔标志变量。当被修改后的被测试代码运行退出时，将统计日志数据写入c.out文件，并打印一部分执行的语句的一个总结。**
    - 如果使用了`-covermode=count`标志参数，那么将在每个代码块插入一个计数器而不是布尔标志量。在统计结果中记录了每个块的执行次数，这可以用于衡量哪些是被频繁执行的热点代码。
    - 
    ```go
    $ go test -run=Coverage -coverprofile=c.out gopl.io/ch7/eval
    ok      gopl.io/ch7/eval         0.032s      coverage: 68.5% of statements
    ```
    ```go
    $ go tool cover -html=c.out
    ```
    ![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/9e7f22fe-c876-4b16-abc4-63de75b68775/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4664PWIWZVC%2F20250718%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250718T165315Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEHgaCXVzLXdlc3QtMiJIMEYCIQClgNNh9nLRYQm4zE1FvoQtbnmJFuQM%2BbJrv4a2zEbCDQIhAN8Wju4n1CvpVnQ8W0Vnh%2FZ0ylAhdhOZzjZU3FCsLc1AKogECJH%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1Igw9rBdEyKh4LtX%2B3JYq3AOAUTgf9DH3ZXIY1j%2BWmZzFAGhSCeeGcytC6GNagBZJQS1Abk1ZCCrJPXv3W91oX0zmiNWBe9xygF9tUkyAW4OeVw0Pd0CXPkdB7h5Qyttpa29WlRdiWA%2FkJWEdk2Lpmb%2FuiPrfb1IdDuH583wdy2kugTzMnWh95bP76jBERNolucBxTFuuyPU25I2iNkDm3WzioVs7gpYn9D0Gezdnzt%2Bp4J2E6R%2F8Zb9VOV%2B0c0h9YHjmNs7FxK%2Bg%2Bvuvd19dfTHLxTpx9s3Ux5XQQBPqiXkv597ofr5zftGcxzgVg6uoTFHxPSsO779S7tm%2FMXTNe6n%2BvrTk3U%2B5yqQhZNLdgRifpSkV9L9F%2F0K3QB4kWd2Fo458lXjyPR1U3wjeYSuCFOHYbZuW%2Fd%2FZyr7e9bIGotcf5teXub3pU%2FJfeoLGb2EVI2jF61klr08yi9sh5jNovoyGodmPNh4vY3EHOXEkUW9YOdAQqlEIPYp%2BmLOE3H6M1RF6L5S7wLJXZLkYNIHLY2VPjXQLp9Eum%2FioMq%2BSsx4cGIsqscXBZaWcmW%2FzBxUDJyR4v%2BLyXTHu9EyX8ccdiMdmpVpWfsKJK1JU19d4GdRZkTqEEyYmaJOaaMjYNgesLc%2BEWBsDnmPkrJuTfTC1zunDBjqkAXDKlDmY%2FnCvdAz5gXPL5QeUvyTix0UxSG83u3DqTANPqFXdF5SDGRvba6%2FpI8yLrA9IwQPS1lS3JDb6IaxQg82UgZWSvmEPOy7uOpS3JepwxlwRQQmvGfncxwMbmuKkNrciiyO3pYP%2F3RIfFGahl%2BS9PcKzwmvGfJt7EJi3D7DDGzDga0u9F9WODqbAjdg06ipK2Li9PDSwUKKGJFZfumzqIcJ9&X-Amz-Signature=a9caf96c29b4a18176e9c0ca7efec901c7d85f43620a00221089943a9f2d2ff1&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)
    - 红色的则表示没有被覆盖到，于是在添加下面的测试用例，确保红色部分的未测试的代码也变成绿色已测试的代码了
    - 实现100%的测试覆盖率听起来很美，但是在具体实践中通常是不可行的，也不是值得推荐的做法。因为那只能说明代码被执行过而已，并不意味着代码就是没有BUG的；因为对于逻辑复杂的语句需要针对不同的输入执行多次。有一些语句，例如上面的panic语句则永远都不会被执行到。另外，还有一些隐晦的错误在现实中很少遇到也很难编写对应的测试代码。测试从本质上来说是一个比较务实的工作，**编写测试代码和编写应用代码的成本对比是需要考虑的**。测试覆盖率工具可以帮助我们快速识别测试薄弱的地方，但是设计好的测试用例和编写应用代码一样需要严密的思考。
### **基准测试函数 Benchmark***

- 基准测试是测量一个程序在固定工作负载下的性能。和普通测试函数写法类似，但以Benchmark为前缀名，并且带有一个`*testing.B`类型的参数；`*testing.B`参数除了提供和`*testing.T`类似的方法，还有额外一些和性能测量相关的方法。它还提供了一个整数N，用于指定操作执行的循环次数。
    ```go
    import "testing"
    func BenchmarkIsPalindrome(b *testing.B) {
        for i := 0; i < b.N; i++ {
            IsPalindrome("A man, a plan, a canal: Panama")
        }
    }
    ```
- 默认情况下不运行任何基准测试，**需要通过**`**-bench**`**命令行标志参数手工指定要运行的基准测试函数**。该参数是一个正则表达式，用于匹配要执行的基准测试函数的名字，默认值是空的。其中“.”模式将可以匹配所有基准测试函数，但因为这里只有一个基准测试函数，因此和`-bench=IsPalindrome`参数是等价的效果。
    - 结果中基准测试名的数字后缀部分，这里是8，表示运行时对应的GOMAXPROCS的值，这对于一些与并发相关的基准测试是重要的信息。
    - 报告显示每次调用IsPalindrome函数花费1.035微秒，是执行1,000,000次（一百万次）的平均时间；
    - 因为基准测试驱动器开始时并不知道每个基准测试函数运行所花的时间，它会尝试在真正运行基准测试前先尝试用较小的N运行测试来估算基准测试函数所需要的时间，然后推断一个较大的时间保证稳定的测量结果。
    - 循环在基准测试函数内实现，而不是放在基准测试框架内实现，这样可以让每个基准测试函数有机会在循环启动前执行初始化代码，这样并不会显著影响每次迭代的平均运行时间。如果还是担心初始化代码部分对测量时间带来干扰，那么可以通过testing.B参数提供的方法来临时关闭或重置计时器，不过这些一般很少会用到。
    ```go
    $ cd $GOPATH/src/gopl.io/ch11/word2
    $ go test -bench=.
    PASS
    BenchmarkIsPalindrome-8 1000000                1035 ns/op
    ok      gopl.io/ch11/word2      2.179s
    ```
- IsPalindrome函数性能优化：
    ```go
    n := len(letters)/2
    for i := 0; i < n; i++ {
        if letters[i] != letters[len(letters)-1-i] {   
            return false     // 避免每个比较都做两次
        }
    }
    return true
    $ go test -bench=.
    PASS
    BenchmarkIsPalindrome-8 1000000              992 ns/op   // 性能+4%
    ok      gopl.io/ch11/word2      2.093s  
    ```
- 再优化：在开始为每个字符预先分配一个足够大的数组，这样就可以避免在append调用时可能会导致内存的多次重新分配
    ```go
    letters := make([]rune, 0, len(s))
    for _, r := range s {
        if unicode.IsLetter(r) {
            letters = append(letters, unicode.ToLower(r))
        }
    }
    $ go test -bench=.
    PASS
    BenchmarkIsPalindrome-8 2000000                      697 ns/op   // 性能+35%
    ok      gopl.io/ch11/word2      1.468s
    ```
- 例子证明：快的程序往往是伴随着较少的内存分配。用一次内存分配代替多次的内存分配，节省了75%的分配调用次数和减少近一半的内存需求。
    ```go
    $ go test -bench=. -benchmem
    PASS
    BenchmarkIsPalindrome    1000000   1026 ns/op    304 B/op  4 allocs/op   // 优化后
    $ go test -bench=. -benchmem   
    PASS
    BenchmarkIsPalindrome    2000000    807 ns/op    128 B/op  1 allocs/op   // 优化后
    ```
- 比较型的基准测试就是普通程序代码。它们通常是单参数的函数，由几个不同数量级的基准测试函数调用，就像这样：
    - 通过函数参数来指定输入的大小，但是参数变量对于每个具体的基准测试都是固定的。要避免直接修改b.N来控制输入的大小。除非你将它作为一个固定大小的迭代计算输入，否则基准测试的结果将毫无意义。
    - 比较型的基准测试反映出的模式在程序设计阶段是很有帮助的，但是即使程序完工了也应当保留基准测试代码。因为随着项目的发展，或者是输入的增加，或者是部署到新的操作系统或不同的处理器，我们可以再次用基准测试来帮助我们改进设计。
    ```go
    func benchmark(b *testing.B, size int) { /* ... */ }
    func Benchmark10(b *testing.B)         { benchmark(b, 10) }
    func Benchmark100(b *testing.B)        { benchmark(b, 100) }
    func Benchmark1000(b *testing.B)       { benchmark(b, 1000) }
    ```
- 性能剖析：
    - 唐纳德·克努斯 的不要过早优化的箴言：
        - 毫无疑问，对效率的片面追求会导致各种滥用。**程序员会浪费大量的时间来思考或担心他们非关键部分的代码的执行速度上**，并且考虑到**程序的调试和维护的时候这些优化的尝试事实上会带来负面的影响**。特别是当调试和维护的时候。**我们必须舍弃微小的性能提升，在97%的情况下：过早的优化是万恶之源**。（费曼法：牺牲代码可维护性而过度优化的hive sql）
        - 然而我们不可以错过那关键的3％的情况。一个好的程序员不会因为这个就自满，明智的方法是：**仅仅在3%的关键代码明确之后，务必仔细地优化关键代码**；通常情况下先入为主地认定程序哪些部分是关键代码是错误的，使用了检测工具的程序员会发现的普遍经验就是他们的直觉是错的。（费曼法：马斯克五步工作法之首：程序员上来先不要优化，而是删减）
    - 基准测试（Benchmark）对于衡量特定操作的性能是有帮助的，但是当我们试图让程序跑的更快的时候，我们通常并不知道从哪里开始优化。
    - **性能剖析**：通过自动化手段在程序执行过程中基于一些性能事件的采样来进行性能评测，然后再从这些采样中推断分析，得到的统计报告就称作为性能剖析(profile)。
        - 每个CPU上面执行的线程都每隔几毫秒会定期地被操作系统中断，在每次中断过程中记录一个性能剖析事件，然后恢复正常执行。
        - 堆性能剖析识别出负责分配最多内存的语句。性能剖析库对协程内部内存分配调用进行采样，因此每个性能剖析事件平均记录了分配的512KB内存。
        - 阻塞性能剖析识别出那些阻塞goroutine最久的操作，例如系统调用，通道发送和接收数据，以及获取锁等。性能分析库在一个goroutine每次被上述操作之一阻塞的时候记录一个事件。
        - 获取待测试代码的性能剖析报告很容易，只需要像下面一样指定一个标记即可。当一次使用多个标记的时候需要注意，获取性能分析报告的机制是当获取其中一个类别的报告时会覆盖掉其他类别的报告。
        - 在我们获取性能剖析结果后，我们需要使用pprof工具来分析它。这是Go发布包的标准部分，但是因为不经常使用，所以通过go tool pprof间接来使用它。它有很多特性和选项，但是基本的用法只有两个参数，产生性能剖析结果的可执行文件和性能剖析日志。
        - 为了提高分析效率和减少空间，分析日志本身并不包含函数的名字；它只包含函数对应的地址。也就是说pprof需要对应的可执行程序来解读剖析数据。虽然`go test`通常在测试完成后就丢弃临时用的测试程序，但是在启用分析的时候会将测试程序保存为foo.test文件，其中foo部分对应待测包的名字。
        - 下面的命令演示了如何收集并展示一个CPU分析文件。我们选择`net/http`包的一个基准测试为例。通常最好是对业务关键代码的部分设计专门的基准测试。因为简单的基准测试几乎没法代表业务场景，因此我们用-run=NONE参数禁止那些简单测试。
            - 参数`-text`用于指定输出格式，在这里每行是一个函数，根据使用CPU的时间长短来排序。其中`-nodecount=10`参数限制了只输出前10行的结果。对于严重的性能问题，这个文本格式基本可以帮助查明原因了。
            - 这个概要文件告诉我们，HTTPS基准测试中`crypto/elliptic.p256ReduceDegree`函数占用了将近一半的CPU资源，对性能占很大比重。相比之下，如果一个概要文件中主要是runtime包的内存分配的函数，那么减少内存消耗可能是一个值得尝试的优化策略。
            - 对于一些更微妙的问题，你可能需要使用pprof的图形显示功能。这个需要安装GraphViz工具，可以从 [http://www.graphviz.org](http://www.graphviz.org/) 下载。参数`-web`用于生成函数的有向图，标注有CPU的使用和最热点的函数等信息。
            ```go
            $ go test -run=NONE -bench=ClientServerParallelTLS64 \
                -cpuprofile=cpu.log net/http
             PASS
             BenchmarkClientServerParallelTLS64-8  1000
                3141325 ns/op  143010 B/op  1747 allocs/op
            ok       net/http       3.395s
            $ go tool pprof -text -nodecount=10 ./http.test cpu.log
            2570ms of 3590ms total (71.59%)
            Dropped 129 nodes (cum <= 17.95ms)
            Showing top 10 nodes out of 166 (cum >= 60ms)
                flat  flat%   sum%     cum   cum%
              1730ms 48.19% 48.19%  1750ms 48.75%  crypto/elliptic.p256ReduceDegree
               230ms  6.41% 54.60%   250ms  6.96%  crypto/elliptic.p256Diff
               120ms  3.34% 57.94%   120ms  3.34%  math/big.addMulVVW
               110ms  3.06% 61.00%   110ms  3.06%  syscall.Syscall
                90ms  2.51% 63.51%  1130ms 31.48%  crypto/elliptic.p256Square
                70ms  1.95% 65.46%   120ms  3.34%  runtime.scanobject
                60ms  1.67% 67.13%   830ms 23.12%  crypto/elliptic.p256Mul
                60ms  1.67% 68.80%   190ms  5.29%  math/big.nat.montgomery
                50ms  1.39% 70.19%    50ms  1.39%  crypto/elliptic.p256ReduceCarry
                50ms  1.39% 71.59%    60ms  1.67%  crypto/elliptic.p256Sum
            ```
### **示例函数 Example***

- 示例函数没有函数参数和返回值。下面是IsPalindrome函数对应的示例函数：
    ```go
    func ExampleIsPalindrome() {
        fmt.Println(IsPalindrome("A man, a plan, a canal: Panama"))
        fmt.Println(IsPalindrome("palindrome"))
        // Output:
        // true
        // false
    }
    ```
- 示例函数有三个用处。
    - 最主要的一个是作为文档
    - 第二个用处是，在`go test`执行测试的时候也会运行示例函数测试。
    - 第三个目的提供一个真实的演练场。如： [http://golang.org](http://golang.org/)


