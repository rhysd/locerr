:x: locerr
==========
[![Build Status][build badge]][travis result]
[![Windows Build status][windows build badge]][appveyor result]
[![Coverage Status][coverage status]][coverage result]

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

```go
package main

import (
	"fmt"
	"github.com/rhysd/locerr"
	"os"
)

func main() {
	// At first you should gain entire source as *locerr.Source instance.

	code :=
		`package main

func main() {
	foo := 42

	foo := true
}`
	src := locerr.NewDummySource(code)

	// You can get *locerr.Source instance from file (locerr.NewSourceFromFile)
	// or stdin (locerr.NewSourceFromStdin) also.

	// Let's say to find an error at some range in the source.

	start := locerr.Pos{
		Offset: 41,
		Line:   6,
		Column: 1,
		File:   src,
	}
	end := locerr.Pos{
		Offset: 52,
		Line:   6,
		Column: 12,
		File:   src,
	}

	// ErrorIn or other factory functions make a new error instance with the range. Error instance implements
	// error interface so it can be handled like other error types.

	err := locerr.ErrorIn(start, end, "Found duplicate symbol 'foo'")

	// Assume that you find additional information (location of variable and its type). Then you can add some
	// notes to the error. Notes can be added by wrapping errors like pkg/errors library.

	prev := locerr.Pos{
		Offset: 26,
		Line:   4,
		Column: 1,
		File:   src,
	}

	err = err.NoteAt(prev, "Defined here at first")
	err = err.NoteAt(prev, "Previously defined as int")

	// Finally you can see the result!

	// Get the error message as string. Note that this is only for non-Windows OS.
	msg := err.Error()
	fmt.Println(msg)

	// Directly writes the error message into given file.
	// This supports Windows. Useful to output from stdout or stderr.
	err.PrintToFile(os.Stderr)
}
```

Above code should show the following output:

```
Error: Found duplicate symbol 'foo' (at <dummy>:6:1)
    Note: Defined here at first (at <dummy>:4:1)
    Note: Previously defined as int (at <dummy>:4:1)

>       foo := true

```

<img src="https://github.com/rhysd/ss/blob/master/loc/output.png?raw=true" width="371" alt="output screenshot"/>

Labels such as 'Error:' or 'Notes:' are colorized. Main error message is emphasized with bold font.
And source code location information (file name, line and column) is added with gray text.
If the error has range information, the error shows code snippet which caused the error at the end
of error message

Please see [documentation][locerr document] to know whole APIs.

[locerr document]: https://godoc.org/github.com/rhysd/locerr
[build badge]: https://travis-ci.org/rhysd/locerr.svg?branch=master
[travis result]: https://travis-ci.org/rhysd/locerr
[coverage status]: https://codecov.io/gh/rhysd/locerr/branch/master/graph/badge.svg
[coverage result]: https://codecov.io/gh/rhysd/locerr
[windows build badge]: https://ci.appveyor.com/api/projects/status/v4ghlgka6e6st2mn/branch/master?svg=true
[appveyor result]: https://ci.appveyor.com/project/rhysd/locerr/branch/master
