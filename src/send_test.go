package main

import (
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
		t.Errorf("len(actual.blocks): expected=%d actual=%d", len(actual.blocks), len(expected.blocks))
	}
	for i, actualBlock := range actual.blocks {
		expectedBlock := expected.blocks[i]
		if actualBlock.blockType != expectedBlock.blockType {
			t.Errorf("sendOptions.blocks[%d].blockType: expected=%s actual=%s",
				i, actualBlock.blockType, expectedBlock.blockType)
		}
		if actualBlock.text != expectedBlock.text {
			t.Errorf("sendOptions.blocks[%d].text: expected=%s actual=%s",
				i, actualBlock.text, expectedBlock.text)
		}
	}
}

func Test_ParseSendArgs_Empty(t *testing.T) {
	var args []string
	expected := sendOptions{
		to:      "",
		subject: "",
		blocks:  []sendBlock{},
	}
	actual, err := parseSendArgs(args)
	exceptOptions(t, expected, actual, err)
}

func Test_ParseSendArgs_Basic(t *testing.T) {
	args := []string{
		"--to", "foobar@example.com",
		"--subject", "example 123",
	}
	expected := sendOptions{
		to:      "foobar@example.com",
		subject: "example 123",
		blocks:  []sendBlock{},
	}
	actual, err := parseSendArgs(args)
	exceptOptions(t, expected, actual, err)
}

func Test_ParseSendArgs_UnknownOption(t *testing.T) {
	args := []string{
		"--foobar", "foobar@example.com",
		"--subject", "example 123",
	}
	_, err := parseSendArgs(args)
	expectError(t, "Unrecognized option: --foobar", err)
}

func Test_ParseSendArgs_MissingValue(t *testing.T) {
	args := []string{
		"--to", "foobar@example.com",
		"--subject",
	}
	_, err := parseSendArgs(args)
	expectError(t, "Missing value for --subject", err)
}

func Test_ParseSendArgs_BlockTypes(t *testing.T) {
	args := []string{
		"--to", "foobar@example.com",
		"--subject", "example 123",
		"--add-heading", "Data processing failed",
		"--add-paragraph", "Log output:",
		"--add-code-block", "foobar",
	}
	expected := sendOptions{
		to:      "foobar@example.com",
		subject: "example 123",
		blocks: []sendBlock{
			sendBlock{blockType: "Heading", text: "Data processing failed"},
			sendBlock{blockType: "Paragraph", text: "Log output:"},
			sendBlock{blockType: "CodeBlock", text: "foobar"},
		},
	}
	actual, err := parseSendArgs(args)
	exceptOptions(t, expected, actual, err)
}

func Test_ParseSendArgs_BlockOrder(t *testing.T) {
	args := []string{
		"--to", "foobar@example.com",
		"--subject", "example 123",
		"--add-heading", "heading 1",
		"--add-paragraph", "paragraph 1",
		"--add-paragraph", "paragraph 2",
		"--add-code-block", "code block 1",
		"--add-code-block", "code block 2",
		"--add-heading", "heading 2",
		"--add-paragraph", "paragraph 3",
	}
	expected := sendOptions{
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
	exceptOptions(t, expected, actual, err)
}