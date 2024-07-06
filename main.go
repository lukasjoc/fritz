package main

import (
	"fmt"
	"os"

	"github.com/lukasjoc/fritz/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
