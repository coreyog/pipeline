package pipeline

import (
	"errors"
)

var ErrNotAFunction = errors.New("not a function")
var ErrParameterMismatch = errors.New("previous functions outputs don't match next functions inputs")
