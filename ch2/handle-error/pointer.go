//go:build ignore
// +build ignore

package main

func Positive1(n int) (bool, bool) {
	if n == 0 {
		return false, false
	}
	return n > -1, true
}

func Positive2(n int) (bool, error) {
	if n == 0 {
		return false, errors.New("undefined")
	}
	return n > -1, nil
}

func Positive3(n int) *bool {
	if n == 0 {
		return nil
	}
	r := n > -1
	return &r
}
