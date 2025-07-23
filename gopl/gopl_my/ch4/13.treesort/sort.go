// Package treesort provides insertion sort using an unbalanced binary tree.
// See page 101.
// 程序使用一个二叉树来实现一个插入排序：
package treesort

// 结构体类型的零值是每个成员都是零值。通常会将零值作为最合理的默认值
// 如，对于bytes.Buffer类型，结构体初始值就是一个随时可用的空缓存
// sync.Mutex的零值也是有效的未锁定状态
// 有时候这种零值可用的特性是自然获得的，但是也有些类型需要一些额外的工作

// 一个命名为S的结构体类型将不能再包含S类型的成员：因为一个聚合的值不能包含它自身。（该限制同样适用于数组。）
// 但S类型的结构体可以包含*S指针类型的成员，这可以让我们创建递归的数据结构，比如链表和树结构等。
// 声明一个struct，二叉树
type tree struct {
	value       int
	left, right *tree
}

// Sort 函数对整数切片进行排序
func Sort(values []int) {
	// 初始化根节点为空
	var root *tree
	// 遍历切片中的每个元素
	for _, v := range values {
		// 将元素插入到二叉树中
		root = add(root, v)
	}
	// 将排序后的元素追加到原始切片的前缀
	appendValues(values[:0], root)
}

// add 函数向二叉树 t 中插入一个值为 value 的节点，并返回插入后的二叉树
func add(t *tree, value int) *tree {
	// 如果树为空，则创建一个新的树节点
	if t == nil {
		// 等价于返回 &tree{value: value}
		t = new(tree)
		t.value = value
		return t
	}
	// 如果值小于当前节点的值，则将其插入到左子树中
	if value < t.value {
		t.left = add(t.left, value)
		// 如果值大于等于当前节点的值，则将其插入到右子树中
	} else {
		t.right = add(t.right, value)
	}
	// 返回插入后的二叉树
	return t
}

// appendValues 函数将二叉树 t 中的元素按顺序追加到 values 切片中，并返回结果切片
func appendValues(values []int, t *tree) []int {
	// 如果树不为空
	if t != nil {
		// 递归地将左子树中的元素追加到 values 中
		values = appendValues(values, t.left)
		// 将当前节点的值追加到 values 中
		values = append(values, t.value)
		// 递归地将右子树中的元素追加到 values 中
		values = appendValues(values, t.right)
	}
	// 返回最终的排序结果
	return values
}
