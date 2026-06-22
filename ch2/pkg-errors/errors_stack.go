//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

func Foo() error {
	return io.EOF
}

func Bar() error {
	err := Foo()
	if err != nil {
		return errors.Wrap(err, "bar")
	}
	return nil
}

func main() {
	err := Bar()
	if err != nil {
		fmt.Printf("err: %+v\n", err)
	}
}
