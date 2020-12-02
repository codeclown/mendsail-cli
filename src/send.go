package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const BlockTypeHeading = "Heading"
const BlockTypeParagraph = "Paragraph"
const BlockTypeList = "List"
const BlockTypeImage = "Image"
const BlockTypeCodeBlock = "CodeBlock"

type sendBlock struct {
	blockType string
	text      string
	items     []string
	url       string
	alt       string
	width     int
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
	optionToBlockType["--heading"] = BlockTypeHeading
	optionToBlockType["--paragraph"] = BlockTypeParagraph
	optionToBlockType["--list"] = BlockTypeList
	optionToBlockType["--image"] = BlockTypeImage
	optionToBlockType["--code-block"] = BlockTypeCodeBlock

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
		case "--heading", "--paragraph", "--code-block":
			blockType := optionToBlockType[arg]
			blocks = append(blocks, sendBlock{
				blockType: blockType,
				text:      value,
			})
		case "--list":
			blockType := optionToBlockType[arg]
			listItems := make([]string, 0)
			for k := i + 1; k < len(args); k += 1 {
				if strings.HasPrefix(args[k], "--") {
					break
				}
				listItems = append(listItems, args[k])
			}
			blocks = append(blocks, sendBlock{
				blockType: blockType,
				items:     listItems,
			})
			i += len(listItems) - 1
		case "--image":
			blockType := optionToBlockType[arg]
			imageOptions := make([]string, 0)
			for k := i + 2; k < len(args); k += 1 {
				if strings.HasPrefix(args[k], "--") || !strings.Contains(args[k], ":") {
					break
				}
				imageOptions = append(imageOptions, args[k])
			}
			imageBlock := sendBlock{
				blockType: blockType,
				url:       value,
			}
			for _, arg := range imageOptions {
				if strings.HasPrefix(arg, "alt:") {
					imageBlock.alt = arg[4:]
				} else if strings.HasPrefix(arg, "width:") {
					width, conversionErr := strconv.Atoi(arg[6:])
					if conversionErr != nil {
						return nil, errors.New("could not parse width as an integer")
					}
					imageBlock.width = width
				} else {
					return nil, errors.New("unknown option: '" + arg + "'")
				}
			}
			blocks = append(blocks, imageBlock)
			i += len(imageOptions)
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
	BlockType string   `json:"type"`
	Text      string   `json:"text,omitempty"`
	Items     []string `json:"items,omitempty"`
	Url       string   `json:"url,omitempty"`
	Alt       string   `json:"alt,omitempty"`
	Width     int      `json:"width,omitempty"`
}

type FullPayload struct {
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
		case BlockTypeList:
			blockPayload.Items = block.items
		case BlockTypeImage:
			blockPayload.Url = block.url
			blockPayload.Alt = block.alt
			blockPayload.Width = block.width
		}
		blocks = append(blocks, blockPayload)
	}
	payload := FullPayload{
		To:      options.to,
		Subject: options.subject,
		Blocks:  blocks,
	}
	return json.Marshal(payload)
}

func getApiEndpoint() string {
	apiBaseUrl := strings.Trim(os.Getenv("MENDSAIL_BASE_URL"), "/")
	if apiBaseUrl == "" {
		apiBaseUrl = "https://api.mendsail.com/v1"
	}
	return apiBaseUrl + "/emails"
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

	apiEndpoint := getApiEndpoint()

	err4 := postJson(apiEndpoint, options.apiKey, payload)
	if err4 != nil {
		return err4
	}

	fmt.Println("Email sent successfully.")

	return nil
}
