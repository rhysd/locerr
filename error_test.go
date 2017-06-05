package locerr

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func TestFunctionsAndMethods(t *testing.T) {
	src := NewDummySource(
		`int main() {
    foo(aaa,
        bbb,
        ccc);
    return 0;
}`,
	)

	s := Pos{21, 2, 9, src}
	e := Pos{50, 4, 11, src}

	snip := `

>     foo(aaa,
>         bbb,
>         ccc);
`
	loc := " (at <dummy>:2:9)"

	cases := []struct {
		what string
		err  *Error
		want string
	}{
		{
			what: "NewError",
			err:  NewError("This is error text"),
			want: "Error: This is error text",
		},
		{
			what: "Errorf",
			err:  Errorf("This is error text: %d", 42),
			want: "Error: This is error text: 42",
		},
		{
			what: "ErrorIn",
			err:  ErrorIn(s, e, "This is error text"),
			want: "Error: This is error text" + loc + snip,
		},
		{
			what: "ErrorfIn",
			err:  ErrorfIn(s, e, "This is error text: %d", 42),
			want: "Error: This is error text: 42" + loc + snip,
		},
		{
			what: "ErrorAt",
			err:  ErrorAt(s, "This is error text"),
			want: "Error: This is error text" + loc,
		},
		{
			what: "ErrorfAt",
			err:  ErrorfAt(s, "This is error text: %d", 42),
			want: "Error: This is error text: 42" + loc,
		},
		{
			what: "WithRange",
			err:  WithRange(s, e, fmt.Errorf("This is error text")),
			want: "Error: This is error text" + loc + snip,
		},
		{
			what: "WithPos",
			err:  WithPos(s, fmt.Errorf("This is error text")),
			want: "Error: This is error text" + loc,
		},
		{
			what: "Note to error",
			err:  Note(fmt.Errorf("This is error text"), "This is note"),
			want: "Error: This is error text\n  Note: This is note",
		},
		{
			what: "Notef to error",
			err:  Notef(fmt.Errorf("This is error text"), "This is note: %d", 42),
			want: "Error: This is error text\n  Note: This is note: 42",
		},
		{
			what: "Note to locerr.Error",
			err:  Note(ErrorIn(s, e, "This is error text"), "This is note"),
			want: "Error: This is error text" + loc + "\n  Note: This is note" + snip,
		},
		{
			what: "Notef to locerr.Error",
			err:  Notef(ErrorIn(s, e, "This is error text"), "This is note: %d", 42),
			want: "Error: This is error text" + loc + "\n  Note: This is note: 42" + snip,
		},
		{
			what: "NoteIn to error",
			err:  NoteIn(s, e, fmt.Errorf("This is error text"), "This is note"),
			want: "Error: This is error text" + loc + "\n  Note: This is note" + snip,
		},
		{
			what: "NotefIn to error",
			err:  NotefIn(s, e, fmt.Errorf("This is error text"), "This is note: %d", 42),
			want: "Error: This is error text" + loc + "\n  Note: This is note: 42" + snip,
		},
		{
			what: "NoteIn to locerr.Error",
			err:  NoteIn(s, e, ErrorIn(s, e, "This is error text"), "This is note"),
			want: "Error: This is error text" + loc + "\n  Note: This is note" + loc + snip,
		},
		{
			what: "NotefIn to locerr.Error",
			err:  NotefIn(s, e, ErrorIn(s, e, "This is error text"), "This is note: %d", 42),
			want: "Error: This is error text" + loc + "\n  Note: This is note: 42" + loc + snip,
		},
		{
			what: "NoteAt to error",
			err:  NoteAt(s, fmt.Errorf("This is error text"), "This is note"),
			want: "Error: This is error text" + loc + "\n  Note: This is note",
		},
		{
			what: "NotefAt to error",
			err:  NotefAt(s, fmt.Errorf("This is error text"), "This is note: %d", 42),
			want: "Error: This is error text" + loc + "\n  Note: This is note: 42",
		},
		{
			what: "NoteAt to locerr.Error",
			err:  NoteAt(s, ErrorIn(s, e, "This is error text"), "This is note"),
			want: "Error: This is error text" + loc + "\n  Note: This is note" + loc + snip,
		},
		{
			what: "NotefAt to locerr.Error",
			err:  NotefAt(s, ErrorIn(s, e, "This is error text"), "This is note: %d", 42),
			want: "Error: This is error text" + loc + "\n  Note: This is note: 42" + loc + snip,
		},
		{
			what: "Note method",
			err:  ErrorIn(s, e, "This is error text").Note("This is note"),
			want: "Error: This is error text" + loc + "\n  Note: This is note" + snip,
		},
		{
			what: "Notef method",
			err:  ErrorIn(s, e, "This is error text").Notef("This is note: %d", 42),
			want: "Error: This is error text" + loc + "\n  Note: This is note: 42" + snip,
		},
		{
			what: "NoteAt method",
			err:  ErrorIn(s, e, "This is error text").NoteAt(s, "This is note"),
			want: "Error: This is error text" + loc + "\n  Note: This is note" + loc + snip,
		},
		{
			what: "NotefAt method",
			err:  ErrorIn(s, e, "This is error text").NotefAt(s, "This is note: %d", 42),
			want: "Error: This is error text" + loc + "\n  Note: This is note: 42" + loc + snip,
		},
		{
			what: "nested notes",
			err:  Note(ErrorIn(s, e, "This is error text"), "This is note").NoteAt(s, "This is note second"),
			want: "Error: This is error text" + loc + "\n  Note: This is note\n  Note: This is note second" + loc + snip,
		},
	}

	for _, tc := range cases {
		t.Run(tc.what, func(t *testing.T) {
			have := tc.err.Error()
			if have != tc.want {
				t.Fatalf("Unexpected error message.\nwant:\n\n'%s'\nhave:\n\n'%s'", tc.want, have)
			}
		})
	}
}

