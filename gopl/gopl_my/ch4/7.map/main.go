package main

import (
	"fmt"
	"sort"
)

func main() {
	// make 创建一个空map
	ages1 := make(map[string]int)
	fmt.Println(ages1, ages1 == nil, len(ages1) == 0) // map[] false true
	ages1["alice"] = 31
	ages1["charlie"] = 34

	// 用map字面值的语法创建一个空map
	ages2 := map[string]int{}
	fmt.Println(ages2, ages2 == nil, len(ages2) == 0) // map[] false true

	// map的零值是nil，也就是没有引用任何hash表
	var ages3 map[string]int
	fmt.Println(ages3, ages3 == nil, len(ages3) == 0) // map[] true true

	// 用map字面值的语法创建一个非空map
	ages := map[string]int{
		"alice":   31,
		"charlie": 34,
	}
	// 通过key对应的下标语法访问map中对应的元素
	fmt.Println(ages["alice"])
	// 使用内置的delete函数删除元素
	delete(ages, "alice")
	// 如果key不存在，那么将返回该map元素类型的零值
	fmt.Println(ages["alice"]) // 这里返回int类型的零值，0
	ages["Bob"] = ages["Bob"] + 1
	ages["Bob"] += 1 //等价写法
	ages["Bob"]++    //等价写法

	// map的大部分操作，包括查找、删除、len和range都可以安全工作在nil值的map上，它们的行为和一个空map类似。
	// 但向一个nil值的map存入元素奖导致一个panic异常
	// ages3["carol"] = 21 // panic: assignment to entry in nil map

	// map中的元素并不是一个变量，因此我们不能对map的元素进行取址操作：
	// 禁止对map元素取址的原因：map可能随着元素数量的增长而分配更大的内存空间，从而可能导致之前的地址无效。
	// _ := &ages["bob"]  // compile error: cannot take address of map element

	// 遍历map中全部的key/value对：可使用range风格的for循环遍历key/value对（同slice遍历语法）
	for name, age := range ages {
		fmt.Printf("%s\t%d\n", name, age)
	}

	// map的迭代顺序是不确定的，并且不同的哈希函数实现可能导致不同的遍历顺序。
	// 在实践中，遍历的顺序是随机的，每一次遍历的顺序都不相同。这是故意的，每次都使用随机的遍历顺序可以强制要求程序不会依赖具体的哈希函数实现。

	// 如果一定要按顺序遍历key/value对，我们必须显式地对map的key进行排序，可以使用sort包的Strings函数对字符串slice进行排序：
	// 创建一个空slice，由于name的最终大小完全确定，容量直接设为map的长度
	names := make([]string, 0, len(ages))
	// 取出map中的所有key
	for name := range ages {
		names = append(names, name)
	}
	// 对key进行排序
	sort.Strings(names)
	// 遍历排序后的slice中的key，根据key取出map中的value
	for _, name := range names {
		fmt.Printf("%s\t%d\n", name, ages[name])
	}
	equal(map[string]int{"A": 0}, map[string]int{"B": 42})

}

// 和slice一样，map之间也不能进行相等比较；唯一的例外是和nil进行比较。
// 要判断两个map是否包含相同的key和value，我们必须通过一个循环实现：

func equal(x, y map[string]int) bool {
	if len(x) != len(y) {
		return false
	}
	for k, xv := range x {
		// 如果key不存在，那么将返回该map元素类型的零值。（这个规则很实用）
		// 但如果元素是一个数字，就需要区分一个已经存在元素的0和不存在元素返回的默认零值的0。
		// 所以我们需要使用一个额外的布尔值变量ok（一般命名为ok）来表示元素是否存在。
		if yv, ok := y[k]; !ok || yv != xv {
			return false
		}
	}
	return true
}
