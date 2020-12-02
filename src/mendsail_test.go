package main

import (
	"errors"
	"reflect"
	"testing"
)

func dummyShowHelp() error {
	return nil
}

func dummyRunSend(args []string) error {
	return nil
}

func Test_runMain_NoArgs(t *testing.T) {
	var calledTimes int
	mockShowHelp := func() error {
		calledTimes += 1
		return errors.New("mocked help")
	}
	args := []string{}
	err := runMain(args, mockShowHelp, dummyRunSend)
	expectError(t, "mocked help", err)
	if calledTimes != 1 {
		t.Errorf("calledWith: expected=%d actual=%d", 1, calledTimes)
	}
}

func Test_runMain_UnknownCommand(t *testing.T) {
	var calledTimes int
	mockShowHelp := func() error {
		calledTimes += 1
		return errors.New("mocked help")
	}
	args := []string{"foobar"}
	err := runMain(args, mockShowHelp, dummyRunSend)
	expectError(t, "mocked help", err)
	if calledTimes != 1 {
		t.Errorf("calledWith: expected=%d actual=%d", 1, calledTimes)
	}
}

func Test_runMain_HelpOption(t *testing.T) {
	var calledTimes int
	mockShowHelp := func() error {
		calledTimes += 1
		return errors.New("mocked help")
	}
	args := []string{"--help"}
	err := runMain(args, mockShowHelp, dummyRunSend)
	expectError(t, "mocked help", err)
	if calledTimes != 1 {
		t.Errorf("calledWith: expected=%d actual=%d", 1, calledTimes)
	}
}
func Test_runMain_CallsRunSend(t *testing.T) {
	var calledWith []string
	mockRunSend := func(args []string) error {
		calledWith = args
		return errors.New("mocked error")
	}
	args := []string{"send", "--to", "foobar@example.com"}
	expectedCalledWith := []string{"--to", "foobar@example.com"}
	err := runMain(args, dummyShowHelp, mockRunSend)
	expectError(t, "mocked error", err)
	if !reflect.DeepEqual(expectedCalledWith, calledWith) {
		t.Errorf("calledWith: expected=%s actual=%s", expectedCalledWith, calledWith)
	}
}
