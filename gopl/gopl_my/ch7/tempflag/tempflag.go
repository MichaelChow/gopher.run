// Tempflag prints the value of its -temp (temperature) flag.
// See page 181.
package main

import (
	"flag"
	"fmt"
	"gopher.run/go/src/ch7/tempconv"
)

var temp = tempconv.CelsiusFlag("temp", 20.0, "the temperature")

func main() {
	flag.Parse()
	fmt.Println(*temp)
}
