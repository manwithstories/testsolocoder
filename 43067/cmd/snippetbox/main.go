package main

import (
	"fmt"
	"os"

	"snippetbox/internal/logger"
)

func main() {
	if err := logger.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to initialize logger: %v\n", err)
	}

	if err := Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
