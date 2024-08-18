package otsukai

import (
	require "github.com/alecthomas/assert/v2"
	"github.com/alecthomas/repr"
	"os"
	"testing"
)

func c(t *testing.T, path string) string {
	content, err := os.ReadFile(path)
	require.NoError(t, err)

	return string(content)
}

func TestExample_DockerCompose(t *testing.T) {
	content := c(t, "./examples/docker-compose/otsukai.rb")
	ret, _ := Parser.ParseString("", content)
	// require.NoError(t, err)

	// require.Equal(t, 4, len(ret.Statements))

	// set default: :deploy
	{
		invocation := ret.Statements[0].Statement.ExpressionStatement.Expression.InvocationExpression
		require.Equal(t, "set", invocation.Expression.IdentifierNameExpression.Identifier.Name)

		args := invocation.ArgumentList.Argument
		lambda := invocation.ArgumentList.LambdaExpression
		require.Equal(t, nil, lambda)
		require.Equal(t, 1, len(args))

		arg := args[0]
		key := *arg.Identifier
		value := arg.Expression.ValueExpression.Value.HashSymbol
		require.Equal(t, "default", key)
		require.Equal(t, "deploy", value.Identifier)
	}

	// set target: { host: "yuuka.natsuneko.net", user: "ubuntu" }
	{
		invocation := ret.Statements[1].Statement.ExpressionStatement.Expression.InvocationExpression
		require.Equal(t, "set", invocation.Expression.IdentifierNameExpression.Identifier.Name)

		args := invocation.ArgumentList.Argument
		lambda := invocation.ArgumentList.LambdaExpression
		require.Equal(t, nil, lambda)
		require.Equal(t, 1, len(args))

		arg := args[0]
		key := *arg.Identifier
		value, _ := arg.Expression.ValueExpression.Value.ToValueObject()
		require.Equal(t, "target", key)

		object, _ := value.ToHashObject()
		host, _ := object["host"].ToString()
		user, _ := object["user"].ToString()
		require.Equal(t, "\"yuuka.natsuneko.net\"", *host)
		require.Equal(t, "\"ubuntu\"", *user)
	}

	repr.Println(ret)
}
