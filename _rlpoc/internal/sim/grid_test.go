package sim

import (
	"strings"
	"testing"
)

func TestNewGridFromSource(t *testing.T) {
	boxSource := `
grid: #####
grid: #...#
grid: #...#
grid: #...#
grid: #####
`
	boxSource = strings.ReplaceAll(boxSource, "\r\n", "\n")
	grid, err := NewGridFromSource(boxSource)
	if err != nil {
		t.Fatalf("parsing grid: %v", err)
	}
	want := strings.TrimLeft(boxSource, " \r\n\t")
	got := grid.Source()
	if want != got {
		t.Logf("want:\n%s", want)
		t.Logf("got:\n%s", got)
		t.Fatalf("input / source mismatch")
	}
}
