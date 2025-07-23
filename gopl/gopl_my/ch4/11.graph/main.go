// Graph shows how to use a map of maps to represent a directed graph.
// See page 99.
package main

import "fmt"

// Map的value类型也可以是一个聚合类型（如是一个map或slice）。
// 从概念上讲，graph将一个字符串类型的key，映射到一组相关的字符串集合，它们指向新的graph的key。
var graph = make(map[string]map[string]bool)

func addEdge(from, to string) {
	// addEdge函数惰性初始化map是一个惯用方式，也就是说在每个值首次作为key时才初始化。
	// addEdge函数显示了如何让map的零值也能正常工作；即使from到to的边不存在，graph[from][to]依然可以返回一个有意义的结果。
	edges := graph[from]
	if edges == nil {
		edges = make(map[string]bool)
		graph[from] = edges
	}
	edges[to] = true
}

func hasEdge(from, to string) bool {
	return graph[from][to]
}

func main() {
	addEdge("a", "b")
	addEdge("c", "d")
	addEdge("a", "d")
	addEdge("d", "a")
	fmt.Println(hasEdge("a", "b"))
	fmt.Println(hasEdge("c", "d"))
	fmt.Println(hasEdge("a", "d"))
	fmt.Println(hasEdge("d", "a"))
	fmt.Println(hasEdge("x", "b"))
	fmt.Println(hasEdge("c", "d"))
	fmt.Println(hasEdge("x", "d"))
	fmt.Println(hasEdge("d", "x"))
}
