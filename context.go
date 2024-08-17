package otsukai

import (
	"errors"
	"slices"
)

type IContext interface {
	SetVar(name string, value IValueObject)
	GetStatements() []Statement
}

// ScopedContext is scoped context, use for in-task variable declarations
type ScopedContext struct {
	Statements []Statement
	Variables  map[string]IValueObject
}

func (ctx ScopedContext) SetVar(name string, value IValueObject) {
	ctx.Variables[name] = value
}

func (ctx ScopedContext) GetStatements() []Statement {
	return ctx.Statements
}

// Context is global context, use for global variable declarations
type Context struct {
	Statements []Statement
	Variables  map[string]IValueObject
}

func (ctx Context) SetVar(name string, value IValueObject) {
	ctx.Variables[name] = value
}

func (ctx Context) GetStatements() []Statement {
	return ctx.Statements
}

func NewContext(ruby *Entry) Context {
	return Context{
		Statements: ruby.Statements,
		Variables:  map[string]IValueObject{},
	}
}

func (ctx *Context) Run() error {
	var err error
	// set default values
	ctx.Variables["default"] = &StringValueObject{val: "deploy"}

	// evaluate set in global scope
	err = runDeclarations(ctx)
	if err != nil {
		return err
	}

	// evaluate task declarations
	tasks, err := runTaskDeclarations(ctx)
	if err != nil {
		return err
	}

	err = assignHookDeclarations(ctx, tasks)
	if err != nil {
		return err
	}

	return nil
}

func runDeclarations(ctx IContext) error {
	for _, statement := range ctx.GetStatements() {
		if statement.SetStatement == nil {
			continue
		}

		set := statement.SetStatement
		identifier := set.Expression.Identifier.Identifier
		value, err := set.Expression.Value.ToValueObject()
		if err != nil {
			return err
		}

		ctx.SetVar(identifier, value)
	}

	return nil
}

func runTaskDeclarations(ctx *Context) ([]Task, error) {
	var tasks []Task

	for _, statement := range ctx.GetStatements() {
		if statement.TaskStatement == nil {
			continue
		}

		task := statement.TaskStatement
		tasks = append(tasks, Task{
			Name:        task.Identifier.Identifier,
			Statements:  task.DoStatement.WithStatements.Statements,
			BeforeHooks: []Hook{},
			AfterHooks:  []Hook{},
		})
	}

	return tasks, nil
}

func assignHookDeclarations(ctx *Context, tasks []Task) error {
	for _, statement := range ctx.GetStatements() {
		if statement.HookStatement == nil {
			continue
		}

		hook := statement.HookStatement
		event := hook.Parameters.Identifier.Identifier
		value, err := hook.Parameters.Value.ToValueObject()
		if err != nil {
			return err
		}
		on, err := value.ToString()
		if err != nil {
			return err // must be string
		}

		idx := slices.IndexFunc(tasks, func(t Task) bool { return t.Name == *on })
		task := &tasks[idx]

		if event == "before" {
			task.BeforeHooks = append(task.BeforeHooks, Hook{Statements: hook.DoStatement.WithStatements.Statements})
		} else if event == "after" {
			task.AfterHooks = append(task.AfterHooks, Hook{Statements: hook.DoStatement.WithStatements.Statements})
		} else {
			return errors.New("unknown hook")
		}

	}

	return nil
}
