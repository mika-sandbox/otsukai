package context

import (
	"github.com/mika-sandbox/otsukai/parser"
	"github.com/mika-sandbox/otsukai/runtime/session"
	"github.com/mika-sandbox/otsukai/runtime/task"
	"github.com/mika-sandbox/otsukai/runtime/value"
)

type IContext interface {
	SetVar(name string, value value.IValueObject)
	GetVar(name string) value.IValueObject
	SetPhase(phase int)
	GetPhase() int
	GetContextFlag() int
	GetStatements() []parser.Statement
	GetTask(name *string) *task.Task
	CreateScope(statements []parser.Statement) IContext
	SetSession(remote session.ISession, local session.ISession)
	GetRemoteSession() session.ISession
	GetLocalSession() session.ISession
	SetLastStatus(status int)
	GetLastStatus() int
}
