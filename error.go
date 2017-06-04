package loc

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"runtime"
	"strings"
)

func init() {
	// FIXME:
	// On Windows, ANSI sequences are not available on cmd.exe.
	// go-colorable cannot help to solve this because it only provides writers to stdout or
	// stderr.
	// At least I need to improve the check. Even if running on windows, WLS would be able to
	// handle ANSI sequences.
	if runtime.GOOS == "windows" {
		color.NoColor = true
	}
}

// SetColor controls font should be colorful or not.
func SetColor(enabled bool) {
	color.NoColor = !enabled || runtime.GOOS == "windows"
}

var (
	bold = color.New(color.Bold).SprintFunc()
)

// stringWriter is an interface to use WriteString method for both *bytes.Buffer and *os.File.
// err.Error uses Buffer to build a string and err.Fprint
type stringWriter interface {
	WriteString(s string) (n int, err error)
}

// Error represents a compilation error with positional information and stacked messages.
type Error struct {
	Start    Pos
	End      Pos
	Messages []string
}

func (err *Error) snippetLines() []string {
	if err.End.File == nil {
		return nil
	}

	code := err.Start.File.Code
	start, end := err.Start.Offset, err.End.Offset
	if start >= end {
		return nil
	}

	// Fix first line and last line. Code snippet should be line-wise.
	for start-1 >= 0 {
		if code[start-1] == '\n' {
			break
		}
		start--
	}
	for end < len(code) {
		if code[end] == '\n' {
			break
		}
		end++
	}

	return strings.Split(string(code[start:end]), "\n")
}

func (err *Error) WriteMessage(w stringWriter) {
	s := err.Start

	// Error: {msg} (at {pos})
	//   {note1}
	//   {note2}
	//   ...
	w.WriteString(color.RedString("Error: "))
	w.WriteString(bold(err.Messages[0]))
	if s.File != nil {
		w.WriteString(color.HiBlackString(" (at %s)", s.String()))
	}
	for _, msg := range err.Messages[1:] {
		w.WriteString(color.GreenString("\n  Note: "))
		w.WriteString(msg)
	}

	lines := err.snippetLines()
	if len(lines) == 0 {
		return
	}

	w.WriteString("\n\n")
	for _, line := range lines {
		w.WriteString("> ")
		w.WriteString(line)
		w.WriteString("\n")
	}
	w.WriteString("\n")

	// TODO:
	// If the code snippet for the token is too long, skip lines with '...' except for starting N lines
	// and ending N lines
}

// Error builds error message for the error.
func (err *Error) Error() string {
	var buf bytes.Buffer
	err.WriteMessage(&buf)
	return buf.String()
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

// NewError makes loc.Error instance without source location information.
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

// NewErrorf makes loc.Error instance without source location information following given format.
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

// Note adds note to the given error. If given error is not loc.Error, it's converted into loc.Error.
func Note(err error, msg string) *Error {
	if err, ok := err.(*Error); ok {
		return err.Note(msg)
	}
	return &Error{Pos{}, Pos{}, []string{err.Error(), msg}}
}

// NoteIn adds range information and stack additional message to the original error. If given error is not loc.Error, it's converted into loc.Error.
func NoteIn(start, end Pos, err error, msg string) *Error {
	if err, ok := err.(*Error); ok {
		return err.NoteAt(start, msg)
	}
	return &Error{start, end, []string{err.Error(), msg}}
}

// NoteAt adds positional information and stack additional message to the original error. If given error is not loc.Error, it's converted into loc.Error.
func NoteAt(pos Pos, err error, msg string) *Error {
	return NoteIn(pos, Pos{}, err, msg)
}

// Notef adds note to the given error. Description will be created following given format and arguments. If given error is not loc.Error, it's converted into loc.Error.
func Notef(err error, format string, args ...interface{}) *Error {
	return Note(err, fmt.Sprintf(format, args...))
}

// NotefIn adds range information and stack additional formatted message to the original error. If given error is not loc.Error, it's converted into loc.Error.
func NotefIn(start, end Pos, err error, format string, args ...interface{}) *Error {
	return NoteIn(start, end, err, fmt.Sprintf(format, args...))
}

// NotefAt adds positional information and stack additional formatted message to the original error If given error is not loc.Error, it's converted into loc.Error..
func NotefAt(pos Pos, err error, format string, args ...interface{}) *Error {
	return NoteIn(pos, Pos{}, err, fmt.Sprintf(format, args...))
}
