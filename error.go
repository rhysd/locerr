package locerr

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

// SetColor controls font should be colorful or not.
func SetColor(enabled bool) {
	color.NoColor = !enabled
}

var (
	bold     = color.New(color.Bold)
	red      = color.New(color.FgRed)
	green    = color.New(color.FgGreen)
	gray     = color.New(color.FgHiBlack)
	emphasis = color.New(color.FgHiGreen, color.Bold, color.Underline)
)

// Error represents a compilation error with positional information and stacked messages.
type Error struct {
	Start    Pos
	End      Pos
	Messages []string
}

func writeSnipLine(w io.Writer, line string) {
	indent, len := 0, len(line)
	for indent < len {
		if line[indent] != ' ' && line[indent] != '\t' {
			break
		}
		indent++
	}
	if indent != 0 {
		// Write indent without emphasis
		fmt.Fprint(w, line[:indent])
	}
	if indent != len {
		// Write code snip with emphasis
		emphasis.Fprint(w, line[indent:])
	}
}

func (err *Error) writeSnip(w io.Writer) {
	fmt.Fprint(w, "\n\n> ")

	code := err.Start.File.Code
	start := err.Start.Offset
	for start-1 >= 0 {
		if code[start-1] == '\n' {
			break
		}
		start--
	}
	if start < err.Start.Offset {
		// Write code before snip in first line
		w.Write(code[start:err.Start.Offset])
	}

	lines := strings.Split(string(code[err.Start.Offset:err.End.Offset]), "\n")

	// First line does not have "> " prefix
	writeSnipLine(w, lines[0])

	for _, line := range lines[1:] {
		fmt.Fprint(w, "\n> ")
		writeSnipLine(w, line)
	}

	end := err.End.Offset
	len := len(code)
	for end < len {
		if code[end] == '\n' {
			break
		}
		end++
	}
	if err.End.Offset < end {
		// Write code after snip in last line
		w.Write(code[err.End.Offset:end])
	}

	fmt.Fprint(w, "\n")

	// TODO:
	// If the code snippet for the token is too long, skip lines with '...' except for starting N lines
	// and ending N lines
}

func lineStartOffset(code []byte, lnum int) int {
	l := 1
	for i, r := range code {
		if l == lnum {
			return i
		}
		if r == '\n' {
			l++
		}
	}
	return -1
}

// Show line based on err.Start.Line. We don't use offset for this because some environment offset
// cannot be obtained (e.g. getting location from runtime.Caller).
func (err *Error) writeOnelineSnip(w io.Writer) {
	code := err.Start.File.Code
	len := len(code)
	if len == 0 {
		return
	}

	start := lineStartOffset(code, err.Start.Line)
	if start == -1 {
		return
	}

	end := start
	for end < len {
		if code[end] == '\n' {
			break
		}
		end++
	}

	if start == end {
		// Snippet is empty. Skipped.
		return
	}

	fmt.Fprint(w, "\n\n> ")
	w.Write(code[start:end])
	w.Write([]byte{'\n'})
}

// WriteMessage writes error message to the given writer
func (err *Error) WriteMessage(w io.Writer) {
	// Error: {msg} (at {pos})
	//   {note1}
	//   {note2}
	//   ...
	red.Fprint(w, "Error: ")
	bold.Fprint(w, err.Messages[0])
	if err.Start.File != nil {
		gray.Fprintf(w, " (at %s)", err.Start.String())
	}
	for _, msg := range err.Messages[1:] {
		green.Fprint(w, "\n  Note: ")
		fmt.Fprint(w, msg)
	}

	if err.Start.File == nil {
		return
	}
	if err.End.File == nil || err.Start.Offset == err.End.Offset {
		err.writeOnelineSnip(w)
		return
	}
	err.writeSnip(w)
}

// Error builds error message for the error.
func (err *Error) Error() string {
	var buf bytes.Buffer
	err.WriteMessage(&buf)
	return buf.String()
}

// PrintToFile prints error message to the given file. This is useful on Windows because Error()
// does not support colorful string on Windows.
func (err *Error) PrintToFile(f *os.File) {
	err.WriteMessage(colorable.NewColorable(f))
}

