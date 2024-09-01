package runtime

import (
	"github.com/mika-sandbox/otsukai/logger"
	"github.com/mika-sandbox/otsukai/parser"
	"github.com/mika-sandbox/otsukai/runtime/context"
	re "github.com/mika-sandbox/otsukai/runtime/errors"
	"github.com/mika-sandbox/otsukai/runtime/value"
)

func InvokeSet(ctx context.IContext, arguments []parser.Argument) (value.IValueObject, error) {
	scope := ctx.GetContextFlag()

	if scope&context.CONTEXT_GLOBAL != context.CONTEXT_GLOBAL {
		logger.Errf("invalid context")
		return nil, re.SYNTAX_ERROR
	}
	if len(arguments) != 1 {
		logger.Errf("set method must be one argument")
		return nil, re.SYNTAX_ERROR
	}

	name := arguments[0].Identifier
	if name == nil {
		logger.Errf("set must be named arguments")
		return nil, re.SYNTAX_ERROR
	}

	value, err := value.ToValueObject(arguments[0].Expression.ValueExpression.Value)
	if err != nil {
		return nil, err
	}

	ctx.SetVar(*name, value)
	return nil, nil
}
