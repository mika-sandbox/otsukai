package context

import (
	"github.com/mika-sandbox/otsukai/logger"
	"github.com/mika-sandbox/otsukai/parser"
	"github.com/mika-sandbox/otsukai/runtime/session"
	"github.com/mika-sandbox/otsukai/runtime/task"
	"github.com/mika-sandbox/otsukai/runtime/value"
)

// Context is global context, use for global variable declarations
type Context struct {
	Statements []parser.Statement
	Variables  map[string]value.IValueObject
	Tasks      map[string]task.Task
	remote     session.ISession
	local      session.ISession
	phase      int
	status     int
}

func (ctx *Context) SetPhase(phase int) {
	ctx.phase = phase
}

func (ctx *Context) GetPhase() int {
	return ctx.phase
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

func (ctx *Context) SetSession(remote session.ISession, local session.ISession) {
	ctx.remote = remote
	ctx.local = local
}

func (ctx *Context) GetRemoteSession() session.ISession {
	return ctx.remote
}

func (ctx *Context) GetLocalSession() session.ISession {
	return ctx.local
}

func (ctx *Context) AddTask(name string, task task.Task) {
	if _, exists := ctx.Tasks[name]; exists {
		logger.Errf("the task '%s' is already declared", name)
		return
	}

	ctx.Tasks[name] = task
}

func (ctx *Context) CreateScope(statements []parser.Statement) IContext {
	context := &ScopedContext{
		Statements: statements,
		Variables:  ctx.Variables,
		Tasks:      &ctx.Tasks,
	}
	context.SetPhase(ctx.GetPhase())
	context.SetSession(ctx.GetRemoteSession(), ctx.GetLocalSession())
	context.SetLastStatus(ctx.GetLastStatus())

	return context
}

func (ctx *Context) SetLastStatus(status int) {
	ctx.status = status
}

func (ctx *Context) GetLastStatus() int {
	return ctx.status
}

//

func NewContext(ruby *parser.Entry) Context {
	return Context{
		phase:      0,
		Statements: ruby.Statements,
		Variables:  map[string]value.IValueObject{},
		Tasks:      map[string]task.Task{},
	}
}
