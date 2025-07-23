// Server3 is a minimal "echo" and counter server.
// See page 20.
package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var mu sync.Mutex
var count int
var urllog []string

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/count", counter)
	fmt.Println("start web server http://localhost:8000 ...")
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// handler echoes the Path component of the requested URL.
func handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	s := fmt.Sprintf("%s: %s\n", time.Now(), r.URL.Path)
	urllog = append(urllog, s)
	mu.Unlock()
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

// counter echoes the number of calls so far.
func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", count)
	fmt.Fprintf(w, "urllog: %v\n", urllog)
	mu.Unlock()
}
