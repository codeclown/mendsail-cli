package main

import (
	"errors"
	"reflect"
	"testing"
)

func dummyRunSend(args []string) error {
	return nil
}

func Test_RunMain_NoArgs(t *testing.T) {
	var args []string
	err := runMain(args, dummyRunSend)
	expectError(t, "expected a subcommand (e.g. 'mendsail send ...')", err)
}

func Test_RunMain_UnknownCommand(t *testing.T) {
	args := []string{"foobar"}
	err := runMain(args, dummyRunSend)
	expectError(t, "expected a subcommand (e.g. 'mendsail send ...')", err)
}

func Test_RunMain_CallsRunSend(t *testing.T) {
	var calledWith []string
	mockRunSend := func(args []string) error {
		calledWith = args
		return errors.New("mocked error")
	}
	args := []string{"send", "--to", "foobar@example.com"}
	expectedCalledWith := []string{"--to", "foobar@example.com"}
	err := runMain(args, mockRunSend)
	expectError(t, "mocked error", err)
	if !reflect.DeepEqual(expectedCalledWith, calledWith) {
		t.Errorf("calledWith: expected=%s actual=%s", expectedCalledWith, calledWith)
	}
}
