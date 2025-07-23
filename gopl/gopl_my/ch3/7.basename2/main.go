// Basename2 reads file names from stdin and prints the base name of each one.
// See page 72.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		fmt.Println(basename(input.Text()))
	}
	// NOTE: ignoring potential errors from input.Err()
}

// basename removes directory components and a trailing .suffix.
// e.g., a => a, a.go => a, a/b/c.go => c, a/b.c.go => b.c
// !+
func basename(s string) string {

	slash := strings.LastIndex(s, "/") // -1 if "/" not found
	// 当前字符串中没有"/"时，slash为-1，此时s[slash+1:]刚好为s[0:]，即整个字符串
	s = s[slash+1:]
	if dot := strings.LastIndex(s, "."); dot >= 0 {
		s = s[:dot]
	}
	return s
}
