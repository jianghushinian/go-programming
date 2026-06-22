//go:build ignore
// +build ignore

package main

import (
	"log"
	"os"
)

func main() {
	log.Print("Print")
	log.Printf("Printf: %s", "print")
	log.Println("Println")

	logger := log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lshortfile)
	logger.Println("custom logger")
}
