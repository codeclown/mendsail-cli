package main

import (
	"reflect"
	"testing"
)

func exceptOptions(t *testing.T, expected sendOptions, actual *sendOptions, err error) {
	if err != nil {
		t.Errorf("err: expected=nil actual=%s", err)
		t.FailNow()
	}
	if actual.to != expected.to {
		t.Errorf("sendOptions.to: expected=%s actual=%s", expected.to, actual.to)
	}
	if actual.subject != expected.subject {
		t.Errorf("sendOptions.subject: expected=%s actual=%s", expected.subject, actual.subject)
	}
	if len(actual.blocks) != len(expected.blocks) {
		t.Errorf("len(actual.blocks): expected=%d actual=%d", len(expected.blocks), len(actual.blocks))
	}
	for i, actualBlock := range actual.blocks {
		expectedBlock := expected.blocks[i]
		if actualBlock.blockType != expectedBlock.blockType {
			t.Errorf("sendOptions.blocks[%d].blockType: expected=%s actual=%s",
				i, actualBlock.blockType, expectedBlock.blockType)
		}
		if actualBlock.text != expectedBlock.text {
			t.Errorf("sendOptions.blocks[%d].text: expected=%s actual=%s",
				i, expectedBlock.text, actualBlock.text)
		}
		if actualBlock.url != expectedBlock.url {
			t.Errorf("sendOptions.blocks[%d].url: expected=%s actual=%s",
				i, expectedBlock.url, actualBlock.url)
		}
		if !reflect.DeepEqual(actualBlock.items, expectedBlock.items) {
			t.Errorf("sendOptions.blocks[%d].items: expected=%s actual=%s",
				i, expectedBlock.items, actualBlock.items)
		}
		if actualBlock.alt != expectedBlock.alt {
			t.Errorf("sendOptions.blocks[%d].alt: expected=%s actual=%s",
				i, expectedBlock.alt, actualBlock.alt)
		}
		if actualBlock.width != expectedBlock.width {
			t.Errorf("sendOptions.blocks[%d].width: expected=%d actual=%d",
				i, expectedBlock.width, actualBlock.width)
		}
		if actualBlock.ghost != expectedBlock.ghost {
			t.Errorf("sendOptions.blocks[%d].ghost: expected=%t actual=%t",
				i, expectedBlock.ghost, actualBlock.ghost)
		}
	}
}

func Test_parseSendArgs_Empty(t *testing.T) {
	var args []string
	expected := sendOptions{
		apiKey:  "",
		to:      "",
		subject: "",
		blocks:  []sendBlock{},
	}
	actual, err := parseSendArgs(args)
	expectNoError(t, err)
	exceptOptions(t, expected, actual, err)
}

func Test_parseSendArgs_Basic(t *testing.T) {
	args := []string{
		"--api-key", "foobar-123",
		"--to", "foobar@example.com",
		"--subject", "example 123",
	}
	expected := sendOptions{
		apiKey:  "foobar-123",
		to:      "foobar@example.com",
		subject: "example 123",
		blocks:  []sendBlock{},
	}
	actual, err := parseSendArgs(args)
	expectNoError(t, err)
	exceptOptions(t, expected, actual, err)
}

func Test_parseSendArgs_UnknownOption(t *testing.T) {
	args := []string{
		"--api-key", "foobar-123",
		"--foobar", "foobar@example.com",
		"--subject", "example 123",
	}
	_, err := parseSendArgs(args)
	expectError(t, "Unrecognized option: --foobar", err)
}

func Test_parseSendArgs_MissingValue(t *testing.T) {
	args := []string{
		"--api-key", "foobar-123",
		"--to", "foobar@example.com",
		"--subject",
	}
	_, err := parseSendArgs(args)
	expectError(t, "missing value for --subject", err)
}

