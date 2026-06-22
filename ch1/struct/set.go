//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"strings"
)

type Set map[string]struct{}

func (s Set) Add(element string) {
	s[element] = struct{}{}
}

func (s Set) Remove(element string) {
	delete(s, element)
}

func (s Set) Contains(element string) bool {
	_, exists := s[element]
	return exists
}

func (s Set) Size() int {
	return len(s)
}

// String implements fmt.Stringer
func (s Set) String() string {
	format := "("
	for element := range s {
		format += element + " "
	}
	format = strings.TrimRight(format, " ") + ")"
	return format
}

func main() {
	s := make(Set)

	s.Add("one")
	s.Add("two")
	s.Add("three")

	fmt.Printf("set: %s\n", s)
	fmt.Printf("set size: %d\n", s.Size())
	fmt.Printf("set contains 'one': %t\n", s.Contains("one"))
	fmt.Printf("set contains 'onex': %t\n", s.Contains("onex"))

	s.Remove("one")

	fmt.Printf("set: %s\n", s)
	fmt.Printf("set size: %d\n", s.Size())
}
