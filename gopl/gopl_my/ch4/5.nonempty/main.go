// Nonempty is an example of an in-place slice algorithm.
// See page 91.
package main

import "fmt"

// nonempty returns a slice holding only the non-empty strings.
// The underlying array is modified during the call.
// 在原有string slice内存空间之上返回不包含空字符串的列表
// 比较微妙的地方是，输入的slice和输出的slice共享一个底层数组。
// 这可以避免分配另一个数组，不过原来的数据将可能会被覆盖
func nonempty(strings []string) []string {
	i := 0
	for _, s := range strings {
		if s != "" {
			strings[i] = s
			i++
		}
	}
	return strings[:i]
}

func main() {
	data := []string{"one", "", "three"}
	fmt.Printf("%q\n", nonempty(data)) // `["one" "three"]`
	// 这可以避免分配另一个数组，不过原来的数据将可能会被覆盖
	fmt.Printf("%q\n", data) // `["one" "three" "three"]`
	// 因此我们通常会这样使用nonempty函数，和append函数类似
	data = nonempty(data)
	fmt.Printf("%q\n", data) // `["one" "three"]`
}

// nonempty函数也可以使用append函数实现
func nonempty2(strings []string) []string {
	out := strings[:0] // zero-length slice of original
	for _, s := range strings {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}
