package main

import (
	"fmt"
	"os"

	"github.com/farion1231/ccmux/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
