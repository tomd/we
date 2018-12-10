# we

"wrapped error" - Yet another package to ease error handling.
Primary purpose is to make a log line that shows the call graph.

Any time you `return err` do `return we.New(err)` instead.

## Usage

`we.New(e)` => `"pkg.func(): e.Error()"`

`we.New(e, 42)` => `"pkg.func(42): e.Error()"`

`we.Newf(e, "foo=%d", 42)` => `"pkg.func(foo=42): e.Error()"`

## License

2018, Tom Spangebu <tom@pogostick.net>

CC0, see LICENSE/COPYING
