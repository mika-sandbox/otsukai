package runtime

import (
	"otsukai"
	"otsukai/parser"
	"otsukai/runtime/context"
	re "otsukai/runtime/errors"
	"otsukai/runtime/value"
)

const PHASE_UNSET = 0
const PHASE_COLLECT = 1
const PHASE_RUN = 2

func Run(ctx *context.Context) error {
	var err error
	// set default values
	ctx.Variables["default"] = &value.StringValueObject{Val: "deploy"}
	ctx.SetPhase(PHASE_COLLECT)
	ctx.SetLastStatus(context.CONTEXT_STATUS_PENDING)

	if err = CollectDeclarations(ctx); err != nil {
		return err
	}

	ctx.SetPhase(PHASE_RUN)
	defaultTaskNameVar := ctx.Variables["default"]
	defaultTaskName, err := defaultTaskNameVar.ToString()
	if err != nil {
		return err
	}

	if err = InvokeTask(ctx, defaultTaskName); err != nil {
		return err
	}

	return nil
}

func InvokeStatements(ctx context.IContext) error {
	for _, statement := range ctx.GetStatements() {
		s := statement.Statement

		// execution
		if s.IfStatement != nil {
			return InvokeIfStatement(ctx.CreateScope(s.IfStatement.Statements), s.IfStatement.Condition)
		}

		if s.BlockStatement != nil {
			return InvokeStatements(ctx.CreateScope(s.BlockStatement.Statements))
		}

		if s.ExpressionStatement != nil {
			expression := s.ExpressionStatement
			err := InvokeExpression(ctx, expression)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func InvokeIfStatement(ctx context.IContext, condition parser.IfStatementOrExpressionConditionExpression) error {
	ret, err := GetConditionResult(ctx, condition)

	if err != nil {
		return err
	}

	if ret {
		err = InvokeStatements(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetConditionResult(ctx context.IContext, condition parser.IfStatementOrExpressionConditionExpression) (bool, error) {
	if condition.InvocationExpression != nil || condition.InvocationExpressionWithParen != nil {
		var identifier string
		var arguments []parser.Argument
		var lambda *parser.LambdaExpression

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

		val, err := InvokeFunction(ctx, identifier, arguments, lambda)
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
		val, err := value.ToValueObject(condition.ValueExpression.Value)
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

func InvokeExpression(ctx context.IContext, expression *parser.ExpressionStatement) error {
	var err error

	expr := expression.Expression

	// if
	if expr.IfExpression != nil {
		err = InvokeIfStatement(ctx.CreateScope(expr.IfExpression.Statements), expr.IfExpression.Condition)
		if err != nil {
			return err
		}
	}

	// invocation
	if expr.InvocationExpression != nil || expr.InvocationExpressionWithParen != nil {
		var identifier string
		var arguments []parser.Argument
		var lambda *parser.LambdaExpression

		if expr.InvocationExpression != nil {
			identifier = expr.InvocationExpression.Expression.IdentifierNameExpression.Identifier.Name
			arguments = expr.InvocationExpression.ArgumentList.Argument
			lambda = expr.InvocationExpression.ArgumentList.LambdaExpression
		}

		if expr.InvocationExpressionWithParen != nil {
			identifier = expr.InvocationExpressionWithParen.Expression.IdentifierNameExpression.Identifier.Name
			arguments = expr.InvocationExpressionWithParen.ArgumentList.Argument
			lambda = expr.InvocationExpressionWithParen.ArgumentList.LambdaExpression
		}

		_, err := InvokeFunction(ctx, identifier, arguments, lambda)
		if err != nil {
			return err
		}

		return nil
	}

	otsukai.Errf("invalid context")
	return re.RUNTIME_ERROR
}

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
