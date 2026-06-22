//go:build ignore
// +build ignore

package main

import "fmt"

func g(i int) (number int, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("f returns err: %v", r)
			}
		}
	}()

	number, err = f(i)
	return number, err
}

func f(i int) (number int, err error) {
	if i%2 == 0 {
		// panic("i is even")
		return 0, fmt.Errorf("%d is even", i)
	}
	return i * 2, nil
}

func main() {
	number, err := g(4)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	} else {
		fmt.Printf("number: %v\n", number)
	}
}
