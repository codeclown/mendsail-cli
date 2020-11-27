package main

import (
	"errors"
	"fmt"
	"os"
)

func runMain(args []string, runSendFn runSendType) error {
	usage := "expected a subcommand (e.g. 'mendsail send ...')"

	if len(args) < 1 {
		return errors.New(usage)
	}

	switch args[0] {
	case "send":
		return runSendFn(args[1:])
	default:
		return errors.New(usage)
	}
}

func main() {
	err := runMain(os.Args[1:], runSend)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
