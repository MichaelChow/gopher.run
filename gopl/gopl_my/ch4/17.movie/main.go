// Movie prints Movies as JSON.
// See page 108.
// 程序负责收集各种电影评论并提供反馈功能
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Movie struct {
	Title string
	// 在结构体声明中，Year和Color成员后面的字符串面值是结构体成员Tag
	// 结构体的成员Tag可以是任意的字符串面值，但是通常是一系列用空格分隔的key:"value"键值对序列
	// 因为值中含有双引号字符，因此成员Tag一般用原生字符串面值（即反引号`包裹）的形式书写，这样不再需要转义字符来表示双引号
	// "json"开头键名对应的值，用于控制encoding/json包的编码和解码的行为，并且encoding/...下面其它的包也遵循这个约定
	// 值的第一部分用于指定JSON对象的名字，如将Go语言中的TotalCount成员对应到JSON中的total_count对象
	// Year名字的成员在编码后变成了released，还有Color成员编码后变成了小写字母开头的color
	// 一个结构体成员Tag是和在编译阶段关联到该成员的元信息字符串
	Year int `json:"released"`
	// 额外的omitempty选项，表示当Go语言结构体成员为空或零值时不生成该JSON对象（这里false为零值）。
	// 即Marshal的json串中"Casablanca"由于其为零值fasle，所以没有Color字段
	Color  bool `json:"color,omitempty"`
	Actors []string
}

var movies = []Movie{
	{Title: "Casablanca", Year: 1942, Color: false,
		Actors: []string{"Humphrey Bogart", "Ingrid Bergman"}},
	{Title: "Cool Hand Luke", Year: 1967, Color: true,
		Actors: []string{"Paul Newman"}},
	{Title: "Bullitt", Year: 1968, Color: true,
		Actors: []string{"Steve McQueen", "Jacqueline Bisset"}},
	// ...
}

func main() {
	{
		// 将一个Go语言中类似movies的 结构体slice 转为 JSON 的过程叫编组（marshaling） 美 /ˈmɑːrʃl/
		// 编组通过调用json.Marshal函数完成，返回一个编码后的字节slice，包含很长的字符串，并且没有空白缩进以紧凑的表示
		// 类似Java中的序列化为Json串
		data, err := json.Marshal(movies)
		if err != nil {
			log.Fatalf("JSON marshaling failed: %s", err)
		}
		fmt.Printf("%s\n", data)
	}

	{
		// 格式化json输出，产生整齐缩进的输出
		// MarshalIndent()有两个额外的字符串参数，用于表示每一行输出的前缀和每一个层级的缩进：
		// 译注：在最后一个成员或元素后面并没有逗号分隔符

		// 在编码时，默认使用Go语言结构体的成员名字作为JSON的对象（通过reflect反射技术，我们将在12.6节讨论）。
		// 只有导出的结构体成员才会被编码，这也就是我们为什么选择用大写字母开头的成员名称。
		data, err := json.MarshalIndent(movies, "", "    ")
		if err != nil {
			log.Fatalf("JSON marshaling failed: %s", err)
		}
		fmt.Printf("%s\n", data)

		// 编码的逆操作是解码，对应将JSON数据解码为Go语言的数据结构，Go语言中一般叫unmarshaling
		// 类似Java中的将json串反序列化为对象
		// 代码将JSON格式的电影数据data解码为一个结构体slice，结构体中只有Title成员
		var titles []struct{ Title string }
		// 通过定义合适的Go语言数据结构，我们可以选择性地解码JSON中感兴趣的成员
		if err := json.Unmarshal(data, &titles); err != nil {
			log.Fatalf("JSON unmarshaling failed: %s", err)
		}
		// 这里的slice将被只含有Title信息的值填充，其它JSON成员将被忽略
		fmt.Println(titles) //  [{Casablanca} {Cool Hand Luke} {Bullitt}]
	}
}

/*
//!+output
[{"Title":"Casablanca","released":1942,"Actors":["Humphrey Bogart","Ingr
id Bergman"]},{"Title":"Cool Hand Luke","released":1967,"color":true,"Ac
tors":["Paul Newman"]},{"Title":"Bullitt","released":1968,"color":true,"
Actors":["Steve McQueen","Jacqueline Bisset"]}]
//!-output
*/

/*
//!+indented
[
    {
        "Title": "Casablanca",
        "released": 1942,
        "Actors": [
            "Humphrey Bogart",
            "Ingrid Bergman"
        ]
    },
    {
        "Title": "Cool Hand Luke",
        "released": 1967,
        "color": true,
        "Actors": [
            "Paul Newman"
        ]
    },
    {
        "Title": "Bullitt",
        "released": 1968,
        "color": true,
        "Actors": [
            "Steve McQueen",
            "Jacqueline Bisset"
        ]
    }
]
//!-indented
*/
