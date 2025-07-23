package main

import (
	"bytes"
	"io"
)

const debug = false

func main() {
	var buf *bytes.Buffer // bug代码
	//var buf io.Writer // 修复代码
	if debug {
		buf = new(bytes.Buffer) // enable collection of output
	}
	f(buf) // NOTE: subtly incorrect!
	if debug {
		// ...use buf...
	}
}

// If out is non-nil, output will be written to it.
func f(out io.Writer) {
	// ...do something...
	if out != nil {
		out.Write([]byte("done!\n"))
	}
}
