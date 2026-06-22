package main

import "log"

type MyError struct {
	msg string
}

func (e *MyError) Error() string {
	return e.msg
}

type ZeroDivisionError struct{}

func (e *ZeroDivisionError) Error() string {
	return "division by zero"
}

// func main() {
// 	var err error
//
// 	if err != nil {
// 		switch err.(type) {
// 		case *MyError:
// 			// ...
// 		case *ZeroDivisionError:
// 			// ...
// 		default:
// 			// ...
// 		}
// 	}
// }

func a() error {
	// ...
	return nil
}

func b() error {
	// ...
	return nil
}

func c() error {
	// ...
	return nil
}

func main() {
	err := a()
	if err != nil {
		log.Fatal(err)
	}
	err = b()
	if err != nil {
		log.Fatal(err)
	}
	err = c()
	if err != nil {
		log.Fatal(err)
	}
}
