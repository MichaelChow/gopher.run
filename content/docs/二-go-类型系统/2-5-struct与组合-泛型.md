---
title: "2.5 struct与组合、泛型"
date: 2025-04-03T08:02:00Z
draft: false
weight: 2005
---

# 2.5 struct与组合、泛型

## struct

struct是**将零个或多个****任意类型的命名变量组合****的聚合数据类型**，每个变量都叫做**结构体的成员。**

**不同的成员顺序为不同的结构体类型****。结构体的成员是变量（可寻址的值），可以取地址。**

**example**：数据处理领域用结构体的经典案例是员工信息记录 **Employee /ɪmˈplɔɪiː/。**

```go
type Employee struct {
	ID int                // 工号
	Name, Address string  // 姓名、地址，相同的成员类型且相关的成员，合并到一行
	DoB       time.Time   // 出生日期 
	Position string      // 职位
	Salary    int        // 薪水
	ManagerID int        // 直属领导
}
```

结构体类型的值可以通过**结构体字面量**来设置，即通过设置结构体的成员变量来设置。

```go
type Point struct { X, Y int }
p := Point{1,2}  // 初始化写法1：要求完全按顺序写；仅在定义结构体的包内部使用，或者是在较小的结构体中使用
p := point.Point{X: 1, Y: 2}  // 初始化写法2：name:value，与顺序无关，默认用成员变量的零值
```

```go
var dilbert Employee
dilbert.Salary -= 5000

position := &dilbert.Position
*position = "Senior " + *position 

// 点操作符也可以和指向结构体的指针一起工作：
var employeeOfTheMonth *Employee = &dilbert

employeeOfTheMonth.Position += "(proactive team player)"
(*employeeOfTheMonth).Position += " (proactive team player)"  // 等价写法
```

**较大的结构体通常会用指针的方式传入和返回，（避免分配一份大结构体的内存再拷贝其数据的操作）；**（在Go语言中，所有的函数参数都是值拷贝传入的，函数参数将不再是函数调用时的原始变量）

```go
func Scale(p point.Point, factor int) point.Point {  // 传结构体，赋值给函数的新的局部变量，再返回结构体再赋值
	return point.Point{p.X * factor, p.Y * factor}
}

func Bonus(p *point.Point, factor int) {   // 传结构体指针，无需返回
    p.X *= factor
	  p.Y *= factor
}
```



**可比较性：****如果结构体的全部成员都是可以比较的，那么结构体也是可以比较的**，可以使用==或!=运算符进行比较两个结构体的每个成员

```go
type Point struct{ X, Y int }

p := Point{1, 2}
q := Point{2, 1}
fmt.Println(p.X == q.X && p.Y == q.Y) // "false"
fmt.Println(p == q)                   // "false"
```

```go
// 可比较的结构体类型（和其他可比较的类型一样）可用于map的key类型
type address struct {
    hostname string
    port     int
}

hits := make(map[address]int)
hits[address{"golang.org", 443}]++
```



**结构体组合和匿名成员**：Go中不同寻常的**结构体嵌套机制**，将一个命名结构体当做另一个结构体类型的匿名成员使用。并可以通过简单的表达式就可以访问嵌套的成员（如x.d.e.f）

```go
type Circle struct {   
    X, Y, Radius int  // 圆形类型，圆心的X、Y坐标、半径
}
type Wheel struct {
    X, Y, Radius, Spokes int // 轮形类型，增加了Spokes 径向辐条数
}
type Point struct {      // 独立相同的属性
    X, Y int
}
type Circle struct {   // 改进后
    Center Point
    Radius int
}
type Wheel struct {   // 改进后
    Circle Circle
    Spokes int
}
```

```go
var w Wheel
w.Circle.Center.X = 8     // 但**访问每个成员变得繁琐**
w.Circle.Center.Y = 8
w.Circle.Radius = 5
w.Spokes = 20
```

Go允许定义不带名称的结构体成员（即**匿名成员**），只需要指定类型即可；**匿名成员实际拥有隐式的名字（即对应类型的名字，如Point、Circle）**，只是在点操作符时可以省略，仅仅是点操作符的语法糖。

```go
type Circle struct {
    Point            // Point类型被嵌入到了Circle结构体
    Radius int
}

type Wheel struct {
    Circle           // Circle类型被嵌入到了Wheel结构体
    Spokes int
}
```

```go
var w Wheel
w.X = 8            // 等价写法：**w.Circle.Point.X = 8**
w.Radius = 5       // 等价写法： w.Circle.Radius = 5
```



