package eval

import "fmt"

type RuntimeError struct {
	Message string
}

func RuntimeErrorf(format string, args ...any) RuntimeError {
	return RuntimeError{Message: fmt.Sprintf(format, args...)}
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("Runtime Error: %s", e.Message)
}

type TypeError struct {
	Message string
}

func TypeErrorf(format string, args ...any) TypeError {
	return TypeError{Message: fmt.Sprintf(format, args...)}
}

func (e TypeError) Error() string {
	return fmt.Sprintf("Type Error: %s", e.Message)
}

type InternalError struct {
	Message string
}

func InternalErrorf(format string, args ...any) InternalError {
	return InternalError{Message: fmt.Sprintf(format, args...)}
}

func (e InternalError) Error() string {
	return fmt.Sprintf("Internal Error: %s", e.Message)
}