func Test_parseSendArgs_BlockTypes(t *testing.T) {
	args := []string{
		"--api-key", "foobar-123",
		"--to", "foobar@example.com",
		"--subject", "example 123",
		"--heading", "Data processing failed",
		"--paragraph", "Log output:",
		"--list", "List item 1",
		"--image", "https://example.com/image.png",
		"--code-block", "code 1",
		"--alert", "alert 1",
		"--link", "https://example.com",
		"--button", "https://example.com", "Button text",
	}
	expected := sendOptions{
		apiKey:  "foobar-123",
		to:      "foobar@example.com",
		subject: "example 123",
		blocks: []sendBlock{
			sendBlock{blockType: "Heading", text: "Data processing failed"},
			sendBlock{blockType: "Paragraph", text: "Log output:"},
			sendBlock{blockType: "List", items: []string{"List item 1"}},
			sendBlock{blockType: "Image", url: "https://example.com/image.png"},
			sendBlock{blockType: "CodeBlock", text: "code 1"},
			sendBlock{blockType: "Alert", text: "alert 1", style: "info"},
			sendBlock{blockType: "Link", url: "https://example.com", text: "https://example.com"},
			sendBlock{blockType: "Button", url: "https://example.com", text: "Button text"},
		},
	}
	actual, err := parseSendArgs(args)
	expectNoError(t, err)
	exceptOptions(t, expected, actual, err)
}

func Test_parseSendArgs_BlockOrder(t *testing.T) {
	args := []string{
		"--api-key", "foobar-123",
		"--to", "foobar@example.com",
		"--subject", "example 123",
		"--heading", "heading 1",
		"--paragraph", "paragraph 1",
		"--paragraph", "paragraph 2",
		"--code-block", "code block 1",
		"--code-block", "code block 2",
		"--heading", "heading 2",
		"--paragraph", "paragraph 3",
	}
	expected := sendOptions{
		apiKey:  "foobar-123",
		to:      "foobar@example.com",
		subject: "example 123",
		blocks: []sendBlock{
			sendBlock{blockType: "Heading", text: "heading 1"},
			sendBlock{blockType: "Paragraph", text: "paragraph 1"},
			sendBlock{blockType: "Paragraph", text: "paragraph 2"},
			sendBlock{blockType: "CodeBlock", text: "code block 1"},
			sendBlock{blockType: "CodeBlock", text: "code block 2"},
			sendBlock{blockType: "Heading", text: "heading 2"},
			sendBlock{blockType: "Paragraph", text: "paragraph 3"},
		},
	}
	actual, err := parseSendArgs(args)
	expectNoError(t, err)
	exceptOptions(t, expected, actual, err)
}

func Test_parseSendArgs_ListMultipleItems(t *testing.T) {
	args := []string{
		"--api-key", "foobar-123",
		"--to", "foobar@example.com",
		"--subject", "example 123",
		"--list", "List item 1", "List item 2", "List item 3",
	}
	expected := sendOptions{
		apiKey:  "foobar-123",
		to:      "foobar@example.com",
		subject: "example 123",
		blocks: []sendBlock{
			sendBlock{blockType: "List", items: []string{"List item 1", "List item 2", "List item 3"}},
		},
	}
	actual, err := parseSendArgs(args)
	expectNoError(t, err)
	exceptOptions(t, expected, actual, err)
}

func Test_parseSendArgs_ImageOptions(t *testing.T) {
	args := []string{
		"--api-key", "foobar-123",
		"--to", "foobar@example.com",
		"--subject", "example 123",
		"--image", "https://example.com/image.png", "alt:Alt text", "width:123",
	}
	expected := sendOptions{
		apiKey:  "foobar-123",
		to:      "foobar@example.com",
		subject: "example 123",
		blocks: []sendBlock{
			sendBlock{blockType: "Image", url: "https://example.com/image.png", alt: "Alt text", width: 123},
		},
	}
	actual, err := parseSendArgs(args)
	expectNoError(t, err)
	exceptOptions(t, expected, actual, err)
}

func Test_parseSendArgs_UnknownImageWidthType(t *testing.T) {
	args := []string{
		"--api-key", "foobar-123",
		"--to", "foobar@example.com",
		"--subject", "example 123",
		"--image", "https://example.com/image.png", "width:asd",
	}
	_, err := parseSendArgs(args)
	expectError(t, "could not parse width as an integer", err)
}

