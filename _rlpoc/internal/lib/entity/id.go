package entity

import "fmt"

type Id struct {
	// hide the id now so we won't have to worry about it when we change it later.
	index int
}

func (id Id) String() string {
	return fmt.Sprintf("%d", id.index)
}

func (id Id) Less(than Id) bool {
	a, b := id, than
	return a.index < b.index
}

func firstId() Id {
	return Id{
		index: 1,
	}
}

func (id Id) next() Id {
	return Id{
		index: id.index + 1,
	}
}
