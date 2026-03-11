package main

import (
	"os"

	"github.com/y0s3ph/k8scope/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
