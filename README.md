:x: locerr
==========
[![Build Status][build badge]][travis result]
[![Windows Build status][windows build badge]][appveyor result]
[![Coverage Status][coverage status]][coverage result]
[![GoDoc][godoc badge]][locerr document]

[locerr][locerr document] is a small library to make a nice-looking locational error in a source code.
It provides a struct to represent a source file, a specific position in code and an error related to
specific range or position in source.

This library is useful to provide a unified look for error messages raised by compilers, interpreters
or translators.

By using `locerr.Source` and `locerr.Position` types as position information, this library can provide
an error type which shows nice look error message.

- It shows the code snippet which caused an error
- Enable to add notes to error by nesting an error instance like [pkg/errors](https://github.com/pkg/errors)
- Proper location is automatically added to error messages and notes
- Colorized label like 'Error:' or 'Note:'
- Windows is supported

It's important to make a good error when compilation or execution errors found. [locerr][locerr document]
helps it. This library is actually used in some my compiler implementation.

## Installation

Please use `go get`.

```console
$ go get -u github.com/rhysd/locerr
```

## Usage

As example, let's say to make a locational error for following pseudo code. In this code, function
`foo` is defined with 1 parameter but called with 3 parameters.

```
function foo(x: bool): int {
  return (if x then 42 else 21)
}

function main() {
  foo(true,
      42,
      "test")
}
```

We can make a locational error with some notes using locerr as following.

```go
package main

import (
	"fmt"
	"os"

	"github.com/rhysd/locerr"
)

func main() {
	// At first you should gain entire source as *locerr.Source instance.

	code :=
		`function foo(x: bool): int {
  return (if x then 42 else 21)
}

function main() {
  foo(true,
      42,
      "test")
}`
	src := locerr.NewDummySource(code)

	// You can get *locerr.Source instance from file (NewSourceFromFile) or stdin (NewSourceFromStdin) also.

	// Let's say to find an error at some range in the source. 'start' indicates the head of the first argument.
    // 'end' indicates the end of the last argument.

	start := locerr.Pos{
		Offset: 88,
		Line:   6,
		Column: 7,
		File:   src,
	}
	end := locerr.Pos{
		Offset: 116,
		Line:   9,
		Column: 12,
		File:   src,
	}

	// NewError or other factory functions make a new error instance with the range. locerr.Error instance
	// implements error interface so it can be handled like other error types.

	err := locerr.ErrorIn(start, end, "Calling 'foo' with wrong number of argument")

	// Assume that you find additional information (location of variable and its type). Then you can add some
	// notes to the error. Notes can be added by wrapping errors like pkg/errors library.

	prev := locerr.Pos{
		Offset: 9,
		Line:   1,
		Column: 10,
		File:   src,
	}

	err = err.NoteAt(prev, "Defined with 1 parameter")
	err = err.NoteAt(prev, "'foo' was defined as 'bool -> int'")

	// Finally you can see the result!

	// Get the error message as string. Note that this is only for non-Windows OS.
	fmt.Println(err)

	// Directly writes the error message into given file.
	// This supports Windows. Useful to output from stdout or stderr.
	err.PrintToFile(os.Stdout)
}
```

Above code should show the following output:

```
Error: Calling 'foo' with wrong number of argument (at <dummy>:6:7)
  Note: Defined with 1 parameter (at <dummy>:1:10)
  Note: 'foo' was defined as 'bool -> int' (at <dummy>:1:10)

>   foo(true,
>       42,
>       "test")

```

<img src="https://github.com/rhysd/ss/blob/master/locerr/output.png?raw=true" width="547" alt="output screenshot"/>

Labels such as 'Error:' or 'Notes:' are colorized. Main error message is emphasized with bold font.
And source code location information (file name, line and column) is added with gray text.
If the error has range information, the error shows code snippet which caused the error at the end
of error message.

If you have only one position information rather than two, 'start' position and 'end' position,
`ErrorAt` is available instead of `ErrorIn`. `ErrorAt` takes one `Pos` instance.

```go
err := locerr.ErrorAt(start, "Calling 'foo' with wrong number of argument")
```

In this case, line snippet is shown in error message. `pos.Line` is used to get line from source text.
`fmt.Println(err)` will show the following.

```
Error: Calling 'foo' with wrong number of argument (at <dummy>:6:7)

>   foo(true,

```


## Development

### How to run tests

```console
$ go test ./
```

Note that `go test -v` may fail because color sequences are not assumed in tests.

### How to run fuzzing test

Fuzzing test using [go-fuzz][].

```console
$ cd ./fuzz
$ go-fuzz-build github.com/rhysd/locerr/fuzz
$ go-fuzz -bin=./locerr_fuzz-fuzz.zip -workdir=fuzz
```

Last command starts fuzzing tests until stopped with `^C`. Every 3 seconds it reports the current
result. It makes 3 directories in `fuzz` directory as the result, `corpus`, `crashers` and
`suppressions`. `crashers` contains the information about the crash caused by fuzzing.

[locerr document]: https://godoc.org/github.com/rhysd/locerr
[build badge]: https://travis-ci.org/rhysd/locerr.svg?branch=master
[travis result]: https://travis-ci.org/rhysd/locerr
[coverage status]: https://codecov.io/gh/rhysd/locerr/branch/master/graph/badge.svg
[coverage result]: https://codecov.io/gh/rhysd/locerr
[windows build badge]: https://ci.appveyor.com/api/projects/status/v4ghlgka6e6st2mn/branch/master?svg=true
[appveyor result]: https://ci.appveyor.com/project/rhysd/locerr/branch/master
[godoc badge]: https://godoc.org/github.com/rhysd/locerr?status.svg
[go-fuzz]: https://github.com/dvyukov/go-fuzz
