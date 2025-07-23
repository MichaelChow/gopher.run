// 程序使用map来记录提交相同的字符串列表的次数。
package main

import "fmt"

// 有时候我们需要一个map或set的key是slice类型，但是map的key必须是可比较的类型，但是slice并不满足这个条件。
// 我们可以绕过这个限制：map的key依然用可比较类型（如string），但每次存取前先将slice转为string类型。
// 1. 定义一个辅助函数k（key），将string slice转为string，作为map的key（string类型）。确保只有x和y相等时k(x) == k(y)才成立。
// 这里使用了fmt.Sprintf函数将字符串列表转换为一个字符串。以用于map的key。通过%q参数忠实地记录每个字符串元素的信息。
func k(list []string) string { return fmt.Sprintf("%q", list) }

// 2. 创建一个key为string类型的map m。在每次对map操作时m[key]，处理成m[k(list)](用k辅助函数将slice转化为string类型)。
var m = make(map[string]int)

func Add(list []string)       { m[k(list)]++ }
func Count(list []string) int { return m[k(list)] }

// 使用同样的技术可以处理任何不可比较的key类型，而不仅仅是slice类型。
// 这种技术对于想使用自定义key比较函数的时候也很有用
// 如在比较字符串的时候忽略大小写。同时，辅助函数k(x)也不一定是字符串类型，它可以返回任何可比较的类型，例如整数、数组或结构体等。
func main() {
	Add([]string{"A", "B", "C"})
	Add([]string{"A", "B", "C"})
	Add([]string{"A", "B", "C"})
	Add([]string{"A", "B", "C"})
	Add([]string{"A", "B", "D"})
	fmt.Println(Count([]string{"A", "B", "C"})) // 2

}
