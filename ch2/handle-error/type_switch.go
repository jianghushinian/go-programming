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
		// 使用 switch type 检查错误类型
		switch e := err.(type) {
		case *os.PathError:
			fmt.Printf("Failed to %s file: %s\n", e.Op, e.Path)
			fmt.Println("Error message:", e.Err)
		default:
			fmt.Println("Error:", err)
		}
	}

	if err != nil {
		switch err.(type) {
		case *os.PathError, *os.LinkError:
			// do something
		default:
			// do something
		}
	}
}
