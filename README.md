:x: loc
=======
[![Build Status][build badge]][travis result]
[![Windows Build status][windows build badge]][appveyor result]
[![Coverage Status][coverage status]][coverage result]

[loc][loc document] is a small library to make a nice-looking error with source location information.
It provides a struct to represent a source file, a specific position in code and an error related to
specific range or position in source.

This library is useful to provide a unified look for error messages of compilers, interpreters or
translators.

By using `Source` and `Position` types as position information, this library provides an error
type which provides nice look error message.

- It shows the code snippet which caused an error
- Enable to add notes to error by nesting an error instance like [pkg/errors](https://github.com/pkg/errors)
- Proper location is automatically added to error messages and notes
- Colorized label like 'Error:' or 'Note:'

It's important to make a good error when compilation or execution errors found. loc helps it.
This library is actually used in some my compiler implementation.

```go
package main

import (
	"fmt"
	"github.com/rhysd/loc"
)

func main() {
	// At first you should gain entire source as *Source instance.

	code :=
		`package main

func main() {
	blah := 42

	blah := true
}`
	src := loc.NewDummySource(code)

	// You can get *Source instance from file (NewSourceFromFile) or stdin (NewSourceFromStdin) also.

	// Let's say to find an error at some range in the source.

	start := loc.Pos{
		Offset: 41,
		Line:   6,
		Column: 1,
		File:   src,
	}
	end := loc.Pos{
		Offset: 54,
		Line:   6,
		Column: 12,
		File:   src,
	}

	// NewError or other factory functions make a new error instance with the range. Error instance implements
	// error interface so it can be handled like other error types.

	err := loc.NewError(start, end, "Found duplicate symbol 'foo'")

	// Assume that you find additional information (location of variable and its type). Then you can add some
	// notes to the error. Notes can be added by wrapping errors like pkg/errors library.

	prev := loc.Pos{
		Offset: 26,
		Line:   4,
		Column: 1,
		File:   src,
	}

	err = err.NoteAt(prev, "Defined here at first")
	err = err.NoteAt(prev, "Previously defined as int")

	// Finally you can see the result!

	fmt.Println(err)
}
```

Above code should show the following output:

```
Error: Found duplicate symbol 'foo' (at <dummy>:6:1)
    Note: Defined here at first (at <dummy>:4:1)
    Note: Previously defined as int (at <dummy>:4:1)

>       blah := true
```

<img src="https://github.com/rhysd/ss/blob/master/loc/output.png?raw=true" width="371" alt="output screenshot"/>

Labels such as 'Error:' or 'Notes:' are colorized. Main error message is emphasized with bold font.
And source code location information (file name, line and column) is added with gray text.
If the error has range information, the error shows code snippet which caused the error at the end
of error message

Please see [documentation][loc document] to know whole APIs.

Note that on Windows always color is disabled because ANSI color sequence is not availale on CMD.exe.

[loc document]: https://godoc.org/github.com/rhysd/loc
[build badge]: https://travis-ci.org/rhysd/loc.svg?branch=master
[travis result]: https://travis-ci.org/rhysd/loc
[coverage status]: https://codecov.io/gh/rhysd/loc/branch/master/graph/badge.svg
[coverage result]: https://codecov.io/gh/rhysd/loc
[windows uild badge]: https://ci.appveyor.com/api/projects/status/4d3bkiabf088gboi/branch/master?svg=true
[appveyor result]: https://ci.appveyor.com/project/rhysd/loc/branch/master
