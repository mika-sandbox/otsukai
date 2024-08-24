package otsukai

import (
	"errors"
)

const CONTEXT_GLOBAL = 1 << 0
const CONTEXT_COMPILATION = 1 << 1
const CONTEXT_TASK = 1 << 2

var SYNTAX_ERROR = errors.New("syntax error")
var RUNTIME_ERROR = errors.New("runtime error")

type IContext interface {
	SetVar(name string, value IValueObject)
	GetContextFlag() int
	GetStatements() []Statement
	GetTask(name *string) *Task
	CreateScope(statements []Statement) IContext
}

// ScopedContext is scoped context, use for in-task variable declarations
type ScopedContext struct {
	Statements []Statement
	Variables  map[string]IValueObject
	Tasks      *map[string]Task
}

func (ctx ScopedContext) SetVar(name string, value IValueObject) {
	ctx.Variables[name] = value
}

func (ctx ScopedContext) GetContextFlag() int {
	return CONTEXT_GLOBAL | CONTEXT_TASK
}

func (ctx ScopedContext) GetStatements() []Statement {
	return ctx.Statements
}

func (ctx ScopedContext) GetTask(name *string) *Task {
	tasks := *ctx.Tasks
	t, ok := tasks[*name]
	if ok {
		return &t
	}

	return nil
}

func (ctx ScopedContext) CreateScope(statements []Statement) IContext {
	return ScopedContext{
		Statements: statements,
		Variables:  ctx.Variables,
		Tasks:      ctx.Tasks,
	}
}

// Context is global context, use for global variable declarations
type Context struct {
	Statements []Statement
	Variables  map[string]IValueObject
	Tasks      map[string]Task
}

func (ctx Context) SetVar(name string, value IValueObject) {
	ctx.Variables[name] = value
}

func (ctx Context) GetContextFlag() int {
	return CONTEXT_GLOBAL | CONTEXT_COMPILATION
}

func (ctx Context) GetStatements() []Statement {
	return ctx.Statements
}

func (ctx Context) GetTask(name *string) *Task {
	t, ok := ctx.Tasks[*name]
	if ok {
		return &t
	}

	return nil
}

func (ctx Context) AddTask(name string, task Task) {
	if _, exists := ctx.Tasks[name]; exists {
		Errf("the task '%s' is already declared", name)
		return
	}

	ctx.Tasks[name] = task
}

func (ctx Context) CreateScope(statements []Statement) IContext {
	return ScopedContext{
		Statements: statements,
		Variables:  ctx.Variables,
		Tasks:      &ctx.Tasks,
	}
}

func NewContext(ruby *Entry) Context {
	return Context{
		Statements: ruby.Statements,
		Variables:  map[string]IValueObject{},
		Tasks:      map[string]Task{},
	}
}

func (ctx Context) Run() error {
	var err error
	// set default values
	ctx.Variables["default"] = &StringValueObject{val: "deploy"}

	if err = declarations(ctx); err != nil {
		return err
	}

	defaultTaskNameVar := ctx.Variables["default"]
	defaultTaskName, err := defaultTaskNameVar.ToString()
	if err != nil {
		return err
	}

	if err = run__task(ctx, defaultTaskName); err != nil {
		return err
	}

	return nil
}

func declarations(ctx Context) error {
	for _, statement := range ctx.GetStatements() {
		s := statement.Statement

		if s.IfStatement != nil || s.BlockStatement != nil {
			Errf("invalid statement: if or block statement could not placed in the compilation block")
			return errors.New("invalid statement: if or block statement could not placed in the compilation block")
		}

		expression := s.ExpressionStatement.Expression
		if expression.IfExpression != nil || expression.ValueExpression != nil || expression.LambdaExpression != nil {
			Errf("invalid statement: if expression, lambda expression or value could not placed in the compilation block")
			return errors.New("invalid statement: if expression, lambda expression or value could not placed in the compilation block")
		}

		identifier := expression.InvocationExpression.Expression.IdentifierNameExpression.Identifier.Name
		arguments := expression.InvocationExpression.ArgumentList.Argument
		lambda := expression.InvocationExpression.ArgumentList.LambdaExpression
		_, err := invoke(ctx, identifier, arguments, lambda)
		if err != nil {
			return err
		}
	}

	return nil
}

