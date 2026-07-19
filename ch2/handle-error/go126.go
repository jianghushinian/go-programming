//go:build ignore
// +build ignore

package main

import (
	"errors"
	"fmt"
	"io"
)

type MyError struct {
	msg string
	err error
}

func (e *MyError) Error() string {
	return e.msg + ": " + e.err.Error()
}

func Foo() error {
	return &MyError{
		msg: "foo",
		err: io.EOF,
	}
}

func Bar() error {
	err := Foo()
	if err != nil {
		return fmt.Errorf("bar: %w", err)
	}
	return nil
}

func main() {
	if err := Bar(); err != nil {
		if myErr, ok := errors.AsType[*MyError](err); ok {
			fmt.Printf("EOF err: %s\n", myErr)
			return
		}
		fmt.Printf("err: %s\n", err)
	}
	return
}
