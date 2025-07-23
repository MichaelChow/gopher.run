package main

import (
	"time"
)

// 声明一个叫Employee的命名的结构体类型
// 结构体类型往往是冗长的，因为它的每个成员可能都会占一行
// 因此，完整的结构体写法通常只在类型声明语句的地方出现
type Employee struct {
	// 通常一行对应一个结构体成员，成员的名字在前，类型在后
	ID int
	// 不过如果相邻的成员类型如果相同的话可以被合并到一行
	// 结构体成员的不同的输入顺序，定义了不同的结构体类型
	Name, Address string
	DoB           time.Time
	// 通常，我们只是将相关的成员写到一起
	Position string
	// 如果结构体成员名字是以大写字母开头的，那么该成员就是导出的
	// 这是Go语言导出规则决定的。一个结构体可能同时包含导出和未导出的成员
	Salary    int
	ManagerID int
}

func main() {
	// 声明了一个Employee类型的变量dilbert
	var dilbert Employee
	// dilbert结构体变量的成员可以通过点操作符访问
	// 因为dilbert是一个变量，它所有的成员也同样是变量，我们可以直接对每个成员赋值：
	dilbert.Salary -= 5000 // demoted, for writing too few lines of code
	// 或者是对成员取地址，然后通过指针访问：
	position := &dilbert.Position
	*position = "Senior " + *position // promoted, for outsourcing to Elbonia

	// 点操作符也可以和指向结构体的指针一起工作：
	var employeeOfTheMonth *Employee = &dilbert
	employeeOfTheMonth.Position += " (proactive team player)"
	// 相当于下面语句
	(*employeeOfTheMonth).Position += " (proactive team player)"

}
