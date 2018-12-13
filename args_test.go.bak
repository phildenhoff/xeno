package main

import (
	"fmt"
	"strings"
	"testing"
	"reflect"
)

func TestNonEmptyArgs(t *testing.T) {
	// Verify that we get a non-empty string array back for an input
	x := sanitiseArguments("ls -la")
	if len(x) != 2 {
		t.Fail()
	}
}

func TestArgContentsMatchInput(t *testing.T) {
	// Verify that what we put in is what we get out
	x := sanitiseArguments("ls -la")
	fmt.Println(x)
	if x[0] != "ls" || x[1] != "-la" {
		t.Fail()
	}
}

func TestGlobQuotes(t *testing.T) {
	// Verify that double quoted strings are one arg
	x := sanitiseArguments("echo \"hello world\"")
	var expectedResult = []string {"echo", "\"hello world\""}
	if !reflect.DeepEqual(x, expectedResult) {
		var errMsg = fmt.Sprintf("Double quoted argument should be one item. Received [%s], expected [%s]",
						strings.Join(x, ","), strings.Join(expectedResult, ","))
		t.Error(errMsg)
	}
}

func TestBrokenGlob(t *testing.T) {
	// Verify that double quoted strings are one arg even if the last quote never appears
	x := sanitiseArguments("echo \"hello world")
	var expectedResult = []string {"echo", "\"hello world"}
	if !reflect.DeepEqual(x, expectedResult) {
		var errMsg = fmt.Sprintf("Double quoted argument should be one item. Received [%s], expected [%s]",
			strings.Join(x, ","), strings.Join(expectedResult, ","))
		t.Error(errMsg)
	}
}