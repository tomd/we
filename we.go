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
		panic("we.New(nil)")
	}
	res := new(wrapped_error)
	// If we wrap a wrapped_error then copy the cause
	// and exit code.  So cause will always be the
	// innermost error (which cannot be a wrapped_error).
	switch e := e.(type) {
	case *wrapped_error:
		res.cause = e.cause
		res.code = e.code
	default:
		res.cause = e
		res.code = DefaultExitCode
	}
	// Find our caller.
	var pc [1]uintptr
	n := runtime.Callers(2, pc[:]) // skip 2: we.New & runtime.Callers
	if n != 1 {
		panic(fmt.Sprintf("we.New(%v): runtime.Callers() == %d", e, n))
	}
	frame, _ := runtime.CallersFrames(pc[:]).Next()
	funcname := frame.Function
	if !MainPrefix && strings.HasPrefix(funcname, "main.") {
		funcname = funcname[5:]
	}
	// Format it.
	args_str := ""
	if len(args_format_and_args) > 0 {
		format := args_format_and_args[0].(string)
		args_str = fmt.Sprintf(format, args_format_and_args[1:]...)
	}
	res.msg = fmt.Sprintf("%s(%s): %s", funcname, args_str, e.Error())
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
		panic("we.WithExitCode(nil)")
	}
	res := new(wrapped_error)
	switch e := e.(type) {
	case *wrapped_error:
		res.cause = e.cause
		res.msg = e.msg
	default:
		res.cause = e
		res.msg = e.Error()
	}
	res.code = code
	return res
}
