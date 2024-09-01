package runtime

import (
	"errors"
	"github.com/mika-sandbox/otsukai"
	"github.com/mika-sandbox/otsukai/parser"
	"github.com/mika-sandbox/otsukai/runtime/context"
	re "github.com/mika-sandbox/otsukai/runtime/errors"
	"github.com/mika-sandbox/otsukai/runtime/session"
	"github.com/mika-sandbox/otsukai/runtime/task"
	"github.com/mika-sandbox/otsukai/runtime/value"
	"time"
)

func getTimeout(ctx *context.Context) *time.Duration {
	timeout := ctx.GetVar("timeout")
	if timeout == nil {
		return nil
	}

	val, err := timeout.ToFloat64()
	if err != nil {
		return nil
	}

	seconds := time.Duration(int(*val)) * time.Second
	return &seconds
}

func InvokeTask(ctx *context.Context, name *string) error {
	var err error

	t := ctx.GetTask(name)
	c := ctx.GetVar("remote")
	timeout := getTimeout(ctx)

	remote, err := session.CreateRemoteSession(&session.CreateRemoteSessionOpts{
		Remote:  c,
		Timeout: timeout,
	})
	if err != nil {
		return err
	}
	defer remote.Close()

	local, err := session.CreateLocalSession()
	if err != nil {
		return err
	}

	ctx.SetSession(remote, local)

	if len(t.BeforeHooks) != 0 {
		hooks := t.BeforeHooks

		for _, h := range hooks {
			err = InvokeStatements(ctx.CreateScope(h.Statements))
			if err != nil {
				return err
			}
		}
	}

	err = InvokeStatements(ctx.CreateScope(t.Statements))
	if err == nil {
		ctx.SetLastStatus(context.CONTEXT_STATUS_SUCCESS)
	} else {
		ctx.SetLastStatus(context.CONTEXT_STATUS_ERROR)
	}

	if len(t.AfterHooks) != 0 {
		hooks := t.AfterHooks

		for _, h := range hooks {
			err = InvokeStatements(ctx.CreateScope(h.Statements))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CollectDeclarations(ctx *context.Context) error {
	for _, statement := range ctx.GetStatements() {
		s := statement.Statement

		if s.IfStatement != nil || s.BlockStatement != nil {
			otsukai.Errf("invalid statement: if or block statement could not placed in the compilation block")
			return errors.New("invalid statement: if or block statement could not placed in the compilation block")
		}

		expression := s.ExpressionStatement.Expression
		if expression.IfExpression != nil || expression.ValueExpression != nil || expression.LambdaExpression != nil {
			otsukai.Errf("invalid statement: if expression, lambda expression or value.go could not placed in the compilation block")
			return errors.New("invalid statement: if expression, lambda expression or value.go could not placed in the compilation block")
		}

		identifier := expression.InvocationExpression.Expression.IdentifierNameExpression.Identifier.Name
		arguments := expression.InvocationExpression.ArgumentList.Argument
		lambda := expression.InvocationExpression.ArgumentList.LambdaExpression
		_, err := InvokeFunction(ctx, identifier, arguments, lambda)
		if err != nil {
			return err
		}
	}

	return nil
}

func CollectTask(ctx context.IContext, arguments []parser.Argument, lambda *parser.LambdaExpression) (value.IValueObject, error) {
	phase := ctx.GetPhase()
	if phase != PHASE_COLLECT {
		return nil, re.RUNTIME_ERROR
	}

	scope := ctx.GetContextFlag()
	if scope&context.CONTEXT_COMPILATION != context.CONTEXT_COMPILATION {
		otsukai.Errf("invalid context")
		return nil, re.SYNTAX_ERROR
	}

	if len(arguments) != 1 {
		otsukai.Errf("task method must have one argument, with lambda")
		return nil, re.SYNTAX_ERROR
	}

	name := arguments[0].Expression.ValueExpression.Value.HashSymbol
	if name == nil {
		otsukai.Errf("the first argument of task must be symbol")
		return nil, re.SYNTAX_ERROR
	}

	context, ok := ctx.(*context.Context)
	if !ok {
		return nil, errors.New("failed to cast context; context must be global")
	}

	t := task.CreateTask(name.Identifier, lambda.Statements)
	context.AddTask(name.Identifier, t)
	return nil, nil
}

func CollectHook(ctx context.IContext, arguments []parser.Argument, lambda *parser.LambdaExpression) (value.IValueObject, error) {
	phase := ctx.GetPhase()
	if phase != PHASE_COLLECT {
		return nil, re.RUNTIME_ERROR
	}

	scope := ctx.GetContextFlag()
	if scope&context.CONTEXT_COMPILATION != context.CONTEXT_COMPILATION {
		otsukai.Errf("invalid context")
		return nil, re.SYNTAX_ERROR
	}

	if len(arguments) != 1 {
		otsukai.Errf("task method must have one argument, with lambda")
		return nil, re.SYNTAX_ERROR
	}

	name := arguments[0].Expression.ValueExpression.Value.HashSymbol
	if name == nil {
		otsukai.Errf("the first argument of task must be symbol")
		return nil, re.SYNTAX_ERROR
	}

	context, ok := ctx.(*context.Context)
	if !ok {
		return nil, errors.New("failed to cast context; context must be global")
	}

	h := *arguments[0].Identifier
	t, ok := context.Tasks[name.Identifier]
	if ok {
		if h == "before" {
			t.BeforeHooks = append(t.BeforeHooks, task.CreateHook(lambda.Statements))
		}

		if h == "after" {
			t.AfterHooks = append(t.AfterHooks, task.CreateHook(lambda.Statements))
		}

		context.Tasks[name.Identifier] = t

		return nil, nil
	}

	otsukai.Errf("invalid hook, hook must be one of before or after")
	return nil, re.RUNTIME_ERROR
}
