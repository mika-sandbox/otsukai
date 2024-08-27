package runtime

import (
	"otsukai/runtime/context"
	re "otsukai/runtime/errors"
	"otsukai/runtime/value"
)

func InvokeTaskSuccess(ctx context.IContext) (value.IValueObject, error) {
	if ctx.GetLastStatus() == context.CONTEXT_STATUS_SUCCESS {
		return value.BooleanValueObject{Val: true}, nil
	}

	if ctx.GetLastStatus() == context.CONTEXT_STATUS_ERROR {
		return value.BooleanValueObject{Val: false}, nil
	}

	return nil, re.EXECUTION_ERROR
}
