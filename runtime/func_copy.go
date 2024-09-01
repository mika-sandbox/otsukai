package runtime

import (
	"github.com/mika-sandbox/otsukai/logger"
	"github.com/mika-sandbox/otsukai/parser"
	"github.com/mika-sandbox/otsukai/runtime/context"
	re "github.com/mika-sandbox/otsukai/runtime/errors"
	"github.com/mika-sandbox/otsukai/runtime/helpers"
	"github.com/mika-sandbox/otsukai/runtime/value"
)

func InvokeCopy(ctx context.IContext, arguments []parser.Argument) (value.IValueObject, error) {
	remote := helpers.GetStringLiteral(helpers.GetNamedArgument(arguments, "remote"))
	if remote == nil {
		logger.Errf("remote named argument must be required")
		return nil, re.RUNTIME_ERROR
	}

	local := helpers.GetStringLiteral(helpers.GetNamedArgument(arguments, "local"))
	if local == nil {
		logger.Errf("local named argument must be required")
		return nil, re.RUNTIME_ERROR
	}

	to := helpers.GetSymbol(helpers.GetNamedArgument(arguments, "to"))
	if to == nil {
		logger.Errf("to named argument must be required")
		return nil, re.RUNTIME_ERROR
	}

	session := ctx.GetRemoteSession()
	if *to == "remote" {
		err := session.CopyToRemote(*local, *remote)
		if err != nil {
			return nil, re.RUNTIME_ERROR
		}
	} else if *to == "local" {
		isDir := helpers.GetBoolLiteral(helpers.GetNamedArgument(arguments, "is_dir"), false)

		err := session.CopyToLocal(*local, *remote, isDir)
		if err != nil {
			return nil, re.RUNTIME_ERROR
		}
	} else {
		return nil, re.RUNTIME_ERROR
	}

	return nil, nil
}
