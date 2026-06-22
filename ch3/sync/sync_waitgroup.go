//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var urls = []string{
		"https://go.dev/",
		"https://www.google.com/",
		"https://jianghushinian.cn/",
	}
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("Failed to fetch %s: %v\n", url, err)
				return
			}
			defer resp.Body.Close()
			fmt.Printf("Fetch url %s status %s\n", url, resp.Status)
		}(url)
	}
	wg.Wait()
}
