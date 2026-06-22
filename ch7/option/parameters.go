//go:build ignore
// +build ignore

package main

import "fmt"

func Sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func main() {
	s := Sum(1, 2, 3) // 直接传入多个参数
	arr := []int{1, 2, 3}
	s2 := Sum(arr...) // 可以展开切片

	fmt.Println(s)
	fmt.Println(s2)
}
