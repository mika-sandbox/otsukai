package runtime

import (
	"github.com/mika-sandbox/otsukai"
	"github.com/mika-sandbox/otsukai/parser"
	"github.com/mika-sandbox/otsukai/runtime/context"
	re "github.com/mika-sandbox/otsukai/runtime/errors"
	"github.com/mika-sandbox/otsukai/runtime/value"
)

func InvokeFunction(ctx context.IContext, identifier string, arguments []parser.Argument, lambda *parser.LambdaExpression) (value.IValueObject, error) {
	switch identifier {
	case "set":
		return InvokeSet(ctx, arguments)

	case "task":
		return CollectTask(ctx, arguments, lambda)

	case "hook":
		return CollectHook(ctx, arguments, lambda)

	case "changed":
		return InvokeChanged(ctx, arguments)

	case "run":
		return InvokeRun(ctx, arguments)

	case "copy":
		return InvokeCopy(ctx, arguments)

	case "task_success":
		return InvokeTaskSuccess(ctx)
	}

	otsukai.Errf("the function `%v` is not declared in context", identifier)
	return nil, re.EXECUTION_ERROR
}