匿名成员的类型不能是匿名类型; **匿名成员的类型必须是 任何一个命名类型或指向命名类型的指针(任何，不仅仅是结构体类型)；**

**嵌套一个非结构体类型（没有子成员）有什么用？以该快捷方式访问匿名成员的内部变量同样适用于访问匿名成员的方法。这个机制就是从简单类型对象组合成复杂的组合类型的主要方式，Go中组合是面向对象编程方式的核心；**

遗憾的是结构体字面值并没有简短表示匿名成员的语法，结构体字面值必须遵循 形状类型声明时的结构；

```go
	w = Wheel{Circle{Point{8, 8}, 5}, 20}
```



### example: treesort 

```go
// 使用一个二叉树来实现一个插入排序
package treesort

// 结构体类型的零值是每个成员都是零值。通常会将零值作为最合理的默认值
// 如，对于bytes.Buffer类型，结构体初始值就是一个随时可用的空缓存
// sync.Mutex的零值也是有效的未锁定状态
// 有时候这种零值可用的特性是自然获得的，但是也有些类型需要一些额外的工作

// 一个命名为S的结构体类型将不能再包含S类型的成员：因为一个聚合的值不能包含它自身。（该限制同样适用于数组。）
// 但S类型的结构体可以包含*S指针类型的成员，这可以让我们创建递归的数据结构，比如链表和树结构等。
// 声明一个struct，二叉树
type tree struct {
	value       int
	left, right *tree
}

// Sort 函数对整数切片进行排序
func Sort(values []int) {
	// 初始化根节点为空
	var root *tree
	// 遍历切片中的每个元素
	for _, v := range values {
		// 将元素插入到二叉树中
		root = add(root, v)
	}
	// 将排序后的元素追加到原始切片的前缀
	appendValues(values[:0], root)
}

// add 函数向二叉树 t 中插入一个值为 value 的节点，并返回插入后的二叉树
func add(t *tree, value int) *tree {
	// 如果树为空，则创建一个新的树节点
	if t == nil {
		// 等价于返回 &tree{value: value}
		t = new(tree)
		t.value = value
		return t
	}
	// 如果值小于当前节点的值，则将其插入到左子树中
	if value < t.value {
		t.left = add(t.left, value)
		// 如果值大于等于当前节点的值，则将其插入到右子树中
	} else {
		t.right = add(t.right, value)
	}
	// 返回插入后的二叉树
	return t
}

// appendValues 函数将二叉树 t 中的元素按顺序追加到 values 切片中，并返回结果切片
func appendValues(values []int, t *tree) []int {
	// 如果树不为空
	if t != nil {
		// 递归地将左子树中的元素追加到 values 中
		values = appendValues(values, t.left)
		// 将当前节点的值追加到 values 中
		values = append(values, t.value)
		// 递归地将右子树中的元素追加到 values 中
		values = appendValues(values, t.right)
	}
	// 返回最终的排序结果
	return values
}
```



```go
// 用单元测试上述代码

package treesort_test

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	treesort "gopher.run/go/src/ch4/13.treesort"
)

func TestSort(t *testing.T) {
	data := make([]int, 50)
	for i := range data {
		data[i] = rand.Int() % 50
	}
	fmt.Println(data)
	treesort.Sort(data)
	fmt.Println(data)
	if !sort.IntsAreSorted(data) {
		t.Errorf("not sorted: %v", data)
	}
}

=== RUN   TestSort
[39 38 30 39 23 1 18 26 33 22 6 26 1 48 32 14 35 38 27 36 35 49 6 44 21 38 32 9 30 8 1 48 29 23 45 27 20 0 37 1 22 11 24 44 24 38 10 21 24 42]
[0 1 1 1 1 6 6 8 9 10 11 14 18 20 21 21 22 22 23 23 24 24 24 26 26 27 27 29 30 30 32 32 33 35 35 36 37 38 38 38 38 39 39 42 44 44 45 48 48 49]
--- PASS: TestSort (0.00s)
PASS
```

```go
// 如果结构体没有任何成员的话就是空结构体，写作struct{}。它的大小为0，也不包含任何信息，但是有时候依然是有价值的。
// 有些Go语言程序员用map来模拟set数据结构时，用它来代替map中布尔类型的value，只是强调key的重要性，但是因为节约的空间有限，而且语法比较复杂，所以我们通常会避免这样的用法。
seen := make(map[string]struct{}) // set of strings
// ...
if _, ok := seen[s]; !ok {
    seen[s] = struct{}{}
    // ...first time seeing s...
}
```





## **泛型**

**generics**

