package main

import (
	"fmt"
	"os"
)

func commandNotFound() {
	fmt.Println("expected a subcommand (e.g. 'mendsail send ...')")
	os.Exit(1)
}

func main() {

	if len(os.Args) < 2 {
		commandNotFound()
	}

	switch os.Args[1] {
	case "send":
		options, err := parseSendArgs(os.Args[2:])
		if err != nil {
			fmt.Println(err)
		}
		runSend(*options)
	default:
		commandNotFound()
	}
}
