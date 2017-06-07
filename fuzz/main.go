package locerrfuzzing

import (
	"github.com/rhysd/locerr"
	"math/rand"
	"strings"
)

func offsetPos(src *locerr.Source, offset int) locerr.Pos {
	o, l, c, end := 0, 1, 1, len(src.Code)
	for o != end {
		if o == offset {
			return locerr.Pos{o, l, c, src}
		}
		if src.Code[o] == '\n' {
			l++
			c = 1
		} else {
			c++
		}
		o++
	}
	return locerr.Pos{o, l, c, src}
}

// Fuzz do fuzzing test using go-fuzz
func Fuzz(data []byte) int {
	src := locerr.NewDummySource(string(data))
	len := len(data)
	if len == 0 {
		p := locerr.Pos{0, 1, 1, src}
		return fuzz(src, p, p)
	}

	o := 0
	if len > 1 {
		o = rand.Intn(len - 1)
	}
	s := offsetPos(src, o)
	o = rand.Intn(len-o) + o
	e := offsetPos(src, o)
	return fuzz(src, s, e)
}

func fuzz(src *locerr.Source, start locerr.Pos, end locerr.Pos) int {
	err := locerr.ErrorIn(start, end, "fuzz")
	err = err.NoteAt(start, "note1")
	err = err.Note("note2")
	msg := err.Error()
	if !strings.Contains(msg, "fuzz") || !strings.Contains(msg, "note1") || !strings.Contains(msg, "note2") {
		panic("Unexpected error message: " + msg)
	}
	return 1 // data is good for fuzzing
}
