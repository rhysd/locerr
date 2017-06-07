/*
Package locerr is a small library to make an error with source code location information.
It provides a struct to represent a source file, a specific position in
code and an error related to specific range or position in source.

It's important to make a good error when compilation or execution errors found. locerr helps it.
This library is actually used in some my compiler implementation.

Repository: https://github.com/rhysd/locerr

At first you should gain entire source as *Source instance.

    code :=
    	`package main

    func main() {
    	foo := 42

    	foo := true
    }
    `
    src := locerr.NewDummySource(code)

You can get *Source instance from file (NewSourceFromFile) or stdin (NewSourceFromStdin) also.

Let's say to find an error at some range in the source.

    start := locerr.Pos{
    	Offset: 41,
    	Line:   6,
    	Column: 2,
    	File:   src,
    }
    end := locerr.Pos{
    	Offset: 52,
    	Line:   6,
    	Column: 12,
    	File:   src,
    }

ErrorIn or other factory functions make a new error instance with the range. Error instance implements
error interface so it can be handled like other error types.

    err := locerr.ErrorIn(start, end, "Found duplicate symbol 'foo'")

Assume that you find additional information (location of variable and its type). Then you can add some
notes to the error. Notes can be added by wrapping errors like pkg/errors library.

    prev := locerr.Pos{
    	Offset: 26,
    	Line:   4,
    	Column: 1,
    	File:   src,
    }

    err = err.NoteAt(prev, "Defined here at first")
    err = err.NoteAt(prev, "Previously defined as int")

Finally you can see the result! err.Error() gets the error message as string. Note that this is only for
non-Windows OS.

    fmt.Println(err)

It should output following:

    Error: Found duplicate symbol 'foo' (at <dummy>:6:1)
      Note: Defined here at first (at <dummy>:4:1)
      Note: Previously defined as int (at <dummy>:4:1)

    >       foo := true


To support Windows, please use PrintToFile() method. It directly writes the error message into given file.
This supports Windows and is useful to output from stdout or stderr.

    err.PrintToFile(os.Stderr)

Labels such as 'Error:' or 'Notes:' are colorized. Main error message is emphasized with bold font.
And source code location information (file name, line and column) is added with gray text.
If the error has range information, the error shows code snippet which caused the error at the end
of error message

Colorized output can be seen at https://github.com/rhysd/ss/blob/master/locerr/output.png?raw=true

If you have only one position information rather than two, 'start' position and 'end' position,
ErrorAt() is available instead of ErrorIn() ErrorAt() takes one Pos instance.

    err = ErrorAt(start, "Calling 'foo' with wrong number of argument")

In this case, line snippet is shown in error message. `pos.Line` is used to get line from source text.

    fmt.Println(err)

It should output following:

    Output:
    Error: Calling 'foo' with wrong number of argument (at <dummy>:6:7)

    >   foo(true,


*/
package locerr
