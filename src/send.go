package main

import (
	"errors"
)

type sendBlock struct {
	blockType string
	text      string
}

type sendOptions struct {
	to      string
	subject string
	blocks  []sendBlock
}

func parseSendArgs(args []string) (*sendOptions, error) {
	options := sendOptions{}
	blocks := make([]sendBlock, 0)

	optionToBlockType := make(map[string]string)
	optionToBlockType["--add-heading"] = "Heading"
	optionToBlockType["--add-paragraph"] = "Paragraph"
	optionToBlockType["--add-code-block"] = "CodeBlock"

	for i := 0; i < len(args); i += 2 {
		arg := args[i]

		if i+1 == len(args) {
			return nil, errors.New("Missing value for " + arg)
		}

		value := args[i+1]

		switch arg {
		case "--to":
			options.to = value
		case "--subject":
			options.subject = value
		case "--add-heading", "--add-paragraph", "--add-code-block":
			blockType := optionToBlockType[arg]
			blocks = append(blocks, sendBlock{
				blockType: blockType,
				text:      value,
			})
		default:
			return nil, errors.New("Unrecognized option: " + arg)
		}
	}

	options.blocks = blocks
	return &options, nil
}

type runSendType func(args []string) error

func runSend(args []string) error {
	// err, options := parseSendArgs(args)
	// fmt.Println(options)
	return nil
}
