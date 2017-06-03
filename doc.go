/*
Package loc is a small library to make an error with source code location information.
It provides a struct to represent a source file, a specific position in
code and an error related to specific range or position in source.

It's important to make a good error when compilation or execution errors found. loc helps it.
This library is actually used in some my compiler implementation.

Repository: https://github.com/rhysd/loc

At first you should gain entire source as *Source instance.

    code :=
    	`package main

    func main() {
    	blah := 42

    	blah := true
    }
    `
    src := loc.NewDummySource(code)

You can get *Source instance from file (NewSourceFromFile) or stdin (NewSourceFromStdin) also.

Let's say to find an error at some range in the source.

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

ErrorIn or other factory functions make a new error instance with the range. Error instance implements
error interface so it can be handled like other error types.

    err := loc.ErrorIn(start, end, "Found duplicate symbol 'foo'")

Assume that you find additional information (location of variable and its type). Then you can add some
notes to the error. Notes can be added by wrapping errors like pkg/errors library.

    prev := loc.Pos{
    	Offset: 26,
    	Line:   4,
    	Column: 1,
    	File:   src,
    }

    err = err.NoteAt(prev, "Defined here at first")
    err = err.NoteAt(prev, "Previously defined as int")

Finally you can see the result!

    fmt.Println(err)
    // Output:
    // Error: Found duplicate symbol 'foo' (at <dummy>:6:1)
    //   Note: Defined here at first (at <dummy>:4:1)
    //   Note: Previously defined as int (at <dummy>:4:1)
    //
    // >       blah := true
    //

Labels such as 'Error:' or 'Notes:' are colorized. Main error message is emphasized with bold font.
And source code location information (file name, line and column) is added with gray text.
If the error has range information, the error shows code snippet which caused the error at the end
of error message

Colorized output can be seen at https://github.com/rhysd/ss/blob/master/loc/output.png?raw=true
*/
package loc
