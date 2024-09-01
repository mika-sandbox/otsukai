package runtime

import (
	"github.com/mika-sandbox/otsukai/runtime/context"
	"github.com/mika-sandbox/otsukai/runtime/value"
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
