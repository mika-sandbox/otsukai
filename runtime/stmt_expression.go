package runtime

import (
	"github.com/mika-sandbox/otsukai"
	"github.com/mika-sandbox/otsukai/parser"
	"github.com/mika-sandbox/otsukai/runtime/context"
	re "github.com/mika-sandbox/otsukai/runtime/errors"
)

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
