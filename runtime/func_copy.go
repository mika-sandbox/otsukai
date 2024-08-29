package runtime

import (
	"otsukai"
	"otsukai/parser"
	"otsukai/runtime/context"
	re "otsukai/runtime/errors"
	"otsukai/runtime/helpers"
	"otsukai/runtime/value"
)

func InvokeCopy(ctx context.IContext, arguments []parser.Argument) (value.IValueObject, error) {
	remote := helpers.GetStringLiteral(helpers.GetNamedArgument(arguments, "remote"))
	if remote == nil {
		otsukai.Errf("remote named argument must be required")
		return nil, re.RUNTIME_ERROR
	}

	local := helpers.GetStringLiteral(helpers.GetNamedArgument(arguments, "local"))
	if local == nil {
		otsukai.Errf("local named argument must be required")
		return nil, re.RUNTIME_ERROR
	}

	to := helpers.GetSymbol(helpers.GetNamedArgument(arguments, "to"))
	if to == nil {
		otsukai.Errf("to named argument must be required")
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
