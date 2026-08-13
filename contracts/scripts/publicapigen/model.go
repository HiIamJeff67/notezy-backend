package main

import "go/ast"

type typeDefinition struct {
	directory string
	expr      ast.Expr
}

type endpoint struct {
	Method       string
	Path         string
	Tag          string
	OperationID  string
	RequestType  string
	ResponseType string
}
