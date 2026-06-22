//go:build ignore
// +build ignore

package main

import (
	"errors"
	"fmt"
	"io"
)

type Error string

func (e Error) Error() string { return string(e) }

const ErrMyEOF = Error("MyEOF")

func main() {
	const ErrMyEOF = Error("MyEOF")
	const ErrNewMyEOF = Error("MyEOF")
	fmt.Println(ErrMyEOF == ErrNewMyEOF) // true

	myEOF := errors.New("EOF")
	fmt.Println(io.EOF == myEOF) // false
}
