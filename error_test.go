package loc

import (
	"fmt"
	"github.com/fatih/color"
	"runtime"
	"testing"
)

func testMakeRange() (Pos, Pos) {
	s := NewDummySource(
		`package prelude

import (
	"testing"
)`,
	)

	start := Pos{4, 1, 4, s}
	end := Pos{20, 3, 3, s}
	return start, end
}

func testMakePos() Pos {
	p, _ := testMakeRange()
	return p
}

func TestNewError(t *testing.T) {
	want :=
		`Error: This is error text (at <dummy>:1:4)

> age prelude
> 
> imp

`

	s, e := testMakeRange()
	errs := []*Error{
		ErrorIn(s, e, "This is error text"),
		ErrorfIn(s, e, "This %s error %s", "is", "text"),
	}
	for _, err := range errs {
		got := err.Error()
		if got != want {
			t.Fatalf("Unexpected error message. want: '%s', got: '%s'", want, got)
		}
	}
}

func TestNewErrorAt(t *testing.T) {
	want := "Error: This is error text (at <dummy>:1:4)"
	for _, err := range []*Error{
		ErrorAt(testMakePos(), "This is error text"),
		ErrorfAt(testMakePos(), "This is %s text", "error"),
	} {
		got := err.Error()
		if got != want {
			t.Fatalf("Unexpected error message. want: '%s', got: '%s'", want, got)
		}
	}
}

func TestWithRange(t *testing.T) {
	want :=
		`Error: This is an error text (at <dummy>:1:4)

> age prelude
> 
> imp

`

	s, e := testMakeRange()
	err := WithRange(s, e, fmt.Errorf("This is an error text"))
	got := err.Error()
	if got != want {
		t.Fatalf("Unexpected error message. want: '%s', got: '%s'", want, got)
	}
}

func TestWithPos(t *testing.T) {
	want := "Error: This is wrapped error text (at <dummy>:1:4)"
	got := WithPos(testMakePos(), fmt.Errorf("This is wrapped error text")).Error()
	if got != want {
		t.Fatalf("Unexpected error message. want: '%s', got: '%s'", want, got)
	}
}

func TestNote(t *testing.T) {
	want :=
		`Error: This is original error text
  Note: This is additional error text`

	errs := []*Error{
		Note(fmt.Errorf("This is original error text"), "This is additional error text"),
		Notef(fmt.Errorf("This is original error text"), "This is %s error text", "additional"),
	}
	for _, err := range errs {
		got := err.Error()
		if got != want {
			t.Fatalf("Unexpected error message. want: '%s', got: '%s'", want, got)
		}
	}

	want =
		`Error: This is original error text (at <dummy>:1:4)
  Note: This is additional error text

> age prelude
> 
> imp

`
	s, e := testMakeRange()
	err := Note(ErrorIn(s, e, "This is original error text"), "This is additional error text")
	got := err.Error()
	if got != want {
		t.Fatalf("Unexpected error message. want: '%s', got: '%s'", want, got)
	}
}

func TestNoteIn(t *testing.T) {
	want :=
		`Error: This is original error text (at <dummy>:1:4)
  Note: This is additional error text

> age prelude
> 
> imp

`

	s, e := testMakeRange()
	errs := []*Error{
		NoteIn(s, e, fmt.Errorf("This is original error text"), "This is additional error text"),
		NotefIn(s, e, fmt.Errorf("This is original error text"), "This is %s error text", "additional"),
	}
	for _, err := range errs {
		got := err.Error()
		if got != want {
			t.Fatalf("Unexpected error message. want: '%s', got: '%s'", want, got)
		}
	}

	want =
		`Error: This is original error text (at <dummy>:1:4)
  Note: This is additional error text (at <dummy>:1:4)

> age prelude
> 
> imp

`
	s, e = testMakeRange()
	err := NoteIn(s, e, ErrorIn(s, e, "This is original error text"), "This is additional error text")
	got := err.Error()
	if got != want {
		t.Fatalf("Unexpected error message. want: '%s', got: '%s'", want, got)
	}
}

func TestNoteAt(t *testing.T) {
	want := "Error: This is original error text (at <dummy>:1:4)\n  Note: This is additional error text"
	pos := testMakePos()
	original := fmt.Errorf("This is original error text")
	for _, err := range []*Error{
		NoteAt(pos, original, "This is additional error text"),
		NotefAt(pos, original, "This is additional %s", "error text"),
	} {
		got := err.Error()
		if got != want {
			t.Fatalf("Unexpected error message. want: '%s', got: '%s'", want, got)
		}
	}
}

func TestNoteMethods(t *testing.T) {
	want :=
		`Error: This is original error text (at <dummy>:1:4)
  Note: This is additional error text

> age prelude
> 
> imp

`

	s, e := testMakeRange()
	errs := []*Error{
		ErrorIn(s, e, "This is original error text").Note("This is additional error text"),
		ErrorIn(s, e, "This is original error text").Notef("This is %s", "additional error text"),
	}
	for _, err := range errs {
		got := err.Error()
		if got != want {
			t.Fatalf("Unexpected error message. want: '%s', got: '%s'", want, got)
		}
	}
}

func TestNoteMethodsWithPos(t *testing.T) {
	want :=
		`Error: This is original error text (at <dummy>:1:4)
  Note: This is additional error text (at <dummy>:1:4)

> age prelude
> 
> imp

`

	s, e := testMakeRange()

	errs := []*Error{
		ErrorIn(s, e, "This is original error text").NoteAt(s, "This is additional error text"),
		ErrorIn(s, e, "This is original error text").NotefAt(s, "This is %s", "additional error text"),
	}
	for _, err := range errs {
		got := err.Error()
		if got != want {
			t.Fatalf("Unexpected error message. want: '%s', got: '%s'", want, got)
		}
	}
}

func TestCodeIsEmpty(t *testing.T) {
	s := NewDummySource("")
	p := Pos{0, 1, 1, s}
	err := ErrorIn(p, p, "This is error text")
	want := "Error: This is error text (at <dummy>:1:1)"
	got := err.Error()

	if want != got {
		t.Fatalf("Unexpected error message. want: '%s', got: '%s'", want, got)
	}
}

func TestSetColor(t *testing.T) {
	defer func() { SetColor(true) }()
	SetColor(false)
	if !color.NoColor {
		t.Fatal("Color should be disabled")
	}
	SetColor(true)
	if runtime.GOOS == "windows" {
		if !color.NoColor {
			t.Fatal("Color should be disabled (2)")
		}
	} else {
		if color.NoColor {
			t.Fatal("Color should be enabled")
		}
	}
	SetColor(false)
	if !color.NoColor {
		t.Fatal("Color should be disabled (2)")
	}
}
