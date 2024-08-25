package helpers

import "otsukai/parser"

func GetNamedArgument(arguments []parser.Argument, name string) *parser.Argument {
	for _, argument := range arguments {
		identifier := argument.Identifier
		if identifier != nil && *identifier == name {
			return &argument
		}
	}

	return nil
}
