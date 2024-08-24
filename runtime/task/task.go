package task

import "otsukai/parser"

type Task struct {
	Name        string
	Statements  []parser.Statement
	BeforeHooks []Hook
	AfterHooks  []Hook
}

type Hook struct {
	Statements []parser.Statement
}

func CreateTask(name string, statements []parser.Statement) Task {
	return Task{
		Name:        name,
		Statements:  statements,
		BeforeHooks: make([]Hook, 0),
		AfterHooks:  make([]Hook, 0),
	}
}

func CreateHook(statements []parser.Statement) Hook {
	return Hook{
		Statements: statements,
	}
}
