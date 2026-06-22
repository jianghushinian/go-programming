//go:build ignore
// +build ignore

package main

func isTimeout(err error) bool {
	type timeout interface {
		Timeout() bool
	}
	te, ok := err.(timeout)
	return ok && te.Timeout()
}
