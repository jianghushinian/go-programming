//go:build ignore
// +build ignore

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	// 典型哨兵错误处理方式
	var err error
	if err != nil {
		if err == bufio.ErrBufferFull {
			// handle ErrBufferFull
		}
		// handle err
	}

	f, err := os.Open("example.txt")
	if err != nil {
		return
	}

	// switch-case 风格哨兵错误处理
	b := bufio.NewReader(f)
	data, err := b.Peek(10)
	if err != nil {
		switch err {
		case bufio.ErrNegativeCount:
			// do something
			return
		case bufio.ErrBufferFull:
			// do something
			return
		default:
			// do something
			return
		}
	}
	fmt.Println(string(data))

	oldEOF := io.EOF
	io.EOF = errors.New("MyEOF")
	fmt.Println(oldEOF == io.EOF) // false
}
