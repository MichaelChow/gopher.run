// The toposort program prints the nodes of a DAG in topological order.
// See page 136.
// 给定一些计算机课程，每个课程都有前置课程，只有完成了前置课程才可以开始当前课程的学习
// 我们的目标是选择出一组课程，这组课程必须确保按顺序学习时，能全部被完成
package main

import (
	"fmt"
	"sort"
)

// !+table
// prereqs maps 记录了每个课程的前置课程
var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},

	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},

	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

//!-table

func main() {
	// 调用 topoSort 函数对 prereqs 进行拓扑排序，并将结果存储在 courses 变量中
	courses := topoSort(prereqs)
	// 遍历排序后的课程列表
	for i, course := range courses {
		// 打印课程的序号和名称
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

// 这类问题被称作拓扑排序。从概念上说，前置条件可以构成有向图。图中的顶点表示课程，边表示课程间的依赖关系
// 显然，图中应该无环，这也就是说从某点出发的边，最终不会回到该点
// 下面的代码用深度优先搜索了整张图，获得了符合要求的课程序列

// topoSort 函数对给定的 map 进行拓扑排序，并返回排序后的字符串切片
func topoSort(m map[string][]string) []string {
	// 用于存储排序后的结果
	var order []string
	// 用于标记已经访问过的节点
	seen := make(map[string]bool)
	// 当匿名函数需要被递归调用时，我们必须首先声明一个变量 visitAll

	// 定义一个匿名函数，用于递归访问所有相关节点
	var visitAll func(items []string)
	// 再将匿名函数赋值给这个变量
	// 如果不分成两步，函数字面量无法与visitAll绑定，我们也无法递归调用该匿名函数

	// 将匿名函数赋值给 visitAll 变量，以便在函数内部调用
	visitAll = func(items []string) {
		// 遍历传入的节点列表
		for _, item := range items {
			// 如果节点没有被访问过
			if !seen[item] {
				// 标记节点为已访问
				seen[item] = true
				// 递归访问该节点的所有依赖项
				visitAll(m[item])
				// 将当前节点添加到排序结果中
				order = append(order, item)
			}
		}
	}

	// 获取 map 中的所有键，即所有节点
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}

	// 对节点列表进行排序，确保结果的稳定性
	sort.Strings(keys)
	// 从排序后的节点列表开始，调用 visitAll 函数进行深度优先搜索，开始拓扑排序 topoSort
	visitAll(keys)
	// 返回排序后的节点列表
	return order
}
