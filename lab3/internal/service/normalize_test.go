package service

import "testing"

func TestNormalizeInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "no spaces", input: "1234", want: "1234"},
		{name: "spaces removed", input: "12 34", want: "1234"},
		{name: "tabs and newlines removed", input: "12\t3\n4", want: "1234"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeInput(tt.input); got != tt.want {
				t.Fatalf("NormalizeInput() = %q, want %q", got, tt.want)
			}
		})
	}
}
