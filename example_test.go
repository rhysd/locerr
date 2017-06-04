package loc

import (
	"fmt"
	"os"
	"testing"
)

func TestExample(t *testing.T) {
	// At first you should gain entire source as *Source instance.

	code :=
		`package main

func main() {
	foo := 42

	foo := true
}`
	src := NewDummySource(code)

	// You can get *Source instance from file (NewSourceFromFile) or stdin (NewSourceFromStdin) also.

	// Let's say to find an error at some range in the source.

	start := Pos{
		Offset: 41,
		Line:   6,
		Column: 1,
		File:   src,
	}
	end := Pos{
		Offset: 54,
		Line:   6,
		Column: 12,
		File:   src,
	}

	// NewError or other factory functions make a new error instance with the range. Error instance implements
	// error interface so it can be handled like other error types.

	err := ErrorIn(start, end, "Found duplicate symbol 'foo'")

	// Assume that you find additional information (location of variable and its type). Then you can add some
	// notes to the error. Notes can be added by wrapping errors like pkg/errors library.

	prev := Pos{
		Offset: 26,
		Line:   4,
		Column: 1,
		File:   src,
	}

	err = err.NoteAt(prev, "Defined here at first")
	err = err.NoteAt(prev, "Previously defined as int")

	// Finally you can see the result!

	// Get the error message as string. Note that this is only for non-Windows OS.
	fmt.Println(err)
	// Output:
	// Error: Found duplicate symbol 'foo' (at <dummy>:6:1)
	//     Note: Defined here at first (at <dummy>:4:1)
	//     Note: Previously defined as int (at <dummy>:4:1)
	//
	// >       foo := true

	// Directly writes the error message into given file.
	// This supports Windows. Useful to output from stdout or stderr.
	err.PrintToFile(os.Stdout)
	// Output:
	// Error: Found duplicate symbol 'foo' (at <dummy>:6:1)
	//     Note: Defined here at first (at <dummy>:4:1)
	//     Note: Previously defined as int (at <dummy>:4:1)
	//
	// >       foo := true
}
