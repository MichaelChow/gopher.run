// Append illustrates the behavior of the built-in append function.
// See page 88.

package main

import "fmt"

func appendslice(x []int, y ...int) []int {
	var z []int
	zlen := len(x) + len(y)
	if zlen <= cap(x) {
		// There is room to expand the slice.
		z = x[:zlen]
	} else {
		// There is insufficient space.
		// Grow by doubling, for amortized linear complexity.
		zcap := zlen
		if zcap < 2*len(x) {
			zcap = 2 * len(x)
		}
		z = make([]int, zlen, zcap)
		copy(z, x)
	}
	copy(z[len(x):], y)
	return z
}

// append函数对于理解slice底层是如何工作的非常重要，appendInt专门用于处理[]int类型的slice
// 函数将y添加到slice x中
func appendInt(x []int, y int) []int {
	var z []int
	zlen := len(x) + 1
	// 必须先检测slice x底层数组是否有足够的容量来保存新添加的元素y
	// 如果有，直接在原有的底层数组x之上扩展slice的len(x)。因此输出的z直接共享输入的x相同的底层数组。
	if zlen <= cap(x) {
		// There is room to grow.  Extend the slice.
		z = x[:zlen]
	} else {
		// There is insufficient space.  Allocate a new array.
		// Grow by doubling, for amortized linear complexity.
		// 如果没有足够的增长空间：
		// 为了提高内存使用效率，新分配的数组一般略大于保存x和y所需要的最低大小。
		// 通过在每次扩展数组时直接将长度翻倍，从而避免了多次内存分配，也确保了添加单个元素的平均时间是一个常数时间。
		zcap := zlen
		if zcap < 2*len(x) {
			zcap = 2 * len(x)
		}
		// 分配一个新的数组z，将长度设置为zlen，容量设置为zcap。因此输出的z与输入的x引用的将是不同的底层数组。
		z = make([]int, zlen, zcap)
		// 虽然通过循环复制元素更直接，不过内置的copy函数可以方便地将一个slice复制另一个相同类型的slice。
		// copy函数的第一个参数是要复制的目标slice，第二个参数是源slice，目标和源的位置顺序和dst = src赋值语句是一致的。两个slice可以共享同一个底层数组，甚至有重叠也没有问题。copy函数将返回成功复制的元素的个数（我们这里没有用到），等于两个slice中较小的长度，所以我们不用担心覆盖会超出目标slice的范围。
		copy(z, x) // a built-in function; see text
	}
	//  最后将新添加的y元素复制到新扩展的空间z[len(x)]，并返回slice。
	z[len(x)] = y
	return z
}

func main() {
	var x, y []int
	for i := 0; i < 10; i++ {
		y = appendInt(x, i)
		fmt.Printf("i=%d  cap=%d\t y=%v\n", i, cap(y), y)
		x = y
	}
}

/*
//!+output
i=0  cap=1       y=[0]
i=1  cap=2       y=[0 1]
i=2  cap=4       y=[0 1 2]
i=3  cap=4       y=[0 1 2 3]
i=4  cap=8       y=[0 1 2 3 4]
i=5  cap=8       y=[0 1 2 3 4 5]
i=6  cap=8       y=[0 1 2 3 4 5 6]
i=7  cap=8       y=[0 1 2 3 4 5 6 7]
i=8  cap=16      y=[0 1 2 3 4 5 6 7 8]
i=9  cap=16      y=[0 1 2 3 4 5 6 7 8 9]
//!-output
*/