func Test_parseSendArgs_UnknownImageOption(t *testing.T) {
	args := []string{
		"--api-key", "foobar-123",
		"--to", "foobar@example.com",
		"--subject", "example 123",
		"--image", "https://example.com/image.png", "foobar:Alt text",
	}
	_, err := parseSendArgs(args)
	expectError(t, "unknown option: 'foobar:Alt text'", err)
}

func Test_parseSendArgs_AlertOptions(t *testing.T) {
	args := []string{
		"--api-key", "foobar-123",
		"--to", "foobar@example.com",
		"--subject", "example 123",
		"--alert", "lorem ipsum", "style:danger",
	}
	expected := sendOptions{
		apiKey:  "foobar-123",
		to:      "foobar@example.com",
		subject: "example 123",
		blocks: []sendBlock{
			sendBlock{blockType: "Alert", text: "lorem ipsum", style: "danger"},
		},
	}
	actual, err := parseSendArgs(args)
	expectNoError(t, err)
	exceptOptions(t, expected, actual, err)
}

func Test_parseSendArgs_UnknownAlertOption(t *testing.T) {
	args := []string{
		"--api-key", "foobar-123",
		"--to", "foobar@example.com",
		"--subject", "example 123",
		"--alert", "lorem ipsum", "foobar:danger",
	}
	_, err := parseSendArgs(args)
	expectError(t, "unknown option: 'foobar:danger'", err)
}

func Test_parseSendArgs_InvalidAlertStyle(t *testing.T) {
	args := []string{
		"--api-key", "foobar-123",
		"--to", "foobar@example.com",
		"--subject", "example 123",
		"--alert", "lorem ipsum", "style:foobar",
	}
	_, err := parseSendArgs(args)
	expectError(t, "invalid style: 'foobar' (should be one of: success, warning, danger, info)", err)
}

func Test_parseSendArgs_LinkText(t *testing.T) {
	args := []string{
		"--api-key", "foobar-123",
		"--to", "foobar@example.com",
		"--subject", "example 123",
		"--link", "https://example.com", "text foobar",
	}
	expected := sendOptions{
		apiKey:  "foobar-123",
		to:      "foobar@example.com",
		subject: "example 123",
		blocks: []sendBlock{
			sendBlock{blockType: "Link", url: "https://example.com", text: "text foobar"},
		},
	}
	actual, err := parseSendArgs(args)
	expectNoError(t, err)
	exceptOptions(t, expected, actual, err)
}

func Test_parseSendArgs_ButtonOptions(t *testing.T) {
	args := []string{
		"--api-key", "foobar-123",
		"--to", "foobar@example.com",
		"--subject", "example 123",
		"--button", "https://example.com", "lorem ipsum", "style:danger", "ghost:true",
	}
	expected := sendOptions{
		apiKey:  "foobar-123",
		to:      "foobar@example.com",
		subject: "example 123",
		blocks: []sendBlock{
			sendBlock{blockType: "Button", url: "https://example.com", text: "lorem ipsum", style: "danger", ghost: true},
		},
	}
	actual, err := parseSendArgs(args)
	expectNoError(t, err)
	exceptOptions(t, expected, actual, err)
}

func Test_parseSendArgs_UnknownMissingText(t *testing.T) {
	args := []string{
		"--api-key", "foobar-123",
		"--to", "foobar@example.com",
		"--subject", "example 123",
		"--button", "https://example.com",
	}
	_, err := parseSendArgs(args)
	expectError(t, "missing button text", err)
}

func Test_parseSendArgs_UnknownButtonOption(t *testing.T) {
	args := []string{
		"--api-key", "foobar-123",
		"--to", "foobar@example.com",
		"--subject", "example 123",
		"--button", "https://example.com", "lorem ipsum", "foobar:danger",
	}
	_, err := parseSendArgs(args)
	expectError(t, "unknown option: 'foobar:danger'", err)
}