func TestCodeSnippet(t *testing.T) {
	cases := []struct {
		what string
		code string
		from int
		to   int
		want []string
	}{
		{
			what: "whole in a line",
			code: "abc",
			from: 0,
			to:   2,
			want: []string{
				"> abc",
			},
		},
		{
			what: "slice in a line",
			code: "abc",
			from: 1,
			to:   2,
			want: []string{
				"> abc",
			},
		},
		{
			what: "slice in a line with indent",
			code: "	 abc",
			from: 3,
			to:   4,
			want: []string{
				"> 	 abc",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.what, func(t *testing.T) {
			src := NewDummySource(tc.code)
			calcPos := func(offset int) Pos {
				code := src.Code
				o, l, c, end := 0, 1, 1, len(code)
				for o != end {
					if o == offset {
						return Pos{o, l, c, src}
					}
					if code[o] == '\n' {
						l++
						c = 1
					} else {
						c++
					}
					o++
				}
				if o != offset {
					t.Fatal("Offsetis illegal")
				}
				return Pos{o, l, c, src}
			}
			err := ErrorIn(calcPos(tc.from), calcPos(tc.to), "text")
			have := strings.SplitN(err.Error(), "\n", 3)[2]
			want := strings.Join(tc.want, "\n") + "\n"
			if have != want {
				t.Fatalf("Unexpected snippet\n\nwant:\n'%s'\nhave:\n'%s'", want, have)
			}
		})
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

func TestSnipIsEmpty(t *testing.T) {
	s := NewDummySource("abc")
	p := Pos{1, 1, 2, s}
	err := ErrorIn(p, p, "This is error text")
	want := "Error: This is error text (at <dummy>:1:2)"
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
	if color.NoColor {
		t.Fatal("Color should be enabled")
	}
	SetColor(false)
	if !color.NoColor {
		t.Fatal("Color should be disabled (2)")
	}
}