**类型参数化**：泛型不是类型，而是类型的"模板"。核心是类型参数T，使用泛型时，需要指定具体的类型。在编译时自动推导。

```go
// Go编译器可以在编译时推导类型
func main() {
    // 编译器自动推导 T 为 int
    result1 := Min(5, 3)
    
    // 编译器自动推导 T 为 float64
    result2 := Min(3.14, 2.71)
    
    // 编译器自动推导 T 为 string
    result3 := Min("apple", "banana")
}

```

**泛型的类型约束：**

```go
// 约束泛型类型参数的行为

// comparable 约束：支持比较操作的类型
func Max[T comparable](a, b T) T {
    if a > b {
        return a
    }
    return b
}

// any 约束：任意类型
func Print[T any](item T) {
    fmt.Printf("%v\n", item)
}

// 自定义约束
type Number interface {
    ~int | ~float64 | ~int64
}

func Sum[T Number](items []T) T {
    var sum T
    for _, item := range items {
        sum += item
    }
    return sum
}
```

**约束组合:**

```go
// 组合多个约束
type Stringer interface {
    String() string
}

type Printable[T any] interface {
    Stringer
    Print()
}

// 使用组合约束
func Process[T Printable[T]](item T) {
    item.Print()
    fmt.Println(item.String())
}
```



```go
// 泛型可以应用于所有四类类型

// 1. 基础类型 + 泛型
func ProcessNumber[T int | float64](n T) T {
    return n * 2
}

// 2. 组合类型 + 泛型
type Array[T any] [10]T

// 3. 引用类型 + 泛型
type Slice[T any] []T
type Map[K comparable, V any] map[K]V

// 4. 接口类型 + 泛型
type Processor[T any] interface {
    Process(item T) error
}
```

**泛型的优势：**

1. **类型安全**：
```go
// 编译时类型检查
func main() {
    numbers := []int{1, 2, 3, 4, 5}
    
    // 类型安全：只能处理int类型
    filtered := Filter(numbers, func(n int) bool {
        return n > 3
    })
    
    // 编译错误：类型不匹配
    // filtered := Filter(numbers, func(s string) bool { return true })
}
```

1. **代码复用：**
```go
// 一个函数处理多种类型
func main() {
    // 处理int切片
    ints := []int{1, 2, 3}
    result1 := Map(ints, func(n int) int { return n * 2 })
    
    // 处理string切片
    strings := []string{"a", "b", "c"}
    result2 := Map(strings, func(s string) string { return strings.ToUpper(s) })
    
    // 处理float64切片
    floats := []float64{1.1, 2.2, 3.3}
    result3 := Map(floats, func(f float64) float64 { return f * 1.5 })
}
```

**3. 性能优化**

```go
// **编译时生成具体类型的代码**
// 没有运行时类型检查的开销

// 对于能在编译器确定的类型，**编译器会为每种具体类型生成专门的代码**
func main() {
    // 生成 int 版本的 Min 函数
    minInt := Min[int]
    
    // 生成 float64 版本的 Min 函数
    minFloat := Min[float64]
    
    // 生成 string 版本的 Min 函数
    minString := Min[string]
}
```



### 泛型的应用场景

**1. 容器类型**

```go
// 泛型切片
type GenericSlice[T any] []T

func (s *GenericSlice[T]) Add(item T) {
    *s = append(*s, item)
}

func (s GenericSlice[T]) Get(index int) T {
    return s[index]
}

// 泛型映射
type GenericMap[K comparable, V any] map[K]V

func (m GenericMap[K, V]) Set(key K, value V) {
    m[key] = value
}

func (m GenericMap[K, V]) Get(key K) (V, bool) {
    value, exists := m[key]
    return value, exists
}
```

**2. 工具函数**

```go
// 泛型工具函数
func Filter[T any](items []T, predicate func(T) bool) []T {
    var result []T
    for _, item := range items {
        if predicate(item) {
            result = append(result, item)
        }
    }
    return result
}

func Map[T, U any](items []T, transform func(T) U) []U {
    result := make([]U, len(items))
    for i, item := range items {
        result[i] = transform(item)
    }
    return result
}

func Reduce[T, U any](items []T, initial U, reducer func(U, T) U) U {
    result := initial
    for _, item := range items {
        result = reducer(result, item)
    }
    return result
}
```



```go
// 泛型函数
func Min[T comparable](a, b T) T {
    if a < b {
        return a
    }
    return b
}

// 泛型结构体
type Stack[T any] struct {
    items []T
}

// 泛型接口
type Container[T any] interface {
    Push(item T)
    Pop() T
    IsEmpty() bool
}
```





