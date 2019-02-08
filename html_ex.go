package blackfriday

import (
	"fmt"
	"io"
)

// FunctionHandler handles a markdown function.
type FunctionHandler func(w io.Writer, name string, args []interface{}) error

// DefaultFunctionHandler is a default implementation of the FunctionHandler.
func DefaultFunctionHandler(w io.Writer, name string, args []interface{}) error {
	return fmt.Errorf("function %s is not defined", name)
}
