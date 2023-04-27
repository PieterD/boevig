package lexer

import "fmt"

type LocationInFile struct {
	File string
	Line int
	Rune int
}

func (lif *LocationInFile) AdvanceLine() {
	lif.Line++
	lif.Rune = 1
}

func (lif *LocationInFile) AdvanceRune() {
	lif.Rune++
}

func (lif LocationInFile) String() string {
	return fmt.Sprintf("%s:%d[%d]", lif.File, lif.Line, lif.Rune)
}
