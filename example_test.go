package locerr

import (
	"fmt"
	"os"
)

const code = `function foo(x: bool): int {
  return (if x then 42 else 21)
}

function main() {
  foo(true,
      42,
      "test")
}`

func ExampleErrorWithRange() {
	// At first you should gain entire source as *Source instance.

	src := NewDummySource(code)

	// You can get *Source instance from file (NewSourceFromFile) or stdin (NewSourceFromStdin) also.

	// Let's say to find an error at some range in the source.

	start := Pos{
		Offset: 88,
		Line:   6,
		Column: 7,
		File:   src,
	}
	end := Pos{
		Offset: 116,
		Line:   9,
		Column: 12,
		File:   src,
	}

	// NewError or other factory functions make a new error instance with the range. Error instance implements
	// error interface so it can be handled like other error types.

	err := ErrorIn(start, end, "Calling 'foo' with wrong number of argument")

	// Assume that you find additional information (location of variable and its type). Then you can add some
	// notes to the error. Notes can be added by wrapping errors like pkg/errors library.

	prev := Pos{
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

func ExampleErrorWithOnePos() {
	src := NewDummySource(code)

	pos := Pos{
		Offset: 88,
		Line:   6,
		Column: 7,
		File:   src,
	}

	// If you have only one position information rather than two, 'start' position and 'end' position,
	// ErrorAt() is available instead of ErrorIn() ErrorAt() takes one Pos instance.
	err := ErrorAt(pos, "Calling 'foo' with wrong number of argument")

	// In this case, line snippet is shown in error message. `pos.Line` is used to get line from source text.
	fmt.Println(err)
}
