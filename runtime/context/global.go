package context

import (
	"otsukai"
	"otsukai/parser"
	"otsukai/runtime/task"
	"otsukai/runtime/value"
)

// Context is global context, use for global variable declarations
type Context struct {
	Statements []parser.Statement
	Variables  map[string]value.IValueObject
	Tasks      map[string]task.Task
	Phase      int
}

func (ctx *Context) SetPhase(phase int) {
	ctx.Phase = phase
}

func (ctx *Context) GetPhase() int {
	return ctx.Phase
}

func (ctx *Context) SetVar(name string, value value.IValueObject) {
	ctx.Variables[name] = value
}

func (ctx *Context) GetVar(name string) value.IValueObject {
	return ctx.Variables[name]
}

func (ctx *Context) GetContextFlag() int {
	return CONTEXT_GLOBAL | CONTEXT_COMPILATION
}

func (ctx *Context) GetStatements() []parser.Statement {
	return ctx.Statements
}

func (ctx *Context) GetTask(name *string) *task.Task {
	t, ok := ctx.Tasks[*name]
	if ok {
		return &t
	}

	return nil
}

func (ctx *Context) AddTask(name string, task task.Task) {
	if _, exists := ctx.Tasks[name]; exists {
		otsukai.Errf("the task '%s' is already declared", name)
		return
	}

	ctx.Tasks[name] = task
}

func (ctx *Context) CreateScope(statements []parser.Statement) IContext {
	return &ScopedContext{
		Phase:      ctx.Phase,
		Statements: statements,
		Variables:  ctx.Variables,
		Tasks:      &ctx.Tasks,
		Remote:     ctx.Remote,
		Local:      ctx.Local,
	}
}

//

func NewContext(ruby *parser.Entry) Context {
	return Context{
		Phase:      0,
		Statements: ruby.Statements,
		Variables:  map[string]value.IValueObject{},
		Tasks:      map[string]task.Task{},
	}
}
