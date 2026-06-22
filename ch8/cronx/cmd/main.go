package main

import (
	"os"

	"cronx/cmd/app"
)

func main() {
	cmd := app.NewCronxCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
