package otsukai

import "errors"

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

func (ctx *Context) SetVar(name string, value IValueObject) {
	ctx.Variables[name] = value
}

func (ctx *Context) GetStatements() []Statement {
	return ctx.Statements
}

func (ctx *Context) CreateScope(statements []Statement) ScopedContext {
	return ScopedContext{
		Statements: statements,
		Variables:  ctx.Variables,
	}
}

func NewContext(ruby *Entry) Context {
	return Context{
		Statements: ruby.Statements,
		Variables:  map[string]IValueObject{},
	}
}

func (ctx *Context) invoke(invocation *InvocationExpression) (IValueObject, error) {
	identifier := invocation.Expression.IdentifierNameExpression.Identifier.Name
	arguments := invocation.ArgumentList.Argument
	// lambda := invocation.ArgumentList.LambdaExpression

	switch identifier {
	case "set":
		if len(arguments) != 1 {
			Err("set method must be one argument")
			return nil, errors.New("set method must be one argument")
		}

		name := arguments[0].Identifier
		if name == nil {
			Err("set must be named arguments")
			return nil, errors.New("set must be named arguments")
		}

		value, err := arguments[0].Expression.ValueExpression.Value.ToValueObject()
		if err != nil {
			return nil, err
		}

		ctx.SetVar(*name, value)
		return nil, nil

	case "task":
		break

	case "hook":
		break
	}

	Err("the function `%v` is not declared in context", identifier)
	return nil, errors.New("the function is not declare in context")
}

func (ctx *Context) Run() error {
	var err error
	// set default values
	ctx.Variables["default"] = &StringValueObject{val: "deploy"}

	for _, statement := range ctx.GetStatements() {
		s := statement.Statement

		if s.IfStatement != nil || s.BlockStatement != nil {
			Err("invalid statement: if or block statement could not placed in the compilation block")
			return errors.New("invalid statement: if or block statement could not placed in the compilation block")
		}

		expression := s.ExpressionStatement.Expression
		if expression.IfExpression != nil || expression.ValueExpression != nil || expression.LambdaExpression != nil {
			Err("invalid statement: if expression, lambda expression or value could not placed in the compilation block")
			return errors.New("invalid statement: if expression, lambda expression or value could not placed in the compilation block")
		}

		_, err = ctx.invoke(expression.InvocationExpression)
		if err != nil {
			return err
		}
	}

	return nil
}
