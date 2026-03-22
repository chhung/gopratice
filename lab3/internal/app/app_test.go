package app

import (
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	input := strings.NewReader("12 34\n")
	var output strings.Builder

	if err := Run(input, &output); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got := output.String(); got != "Input: Normalized Input: 1234\n" {
		t.Fatalf("Run() output = %q, want %q", got, "Input: Normalized Input: 1234\n")
	}
}
