//go:build ignore
// +build ignore

package main

import (
	"io"
	"os"
)

type file struct {
	name string
}

func (f *file) Close() {}

func CopyFile1(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}

	dst, err := os.Create(dstName)
	if err != nil {
		return
	}

	written, err = io.Copy(dst, src)
	dst.Close()
	src.Close()
	return
}

func CopyFile2(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

func processFile1() {
	f := file{name: "f1"}
	defer f.Close()

	f = file{name: "f2"}
	defer f.Close()

	return
}

func processFile2() {
	f := file{name: "f1"}
	defer func(f file) {
		f.Close()
	}(f)

	f = file{name: "f2"}
	defer func(f file) {
		f.Close()
	}(f)

	return
}

func main() {}
