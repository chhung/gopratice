package model

type Operation interface {
	Execute(a, b int) int
}

type AddOperation struct{}

func (AddOperation) Execute(a, b int) int {
	return a + b
}

type SubtractOperation struct{}

func (SubtractOperation) Execute(a, b int) int {
	return a - b
}