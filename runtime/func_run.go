package runtime

import (
	"github.com/mika-sandbox/otsukai"
	"github.com/mika-sandbox/otsukai/parser"
	"github.com/mika-sandbox/otsukai/runtime/context"
	re "github.com/mika-sandbox/otsukai/runtime/errors"
	"github.com/mika-sandbox/otsukai/runtime/helpers"
	"github.com/mika-sandbox/otsukai/runtime/value"
)

func InvokeRun(ctx context.IContext, arguments []parser.Argument) (value.IValueObject, error) {
	remote := helpers.GetNamedArgument(arguments, "remote")
	local := helpers.GetNamedArgument(arguments, "local")
	stdout := helpers.GetNamedArgument(arguments, "stdout")
	redirectToStdOut := false

	if stdout != nil {
		val := stdout.Expression.ValueExpression
		if val == nil || val.Value.Literal == nil {
			otsukai.Errf("the argument of run must be boolean literal")
			return nil, re.SYNTAX_ERROR
		}

		redirectToStdOut = val.Value.Literal.True != nil
	}

	if remote != nil && local != nil {
		otsukai.Errf("invalid argument: could not specify both of remote and local")
		return nil, re.RUNTIME_ERROR
	}

	if remote != nil {
		val := remote.Expression.ValueExpression
		if val == nil || val.Value.Literal == nil {
			otsukai.Errf("the argument of run must be string literal")
			return nil, re.SYNTAX_ERROR
		}

		command := val.Value.Literal.String
		session := ctx.GetRemoteSession()
		err := session.Run(*command, redirectToStdOut)
		if err != nil {
			return nil, err
		}
	}

	if local != nil {
		val := local.Expression.ValueExpression
		if val == nil || val.Value.Literal == nil {
			otsukai.Errf("the argument of run must be string literal")
			return nil, re.SYNTAX_ERROR
		}

		command := val.Value.Literal.String
		session := ctx.GetLocalSession()
		err := session.Run(*command, redirectToStdOut)
		if err != nil {
			otsukai.Errf("failed to execute command: %s", err)
			return nil, re.EXECUTION_ERROR
		}
	}

	return nil, nil
}
