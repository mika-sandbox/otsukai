package context

import (
	"github.com/mika-sandbox/otsukai/parser"
	"github.com/mika-sandbox/otsukai/runtime/session"
	"github.com/mika-sandbox/otsukai/runtime/task"
	"github.com/mika-sandbox/otsukai/runtime/value"
)

// ScopedContext is scoped context, use for in-task variable declarations
type ScopedContext struct {
	Statements []parser.Statement
	Variables  map[string]value.IValueObject
	Tasks      *map[string]task.Task
	phase      int
	remote     session.ISession
	local      session.ISession
	status     int
}

func (ctx *ScopedContext) SetPhase(phase int) {
	ctx.phase = phase
}

func (ctx *ScopedContext) GetPhase() int {
	return ctx.phase
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

func (ctx *ScopedContext) SetSession(remote session.ISession, local session.ISession) {
	ctx.remote = remote
	ctx.local = local
}

func (ctx *ScopedContext) GetRemoteSession() session.ISession {
	return ctx.remote
}

func (ctx *ScopedContext) GetLocalSession() session.ISession {
	return ctx.local
}

func (ctx *ScopedContext) SetLastStatus(status int) {
	ctx.status = status
}

func (ctx *ScopedContext) GetLastStatus() int {
	return ctx.status
}

func (ctx *ScopedContext) CreateScope(statements []parser.Statement) IContext {
	return &ScopedContext{

		Statements: statements,
		Variables:  ctx.Variables,
		Tasks:      ctx.Tasks,
		phase:      ctx.phase,
		remote:     ctx.remote,
		local:      ctx.local,
		status:     ctx.GetLastStatus(),
	}
}
