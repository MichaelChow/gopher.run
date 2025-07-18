---
title: "1.8 控制结构: if、switch、select、for"
date: 2025-03-31T23:18:00Z
draft: false
weight: 1008
---

# 1.8 控制结构: if、switch、select、for

### if

- `if` 、`switch` ，像 `for`一样**可接受可选的初始化语句**设置局部变量；
    ```go
    if err := file.Chmod(0664); err != nil {
    	log.Print(err)
    	return err
    }
    ```
- 语法上其主体强制大括号促使语句分成多行而更简洁清晰，没有圆括号；
    ```go
    if x > 0 {
    	return y
    }
    ```
- **Go采用排错式风格(如err、nil未处理，导致线上运行时故障)，若控制流成功继续，则说明程序已排除错误，正常直行的语句块不会被if else缩进的支离破碎。**由于出错时将以`return`、`break`、`continue`、`goto` 结束， 之后的代码也就无需`else`了。
    - 由于两个err在同一个词法域，且d为声明，Go在短变量声明语法上出于实用性设计，后一个err声明为重新赋值，而不必使用err1、err2、err3…
    ```go
    f, err := os.Open(name)
    if err != nil {
        return err
    }
    d, err := f.Stat()  
    if err != nil {
        f.Close()
        return err
    }
    codeUsing(f, d)
    ```
- Go只有一个更通用的 `for`循环，不再使用 `do` 或 `while` 循环；
- `select`包含类型选择和多路通信复用器的新控制结构；
- Go中的函数形参和返回值的作用域，和在词法上处于大括号内的函数体一致；


### **switch**

- Go风格将`if-elif-else` 链写成一个 `switch`，其表达式无需为常量或整数，`case` 语句会自上而下逐一进行求值直到匹配为止，最后一个没有case将匹配为true。
    ```go
    func unhex(c byte) byte {
    	switch {
    	case '0' <= c && c <= '9':
    		return c - '0'
    	case 'a' <= c && c <= 'f':
    		return c - 'a' + 10
    	case 'A' <= c && c <= 'F':
    		return c - 'A' + 10
    	}
    	return 0
    }
    ```
- `switch` 默认不会自动下溯（case匹配执行后自动跳出switch），`case`可通过逗号分隔来列举相同的处理条件。显式下溯使用fallthrough关键字；
    ```go
    func shouldEscape(c byte) bool {
    	switch c {
    	case ' ', '?', '&', '=', '#', '+', '%':
    		return true
    	}
    	return false
    }
    ```
- `break` 语句可以使 `switch` 打破层层的循环提前终止。在Go中，我们只需将标签放置到循环外，然后 “蹦”到那里即可。`continue` 语句也能接受一个可选的标签，不过它只能在循环中使用。
    ```go
    **Loop**:
    	for n := 0; n < len(src); n += size {
    		switch {
    		case src[n] < sizeOne:
    			if validateOnly {
    				break
    			}
    			size = 1
    			update(src[n])
    		case src[n] < sizeTwo:
    			if n+1 >= len(src) {
    				err = errShortInput
    				break **Loop**
    			}
    			if validateOnly {
    				break
    			}
    			size = 2
    			update(src[n] + src[n+1]<<shift)
    		}
    	}
    ```
- 作为这一节的结束，此程序通过使用两个 `switch` 语句对字节数组进行比较：
    ```go
    // Compare 按字典顺序比较两个字节切片并返回一个整数。
    // 若 a == b，则结果为零；若 a < b；则结果为 -1；若 a > b，则结果为 +1。
    func Compare(a, b []byte) int {
    	for i := 0; i < len(a) && i < len(b); i++ {
    		switch {
    		case a[i] > b[i]:
    			return 1
    		case a[i] < b[i]:
    			return -1
    		}
    	}
    	switch {
    	case len(a) > len(b):
    		return 1
    	case len(a) < len(b):
    		return -1
    	}
    	return 0
    }
    ```
- `switch` 也可用于判断接口变量的动态类型。如 **类型选择** 通过圆括号中的关键字 `type` 使用类型断言语法。若 `switch` 在表达式中声明了一个变量，那么该变量的每个子句中都将有该变量对应的类型。
    ```go
    var t interface{}
    t = functionOfSomeType()
    switch t := t.(type) {
    default:
    	fmt.Printf("unexpected type %T", t)       // %T 输出 t 是什么类型
    case bool:
    	fmt.Printf("boolean %t\n", t)             // t 是 bool 类型
    case int:
    	fmt.Printf("integer %d\n", t)             // t 是 int 类型
    case *bool:
    	fmt.Printf("pointer to boolean %t\n", *t) // t 是 *bool 类型
    case *int:
    	fmt.Printf("pointer to integer %d\n", *t) // t 是 *int 类型
    }
    ```


### for

- Go的 `for` 循环类似于C，不再有`while`、 `do-while` 了
    ```go
    // 形式一：如同C的for循环
    for init; condition; post { }
    // 形式二：降级为C的while循环
    for condition { }
    // 形式三：降级为空子句的C的for(;;)循环
    for { }
    ```
- 短变量声明能方便的在循环中声明下标变量：
    ```go
    sum := 0
    for i := 0; i < 10; i++ {
    	sum += i
    }
    ```
- `range` 子句遍历数组、切片、字符串或者映射，或从信道中读取消息
    ```go
    for key, value := range oldMap {
    	newMap[key] = value
    }
    ```
    - 若你只需要该遍历中的第一个项（键或下标），去掉第二个就行了：
        ```go
        for key := range m {
        	if key.expired() {
        		delete(m, key)
        	}
        }
        ```
    - 若你需要丢弃第一个项（键或下标）请使用**空白标识符**，即下划线来丢弃第一个值：
        ```go
        sum := 0
        for _, value := range array {
        	sum += value
        }
        ```
    - 对于字符串，`range` 能够提供更多便利。它能通过解析UTF-8， 将每个独立的Unicode码点分离出来。错误的编码将占用一个字节，并以符文U+FFFD来代替。 （名称“符文”和内建类型 `rune` 是Go对单个Unicode码点的成称谓。 详情见[语言规范](http://golang.org/ref/spec#%E7%AC%A6%E6%96%87%E5%AD%97%E9%9D%A2)）。循环
        ```go
        for pos, char := range "日本\x80語" { // \x80 是个非法的UTF-8编码
        	fmt.Printf("字符 %#U 始于字节位置 %d\n", char, pos)
        }
        字符 U+65E5 '日' 始于字节位置 0
        字符 U+672C '本' 始于字节位置 3
        字符 U+FFFD '�' 始于字节位置 6
        字符 U+8A9E '語' 始于字节位置 7
        ```
- Go中没有逗号操作符，`++` 和 `--` 为语句而非表达式，所以 `for ; ; i++,j++`被拒绝。所以如要在 `for` 中使用多个变量，应采用**平行赋值**（i, j = i+1, j--）的方式 
    ```go
    // 反转 a
    for i, j := 0, len(a)-1; i < j; **i, j = i+1, j-1** {
    	a[i], a[j] = a[j], a[i]
    }
    ```


