// "we" wraps an existing error, prepending a call stack to the message
// and/or adds an exit code.  Basic use is `return we.New(err)`.
//
//	we.New(e) => "pkg.func(): e.Error()"
//
//	we.New(e, 42) => "pkg.func(42): e.Error()"
//
//	we.Newf(e, "foo=%d", 42) => "pkg.func(foo=42): e.Error()"
//
// "we" can also keep track of an exit code the caller may
// extract with we.ExitCode():
//
//	we.ExitCode(we.WithExitCode(42, e)) == 42
//
// Note for New/Newf/NewEC/NewfEC/WithExitCode:
// if given a wrapped_error these functions will actually just mutate
// the wrapped_error and return it.
//
package we

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// MainPrefix: set to true to keep "main." package prefix.
var MainPrefix bool = false

// DefaultExitCode: exit code when WithExitCode() not used.
var DefaultExitCode int = 1

// wrapped_error is the error type of "we".
type wrapped_error struct {
	msg   string
	cause error
	code  int
}

// wrapped_error implements the error interface.
func (self *wrapped_error) Error() string {
	return self.msg
}

// Cause returns the original error.
func Cause(e error) error {
	if e, ok := e.(*wrapped_error); ok {
		return e.cause
	}
	return e
}

// new_newf is the implementation of New and Newf.
func new_newf(f bool, e error, args ...interface{}) error {
	if e == nil {
		return nil
	}

	funcname := caller(4) // skip 4: New[f][EC](), new_newf(), caller() and runtime.Callers()
	if !MainPrefix && strings.HasPrefix(funcname, "main.") {
		funcname = funcname[5:]
	}

	args_str := ""
	if len(args) > 0 {
		if f {
			fmt_str := args[0].(string)
			args_str = fmt.Sprintf(fmt_str, args[1:]...)
		} else {
			fmt_str := strings.Repeat("%v,", len(args))
			args_str = fmt.Sprintf(fmt_str[:len(fmt_str)-1], args...)
		}
	}
	msg := fmt.Sprintf("%s(%s): %s", funcname, args_str, e.Error())

	if e, ok := e.(*wrapped_error); ok {
		e.msg = msg
		return e
	}
	res := new(wrapped_error)
	res.msg = msg
	res.cause = e
	res.code = DefaultExitCode
	return res
}

// New create a new wrapped_error with the given arguments.
func New(e error, args ...interface{}) error {
	return new_newf(false, e, args...)
}

// Newf create a new wrapped_error with the given format and arguments.
func Newf(e error, format_and_args ...interface{}) error {
	return new_newf(true, e, format_and_args...)
}

// NewEC == New + WithExitCode
func NewEC(code int, e error, args ...interface{}) error {
	res := new_newf(false, e, args...)
	return WithExitCode(code, res)
}

// NewfEC == Newf + WithExitCode
func NewfEC(code int, e error, format_and_args ...interface{}) error {
	res := new_newf(true, e, format_and_args...)
	return WithExitCode(code, res)
}

// ExitCode extracts the exit code if e is a wrapped_error, otherwise returns DefaultExitCode.
func ExitCode(e error) int {
	if e, ok := e.(*wrapped_error); ok {
		return e.code
	}
	return DefaultExitCode
}

// WithExitCode create a new wrapped_error with the given exit code.
func WithExitCode(code int, e error) error {
	if e == nil {
		return nil
	}
	if e, ok := e.(*wrapped_error); ok {
		e.code = code
		return e
	}
	res := new(wrapped_error)
	res.msg = e.Error()
	res.cause = e
	res.code = code
	return res
}

// callers returns the name of the function "skip" frames above runtime.Callers.
func caller(skip int) string {
	var pc [1]uintptr
	n := runtime.Callers(skip, pc[:])
	if n != 1 {
		panic(fmt.Sprintf("we.caller(): runtime.Callers() == %d", n))
	}
	frame, _ := runtime.CallersFrames(pc[:]).Next()
	return frame.Function
}

// Errorf == errors.New / fmt.Errorf.  When you don't want to import errors
// or fmt just to create a new error.
func Errorf(format string, a ...interface{}) error {
	if len(a) == 0 {
		return errors.New(format)
	}
	return fmt.Errorf(format, a...)
}
