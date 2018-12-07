//
// To wrap an error:
//	we.New(err)
//		=> "package.funcname(): err"
// or
//	we.New(err, "foo=%d,bar=%s", 42, "plugh")
//		=> "package.funcname(foo=42,bar=plugh): err"
//
// Typical use:
//	if err != nil {
//		return we.New(err)
//	}
// or
//	if !ok {
//		return we.New(fmt.Errorf("not ok: %v", xyzzy))
//	}
//
// Also the error site may indicate a preferred exit code:
//	return we.WithExitCode(29, err)
//
// Which may be recovered with:
//	os.Exit(we.ExitCode(err))
//
package we

import (
	"fmt"
	"runtime"
	"strings"
)

// Set to true to keep "main." package prefix.
var MainPrefix bool = false

// Default exit code when WithExitCode() not used.
var DefaultExitCode int = 1

type wrapped_error struct {
	msg   string
	cause error
	code  int
}

func (self *wrapped_error) Error() string {
	return self.msg
}

func Cause(e error) error {
	if e, ok := e.(*wrapped_error); ok {
		return e.cause
	}
	return e
}

func New(e error, args_format_and_args ...interface{}) error {
	if e == nil {
		return nil
	}

	funcname := caller(3) // skip 3: we.New(), we.caller() and runtime.Callers()
	if !MainPrefix && strings.HasPrefix(funcname, "main.") {
		funcname = funcname[5:]
	}

	args_str := ""
	if len(args_format_and_args) > 0 {
		format := args_format_and_args[0].(string)
		args_str = fmt.Sprintf(format, args_format_and_args[1:]...)
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

func ExitCode(e error) int {
	if e, ok := e.(*wrapped_error); ok {
		return e.code
	}
	return DefaultExitCode
}

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

func caller(skip int) string {
	var pc [1]uintptr
	n := runtime.Callers(skip, pc[:])
	if n != 1 {
		panic(fmt.Sprintf("we.caller(): runtime.Callers() == %d", n))
	}
	frame, _ := runtime.CallersFrames(pc[:]).Next()
	return frame.Function
}
