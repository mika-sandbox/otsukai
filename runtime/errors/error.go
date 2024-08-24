package error

import "errors"

var CAST_ERROR = errors.New("failed to cast value.go")
var SYNTAX_ERROR = errors.New("syntax error")
var RUNTIME_ERROR = errors.New("runtime error")