func Test_parseSendArgs_InvalidButtonStyle(t *testing.T) {
	args := []string{
		"--api-key", "foobar-123",
		"--to", "foobar@example.com",
		"--subject", "example 123",
		"--button", "https://example.com", "lorem ipsum", "style:foobar",
	}
	_, err := parseSendArgs(args)
	expectError(t, "invalid style: 'foobar' (should be one of: success, warning, danger, info)", err)
}

func Test_validateSendOptions_Valid(t *testing.T) {
	options := sendOptions{
		apiKey:  "foobar-123",
		to:      "foobar@example.com",
		subject: "example 123",
	}
	err := validateSendOptions(options)
	expectNoError(t, err)
}

func Test_validateSendOptions_ApiKey(t *testing.T) {
	options := sendOptions{
		apiKey:  "",
		to:      "foobar@example.com",
		subject: "example 123",
	}
	err := validateSendOptions(options)
	expectError(t, "missing option: --api-key", err)
}

func Test_validateSendOptions_To(t *testing.T) {
	options := sendOptions{
		apiKey:  "foobar-123",
		to:      "",
		subject: "example 123",
	}
	err := validateSendOptions(options)
	expectError(t, "missing option: --to", err)
}

func Test_validateSendOptions_Subject(t *testing.T) {
	options := sendOptions{
		apiKey:  "foobar-123",
		to:      "foobar@example.com",
		subject: "",
	}
	err := validateSendOptions(options)
	expectError(t, "missing option: --subject", err)
}

func Test_sendOptionsToJsonPayload_works(t *testing.T) {
	options := sendOptions{
		apiKey:  "foobar-123",
		to:      "foobar@example.com",
		subject: "example 123",
		blocks: []sendBlock{
			sendBlock{blockType: "Heading", text: "heading 1"},
			sendBlock{blockType: "Paragraph", text: "paragraph 1"},
			sendBlock{blockType: "Paragraph", text: "paragraph 2"},
			sendBlock{blockType: "CodeBlock", text: "code block 1"},
			sendBlock{blockType: "CodeBlock", text: "code block 2"},
			sendBlock{blockType: "Heading", text: "heading 2"},
			sendBlock{blockType: "Paragraph", text: "paragraph 3"},
			sendBlock{blockType: "List", items: []string{"item 1", "item 2"}},
			sendBlock{blockType: "Image", url: "image.png"},
			sendBlock{blockType: "Image", url: "image.png", alt: "alt text", width: 123},
			sendBlock{blockType: "Alert", text: "alert 1", style: "info"},
			sendBlock{blockType: "Link", url: "https://example.com", text: "lorem ipsum"},
		},
	}
	expected := "{" +
		"\"to\":\"foobar@example.com\"," +
		"\"subject\":\"example 123\"," +
		"\"blocks\":[" +
		"{\"type\":\"Heading\",\"text\":\"heading 1\"}," +
		"{\"type\":\"Paragraph\",\"text\":\"paragraph 1\"}," +
		"{\"type\":\"Paragraph\",\"text\":\"paragraph 2\"}," +
		"{\"type\":\"CodeBlock\",\"text\":\"code block 1\"}," +
		"{\"type\":\"CodeBlock\",\"text\":\"code block 2\"}," +
		"{\"type\":\"Heading\",\"text\":\"heading 2\"}," +
		"{\"type\":\"Paragraph\",\"text\":\"paragraph 3\"}," +
		"{\"type\":\"List\",\"items\":[\"item 1\",\"item 2\"]}," +
		"{\"type\":\"Image\",\"url\":\"image.png\"}," +
		"{\"type\":\"Image\",\"url\":\"image.png\",\"alt\":\"alt text\",\"width\":123}," +
		"{\"type\":\"Alert\",\"text\":\"alert 1\",\"style\":\"info\"}," +
		"{\"type\":\"Link\",\"text\":\"lorem ipsum\",\"url\":\"https://example.com\"}" +
		"]" +
		"}"
	actual, err := sendOptionsToJsonPayload(options)
	expectNoError(t, err)
	exceptStringsEqual(t, expected, string(actual))
}
