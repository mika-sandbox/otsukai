package helpers

import (
	"otsukai/parser"
)

func GetNamedArgument(arguments []parser.Argument, name string) *parser.Argument {
	for _, argument := range arguments {
		identifier := argument.Identifier
		if identifier != nil && *identifier == name {
			return &argument
		}
	}

	return nil
}

func GetStringLiteral(argument *parser.Argument) *string {
	if argument.Expression.ValueExpression == nil || argument.Expression.ValueExpression.Value.Literal == nil {
		return nil
	}

	return argument.Expression.ValueExpression.Value.Literal.String
}

func GetBoolLiteral(argument *parser.Argument, defaultValue bool) bool {
	if argument == nil || argument.Expression.ValueExpression == nil || argument.Expression.ValueExpression.Value.Literal == nil {
		return defaultValue
	}

	val := argument.Expression.ValueExpression.Value.Literal.True != nil
	return val
}

func GetSymbol(argument *parser.Argument) *string {
	if argument.Expression.ValueExpression == nil || argument.Expression.ValueExpression.Value.HashSymbol == nil {
		return nil
	}

	return &argument.Expression.ValueExpression.Value.HashSymbol.Identifier
}
