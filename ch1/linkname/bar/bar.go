package bar

import _ "unsafe"

func add(a, b int) int {
	return a + b
}

//go:linkname div linkname/foo.Div
func div(a, b int) int {
	return a / b
}

//go:linkname hello
func hello(name string) string {
	return "Hello " + name + "!"
}
