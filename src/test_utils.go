package main

import "testing"

func expectNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("err: expected=nil actual=%s", err)
		t.FailNow()
	}
}

func expectError(t *testing.T, expected string, err error) {
	if err == nil {
		t.Errorf("err: expected=%s actual=nil", expected)
		t.FailNow()
	}
	if err.Error() != expected {
		t.Errorf("err: expected=%s actual=%s", expected, err)
		t.FailNow()
	}
}
