package cmd

import (
	"fmt"
	"os"

	"github.com/lukasjoc/fritz/internal"
)

// TODO: dont like the arch of this cmd. Should use sth else.

func usage(program string, cmd string) {
	if len(cmd) > 0 {
		fmt.Fprintf(os.Stderr, "unknown command \"%s\" for \"%s\"\n\n", cmd, program)
	}
	fmt.Fprintf(os.Stderr, "Usage: %s <command>\n", program)
	fmt.Fprintf(os.Stderr, "%-10s %s\n", "info", "Print info about the box and its configuration")
	fmt.Fprintf(os.Stderr, "%-10s %s\n", "reconnect", "Quickly disconnect and reconnect again")
	fmt.Fprintf(os.Stderr, "%-10s %s\n", "reboot", "Quickly reboot the box")
	os.Exit(0)
}

func Run() error {
	program := os.Args[0]
	if len(os.Args) < 2 {
		usage(program, "")
	}
	fritz, err := internal.NewFritz()
	if err != nil {
		panic(err)
	}
	cmd := os.Args[1]
	switch cmd {
	case "reboot":
		if err := fritz.Reboot(); err != nil {
			return err
		}
	case "reconnect":
		if err := fritz.Reconnect(); err != nil {
			return err
		}
	case "info":
		if err := fritz.Info(); err != nil {
			return err
		}
	default:
		usage(program, cmd)
	}
	return nil
}
