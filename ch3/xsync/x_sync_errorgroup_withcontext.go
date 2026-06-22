//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	var urls = []string{
		"https://go.dev/",
		"https://www.google.com/",
		"https://jianghushinian.cn/",
	}

	g, ctx := errgroup.WithContext(context.Background())

	for i, url := range urls {
		g.Go(func() error {
			ctx, cancel := context.WithTimeout(ctx, time.Duration(i+1)*time.Second)
			defer cancel()
			req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Printf("Go: Failed to fetch %v\n", err)
				return err
			}
			defer resp.Body.Close()
			fmt.Printf("Fetch url %s status %s\n", url, resp.Status)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("Wait: Failed to fetch %v\n", err)
	}
}
