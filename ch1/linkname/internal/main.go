package main

import (
	"fmt"
	_ "unsafe"
)

//go:linkname TooLarge fmt.tooLarge
func TooLarge(x int) bool

func main() {
	fmt.Println("TooLarge(1e6 + 1):", TooLarge(1e6+1))
}