func run__task(ctx IContext, name *string) error {
	var err error

	t := ctx.GetTask(name)
	if len(t.BeforeHooks) != 0 {
		hooks := t.BeforeHooks

		for _, h := range hooks {
			err = run__statements(ctx.CreateScope(h.Statements))
			if err != nil {
				return err
			}
		}
	}

	err = run__statements(ctx.CreateScope(t.Statements))
	if err != nil {
		return err
	}

	if len(t.AfterHooks) != 0 {
		hooks := t.AfterHooks

		for _, h := range hooks {
			err = run__statements(ctx.CreateScope(h.Statements))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func run__statements(ctx IContext) error {
	for _, statement := range ctx.GetStatements() {
		s := statement.Statement

		// execution
		if s.IfStatement != nil {
			return run__if(ctx.CreateScope(s.IfStatement.Statements), s.IfStatement.Condition)
		}

		if s.BlockStatement != nil {
			return run__statements(ctx.CreateScope(s.BlockStatement.Statements))
		}

		if s.ExpressionStatement != nil {
			expression := s.ExpressionStatement
			err := run__expression(ctx, expression)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func run__if(ctx IContext, condition IfStatementOrExpressionConditionExpression) error {
	ret, err := run__condition(ctx, condition)

	if err != nil {
		return err
	}

	if ret {
		err = run__statements(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func run__condition(ctx IContext, condition IfStatementOrExpressionConditionExpression) (bool, error) {
	if condition.InvocationExpression != nil || condition.InvocationExpressionWithParen != nil {
		var identifier string
		var arguments []Argument
		var lambda *LambdaExpression

		if condition.InvocationExpression != nil {
			identifier = condition.InvocationExpression.Expression.IdentifierNameExpression.Identifier.Name
			arguments = condition.InvocationExpression.ArgumentList.Argument
			lambda = condition.InvocationExpression.ArgumentList.LambdaExpression
		}

		if condition.InvocationExpressionWithParen != nil {
			identifier = condition.InvocationExpressionWithParen.Expression.IdentifierNameExpression.Identifier.Name
			arguments = condition.InvocationExpressionWithParen.ArgumentList.Argument
			lambda = condition.InvocationExpressionWithParen.ArgumentList.LambdaExpression
		}

		val, err := invoke(ctx, identifier, arguments, lambda)
		if err != nil {
			return false, err
		}

		b, err := val.ToBoolean()
		if err != nil {
			return false, err
		}

		return *b, nil
	}

	if condition.ValueExpression != nil {
		val, err := condition.ValueExpression.Value.ToValueObject()
		if err != nil {
			return false, err
		}

		b, err := val.ToBoolean()
		if err != nil {
			return false, err
		}

		return *b, nil
	}

	return false, nil
}

func run__expression(ctx IContext, expression *ExpressionStatement) error {
	return nil
}

func set(ctx IContext, arguments []Argument) (IValueObject, error) {
	scope := ctx.GetContextFlag()

	if scope&CONTEXT_GLOBAL != CONTEXT_GLOBAL {
		Errf("invalid context")
		return nil, SYNTAX_ERROR
	}
	if len(arguments) != 1 {
		Errf("set method must be one argument")
		return nil, SYNTAX_ERROR
	}

	name := arguments[0].Identifier
	if name == nil {
		Errf("set must be named arguments")
		return nil, SYNTAX_ERROR
	}

	value, err := arguments[0].Expression.ValueExpression.Value.ToValueObject()
	if err != nil {
		return nil, err
	}

	ctx.SetVar(*name, value)
	return nil, nil
}

func task(ctx IContext, arguments []Argument, lambda *LambdaExpression) (IValueObject, error) {
	scope := ctx.GetContextFlag()

	if scope&CONTEXT_COMPILATION != CONTEXT_COMPILATION {
		Errf("invalid context")
		return nil, SYNTAX_ERROR
	}

	if len(arguments) != 1 {
		Errf("task method must have one argument, with lambda")
		return nil, SYNTAX_ERROR
	}

	name := arguments[0].Expression.ValueExpression.Value.HashSymbol
	if name == nil {
		Errf("the first argument of task must be symbol")
		return nil, SYNTAX_ERROR
	}

	context, ok := ctx.(Context)
	if !ok {
		return nil, errors.New("failed to cast context; context must be global")
	}

	t := CreateTask(name.Identifier, lambda.Statements)
	context.AddTask(name.Identifier, t)
	return nil, nil
}

func hook(ctx IContext, arguments []Argument, lambda *LambdaExpression) (IValueObject, error) {
	scope := ctx.GetContextFlag()

	if scope&CONTEXT_COMPILATION != CONTEXT_COMPILATION {
		Errf("invalid context")
		return nil, SYNTAX_ERROR
	}

	if len(arguments) != 1 {
		Errf("task method must have one argument, with lambda")
		return nil, SYNTAX_ERROR
	}

	name := arguments[0].Expression.ValueExpression.Value.HashSymbol
	if name == nil {
		Errf("the first argument of task must be symbol")
		return nil, SYNTAX_ERROR
	}

	context, ok := ctx.(Context)
	if !ok {
		return nil, errors.New("failed to cast context; context must be global")
	}

	h := *arguments[0].Identifier
	t, ok := context.Tasks[name.Identifier]
	if ok {
		if h == "before" {
			t.BeforeHooks = append(t.BeforeHooks, CreateHook(lambda.Statements))
		}

		if h == "after" {
			t.AfterHooks = append(t.AfterHooks, CreateHook(lambda.Statements))
		}

		context.Tasks[name.Identifier] = t

		return nil, nil
	}

	Errf("invalid hooks, hooks must be one of before or after")
	return nil, RUNTIME_ERROR
}

func invoke(ctx IContext, identifier string, arguments []Argument, lambda *LambdaExpression) (IValueObject, error) {
	switch identifier {
	case "set":
		return set(ctx, arguments)

	case "task":
		return task(ctx, arguments, lambda)

	case "hook":
		return hook(ctx, arguments, lambda)
	}

	Errf("the function `%v` is not declared in context", identifier)
	return nil, errors.New("the function is not declare in context")
}
