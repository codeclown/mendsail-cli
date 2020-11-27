package main

import (
	"encoding/json"
	"errors"
)

const BlockTypeHeading = "Heading"
const BlockTypeParagraph = "Paragraph"
const BlockTypeCodeBlock = "CodeBlock"

type sendBlock struct {
	blockType string
	text      string
}

type sendOptions struct {
	apiKey  string
	to      string
	subject string
	blocks  []sendBlock
}

func parseSendArgs(args []string) (*sendOptions, error) {
	options := sendOptions{}
	blocks := make([]sendBlock, 0)

	optionToBlockType := make(map[string]string)
	optionToBlockType["--add-heading"] = BlockTypeHeading
	optionToBlockType["--add-paragraph"] = BlockTypeParagraph
	optionToBlockType["--add-code-block"] = BlockTypeCodeBlock

	for i := 0; i < len(args); i += 2 {
		arg := args[i]

		if i+1 == len(args) {
			return nil, errors.New("Missing value for " + arg)
		}

		value := args[i+1]

		switch arg {
		case "--api-key":
			options.apiKey = value
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

func validateSendOptions(options sendOptions) error {
	if options.apiKey == "" {
		return errors.New("missing option: --api-key")
	}
	if options.to == "" {
		return errors.New("missing option: --to")
	}
	if options.subject == "" {
		return errors.New("missing option: --subject")
	}
	return nil
}

type BlockPayload struct {
	BlockType string `json:"type"`
	Text      string `json:"text,omitempty"`
}

type AsdPayload struct {
	To      string         `json:"to"`
	Subject string         `json:"subject"`
	Blocks  []BlockPayload `json:"blocks"`
}

func sendOptionsToJsonPayload(options sendOptions) ([]byte, error) {
	blocks := []BlockPayload{}
	for _, block := range options.blocks {
		blockPayload := BlockPayload{
			BlockType: block.blockType,
		}
		switch block.blockType {
		case BlockTypeHeading:
			blockPayload.Text = block.text
		case BlockTypeParagraph:
			blockPayload.Text = block.text
		case BlockTypeCodeBlock:
			blockPayload.Text = block.text
		}
		blocks = append(blocks, blockPayload)
	}
	payload := AsdPayload{
		To:      options.to,
		Subject: options.subject,
		Blocks:  blocks,
	}
	return json.Marshal(payload)
}

type runSendType func(args []string) error

func runSend(args []string) error {
	options, err1 := parseSendArgs(args)
	if err1 != nil {
		return err1
	}

	err2 := validateSendOptions(*options)
	if err2 != nil {
		return err2
	}

	payload, err3 := sendOptionsToJsonPayload(*options)
	if err3 != nil {
		return err3
	}

	err4 := postJson("https://reqres.in/api/users", options.apiKey, payload)
	if err4 != nil {
		return err4
	}

	return nil
}
