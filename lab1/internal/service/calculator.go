package service

import "lab1/internal/model"

func Execute(op model.Operation, a, b int) int {
	return op.Execute(a, b)
}