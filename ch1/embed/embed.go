package main

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed hello.txt
var content string

//go:embed hello.txt
var contentBytes []byte

//go:embed hello.txt
var fileFS embed.FS
var data, _ = fileFS.ReadFile("hello.txt")

//go:embed file/hello1.txt
//go:embed file/hello2.txt
var helloFS embed.FS

/*
//go:embed file/hello1.txt file/hello2.txt
var helloFS embed.FS
*/

func main() {
	fmt.Printf("content: %s\n", content)
	fmt.Printf("contentBytes: %s\n", contentBytes)
	fmt.Printf("data: %s\n", data)

	subFS, _ := fs.Sub(helloFS, "file")
	f, _ := subFS.Open("hello2.txt")

	hello2 := make([]byte, 6)
	_, _ = f.Read(hello2)
	fmt.Printf("hello2: %s\n", hello2)
}
