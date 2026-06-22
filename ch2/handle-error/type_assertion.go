//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"
)

func main() {
	_, err := os.Open("nonexistent.txt")
	if err != nil {
		// 使用类型断言检查是否为 *os.PathError 类型
		if pathErr, ok := err.(*os.PathError); ok {
			fmt.Printf("Failed to %s file: %s\n", pathErr.Op, pathErr.Path)
			fmt.Println("Error message:", pathErr.Err)
		} else {
			fmt.Println("Error:", err)
		}
	}
}
