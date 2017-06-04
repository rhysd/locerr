package locerr

import (
	"fmt"
	"os"
	"path/filepath"
)

var currentDir string

func init() {
	currentDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
}

// Pos represents some point in a source code.
type Pos struct {
	// Offset from the beginning of code.
	Offset int
	// Line number.
	Line int
	// Column number.
	Column int
	// File of this position.
	File *Source
}

// String makes a string representation of the position. Format is 'file:line:column'.
func (p Pos) String() string {
	if p.File == nil {
		return "<unknown>:0:0"
	}
	f := p.File.Path
	if p.File.Exists && currentDir != "" && filepath.HasPrefix(f, currentDir) {
		f, _ = filepath.Rel(currentDir, f)
	}
	return fmt.Sprintf("%s:%d:%d", f, p.Line, p.Column)
}
