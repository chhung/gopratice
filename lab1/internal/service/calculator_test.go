package service

import (
	"testing"

	"lab1/internal/model"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name string
		op   model.Operation
		a    int
		b    int
		want int
	}{
		{name: "add", op: model.AddOperation{}, a: 10, b: 20, want: 30},
		{name: "subtract", op: model.SubtractOperation{}, a: 10, b: 20, want: -10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Execute(tt.op, tt.a, tt.b); got != tt.want {
				t.Fatalf("Execute() = %d, want %d", got, tt.want)
			}
		})
	}
}