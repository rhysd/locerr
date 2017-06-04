package locerr

import (
	"path/filepath"
	"testing"
)

func TestStringizePos(t *testing.T) {
	src := NewDummySource("test")
	p := Pos{4, 1, 3, src}
	want := "<dummy>:1:3"
	if p.String() != want {
		t.Fatal("Unknown position format: ", p.String(), "wanted", want)
	}
}

func TestStringizeUnknownFile(t *testing.T) {
	p := Pos{}
	want := "<unknown>:0:0"
	if p.String() != want {
		t.Fatal("Unexpected position", p.String(), "wanted", want)
	}
}

func TestPosStringCanonicalPath(t *testing.T) {
	f, err := filepath.Abs("position_test.go")
	if err != nil {
		panic(err)
	}

	src, err := NewSourceFromFile("position_test.go")
	saved := currentDir
	currentDir = filepath.Dir(f)
	defer func() { currentDir = saved }()
	if err != nil {
		t.Fatal(err)
	}
	have := Pos{0, 1, 1, src}.String()
	want := "position_test.go:1:1" // Prefix was stripped
	if have != want {
		t.Fatal(want, "was wanted but have", have)
	}
}