// Note stacks the additional message upon current error.
func (err *Error) Note(msg string) *Error {
	err.Messages = append(err.Messages, msg)
	return err
}

// Notef stacks the additional formatted message upon current error.
func (err *Error) Notef(format string, args ...interface{}) *Error {
	err.Messages = append(err.Messages, fmt.Sprintf(format, args...))
	return err
}

// NoteAt stacks the additional message upon current error with position.
func (err *Error) NoteAt(pos Pos, msg string) *Error {
	at := color.HiBlackString("(at %s)", pos.String())
	err.Messages = append(err.Messages, fmt.Sprintf("%s %s", msg, at))
	return err
}

// NotefAt stacks the additional formatted message upon current error with poisition.
func (err *Error) NotefAt(pos Pos, format string, args ...interface{}) *Error {
	return err.NoteAt(pos, fmt.Sprintf(format, args...))
}

// NewError makes locerr.Error instance without source location information.
func NewError(msg string) *Error {
	return &Error{Pos{}, Pos{}, []string{msg}}
}

// ErrorIn makes a new compilation error with the range.
func ErrorIn(start, end Pos, msg string) *Error {
	return &Error{start, end, []string{msg}}
}

// ErrorAt makes a new compilation error with the position.
func ErrorAt(pos Pos, msg string) *Error {
	return ErrorIn(pos, Pos{}, msg)
}

// Errorf makes locerr.Error instance without source location information following given format.
func Errorf(format string, args ...interface{}) *Error {
	return NewError(fmt.Sprintf(format, args...))
}

// ErrorfIn makes a new compilation error with the range and formatted message.
func ErrorfIn(start, end Pos, format string, args ...interface{}) *Error {
	return ErrorIn(start, end, fmt.Sprintf(format, args...))
}

// ErrorfAt makes a new compilation error with the position and formatted message.
func ErrorfAt(pos Pos, format string, args ...interface{}) *Error {
	return ErrorIn(pos, Pos{}, fmt.Sprintf(format, args...))
}

// WithRange adds range information to the passed error.
func WithRange(start, end Pos, err error) *Error {
	return ErrorIn(start, end, err.Error())
}

// WithPos adds positional information to the passed error.
func WithPos(pos Pos, err error) *Error {
	return ErrorAt(pos, err.Error())
}

// Note adds note to the given error. If given error is not locerr.Error, it's converted into locerr.Error.
func Note(err error, msg string) *Error {
	if err, ok := err.(*Error); ok {
		return err.Note(msg)
	}
	return &Error{Pos{}, Pos{}, []string{err.Error(), msg}}
}

// NoteIn adds range information and stack additional message to the original error. If given error is not locerr.Error, it's converted into locerr.Error.
func NoteIn(start, end Pos, err error, msg string) *Error {
	if err, ok := err.(*Error); ok {
		return err.NoteAt(start, msg)
	}
	return &Error{start, end, []string{err.Error(), msg}}
}

// NoteAt adds positional information and stack additional message to the original error. If given error is not locerr.Error, it's converted into locerr.Error.
func NoteAt(pos Pos, err error, msg string) *Error {
	return NoteIn(pos, Pos{}, err, msg)
}

// Notef adds note to the given error. Description will be created following given format and arguments. If given error is not locerr.Error, it's converted into locerr.Error.
func Notef(err error, format string, args ...interface{}) *Error {
	return Note(err, fmt.Sprintf(format, args...))
}

// NotefIn adds range information and stack additional formatted message to the original error. If given error is not locerr.Error, it's converted into locerr.Error.
func NotefIn(start, end Pos, err error, format string, args ...interface{}) *Error {
	return NoteIn(start, end, err, fmt.Sprintf(format, args...))
}

// NotefAt adds positional information and stack additional formatted message to the original error If given error is not locerr.Error, it's converted into locerr.Error.
func NotefAt(pos Pos, err error, format string, args ...interface{}) *Error {
	return NoteIn(pos, Pos{}, err, fmt.Sprintf(format, args...))
}
