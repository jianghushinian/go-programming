package main

import (
	"fmt"

	"linkname/foo"
)

func main() {
	fmt.Println("foo.Add(1, 2):", foo.Add(1, 2))
	fmt.Println("foo.Div(2, 1):", foo.Div(2, 1))
	fmt.Println(`foo.Hello("江湖十年"):`, foo.Hello("江湖十年"))
}
