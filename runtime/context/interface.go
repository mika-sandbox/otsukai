package context

import (
	"otsukai/parser"
	"otsukai/runtime/session"
	"otsukai/runtime/task"
	"otsukai/runtime/value"
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
}
