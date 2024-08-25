package context

import (
	"otsukai/parser"
	"otsukai/runtime/task"
	"otsukai/runtime/value"
)

// ScopedContext is scoped context, use for in-task variable declarations
type ScopedContext struct {
	Statements []parser.Statement
	Variables  map[string]value.IValueObject
	Tasks      *map[string]task.Task
	Phase      int
}

func (ctx *ScopedContext) SetPhase(phase int) {
	ctx.Phase = phase
}

func (ctx *ScopedContext) GetPhase() int {
	return ctx.Phase
}

func (ctx *ScopedContext) SetVar(name string, value value.IValueObject) {
	ctx.Variables[name] = value
}

func (ctx *ScopedContext) GetVar(name string) value.IValueObject {
	return ctx.Variables[name]
}

func (ctx *ScopedContext) GetContextFlag() int {
	return CONTEXT_GLOBAL | CONTEXT_TASK
}

func (ctx *ScopedContext) GetStatements() []parser.Statement {
	return ctx.Statements
}

func (ctx *ScopedContext) GetTask(name *string) *task.Task {
	tasks := *ctx.Tasks
	t, ok := tasks[*name]
	if ok {
		return &t
	}

	return nil
}


func (ctx *ScopedContext) CreateScope(statements []parser.Statement) IContext {
	return &ScopedContext{
		Phase:      ctx.Phase,
		Statements: statements,
		Variables:  ctx.Variables,
		Tasks:      ctx.Tasks,
		Remote:     ctx.Remote,
		Local:      ctx.Local,
	}
}
