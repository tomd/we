# we

"wrapped error" - Yet another package to ease error handling.
Primary purpose is to make a log line that shows the call graph.

Any time you `return err;` do `return we.New(err);` instead.

## Usage

`we.New(err)` => `"pkg.funcname(): err"`

`we.New(err, "%d", 42)` => `"pkg.funcname(42): err"`

## License

2018, Tom Spangebu <tom@pogostick.net>
CC0, see LICENSE/COPYING
