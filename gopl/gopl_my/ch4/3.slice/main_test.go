package main_test

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	//months := [...]string{1: "January", 2: "February", 3: "March", 4: "April", 5: "May", 6: "June", 7: "July", 8: "Augest", 9: "September", 10: "October", 11: "November", 12: "December"}
	//fmt.Printf("%T\n", months)
	//Q2 := months[4:7]
	//fmt.Printf("%T\n", Q2)
	//summer := months[6:9]
	//fmt.Printf("%T\n", summer)
	//fmt.Println(Q2, len(Q2), cap(Q2)) // [April May June] 3 9
	//fmt.Println(Q2[:9])               // [April May June July Augest September October November December]
	//fmt.Println(Q2[:10])              // panic: runtime error: slice bounds out of range [:10] with capacity 9

	//age := make(map[string]int)
	//var age map[string]int
	age := map[string]int{}
	fmt.Println(age == nil)
	age["Jan"] = 1
}
