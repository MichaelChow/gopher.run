package main

import "fmt"

func main() {
	// 无论如何实现，以这种方式重用一个slice一般都要求最多为每个输入值产生一个输出值
	// 事实上很多这类算法都是用来过滤或合并序列中相邻的元素。这种slice用法是比较复杂的技巧，虽然使用到了slice的一些技巧，但是对于某些场合是比较清晰和有效的。
	// 一个slice可以用来模拟一个stack。最初给定的空slice对应一个空的stack，然后可以使用append函数将新的值压入stack：
	var stack = make([]int, 0)
	v := 1
	stack = append(stack, v) // push v
	// stack的顶部位置对应slice的最后一个元素：
	top := stack[len(stack)-1] // top of stack
	// 通过收缩stack可以弹出栈顶的元素
	stack = stack[:len(stack)-1] // pop
	fmt.Println(top)

	s := []int{5, 6, 7, 8, 9}
	fmt.Println(remove(s, 2)) // [5 6 8 9]
	s = []int{5, 6, 7, 8, 9}
	fmt.Println(remove2(s, 2)) // [5 6 9 8]
}

func remove(slice []int, i int) []int {
	// // 要删除slice中间的某个元素并保存原有的元素顺序，可以通过内置的copy函数将后面的子slice向前依次移动一位完成：
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]

}

// 如果删除元素后不用保持原来顺序的话，我们可以简单的用最后一个元素覆盖被删除的元素：
func remove2(slice []int, i int) []int {
	slice[i] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}
