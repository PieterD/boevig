package cardinal

import "testing"

func TestBox_Contains(t *testing.T) {
	box5x5 := Box{
		TL: Coord{X: -2, Y: -2},
		BR: Coord{X: 2, Y: 2},
	}
	tests := []struct {
		desc   string
		box    Box
		at     Coord
		expect bool
	}{
		{"center", box5x5, Coord{0, 0}, true},
		{"left", box5x5, Coord{-2, 0}, true},
		{"left over", box5x5, Coord{-3, 0}, false},
		{"right", box5x5, Coord{2, 0}, true},
		{"right over", box5x5, Coord{3, 0}, false},
		{"top", box5x5, Coord{0, -2}, true},
		{"top over", box5x5, Coord{0, -3}, false},
		{"bottom", box5x5, Coord{0, 2}, true},
		{"bottom over", box5x5, Coord{0, 3}, false},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got := test.box.Contains(test.at)
			want := test.expect
			if got != want {
				t.Errorf("invalid result: got %t, want %t", got, want)
			}
		})
	}
}

func TestBox_Include(t *testing.T) {
	box5x5 := Box{
		TL: Coord{X: -2, Y: -2},
		BR: Coord{X: 2, Y: 2},
	}
	tests := []struct {
		desc    string
		initial Box
		at      Coord
		expect  Box
	}{
		{"center", box5x5, Coord{0, 0}, box5x5},
		{"left", box5x5, Coord{-2, 0}, box5x5},
		{"right", box5x5, Coord{2, 0}, box5x5},
		{"top", box5x5, Coord{0, -2}, box5x5},
		{"bottom", box5x5, Coord{0, 2}, box5x5},
		{"left over", box5x5, Coord{-3, 0}, Box{
			TL: Coord{X: -3, Y: -2},
			BR: Coord{X: 2, Y: 2},
		}},
		{"right over", box5x5, Coord{3, 0}, Box{
			TL: Coord{X: -2, Y: -2},
			BR: Coord{X: 3, Y: 2},
		}},
		{"top over", box5x5, Coord{0, -3}, Box{
			TL: Coord{X: -2, Y: -3},
			BR: Coord{X: 2, Y: 2},
		}},
		{"bottom over", box5x5, Coord{0, 3}, Box{
			TL: Coord{X: -2, Y: -2},
			BR: Coord{X: 2, Y: 3},
		}},
		{"bottom right over", box5x5, Coord{3, 3}, Box{
			TL: Coord{X: -2, Y: -2},
			BR: Coord{X: 3, Y: 3},
		}},
		{"bottom left over", box5x5, Coord{-3, 3}, Box{
			TL: Coord{X: -3, Y: -2},
			BR: Coord{X: 2, Y: 3},
		}},
		{"top right over", box5x5, Coord{3, -3}, Box{
			TL: Coord{X: -2, Y: -3},
			BR: Coord{X: 3, Y: 2},
		}},
		{"top left over", box5x5, Coord{-3, -3}, Box{
			TL: Coord{X: -3, Y: -3},
			BR: Coord{X: 2, Y: 2},
		}},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got := test.initial.Include(test.at)
			want := test.expect
			if got != want {
				t.Errorf("invalid result: got %#v, want %#v", got, want)
			}
		})
	}
}
