package runtime

import (
	"github.com/mika-sandbox/otsukai/parser"
	"github.com/mika-sandbox/otsukai/runtime/context"
	"github.com/mika-sandbox/otsukai/runtime/value"
)

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
