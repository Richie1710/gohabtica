package main

import (
	"os"
)

func main() {
	if err := execute(os.Args[1:]); err != nil {
		// Simple stderr write without adding extra dependencies
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

