package main

import (
	"bufio"
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
const BlockTypeAlert = "Alert"
const BlockTypeLink = "Link"
const BlockTypeButton = "Button"

type sendBlock struct {
	blockType string
	text      string
	items     []string
	url       string
	style     string
	alt       string
	width     int
	ghost     bool
}

type sendOptions struct {
	apiKey  string
	to      string
	subject string
	blocks  []sendBlock
	dump    bool
}

func readStdin() (bool, []byte, error) {
	stat, _ := os.Stdin.Stat()
	var buf []byte
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			buf = append(buf, scanner.Bytes()...)
			buf = append(buf, '\n')
		}
		if err := scanner.Err(); err != nil {
			return true, buf, err
		}
		return true, buf, nil
	}
	return false, buf, nil
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
	optionToBlockType["--alert"] = BlockTypeAlert
	optionToBlockType["--link"] = BlockTypeLink
	optionToBlockType["--button"] = BlockTypeButton

	for i := 0; i < len(args); i += 2 {
		arg := args[i]

		if arg == "--dump" {
			options.dump = true
			i -= 1
			continue
		}

		if i+1 == len(args) {
			return nil, errors.New("missing value for " + arg)
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
		case "--alert":
			blockType := optionToBlockType[arg]
			alertOptions := make([]string, 0)
			for k := i + 2; k < len(args); k += 1 {
				if strings.HasPrefix(args[k], "--") || !strings.Contains(args[k], ":") {
					break
				}
				alertOptions = append(alertOptions, args[k])
			}
			alertBlock := sendBlock{
				blockType: blockType,
				text:      value,
			}
			for _, arg := range alertOptions {
				if strings.HasPrefix(arg, "style:") {
					style := arg[6:]
					if style != "success" && style != "warning" && style != "danger" && style != "info" {
						return nil, errors.New("invalid style: '" + style + "' (should be one of: success, warning, danger, info)")
					}
					alertBlock.style = style
				} else {
					return nil, errors.New("unknown option: '" + arg + "'")
				}
			}
			blocks = append(blocks, alertBlock)
			i += len(alertOptions)
		case "--button":
			blockType := optionToBlockType[arg]
			buttonOptions := make([]string, 0)
			url := value
			if i+2 == len(args) {
				return nil, errors.New("missing button text")
			}
			text := args[i+2]
			for k := i + 3; k < len(args); k += 1 {
				if strings.HasPrefix(args[k], "--") || !strings.Contains(args[k], ":") {
					break
				}
				buttonOptions = append(buttonOptions, args[k])
			}
			buttonBlock := sendBlock{
				blockType: blockType,
				url:       url,
				text:      text,
			}
			for _, arg := range buttonOptions {
				if strings.HasPrefix(arg, "style:") {
					style := arg[6:]
					if style != "success" && style != "warning" && style != "danger" && style != "info" {
						return nil, errors.New("invalid style: '" + style + "' (should be one of: success, warning, danger, info)")
					}
					buttonBlock.style = style
				} else if strings.HasPrefix(arg, "ghost:") {
					value := arg[6:]
					buttonBlock.ghost = value != "" && value != "false"
				} else {
					return nil, errors.New("unknown option: '" + arg + "'")
				}
			}
			blocks = append(blocks, buttonBlock)
			i += 1 // text
			i += len(buttonOptions)
		case "--link":
			blockType := optionToBlockType[arg]
			text := value
			if i < len(args)-2 && !strings.HasPrefix(args[i+2], "-") {
				text = args[i+2]
				i += 1
			}
			blocks = append(blocks, sendBlock{
				blockType: blockType,
				url:       value,
				text:      text,
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
	BlockType string   `json:"type"`
	Text      string   `json:"text,omitempty"`
	Items     []string `json:"items,omitempty"`
	Url       string   `json:"url,omitempty"`
	Alt       string   `json:"alt,omitempty"`
	Width     int      `json:"width,omitempty"`
	Style     string   `json:"style,omitempty"`
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
		case BlockTypeAlert:
			blockPayload.Text = block.text
			blockPayload.Style = block.style
		case BlockTypeLink:
			blockPayload.Url = block.url
			blockPayload.Text = block.text
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

func envOrDevault(envName string, defaultValue string, preferDefault bool) string {
	if preferDefault && defaultValue != "" {
		return defaultValue
	}
	fromEnv := os.Getenv(envName)
	if fromEnv == "" {
		return defaultValue
	}
	return fromEnv
}

type runSendType func(args []string) error

func runSend(args []string) error {
	didReadStdin, stdinContent, stdinError := readStdin()
	if stdinError != nil {
		return stdinError
	}

	options, err1 := parseSendArgs(args)
	if err1 != nil {
		return err1
	}

	options.apiKey = envOrDevault("MENDSAIL_API_KEY", options.apiKey, true)
	options.to = envOrDevault("MENDSAIL_TO", options.to, true)
	options.subject = envOrDevault("MENDSAIL_SUBJECT", options.subject, true)

	err2 := validateSendOptions(*options)
	if err2 != nil {
		return err2
	}

	if didReadStdin {
		options.blocks = append(options.blocks, sendBlock{
			blockType: "CodeBlock",
			text:      string(stdinContent),
		})
	}

	payload, err3 := sendOptionsToJsonPayload(*options)
	if err3 != nil {
		return err3
	}

	if options.dump {
		fmt.Println("Begin JSON payload")
		fmt.Println(string(payload))
		fmt.Println("End JSON payload")
		return errors.New("--dump was specified, aborting after printing JSON")
	}

	apiBaseUrl := envOrDevault("MENDSAIL_BASE_URL", "https://api.mendsail.com/v1", false)
	apiBaseUrl = strings.Trim(apiBaseUrl, "/")
	apiEndpoint := apiBaseUrl + "/emails"

	err4 := postJson(apiEndpoint, options.apiKey, payload)
	if err4 != nil {
		return err4
	}

	fmt.Println("Email sent successfully.")

	return nil
}
