package otsukai

import (
	require "github.com/alecthomas/assert/v2"
	"os"
	"testing"
)

func c(t *testing.T, path string) string {
	content, err := os.ReadFile(path)
	require.NoError(t, err)

	return string(content) + "\n"
}

func TestExample_DockerCompose(t *testing.T) {
	content := c(t, "./examples/docker-compose/otsukai.rb")
	ret, err := Parser.ParseString("", content)
	require.NoError(t, err)
	require.Equal(t, 4, len(ret.Statements))

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

	// task :deploy do ... end
	{
		invocation := ret.Statements[2].Statement.ExpressionStatement.Expression.InvocationExpression
		require.Equal(t, "task", invocation.Expression.IdentifierNameExpression.Identifier.Name)

		args := invocation.ArgumentList.Argument
		lambda := invocation.ArgumentList.LambdaExpression

		arg := args[0]
		require.Equal(t, "deploy", arg.Expression.ValueExpression.Value.HashSymbol.Identifier)
		require.NotEqual(t, nil, lambda)

		statements := lambda.Statements
		require.Equal(t, 1, len(statements))

		statement := statements[0].Statement.IfStatement
		require.NotEqual(t, nil, statement)

		// if
		{
			condition := statement.Condition.InvocationExpression
			require.Equal(t, "changed", condition.Expression.IdentifierNameExpression.Identifier.Name)

			args = condition.ArgumentList.Argument
			require.Equal(t, 2, len(args))

			a := args[0]
			require.Equal(t, "\"/path/to/docker-compose.yml\"", *a.Expression.ValueExpression.Value.Literal.String)

			b := args[1]
			require.Equal(t, "from", *b.Identifier)
			require.Equal(t, "last_commit", b.Expression.ValueExpression.Value.HashSymbol.Identifier)

			// run_as
			{
				statements = statement.Statements
				require.Equal(t, 1, len(statements))

				statement := statements[0].Statement
				invocation = statement.ExpressionStatement.Expression.InvocationExpression
				require.NotEqual(t, nil, invocation)
				require.Equal(t, "run_as", invocation.Expression.IdentifierNameExpression.Identifier.Name)

				args = invocation.ArgumentList.Argument
				require.Equal(t, 1, len(args))
				require.Equal(t, "sudo", args[0].Expression.ValueExpression.Value.HashSymbol.Identifier)

				lambda = invocation.ArgumentList.LambdaExpression
				require.NotEqual(t, nil, lambda)

				// do ...
				{
					statements = lambda.Statements
					require.Equal(t, 3, len(statements))
				}
			}

		}
	}

	// hook after: :deploy do ... end
	{
		invocation := ret.Statements[3].Statement.ExpressionStatement.Expression.InvocationExpression
		require.Equal(t, "hook", invocation.Expression.IdentifierNameExpression.Identifier.Name)

		args := invocation.ArgumentList.Argument
		lambda := invocation.ArgumentList.LambdaExpression
		require.Equal(t, 1, len(args))
		require.NotEqual(t, nil, lambda)

		statement := lambda.Statements[0].Statement.IfStatement
		require.NotEqual(t, nil, statement)

		condition := statement.Condition
		statements := statement.Statements

		invocation = condition.InvocationExpression
		require.Equal(t, "task_success", invocation.Expression.IdentifierNameExpression.Identifier.Name)
		require.Equal(t, 0, len(invocation.ArgumentList.Argument))
		require.Equal(t, 1, len(statements))
	}

	// epr.Println(ret)
}
