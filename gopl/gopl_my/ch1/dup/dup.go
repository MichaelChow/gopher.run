// Package dup prints the text of each line that appears more than. Press Ctrl+D (on Unix-like systems) or Ctrl+Z (on Windows systems) in the terminal to indicate the end of input.
package dup

import (
	"bufio"
	"fmt"
	"os"
)

func v1() {
	// The make built-in function allocates and initializes an object of type
	// slice, map, or chan (only). Like new, the first argument is a type, not a
	// value. Unlike new, make's return type is the same as the type of its
	// argument, not a pointer to it. The specification of the result depends on
	// the type:
	//
	//	Slice: The size specifies the length. The capacity of the slice is
	//	equal to its length. A second integer argument may be provided to
	//	specify a different capacity; it must be no smaller than the
	//	length. For example, make([]int, 0, 10) allocates an underlying array
	//	of size 10 and returns a slice of length 0 and capacity 10 that is
	//	backed by this underlying array.
	//	Map: An empty map is allocated with enough space to hold the
	//	specified number of elements. The size may be omitted, in which case
	//	a small starting size is allocated.
	//	Channel: The channel's buffer is initialized with the specified
	//	buffer capacity. If zero, or the size is omitted, the channel is
	//	unbuffered.
	counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		counts[input.Text()]++
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%s\t%d\n", line, n)
		}
	}
}

//
//func main() {
//	v1()
//}
