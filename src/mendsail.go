package main

import (
	"errors"
	"fmt"
	"os"
)

type showHelpType func() error

func showHelp() error {
	usage := "Usage: mendsail send <options> <blocks>\n" +
		"\n" +
		"Sending options:\n" +
		"  --api-key  <string>  API key for authentication\n" +
		"  --to       <string>  Recipient email address\n" +
		"  --subject  <string>  Subject line\n" +
		"\n" +
		"Blocks:\n" +
		"  --alert              <text> [style:success|warning|danger|info]\n" +
		"  --code-block         <text>\n" +
		"  --heading            <text>\n" +
		"  --image              <url> [alt:text] [width:number]\n" +
		"  --link               <url> [text]\n" +
		"  --list               <item1> <item2> ... <itemN>\n" +
		"  --paragraph          <text>\n" +
		"\n" +
		"Other options:\n" +
		"  --help               Show this help message\n" +
		"\n" +
		"Supported environment variables:\n" +
		"  MENDSAIL_API_KEY MENDSAIL_TO MENDSAIL_SUBJECT\n" +
		"\n" +
		"Links:\n" +
		"  - Documentation:     https://mendsail.com/docs\n" +
		"  - Source code:       https://github.com/codeclown/mendsail-cli\n" +
		""
	return errors.New(usage)
}

func runMain(args []string, showHelpFn showHelpType, runSendFn runSendType) error {
	if len(args) < 1 {
		return showHelpFn()
	}

	switch args[0] {
	case "send":
		return runSendFn(args[1:])
	default:
		return showHelpFn()
	}
}

func main() {
	err := runMain(os.Args[1:], showHelp, runSend)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
