//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log/slog"

	"github.com/pkg/errors"
)

/*
func Foo() error {
	return nil
}

func main() {
	err := Foo()
	fmt.Printf("call foo: %s\n", err)
}
*/

func Foo() error {
	return errors.New("foo error")
}

func Bar() error {
	err := Foo()
	return errors.WithMessage(err, "Bar")
}

func main() {
	err := Bar()
	if err != nil {
		slog.Error(fmt.Sprintf("%+v", err))
	}
}
