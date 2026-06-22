//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func main() {
	var eg errgroup.Group
	var urls = []string{
		"https://go.dev/",
		"https://www.google.com/",
		"https://jianghushinian.cn/",
	}

	for _, url := range urls {
		eg.Go(func() error {
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("Go: Failed to fetch %v\n", err)
				return err
			}
			defer resp.Body.Close()
			fmt.Printf("Fetch url %s status %s\n", url, resp.Status)
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		fmt.Printf("Wait: Failed to fetch %v\n", err)
	}
}
